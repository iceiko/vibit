package connection

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type ConnectionID string
type ConnectionEpoch uint64
type PlayerID string
type RuntimeSessionID string
type AccessTokenRecordID string

type State string

const (
	StateOpenUnbound State = "open_unbound"
	StateBound       State = "bound"
	StateClosed      State = "closed"
	StateInvalidated State = "invalidated"
)

func (s State) IsActive() bool {
	return s == StateOpenUnbound || s == StateBound
}

type ActorKind string

const (
	ActorKindPlayer ActorKind = "player"
)

type Clock interface {
	Now() time.Time
}

type ErrorCode string

const (
	ErrorCodeConnectionInvalid      ErrorCode = "connection_invalid"
	ErrorCodeConnectionAlreadyOpen  ErrorCode = "connection_already_open"
	ErrorCodeConnectionNotFound     ErrorCode = "connection_not_found"
	ErrorCodeConnectionNotActive    ErrorCode = "connection_not_active"
	ErrorCodeIdentityInvalid        ErrorCode = "identity_invalid"
	ErrorCodeInvalidationInvalid    ErrorCode = "invalidation_invalid"
	ErrorCodeClockUnavailable       ErrorCode = "clock_unavailable"
	ErrorCodeUnsupportedActorKind   ErrorCode = "unsupported_actor_kind"
	ErrorCodeInternalInvariantError ErrorCode = "internal_invariant_error"
)

type RegistryError struct {
	Code ErrorCode
	Err  error
}

func (e *RegistryError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("connection registry: %s", e.Code)
}

func (e *RegistryError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *RegistryError) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

