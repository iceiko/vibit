package bootstrap

import (
	"context"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

func TestInventoryDispatcherUsesUnitOfWorkRepositoryForCommands(t *testing.T) {
	commandRepository := inventory.NewMemoryRepository()
	queryRepository := inventory.NewMemoryRepository()
	provider := &recordingInventoryRepositoryProvider{
		commandRepository: commandRepository,
		queryRepository:   queryRepository,
	}
	dispatcher, err := NewInventoryDispatcher(testInventoryOptions(provider))
	if err != nil {
		t.Fatalf("NewInventoryDispatcher() error = %v, want nil", err)
	}
	transactional := app.TransactionalDispatcher{
		Dispatcher: dispatcher,
		Runner:     tx.NoopRunner{},
	}

	result, err := transactional.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     inventory.GrantItemRoute(),
		Payload: inventory.GrantItemRequest{
			PlayerID:    "player-1",
			ItemID:      "item-1",
			Quantity:    3,
			Reason:      "test",
			RequestedBy: "maintainer",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch(GrantItem) error = %v, want nil", err)
	}
	if result.PayloadType != "inventory.GrantItemResponse" {
		t.Fatalf("GrantItem PayloadType = %q, want inventory.GrantItemResponse", result.PayloadType)
	}
	if provider.commandCalls != 1 {
		t.Fatalf("command provider calls = %d, want 1", provider.commandCalls)
	}
	if provider.queryCalls != 0 {
		t.Fatalf("query provider calls = %d, want 0", provider.queryCalls)
	}
	if provider.lastCommandUnit == nil {
		t.Fatal("command provider did not receive a unit of work")
	}

	items, err := commandRepository.GetInventory(context.Background(), "player-1")
	if err != nil {
		t.Fatalf("command repository GetInventory() error = %v, want nil", err)
	}
	if len(items) != 1 || items[0].ItemID != "item-1" || items[0].Quantity != 3 {
		t.Fatalf("command repository items = %#v, want item-1 quantity 3", items)
	}
	queryItems, err := queryRepository.GetInventory(context.Background(), "player-1")
	if err != nil {
		t.Fatalf("query repository GetInventory() error = %v, want nil", err)
	}
	if len(queryItems) != 0 {
		t.Fatalf("query repository items = %#v, want empty", queryItems)
	}
}

func TestInventoryDispatcherUsesQueryRepositoryWithoutUnitOfWork(t *testing.T) {
	queryRepository := inventory.NewMemoryRepository()
	lock, err := queryRepository.LockInventoryForMutation(context.Background(), "player-1")
	if err != nil {
		t.Fatalf("LockInventoryForMutation() error = %v, want nil", err)
	}
	if _, err := lock.GrantItem(context.Background(), inventory.GrantItemMutation{
		EventID:    "event-1",
		OccurredAt: testClockTime,
		PlayerID:   "player-1",
		ItemID:     "item-1",
		Quantity:   2,
		Reason:     "seed",
	}); err != nil {
		t.Fatalf("GrantItem() seed error = %v, want nil", err)
	}
	lock.Release()

	provider := &recordingInventoryRepositoryProvider{
		commandRepository: inventory.NewMemoryRepository(),
		queryRepository:   queryRepository,
	}
	dispatcher, err := NewInventoryDispatcher(testInventoryOptions(provider))
	if err != nil {
		t.Fatalf("NewInventoryDispatcher() error = %v, want nil", err)
	}

	result, err := dispatcher.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     inventory.GetInventoryRoute(),
		Payload: inventory.GetInventoryRequest{
			PlayerID:    "player-1",
			RequestedBy: "maintainer",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch(GetInventory) error = %v, want nil", err)
	}
	response, ok := result.Payload.(inventory.GetInventoryResponse)
	if !ok {
		t.Fatalf("GetInventory payload type = %T, want inventory.GetInventoryResponse", result.Payload)
	}
	if len(response.Items) != 1 || response.Items[0].ItemID != "item-1" || response.Items[0].Quantity != 2 {
		t.Fatalf("GetInventory items = %#v, want item-1 quantity 2", response.Items)
	}
	if provider.commandCalls != 0 {
		t.Fatalf("command provider calls = %d, want 0", provider.commandCalls)
	}
	if provider.queryCalls != 1 {
		t.Fatalf("query provider calls = %d, want 1", provider.queryCalls)
	}
}

func TestInventoryDispatcherPassesRequestIdentityToPermissionPolicy(t *testing.T) {
	repository := inventory.NewMemoryRepository()
	provider := &recordingInventoryRepositoryProvider{
		commandRepository: repository,
		queryRepository:   repository,
	}
	policy := &recordingInventoryPermissionPolicy{grantAllowed: true, readAllowed: true}
	dispatcher, err := NewInventoryDispatcher(InventoryOptions{
		Repositories:     provider,
		PermissionPolicy: policy,
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 256},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "test-event"},
		Clock:            fixedInventoryClock{},
	})
	if err != nil {
		t.Fatalf("NewInventoryDispatcher() error = %v, want nil", err)
	}

	identity := app.ValidatedPlayerIdentity("player-1", app.Session{
		SessionID: "session-1",
		PlayerID:  "player-1",
	})
	_, err = dispatcher.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     inventory.GrantItemRoute(),
		Session:   app.Session{SessionID: "session-1", PlayerID: "player-1"},
		Identity:  identity,
		Payload: inventory.GrantItemRequest{
			PlayerID:    "player-1",
			ItemID:      "item-1",
			Quantity:    1,
			Reason:      "test",
			RequestedBy: "actor-1",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch(GrantItem) error = %v, want nil", err)
	}

	if policy.lastGrant.Identity != identity {
		t.Fatalf("grant Identity = %#v, want %#v", policy.lastGrant.Identity, identity)
	}
	if policy.lastGrant.Permission != inventory.PermissionGrantItem {
		t.Fatalf("grant Permission = %q, want %q", policy.lastGrant.Permission, inventory.PermissionGrantItem)
	}

	_, err = dispatcher.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-2",
		Route:     inventory.GetInventoryRoute(),
		Session:   app.Session{SessionID: "session-1", PlayerID: "player-1"},
		Identity:  identity,
		Payload: inventory.GetInventoryRequest{
			PlayerID:    "player-1",
			RequestedBy: "actor-1",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch(GetInventory) error = %v, want nil", err)
	}
	if policy.lastRead.Identity != identity {
		t.Fatalf("read Identity = %#v, want %#v", policy.lastRead.Identity, identity)
	}
	if policy.lastRead.Permission != inventory.PermissionRead {
		t.Fatalf("read Permission = %q, want %q", policy.lastRead.Permission, inventory.PermissionRead)
	}
}

