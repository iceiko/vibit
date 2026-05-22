package bootstrap

import (
	"context"
	"errors"

	"github.com/iceiko/vibit/runtime/internal/app"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
)

type StorageService interface {
	GetOwnStorageObject(context.Context, appstorage.GetOwnStorageObjectRequest) (appstorage.StorageObjectResult, error)
	ListOwnStorageObjects(context.Context, appstorage.ListOwnStorageObjectsRequest) (appstorage.StorageObjectResult, error)
	PutOwnStorageObject(context.Context, appstorage.PutOwnStorageObjectRequest) (appstorage.StorageObjectResult, error)
	DeleteOwnStorageObject(context.Context, appstorage.DeleteOwnStorageObjectRequest) (appstorage.StorageObjectResult, error)
}

type StorageRouteHandlers struct {
	Service StorageService
}

func (h StorageRouteHandlers) RegisterRoutes(dispatcher *app.Dispatcher) error {
	if dispatcher == nil {
		return errors.New("storage bootstrap: dispatcher is nil")
	}
	if err := dispatcher.Register(appstorage.GetOwnStorageObjectRoute(), app.HandlerFunc(h.HandleGetOwnStorageObjectRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appstorage.ListOwnStorageObjectsRoute(), app.HandlerFunc(h.HandleListOwnStorageObjectsRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appstorage.PutOwnStorageObjectRoute(), app.HandlerFunc(h.HandlePutOwnStorageObjectRoute)); err != nil {
		return err
	}
	return dispatcher.Register(appstorage.DeleteOwnStorageObjectRoute(), app.HandlerFunc(h.HandleDeleteOwnStorageObjectRoute))
}

func (h StorageRouteHandlers) HandleGetOwnStorageObjectRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectUnavailable)
	}

	payload, ok := request.Payload.(appstorage.GetOwnStorageObjectRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*appstorage.GetOwnStorageObjectRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.GetOwnStorageObject(ctx, payload)
	return storageResultForRequest(request, result, err)
}

func (h StorageRouteHandlers) HandleListOwnStorageObjectsRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectUnavailable)
	}

	payload, ok := request.Payload.(appstorage.ListOwnStorageObjectsRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*appstorage.ListOwnStorageObjectsRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.ListOwnStorageObjects(ctx, payload)
	return storageResultForRequest(request, result, err)
}

func (h StorageRouteHandlers) HandlePutOwnStorageObjectRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectUnavailable)
	}

	payload, ok := request.Payload.(appstorage.PutOwnStorageObjectRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*appstorage.PutOwnStorageObjectRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.PutOwnStorageObject(ctx, payload)
	return storageResultForRequest(request, result, err)
}

func (h StorageRouteHandlers) HandleDeleteOwnStorageObjectRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectUnavailable)
	}

	payload, ok := request.Payload.(appstorage.DeleteOwnStorageObjectRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*appstorage.DeleteOwnStorageObjectRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		return storageApplicationErrorResult(request, appstorage.PublicErrorStorageObjectInvalidRequest)
	}

	payload.Identity = request.Identity
	result, err := service.DeleteOwnStorageObject(ctx, payload)
	return storageResultForRequest(request, result, err)
}

func storageResultForRequest(request app.RouteRequest, storageResult appstorage.StorageObjectResult, err error) (app.ApplicationResult, error) {
	result := resultForRequest(request)
	if err != nil {
		appErr := storageApplicationError(request.Route, storagePublicErrorCode(storageResult, err))
		result.Error = appErr
		return result, appErr
	}

	result.PayloadType = "storage.StorageObjectResult"
	result.Payload = storageResult
	return result, nil
}

func storagePublicErrorCode(result appstorage.StorageObjectResult, err error) appstorage.PublicErrorCode {
	if result.PublicErrorCode != "" {
		return result.PublicErrorCode
	}

	var serviceErr *appstorage.ServiceError
	if errors.As(err, &serviceErr) && serviceErr.PublicCode != "" {
		return serviceErr.PublicCode
	}

	return appstorage.PublicErrorStorageObjectUnavailable
}

func storageApplicationErrorResult(request app.RouteRequest, publicCode appstorage.PublicErrorCode) (app.ApplicationResult, error) {
	result := resultForRequest(request)
	appErr := storageApplicationError(request.Route, publicCode)
	result.Error = appErr
	return result, appErr
}

func storageApplicationError(route app.RouteKey, publicCode appstorage.PublicErrorCode) *app.ApplicationError {
	code := app.ErrorCode(publicCode)
	if code == "" {
		code = app.ErrorCode(appstorage.PublicErrorStorageObjectUnavailable)
	}
	return &app.ApplicationError{
		Code:    code,
		Message: storageErrorMessage(publicCode),
		Route:   route,
	}
}

func storageErrorMessage(code appstorage.PublicErrorCode) string {
	switch code {
	case appstorage.PublicErrorStorageObjectInvalidRequest:
		return "storage object request is invalid"
	case appstorage.PublicErrorStorageObjectNotFound:
		return "storage object was not found"
	case appstorage.PublicErrorStorageObjectAlreadyExists:
		return "storage object already exists"
	case appstorage.PublicErrorStorageObjectVersionMismatch:
		return "storage object version mismatch"
	case appstorage.PublicErrorStorageObjectForbidden:
		return "storage object request is forbidden"
	default:
		return "storage object service is unavailable"
	}
}
