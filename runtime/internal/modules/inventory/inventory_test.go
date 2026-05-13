package inventory

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
)

func TestGrantItemRecordsItemAndEmitsEvent(t *testing.T) {
	repository := newMemoryRepository()
	handlers := testHandlers(repository)

	response, events, err := handlers.GrantItem(context.Background(), GrantItemRequest{
		PlayerID:    " player-1 ",
		ItemID:      " item-1 ",
		Quantity:    3,
		Reason:      " reward ",
		RequestedBy: " admin-1 ",
	})
	if err != nil {
		t.Fatalf("GrantItem() error = %v, want nil", err)
	}

	if response.PlayerID != "player-1" || response.ItemID != "item-1" || response.Quantity != 3 || response.NewQuantity != 3 {
		t.Fatalf("GrantItem() response = %#v, want normalized grant response", response)
	}
	if response.Event != EventItemGranted {
		t.Fatalf("GrantItem() response event = %q, want %q", response.Event, EventItemGranted)
	}

	items, err := repository.GetInventory(context.Background(), "player-1")
	if err != nil {
		t.Fatalf("GetInventory() repository error = %v", err)
	}
	if len(items) != 1 || items[0].ItemID != "item-1" || items[0].Quantity != 3 {
		t.Fatalf("repository items = %#v, want item-1 quantity 3", items)
	}

	if len(events) != 1 {
		t.Fatalf("events len = %d, want 1", len(events))
	}
	event := events[0]
	if event.EventID != "event-1" {
		t.Fatalf("event id = %q, want event-1", event.EventID)
	}
	if event.OccurredAt != fixedTime {
		t.Fatalf("event occurred_at = %s, want %s", event.OccurredAt, fixedTime)
	}
	if event.PlayerID != "player-1" || event.ItemID != "item-1" || event.Quantity != 3 || event.NewQuantity != 3 || event.Reason != "reward" {
		t.Fatalf("event = %#v, want normalized ItemGranted event", event)
	}
	if repository.lastGrantMutation.EventID != event.EventID {
		t.Fatalf("mutation event id = %q, want %q", repository.lastGrantMutation.EventID, event.EventID)
	}
	if repository.lastGrantMutation.OccurredAt != event.OccurredAt {
		t.Fatalf("mutation occurred_at = %s, want %s", repository.lastGrantMutation.OccurredAt, event.OccurredAt)
	}
	if repository.lastGrantMutation.Reason != event.Reason {
		t.Fatalf("mutation reason = %q, want %q", repository.lastGrantMutation.Reason, event.Reason)
	}
}

func TestGrantItemAccumulatesExistingItemQuantity(t *testing.T) {
	repository := newMemoryRepository()
	handlers := testHandlers(repository)

	_, _, err := handlers.GrantItem(context.Background(), GrantItemRequest{
		PlayerID:    "player-1",
		ItemID:      "item-1",
		Quantity:    3,
		Reason:      "first",
		RequestedBy: "admin-1",
	})
	if err != nil {
		t.Fatalf("GrantItem() first error = %v, want nil", err)
	}

	response, events, err := handlers.GrantItem(context.Background(), GrantItemRequest{
		PlayerID:    "player-1",
		ItemID:      "item-1",
		Quantity:    2,
		Reason:      "second",
		RequestedBy: "admin-1",
	})
	if err != nil {
		t.Fatalf("GrantItem() second error = %v, want nil", err)
	}
	if response.NewQuantity != 5 {
		t.Fatalf("NewQuantity = %d, want 5", response.NewQuantity)
	}
	if len(events) != 1 || events[0].NewQuantity != 5 {
		t.Fatalf("events = %#v, want one event with new quantity 5", events)
	}
}

