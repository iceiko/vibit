package protobuf

import (
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type DecodeEnvelopeError struct {
	Message string
}

func (e *DecodeEnvelopeError) Error() string {
	return e.Message
}

func RouteRequestFromEnvelope(envelope *protocolv1.Envelope) (app.RouteRequest, error) {
	request, err := RouteRequestMetadataFromEnvelope(envelope)
	if err != nil {
		return app.RouteRequest{}, err
	}
	payload, err := DecodePayload(request.PayloadType, request.PayloadBytes)
	if err != nil {
		return app.RouteRequest{}, err
	}
	request.Payload = payload
	return request, nil
}

func RouteRequestMetadataFromEnvelope(envelope *protocolv1.Envelope) (app.RouteRequest, error) {
	if envelope == nil {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "envelope is nil"}
	}

	module := strings.TrimSpace(envelope.GetModule())
	name := strings.TrimSpace(envelope.GetName())
	if module == "" {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "envelope module is empty"}
	}
	if name == "" {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "envelope name is empty"}
	}
	kind := appMessageKind(envelope.GetKind())
	if kind == "" {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "envelope kind is unsupported"}
	}

	target := envelope.GetTarget()
	if target == nil {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "envelope target is nil"}
	}

	session := envelope.GetSession()
	if session == nil {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "envelope session is nil"}
	}

	payloadType := strings.TrimSpace(envelope.GetPayloadType())
	if payloadType == "" {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "envelope payload_type is empty"}
	}
	payloadBytes := append([]byte(nil), envelope.GetPayload()...)

	appSession := app.Session{
		ConnectionID:    strings.TrimSpace(session.GetConnectionId()),
		SessionID:       strings.TrimSpace(session.GetSessionId()),
		PlayerID:        strings.TrimSpace(session.GetPlayerId()),
		ConnectionEpoch: session.GetConnectionEpoch(),
	}

	return app.RouteRequest{
		RequestID: strings.TrimSpace(envelope.GetRequestId()),
		Route: app.RouteKey{
			Kind:   kind,
			Module: module,
			Name:   name,
		},
		Target: app.Target{
			Scope: appTargetScope(target.GetScope()),
			ID:    strings.TrimSpace(target.GetId()),
		},
		Session:      appSession,
		Identity:     app.MetadataOnlyIdentityFromSession(appSession),
		PayloadType:  payloadType,
		PayloadBytes: payloadBytes,
	}, nil
}

func DecodePayload(payloadType string, payload []byte) (proto.Message, error) {
	messageName := protoreflect.FullName(strings.TrimSpace(payloadType))
	if !messageName.IsValid() {
		return nil, &DecodeEnvelopeError{Message: "payload_type is invalid"}
	}

	messageType, err := protoregistry.GlobalTypes.FindMessageByName(messageName)
	if err != nil {
		return nil, &DecodeEnvelopeError{Message: "payload_type is unknown"}
	}

	message := messageType.New().Interface()
	if err := proto.Unmarshal(payload, message); err != nil {
		return nil, err
	}

	return message, nil
}

