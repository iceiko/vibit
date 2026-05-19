package connection

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type CloseReasonClass string

const (
	CloseReasonTokenRevoked              CloseReasonClass = "token_revoked"
	CloseReasonLogoutPresentedToken      CloseReasonClass = "logout_presented_token"
	CloseReasonSessionRevoked            CloseReasonClass = "session_revoked"
	CloseReasonDuplicateConnectionPolicy CloseReasonClass = "duplicate_connection_policy"
	CloseReasonServerShutdownOrDrain     CloseReasonClass = "server_shutdown_or_drain"
	CloseReasonPolicyViolation           CloseReasonClass = "policy_violation"
	CloseReasonAdministrativeAction      CloseReasonClass = "administrative_action"
	CloseReasonProtocolError             CloseReasonClass = "protocol_error"
	CloseReasonIdleTimeout               CloseReasonClass = "idle_timeout"
	CloseReasonUnknownInternal           CloseReasonClass = "unknown_internal"
)

type CloseTargetKind string

const (
	CloseTargetConnectionID        CloseTargetKind = "connection_id_and_epoch"
	CloseTargetPlayerID            CloseTargetKind = "player_id"
	CloseTargetRuntimeSessionID    CloseTargetKind = "runtime_session_id"
	CloseTargetAccessTokenRecordID CloseTargetKind = "access_token_record_id"
)

type CloseTransportAction string

const (
	CloseTransportActionMarkInvalidatedOnly CloseTransportAction = "mark_invalidated_only"
)

type CloseRetryability string

const (
	CloseRetryabilityRetryable    CloseRetryability = "retryable"
	CloseRetryabilityNotRetryable CloseRetryability = "not_retryable"
	CloseRetryabilityUnknown      CloseRetryability = "unknown"
)

type ClosePublicVisibility string

const (
	ClosePublicVisibilitySilent                ClosePublicVisibility = "silent"
	ClosePublicVisibilityGenericDisconnect     ClosePublicVisibility = "generic_disconnect"
	ClosePublicVisibilityGenericReauthRequired ClosePublicVisibility = "generic_reauth_required"
)

type CloseOutcome string

const (
	CloseOutcomeInvalidated CloseOutcome = "invalidated"
	CloseOutcomeSkipped     CloseOutcome = "skipped"
)

type ClosePolicyErrorCode string

const (
	ClosePolicyErrorCodeTargetInvalid       ClosePolicyErrorCode = "target_invalid"
	ClosePolicyErrorCodeReasonInvalid       ClosePolicyErrorCode = "reason_invalid"
	ClosePolicyErrorCodeRegistryUnavailable ClosePolicyErrorCode = "registry_unavailable"
	ClosePolicyErrorCodeClockUnavailable    ClosePolicyErrorCode = "clock_unavailable"
	ClosePolicyErrorCodeInternal            ClosePolicyErrorCode = "internal"
)

type ClosePolicyError struct {
	Code ClosePolicyErrorCode
	Err  error
}

func (e *ClosePolicyError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("connection close policy: %s", e.Code)
}

func (e *ClosePolicyError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *ClosePolicyError) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

type ClosePolicy struct {
	registry *InMemoryRegistry
	clock    Clock
}

type ClosePolicyOption func(*ClosePolicy)

func WithClosePolicyClock(clock Clock) ClosePolicyOption {
	return func(policy *ClosePolicy) {
		policy.clock = clock
	}
}

func NewClosePolicy(registry *InMemoryRegistry, options ...ClosePolicyOption) *ClosePolicy {
	policy := &ClosePolicy{registry: registry}
	for _, option := range options {
		if option != nil {
			option(policy)
		}
	}
	return policy
}

type CloseTarget struct {
	Kind                CloseTargetKind
	ConnectionID        ConnectionID
	ConnectionEpoch     ConnectionEpoch
	PlayerID            PlayerID
	RuntimeSessionID    RuntimeSessionID
	AccessTokenRecordID AccessTokenRecordID
}

func TargetConnection(connectionID ConnectionID, connectionEpoch ConnectionEpoch) CloseTarget {
	return CloseTarget{
		Kind:            CloseTargetConnectionID,
		ConnectionID:    connectionID,
		ConnectionEpoch: connectionEpoch,
	}
}

