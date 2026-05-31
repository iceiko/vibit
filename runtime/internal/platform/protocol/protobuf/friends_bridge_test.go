package protobuf

import (
	"errors"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appfriends "github.com/iceiko/vibit/runtime/internal/app/friends"
	friendsv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/friends/v1"
	friendsmodule "github.com/iceiko/vibit/runtime/internal/modules/friends"
	"google.golang.org/protobuf/proto"
)

func TestRouteRequestWithFriendsPayloadMapsRelationshipRequests(t *testing.T) {
	expectedVersion := int64(7)

	tests := []struct {
		name    string
		request app.RouteRequest
		assert  func(t *testing.T, payload any)
	}{
		{
			name: "send",
			request: app.RouteRequest{
				Route: appfriends.SendFriendRequestRoute(),
				Payload: &friendsv1.SendFriendRequestRequest{
					TargetPlayerId: "player-b",
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.SendFriendRequestRequest)
				if !ok {
					t.Fatalf("payload = %T, want SendFriendRequestRequest", payload)
				}
				if got.TargetPlayerID != "player-b" {
					t.Fatalf("payload = %#v, want mapped target player id", got)
				}
			},
		},
		{
			name: "accept",
			request: app.RouteRequest{
				Route: appfriends.AcceptFriendRequestRoute(),
				Payload: &friendsv1.AcceptFriendRequestRequest{
					TargetPlayerId:  "player-b",
					ExpectedVersion: &expectedVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.AcceptFriendRequestRequest)
				if !ok {
					t.Fatalf("payload = %T, want AcceptFriendRequestRequest", payload)
				}
				assertFriendsExpectedVersion(t, got.ExpectedVersion, expectedVersion)
			},
		},
		{
			name: "reject",
			request: app.RouteRequest{
				Route: appfriends.RejectFriendRequestRoute(),
				Payload: &friendsv1.RejectFriendRequestRequest{
					TargetPlayerId:  "player-b",
					ExpectedVersion: &expectedVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.RejectFriendRequestRequest)
				if !ok {
					t.Fatalf("payload = %T, want RejectFriendRequestRequest", payload)
				}
				assertFriendsExpectedVersion(t, got.ExpectedVersion, expectedVersion)
			},
		},
		{
			name: "remove",
			request: app.RouteRequest{
				Route: appfriends.RemoveFriendRoute(),
				Payload: &friendsv1.RemoveFriendRequest{
					TargetPlayerId:  "player-b",
					ExpectedVersion: &expectedVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.RemoveFriendRequest)
				if !ok {
					t.Fatalf("payload = %T, want RemoveFriendRequest", payload)
				}
				assertFriendsExpectedVersion(t, got.ExpectedVersion, expectedVersion)
			},
		},
		{
			name: "block",
			request: app.RouteRequest{
				Route: appfriends.BlockPlayerRoute(),
				Payload: &friendsv1.BlockPlayerRequest{
					TargetPlayerId:  "player-b",
					ExpectedVersion: &expectedVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.BlockPlayerRequest)
				if !ok {
					t.Fatalf("payload = %T, want BlockPlayerRequest", payload)
				}
				assertFriendsExpectedVersion(t, got.ExpectedVersion, expectedVersion)
			},
		},
		{
			name: "unblock",
			request: app.RouteRequest{
				Route: appfriends.UnblockPlayerRoute(),
				Payload: &friendsv1.UnblockPlayerRequest{
					TargetPlayerId:  "player-b",
					ExpectedVersion: &expectedVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.UnblockPlayerRequest)
				if !ok {
					t.Fatalf("payload = %T, want UnblockPlayerRequest", payload)
				}
				assertFriendsExpectedVersion(t, got.ExpectedVersion, expectedVersion)
			},
		},
		{
			name: "list",
			request: app.RouteRequest{
				Route: appfriends.ListFriendRelationshipsRoute(),
				Payload: &friendsv1.ListFriendRelationshipsRequest{
					Status:         "friends",
					Limit:          25,
					AfterPairToken: "pair-page-1",
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.ListFriendRelationshipsRequest)
				if !ok {
					t.Fatalf("payload = %T, want ListFriendRelationshipsRequest", payload)
				}
				if got.Status != friendsmodule.FriendRelationshipStatusFriends ||
					got.Limit != 25 ||
					got.AfterPairToken != "pair-page-1" {
					t.Fatalf("payload = %#v, want mapped list request", got)
				}
			},
		},
		{
			name: "status",
			request: app.RouteRequest{
				Route: appfriends.GetFriendRelationshipStatusRoute(),
				Payload: &friendsv1.GetFriendRelationshipStatusRequest{
					TargetPlayerId: "player-b",
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appfriends.GetFriendRelationshipStatusRequest)
				if !ok {
					t.Fatalf("payload = %T, want GetFriendRelationshipStatusRequest", payload)
				}
				if got.TargetPlayerID != "player-b" {
					t.Fatalf("payload = %#v, want mapped target player id", got)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mapped, err := RouteRequestWithDomainPayload(tc.request)
			if err != nil {
				t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
			}
			tc.assert(t, mapped.Payload)
		})
	}
}

func TestRouteRequestWithFriendsPayloadPreservesMissingExpectedVersion(t *testing.T) {
	mapped, err := RouteRequestWithDomainPayload(app.RouteRequest{
		Route: appfriends.AcceptFriendRequestRoute(),
		Payload: &friendsv1.AcceptFriendRequestRequest{
			TargetPlayerId: "player-b",
		},
	})
	if err != nil {
		t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
	}
	payload, ok := mapped.Payload.(appfriends.AcceptFriendRequestRequest)
	if !ok {
		t.Fatalf("Payload = %T, want AcceptFriendRequestRequest", mapped.Payload)
	}
	if payload.ExpectedVersion != nil {
		t.Fatalf("ExpectedVersion = %#v, want nil when optional field is absent", payload.ExpectedVersion)
	}
}

func TestRouteRequestWithFriendsPayloadRejectsWrongPayload(t *testing.T) {
	_, err := RouteRequestWithDomainPayload(app.RouteRequest{
		Route:   appfriends.SendFriendRequestRoute(),
		Payload: &friendsv1.ListFriendRelationshipsRequest{},
	})
	if err == nil {
		t.Fatal("RouteRequestWithDomainPayload() error = nil, want bridge error")
	}
	var bridgeErr *PayloadBridgeError
	if !errors.As(err, &bridgeErr) {
		t.Fatalf("error = %T %v, want *PayloadBridgeError", err, err)
	}
}

func TestProtoPayloadFromFriendsResultMapsResponses(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 8, 30, 0, 123, time.FixedZone("test", 8*60*60))
	updatedAt := createdAt.Add(time.Minute)
	first := friendsBridgeRelationship("relationship-1", "player-a", "player-b", "player-a", friendsmodule.FriendRelationshipLifecyclePending, 7, createdAt, updatedAt)
	second := friendsBridgeRelationship("relationship-2", "player-a", "player-c", "player-c", friendsmodule.FriendRelationshipLifecycleFriends, 8, createdAt.Add(time.Hour), updatedAt.Add(time.Hour))

	sendPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appfriends.SendFriendRequestRoute(),
		Payload: appfriends.FriendRelationshipResult{
			Status:       appfriends.FriendRelationshipOperationStatusSent,
			PublicStatus: appfriends.FriendRelationshipPublicStatusOutgoingRequestPending,
			Relationship: first,
			Version:      first.Version,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(send) error = %v, want nil", err)
	}
	sendResponse, ok := sendPayload.(*friendsv1.SendFriendRequestResponse)
	if !ok {
		t.Fatalf("send payload = %T, want SendFriendRequestResponse", sendPayload)
	}
	if sendResponse.GetStatus() != string(appfriends.FriendRelationshipOperationStatusSent) ||
		sendResponse.GetVersion() != int64(first.Version) {
		t.Fatalf("send response = %#v, want sent status and version", sendResponse)
	}
	assertProtoFriendRelationship(t, sendResponse.GetRelationship(), first, appfriends.FriendRelationshipPublicStatusOutgoingRequestPending)

	listPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appfriends.ListFriendRelationshipsRoute(),
		Payload: appfriends.FriendRelationshipResult{
			Status: appfriends.FriendRelationshipOperationStatusListed,
			Relationships: []appfriends.FriendRelationshipView{
				{Relationship: first, PublicStatus: appfriends.FriendRelationshipPublicStatusOutgoingRequestPending, Version: first.Version},
				{Relationship: second, PublicStatus: appfriends.FriendRelationshipPublicStatusFriends, Version: second.Version},
			},
			NextPairToken: "pair-page-2",
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(list) error = %v, want nil", err)
	}
	listResponse, ok := listPayload.(*friendsv1.ListFriendRelationshipsResponse)
	if !ok {
		t.Fatalf("list payload = %T, want ListFriendRelationshipsResponse", listPayload)
	}
	if listResponse.GetPage().GetNextPairToken() != "pair-page-2" ||
		len(listResponse.GetPage().GetRelationships()) != 2 {
		t.Fatalf("list response = %#v, want page with two relationships", listResponse)
	}
	assertProtoFriendRelationship(t, listResponse.GetPage().GetRelationships()[1], second, appfriends.FriendRelationshipPublicStatusFriends)

	statusPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appfriends.GetFriendRelationshipStatusRoute(),
		Payload: &appfriends.FriendRelationshipResult{
			Status:       appfriends.FriendRelationshipOperationStatusFound,
			PublicStatus: appfriends.FriendRelationshipPublicStatusNone,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(status) error = %v, want nil", err)
	}
	statusResponse, ok := statusPayload.(*friendsv1.GetFriendRelationshipStatusResponse)
	if !ok {
		t.Fatalf("status payload = %T, want GetFriendRelationshipStatusResponse", statusPayload)
	}
	if statusResponse.GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusNone) ||
		statusResponse.GetRelationship() != nil ||
		statusResponse.GetVersion() != 0 {
		t.Fatalf("status response = %#v, want none status without relationship", statusResponse)
	}
}