func TestPostgresInventoryRepositoryProviderRequiresRepositoryFactoryForCommands(t *testing.T) {
	provider := PostgresInventoryRepositoryProvider{
		QueryRepository: inventory.NewMemoryRepository(),
	}

	_, err := provider.ForCommand(context.Background(), tx.NoopUnitOfWork{})
	if err == nil {
		t.Fatal("ForCommand() error = nil, want missing repository factory error")
	}
}

func TestPostgresInventoryRepositoryProviderUsesUnitOfWorkFactoryForCommands(t *testing.T) {
	repository := inventory.NewMemoryRepository()
	unit := fakeInventoryRepositoryFactoryUnit{repository: repository}
	provider := PostgresInventoryRepositoryProvider{
		QueryRepository: inventory.NewMemoryRepository(),
	}

	got, err := provider.ForCommand(context.Background(), unit)
	if err != nil {
		t.Fatalf("ForCommand() error = %v, want nil", err)
	}
	if got != repository {
		t.Fatalf("ForCommand() repository = %#v, want unit-of-work repository", got)
	}
}

func TestPostgresInventoryRepositoryProviderUsesQueryRepository(t *testing.T) {
	repository := inventory.NewMemoryRepository()
	provider := PostgresInventoryRepositoryProvider{
		QueryRepository: repository,
	}

	got, err := provider.ForQuery(context.Background())
	if err != nil {
		t.Fatalf("ForQuery() error = %v, want nil", err)
	}
	if got != repository {
		t.Fatalf("ForQuery() repository = %#v, want configured repository", got)
	}
}

func testInventoryOptions(provider InventoryRepositoryProvider) InventoryOptions {
	return InventoryOptions{
		Repositories:     provider,
		PermissionPolicy: inventory.StaticPermissionPolicy{GrantAllowed: true, ReadAllowed: true},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 256},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "test-event"},
		Clock:            fixedInventoryClock{},
	}
}

type recordingInventoryRepositoryProvider struct {
	commandRepository inventory.Repository
	queryRepository   inventory.Repository
	commandCalls      int
	queryCalls        int
	lastCommandUnit   tx.UnitOfWork
}

type recordingInventoryPermissionPolicy struct {
	grantAllowed bool
	readAllowed  bool
	lastGrant    inventory.PermissionContext
	lastRead     inventory.PermissionContext
}

func (p *recordingInventoryPermissionPolicy) CanGrantItem(_ context.Context, ctx inventory.PermissionContext) (bool, error) {
	p.lastGrant = ctx
	return p.grantAllowed, nil
}

func (p *recordingInventoryPermissionPolicy) CanReadInventory(_ context.Context, ctx inventory.PermissionContext) (bool, error) {
	p.lastRead = ctx
	return p.readAllowed, nil
}

func (p *recordingInventoryRepositoryProvider) ForCommand(_ context.Context, unit tx.UnitOfWork) (inventory.Repository, error) {
	p.commandCalls += 1
	p.lastCommandUnit = unit
	return p.commandRepository, nil
}

func (p *recordingInventoryRepositoryProvider) ForQuery(context.Context) (inventory.Repository, error) {
	p.queryCalls += 1
	return p.queryRepository, nil
}

type fixedInventoryClock struct{}

var testClockTime = time.Date(2026, 5, 13, 10, 30, 0, 0, time.UTC)

func (fixedInventoryClock) Now() time.Time {
	return testClockTime
}

type fakeInventoryRepositoryFactoryUnit struct {
	repository inventory.Repository
}

func (u fakeInventoryRepositoryFactoryUnit) Context() context.Context {
	return context.Background()
}

func (u fakeInventoryRepositoryFactoryUnit) NewInventoryRepository() (inventory.Repository, error) {
	return u.repository, nil
}
