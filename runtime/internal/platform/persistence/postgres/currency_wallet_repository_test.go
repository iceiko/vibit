package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/currency"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCurrencyWalletRepositoryImplementsCurrencyRepository(t *testing.T) {
	var _ currency.Repository = (*CurrencyWalletRepository)(nil)
}

func TestCurrencyWalletRepositoryCreateInsertsActiveWallet(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 8, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{values: activeCurrencyWalletRowValues(createdAt.UTC())},
		},
	}
	repository := NewCurrencyWalletRepositoryForUnitOfWork(executor)

	record, err := repository.CreateCurrencyWallet(context.Background(), currency.CreateCurrencyWalletInput{
		WalletID:       " wallet-1 ",
		Owner:          currency.CurrencyWalletOwner{Kind: " player ", ID: " player-1 "},
		InitialState:   " active ",
		InitialVersion: 1,
		RequestedBy:    " currency_service ",
	})
	if err != nil {
		t.Fatalf("CreateCurrencyWallet() error = %v, want nil", err)
	}
	if record.WalletID != "wallet-1" ||
		record.Owner.Kind != currency.CurrencyWalletOwnerKindPlayer ||
		record.Owner.ID != "player-1" ||
		record.LifecycleState != currency.CurrencyWalletLifecycleActive ||
		record.WalletVersion != currency.InitialCurrencyWalletVersion {
		t.Fatalf("CreateCurrencyWallet() record = %#v, want normalized active wallet", record)
	}
	if record.CreatedAt.Location() != time.UTC || record.UpdatedAt.Location() != time.UTC || record.StateChangedAt.Location() != time.UTC {
		t.Fatalf("CreateCurrencyWallet() timestamps = %#v, want UTC", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "INSERT INTO currency_wallets")
	assertSQLContains(t, call.sql, "wallet_id")
	assertSQLContains(t, call.sql, "owner_kind")
	assertSQLContains(t, call.sql, "owner_id")
	assertSQLContains(t, call.sql, "lifecycle_state")
	assertSQLContains(t, call.sql, "wallet_version")
	assertSQLContains(t, call.sql, "RETURNING")
	assertArgs(t, call.args, "wallet-1", "player", "player-1", "active", int64(1))
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestCurrencyWalletRepositoryGetSelectsByWalletID(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{values: activeCurrencyWalletRowValues(createdAt)},
		},
	}
	repository := NewCurrencyWalletRepositoryForUnitOfWork(executor)

	record, err := repository.GetCurrencyWallet(context.Background(), currency.GetCurrencyWalletInput{
		WalletID: " wallet-1 ",
	})
	if err != nil {
		t.Fatalf("GetCurrencyWallet() error = %v, want nil", err)
	}
	if record.WalletID != "wallet-1" || record.LifecycleState != currency.CurrencyWalletLifecycleActive {
		t.Fatalf("GetCurrencyWallet() record = %#v, want active wallet-1", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "FROM currency_wallets")
	assertSQLContains(t, call.sql, "WHERE wallet_id = $1")
	assertArgs(t, call.args, "wallet-1")
}

func TestCurrencyWalletRepositoryGetForOwnerSelectsByOwnerIdentity(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{values: activeCurrencyWalletRowValues(createdAt)},
		},
	}
	repository := NewCurrencyWalletRepositoryForUnitOfWork(executor)

	record, err := repository.GetCurrencyWalletForOwner(context.Background(), currency.GetCurrencyWalletForOwnerInput{
		Owner: currency.CurrencyWalletOwner{Kind: " player ", ID: " player-1 "},
	})
	if err != nil {
		t.Fatalf("GetCurrencyWalletForOwner() error = %v, want nil", err)
	}
	if record.Owner.ID != "player-1" {
		t.Fatalf("GetCurrencyWalletForOwner() owner = %#v, want player-1", record.Owner)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "FROM currency_wallets")
	assertSQLContains(t, call.sql, "owner_kind = $1")
	assertSQLContains(t, call.sql, "owner_id = $2")
	assertArgs(t, call.args, "player", "player-1")
}