func TestGrantItemUsesMutationLockBeforeReadingAndGranting(t *testing.T) {
	repository := newMemoryRepository()
	handlers := testHandlers(repository)

	_, _, err := handlers.GrantItem(context.Background(), GrantItemRequest{
		PlayerID:    "player-1",
		ItemID:      "item-1",
		Quantity:    3,
		Reason:      "test",
		RequestedBy: "admin-1",
	})
	if err != nil {
		t.Fatalf("GrantItem() error = %v, want nil", err)
	}

	got := strings.Join(repository.operations, ",")
	want := "lock:player-1,locked_get:player-1,locked_grant:player-1:item-1,release:player-1"
	if got != want {
		t.Fatalf("repository operations = %s, want %s", got, want)
	}
}

func TestGrantItemRejectsInvalidQuantityBeforeMutation(t *testing.T) {
	repository := newMemoryRepository()
	handlers := testHandlers(repository)

	_, _, err := handlers.GrantItem(context.Background(), GrantItemRequest{
		PlayerID:    "player-1",
		ItemID:      "item-1",
		Quantity:    0,
		Reason:      "test",
		RequestedBy: "admin-1",
	})
	assertInventoryErrorCode(t, err, ErrorCodeInvalidItemQuantity)

	if repository.grantCalls != 0 {
		t.Fatalf("repository grant calls = %d, want 0", repository.grantCalls)
	}
	if repository.lockCalls != 0 {
		t.Fatalf("repository lock calls = %d, want 0", repository.lockCalls)
	}
}

func TestGrantItemRejectsCapacityBeforeMutation(t *testing.T) {
	repository := newMemoryRepository()
	repository.items["player-1"] = map[string]int64{"item-1": 1}
	handlers := testHandlers(repository)
	handlers.CapacityPolicy = MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 1}

	_, _, err := handlers.GrantItem(context.Background(), GrantItemRequest{
		PlayerID:    "player-1",
		ItemID:      "item-2",
		Quantity:    1,
		Reason:      "test",
		RequestedBy: "admin-1",
	})
	assertInventoryErrorCode(t, err, ErrorCodeInventoryCapacity)

	if repository.grantCalls != 0 {
		t.Fatalf("repository grant calls = %d, want 0", repository.grantCalls)
	}
	if repository.lockCalls != 1 {
		t.Fatalf("repository lock calls = %d, want 1", repository.lockCalls)
	}
	if repository.releaseCalls != 1 {
		t.Fatalf("repository release calls = %d, want 1", repository.releaseCalls)
	}
}

func TestGrantItemRejectsPermissionBeforeMutation(t *testing.T) {
	repository := newMemoryRepository()
	handlers := testHandlers(repository)
	handlers.PermissionPolicy = staticPermissionPolicy{grantAllowed: false, readAllowed: true}

	_, _, err := handlers.GrantItem(context.Background(), GrantItemRequest{
		PlayerID:    "player-1",
		ItemID:      "item-1",
		Quantity:    1,
		Reason:      "test",
		RequestedBy: "admin-1",
	})
	assertInventoryErrorCode(t, err, ErrorCodeInventoryPermission)

	if repository.grantCalls != 0 {
		t.Fatalf("repository grant calls = %d, want 0", repository.grantCalls)
	}
	if repository.lockCalls != 0 {
		t.Fatalf("repository lock calls = %d, want 0", repository.lockCalls)
	}
}

func TestGetInventoryReadsWithoutMutation(t *testing.T) {
	repository := newMemoryRepository()
	repository.items["player-1"] = map[string]int64{
		"item-2": 2,
		"item-1": 1,
	}
	handlers := testHandlers(repository)

	response, err := handlers.GetInventory(context.Background(), GetInventoryRequest{
		PlayerID:    " player-1 ",
		RequestedBy: " player-1 ",
	})
	if err != nil {
		t.Fatalf("GetInventory() error = %v, want nil", err)
	}

	if response.PlayerID != "player-1" {
		t.Fatalf("PlayerID = %q, want player-1", response.PlayerID)
	}
	if len(response.Items) != 2 {
		t.Fatalf("items len = %d, want 2", len(response.Items))
	}
	if response.Items[0].ItemID != "item-1" || response.Items[1].ItemID != "item-2" {
		t.Fatalf("items = %#v, want sorted item ids", response.Items)
	}
	if repository.grantCalls != 0 {
		t.Fatalf("repository grant calls = %d, want 0", repository.grantCalls)
	}
	if repository.lockCalls != 0 {
		t.Fatalf("repository lock calls = %d, want 0", repository.lockCalls)
	}
}

