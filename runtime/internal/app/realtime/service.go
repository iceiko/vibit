package realtime

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/connection"
)

type MessageIntentKind string

const (
	MessageIntentServerNotice    MessageIntentKind = "server_notice"
	MessageIntentDomainEventPush MessageIntentKind = "domain_event_push"
	MessageIntentStreamMessage   MessageIntentKind = "stream_message"
	MessageIntentPresenceSignal  MessageIntentKind = "presence_signal"
)

type TargetKind string

const (
	TargetConnectionIDAndEpoch TargetKind = "connection_id_and_epoch"
	TargetPlayerCurrentConns   TargetKind = "player_current_connections"
	TargetStreamSubscribers    TargetKind = "stream_subscribers"
)

type DeliveryOutcome string

const (
	DeliveryOutcomeAccepted               DeliveryOutcome = "accepted"
	DeliveryOutcomeNoActiveRecipient      DeliveryOutcome = "no_active_recipient"
	DeliveryOutcomeRecipientNotAuthorized DeliveryOutcome = "recipient_not_authorized"
	DeliveryOutcomePayloadInvalid         DeliveryOutcome = "payload_invalid"
	DeliveryOutcomeDeliveryUnavailable    DeliveryOutcome = "delivery_unavailable"
)

type FailureClass string

const (
	FailureClassInvalidIntent         FailureClass = "invalid_intent"
	FailureClassForbidden             FailureClass = "forbidden"
	FailureClassNoActiveRecipient     FailureClass = "no_active_recipient"
	FailureClassDependencyUnavailable FailureClass = "dependency_unavailable"
)

type ErrorCode string

const (
	ErrorCodeIntentInvalid       ErrorCode = "intent_invalid"
	ErrorCodeRecipientForbidden  ErrorCode = "recipient_forbidden"
	ErrorCodeNoActiveRecipient   ErrorCode = "no_active_recipient"
	ErrorCodeDeliveryUnavailable ErrorCode = "delivery_unavailable"
)

type ServiceError struct {
	Code  ErrorCode
	Class FailureClass
	Err   error
}

func (e *ServiceError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("realtime service: %s", e.Code)
}

func (e *ServiceError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *ServiceError) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

type Registry interface {
	FindConnectionByID(context.Context, connection.ConnectionID, connection.ConnectionEpoch) (connection.Record, bool)
	ListConnectionsByPlayerID(context.Context, connection.PlayerID) []connection.Record
}

type Clock interface {
	Now() time.Time
}

type Service struct {
	registry Registry
	clock    Clock
}

type ServiceDependencies struct {
	Registry Registry
	Clock    Clock
}

func NewService(dependencies ServiceDependencies) (Service, error) {
	if isNilInterface(dependencies.Registry) {
		return Service{}, serviceError(ErrorCodeDeliveryUnavailable, FailureClassDependencyUnavailable, errMissingRegistry)
	}
	return Service{
		registry: dependencies.Registry,
		clock:    dependencies.Clock,
	}, nil
}

type MessageIntent struct {
	Kind           MessageIntentKind
	Target         Target
	SenderIdentity app.RequestIdentity
	PayloadType    string
	PayloadBytes   []byte
	MessageID      string
	CreatedAt      time.Time
}

type Target struct {
	Kind            TargetKind
	ConnectionID    connection.ConnectionID
	ConnectionEpoch connection.ConnectionEpoch
	PlayerID        connection.PlayerID
	StreamID        string
}

func TargetConnection(connectionID connection.ConnectionID, connectionEpoch connection.ConnectionEpoch) Target {
	return Target{
		Kind:            TargetConnectionIDAndEpoch,
		ConnectionID:    connectionID,
		ConnectionEpoch: connectionEpoch,
	}
}

func TargetPlayer(playerID connection.PlayerID) Target {
	return Target{
		Kind:     TargetPlayerCurrentConns,
		PlayerID: playerID,
	}
}

func TargetStream(streamID string) Target {
	return Target{
		Kind:     TargetStreamSubscribers,
		StreamID: streamID,
	}
}

type DeliveryResult struct {
	Outcome      DeliveryOutcome
	FailureClass FailureClass
	ErrorCode    ErrorCode
	IntentKind   MessageIntentKind
	Target       Target
	MessageID    string
	PayloadType  string
	AcceptedAt   time.Time
	Intents      []DeliveryIntent
}

