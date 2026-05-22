package postgres

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestStorageObjectRepositoryCreateInsertsActiveObject(t *testing.T) {
	createdAt := time.Date(2026, 5, 22, 9, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{values: activeStorageObjectRowValues(createdAt.UTC())},
		},
	}
	repository := NewStorageObjectRepositoryForUnitOfWork(executor)

	value := []byte(`{"level":3}`)
	record, err := repository.CreateStorageObject(context.Background(), storage.CreateStorageObjectInput{
		ObjectID:    " object-1 ",
		Owner:       storage.StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Identity:    storage.StorageObjectIdentity{Collection: " progress ", Key: " tutorial "},
		Value:       storage.StorageObjectValue{JSON: value},
		RequestedBy: " storage_service ",
	})
	if err != nil {
		t.Fatalf("CreateStorageObject() error = %v, want nil", err)
	}
	if record.ObjectID != "object-1" ||
		record.Owner.Kind != storage.OwnerKindPlayer ||
		record.Owner.ID != "player-1" ||
		record.Identity.Collection != "progress" ||
		record.Identity.Key != "tutorial" ||
		record.Version != storage.InitialStorageObjectVersion ||
		record.Status != storage.StorageObjectStatusActive {
		t.Fatalf("CreateStorageObject() record = %#v, want active object", record)
	}
	if record.CreatedAt.Location() != time.UTC || record.UpdatedAt.Location() != time.UTC {
		t.Fatalf("CreateStorageObject() timestamps = %#v, want UTC", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "INSERT INTO storage_objects")
	assertSQLContains(t, call.sql, "RETURNING")
	assertArgs(t,
		call.args,
		"object-1",
		"player",
		"player-1",
		"progress",
		"tutorial",
		[]byte(`{"level":3}`),
		int64(1),
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}

	value[0] = '['
	if string(record.Value.JSON) != `{"level":3}` {
		t.Fatalf("record.Value.JSON = %q, want copied value", string(record.Value.JSON))
	}
}

func TestStorageObjectRepositoryGetSelectsActiveObjectByOwnerIdentity(t *testing.T) {
	createdAt := time.Date(2026, 5, 22, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{values: activeStorageObjectRowValues(createdAt)},
		},
	}
	repository := NewStorageObjectRepositoryForUnitOfWork(executor)

	record, err := repository.GetStorageObject(context.Background(), storage.GetStorageObjectInput{
		Owner:    storage.StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Identity: storage.StorageObjectIdentity{Collection: " progress ", Key: " tutorial "},
	})
	if err != nil {
		t.Fatalf("GetStorageObject() error = %v, want nil", err)
	}
	if record.ObjectID != "object-1" || record.Status != storage.StorageObjectStatusActive {
		t.Fatalf("GetStorageObject() record = %#v, want active object-1", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "FROM storage_objects")
	assertSQLContains(t, call.sql, "owner_kind = $1")
	assertSQLContains(t, call.sql, "owner_id = $2")
	assertSQLContains(t, call.sql, "collection = $3")
	assertSQLContains(t, call.sql, "object_key = $4")
	assertSQLContains(t, call.sql, "deleted_at IS NULL")
	assertArgs(t, call.args, "player", "player-1", "progress", "tutorial")
}