func TargetPlayer(playerID PlayerID) CloseTarget {
	return CloseTarget{
		Kind:     CloseTargetPlayerID,
		PlayerID: playerID,
	}
}

func TargetRuntimeSession(sessionID RuntimeSessionID) CloseTarget {
	return CloseTarget{
		Kind:             CloseTargetRuntimeSessionID,
		RuntimeSessionID: sessionID,
	}
}

func TargetAccessTokenRecord(tokenRecordID AccessTokenRecordID) CloseTarget {
	return CloseTarget{
		Kind:                CloseTargetAccessTokenRecordID,
		AccessTokenRecordID: tokenRecordID,
	}
}

type CloseConnectionsCommand struct {
	Target           CloseTarget
	ReasonClass      CloseReasonClass
	PublicVisibility ClosePublicVisibility
	Retryability     CloseRetryability
	RequestedAt      time.Time
}

type ClosePolicyResult struct {
	Target           CloseTarget
	ReasonClass      CloseReasonClass
	TransportAction  CloseTransportAction
	PublicVisibility ClosePublicVisibility
	Retryability     CloseRetryability
	RequestedAt      time.Time
	Intents          []CloseIntent
	Skipped          []CloseSkipped
}

type CloseIntent struct {
	ConnectionID        ConnectionID
	ConnectionEpoch     ConnectionEpoch
	TargetKind          CloseTargetKind
	ActorKind           ActorKind
	PlayerID            PlayerID
	RuntimeSessionID    RuntimeSessionID
	AccessTokenRecordID AccessTokenRecordID
	ReasonClass         CloseReasonClass
	TransportAction     CloseTransportAction
	PublicVisibility    ClosePublicVisibility
	Retryability        CloseRetryability
	Outcome             CloseOutcome
	RecordedAt          time.Time
}

type CloseSkipped struct {
	ConnectionID    ConnectionID
	ConnectionEpoch ConnectionEpoch
	TargetKind      CloseTargetKind
	ReasonClass     CloseReasonClass
	Outcome         CloseOutcome
	FailureClass    string
	RecordedAt      time.Time
}

func (p *ClosePolicy) RequestClose(ctx context.Context, command CloseConnectionsCommand) (ClosePolicyResult, error) {
	if err := ctxErr(ctx); err != nil {
		return ClosePolicyResult{}, err
	}
	if p == nil || p.registry == nil {
		return ClosePolicyResult{}, closePolicyError(ClosePolicyErrorCodeRegistryUnavailable, nil)
	}
	command, err := normalizeCloseConnectionsCommand(command, p.now)
	if err != nil {
		return ClosePolicyResult{}, err
	}

	targetRecords, err := p.resolveTargetRecords(ctx, command.Target)
	if err != nil {
		return ClosePolicyResult{}, err
	}

	result := ClosePolicyResult{
		Target:           command.Target,
		ReasonClass:      command.ReasonClass,
		TransportAction:  CloseTransportActionMarkInvalidatedOnly,
		PublicVisibility: command.PublicVisibility,
		Retryability:     command.Retryability,
		RequestedAt:      command.RequestedAt,
		Intents:          make([]CloseIntent, 0, len(targetRecords)),
		Skipped:          make([]CloseSkipped, 0),
	}

	for _, record := range targetRecords {
		invalidated, err := p.registry.MarkConnectionInvalidated(ctx, Invalidation{
			ConnectionID:      record.ConnectionID,
			ConnectionEpoch:   record.ConnectionEpoch,
			InvalidatedAt:     command.RequestedAt,
			InvalidationClass: string(command.ReasonClass),
		})
		if err != nil {
			if hasRegistryErrorCode(err, ErrorCodeConnectionNotFound) || hasRegistryErrorCode(err, ErrorCodeConnectionNotActive) {
				result.Skipped = append(result.Skipped, CloseSkipped{
					ConnectionID:    record.ConnectionID,
					ConnectionEpoch: record.ConnectionEpoch,
					TargetKind:      command.Target.Kind,
					ReasonClass:     command.ReasonClass,
					Outcome:         CloseOutcomeSkipped,
					FailureClass:    "target_not_active",
					RecordedAt:      command.RequestedAt,
				})
				continue
			}
			return ClosePolicyResult{}, closePolicyError(ClosePolicyErrorCodeInternal, err)
		}

		result.Intents = append(result.Intents, CloseIntent{
			ConnectionID:        invalidated.ConnectionID,
			ConnectionEpoch:     invalidated.ConnectionEpoch,
			TargetKind:          command.Target.Kind,
			ActorKind:           invalidated.ActorKind,
			PlayerID:            invalidated.PlayerID,
			RuntimeSessionID:    invalidated.RuntimeSessionID,
			AccessTokenRecordID: invalidated.AccessTokenRecordID,
			ReasonClass:         command.ReasonClass,
			TransportAction:     CloseTransportActionMarkInvalidatedOnly,
			PublicVisibility:    command.PublicVisibility,
			Retryability:        command.Retryability,
			Outcome:             CloseOutcomeInvalidated,
			RecordedAt:          command.RequestedAt,
		})
	}

	return copyClosePolicyResult(result), nil
}

