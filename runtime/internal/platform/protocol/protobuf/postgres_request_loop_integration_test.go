package protobuf

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/bootstrap"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/platform/migrations"
	"github.com/iceiko/vibit/runtime/internal/platform/persistence/postgres"
)

const (
	envPostgresTestDSN              = "VIBIT_POSTGRES_TEST_DSN"
	envPostgresTestCleanup          = "VIBIT_POSTGRES_TEST_CLEANUP"
	envPostgresTestAllowDestructive = "VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE"
)

func TestPostgresPersistentInventoryRequestLoop(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv(envPostgresTestDSN))
	if dsn == "" {
		t.Skip("live PostgreSQL verification skipped because VIBIT_POSTGRES_TEST_DSN was not set")
	}

	cleanupMode := postgresTestCleanupMode()
	if cleanupMode != "drop_schema" {
		t.Skip("live PostgreSQL verification skipped because TestPostgresPersistentInventoryRequestLoop requires VIBIT_POSTGRES_TEST_CLEANUP=drop_schema")
	}
	if os.Getenv(envPostgresTestAllowDestructive) != "1" {
		t.Skip("live PostgreSQL verification skipped because destructive cleanup requires VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1")
	}

	ctx := context.Background()
	pool, err := postgres.OpenPool(ctx, postgres.Config{DSN: dsn, MaxConns: 4, MinConns: 1})
	if err != nil {
		t.Fatalf("OpenPool() error = %v, want nil", err)
	}
	defer pool.Close()

	sqlDB, err := postgres.OpenSQLDBFromPool(pool)
	if err != nil {
		t.Fatalf("OpenSQLDBFromPool() error = %v, want nil", err)
	}
	defer sqlDB.Close()

	postgresCleanupSchema(t, ctx, pool)
	defer postgresCleanupSchema(t, ctx, pool)

	runner, err := migrations.NewPostgresRunnerFromDir(sqlDB, postgresMigrationDir(t))
	if err != nil {
		t.Fatalf("NewPostgresRunnerFromDir() error = %v, want nil", err)
	}

	statusBefore, err := runner.Status(ctx)
	if err != nil {
		t.Fatalf("Status() before apply error = %v, want nil", err)
	}
	if !migrationStatusContains(t, statusBefore, 1, "pending") {
		t.Fatalf("Status() before apply = %#v, want migration 1 pending", statusBefore)
	}

	applied, err := runner.Apply(ctx)
	if err != nil {
		t.Fatalf("Apply() error = %v, want nil", err)
	}
	if len(applied) != 1 || applied[0].Source.Version != 1 || applied[0].Direction != "up" {
		t.Fatalf("Apply() = %#v, want migration 1 applied up", applied)
	}

	statusAfter, err := runner.Status(ctx)
	if err != nil {
		t.Fatalf("Status() after apply error = %v, want nil", err)
	}
	if !migrationStatusContains(t, statusAfter, 1, "applied") {
		t.Fatalf("Status() after apply = %#v, want migration 1 applied", statusAfter)
	}

	dispatcher, err := bootstrap.NewInventoryDispatcher(bootstrap.InventoryOptions{
		Repositories: bootstrap.PostgresInventoryRepositoryProvider{
			QueryRepository: postgres.NewInventoryRepositoryForUnitOfWork(pool),
		},
		PermissionPolicy: inventory.StaticPermissionPolicy{GrantAllowed: true, ReadAllowed: true},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 16},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "postgres-request-loop-event"},
		Clock:            requestLoopClock{},
	})
	if err != nil {
		t.Fatalf("NewInventoryDispatcher() error = %v, want nil", err)
	}

	handler := FrameHandler{
		Dispatcher: app.TransactionalDispatcher{
			Dispatcher: dispatcher,
			Runner:     postgres.NewPoolRunner(pool),
		},
	}

	grantFrame := mustMarshalEnvelope(t, inventory.GrantItemRoute(), &inventoryv1.GrantItemRequest{
		PlayerId:    "player-1",
		ItemId:      "item-1",
		Quantity:    3,
		Reason:      "integration-test",
		RequestedBy: "admin-1",
	})
	grantResponses, err := handler.HandleFrame(ctx, FrameRequest{
		ConnectionID: "postgres-ws-1",
		RemoteAddr:   "127.0.0.1:1",
		Payload:      grantFrame,
	})
	if err != nil {
		t.Fatalf("HandleFrame(GrantItem) error = %v, want nil", err)
	}
	grantEnvelope := singleResponseEnvelope(t, grantResponses)
	if grantEnvelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_COMMAND {
		t.Fatalf("GrantItem response kind = %v, want command", grantEnvelope.GetKind())
	}
	grantPayload, err := DecodePayload(grantEnvelope.GetPayloadType(), grantEnvelope.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(GrantItem) error = %v, want nil", err)
	}
	grantResponse, ok := grantPayload.(*inventoryv1.GrantItemResponse)
	if !ok {
		t.Fatalf("GrantItem response payload = %T, want *inventoryv1.GrantItemResponse", grantPayload)
	}
	if grantResponse.GetNewQuantity() != 3 {
		t.Fatalf("GrantItem new quantity = %d, want 3", grantResponse.GetNewQuantity())
	}

	getFrame := mustMarshalEnvelope(t, inventory.GetInventoryRoute(), &inventoryv1.GetInventoryRequest{
		PlayerId:    "player-1",
		RequestedBy: "player-1",
	})
	getResponses, err := handler.HandleFrame(ctx, FrameRequest{
		ConnectionID: "postgres-ws-1",
		RemoteAddr:   "127.0.0.1:1",
		Payload:      getFrame,
	})
	if err != nil {
		t.Fatalf("HandleFrame(GetInventory) error = %v, want nil", err)
	}
	getEnvelope := singleResponseEnvelope(t, getResponses)
	if getEnvelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_QUERY {
		t.Fatalf("GetInventory response kind = %v, want query", getEnvelope.GetKind())
	}
	getPayload, err := DecodePayload(getEnvelope.GetPayloadType(), getEnvelope.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(GetInventory) error = %v, want nil", err)
	}
	getResponse, ok := getPayload.(*inventoryv1.GetInventoryResponse)
	if !ok {
		t.Fatalf("GetInventory response payload = %T, want *inventoryv1.GetInventoryResponse", getPayload)
	}
	if getResponse.GetPlayerId() != "player-1" {
		t.Fatalf("GetInventory player_id = %q, want player-1", getResponse.GetPlayerId())
	}
	if len(getResponse.GetItems()) != 1 {
		t.Fatalf("GetInventory item count = %d, want 1", len(getResponse.GetItems()))
	}
	item := getResponse.GetItems()[0]
	if item.GetItemId() != "item-1" || item.GetQuantity() != 3 {
		t.Fatalf("GetInventory item = %#v, want item-1 quantity 3", item)
	}
}