func TestStorageObjectRepositoryListIsOwnerCollectionScopedAndOrdered(t *testing.T) {
	createdAt := time.Date(2026, 5, 22, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowsResponses: []pgx.Rows{
			&storageObjectRows{
				values: [][]any{
					activeStorageObjectRowValues(createdAt, withStorageObjectKey("alpha"), withStorageObjectValue([]byte(`{"slot":1}`))),
					activeStorageObjectRowValues(createdAt, withStorageObjectID("object-2"), withStorageObjectKey("bravo"), withStorageObjectValue([]byte(`{"slot":2}`))),
					activeStorageObjectRowValues(createdAt, withStorageObjectID("object-3"), withStorageObjectKey("charlie"), withStorageObjectValue([]byte(`{"slot":3}`))),
				},
			},
		},
	}
	repository := NewStorageObjectRepositoryForUnitOfWork(executor)

	result, err := repository.ListStorageObjects(context.Background(), storage.ListStorageObjectsInput{
		Owner:          storage.StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Collection:     " progress ",
		Limit:          2,
		AfterObjectKey: " alpha ",
	})
	if err != nil {
		t.Fatalf("ListStorageObjects() error = %v, want nil", err)
	}
	if len(result.Objects) != 2 || result.Objects[0].Identity.Key != "alpha" || result.Objects[1].Identity.Key != "bravo" {
		t.Fatalf("ListStorageObjects() objects = %#v, want first two ordered rows", result.Objects)
	}
	if result.NextObjectKey != "charlie" {
		t.Fatalf("NextObjectKey = %q, want charlie overflow cursor", result.NextObjectKey)
	}

	if len(executor.queries) != 1 {
		t.Fatalf("queries len = %d, want 1", len(executor.queries))
	}
	call := executor.queries[0]
	assertSQLContains(t, call.sql, "FROM storage_objects")
	assertSQLContains(t, call.sql, "owner_kind = $1")
	assertSQLContains(t, call.sql, "owner_id = $2")
	assertSQLContains(t, call.sql, "collection = $3")
	assertSQLContains(t, call.sql, "object_key > $4")
	assertSQLContains(t, call.sql, "ORDER BY object_key")
	assertSQLContains(t, call.sql, "LIMIT $5")
	assertArgs(t, call.args, "player", "player-1", "progress", "alpha", int32(3))
}

func TestStorageObjectRepositoryUpdateChecksExpectedVersionAndIncrementsVersion(t *testing.T) {
	updatedAt := time.Date(2026, 5, 22, 2, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{values: activeStorageObjectRowValues(updatedAt.UTC(), withStorageObjectVersion(4), withStorageObjectValue([]byte(`{"level":4}`)))},
		},
	}
	repository := NewStorageObjectRepositoryForUnitOfWork(executor)
	expected := storage.StorageObjectVersion(3)

	record, err := repository.UpdateStorageObject(context.Background(), storage.UpdateStorageObjectInput{
		Owner:           storage.StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Identity:        storage.StorageObjectIdentity{Collection: " progress ", Key: " tutorial "},
		Value:           storage.StorageObjectValue{JSON: []byte(`{"level":4}`)},
		ExpectedVersion: &expected,
		RequestedBy:     " storage_service ",
	})
	if err != nil {
		t.Fatalf("UpdateStorageObject() error = %v, want nil", err)
	}
	if record.Version != 4 || string(record.Value.JSON) != `{"level":4}` {
		t.Fatalf("UpdateStorageObject() record = %#v, want incremented version and replacement value", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "UPDATE storage_objects")
	assertSQLContains(t, call.sql, "value_json = $5")
	assertSQLContains(t, call.sql, "version = version + 1")
	assertSQLContains(t, call.sql, "AND version = $6")
	assertSQLContains(t, call.sql, "deleted_at IS NULL")
	assertArgs(t, call.args, "player", "player-1", "progress", "tutorial", []byte(`{"level":4}`), int64(3))
}

func TestStorageObjectRepositoryDeleteSoftDeletesAndChecksExpectedVersion(t *testing.T) {
	deletedAt := time.Date(2026, 5, 22, 2, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{values: activeStorageObjectRowValues(deletedAt.UTC(), withStorageObjectVersion(5), withStorageObjectDeletedAt(deletedAt.UTC()))},
		},
	}
	repository := NewStorageObjectRepositoryForUnitOfWork(executor)
	expected := storage.StorageObjectVersion(4)

	record, err := repository.DeleteStorageObject(context.Background(), storage.DeleteStorageObjectInput{
		Owner:           storage.StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Identity:        storage.StorageObjectIdentity{Collection: " progress ", Key: " tutorial "},
		ExpectedVersion: &expected,
		RequestedBy:     " storage_service ",
	})
	if err != nil {
		t.Fatalf("DeleteStorageObject() error = %v, want nil", err)
	}
	if record.Status != storage.StorageObjectStatusDeleted || record.DeletedAt == nil {
		t.Fatalf("DeleteStorageObject() record = %#v, want soft-deleted object", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "UPDATE storage_objects")
	assertSQLContains(t, call.sql, "deleted_at = now()")
	assertSQLContains(t, call.sql, "version = version + 1")
	assertSQLContains(t, call.sql, "AND version = $5")
	assertSQLContains(t, call.sql, "deleted_at IS NULL")
	assertArgs(t, call.args, "player", "player-1", "progress", "tutorial", int64(4))
}

