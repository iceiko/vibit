package protobuf

import (
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appfriends "github.com/iceiko/vibit/runtime/internal/app/friends"
	friendsv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/friends/v1"
	friendsmodule "github.com/iceiko/vibit/runtime/internal/modules/friends"
	"google.golang.org/protobuf/proto"
)

func routeRequestWithFriendsPayload(request app.RouteRequest) (app.RouteRequest, bool, error) {
	switch request.Route {
	case appfriends.SendFriendRequestRoute():
		payload, ok := request.Payload.(*friendsv1.SendFriendRequestRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.SendFriendRequestRequest")
		}
		request.Payload = appfriends.SendFriendRequestRequest{
			TargetPlayerID: payload.GetTargetPlayerId(),
		}
		return request, true, nil

	case appfriends.AcceptFriendRequestRoute():
		payload, ok := request.Payload.(*friendsv1.AcceptFriendRequestRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.AcceptFriendRequestRequest")
		}
		request.Payload = appfriends.AcceptFriendRequestRequest{
			TargetPlayerID:  payload.GetTargetPlayerId(),
			ExpectedVersion: friendRelationshipVersionFromOptionalInt64(payload.ExpectedVersion),
		}
		return request, true, nil

	case appfriends.RejectFriendRequestRoute():
		payload, ok := request.Payload.(*friendsv1.RejectFriendRequestRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.RejectFriendRequestRequest")
		}
		request.Payload = appfriends.RejectFriendRequestRequest{
			TargetPlayerID:  payload.GetTargetPlayerId(),
			ExpectedVersion: friendRelationshipVersionFromOptionalInt64(payload.ExpectedVersion),
		}
		return request, true, nil

	case appfriends.RemoveFriendRoute():
		payload, ok := request.Payload.(*friendsv1.RemoveFriendRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.RemoveFriendRequest")
		}
		request.Payload = appfriends.RemoveFriendRequest{
			TargetPlayerID:  payload.GetTargetPlayerId(),
			ExpectedVersion: friendRelationshipVersionFromOptionalInt64(payload.ExpectedVersion),
		}
		return request, true, nil

	case appfriends.BlockPlayerRoute():
		payload, ok := request.Payload.(*friendsv1.BlockPlayerRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.BlockPlayerRequest")
		}
		request.Payload = appfriends.BlockPlayerRequest{
			TargetPlayerID:  payload.GetTargetPlayerId(),
			ExpectedVersion: friendRelationshipVersionFromOptionalInt64(payload.ExpectedVersion),
		}
		return request, true, nil

	case appfriends.UnblockPlayerRoute():
		payload, ok := request.Payload.(*friendsv1.UnblockPlayerRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.UnblockPlayerRequest")
		}
		request.Payload = appfriends.UnblockPlayerRequest{
			TargetPlayerID:  payload.GetTargetPlayerId(),
			ExpectedVersion: friendRelationshipVersionFromOptionalInt64(payload.ExpectedVersion),
		}
		return request, true, nil

	case appfriends.ListFriendRelationshipsRoute():
		payload, ok := request.Payload.(*friendsv1.ListFriendRelationshipsRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.ListFriendRelationshipsRequest")
		}
		request.Payload = appfriends.ListFriendRelationshipsRequest{
			Status:         friendsmodule.FriendRelationshipStatus(payload.GetStatus()),
			Limit:          int(payload.GetLimit()),
			AfterPairToken: payload.GetAfterPairToken(),
		}
		return request, true, nil

	case appfriends.GetFriendRelationshipStatusRoute():
		payload, ok := request.Payload.(*friendsv1.GetFriendRelationshipStatusRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.friends.v1.GetFriendRelationshipStatusRequest")
		}
		request.Payload = appfriends.GetFriendRelationshipStatusRequest{
			TargetPlayerID: payload.GetTargetPlayerId(),
		}
		return request, true, nil

	default:
		return request, false, nil
	}
}