func TestCurrencyWalletRepositoryListBalancesIsWalletScopedOrderedAndBounded(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowsResponses: []pgx.Rows{
			&currencyWalletRows{
				values: [][]any{
					currencyBalanceRowValues(createdAt, withCurrencyBalanceCode("coins"), withCurrencyBalanceAmount(100)),
					currencyBalanceRowValues(createdAt, withCurrencyBalanceCode("gems"), withCurrencyBalanceAmount(5)),
					currencyBalanceRowValues(createdAt, withCurrencyBalanceCode("tickets"), withCurrencyBalanceAmount(9)),
				},
			},
		},
	}
	repository := NewCurrencyWalletRepositoryForUnitOfWork(executor)

	result, err := repository.ListCurrencyWalletBalances(context.Background(), currency.ListCurrencyWalletBalancesInput{
		WalletID:          " wallet-1 ",
		Limit:             2,
		AfterCurrencyCode: " coins ",
	})
	if err != nil {
		t.Fatalf("ListCurrencyWalletBalances() error = %v, want nil", err)
	}
	if len(result.Balances) != 2 || result.Balances[0].CurrencyCode != "coins" || result.Balances[1].CurrencyCode != "gems" {
		t.Fatalf("ListCurrencyWalletBalances() balances = %#v, want first two ordered rows", result.Balances)
	}
	if result.NextCurrencyCode != "tickets" {
		t.Fatalf("NextCurrencyCode = %q, want overflow cursor tickets", result.NextCurrencyCode)
	}

	if len(executor.queries) != 1 {
		t.Fatalf("queries len = %d, want 1", len(executor.queries))
	}
	call := executor.queries[0]
	assertSQLContains(t, call.sql, "FROM currency_wallet_balances")
	assertSQLContains(t, call.sql, "wallet_id = $1")
	assertSQLContains(t, call.sql, "currency_code > $2")
	assertSQLContains(t, call.sql, "ORDER BY currency_code")
	assertSQLContains(t, call.sql, "LIMIT $3")
	assertArgs(t, call.args, "wallet-1", "coins", int32(3))
}