func TestStorageObjectRepositoryMapsErrorsAndRedactsDetails(t *testing.T) {
	repository := NewStorageObjectRepositoryForUnitOfWork(&recordingExecutor{})

	_, err := repository.GetStorageObject(context.Background(), storage.GetStorageObjectInput{
		Owner:    storage.StorageObjectOwner{Kind: storage.OwnerKindPlayer, ID: "player-1"},
		Identity: storage.StorageObjectIdentity{Collection: "progress", Key: "tutorial"},
	})
	assertStorageConflictClass(t, err, storage.StorageObjectConflictNotFound)

	duplicate := &pgconn.PgError{Code: "23505", ConstraintName: "storage_objects_active_identity_uq", Detail: `value_json {"secret":true}`}
	repository = NewStorageObjectRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{err: duplicate},
		},
	})
	_, err = repository.CreateStorageObject(context.Background(), validCreateStorageObjectInput())
	assertStorageConflictClass(t, err, storage.StorageObjectConflictAlreadyExists)
	assertStorageObjectErrorRedacted(t, err)

	constraint := &pgconn.PgError{Code: "23503", ConstraintName: "storage_objects_owner_player_fk", Detail: "player-1"}
	repository = NewStorageObjectRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{err: constraint},
		},
	})
	_, err = repository.CreateStorageObject(context.Background(), validCreateStorageObjectInput())
	assertStorageConflictClass(t, err, storage.StorageObjectConflictStorageUnavailable)
	assertStorageObjectErrorRedacted(t, err)

	expected := storage.StorageObjectVersion(7)
	repository = NewStorageObjectRepositoryForUnitOfWork(&recordingExecutor{})
	_, err = repository.UpdateStorageObject(context.Background(), storage.UpdateStorageObjectInput{
		Owner:           storage.StorageObjectOwner{Kind: storage.OwnerKindPlayer, ID: "player-1"},
		Identity:        storage.StorageObjectIdentity{Collection: "progress", Key: "tutorial"},
		Value:           storage.StorageObjectValue{JSON: []byte(`{"level":4}`)},
		ExpectedVersion: &expected,
		RequestedBy:     "storage_service",
	})
	assertStorageConflictClass(t, err, storage.StorageObjectConflictVersionMismatch)
}

func TestStorageObjectRepositoryRejectsInvalidRows(t *testing.T) {
	values := activeStorageObjectRowValues(time.Date(2026, 5, 22, 1, 2, 3, 0, time.UTC))
	values[5] = []byte(`[]`)
	repository := NewStorageObjectRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{values: values},
		},
	})

	_, err := repository.GetStorageObject(context.Background(), storage.GetStorageObjectInput{
		Owner:    storage.StorageObjectOwner{Kind: storage.OwnerKindPlayer, ID: "player-1"},
		Identity: storage.StorageObjectIdentity{Collection: "progress", Key: "tutorial"},
	})
	assertStorageConflictClass(t, err, storage.StorageObjectConflictStorageUnavailable)
	assertStorageObjectErrorRedacted(t, err)
}