func TestFriendsProtoPayloadShapeOmitsActorProofSecretsAndBroadSocialFields(t *testing.T) {
	messages := []proto.Message{
		&friendsv1.FriendRelationship{},
		&friendsv1.SendFriendRequestRequest{},
		&friendsv1.AcceptFriendRequestRequest{},
		&friendsv1.RejectFriendRequestRequest{},
		&friendsv1.RemoveFriendRequest{},
		&friendsv1.BlockPlayerRequest{},
		&friendsv1.UnblockPlayerRequest{},
		&friendsv1.ListFriendRelationshipsRequest{},
		&friendsv1.GetFriendRelationshipStatusRequest{},
	}
	for _, message := range messages {
		descriptor := message.ProtoReflect().Descriptor()
		fields := descriptor.Fields()
		fieldNames := map[string]bool{}
		for i := 0; i < fields.Len(); i++ {
			fieldNames[string(fields.Get(i).Name())] = true
		}
		for _, forbidden := range []string{
			"actor_id",
			"session_id",
			"access_token",
			"credential_proof",
			"lookup_digest",
			"verifier_digest",
			"sql",
			"chat_channel_id",
			"group_id",
			"party_id",
			"match_id",
		} {
			if fieldNames[forbidden] {
				t.Fatalf("%s has forbidden field %q", descriptor.FullName(), forbidden)
			}
		}
	}
}