func TestCurrencyWalletRepositoryRecordGrantInsertsTransactionAndUpdatesBalance(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 2, 0, 0, 0, time.UTC)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{values: currencyTransactionRowValues(createdAt, withCurrencyTransactionAmount(25), withCurrencyTransactionBalanceAfter(125))},
		},
	}
	repository := NewCurrencyWalletRepositoryForUnitOfWork(executor)
	walletVersion := currency.CurrencyWalletVersion(7)
	balanceVersion := currency.CurrencyBalanceVersion(3)
	metadata := []byte(`{"source":"quest"}`)

	record, err := repository.RecordCurrencyGrant(context.Background(), currency.RecordCurrencyGrantInput{
		WalletID:               " wallet-1 ",
		TransactionID:          " txn-1 ",
		CurrencyCode:           " coins ",
		Amount:                 25,
		IdempotencyKey:         " grant-key ",
		IdempotencyScope:       " quest ",
		Actor:                  currency.CurrencyWalletActor{Kind: " system ", ID: " rewards "},
		ReasonCode:             " quest_reward ",
		ExternalReference:      " quest-1 ",
		MetadataJSON:           metadata,
		ExpectedWalletVersion:  &walletVersion,
		ExpectedBalanceVersion: &balanceVersion,
	})
	if err != nil {
		t.Fatalf("RecordCurrencyGrant() error = %v, want nil", err)
	}
	if record.TransactionKind != currency.CurrencyWalletTransactionGrant || record.AmountDelta != 25 || record.BalanceAfter != 125 {
		t.Fatalf("RecordCurrencyGrant() record = %#v, want grant transaction", record)
	}
	metadata[0] = '['
	if string(record.MetadataJSON) != `{"source":"quest"}` {
		t.Fatalf("record.MetadataJSON = %q, want copied metadata", string(record.MetadataJSON))
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "INSERT INTO currency_wallet_transactions")
	assertSQLContains(t, call.sql, "INSERT INTO currency_wallet_balances")
	assertSQLContains(t, call.sql, "ON CONFLICT (wallet_id, currency_code)")
	assertSQLContains(t, call.sql, "balance_amount = currency_wallet_balances.balance_amount + EXCLUDED.balance_amount")
	assertSQLContains(t, call.sql, "transaction_kind")
	assertSQLContains(t, call.sql, "balance_after")
	assertSQLContains(t, call.sql, "wallet_version = $14")
	assertSQLContains(t, call.sql, "balance_version = $15")
	assertArgs(t,
		call.args,
		"txn-1",
		"wallet-1",
		"coins",
		"grant",
		int64(25),
		"grant-key",
		"quest",
		"system",
		pgtype.Text{String: "rewards", Valid: true},
		pgtype.Text{String: "quest_reward", Valid: true},
		pgtype.Text{String: "quest-1", Valid: true},
		[]byte(`{"source":"quest"}`),
		int64(25),
		int64(7),
		int64(3),
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestCurrencyWalletRepositoryRecordSpendChecksBalanceAndMapsInsufficientBalance(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 2, 30, 0, 0, time.UTC)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{values: currencyTransactionRowValues(createdAt, withCurrencyTransactionKind("spend"), withCurrencyTransactionAmount(-40), withCurrencyTransactionBalanceAfter(60))},
		},
	}
	repository := NewCurrencyWalletRepositoryForUnitOfWork(executor)

	record, err := repository.RecordCurrencySpend(context.Background(), validRecordCurrencySpendInput())
	if err != nil {
		t.Fatalf("RecordCurrencySpend() error = %v, want nil", err)
	}
	if record.TransactionKind != currency.CurrencyWalletTransactionSpend || record.AmountDelta != -40 || record.BalanceAfter != 60 {
		t.Fatalf("RecordCurrencySpend() record = %#v, want spend transaction", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "UPDATE currency_wallet_balances")
	assertSQLContains(t, call.sql, "balance_amount = balance_amount - $13")
	assertSQLContains(t, call.sql, "balance_amount >= $13")
	assertSQLContains(t, call.sql, "INSERT INTO currency_wallet_transactions")
	assertArgs(t,
		call.args,
		"txn-2",
		"wallet-1",
		"coins",
		"spend",
		int64(-40),
		"spend-key",
		"shop",
		"player",
		pgtype.Text{String: "player-1", Valid: true},
		pgtype.Text{String: "shop_purchase", Valid: true},
		pgtype.Text{String: "sku-1", Valid: true},
		[]byte(`{"sku":"sku-1"}`),
		int64(40),
		nil,
		nil,
	)

	repository = NewCurrencyWalletRepositoryForUnitOfWork(&recordingExecutor{})
	_, err = repository.RecordCurrencySpend(context.Background(), validRecordCurrencySpendInput())
	assertCurrencyWalletConflictClass(t, err, currency.CurrencyWalletConflictInsufficientBalance)
	assertCurrencyWalletErrorRedacted(t, err)
}

func TestCurrencyWalletRepositoryListTransactionsIsWalletScopedFilteredOrderedAndBounded(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 3, 0, 0, 0, time.UTC)
	executor := &recordingExecutor{
		rowsResponses: []pgx.Rows{
			&currencyWalletRows{
				values: [][]any{
					currencyTransactionRowValues(createdAt, withCurrencyTransactionID("txn-1")),
					currencyTransactionRowValues(createdAt.Add(time.Second), withCurrencyTransactionID("txn-2")),
					currencyTransactionRowValues(createdAt.Add(2*time.Second), withCurrencyTransactionID("txn-3")),
				},
			},
		},
	}
	repository := NewCurrencyWalletRepositoryForUnitOfWork(executor)

	result, err := repository.ListCurrencyWalletTransactions(context.Background(), currency.ListCurrencyWalletTransactionsInput{
		WalletID:             " wallet-1 ",
		CurrencyCode:         " coins ",
		Limit:                2,
		AfterTransactionID:   " txn-0 ",
		AfterTransactionTime: createdAt.Add(-time.Second),
	})
	if err != nil {
		t.Fatalf("ListCurrencyWalletTransactions() error = %v, want nil", err)
	}
	if len(result.Transactions) != 2 || result.Transactions[0].TransactionID != "txn-1" || result.Transactions[1].TransactionID != "txn-2" {
		t.Fatalf("ListCurrencyWalletTransactions() transactions = %#v, want first two rows", result.Transactions)
	}
	if result.NextTransactionID != "txn-3" || !result.NextTransactionCreateAt.Equal(createdAt.Add(2*time.Second)) {
		t.Fatalf("next cursor = (%q, %s), want txn-3 overflow cursor", result.NextTransactionID, result.NextTransactionCreateAt)
	}

	if len(executor.queries) != 1 {
		t.Fatalf("queries len = %d, want 1", len(executor.queries))
	}
	call := executor.queries[0]
	assertSQLContains(t, call.sql, "FROM currency_wallet_transactions")
	assertSQLContains(t, call.sql, "wallet_id = $1")
	assertSQLContains(t, call.sql, "($2::text = '' OR currency_code = $2)")
	assertSQLContains(t, call.sql, "(created_at, transaction_id) > ($3, $4)")
	assertSQLContains(t, call.sql, "ORDER BY created_at, transaction_id")
	assertSQLContains(t, call.sql, "LIMIT $5")
	assertArgs(t, call.args, "wallet-1", "coins", createdAt.Add(-time.Second).UTC(), "txn-0", int32(3))
}