func postgresTestCleanupMode() string {
	cleanupMode := strings.ToLower(strings.TrimSpace(os.Getenv(envPostgresTestCleanup)))
	if cleanupMode == "" {
		return "drop_schema"
	}
	switch cleanupMode {
	case "drop_schema", "truncate", "keep":
		return cleanupMode
	default:
		return "drop_schema"
	}
}

func postgresCleanupSchema(t *testing.T, ctx context.Context, executor postgres.Executor) {
	t.Helper()

	for _, query := range []string{
		"DROP TABLE IF EXISTS inventory_item_grants",
		"DROP TABLE IF EXISTS inventory_items",
		"DROP TABLE IF EXISTS inventory_accounts",
		"DROP TABLE IF EXISTS goose_db_version",
	} {
		if _, err := executor.Exec(ctx, query); err != nil {
			t.Fatalf("cleanup query %q error = %v, want nil", query, err)
		}
	}
}

func postgresMigrationDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed")
	}
	dir := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", "migrations", "postgres"))
	if _, err := os.Stat(filepath.Join(dir, "000001_create_inventory_state.sql")); err != nil {
		t.Fatalf("migration directory %q is unavailable: %v", dir, err)
	}
	return dir
}

func migrationStatusContains(t *testing.T, statuses []migrations.MigrationStatus, version int64, state string) bool {
	t.Helper()

	for _, status := range statuses {
		if status.Source.Version == version && status.State == state {
			return true
		}
	}
	return false
}

func singleResponseEnvelope(t *testing.T, responses [][]byte) *protocolv1.Envelope {
	t.Helper()

	if len(responses) != 1 {
		t.Fatalf("responses len = %d, want 1", len(responses))
	}
	return mustUnmarshalFrameEnvelope(t, responses[0])
}
