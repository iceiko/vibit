package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/player"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestPlayerAccountRepositoryCreatePlayerAccountInsertsAccountAndEvent(t *testing.T) {
	occurredAt := time.Date(2026, 5, 14, 13, 20, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	createdAt := occurredAt.Add(2 * time.Second)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			playerAccountRow{
				values: []any{
					"player-1",
					"Player One",
					"active",
					createdAt,
					createdAt,
					nil,
					nil,
				},
			},
		},
	}
	repository := NewPlayerAccountRepositoryForUnitOfWork(executor)

	account, err := repository.CreatePlayerAccount(context.Background(), player.CreatePlayerAccountMutation{
		EventID:     " event-1 ",
		OccurredAt:  occurredAt,
		PlayerID:    " player-1 ",
		DisplayName: " Player One ",
		RequestedBy: " maintainer ",
	})
	if err != nil {
		t.Fatalf("CreatePlayerAccount() error = %v, want nil", err)
	}

	if account.PlayerID != "player-1" || account.DisplayName != "Player One" || account.AccountState != player.AccountStateActive {
		t.Fatalf("CreatePlayerAccount() account = %#v, want normalized active player account", account)
	}
	if account.CreatedAt.Location() != time.UTC || account.UpdatedAt.Location() != time.UTC {
		t.Fatalf("CreatePlayerAccount() timestamps = %s/%s, want UTC", account.CreatedAt.Location(), account.UpdatedAt.Location())
	}
	if account.DisabledAt != nil || account.DeletedAt != nil {
		t.Fatalf("CreatePlayerAccount() nullable lifecycle timestamps = %#v/%#v, want nil", account.DisabledAt, account.DeletedAt)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	accountInsert := executor.queryRowCalls[0]
	assertSQLContains(t, accountInsert.sql, "INSERT INTO player_accounts")
	assertSQLContains(t, accountInsert.sql, "RETURNING player_id, display_name, account_state, created_at, updated_at, disabled_at, deleted_at")
	assertArgs(t, accountInsert.args, "player-1", "Player One", string(player.AccountStateActive), occurredAt.UTC())

	if len(executor.execs) != 1 {
		t.Fatalf("execs len = %d, want 1", len(executor.execs))
	}
	eventInsert := executor.execs[0]
	assertSQLContains(t, eventInsert.sql, "INSERT INTO player_account_events")
	assertArgs(
		t,
		eventInsert.args,
		"event-1",
		player.EventPlayerAccountCreated,
		occurredAt.UTC(),
		"player-1",
		"maintainer",
		string(player.AccountStateActive),
		"Player One",
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestPlayerAccountRepositoryGetPlayerAccountSelectsAndMapsLifecycleRow(t *testing.T) {
	createdAt := time.Date(2026, 5, 14, 8, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	updatedAt := createdAt.Add(30 * time.Minute)
	disabledAt := updatedAt.Add(10 * time.Minute)
	deletedAt := updatedAt.Add(20 * time.Minute)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			playerAccountRow{
				values: []any{
					"player-1",
					"Player One",
					"deleted",
					createdAt,
					updatedAt,
					disabledAt,
					deletedAt,
				},
			},
		},
	}
	repository := NewPlayerAccountRepositoryForUnitOfWork(executor)

	account, err := repository.GetPlayerAccount(context.Background(), " player-1 ")
	if err != nil {
		t.Fatalf("GetPlayerAccount() error = %v, want nil", err)
	}

	if account.PlayerID != "player-1" || account.DisplayName != "Player One" || account.AccountState != player.AccountStateDeleted {
		t.Fatalf("GetPlayerAccount() account = %#v, want mapped deleted account", account)
	}
	if !account.CreatedAt.Equal(createdAt.UTC()) || account.CreatedAt.Location() != time.UTC {
		t.Fatalf("CreatedAt = %s, want UTC %s", account.CreatedAt, createdAt.UTC())
	}
	if !account.UpdatedAt.Equal(updatedAt.UTC()) || account.UpdatedAt.Location() != time.UTC {
		t.Fatalf("UpdatedAt = %s, want UTC %s", account.UpdatedAt, updatedAt.UTC())
	}
	if account.DisabledAt == nil || !account.DisabledAt.Equal(disabledAt.UTC()) || account.DisabledAt.Location() != time.UTC {
		t.Fatalf("DisabledAt = %#v, want UTC %s", account.DisabledAt, disabledAt.UTC())
	}
	if account.DeletedAt == nil || !account.DeletedAt.Equal(deletedAt.UTC()) || account.DeletedAt.Location() != time.UTC {
		t.Fatalf("DeletedAt = %#v, want UTC %s", account.DeletedAt, deletedAt.UTC())
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	query := executor.queryRowCalls[0]
	assertSQLContains(t, query.sql, "FROM player_accounts")
	assertSQLContains(t, query.sql, "WHERE player_id = $1")
	assertArgs(t, query.args, "player-1")
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestPlayerAccountRepositoryMapsMissingRowToNotFound(t *testing.T) {
	repository := NewPlayerAccountRepositoryForUnitOfWork(&recordingExecutor{})

	_, err := repository.GetPlayerAccount(context.Background(), "player-1")
	if !errors.Is(err, ErrPlayerAccountNotFound) {
		t.Fatalf("GetPlayerAccount() error = %v, want ErrPlayerAccountNotFound", err)
	}
}

func TestPlayerAccountRepositoryMapsDuplicateAccountOrEventToConflict(t *testing.T) {
	duplicate := &pgconn.PgError{Code: "23505", ConstraintName: "player_accounts_pkey"}
	repository := NewPlayerAccountRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			playerAccountRow{err: duplicate},
		},
	})

	_, err := repository.CreatePlayerAccount(context.Background(), validCreatePlayerAccountMutation())
	if !errors.Is(err, ErrPlayerAccountConflict) {
		t.Fatalf("CreatePlayerAccount() account insert error = %v, want ErrPlayerAccountConflict", err)
	}
	assertDoesNotLeakPgError(t, err)

	eventDuplicate := &pgconn.PgError{Code: "23505", ConstraintName: "player_account_events_pkey"}
	repository = NewPlayerAccountRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			playerAccountRow{values: activePlayerAccountRowValues()},
		},
		execErr: eventDuplicate,
	})

	_, err = repository.CreatePlayerAccount(context.Background(), validCreatePlayerAccountMutation())
	if !errors.Is(err, ErrPlayerAccountConflict) {
		t.Fatalf("CreatePlayerAccount() event insert error = %v, want ErrPlayerAccountConflict", err)
	}
	assertDoesNotLeakPgError(t, err)
}