func TestCurrencyWalletRepositoryMapsErrorsAndRedactsDetails(t *testing.T) {
	repository := NewCurrencyWalletRepositoryForUnitOfWork(&recordingExecutor{})

	_, err := repository.GetCurrencyWallet(context.Background(), currency.GetCurrencyWalletInput{WalletID: "wallet-1"})
	assertCurrencyWalletConflictClass(t, err, currency.CurrencyWalletConflictWalletNotFound)

	duplicateOwner := &pgconn.PgError{Code: "23505", ConstraintName: "currency_wallets_owner_uq", Detail: "player-1 wallet-1"}
	repository = NewCurrencyWalletRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{err: duplicateOwner},
		},
	})
	_, err = repository.CreateCurrencyWallet(context.Background(), validCreateCurrencyWalletInput())
	assertCurrencyWalletConflictClass(t, err, currency.CurrencyWalletConflictWalletAlreadyExists)
	assertCurrencyWalletErrorRedacted(t, err)

	duplicateTransaction := &pgconn.PgError{Code: "23505", ConstraintName: "currency_wallet_transactions_idempotency_uq", Detail: "wallet-1 grant-key"}
	repository = NewCurrencyWalletRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{err: duplicateTransaction},
		},
	})
	_, err = repository.RecordCurrencyGrant(context.Background(), validRecordCurrencyGrantInput())
	assertCurrencyWalletConflictClass(t, err, currency.CurrencyWalletConflictDuplicateTransaction)
	assertCurrencyWalletErrorRedacted(t, err)

	constraint := &pgconn.PgError{Code: "23503", ConstraintName: "currency_wallets_owner_player_fk", Detail: "player-1"}
	repository = NewCurrencyWalletRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{err: constraint},
		},
	})
	_, err = repository.CreateCurrencyWallet(context.Background(), validCreateCurrencyWalletInput())
	assertCurrencyWalletConflictClass(t, err, currency.CurrencyWalletConflictStorageUnavailable)
	assertCurrencyWalletErrorRedacted(t, err)
}

func TestCurrencyWalletRepositoryRejectsInvalidRows(t *testing.T) {
	values := activeCurrencyWalletRowValues(time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC))
	values[3] = "archived"
	repository := NewCurrencyWalletRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{values: values},
		},
	})

	_, err := repository.GetCurrencyWallet(context.Background(), currency.GetCurrencyWalletInput{WalletID: "wallet-1"})
	assertCurrencyWalletConflictClass(t, err, currency.CurrencyWalletConflictStorageUnavailable)
	assertCurrencyWalletErrorRedacted(t, err)
}

func TestCurrencyWalletRepositoryRequiresUnitOfWorkExecutor(t *testing.T) {
	repository := NewCurrencyWalletRepositoryForUnitOfWork(nil)

	_, err := repository.CreateCurrencyWallet(context.Background(), validCreateCurrencyWalletInput())
	if err == nil {
		t.Fatal("CreateCurrencyWallet() error = nil, want executor error")
	}

	_, err = repository.ListCurrencyWalletBalances(context.Background(), currency.ListCurrencyWalletBalancesInput{
		WalletID: "wallet-1",
		Limit:    1,
	})
	if err == nil {
		t.Fatal("ListCurrencyWalletBalances() error = nil, want executor error")
	}
}