func assertFriendsExpectedVersion(t *testing.T, got *friendsmodule.FriendRelationshipVersion, want int64) {
	t.Helper()
	if got == nil || *got != friendsmodule.FriendRelationshipVersion(want) {
		t.Fatalf("ExpectedVersion = %#v, want %d", got, want)
	}
}

func friendsBridgeRelationship(relationshipID, low, high, requestedBy string, lifecycle friendsmodule.FriendRelationshipLifecycleState, version friendsmodule.FriendRelationshipVersion, createdAt, updatedAt time.Time) friendsmodule.FriendRelationship {
	return friendsmodule.FriendRelationship{
		RelationshipID: friendsmodule.FriendRelationshipID(relationshipID),
		Pair: friendsmodule.FriendRelationshipPair{
			PlayerLowID:  low,
			PlayerHighID: high,
		},
		RequestedByPlayerID: requestedBy,
		LifecycleState:      lifecycle,
		Version:             version,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		StateChangedAt:      updatedAt,
	}
}

func assertProtoFriendRelationship(t *testing.T, got *friendsv1.FriendRelationship, want friendsmodule.FriendRelationship, publicStatus appfriends.FriendRelationshipPublicStatus) {
	t.Helper()
	if got == nil {
		t.Fatal("FriendRelationship = nil, want mapped relationship")
	}
	if got.GetRelationshipId() != string(want.RelationshipID) ||
		got.GetPlayerLowId() != want.Pair.PlayerLowID ||
		got.GetPlayerHighId() != want.Pair.PlayerHighID ||
		got.GetRequestedByPlayerId() != want.RequestedByPlayerID ||
		got.GetLifecycleState() != string(want.LifecycleState) ||
		got.GetPublicStatus() != string(publicStatus) ||
		got.GetVersion() != int64(want.Version) ||
		got.GetCreatedAt() != want.CreatedAt.UTC().Format(time.RFC3339Nano) ||
		got.GetUpdatedAt() != want.UpdatedAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("FriendRelationship = %#v, want mapped relationship %#v with status %s", got, want, publicStatus)
	}
}
