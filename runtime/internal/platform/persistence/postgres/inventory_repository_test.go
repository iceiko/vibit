package postgres

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestInventoryRepositoryGetInventoryMapsRowsInDatabaseOrder(t *testing.T) {
	executor := &recordingExecutor{
		rowsResponses: []pgx.Rows{
			&fakeRows{
				values: [][]any{
					{"item-1", int64(2)},
					{"item-2", int64(5)},
				},
			},
		},
	}
	repository := NewInventoryRepositoryForUnitOfWork(executor)

	items, err := repository.GetInventory(context.Background(), " player-1 ")
	if err != nil {
		t.Fatalf("GetInventory() error = %v, want nil", err)
	}

	want := []inventory.Item{
		{ItemID: "item-1", Quantity: 2},
		{ItemID: "item-2", Quantity: 5},
	}
	if !reflect.DeepEqual(items, want) {
		t.Fatalf("GetInventory() items = %#v, want %#v", items, want)
	}

	if len(executor.queries) != 1 {
		t.Fatalf("queries len = %d, want 1", len(executor.queries))
	}
	assertSQLContains(t, executor.queries[0].sql, "FROM inventory_items")
	assertSQLContains(t, executor.queries[0].sql, "WHERE player_id = $1")
	assertSQLContains(t, executor.queries[0].sql, "ORDER BY item_id")
	assertArgs(t, executor.queries[0].args, "player-1")
}

func TestInventoryRepositoryLockInventoryForMutationUsesAccountRowLock(t *testing.T) {
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			fakeRow{values: []any{"player-1"}},
		},
	}
	repository := NewInventoryRepositoryForUnitOfWork(executor)

	lock, err := repository.LockInventoryForMutation(context.Background(), " player-1 ")
	if err != nil {
		t.Fatalf("LockInventoryForMutation() error = %v, want nil", err)
	}
	defer lock.Release()

	if len(executor.execs) != 1 {
		t.Fatalf("execs len = %d, want 1", len(executor.execs))
	}
	assertSQLContains(t, executor.execs[0].sql, "INSERT INTO inventory_accounts")
	assertSQLContains(t, executor.execs[0].sql, "ON CONFLICT (player_id) DO NOTHING")
	assertArgs(t, executor.execs[0].args, "player-1")

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	assertSQLContains(t, executor.queryRowCalls[0].sql, "FROM inventory_accounts")
	assertSQLContains(t, executor.queryRowCalls[0].sql, "FOR UPDATE")
	assertArgs(t, executor.queryRowCalls[0].args, "player-1")
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}

	_, err = lock.GetInventory(context.Background(), "other-player")
	if err == nil {
		t.Fatal("locked GetInventory() error = nil, want player mismatch")
	}
}

func TestInventoryMutationLockGrantItemUpsertsQuantityAndRecordsGrant(t *testing.T) {
	occurredAt := time.Date(2026, 5, 13, 10, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			fakeRow{values: []any{"player-1"}},
			fakeRow{values: []any{"item-1", int64(8)}},
		},
	}
	repository := NewInventoryRepositoryForUnitOfWork(executor)
	lock, err := repository.LockInventoryForMutation(context.Background(), "player-1")
	if err != nil {
		t.Fatalf("LockInventoryForMutation() error = %v, want nil", err)
	}

	granted, err := lock.GrantItem(context.Background(), inventory.GrantItemMutation{
		EventID:    " event-1 ",
		OccurredAt: occurredAt,
		PlayerID:   " player-1 ",
		ItemID:     " item-1 ",
		Quantity:   3,
		Reason:     " reward ",
	})
	if err != nil {
		t.Fatalf("GrantItem() error = %v, want nil", err)
	}
	if granted != (inventory.Item{ItemID: "item-1", Quantity: 8}) {
		t.Fatalf("GrantItem() item = %#v, want item-1 quantity 8", granted)
	}

	if len(executor.queryRowCalls) != 2 {
		t.Fatalf("query rows len = %d, want 2", len(executor.queryRowCalls))
	}
	grantQuery := executor.queryRowCalls[1]
	assertSQLContains(t, grantQuery.sql, "INSERT INTO inventory_items")
	assertSQLContains(t, grantQuery.sql, "ON CONFLICT (player_id, item_id)")
	assertSQLContains(t, grantQuery.sql, "quantity = inventory_items.quantity + EXCLUDED.quantity")
	assertSQLContains(t, grantQuery.sql, "RETURNING item_id, quantity")
	assertArgs(t, grantQuery.args, "player-1", "item-1", int64(3))

	if len(executor.execs) != 2 {
		t.Fatalf("execs len = %d, want 2", len(executor.execs))
	}
	grantRecordExec := executor.execs[1]
	assertSQLContains(t, grantRecordExec.sql, "INSERT INTO inventory_item_grants")
	assertArgs(
		t,
		grantRecordExec.args,
		"event-1",
		occurredAt.UTC(),
		"player-1",
		"item-1",
		int64(3),
		int64(8),
		"reward",
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestInventoryMutationLockRejectsMissingGrantRecordFields(t *testing.T) {
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			fakeRow{values: []any{"player-1"}},
		},
	}
	repository := NewInventoryRepositoryForUnitOfWork(executor)
	lock, err := repository.LockInventoryForMutation(context.Background(), "player-1")
	if err != nil {
		t.Fatalf("LockInventoryForMutation() error = %v, want nil", err)
	}

	_, err = lock.GrantItem(context.Background(), inventory.GrantItemMutation{
		PlayerID:   "player-1",
		ItemID:     "item-1",
		Quantity:   1,
		Reason:     "reward",
		OccurredAt: time.Now().UTC(),
	})
	if err == nil {
		t.Fatal("GrantItem() error = nil, want missing event_id error")
	}
	if !strings.Contains(err.Error(), "event_id") {
		t.Fatalf("GrantItem() error = %v, want event_id error", err)
	}
	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want only account lock query before validation failure", len(executor.queryRowCalls))
	}
}