func TestCurrencyWalletRepositoryDefaultTestsDoNotRequireLivePostgreSQL(t *testing.T) {
	if os.Getenv("VIBIT_POSTGRES_TEST_DSN") != "" {
		t.Skip("live PostgreSQL environment is opt-in and not needed for this fake-executor test")
	}

	repository := NewCurrencyWalletRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			currencyWalletRow{values: activeCurrencyWalletRowValues(time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC))},
		},
	})

	if _, err := repository.GetCurrencyWallet(context.Background(), currency.GetCurrencyWalletInput{WalletID: "wallet-1"}); err != nil {
		t.Fatalf("GetCurrencyWallet() error = %v, want nil without live PostgreSQL", err)
	}
}

func TestPostgresUnitOfWorkCreatesCurrencyWalletRepository(t *testing.T) {
	executor := &recordingExecutor{}
	unit := UnitOfWork{executor: executor}

	repository, err := unit.NewCurrencyWalletRepository()
	if err != nil {
		t.Fatalf("NewCurrencyWalletRepository() error = %v, want nil", err)
	}
	if repository == nil {
		t.Fatal("NewCurrencyWalletRepository() = nil, want repository")
	}
}

func validCreateCurrencyWalletInput() currency.CreateCurrencyWalletInput {
	return currency.CreateCurrencyWalletInput{
		WalletID:       "wallet-1",
		Owner:          currency.CurrencyWalletOwner{Kind: currency.CurrencyWalletOwnerKindPlayer, ID: "player-1"},
		InitialState:   currency.CurrencyWalletLifecycleActive,
		InitialVersion: currency.InitialCurrencyWalletVersion,
		RequestedBy:    "currency_service",
	}
}

func validRecordCurrencyGrantInput() currency.RecordCurrencyGrantInput {
	return currency.RecordCurrencyGrantInput{
		WalletID:          "wallet-1",
		TransactionID:     "txn-1",
		CurrencyCode:      "coins",
		Amount:            25,
		IdempotencyKey:    "grant-key",
		IdempotencyScope:  "quest",
		Actor:             currency.CurrencyWalletActor{Kind: currency.CurrencyWalletActorSystem, ID: "rewards"},
		ReasonCode:        "quest_reward",
		ExternalReference: "quest-1",
		MetadataJSON:      []byte(`{"source":"quest"}`),
	}
}

func validRecordCurrencySpendInput() currency.RecordCurrencySpendInput {
	return currency.RecordCurrencySpendInput{
		WalletID:          "wallet-1",
		TransactionID:     "txn-2",
		CurrencyCode:      "coins",
		Amount:            40,
		IdempotencyKey:    "spend-key",
		IdempotencyScope:  "shop",
		Actor:             currency.CurrencyWalletActor{Kind: currency.CurrencyWalletActorPlayer, ID: "player-1"},
		ReasonCode:        "shop_purchase",
		ExternalReference: "sku-1",
		MetadataJSON:      []byte(`{"sku":"sku-1"}`),
	}
}

type currencyWalletRowOption func([]any)

func withCurrencyWalletID(walletID string) currencyWalletRowOption {
	return func(values []any) {
		values[0] = walletID
	}
}

func withCurrencyBalanceCode(code string) currencyWalletRowOption {
	return func(values []any) {
		values[1] = code
	}
}

func withCurrencyBalanceAmount(amount currency.CurrencyAmount) currencyWalletRowOption {
	return func(values []any) {
		values[2] = int64(amount)
	}
}

func withCurrencyTransactionID(transactionID string) currencyWalletRowOption {
	return func(values []any) {
		values[0] = transactionID
	}
}

func withCurrencyTransactionKind(kind string) currencyWalletRowOption {
	return func(values []any) {
		values[3] = kind
	}
}

func withCurrencyTransactionAmount(amount currency.CurrencyAmount) currencyWalletRowOption {
	return func(values []any) {
		values[4] = int64(amount)
	}
}

func withCurrencyTransactionBalanceAfter(amount currency.CurrencyAmount) currencyWalletRowOption {
	return func(values []any) {
		values[5] = int64(amount)
	}
}