func TestGetInventoryRejectsPermission(t *testing.T) {
	repository := newMemoryRepository()
	handlers := testHandlers(repository)
	handlers.PermissionPolicy = staticPermissionPolicy{grantAllowed: true, readAllowed: false}

	_, err := handlers.GetInventory(context.Background(), GetInventoryRequest{
		PlayerID:    "player-1",
		RequestedBy: "user-2",
	})
	assertInventoryErrorCode(t, err, ErrorCodeInventoryPermission)
}

func TestRegisterRoutesDispatchesInventoryHandlers(t *testing.T) {
	repository := newMemoryRepository()
	handlers := testHandlers(repository)
	dispatcher := app.NewDispatcher()

	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}

	request := app.RouteRequest{
		RequestID: "request-1",
		Route:     GrantItemRoute(),
		Target:    app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		Session:   app.Session{SessionID: "session-1", PlayerID: "player-1"},
		Payload: GrantItemRequest{
			PlayerID:    "player-1",
			ItemID:      "item-1",
			Quantity:    2,
			Reason:      "test",
			RequestedBy: "admin-1",
		},
	}

	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch() grant error = %v, want nil", err)
	}
	if result.RequestID != request.RequestID || result.Route != request.Route || result.Target != request.Target || result.Session != request.Session {
		t.Fatalf("result metadata = %#v, want request metadata", result)
	}
	response, ok := result.Payload.(GrantItemResponse)
	if !ok {
		t.Fatalf("result payload = %T, want GrantItemResponse", result.Payload)
	}
	if response.NewQuantity != 2 {
		t.Fatalf("NewQuantity = %d, want 2", response.NewQuantity)
	}
	if len(result.Events) != 1 {
		t.Fatalf("events len = %d, want 1", len(result.Events))
	}
	if result.Events[0].Route != (app.RouteKey{Kind: app.MessageKindEvent, Module: ModuleName, Name: EventItemGranted}) {
		t.Fatalf("event route = %#v, want ItemGranted route", result.Events[0].Route)
	}

	queryResult, err := dispatcher.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-2",
		Route:     GetInventoryRoute(),
		Payload: GetInventoryRequest{
			PlayerID:    "player-1",
			RequestedBy: "player-1",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() query error = %v, want nil", err)
	}
	queryResponse, ok := queryResult.Payload.(GetInventoryResponse)
	if !ok {
		t.Fatalf("query payload = %T, want GetInventoryResponse", queryResult.Payload)
	}
	if len(queryResponse.Items) != 1 || queryResponse.Items[0].Quantity != 2 {
		t.Fatalf("query response = %#v, want granted item", queryResponse)
	}
}

func TestHandleGrantItemRouteRejectsInvalidPayload(t *testing.T) {
	handlers := testHandlers(newMemoryRepository())

	_, err := handlers.HandleGrantItemRoute(context.Background(), app.RouteRequest{
		Route:   GrantItemRoute(),
		Payload: GetInventoryRequest{},
	})
	if err == nil {
		t.Fatal("HandleGrantItemRoute() error = nil, want invalid payload error")
	}
	var appErr *app.ApplicationError
	if errors.As(err, &appErr) {
		t.Fatalf("error = %#v, want internal invalid payload error, not public app error", appErr)
	}
}

func testHandlers(repository *memoryRepository) Handlers {
	return Handlers{
		Repository:       repository,
		PermissionPolicy: staticPermissionPolicy{grantAllowed: true, readAllowed: true},
		CapacityPolicy:   MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 16},
		EventIDs:         &fixedEventIDs{},
		Clock:            fixedClock{},
	}
}

