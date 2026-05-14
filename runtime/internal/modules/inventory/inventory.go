package inventory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
)

const (
	ModuleName = "inventory"

	CommandGrantItem  = "GrantItem"
	QueryGetInventory = "GetInventory"

	PermissionGrantItem = "inventory_grant_item"
	PermissionRead      = "inventory_read"

	EventItemGranted = "ItemGranted"
)

type ErrorCode string

const (
	ErrorCodeInvalidItemQuantity ErrorCode = "INVALID_ITEM_QUANTITY"
	ErrorCodeInventoryCapacity   ErrorCode = "INVENTORY_CAPACITY_EXCEEDED"
	ErrorCodeInventoryPermission ErrorCode = "INVENTORY_PERMISSION_DENIED"
)

type InventoryError struct {
	Code    ErrorCode
	Message string
}

func (e *InventoryError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return string(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

type Item struct {
	ItemID   string
	Quantity int64
}

type GrantItemRequest struct {
	PlayerID    string
	ItemID      string
	Quantity    int64
	Reason      string
	RequestedBy string
}

type GrantItemResponse struct {
	PlayerID    string
	ItemID      string
	Quantity    int64
	NewQuantity int64
	Event       string
}

type GetInventoryRequest struct {
	PlayerID    string
	RequestedBy string
}

type GetInventoryResponse struct {
	PlayerID string
	Items    []Item
}

type ItemGrantedEvent struct {
	EventID     string
	OccurredAt  time.Time
	PlayerID    string
	ItemID      string
	Quantity    int64
	NewQuantity int64
	Reason      string
}

type Repository interface {
	GetInventory(context.Context, string) ([]Item, error)
	LockInventoryForMutation(context.Context, string) (MutationLock, error)
}

type MutationRepository interface {
	GetInventory(context.Context, string) ([]Item, error)
	GrantItem(context.Context, GrantItemMutation) (Item, error)
}

type MutationLock interface {
	MutationRepository
	Release()
}

type GrantItemMutation struct {
	EventID    string
	OccurredAt time.Time
	PlayerID   string
	ItemID     string
	Quantity   int64
	Reason     string
}

type PermissionPolicy interface {
	CanGrantItem(context.Context, PermissionContext) (bool, error)
	CanReadInventory(context.Context, PermissionContext) (bool, error)
}

type PermissionContext struct {
	Permission  string
	RequestedBy string
	PlayerID    string
	Identity    app.RequestIdentity
}

type CapacityPolicy interface {
	CanGrantItem(context.Context, string, []Item, GrantItemRequest) (bool, error)
}

type EventIDGenerator interface {
	NewEventID() string
}

type Clock interface {
	Now() time.Time
}

type MaxUniqueItemsCapacityPolicy struct {
	MaxUniqueItems int
}

func (p MaxUniqueItemsCapacityPolicy) CanGrantItem(_ context.Context, _ string, current []Item, request GrantItemRequest) (bool, error) {
	if p.MaxUniqueItems <= 0 {
		return true, nil
	}

	seen := make(map[string]struct{}, len(current)+1)
	for _, item := range current {
		seen[item.ItemID] = struct{}{}
	}
	seen[request.ItemID] = struct{}{}

	return len(seen) <= p.MaxUniqueItems, nil
}

type IncrementingEventIDGenerator struct {
	Prefix string
	mu     sync.Mutex
	next   uint64
}

func (g *IncrementingEventIDGenerator) NewEventID() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.next += 1
	prefix := strings.TrimSpace(g.Prefix)
	if prefix == "" {
		prefix = "inventory-event"
	}
	return fmt.Sprintf("%s-%d", prefix, g.next)
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}

type Handlers struct {
	Repository       Repository
	PermissionPolicy PermissionPolicy
	CapacityPolicy   CapacityPolicy
	EventIDs         EventIDGenerator
	Clock            Clock
}

func (h Handlers) RegisterRoutes(dispatcher *app.Dispatcher) error {
	if dispatcher == nil {
		return errors.New("inventory: dispatcher is nil")
	}
	if err := dispatcher.Register(GrantItemRoute(), app.HandlerFunc(h.HandleGrantItemRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(GetInventoryRoute(), app.HandlerFunc(h.HandleGetInventoryRoute)); err != nil {
		return err
	}
	return nil
}

func GrantItemRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: ModuleName, Name: CommandGrantItem}
}

func GetInventoryRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: ModuleName, Name: QueryGetInventory}
}