func activeCurrencyWalletRowValues(timestamp time.Time, options ...currencyWalletRowOption) []any {
	values := []any{
		"wallet-1",
		"player",
		"player-1",
		"active",
		int64(1),
		timestamp,
		timestamp,
		timestamp,
		nil,
		nil,
	}
	for _, option := range options {
		option(values)
	}
	return values
}

func currencyBalanceRowValues(timestamp time.Time, options ...currencyWalletRowOption) []any {
	values := []any{
		"wallet-1",
		"coins",
		int64(100),
		int64(1),
		timestamp,
		timestamp,
	}
	for _, option := range options {
		option(values)
	}
	return values
}

func currencyTransactionRowValues(timestamp time.Time, options ...currencyWalletRowOption) []any {
	values := []any{
		"txn-1",
		"wallet-1",
		"coins",
		"grant",
		int64(25),
		int64(125),
		"grant-key",
		"quest",
		"system",
		"rewards",
		"quest_reward",
		"quest-1",
		[]byte(`{"source":"quest"}`),
		timestamp,
	}
	for _, option := range options {
		option(values)
	}
	return values
}

type currencyWalletRow struct {
	values []any
	err    error
}

func (r currencyWalletRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	assignCurrencyWalletValues("currency wallet row", dest, r.values)
	return nil
}

type currencyWalletRows struct {
	values [][]any
	index  int
	err    error
	closed bool
}

func (r *currencyWalletRows) Close() {
	r.closed = true
}

func (r *currencyWalletRows) Err() error {
	return r.err
}

func (r *currencyWalletRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}

func (r *currencyWalletRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *currencyWalletRows) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index += 1
	return true
}

func (r *currencyWalletRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.values) {
		return errors.New("currency wallet rows: scan without current row")
	}
	assignCurrencyWalletValues("currency wallet rows", dest, r.values[r.index-1])
	return nil
}

func (r *currencyWalletRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.values) {
		return nil, errors.New("currency wallet rows: values without current row")
	}
	return append([]any(nil), r.values[r.index-1]...), nil
}

func (r *currencyWalletRows) RawValues() [][]byte {
	return nil
}

func (r *currencyWalletRows) Conn() *pgx.Conn {
	return nil
}

func assignCurrencyWalletValues(label string, dest []any, values []any) {
	if len(dest) != len(values) {
		panic(label + ": destination count mismatch")
	}
	for i := range dest {
		assignCurrencyWalletValue(label, dest[i], values[i])
	}
}

func assignCurrencyWalletValue(label string, dest any, value any) {
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
	case *pgtype.Text:
		switch text := value.(type) {
		case nil:
			*pointer = pgtype.Text{}
		case string:
			*pointer = pgtype.Text{String: text, Valid: true}
		case pgtype.Text:
			*pointer = text
		default:
			panic(label + ": unsupported text value")
		}
	default:
		panic(label + ": unsupported destination type")
	}
}

func assertCurrencyWalletConflictClass(t *testing.T, err error, want currency.CurrencyWalletConflictClass) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want currency wallet conflict class %q", want)
	}
	var repositoryErr *currency.CurrencyWalletRepositoryError
	if !errors.As(err, &repositoryErr) {
		t.Fatalf("error = %T %[1]v, want CurrencyWalletRepositoryError", err)
	}
	if !errors.Is(err, currency.ErrCurrencyWalletConflict) {
		t.Fatalf("errors.Is(err, ErrCurrencyWalletConflict) = false, want true for %v", err)
	}
	if repositoryErr.Conflict.Class != want {
		t.Fatalf("conflict class = %q, want %q", repositoryErr.Conflict.Class, want)
	}
}

func assertCurrencyWalletErrorRedacted(t *testing.T, err error) {
	t.Helper()
	text := err.Error()
	for _, forbidden := range []string{
		"player-1",
		"wallet-1",
		"txn-1",
		"grant-key",
		"spend-key",
		"currency_wallets_owner_uq",
		"currency_wallet_transactions_idempotency_uq",
		"currency_wallets_owner_player_fk",
		"balance_amount",
		"metadata_json",
		"quest",
		"sku",
		"{",
		"}",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("error %q leaks forbidden fragment %q", text, forbidden)
		}
	}
}