func (p *ClosePolicy) resolveTargetRecords(ctx context.Context, target CloseTarget) ([]Record, error) {
	switch target.Kind {
	case CloseTargetConnectionID:
		record, ok := p.registry.FindConnectionByID(ctx, target.ConnectionID, target.ConnectionEpoch)
		if !ok || record.State != StateBound {
			return nil, nil
		}
		return []Record{record}, nil
	case CloseTargetPlayerID:
		return p.registry.ListConnectionsByPlayerID(ctx, target.PlayerID), nil
	case CloseTargetRuntimeSessionID:
		return p.registry.ListConnectionsByRuntimeSessionID(ctx, target.RuntimeSessionID), nil
	case CloseTargetAccessTokenRecordID:
		return p.registry.ListConnectionsByAccessTokenRecordID(ctx, target.AccessTokenRecordID), nil
	default:
		return nil, closePolicyError(ClosePolicyErrorCodeTargetInvalid, nil)
	}
}

func normalizeCloseConnectionsCommand(command CloseConnectionsCommand, now func() (time.Time, error)) (CloseConnectionsCommand, error) {
	target, err := normalizeCloseTarget(command.Target)
	if err != nil {
		return CloseConnectionsCommand{}, err
	}
	command.Target = target

	reason := CloseReasonClass(normalizeReasonClass(string(command.ReasonClass)))
	if !isAllowedCloseReason(reason) {
		return CloseConnectionsCommand{}, closePolicyError(ClosePolicyErrorCodeReasonInvalid, nil)
	}
	command.ReasonClass = reason

	if command.PublicVisibility == "" {
		command.PublicVisibility = ClosePublicVisibilitySilent
	}
	if !isAllowedClosePublicVisibility(command.PublicVisibility) {
		return CloseConnectionsCommand{}, closePolicyError(ClosePolicyErrorCodeReasonInvalid, nil)
	}

	if command.Retryability == "" {
		command.Retryability = CloseRetryabilityUnknown
	}
	if !isAllowedCloseRetryability(command.Retryability) {
		return CloseConnectionsCommand{}, closePolicyError(ClosePolicyErrorCodeReasonInvalid, nil)
	}

	command.RequestedAt, err = normalizeOptionalObservedAt(command.RequestedAt, now)
	if err != nil {
		return CloseConnectionsCommand{}, closePolicyError(ClosePolicyErrorCodeClockUnavailable, err)
	}
	return command, nil
}