func (h Handlers) HandleGrantItemRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	payload, ok := request.Payload.(GrantItemRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*GrantItemRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		return baseResult(request), errors.New("inventory: GrantItem payload must be inventory.GrantItemRequest")
	}

	response, events, err := h.GrantItem(ctx, payload, request.Identity)
	if err != nil {
		return inventoryErrorResult(request, err)
	}

	result := baseResult(request)
	result.PayloadType = "inventory.GrantItemResponse"
	result.Payload = response
	result.Events = make([]app.ApplicationEvent, 0, len(events))
	for _, event := range events {
		result.Events = append(result.Events, app.ApplicationEvent{
			Route:       app.RouteKey{Kind: app.MessageKindEvent, Module: ModuleName, Name: EventItemGranted},
			PayloadType: "inventory.ItemGrantedEvent",
			Payload:     event,
		})
	}
	return result, nil
}

func (h Handlers) HandleGetInventoryRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	payload, ok := request.Payload.(GetInventoryRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*GetInventoryRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		return baseResult(request), errors.New("inventory: GetInventory payload must be inventory.GetInventoryRequest")
	}

	response, err := h.GetInventory(ctx, payload, request.Identity)
	if err != nil {
		return inventoryErrorResult(request, err)
	}

	result := baseResult(request)
	result.PayloadType = "inventory.GetInventoryResponse"
	result.Payload = response
	return result, nil
}

func (h Handlers) GrantItem(ctx context.Context, request GrantItemRequest, identity ...app.RequestIdentity) (GrantItemResponse, []ItemGrantedEvent, error) {
	request = normalizeGrantItemRequest(request)
	if request.Quantity <= 0 {
		return GrantItemResponse{}, nil, inventoryError(ErrorCodeInvalidItemQuantity, "item quantity must be positive")
	}
	if request.PlayerID == "" {
		return GrantItemResponse{}, nil, errors.New("inventory: player_id is required")
	}
	if request.ItemID == "" {
		return GrantItemResponse{}, nil, errors.New("inventory: item_id is required")
	}
	if request.Reason == "" {
		return GrantItemResponse{}, nil, errors.New("inventory: reason is required")
	}
	if request.RequestedBy == "" {
		return GrantItemResponse{}, nil, errors.New("inventory: requested_by is required")
	}

	repository, permissionPolicy, capacityPolicy, eventIDs, clock, err := h.dependencies()
	if err != nil {
		return GrantItemResponse{}, nil, err
	}

	allowed, err := permissionPolicy.CanGrantItem(ctx, PermissionContext{
		Permission:  PermissionGrantItem,
		RequestedBy: request.RequestedBy,
		PlayerID:    request.PlayerID,
		Identity:    firstIdentity(identity),
	})
	if err != nil {
		return GrantItemResponse{}, nil, err
	}
	if !allowed {
		return GrantItemResponse{}, nil, inventoryError(ErrorCodeInventoryPermission, "actor cannot grant inventory items")
	}

	mutationLock, err := repository.LockInventoryForMutation(ctx, request.PlayerID)
	if err != nil {
		return GrantItemResponse{}, nil, err
	}
	if mutationLock == nil {
		return GrantItemResponse{}, nil, errors.New("inventory: mutation lock is required")
	}
	defer mutationLock.Release()

	current, err := mutationLock.GetInventory(ctx, request.PlayerID)
	if err != nil {
		return GrantItemResponse{}, nil, err
	}
	current = normalizeItems(current)

	hasCapacity, err := capacityPolicy.CanGrantItem(ctx, request.PlayerID, current, request)
	if err != nil {
		return GrantItemResponse{}, nil, err
	}
	if !hasCapacity {
		return GrantItemResponse{}, nil, inventoryError(ErrorCodeInventoryCapacity, "inventory capacity would be exceeded")
	}

	event := ItemGrantedEvent{
		EventID:    eventIDs.NewEventID(),
		OccurredAt: clock.Now().UTC(),
		PlayerID:   request.PlayerID,
		ItemID:     request.ItemID,
		Quantity:   request.Quantity,
		Reason:     request.Reason,
	}
	if strings.TrimSpace(event.EventID) == "" {
		return GrantItemResponse{}, nil, errors.New("inventory: event_id is required")
	}
	if event.OccurredAt.IsZero() {
		return GrantItemResponse{}, nil, errors.New("inventory: event occurred_at is required")
	}

	granted, err := mutationLock.GrantItem(ctx, GrantItemMutation{
		EventID:    event.EventID,
		OccurredAt: event.OccurredAt,
		PlayerID:   request.PlayerID,
		ItemID:     request.ItemID,
		Quantity:   request.Quantity,
		Reason:     request.Reason,
	})
	if err != nil {
		return GrantItemResponse{}, nil, err
	}
	granted.ItemID = strings.TrimSpace(granted.ItemID)
	if granted.ItemID == "" {
		return GrantItemResponse{}, nil, errors.New("inventory: repository returned empty item_id")
	}
	if granted.Quantity < request.Quantity {
		return GrantItemResponse{}, nil, errors.New("inventory: repository returned invalid new quantity")
	}
	event.ItemID = granted.ItemID
	event.NewQuantity = granted.Quantity

	response := GrantItemResponse{
		PlayerID:    request.PlayerID,
		ItemID:      granted.ItemID,
		Quantity:    request.Quantity,
		NewQuantity: granted.Quantity,
		Event:       EventItemGranted,
	}
	return response, []ItemGrantedEvent{event}, nil
}