func BuildEnvelope(route app.RouteKey, requestID string, target app.Target, session app.Session, payload proto.Message) (*protocolv1.Envelope, error) {
	if payload == nil {
		return nil, &DecodeEnvelopeError{Message: "payload is nil"}
	}
	kind := protoMessageKind(route.Kind)
	if kind == protocolv1.MessageKind_MESSAGE_KIND_UNSPECIFIED {
		return nil, &DecodeEnvelopeError{Message: "route kind is unsupported"}
	}
	module := strings.TrimSpace(route.Module)
	if module == "" {
		return nil, &DecodeEnvelopeError{Message: "route module is empty"}
	}
	name := strings.TrimSpace(route.Name)
	if name == "" {
		return nil, &DecodeEnvelopeError{Message: "route name is empty"}
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &protocolv1.Envelope{
		ProtocolVersion: 1,
		Kind:            kind,
		RequestId:       strings.TrimSpace(requestID),
		Module:          module,
		Name:            name,
		Target: &protocolv1.Target{
			Scope: protoTargetScope(target.Scope),
			Id:    strings.TrimSpace(target.ID),
		},
		Session: &protocolv1.Session{
			ConnectionId:    strings.TrimSpace(session.ConnectionID),
			SessionId:       strings.TrimSpace(session.SessionID),
			PlayerId:        strings.TrimSpace(session.PlayerID),
			ConnectionEpoch: session.ConnectionEpoch,
		},
		PayloadType: PayloadType(payload),
		Payload:     payloadBytes,
	}, nil
}

func PayloadType(message proto.Message) string {
	if message == nil {
		return ""
	}
	return string(message.ProtoReflect().Descriptor().FullName())
}

func PayloadTypeName(message protoreflect.MessageDescriptor) string {
	if message == nil {
		return ""
	}
	return string(message.FullName())
}

func appMessageKind(kind protocolv1.MessageKind) app.MessageKind {
	switch kind {
	case protocolv1.MessageKind_MESSAGE_KIND_COMMAND:
		return app.MessageKindCommand
	case protocolv1.MessageKind_MESSAGE_KIND_QUERY:
		return app.MessageKindQuery
	case protocolv1.MessageKind_MESSAGE_KIND_EVENT:
		return app.MessageKindEvent
	case protocolv1.MessageKind_MESSAGE_KIND_ERROR:
		return app.MessageKindError
	case protocolv1.MessageKind_MESSAGE_KIND_SYSTEM:
		return app.MessageKindSystem
	case protocolv1.MessageKind_MESSAGE_KIND_ACK:
		return app.MessageKindAck
	case protocolv1.MessageKind_MESSAGE_KIND_HEARTBEAT:
		return app.MessageKindHeartbeat
	case protocolv1.MessageKind_MESSAGE_KIND_INPUT:
		return app.MessageKindInput
	case protocolv1.MessageKind_MESSAGE_KIND_STATE:
		return app.MessageKindState
	default:
		return ""
	}
}

func protoMessageKind(kind app.MessageKind) protocolv1.MessageKind {
	switch kind {
	case app.MessageKindCommand:
		return protocolv1.MessageKind_MESSAGE_KIND_COMMAND
	case app.MessageKindQuery:
		return protocolv1.MessageKind_MESSAGE_KIND_QUERY
	case app.MessageKindEvent:
		return protocolv1.MessageKind_MESSAGE_KIND_EVENT
	case app.MessageKindError:
		return protocolv1.MessageKind_MESSAGE_KIND_ERROR
	case app.MessageKindSystem:
		return protocolv1.MessageKind_MESSAGE_KIND_SYSTEM
	case app.MessageKindAck:
		return protocolv1.MessageKind_MESSAGE_KIND_ACK
	case app.MessageKindHeartbeat:
		return protocolv1.MessageKind_MESSAGE_KIND_HEARTBEAT
	case app.MessageKindInput:
		return protocolv1.MessageKind_MESSAGE_KIND_INPUT
	case app.MessageKindState:
		return protocolv1.MessageKind_MESSAGE_KIND_STATE
	default:
		return protocolv1.MessageKind_MESSAGE_KIND_UNSPECIFIED
	}
}

func appTargetScope(scope protocolv1.TargetScope) app.TargetScope {
	switch scope {
	case protocolv1.TargetScope_TARGET_SCOPE_GLOBAL:
		return app.TargetScopeGlobal
	case protocolv1.TargetScope_TARGET_SCOPE_PLAYER:
		return app.TargetScopePlayer
	case protocolv1.TargetScope_TARGET_SCOPE_PARTY:
		return app.TargetScopeParty
	case protocolv1.TargetScope_TARGET_SCOPE_ROOM:
		return app.TargetScopeRoom
	case protocolv1.TargetScope_TARGET_SCOPE_MATCH:
		return app.TargetScopeMatch
	case protocolv1.TargetScope_TARGET_SCOPE_STREAM:
		return app.TargetScopeStream
	case protocolv1.TargetScope_TARGET_SCOPE_SYSTEM:
		return app.TargetScopeSystem
	default:
		return ""
	}
}

func protoTargetScope(scope app.TargetScope) protocolv1.TargetScope {
	switch scope {
	case app.TargetScopeGlobal:
		return protocolv1.TargetScope_TARGET_SCOPE_GLOBAL
	case app.TargetScopePlayer:
		return protocolv1.TargetScope_TARGET_SCOPE_PLAYER
	case app.TargetScopeParty:
		return protocolv1.TargetScope_TARGET_SCOPE_PARTY
	case app.TargetScopeRoom:
		return protocolv1.TargetScope_TARGET_SCOPE_ROOM
	case app.TargetScopeMatch:
		return protocolv1.TargetScope_TARGET_SCOPE_MATCH
	case app.TargetScopeStream:
		return protocolv1.TargetScope_TARGET_SCOPE_STREAM
	case app.TargetScopeSystem:
		return protocolv1.TargetScope_TARGET_SCOPE_SYSTEM
	default:
		return protocolv1.TargetScope_TARGET_SCOPE_UNSPECIFIED
	}
}