type DeliveryIntent struct {
	IntentKind          MessageIntentKind
	ConnectionID        connection.ConnectionID
	ConnectionEpoch     connection.ConnectionEpoch
	ActorKind           connection.ActorKind
	PlayerID            connection.PlayerID
	RuntimeSessionID    connection.RuntimeSessionID
	AccessTokenRecordID connection.AccessTokenRecordID
	MessageID           string
	PayloadType         string
	PayloadBytes        []byte
	AcceptedAt          time.Time
}

func (s Service) AcceptServerMessage(ctx context.Context, intent MessageIntent) (DeliveryResult, error) {
	if err := ctx.Err(); err != nil {
		return DeliveryResult{}, err
	}
	if isNilInterface(s.registry) {
		return rejectedResult(intent, DeliveryOutcomeDeliveryUnavailable, FailureClassDependencyUnavailable, ErrorCodeDeliveryUnavailable),
			serviceError(ErrorCodeDeliveryUnavailable, FailureClassDependencyUnavailable, errMissingRegistry)
	}

	normalized, err := s.normalizeIntent(intent)
	if err != nil {
		result := rejectedResult(intent, DeliveryOutcomePayloadInvalid, FailureClassInvalidIntent, ErrorCodeIntentInvalid)
		var serviceErr *ServiceError
		if errors.As(err, &serviceErr) {
			switch serviceErr.Class {
			case FailureClassForbidden:
				result = rejectedResult(intent, DeliveryOutcomeRecipientNotAuthorized, FailureClassForbidden, ErrorCodeRecipientForbidden)
			case FailureClassDependencyUnavailable:
				result = rejectedResult(intent, DeliveryOutcomeDeliveryUnavailable, FailureClassDependencyUnavailable, ErrorCodeDeliveryUnavailable)
			}
		}
		return result, err
	}

	recipients, err := s.resolveRecipients(ctx, normalized)
	if err != nil {
		result := rejectedResult(normalized, DeliveryOutcomeRecipientNotAuthorized, FailureClassForbidden, ErrorCodeRecipientForbidden)
		var serviceErr *ServiceError
		if errors.As(err, &serviceErr) && serviceErr.Class == FailureClassDependencyUnavailable {
			result = rejectedResult(normalized, DeliveryOutcomeDeliveryUnavailable, FailureClassDependencyUnavailable, ErrorCodeDeliveryUnavailable)
		}
		return result, err
	}
	if len(recipients) == 0 {
		return rejectedResult(normalized, DeliveryOutcomeNoActiveRecipient, FailureClassNoActiveRecipient, ErrorCodeNoActiveRecipient),
			serviceError(ErrorCodeNoActiveRecipient, FailureClassNoActiveRecipient, nil)
	}

	result := DeliveryResult{
		Outcome:     DeliveryOutcomeAccepted,
		IntentKind:  normalized.Kind,
		Target:      normalized.Target,
		MessageID:   normalized.MessageID,
		PayloadType: normalized.PayloadType,
		AcceptedAt:  normalized.CreatedAt,
		Intents:     make([]DeliveryIntent, 0, len(recipients)),
	}
	for _, recipient := range recipients {
		result.Intents = append(result.Intents, DeliveryIntent{
			IntentKind:          normalized.Kind,
			ConnectionID:        recipient.ConnectionID,
			ConnectionEpoch:     recipient.ConnectionEpoch,
			ActorKind:           recipient.ActorKind,
			PlayerID:            recipient.PlayerID,
			RuntimeSessionID:    recipient.RuntimeSessionID,
			AccessTokenRecordID: recipient.AccessTokenRecordID,
			MessageID:           normalized.MessageID,
			PayloadType:         normalized.PayloadType,
			PayloadBytes:        append([]byte(nil), normalized.PayloadBytes...),
			AcceptedAt:          normalized.CreatedAt,
		})
	}
	return copyDeliveryResult(result), nil
}

func (s Service) normalizeIntent(intent MessageIntent) (MessageIntent, error) {
	intent.Kind = MessageIntentKind(strings.TrimSpace(string(intent.Kind)))
	if !allowedIntentKind(intent.Kind) {
		return MessageIntent{}, serviceError(ErrorCodeIntentInvalid, FailureClassInvalidIntent, nil)
	}
	intent.PayloadType = strings.TrimSpace(intent.PayloadType)
	if intent.PayloadType == "" || len(intent.PayloadBytes) == 0 {
		return MessageIntent{}, serviceError(ErrorCodeIntentInvalid, FailureClassInvalidIntent, nil)
	}
	intent.PayloadBytes = append([]byte(nil), intent.PayloadBytes...)
	intent.MessageID = strings.TrimSpace(intent.MessageID)
	intent.Target = normalizeTarget(intent.Target)
	if !targetHasRequiredFields(intent.Target) {
		return MessageIntent{}, serviceError(ErrorCodeIntentInvalid, FailureClassInvalidIntent, nil)
	}
	if intent.CreatedAt.IsZero() {
		intent.CreatedAt = s.now()
	} else {
		intent.CreatedAt = intent.CreatedAt.UTC()
	}
	if intent.CreatedAt.IsZero() {
		return MessageIntent{}, serviceError(ErrorCodeDeliveryUnavailable, FailureClassDependencyUnavailable, errClockUnavailable)
	}
	if !identityIsServerAuthorized(intent.SenderIdentity) {
		return MessageIntent{}, serviceError(ErrorCodeRecipientForbidden, FailureClassForbidden, nil)
	}
	return intent, nil
}