func assertInventoryErrorCode(t *testing.T, err error, code ErrorCode) {
	t.Helper()

	var inventoryErr *InventoryError
	if !errors.As(err, &inventoryErr) {
		t.Fatalf("error = %v, want *InventoryError", err)
	}
	if inventoryErr.Code != code {
		t.Fatalf("InventoryError.Code = %s, want %s", inventoryErr.Code, code)
	}
}

type staticPermissionPolicy struct {
	grantAllowed bool
	readAllowed  bool
}

func (p staticPermissionPolicy) CanGrantItem(context.Context, string, string) (bool, error) {
	return p.grantAllowed, nil
}

func (p staticPermissionPolicy) CanReadInventory(context.Context, string, string) (bool, error) {
	return p.readAllowed, nil
}

type memoryRepository struct {
	items             map[string]map[string]int64
	getCalls          int
	lockCalls         int
	releaseCalls      int
	grantCalls        int
	operations        []string
	lastGrantMutation GrantItemMutation
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{
		items: make(map[string]map[string]int64),
	}
}

func (r *memoryRepository) GetInventory(_ context.Context, playerID string) ([]Item, error) {
	r.getCalls += 1
	r.operations = append(r.operations, "get:"+playerID)
	return r.getInventory(playerID), nil
}

func (r *memoryRepository) LockInventoryForMutation(_ context.Context, playerID string) (MutationLock, error) {
	r.lockCalls += 1
	r.operations = append(r.operations, "lock:"+playerID)
	return &recordingMutationLock{
		repository: r,
		playerID:   playerID,
	}, nil
}

func (r *memoryRepository) getInventory(playerID string) []Item {
	playerItems := r.items[playerID]
	items := make([]Item, 0, len(playerItems))
	for itemID, quantity := range playerItems {
		items = append(items, Item{ItemID: itemID, Quantity: quantity})
	}
	return items
}

func (r *memoryRepository) grantItem(mutation GrantItemMutation) Item {
	r.grantCalls += 1
	r.lastGrantMutation = mutation
	if r.items[mutation.PlayerID] == nil {
		r.items[mutation.PlayerID] = make(map[string]int64)
	}
	r.items[mutation.PlayerID][mutation.ItemID] += mutation.Quantity
	return Item{
		ItemID:   mutation.ItemID,
		Quantity: r.items[mutation.PlayerID][mutation.ItemID],
	}
}

type recordingMutationLock struct {
	repository *memoryRepository
	playerID   string
	released   bool
}

func (l *recordingMutationLock) GetInventory(_ context.Context, playerID string) ([]Item, error) {
	if err := l.ensureUsable(playerID); err != nil {
		return nil, err
	}
	l.repository.operations = append(l.repository.operations, "locked_get:"+playerID)
	return l.repository.getInventory(playerID), nil
}

func (l *recordingMutationLock) GrantItem(_ context.Context, mutation GrantItemMutation) (Item, error) {
	if err := l.ensureUsable(mutation.PlayerID); err != nil {
		return Item{}, err
	}
	l.repository.operations = append(l.repository.operations, "locked_grant:"+mutation.PlayerID+":"+mutation.ItemID)
	return l.repository.grantItem(mutation), nil
}

func (l *recordingMutationLock) Release() {
	if l == nil || l.released {
		return
	}
	l.released = true
	l.repository.releaseCalls += 1
	l.repository.operations = append(l.repository.operations, "release:"+l.playerID)
}

func (l *recordingMutationLock) ensureUsable(playerID string) error {
	if l == nil || l.repository == nil {
		return errors.New("inventory test: mutation lock is not initialized")
	}
	if l.released {
		return errors.New("inventory test: mutation lock was released")
	}
	if playerID != l.playerID {
		return errors.New("inventory test: mutation lock player_id mismatch")
	}
	return nil
}

type fixedEventIDs struct {
	next int
}

func (g *fixedEventIDs) NewEventID() string {
	g.next += 1
	return fmt.Sprintf("event-%d", g.next)
}

var fixedTime = time.Date(2026, 5, 13, 9, 30, 0, 0, time.UTC)

type fixedClock struct{}

func (fixedClock) Now() time.Time {
	return fixedTime
}