type Record struct {
	ConnectionID        ConnectionID
	ConnectionEpoch     ConnectionEpoch
	State               State
	ActorKind           ActorKind
	PlayerID            PlayerID
	RuntimeSessionID    RuntimeSessionID
	AccessTokenRecordID AccessTokenRecordID
	OpenedAt            time.Time
	BoundAt             *time.Time
	LastSeenAt          *time.Time
	ClosedAt            *time.Time
	CloseReasonClass    string
	InvalidatedAt       *time.Time
	InvalidationClass   string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type OpenConnection struct {
	ConnectionID    ConnectionID
	ConnectionEpoch ConnectionEpoch
	OpenedAt        time.Time
}

type BindIdentity struct {
	ConnectionID        ConnectionID
	ConnectionEpoch     ConnectionEpoch
	ActorKind           ActorKind
	PlayerID            PlayerID
	RuntimeSessionID    RuntimeSessionID
	AccessTokenRecordID AccessTokenRecordID
	ValidatedAt         time.Time
}

type MarkClosed struct {
	ConnectionID     ConnectionID
	ConnectionEpoch  ConnectionEpoch
	ClosedAt         time.Time
	CloseReasonClass string
}

type Invalidation struct {
	ConnectionID      ConnectionID
	ConnectionEpoch   ConnectionEpoch
	InvalidatedAt     time.Time
	InvalidationClass string
}

type InMemoryRegistry struct {
	mu      sync.RWMutex
	clock   Clock
	records map[recordKey]Record
}

type recordKey struct {
	connectionID    ConnectionID
	connectionEpoch ConnectionEpoch
}

func NewInMemoryRegistry(clock Clock) *InMemoryRegistry {
	return &InMemoryRegistry{
		clock:   clock,
		records: make(map[recordKey]Record),
	}
}

func (r *InMemoryRegistry) RegisterOpenConnection(ctx context.Context, command OpenConnection) (Record, error) {
	if err := ctxErr(ctx); err != nil {
		return Record{}, err
	}
	command, err := normalizeOpenConnection(command, r.now)
	if err != nil {
		return Record{}, err
	}

	key := recordKey{connectionID: command.ConnectionID, connectionEpoch: command.ConnectionEpoch}
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.records[key]; ok && existing.State.IsActive() {
		return Record{}, registryError(ErrorCodeConnectionAlreadyOpen, nil)
	}

	record := Record{
		ConnectionID:    command.ConnectionID,
		ConnectionEpoch: command.ConnectionEpoch,
		State:           StateOpenUnbound,
		OpenedAt:        command.OpenedAt,
		CreatedAt:       command.OpenedAt,
		UpdatedAt:       command.OpenedAt,
	}
	r.records[key] = record
	return copyRecord(record), nil
}

func (r *InMemoryRegistry) BindConnectionIdentity(ctx context.Context, command BindIdentity) (Record, error) {
	if err := ctxErr(ctx); err != nil {
		return Record{}, err
	}
	command, err := normalizeBindIdentity(command, r.now)
	if err != nil {
		return Record{}, err
	}

	key := recordKey{connectionID: command.ConnectionID, connectionEpoch: command.ConnectionEpoch}
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.records[key]
	if !ok {
		return Record{}, registryError(ErrorCodeConnectionNotFound, nil)
	}
	if !record.State.IsActive() {
		return Record{}, registryError(ErrorCodeConnectionNotActive, nil)
	}

	record.State = StateBound
	record.ActorKind = command.ActorKind
	record.PlayerID = command.PlayerID
	record.RuntimeSessionID = command.RuntimeSessionID
	record.AccessTokenRecordID = command.AccessTokenRecordID
	record.BoundAt = copyTime(command.ValidatedAt)
	record.LastSeenAt = copyTime(command.ValidatedAt)
	record.UpdatedAt = command.ValidatedAt
	r.records[key] = record
	return copyRecord(record), nil
}

func (r *InMemoryRegistry) MarkConnectionClosed(ctx context.Context, command MarkClosed) (Record, error) {
	if err := ctxErr(ctx); err != nil {
		return Record{}, err
	}
	command, err := normalizeMarkClosed(command, r.now)
	if err != nil {
		return Record{}, err
	}

	key := recordKey{connectionID: command.ConnectionID, connectionEpoch: command.ConnectionEpoch}
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.records[key]
	if !ok {
		return Record{}, registryError(ErrorCodeConnectionNotFound, nil)
	}
	if !record.State.IsActive() {
		return Record{}, registryError(ErrorCodeConnectionNotActive, nil)
	}

	record.State = StateClosed
	record.ClosedAt = copyTime(command.ClosedAt)
	record.CloseReasonClass = command.CloseReasonClass
	record.UpdatedAt = command.ClosedAt
	r.records[key] = record
	return copyRecord(record), nil
}

func (r *InMemoryRegistry) MarkConnectionInvalidated(ctx context.Context, command Invalidation) (Record, error) {
	if err := ctxErr(ctx); err != nil {
		return Record{}, err
	}
	command, err := normalizeInvalidation(command, r.now)
	if err != nil {
		return Record{}, err
	}

	key := recordKey{connectionID: command.ConnectionID, connectionEpoch: command.ConnectionEpoch}
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.records[key]
	if !ok {
		return Record{}, registryError(ErrorCodeConnectionNotFound, nil)
	}
	if !record.State.IsActive() {
		return Record{}, registryError(ErrorCodeConnectionNotActive, nil)
	}

	record.State = StateInvalidated
	record.InvalidatedAt = copyTime(command.InvalidatedAt)
	record.InvalidationClass = command.InvalidationClass
	record.UpdatedAt = command.InvalidatedAt
	r.records[key] = record
	return copyRecord(record), nil
}

func (r *InMemoryRegistry) FindConnectionByID(ctx context.Context, connectionID ConnectionID, connectionEpoch ConnectionEpoch) (Record, bool) {
	if err := ctxErr(ctx); err != nil {
		return Record{}, false
	}
	connectionID = ConnectionID(strings.TrimSpace(string(connectionID)))
	if connectionID == "" || connectionEpoch == 0 {
		return Record{}, false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.records[recordKey{connectionID: connectionID, connectionEpoch: connectionEpoch}]
	if !ok {
		return Record{}, false
	}
	return copyRecord(record), true
}

func (r *InMemoryRegistry) ListConnectionsByPlayerID(ctx context.Context, playerID PlayerID) []Record {
	playerID = PlayerID(strings.TrimSpace(string(playerID)))
	if ctxErr(ctx) != nil || playerID == "" {
		return nil
	}
	return r.listActive(func(record Record) bool {
		return record.State == StateBound && record.PlayerID == playerID
	})
}

func (r *InMemoryRegistry) ListConnectionsByRuntimeSessionID(ctx context.Context, sessionID RuntimeSessionID) []Record {
	sessionID = RuntimeSessionID(strings.TrimSpace(string(sessionID)))
	if ctxErr(ctx) != nil || sessionID == "" {
		return nil
	}
	return r.listActive(func(record Record) bool {
		return record.State == StateBound && record.RuntimeSessionID == sessionID
	})
}

func (r *InMemoryRegistry) ListConnectionsByAccessTokenRecordID(ctx context.Context, tokenRecordID AccessTokenRecordID) []Record {
	tokenRecordID = AccessTokenRecordID(strings.TrimSpace(string(tokenRecordID)))
	if ctxErr(ctx) != nil || tokenRecordID == "" {
		return nil
	}
	return r.listActive(func(record Record) bool {
		return record.State == StateBound && record.AccessTokenRecordID == tokenRecordID
	})
}

func (r *InMemoryRegistry) listActive(match func(Record) bool) []Record {
	r.mu.RLock()
	defer r.mu.RUnlock()

	records := make([]Record, 0)
	for _, record := range r.records {
		if record.State.IsActive() && match(record) {
			records = append(records, copyRecord(record))
		}
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].ConnectionID == records[j].ConnectionID {
			return records[i].ConnectionEpoch < records[j].ConnectionEpoch
		}
		return records[i].ConnectionID < records[j].ConnectionID
	})
	return records
}