func TestStorageObjectRepositoryRequiresUnitOfWorkExecutor(t *testing.T) {
	repository := NewStorageObjectRepositoryForUnitOfWork(nil)

	_, err := repository.CreateStorageObject(context.Background(), validCreateStorageObjectInput())
	if err == nil {
		t.Fatal("CreateStorageObject() error = nil, want executor error")
	}

	_, err = repository.ListStorageObjects(context.Background(), storage.ListStorageObjectsInput{
		Owner:      storage.StorageObjectOwner{Kind: storage.OwnerKindPlayer, ID: "player-1"},
		Collection: "progress",
		Limit:      1,
	})
	if err == nil {
		t.Fatal("ListStorageObjects() error = nil, want executor error")
	}
}

func TestStorageObjectRepositoryDefaultTestsDoNotRequireLivePostgreSQL(t *testing.T) {
	if os.Getenv("VIBIT_POSTGRES_TEST_DSN") != "" {
		t.Skip("live PostgreSQL environment is opt-in and not needed for this fake-executor test")
	}

	repository := NewStorageObjectRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{values: activeStorageObjectRowValues(time.Date(2026, 5, 22, 1, 2, 3, 0, time.UTC))},
		},
	})

	if _, err := repository.GetStorageObject(context.Background(), storage.GetStorageObjectInput{
		Owner:    storage.StorageObjectOwner{Kind: storage.OwnerKindPlayer, ID: "player-1"},
		Identity: storage.StorageObjectIdentity{Collection: "progress", Key: "tutorial"},
	}); err != nil {
		t.Fatalf("GetStorageObject() error = %v, want nil without live PostgreSQL", err)
	}
}

func TestPostgresUnitOfWorkCreatesStorageObjectRepository(t *testing.T) {
	executor := &recordingExecutor{}
	unit := UnitOfWork{executor: executor}

	repository, err := unit.NewStorageObjectRepository()
	if err != nil {
		t.Fatalf("NewStorageObjectRepository() error = %v, want nil", err)
	}
	if repository == nil {
		t.Fatal("NewStorageObjectRepository() = nil, want repository")
	}
}

func validCreateStorageObjectInput() storage.CreateStorageObjectInput {
	return storage.CreateStorageObjectInput{
		ObjectID:    "object-1",
		Owner:       storage.StorageObjectOwner{Kind: storage.OwnerKindPlayer, ID: "player-1"},
		Identity:    storage.StorageObjectIdentity{Collection: "progress", Key: "tutorial"},
		Value:       storage.StorageObjectValue{JSON: []byte(`{"level":3}`)},
		RequestedBy: "storage_service",
	}
}

type storageObjectRowOption func([]any)

func withStorageObjectID(objectID string) storageObjectRowOption {
	return func(values []any) {
		values[0] = objectID
	}
}

func withStorageObjectKey(key string) storageObjectRowOption {
	return func(values []any) {
		values[4] = key
	}
}

func withStorageObjectValue(value []byte) storageObjectRowOption {
	return func(values []any) {
		values[5] = append([]byte(nil), value...)
	}
}

func withStorageObjectVersion(version storage.StorageObjectVersion) storageObjectRowOption {
	return func(values []any) {
		values[6] = int64(version)
	}
}

func withStorageObjectDeletedAt(deletedAt time.Time) storageObjectRowOption {
	return func(values []any) {
		values[9] = deletedAt
	}
}

func activeStorageObjectRowValues(timestamp time.Time, options ...storageObjectRowOption) []any {
	values := []any{
		"object-1",
		"player",
		"player-1",
		"progress",
		"tutorial",
		[]byte(`{"level":3}`),
		int64(1),
		timestamp,
		timestamp,
		nil,
	}
	for _, option := range options {
		option(values)
	}
	return values
}

type storageObjectRow struct {
	values []any
	err    error
}

func (r storageObjectRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	assignStorageObjectValues("storage object row", dest, r.values)
	return nil
}

type storageObjectRows struct {
	values [][]any
	index  int
	err    error
	closed bool
}

func (r *storageObjectRows) Close() {
	r.closed = true
}

func (r *storageObjectRows) Err() error {
	return r.err
}