func TestInventoryMutationLockReleaseDoesNotCommitOrRollback(t *testing.T) {
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			fakeRow{values: []any{"player-1"}},
		},
	}
	repository := NewInventoryRepositoryForUnitOfWork(executor)
	lock, err := repository.LockInventoryForMutation(context.Background(), "player-1")
	if err != nil {
		t.Fatalf("LockInventoryForMutation() error = %v, want nil", err)
	}

	lock.Release()
	lock.Release()

	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
	_, err = lock.GetInventory(context.Background(), "player-1")
	if err == nil {
		t.Fatal("GetInventory() after Release error = nil, want released lock error")
	}
}

func TestInventoryRepositoryRequiresUnitOfWorkExecutor(t *testing.T) {
	repository := NewInventoryRepositoryForUnitOfWork(nil)

	_, err := repository.GetInventory(context.Background(), "player-1")
	if err == nil {
		t.Fatal("GetInventory() error = nil, want executor error")
	}

	_, err = repository.LockInventoryForMutation(context.Background(), "player-1")
	if err == nil {
		t.Fatal("LockInventoryForMutation() error = nil, want executor error")
	}
}

type recordedCall struct {
	sql  string
	args []any
}

type recordingExecutor struct {
	execs          []recordedCall
	queries        []recordedCall
	queryRowCalls  []recordedCall
	rowsErr        error
	rowsResponses  []pgx.Rows
	rowsIdx        int
	rowResponses   []pgx.Row
	rowIdx         int
	execErr        error
	queryErr       error
	execCommandTag string
}

func (e *recordingExecutor) Exec(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	e.execs = append(e.execs, recordedCall{sql: sql, args: append([]any(nil), args...)})
	if e.execErr != nil {
		return pgconn.CommandTag{}, e.execErr
	}
	if e.execCommandTag != "" {
		return pgconn.NewCommandTag(e.execCommandTag), nil
	}
	return pgconn.NewCommandTag("INSERT 1"), nil
}

func (e *recordingExecutor) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
	e.queries = append(e.queries, recordedCall{sql: sql, args: append([]any(nil), args...)})
	if e.queryErr != nil {
		return nil, e.queryErr
	}
	if e.rowsIdx >= len(e.rowsResponses) {
		return &fakeRows{err: e.rowsErr}, nil
	}
	rows := e.rowsResponses[e.rowsIdx]
	e.rowsIdx += 1
	return rows, nil
}

func (e *recordingExecutor) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	e.queryRowCalls = append(e.queryRowCalls, recordedCall{sql: sql, args: append([]any(nil), args...)})
	if e.rowIdx >= len(e.rowResponses) {
		return fakeRow{err: pgx.ErrNoRows}
	}
	row := e.rowResponses[e.rowIdx]
	e.rowIdx += 1
	return row
}

func (e *recordingExecutor) allSQL() []string {
	sql := make([]string, 0, len(e.execs)+len(e.queries)+len(e.queryRowCalls))
	for _, call := range e.execs {
		sql = append(sql, call.sql)
	}
	for _, call := range e.queries {
		sql = append(sql, call.sql)
	}
	for _, call := range e.queryRowCalls {
		sql = append(sql, call.sql)
	}
	return sql
}

type fakeRow struct {
	values []any
	err    error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("fake row: destination count mismatch")
	}
	for i := range dest {
		assignValue(dest[i], r.values[i])
	}
	return nil
}

type fakeRows struct {
	values [][]any
	index  int
	err    error
	closed bool
}

func (r *fakeRows) Close() {
	r.closed = true
}

func (r *fakeRows) Err() error {
	return r.err
}

func (r *fakeRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}

func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *fakeRows) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index += 1
	return true
}

func (r *fakeRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.values) {
		return errors.New("fake rows: scan without current row")
	}
	row := r.values[r.index-1]
	if len(dest) != len(row) {
		return errors.New("fake rows: destination count mismatch")
	}
	for i := range dest {
		assignValue(dest[i], row[i])
	}
	return nil
}

func (r *fakeRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.values) {
		return nil, errors.New("fake rows: values without current row")
	}
	return append([]any(nil), r.values[r.index-1]...), nil
}

func (r *fakeRows) RawValues() [][]byte {
	return nil
}

func (r *fakeRows) Conn() *pgx.Conn {
	return nil
}

func assignValue(dest any, value any) {
	switch pointer := dest.(type) {
	case *string:
		*pointer = value.(string)
	case *int64:
		*pointer = value.(int64)
	case *time.Time:
		*pointer = value.(time.Time)
	default:
		panic("fake rows: unsupported destination type")
	}
}

func assertSQLContains(t *testing.T, sql string, want string) {
	t.Helper()
	if !strings.Contains(sql, want) {
		t.Fatalf("sql %q does not contain %q", sql, want)
	}
}

func assertArgs(t *testing.T, got []any, want ...any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args = %#v, want %#v", got, want)
	}
}

func hasTransactionControlSQL(sqlStatements []string) bool {
	for _, sql := range sqlStatements {
		sql = strings.ToUpper(strings.TrimSpace(sql))
		if strings.HasPrefix(sql, "BEGIN") || strings.HasPrefix(sql, "COMMIT") || strings.HasPrefix(sql, "ROLLBACK") {
			return true
		}
	}
	return false
}
