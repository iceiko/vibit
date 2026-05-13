package protobuf

import (
	"context"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"google.golang.org/protobuf/proto"
)

type requestLoopFixture struct {
	Dispatcher *app.Dispatcher
	Repository *inventory.MemoryRepository
}

func newInventoryRequestLoopFixture(t *testing.T, grantAllowed bool, readAllowed bool) requestLoopFixture {
	t.Helper()

	repository := inventory.NewMemoryRepository()
	handlers := inventory.Handlers{
		Repository:       repository,
		PermissionPolicy: inventory.StaticPermissionPolicy{GrantAllowed: grantAllowed, ReadAllowed: readAllowed},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 16},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "request-loop-event"},
		Clock:            requestLoopClock{},
	}
	dispatcher := app.NewDispatcher()
	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}

	return requestLoopFixture{
		Dispatcher: dispatcher,
		Repository: repository,
	}
}

func mustBuildEnvelope(t *testing.T, route app.RouteKey, requestID string, target app.Target, session app.Session, payload proto.Message) *protocolv1.Envelope {
	t.Helper()

	envelope, err := BuildEnvelope(route, requestID, target, session, payload)
	if err != nil {
		t.Fatalf("BuildEnvelope() error = %v, want nil", err)
	}
	return envelope
}

func mustMarshalEnvelope(t *testing.T, route app.RouteKey, payload proto.Message) []byte {
	t.Helper()

	envelope := mustBuildEnvelope(
		t,
		route,
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{SessionID: "session-1", PlayerID: "player-1", ConnectionEpoch: 7},
		payload,
	)
	encoded, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal(envelope) error = %v, want nil", err)
	}
	return encoded
}

func mustUnmarshalFrameEnvelope(t *testing.T, payload []byte) *protocolv1.Envelope {
	t.Helper()

	var envelope protocolv1.Envelope
	if err := proto.Unmarshal(payload, &envelope); err != nil {
		t.Fatalf("proto.Unmarshal(response) error = %v, want nil", err)
	}
	return &envelope
}

func mustRouteRequestFromEnvelope(t *testing.T, envelope *protocolv1.Envelope) app.RouteRequest {
	t.Helper()

	request, err := RouteRequestFromEnvelopeForDispatch(envelope)
	if err != nil {
		t.Fatalf("RouteRequestFromEnvelopeForDispatch() error = %v, want nil", err)
	}
	return request
}

type requestLoopDispatcherFunc func(context.Context, app.RouteRequest) (app.ApplicationResult, error)

func (f requestLoopDispatcherFunc) Dispatch(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	return f(ctx, request)
}

type requestLoopClock struct{}

func (requestLoopClock) Now() time.Time {
	return time.Date(2026, 5, 13, 12, 0, 0, 0, time.UTC)
}