func (h Handlers) GetInventory(ctx context.Context, request GetInventoryRequest, identity ...app.RequestIdentity) (GetInventoryResponse, error) {
	request = normalizeGetInventoryRequest(request)
	if request.PlayerID == "" {
		return GetInventoryResponse{}, errors.New("inventory: player_id is required")
	}
	if request.RequestedBy == "" {
		return GetInventoryResponse{}, errors.New("inventory: requested_by is required")
	}

	repository, permissionPolicy, err := h.readDependencies()
	if err != nil {
		return GetInventoryResponse{}, err
	}

	allowed, err := permissionPolicy.CanReadInventory(ctx, PermissionContext{
		Permission:  PermissionRead,
		RequestedBy: request.RequestedBy,
		PlayerID:    request.PlayerID,
		Identity:    firstIdentity(identity),
	})
	if err != nil {
		return GetInventoryResponse{}, err
	}
	if !allowed {
		return GetInventoryResponse{}, inventoryError(ErrorCodeInventoryPermission, "actor cannot read inventory")
	}

	items, err := repository.GetInventory(ctx, request.PlayerID)
	if err != nil {
		return GetInventoryResponse{}, err
	}

	return GetInventoryResponse{
		PlayerID: request.PlayerID,
		Items:    normalizeItems(items),
	}, nil
}

func (h Handlers) dependencies() (Repository, PermissionPolicy, CapacityPolicy, EventIDGenerator, Clock, error) {
	repository, permissionPolicy, err := h.readDependencies()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	capacityPolicy := h.CapacityPolicy
	if capacityPolicy == nil {
		return nil, nil, nil, nil, nil, errors.New("inventory: capacity policy is required")
	}

	eventIDs := h.EventIDs
	if eventIDs == nil {
		return nil, nil, nil, nil, nil, errors.New("inventory: event id generator is required")
	}

	clock := h.Clock
	if clock == nil {
		clock = SystemClock{}
	}

	return repository, permissionPolicy, capacityPolicy, eventIDs, clock, nil
}

func (h Handlers) readDependencies() (Repository, PermissionPolicy, error) {
	if h.Repository == nil {
		return nil, nil, errors.New("inventory: repository is required")
	}

	permissionPolicy := h.PermissionPolicy
	if permissionPolicy == nil {
		return nil, nil, errors.New("inventory: permission policy is required")
	}

	return h.Repository, permissionPolicy, nil
}

func normalizeGrantItemRequest(request GrantItemRequest) GrantItemRequest {
	request.PlayerID = strings.TrimSpace(request.PlayerID)
	request.ItemID = strings.TrimSpace(request.ItemID)
	request.Reason = strings.TrimSpace(request.Reason)
	request.RequestedBy = strings.TrimSpace(request.RequestedBy)
	return request
}

func normalizeGetInventoryRequest(request GetInventoryRequest) GetInventoryRequest {
	request.PlayerID = strings.TrimSpace(request.PlayerID)
	request.RequestedBy = strings.TrimSpace(request.RequestedBy)
	return request
}

func firstIdentity(identity []app.RequestIdentity) app.RequestIdentity {
	if len(identity) == 0 {
		return app.RequestIdentity{}
	}
	return identity[0]
}

func normalizeItems(items []Item) []Item {
	normalized := make([]Item, 0, len(items))
	for _, item := range items {
		item.ItemID = strings.TrimSpace(item.ItemID)
		if item.ItemID == "" || item.Quantity <= 0 {
			continue
		}
		normalized = append(normalized, item)
	}
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].ItemID < normalized[j].ItemID
	})
	return normalized
}

func inventoryError(code ErrorCode, message string) *InventoryError {
	return &InventoryError{Code: code, Message: message}
}

func inventoryErrorResult(request app.RouteRequest, err error) (app.ApplicationResult, error) {
	var inventoryErr *InventoryError
	if errors.As(err, &inventoryErr) {
		return applicationErrorResult(request, inventoryErr.Code, inventoryErr.Message)
	}

	result := baseResult(request)
	return result, err
}

func applicationErrorResult(request app.RouteRequest, code ErrorCode, message string) (app.ApplicationResult, error) {
	result := baseResult(request)
	appErr := &app.ApplicationError{
		Code:    app.ErrorCode(code),
		Message: message,
		Route:   request.Route,
	}
	result.Error = appErr
	return result, appErr
}

func baseResult(request app.RouteRequest) app.ApplicationResult {
	return app.ApplicationResult{
		RequestID: request.RequestID,
		Route:     request.Route,
		Target:    request.Target,
		Session:   request.Session,
		Identity:  request.Identity,
	}
}
