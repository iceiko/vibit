package bootstrap

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appfriends "github.com/iceiko/vibit/runtime/internal/app/friends"
	friendsmodule "github.com/iceiko/vibit/runtime/internal/modules/friends"
)

func TestFriendsRouteHandlersRegisterRelationshipRoutesAndPassValidatedIdentity(t *testing.T) {
	dispatcher := app.NewDispatcher()
	service := &recordingFriendsService{}
	handlers := FriendsRouteHandlers{Service: service}

	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}

	identity := app.ValidatedPlayerIdentity("player-a", app.Session{
		ConnectionID:    "connection-1",
		SessionID:       "session-1",
		PlayerID:        "player-a",
		ConnectionEpoch: 7,
	})
	expectedVersion := friendsmodule.FriendRelationshipVersion(3)

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-send",
		Route:     appfriends.SendFriendRequestRoute(),
		Identity:  identity,
		Payload: appfriends.SendFriendRequestRequest{
			TargetPlayerID: "player-b",
		},
	})
	if !service.sendCalled || service.sendRequest.Identity != identity ||
		service.sendRequest.TargetPlayerID != "player-b" {
		t.Fatalf("send request = %#v, want validated identity and target", service.sendRequest)
	}

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-accept",
		Route:     appfriends.AcceptFriendRequestRoute(),
		Identity:  identity,
		Payload: appfriends.AcceptFriendRequestRequest{
			TargetPlayerID:  "player-b",
			ExpectedVersion: &expectedVersion,
		},
	})
	if !service.acceptCalled || service.acceptRequest.Identity != identity ||
		service.acceptRequest.ExpectedVersion == nil ||
		*service.acceptRequest.ExpectedVersion != expectedVersion {
		t.Fatalf("accept request = %#v, want validated identity and expected version", service.acceptRequest)
	}

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-reject",
		Route:     appfriends.RejectFriendRequestRoute(),
		Identity:  identity,
		Payload: appfriends.RejectFriendRequestRequest{
			TargetPlayerID:  "player-b",
			ExpectedVersion: &expectedVersion,
		},
	})
	if !service.rejectCalled || service.rejectRequest.Identity != identity ||
		service.rejectRequest.ExpectedVersion == nil ||
		*service.rejectRequest.ExpectedVersion != expectedVersion {
		t.Fatalf("reject request = %#v, want validated identity and expected version", service.rejectRequest)
	}

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-remove",
		Route:     appfriends.RemoveFriendRoute(),
		Identity:  identity,
		Payload: appfriends.RemoveFriendRequest{
			TargetPlayerID:  "player-b",
			ExpectedVersion: &expectedVersion,
		},
	})
	if !service.removeCalled || service.removeRequest.Identity != identity ||
		service.removeRequest.ExpectedVersion == nil ||
		*service.removeRequest.ExpectedVersion != expectedVersion {
		t.Fatalf("remove request = %#v, want validated identity and expected version", service.removeRequest)
	}

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-block",
		Route:     appfriends.BlockPlayerRoute(),
		Identity:  identity,
		Payload: appfriends.BlockPlayerRequest{
			TargetPlayerID:  "player-b",
			ExpectedVersion: &expectedVersion,
		},
	})
	if !service.blockCalled || service.blockRequest.Identity != identity ||
		service.blockRequest.ExpectedVersion == nil ||
		*service.blockRequest.ExpectedVersion != expectedVersion {
		t.Fatalf("block request = %#v, want validated identity and expected version", service.blockRequest)
	}

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-unblock",
		Route:     appfriends.UnblockPlayerRoute(),
		Identity:  identity,
		Payload: appfriends.UnblockPlayerRequest{
			TargetPlayerID:  "player-b",
			ExpectedVersion: &expectedVersion,
		},
	})
	if !service.unblockCalled || service.unblockRequest.Identity != identity ||
		service.unblockRequest.ExpectedVersion == nil ||
		*service.unblockRequest.ExpectedVersion != expectedVersion {
		t.Fatalf("unblock request = %#v, want validated identity and expected version", service.unblockRequest)
	}

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-list",
		Route:     appfriends.ListFriendRelationshipsRoute(),
		Identity:  identity,
		Payload: appfriends.ListFriendRelationshipsRequest{
			Status:         friendsmodule.FriendRelationshipStatusFriends,
			Limit:          25,
			AfterPairToken: "pair-page-1",
		},
	})
	if !service.listCalled || service.listRequest.Identity != identity ||
		service.listRequest.Status != friendsmodule.FriendRelationshipStatusFriends ||
		service.listRequest.Limit != 25 ||
		service.listRequest.AfterPairToken != "pair-page-1" {
		t.Fatalf("list request = %#v, want validated identity and pagination", service.listRequest)
	}

	dispatchFriendsRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-status",
		Route:     appfriends.GetFriendRelationshipStatusRoute(),
		Identity:  identity,
		Payload: appfriends.GetFriendRelationshipStatusRequest{
			TargetPlayerID: "player-b",
		},
	})
	if !service.statusCalled || service.statusRequest.Identity != identity ||
		service.statusRequest.TargetPlayerID != "player-b" {
		t.Fatalf("status request = %#v, want validated identity and target", service.statusRequest)
	}
}

