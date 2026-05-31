package bootstrap

import (
	"context"
	"errors"

	"github.com/iceiko/vibit/runtime/internal/app"
	appfriends "github.com/iceiko/vibit/runtime/internal/app/friends"
)

type FriendsService interface {
	SendFriendRequest(context.Context, appfriends.SendFriendRequestRequest) (appfriends.FriendRelationshipResult, error)
	AcceptFriendRequest(context.Context, appfriends.AcceptFriendRequestRequest) (appfriends.FriendRelationshipResult, error)
	RejectFriendRequest(context.Context, appfriends.RejectFriendRequestRequest) (appfriends.FriendRelationshipResult, error)
	RemoveFriend(context.Context, appfriends.RemoveFriendRequest) (appfriends.FriendRelationshipResult, error)
	BlockPlayer(context.Context, appfriends.BlockPlayerRequest) (appfriends.FriendRelationshipResult, error)
	UnblockPlayer(context.Context, appfriends.UnblockPlayerRequest) (appfriends.FriendRelationshipResult, error)
	ListFriendRelationships(context.Context, appfriends.ListFriendRelationshipsRequest) (appfriends.FriendRelationshipResult, error)
	GetFriendRelationshipStatus(context.Context, appfriends.GetFriendRelationshipStatusRequest) (appfriends.FriendRelationshipResult, error)
}

type FriendsRouteHandlers struct {
	Service FriendsService
}

func (h FriendsRouteHandlers) RegisterRoutes(dispatcher *app.Dispatcher) error {
	if dispatcher == nil {
		return errors.New("friends bootstrap: dispatcher is nil")
	}
	if err := dispatcher.Register(appfriends.SendFriendRequestRoute(), app.HandlerFunc(h.HandleSendFriendRequestRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appfriends.AcceptFriendRequestRoute(), app.HandlerFunc(h.HandleAcceptFriendRequestRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appfriends.RejectFriendRequestRoute(), app.HandlerFunc(h.HandleRejectFriendRequestRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appfriends.RemoveFriendRoute(), app.HandlerFunc(h.HandleRemoveFriendRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appfriends.BlockPlayerRoute(), app.HandlerFunc(h.HandleBlockPlayerRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appfriends.UnblockPlayerRoute(), app.HandlerFunc(h.HandleUnblockPlayerRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appfriends.ListFriendRelationshipsRoute(), app.HandlerFunc(h.HandleListFriendRelationshipsRoute)); err != nil {
		return err
	}
	return dispatcher.Register(appfriends.GetFriendRelationshipStatusRoute(), app.HandlerFunc(h.HandleGetFriendRelationshipStatusRoute))
}

func (h FriendsRouteHandlers) HandleSendFriendRequestRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.SendFriendRequestRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.SendFriendRequest(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func (h FriendsRouteHandlers) HandleAcceptFriendRequestRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.AcceptFriendRequestRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.AcceptFriendRequest(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func (h FriendsRouteHandlers) HandleRejectFriendRequestRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.RejectFriendRequestRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.RejectFriendRequest(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func (h FriendsRouteHandlers) HandleRemoveFriendRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.RemoveFriendRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.RemoveFriend(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func (h FriendsRouteHandlers) HandleBlockPlayerRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.BlockPlayerRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.BlockPlayer(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func (h FriendsRouteHandlers) HandleUnblockPlayerRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.UnblockPlayerRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.UnblockPlayer(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func (h FriendsRouteHandlers) HandleListFriendRelationshipsRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.ListFriendRelationshipsRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.ListFriendRelationships(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func (h FriendsRouteHandlers) HandleGetFriendRelationshipStatusRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipUnavailable)
	}

	payload, ok := friendsPayload[appfriends.GetFriendRelationshipStatusRequest](request.Payload)
	if !ok {
		return friendsApplicationErrorResult(request, appfriends.PublicErrorFriendshipInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.GetFriendRelationshipStatus(ctx, payload)
	return friendsResultForRequest(request, result, err)
}

func friendsPayload[T any](payload any) (T, bool) {
	value, ok := payload.(T)
	if !ok {
		if pointerValue, pointerOK := payload.(*T); pointerOK && pointerValue != nil {
			value = *pointerValue
			ok = true
		}
	}
	return value, ok
}

func friendsResultForRequest(request app.RouteRequest, friendsResult appfriends.FriendRelationshipResult, err error) (app.ApplicationResult, error) {
	result := resultForRequest(request)
	if err != nil {
		appErr := friendsApplicationError(request.Route, friendsPublicErrorCode(friendsResult, err))
		result.Error = appErr
		return result, appErr
	}

	result.PayloadType = "friends.FriendRelationshipResult"
	result.Payload = friendsResult
	return result, nil
}

func friendsPublicErrorCode(result appfriends.FriendRelationshipResult, err error) appfriends.PublicErrorCode {
	if result.PublicErrorCode != "" {
		return result.PublicErrorCode
	}

	var serviceErr *appfriends.ServiceError
	if errors.As(err, &serviceErr) && serviceErr.PublicCode != "" {
		return serviceErr.PublicCode
	}

	return appfriends.PublicErrorFriendshipUnavailable
}

func friendsApplicationErrorResult(request app.RouteRequest, publicCode appfriends.PublicErrorCode) (app.ApplicationResult, error) {
	result := resultForRequest(request)
	appErr := friendsApplicationError(request.Route, publicCode)
	result.Error = appErr
	return result, appErr
}

func friendsApplicationError(route app.RouteKey, publicCode appfriends.PublicErrorCode) *app.ApplicationError {
	code := app.ErrorCode(publicCode)
	if code == "" {
		code = app.ErrorCode(appfriends.PublicErrorFriendshipUnavailable)
	}
	return &app.ApplicationError{
		Code:    code,
		Message: friendsErrorMessage(publicCode),
		Route:   route,
	}
}

func friendsErrorMessage(code appfriends.PublicErrorCode) string {
	switch code {
	case appfriends.PublicErrorFriendshipInvalidRequest:
		return "friendship request is invalid"
	case appfriends.PublicErrorFriendshipUnauthenticated:
		return "friendship request is unauthenticated"
	case appfriends.PublicErrorFriendshipForbidden:
		return "friendship request is forbidden"
	case appfriends.PublicErrorFriendshipTargetNotFound:
		return "friendship target was not found"
	case appfriends.PublicErrorFriendshipRelationshipNotFound:
		return "friendship relationship was not found"
	case appfriends.PublicErrorFriendshipDuplicateRequest:
		return "friendship request already exists"
	case appfriends.PublicErrorFriendshipAlreadyFriends:
		return "players are already friends"
	case appfriends.PublicErrorFriendshipBlockedRelationship:
		return "friendship relationship is blocked"
	case appfriends.PublicErrorFriendshipInvalidTransition:
		return "friendship transition is invalid"
	case appfriends.PublicErrorFriendshipVersionMismatch:
		return "friendship version mismatch"
	default:
		return "friendship service is unavailable"
	}
}