func normalizeCloseTarget(target CloseTarget) (CloseTarget, error) {
	switch target.Kind {
	case CloseTargetConnectionID:
		connectionID, err := normalizeConnectionID(target.ConnectionID)
		if err != nil {
			return CloseTarget{}, closePolicyError(ClosePolicyErrorCodeTargetInvalid, err)
		}
		if target.ConnectionEpoch == 0 {
			return CloseTarget{}, closePolicyError(ClosePolicyErrorCodeTargetInvalid, nil)
		}
		return CloseTarget{
			Kind:            CloseTargetConnectionID,
			ConnectionID:    connectionID,
			ConnectionEpoch: target.ConnectionEpoch,
		}, nil
	case CloseTargetPlayerID:
		playerID := PlayerID(strings.TrimSpace(string(target.PlayerID)))
		if playerID == "" {
			return CloseTarget{}, closePolicyError(ClosePolicyErrorCodeTargetInvalid, nil)
		}
		return CloseTarget{
			Kind:     CloseTargetPlayerID,
			PlayerID: playerID,
		}, nil
	case CloseTargetRuntimeSessionID:
		sessionID := RuntimeSessionID(strings.TrimSpace(string(target.RuntimeSessionID)))
		if sessionID == "" {
			return CloseTarget{}, closePolicyError(ClosePolicyErrorCodeTargetInvalid, nil)
		}
		return CloseTarget{
			Kind:             CloseTargetRuntimeSessionID,
			RuntimeSessionID: sessionID,
		}, nil
	case CloseTargetAccessTokenRecordID:
		tokenRecordID := AccessTokenRecordID(strings.TrimSpace(string(target.AccessTokenRecordID)))
		if tokenRecordID == "" {
			return CloseTarget{}, closePolicyError(ClosePolicyErrorCodeTargetInvalid, nil)
		}
		return CloseTarget{
			Kind:                CloseTargetAccessTokenRecordID,
			AccessTokenRecordID: tokenRecordID,
		}, nil
	default:
		return CloseTarget{}, closePolicyError(ClosePolicyErrorCodeTargetInvalid, nil)
	}
}

func isAllowedCloseReason(reason CloseReasonClass) bool {
	switch reason {
	case CloseReasonTokenRevoked,
		CloseReasonLogoutPresentedToken,
		CloseReasonSessionRevoked,
		CloseReasonDuplicateConnectionPolicy,
		CloseReasonServerShutdownOrDrain,
		CloseReasonPolicyViolation,
		CloseReasonAdministrativeAction,
		CloseReasonProtocolError,
		CloseReasonIdleTimeout,
		CloseReasonUnknownInternal:
		return true
	default:
		return false
	}
}

func isAllowedClosePublicVisibility(visibility ClosePublicVisibility) bool {
	switch visibility {
	case ClosePublicVisibilitySilent,
		ClosePublicVisibilityGenericDisconnect,
		ClosePublicVisibilityGenericReauthRequired:
		return true
	default:
		return false
	}
}

func isAllowedCloseRetryability(retryability CloseRetryability) bool {
	switch retryability {
	case CloseRetryabilityRetryable,
		CloseRetryabilityNotRetryable,
		CloseRetryabilityUnknown:
		return true
	default:
		return false
	}
}

func (p *ClosePolicy) now() (time.Time, error) {
	if p.clock != nil {
		value := p.clock.Now()
		if value.IsZero() {
			return time.Time{}, closePolicyError(ClosePolicyErrorCodeClockUnavailable, nil)
		}
		return value.UTC(), nil
	}
	if p.registry != nil {
		return p.registry.now()
	}
	return time.Time{}, closePolicyError(ClosePolicyErrorCodeRegistryUnavailable, nil)
}

func copyClosePolicyResult(result ClosePolicyResult) ClosePolicyResult {
	result.Intents = append([]CloseIntent(nil), result.Intents...)
	result.Skipped = append([]CloseSkipped(nil), result.Skipped...)
	sort.Slice(result.Intents, func(i, j int) bool {
		if result.Intents[i].ConnectionID == result.Intents[j].ConnectionID {
			return result.Intents[i].ConnectionEpoch < result.Intents[j].ConnectionEpoch
		}
		return result.Intents[i].ConnectionID < result.Intents[j].ConnectionID
	})
	sort.Slice(result.Skipped, func(i, j int) bool {
		if result.Skipped[i].ConnectionID == result.Skipped[j].ConnectionID {
			return result.Skipped[i].ConnectionEpoch < result.Skipped[j].ConnectionEpoch
		}
		return result.Skipped[i].ConnectionID < result.Skipped[j].ConnectionID
	})
	return result
}

func hasRegistryErrorCode(err error, code ErrorCode) bool {
	var registryErr *RegistryError
	return errors.As(err, &registryErr) && registryErr.Code == code
}

func closePolicyError(code ClosePolicyErrorCode, err error) *ClosePolicyError {
	return &ClosePolicyError{Code: code, Err: err}
}