func normalizeOpenConnection(command OpenConnection, now func() (time.Time, error)) (OpenConnection, error) {
	var err error
	command.ConnectionID, err = normalizeConnectionID(command.ConnectionID)
	if err != nil {
		return OpenConnection{}, err
	}
	if command.ConnectionEpoch == 0 {
		return OpenConnection{}, registryError(ErrorCodeConnectionInvalid, nil)
	}
	command.OpenedAt, err = normalizeOptionalObservedAt(command.OpenedAt, now)
	if err != nil {
		return OpenConnection{}, err
	}
	return command, nil
}

func normalizeBindIdentity(command BindIdentity, now func() (time.Time, error)) (BindIdentity, error) {
	var err error
	command.ConnectionID, err = normalizeConnectionID(command.ConnectionID)
	if err != nil {
		return BindIdentity{}, err
	}
	if command.ConnectionEpoch == 0 {
		return BindIdentity{}, registryError(ErrorCodeConnectionInvalid, nil)
	}
	if command.ActorKind == "" {
		command.ActorKind = ActorKindPlayer
	}
	if command.ActorKind != ActorKindPlayer {
		return BindIdentity{}, registryError(ErrorCodeUnsupportedActorKind, nil)
	}
	command.PlayerID = PlayerID(strings.TrimSpace(string(command.PlayerID)))
	if command.PlayerID == "" {
		return BindIdentity{}, registryError(ErrorCodeIdentityInvalid, nil)
	}
	command.RuntimeSessionID = RuntimeSessionID(strings.TrimSpace(string(command.RuntimeSessionID)))
	command.AccessTokenRecordID = AccessTokenRecordID(strings.TrimSpace(string(command.AccessTokenRecordID)))
	command.ValidatedAt, err = normalizeOptionalObservedAt(command.ValidatedAt, now)
	if err != nil {
		return BindIdentity{}, err
	}
	return command, nil
}

func normalizeMarkClosed(command MarkClosed, now func() (time.Time, error)) (MarkClosed, error) {
	var err error
	command.ConnectionID, err = normalizeConnectionID(command.ConnectionID)
	if err != nil {
		return MarkClosed{}, err
	}
	if command.ConnectionEpoch == 0 {
		return MarkClosed{}, registryError(ErrorCodeConnectionInvalid, nil)
	}
	command.CloseReasonClass = normalizeReasonClass(command.CloseReasonClass)
	command.ClosedAt, err = normalizeOptionalObservedAt(command.ClosedAt, now)
	if err != nil {
		return MarkClosed{}, err
	}
	return command, nil
}

func normalizeInvalidation(command Invalidation, now func() (time.Time, error)) (Invalidation, error) {
	var err error
	command.ConnectionID, err = normalizeConnectionID(command.ConnectionID)
	if err != nil {
		return Invalidation{}, err
	}
	if command.ConnectionEpoch == 0 {
		return Invalidation{}, registryError(ErrorCodeConnectionInvalid, nil)
	}
	command.InvalidationClass = normalizeReasonClass(command.InvalidationClass)
	if command.InvalidationClass == "" {
		return Invalidation{}, registryError(ErrorCodeInvalidationInvalid, nil)
	}
	command.InvalidatedAt, err = normalizeOptionalObservedAt(command.InvalidatedAt, now)
	if err != nil {
		return Invalidation{}, err
	}
	return command, nil
}

func normalizeConnectionID(connectionID ConnectionID) (ConnectionID, error) {
	connectionID = ConnectionID(strings.TrimSpace(string(connectionID)))
	if connectionID == "" {
		return "", registryError(ErrorCodeConnectionInvalid, nil)
	}
	return connectionID, nil
}

func normalizeOptionalObservedAt(value time.Time, now func() (time.Time, error)) (time.Time, error) {
	if value.IsZero() {
		return now()
	}
	return value.UTC(), nil
}

func normalizeReasonClass(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.Join(strings.Fields(value), "_")
}

func (r *InMemoryRegistry) now() (time.Time, error) {
	if r.clock == nil {
		return time.Now().UTC(), nil
	}
	value := r.clock.Now()
	if value.IsZero() {
		return time.Time{}, registryError(ErrorCodeClockUnavailable, nil)
	}
	return value.UTC(), nil
}

func copyRecord(record Record) Record {
	record.BoundAt = copyTimeValue(record.BoundAt)
	record.LastSeenAt = copyTimeValue(record.LastSeenAt)
	record.ClosedAt = copyTimeValue(record.ClosedAt)
	record.InvalidatedAt = copyTimeValue(record.InvalidatedAt)
	return record
}

func copyTimeValue(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := value.UTC()
	return &copied
}

func copyTime(value time.Time) *time.Time {
	value = value.UTC()
	return &value
}

func registryError(code ErrorCode, err error) *RegistryError {
	return &RegistryError{Code: code, Err: err}
}

func ctxErr(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}