func (r *storageObjectRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}

func (r *storageObjectRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *storageObjectRows) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index += 1
	return true
}

func (r *storageObjectRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.values) {
		return errors.New("storage object rows: scan without current row")
	}
	assignStorageObjectValues("storage object rows", dest, r.values[r.index-1])
	return nil
}

func (r *storageObjectRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.values) {
		return nil, errors.New("storage object rows: values without current row")
	}
	return append([]any(nil), r.values[r.index-1]...), nil
}

func (r *storageObjectRows) RawValues() [][]byte {
	return nil
}

func (r *storageObjectRows) Conn() *pgx.Conn {
	return nil
}

func assignStorageObjectValues(label string, dest []any, values []any) {
	if len(dest) != len(values) {
		panic(label + ": destination count mismatch")
	}
	for i := range dest {
		assignStorageObjectValue(label, dest[i], values[i])
	}
}

func assignStorageObjectValue(label string, dest any, value any) {
	switch pointer := dest.(type) {
	case *string:
		*pointer = value.(string)
	case *[]byte:
		*pointer = append([]byte(nil), value.([]byte)...)
	case *int64:
		*pointer = value.(int64)
	case *time.Time:
		*pointer = value.(time.Time)
	case *pgtype.Timestamptz:
		switch timestamp := value.(type) {
		case nil:
			*pointer = pgtype.Timestamptz{}
		case time.Time:
			*pointer = pgtype.Timestamptz{Time: timestamp, Valid: true}
		case pgtype.Timestamptz:
			*pointer = timestamp
		default:
			panic(label + ": unsupported timestamptz value")
		}
	default:
		panic(label + ": unsupported destination type")
	}
}

func assertStorageConflictClass(t *testing.T, err error, want storage.StorageObjectConflictClass) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want storage conflict class %q", want)
	}
	var repositoryErr *storage.StorageObjectRepositoryError
	if !errors.As(err, &repositoryErr) {
		t.Fatalf("error = %T %[1]v, want StorageObjectRepositoryError", err)
	}
	if !errors.Is(err, storage.ErrStorageObjectConflict) {
		t.Fatalf("errors.Is(err, ErrStorageObjectConflict) = false, want true for %v", err)
	}
	if repositoryErr.Conflict.Class != want {
		t.Fatalf("conflict class = %q, want %q", repositoryErr.Conflict.Class, want)
	}
}

func assertStorageObjectErrorRedacted(t *testing.T, err error) {
	t.Helper()
	text := err.Error()
	for _, forbidden := range []string{
		"player-1",
		"object-1",
		"storage_objects_active_identity_uq",
		"storage_objects_owner_player_fk",
		"value_json",
		"secret",
		"{",
		"}",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("error %q leaks forbidden fragment %q", text, forbidden)
		}
	}
}

func TestStorageObjectRepositoryImplementsStorageRepository(t *testing.T) {
	var _ storage.Repository = (*StorageObjectRepository)(nil)
}

func TestStorageObjectRepositoryPreservesByteCopies(t *testing.T) {
	timestamp := time.Date(2026, 5, 22, 1, 2, 3, 0, time.UTC)
	source := []byte(`{"level":3}`)
	repository := NewStorageObjectRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			storageObjectRow{values: activeStorageObjectRowValues(timestamp, withStorageObjectValue(source))},
		},
	})

	record, err := repository.GetStorageObject(context.Background(), storage.GetStorageObjectInput{
		Owner:    storage.StorageObjectOwner{Kind: storage.OwnerKindPlayer, ID: "player-1"},
		Identity: storage.StorageObjectIdentity{Collection: "progress", Key: "tutorial"},
	})
	if err != nil {
		t.Fatalf("GetStorageObject() error = %v, want nil", err)
	}
	source[0] = '['
	if !reflect.DeepEqual(record.Value.JSON, []byte(`{"level":3}`)) {
		t.Fatalf("record.Value.JSON = %q, want copied JSON value", string(record.Value.JSON))
	}
}