func TestFriendsRouteHandlersRejectMalformedPayloadBeforeService(t *testing.T) {
	service := &recordingFriendsService{}
	handlers := FriendsRouteHandlers{Service: service}

	result, err := handlers.HandleSendFriendRequestRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appfriends.SendFriendRequestRoute(),
		Payload:   "not a friends request",
	})
	if err == nil {
		t.Fatal("HandleSendFriendRequestRoute() error = nil, want invalid request error")
	}
	if service.sendCalled {
		t.Fatal("friends service was called for malformed payload")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appfriends.PublicErrorFriendshipInvalidRequest) {
		t.Fatalf("result error = %#v, want friendship invalid request", result.Error)
	}
}

func TestFriendsRouteHandlersMapServiceErrorsWithoutPrivateDetailLeak(t *testing.T) {
	secretDetail := "internal sql friend_relationships player-a player-b access-token"
	service := &recordingFriendsService{
		sendResult: appfriends.FriendRelationshipResult{
			Status:          appfriends.FriendRelationshipOperationStatusRejected,
			PublicErrorCode: appfriends.PublicErrorFriendshipVersionMismatch,
			FailureClass:    appfriends.FailureClassVersionMismatch,
		},
		sendErr: &appfriends.ServiceError{
			Operation:  appfriends.OperationSendFriendRequest,
			Class:      appfriends.FailureClassVersionMismatch,
			PublicCode: appfriends.PublicErrorFriendshipVersionMismatch,
			Err:        errors.New(secretDetail),
		},
	}
	handlers := FriendsRouteHandlers{Service: service}

	result, err := handlers.HandleSendFriendRequestRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appfriends.SendFriendRequestRoute(),
		Payload: appfriends.SendFriendRequestRequest{
			TargetPlayerID: "player-b",
		},
	})
	if err == nil {
		t.Fatal("HandleSendFriendRequestRoute() error = nil, want public friends error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appfriends.PublicErrorFriendshipVersionMismatch) {
		t.Fatalf("result error = %#v, want version mismatch code", result.Error)
	}
	errorText := result.Error.Error()
	for _, leak := range []string{secretDetail, "friend_relationships", "player-a", "player-b", "access-token", "sql"} {
		if strings.Contains(errorText, leak) {
			t.Fatalf("error %q leaked private detail %q", errorText, leak)
		}
	}
	if result.Payload != nil {
		t.Fatalf("result Payload = %#v, want nil on friends failure", result.Payload)
	}
}

func TestFriendsRouteHandlersRequireService(t *testing.T) {
	handlers := FriendsRouteHandlers{}

	result, err := handlers.HandleGetFriendRelationshipStatusRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appfriends.GetFriendRelationshipStatusRoute(),
		Payload:   appfriends.GetFriendRelationshipStatusRequest{TargetPlayerID: "player-b"},
	})
	if err == nil {
		t.Fatal("HandleGetFriendRelationshipStatusRoute() error = nil, want unavailable error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appfriends.PublicErrorFriendshipUnavailable) {
		t.Fatalf("result error = %#v, want friendship unavailable", result.Error)
	}
}

func dispatchFriendsRoute(t *testing.T, dispatcher *app.Dispatcher, request app.RouteRequest) app.ApplicationResult {
	t.Helper()
	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch(%s) error = %v, want nil", app.RenderRouteKey(request.Route), err)
	}
	if result.PayloadType != "friends.FriendRelationshipResult" {
		t.Fatalf("PayloadType = %q, want friends.FriendRelationshipResult", result.PayloadType)
	}
	if _, ok := result.Payload.(appfriends.FriendRelationshipResult); !ok {
		t.Fatalf("Payload = %T, want friends.FriendRelationshipResult", result.Payload)
	}
	return result
}

type recordingFriendsService struct {
	sendCalled     bool
	sendRequest    appfriends.SendFriendRequestRequest
	sendResult     appfriends.FriendRelationshipResult
	sendErr        error
	acceptCalled   bool
	acceptRequest  appfriends.AcceptFriendRequestRequest
	acceptResult   appfriends.FriendRelationshipResult
	acceptErr      error
	rejectCalled   bool
	rejectRequest  appfriends.RejectFriendRequestRequest
	rejectResult   appfriends.FriendRelationshipResult
	rejectErr      error
	removeCalled   bool
	removeRequest  appfriends.RemoveFriendRequest
	removeResult   appfriends.FriendRelationshipResult
	removeErr      error
	blockCalled    bool
	blockRequest   appfriends.BlockPlayerRequest
	blockResult    appfriends.FriendRelationshipResult
	blockErr       error
	unblockCalled  bool
	unblockRequest appfriends.UnblockPlayerRequest
	unblockResult  appfriends.FriendRelationshipResult
	unblockErr     error
	listCalled     bool
	listRequest    appfriends.ListFriendRelationshipsRequest
	listResult     appfriends.FriendRelationshipResult
	listErr        error
	statusCalled   bool
	statusRequest  appfriends.GetFriendRelationshipStatusRequest
	statusResult   appfriends.FriendRelationshipResult
	statusErr      error
}