func protoPayloadFromFriendsRoute(route app.RouteKey, payload any) (proto.Message, bool, error) {
	switch route {
	case appfriends.SendFriendRequestRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		return &friendsv1.SendFriendRequestResponse{
			Relationship: protoFriendRelationship(result.Relationship, result.PublicStatus),
			Status:       string(result.Status),
			Version:      int64(result.Version),
		}, true, nil

	case appfriends.AcceptFriendRequestRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		return &friendsv1.AcceptFriendRequestResponse{
			Relationship: protoFriendRelationship(result.Relationship, result.PublicStatus),
			Status:       string(result.Status),
			Version:      int64(result.Version),
		}, true, nil

	case appfriends.RejectFriendRequestRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		return &friendsv1.RejectFriendRequestResponse{
			Relationship: protoFriendRelationship(result.Relationship, result.PublicStatus),
			Status:       string(result.Status),
			Version:      int64(result.Version),
		}, true, nil

	case appfriends.RemoveFriendRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		return &friendsv1.RemoveFriendResponse{
			Relationship: protoFriendRelationship(result.Relationship, result.PublicStatus),
			Status:       string(result.Status),
			Version:      int64(result.Version),
		}, true, nil

	case appfriends.BlockPlayerRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		return &friendsv1.BlockPlayerResponse{
			Relationship: protoFriendRelationship(result.Relationship, result.PublicStatus),
			Status:       string(result.Status),
			Version:      int64(result.Version),
		}, true, nil

	case appfriends.UnblockPlayerRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		return &friendsv1.UnblockPlayerResponse{
			Relationship: protoFriendRelationship(result.Relationship, result.PublicStatus),
			Status:       string(result.Status),
			Version:      int64(result.Version),
		}, true, nil

	case appfriends.ListFriendRelationshipsRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		relationships := make([]*friendsv1.FriendRelationship, 0, len(result.Relationships))
		for _, relationship := range result.Relationships {
			relationships = append(relationships, protoFriendRelationship(relationship.Relationship, relationship.PublicStatus))
		}
		return &friendsv1.ListFriendRelationshipsResponse{
			Page: &friendsv1.FriendRelationshipPage{
				Relationships: relationships,
				NextPairToken: result.NextPairToken,
			},
		}, true, nil

	case appfriends.GetFriendRelationshipStatusRoute():
		result, ok := friendsResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be friends.FriendRelationshipResult")
		}
		return &friendsv1.GetFriendRelationshipStatusResponse{
			PublicStatus: string(result.PublicStatus),
			Relationship: protoFriendRelationship(
				result.Relationship,
				result.PublicStatus,
			),
			Version: int64(result.Version),
		}, true, nil

	default:
		return nil, false, nil
	}
}

func friendsResultPayload(payload any) (appfriends.FriendRelationshipResult, bool) {
	result, ok := payload.(appfriends.FriendRelationshipResult)
	if !ok {
		if pointerResult, pointerOK := payload.(*appfriends.FriendRelationshipResult); pointerOK && pointerResult != nil {
			result = *pointerResult
			ok = true
		}
	}
	return result, ok
}

func friendRelationshipVersionFromOptionalInt64(value *int64) *friendsmodule.FriendRelationshipVersion {
	if value == nil {
		return nil
	}
	version := friendsmodule.FriendRelationshipVersion(*value)
	return &version
}

func protoFriendRelationship(relationship friendsmodule.FriendRelationship, publicStatus appfriends.FriendRelationshipPublicStatus) *friendsv1.FriendRelationship {
	if relationship.RelationshipID == "" {
		return nil
	}
	return &friendsv1.FriendRelationship{
		RelationshipId:      string(relationship.RelationshipID),
		PlayerLowId:         relationship.Pair.PlayerLowID,
		PlayerHighId:        relationship.Pair.PlayerHighID,
		RequestedByPlayerId: relationship.RequestedByPlayerID,
		LifecycleState:      string(relationship.LifecycleState),
		PublicStatus:        string(publicStatus),
		Version:             int64(relationship.Version),
		CreatedAt:           formatFriendRelationshipTime(relationship.CreatedAt),
		UpdatedAt:           formatFriendRelationshipTime(relationship.UpdatedAt),
	}
}

func formatFriendRelationshipTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
