package protobuf

import (
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
)

func TestBuildErrorEnvelopeFromApplicationResult(t *testing.T) {
	result := app.ApplicationResult{
		RequestID: " request-1 ",
		Route:     inventory.GrantItemRoute(),
		Target:    app.Target{Scope: app.TargetScopePlayer, ID: " player-1 "},
		Session: app.Session{
			ConnectionID:    " connection-1 ",
			SessionID:       " session-1 ",
			PlayerID:        " player-1 ",
			ConnectionEpoch: 9,
		},
		Error: &app.ApplicationError{
			Code:    app.ErrorCode(inventory.ErrorCodeInventoryPermission),
			Message: " actor cannot grant inventory items ",
			Route:   inventory.GrantItemRoute(),
		},
	}

	envelope, err := BuildErrorEnvelopeFromApplicationResult(result)
	if err != nil {
		t.Fatalf("BuildErrorEnvelopeFromApplicationResult() error = %v, want nil", err)
	}

	if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_ERROR {
		t.Fatalf("kind = %v, want MESSAGE_KIND_ERROR", envelope.GetKind())
	}
	if envelope.GetRequestId() != "request-1" {
		t.Fatalf("request_id = %q, want request-1", envelope.GetRequestId())
	}
	if envelope.GetModule() != inventory.ModuleName || envelope.GetName() != inventory.CommandGrantItem {
		t.Fatalf("route = %s.%s, want inventory.GrantItem", envelope.GetModule(), envelope.GetName())
	}
	if envelope.GetPayloadType() != "" || len(envelope.GetPayload()) != 0 {
		t.Fatalf("payload_type = %q payload len = %d, want empty error payload", envelope.GetPayloadType(), len(envelope.GetPayload()))
	}
	if envelope.GetTarget().GetScope() != protocolv1.TargetScope_TARGET_SCOPE_PLAYER || envelope.GetTarget().GetId() != "player-1" {
		t.Fatalf("target = %#v, want normalized player target", envelope.GetTarget())
	}
	if envelope.GetSession().GetConnectionId() != "connection-1" || envelope.GetSession().GetSessionId() != "session-1" || envelope.GetSession().GetPlayerId() != "player-1" || envelope.GetSession().GetConnectionEpoch() != 9 {
		t.Fatalf("session = %#v, want normalized session", envelope.GetSession())
	}

	wireError := envelope.GetError()
	if wireError == nil {
		t.Fatal("error = nil, want protocol error")
	}
	if wireError.GetCode() != string(inventory.ErrorCodeInventoryPermission) {
		t.Fatalf("error code = %q, want %s", wireError.GetCode(), inventory.ErrorCodeInventoryPermission)
	}
	if wireError.GetMessage() != "actor cannot grant inventory items" {
		t.Fatalf("error message = %q, want trimmed public message", wireError.GetMessage())
	}
	if wireError.GetRequestId() != "request-1" {
		t.Fatalf("error request_id = %q, want request-1", wireError.GetRequestId())
	}
	if wireError.GetRetryable() {
		t.Fatal("retryable = true, want false")
	}
}

func TestBuildErrorEnvelopeFromApplicationErrorUsesFallbackRoute(t *testing.T) {
	envelope, err := BuildErrorEnvelopeFromApplicationError("request-1", &app.ApplicationError{
		Code:    app.ErrorCodeInvalidRoute,
		Message: "route is invalid",
	}, app.Target{}, app.Session{})
	if err != nil {
		t.Fatalf("BuildErrorEnvelopeFromApplicationError() error = %v, want nil", err)
	}

	if envelope.GetModule() != "system" || envelope.GetName() != "ApplicationError" {
		t.Fatalf("route = %s.%s, want system.ApplicationError", envelope.GetModule(), envelope.GetName())
	}
}

func TestBuildErrorEnvelopeFromApplicationResultRejectsMissingError(t *testing.T) {
	_, err := BuildErrorEnvelopeFromApplicationResult(app.ApplicationResult{})
	if err == nil {
		t.Fatal("BuildErrorEnvelopeFromApplicationResult() error = nil, want validation error")
	}
}