func (s *recordingFriendsService) SendFriendRequest(_ context.Context, request appfriends.SendFriendRequestRequest) (appfriends.FriendRelationshipResult, error) {
	s.sendCalled = true
	s.sendRequest = request
	return s.resultOrDefault(s.sendResult, appfriends.FriendRelationshipOperationStatusSent), s.sendErr
}

func (s *recordingFriendsService) AcceptFriendRequest(_ context.Context, request appfriends.AcceptFriendRequestRequest) (appfriends.FriendRelationshipResult, error) {
	s.acceptCalled = true
	s.acceptRequest = request
	return s.resultOrDefault(s.acceptResult, appfriends.FriendRelationshipOperationStatusAccepted), s.acceptErr
}

func (s *recordingFriendsService) RejectFriendRequest(_ context.Context, request appfriends.RejectFriendRequestRequest) (appfriends.FriendRelationshipResult, error) {
	s.rejectCalled = true
	s.rejectRequest = request
	return s.resultOrDefault(s.rejectResult, appfriends.FriendRelationshipOperationStatusRequestRejected), s.rejectErr
}

func (s *recordingFriendsService) RemoveFriend(_ context.Context, request appfriends.RemoveFriendRequest) (appfriends.FriendRelationshipResult, error) {
	s.removeCalled = true
	s.removeRequest = request
	return s.resultOrDefault(s.removeResult, appfriends.FriendRelationshipOperationStatusRemoved), s.removeErr
}

func (s *recordingFriendsService) BlockPlayer(_ context.Context, request appfriends.BlockPlayerRequest) (appfriends.FriendRelationshipResult, error) {
	s.blockCalled = true
	s.blockRequest = request
	return s.resultOrDefault(s.blockResult, appfriends.FriendRelationshipOperationStatusBlocked), s.blockErr
}

func (s *recordingFriendsService) UnblockPlayer(_ context.Context, request appfriends.UnblockPlayerRequest) (appfriends.FriendRelationshipResult, error) {
	s.unblockCalled = true
	s.unblockRequest = request
	return s.resultOrDefault(s.unblockResult, appfriends.FriendRelationshipOperationStatusUnblocked), s.unblockErr
}

func (s *recordingFriendsService) ListFriendRelationships(_ context.Context, request appfriends.ListFriendRelationshipsRequest) (appfriends.FriendRelationshipResult, error) {
	s.listCalled = true
	s.listRequest = request
	return s.resultOrDefault(s.listResult, appfriends.FriendRelationshipOperationStatusListed), s.listErr
}

func (s *recordingFriendsService) GetFriendRelationshipStatus(_ context.Context, request appfriends.GetFriendRelationshipStatusRequest) (appfriends.FriendRelationshipResult, error) {
	s.statusCalled = true
	s.statusRequest = request
	return s.resultOrDefault(s.statusResult, appfriends.FriendRelationshipOperationStatusFound), s.statusErr
}

func (s *recordingFriendsService) resultOrDefault(result appfriends.FriendRelationshipResult, status appfriends.FriendRelationshipOperationStatus) appfriends.FriendRelationshipResult {
	if result.Status != "" {
		return result
	}
	relationship := friendsmodule.FriendRelationship{
		RelationshipID: friendsmodule.FriendRelationshipID("relationship-1"),
		Pair: friendsmodule.FriendRelationshipPair{
			PlayerLowID:  "player-a",
			PlayerHighID: "player-b",
		},
		RequestedByPlayerID: "player-a",
		LifecycleState:      friendsmodule.FriendRelationshipLifecycleFriends,
		Version:             friendsmodule.FriendRelationshipVersion(4),
		CreatedAt:           time.Date(2026, 5, 26, 8, 30, 0, 0, time.UTC),
		UpdatedAt:           time.Date(2026, 5, 26, 8, 31, 0, 0, time.UTC),
		StateChangedAt:      time.Date(2026, 5, 26, 8, 31, 0, 0, time.UTC),
	}
	return appfriends.FriendRelationshipResult{
		Status:       status,
		PublicStatus: appfriends.FriendRelationshipPublicStatusFriends,
		Relationship: relationship,
		Version:      relationship.Version,
	}
}