func (s Service) resolveRecipients(ctx context.Context, intent MessageIntent) ([]connection.Record, error) {
	switch intent.Target.Kind {
	case TargetConnectionIDAndEpoch:
		record, ok := s.registry.FindConnectionByID(ctx, intent.Target.ConnectionID, intent.Target.ConnectionEpoch)
		if !ok || record.State != connection.StateBound {
			return nil, nil
		}
		return []connection.Record{record}, nil
	case TargetPlayerCurrentConns:
		return activeBoundRecords(s.registry.ListConnectionsByPlayerID(ctx, intent.Target.PlayerID)), nil
	case TargetStreamSubscribers:
		return nil, serviceError(ErrorCodeRecipientForbidden, FailureClassForbidden, nil)
	default:
		return nil, serviceError(ErrorCodeIntentInvalid, FailureClassInvalidIntent, nil)
	}
}

func activeBoundRecords(records []connection.Record) []connection.Record {
	if len(records) == 0 {
		return nil
	}
	active := make([]connection.Record, 0, len(records))
	for _, record := range records {
		if record.State == connection.StateBound {
			active = append(active, record)
		}
	}
	return active
}

func identityIsServerAuthorized(identity app.RequestIdentity) bool {
	return identity.Status == app.IdentityValidationValidated &&
		(identity.ActorKind == app.ActorKindService || identity.ActorKind == app.ActorKindAdmin) &&
		strings.TrimSpace(identity.ActorID) != ""
}

func allowedIntentKind(kind MessageIntentKind) bool {
	switch kind {
	case MessageIntentServerNotice,
		MessageIntentDomainEventPush,
		MessageIntentStreamMessage,
		MessageIntentPresenceSignal:
		return true
	default:
		return false
	}
}

func normalizeTarget(target Target) Target {
	target.Kind = TargetKind(strings.TrimSpace(string(target.Kind)))
	target.ConnectionID = connection.ConnectionID(strings.TrimSpace(string(target.ConnectionID)))
	target.PlayerID = connection.PlayerID(strings.TrimSpace(string(target.PlayerID)))
	target.StreamID = strings.TrimSpace(target.StreamID)
	return target
}

func targetHasRequiredFields(target Target) bool {
	switch target.Kind {
	case TargetConnectionIDAndEpoch:
		return target.ConnectionID != "" && target.ConnectionEpoch != 0
	case TargetPlayerCurrentConns:
		return target.PlayerID != ""
	case TargetStreamSubscribers:
		return target.StreamID != ""
	default:
		return false
	}
}

func rejectedResult(intent MessageIntent, outcome DeliveryOutcome, class FailureClass, code ErrorCode) DeliveryResult {
	return DeliveryResult{
		Outcome:      outcome,
		FailureClass: class,
		ErrorCode:    code,
		IntentKind:   MessageIntentKind(strings.TrimSpace(string(intent.Kind))),
		Target:       normalizeTarget(intent.Target),
		MessageID:    strings.TrimSpace(intent.MessageID),
		PayloadType:  strings.TrimSpace(intent.PayloadType),
	}
}

func serviceError(code ErrorCode, class FailureClass, err error) *ServiceError {
	return &ServiceError{
		Code:  code,
		Class: class,
		Err:   err,
	}
}

func (s Service) now() time.Time {
	if s.clock == nil {
		return time.Now().UTC()
	}
	return s.clock.Now().UTC()
}

func copyDeliveryResult(result DeliveryResult) DeliveryResult {
	result.Intents = append([]DeliveryIntent(nil), result.Intents...)
	for index := range result.Intents {
		result.Intents[index].PayloadBytes = append([]byte(nil), result.Intents[index].PayloadBytes...)
	}
	return result
}

func isNilInterface(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

var (
	errMissingRegistry  = errors.New("missing realtime registry")
	errClockUnavailable = errors.New("realtime clock unavailable")
)