func TestPlayerAccountRepositoryMapsConstraintViolations(t *testing.T) {
	constraint := &pgconn.PgError{Code: "23514", ConstraintName: "player_accounts_account_state_valid"}
	repository := NewPlayerAccountRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			playerAccountRow{err: constraint},
		},
	})

	_, err := repository.CreatePlayerAccount(context.Background(), validCreatePlayerAccountMutation())
	if !errors.Is(err, ErrPlayerAccountConstraint) {
		t.Fatalf("CreatePlayerAccount() error = %v, want ErrPlayerAccountConstraint", err)
	}
	assertDoesNotLeakPgError(t, err)

	repository = NewPlayerAccountRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			playerAccountRow{values: []any{
				"player-1",
				"Player One",
				"suspended",
				time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC),
				time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC),
				nil,
				nil,
			}},
		},
	})
	_, err = repository.GetPlayerAccount(context.Background(), "player-1")
	if !errors.Is(err, ErrPlayerAccountConstraint) {
		t.Fatalf("GetPlayerAccount() invalid row error = %v, want ErrPlayerAccountConstraint", err)
	}
}

func TestPlayerAccountRepositoryRequiresUnitOfWorkExecutor(t *testing.T) {
	repository := NewPlayerAccountRepositoryForUnitOfWork(nil)

	_, err := repository.CreatePlayerAccount(context.Background(), validCreatePlayerAccountMutation())
	if err == nil {
		t.Fatal("CreatePlayerAccount() error = nil, want executor error")
	}

	_, err = repository.GetPlayerAccount(context.Background(), "player-1")
	if err == nil {
		t.Fatal("GetPlayerAccount() error = nil, want executor error")
	}
}

func TestPlayerAccountRepositoryDefaultTestsDoNotRequireLivePostgreSQL(t *testing.T) {
	if os.Getenv("VIBIT_POSTGRES_TEST_DSN") != "" {
		t.Skip("live PostgreSQL environment is opt-in and not needed for this fake-executor test")
	}

	repository := NewPlayerAccountRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			playerAccountRow{values: activePlayerAccountRowValues()},
		},
	})

	if _, err := repository.GetPlayerAccount(context.Background(), "player-1"); err != nil {
		t.Fatalf("GetPlayerAccount() error = %v, want nil without live PostgreSQL", err)
	}
}

func TestPostgresUnitOfWorkCreatesPlayerAccountRepository(t *testing.T) {
	executor := &recordingExecutor{}
	unit := UnitOfWork{executor: executor}

	repository, err := unit.NewPlayerAccountRepository()
	if err != nil {
		t.Fatalf("NewPlayerAccountRepository() error = %v, want nil", err)
	}
	if repository == nil {
		t.Fatal("NewPlayerAccountRepository() = nil, want repository")
	}
}

func validCreatePlayerAccountMutation() player.CreatePlayerAccountMutation {
	return player.CreatePlayerAccountMutation{
		EventID:     "event-1",
		OccurredAt:  time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC),
		PlayerID:    "player-1",
		DisplayName: "Player One",
		RequestedBy: "maintainer",
	}
}

func activePlayerAccountRowValues() []any {
	timestamp := time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC)
	return []any{
		"player-1",
		"Player One",
		"active",
		timestamp,
		timestamp,
		nil,
		nil,
	}
}

func assertDoesNotLeakPgError(t *testing.T, err error) {
	t.Helper()
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		t.Fatalf("error %v leaks pgconn.PgError", err)
	}
}

type playerAccountRow struct {
	values []any
	err    error
}

func (r playerAccountRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("player account row: destination count mismatch")
	}
	for i := range dest {
		assignPlayerAccountRowValue(dest[i], r.values[i])
	}
	return nil
}

func assignPlayerAccountRowValue(dest any, value any) {
	switch pointer := dest.(type) {
	case *string:
		*pointer = value.(string)
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
			panic("player account row: unsupported timestamptz value")
		}
	default:
		panic("player account row: unsupported destination type")
	}
}
