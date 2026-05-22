package bootstrap

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
)

func TestStorageRouteHandlersRegisterOwnObjectRoutesAndPassValidatedIdentity(t *testing.T) {
	dispatcher := app.NewDispatcher()
	service := &recordingStorageService{}
	handlers := StorageRouteHandlers{Service: service}

	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}

	identity := app.ValidatedPlayerIdentity("player-1", app.Session{
		ConnectionID:    "connection-1",
		SessionID:       "session-1",
		PlayerID:        "player-1",
		ConnectionEpoch: 7,
	})
	expectedVersion := storagemodule.StorageObjectVersion(3)

	dispatchStorageRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-get",
		Route:     appstorage.GetOwnStorageObjectRoute(),
		Identity:  identity,
		Payload: appstorage.GetOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
		},
	})
	if !service.getCalled || service.getRequest.Identity != identity ||
		service.getRequest.Collection != "progress" ||
		service.getRequest.Key != "tutorial" {
		t.Fatalf("get request = %#v, want validated identity and collection/key", service.getRequest)
	}

	dispatchStorageRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-list",
		Route:     appstorage.ListOwnStorageObjectsRoute(),
		Identity:  identity,
		Payload: appstorage.ListOwnStorageObjectsRequest{
			Collection:     "progress",
			Limit:          25,
			AfterObjectKey: "tutorial",
		},
	})
	if !service.listCalled || service.listRequest.Identity != identity ||
		service.listRequest.Limit != 25 ||
		service.listRequest.AfterObjectKey != "tutorial" {
		t.Fatalf("list request = %#v, want validated identity and pagination", service.listRequest)
	}

	dispatchStorageRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-put",
		Route:     appstorage.PutOwnStorageObjectRoute(),
		Identity:  identity,
		Payload: appstorage.PutOwnStorageObjectRequest{
			Collection:      "progress",
			Key:             "tutorial",
			ValueJSON:       []byte(`{"level":4}`),
			ExpectedVersion: &expectedVersion,
		},
	})
	if !service.putCalled || service.putRequest.Identity != identity ||
		string(service.putRequest.ValueJSON) != `{"level":4}` ||
		service.putRequest.ExpectedVersion == nil ||
		*service.putRequest.ExpectedVersion != expectedVersion {
		t.Fatalf("put request = %#v, want validated identity, JSON value, and expected version", service.putRequest)
	}

	dispatchStorageRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-delete",
		Route:     appstorage.DeleteOwnStorageObjectRoute(),
		Identity:  identity,
		Payload: appstorage.DeleteOwnStorageObjectRequest{
			Collection:      "progress",
			Key:             "tutorial",
			ExpectedVersion: &expectedVersion,
		},
	})
	if !service.deleteCalled || service.deleteRequest.Identity != identity ||
		service.deleteRequest.ExpectedVersion == nil ||
		*service.deleteRequest.ExpectedVersion != expectedVersion {
		t.Fatalf("delete request = %#v, want validated identity and expected version", service.deleteRequest)
	}
}

func TestStorageRouteHandlersRejectMalformedPayloadBeforeService(t *testing.T) {
	service := &recordingStorageService{}
	handlers := StorageRouteHandlers{Service: service}

	result, err := handlers.HandlePutOwnStorageObjectRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appstorage.PutOwnStorageObjectRoute(),
		Payload:   "not a storage request",
	})
	if err == nil {
		t.Fatal("HandlePutOwnStorageObjectRoute() error = nil, want invalid request error")
	}
	if service.putCalled {
		t.Fatal("storage service was called for malformed payload")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appstorage.PublicErrorStorageObjectInvalidRequest) {
		t.Fatalf("result error = %#v, want storage invalid request", result.Error)
	}
}

func TestStorageRouteHandlersMapServiceErrorsWithoutValueLeak(t *testing.T) {
	secretDetail := `internal sql value_json={"secret":true}`
	service := &recordingStorageService{
		putResult: appstorage.StorageObjectResult{
			Status:          appstorage.StorageObjectOperationStatusRejected,
			PublicErrorCode: appstorage.PublicErrorStorageObjectVersionMismatch,
			FailureClass:    appstorage.FailureClassVersionMismatch,
		},
		putErr: &appstorage.ServiceError{
			Operation:  appstorage.OperationPutOwnStorageObject,
			Class:      appstorage.FailureClassVersionMismatch,
			PublicCode: appstorage.PublicErrorStorageObjectVersionMismatch,
			Err:        errors.New(secretDetail),
		},
	}
	handlers := StorageRouteHandlers{Service: service}

	result, err := handlers.HandlePutOwnStorageObjectRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appstorage.PutOwnStorageObjectRoute(),
		Payload: appstorage.PutOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
			ValueJSON:  []byte(`{"secret":true}`),
		},
	})
	if err == nil {
		t.Fatal("HandlePutOwnStorageObjectRoute() error = nil, want public storage error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appstorage.PublicErrorStorageObjectVersionMismatch) {
		t.Fatalf("result error = %#v, want version mismatch code", result.Error)
	}
	errorText := result.Error.Error()
	if strings.Contains(errorText, secretDetail) ||
		strings.Contains(errorText, "value_json") ||
		strings.Contains(errorText, "secret") {
		t.Fatalf("error %q leaked storage details", errorText)
	}
	if result.Payload != nil {
		t.Fatalf("result Payload = %#v, want nil on storage failure", result.Payload)
	}
}

func TestStorageRouteHandlersRequireService(t *testing.T) {
	handlers := StorageRouteHandlers{}

	result, err := handlers.HandleGetOwnStorageObjectRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appstorage.GetOwnStorageObjectRoute(),
		Payload:   appstorage.GetOwnStorageObjectRequest{Collection: "progress", Key: "tutorial"},
	})
	if err == nil {
		t.Fatal("HandleGetOwnStorageObjectRoute() error = nil, want unavailable error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appstorage.PublicErrorStorageObjectUnavailable) {
		t.Fatalf("result error = %#v, want storage unavailable", result.Error)
	}
}

func dispatchStorageRoute(t *testing.T, dispatcher *app.Dispatcher, request app.RouteRequest) app.ApplicationResult {
	t.Helper()
	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch(%s) error = %v, want nil", app.RenderRouteKey(request.Route), err)
	}
	if result.PayloadType != "storage.StorageObjectResult" {
		t.Fatalf("PayloadType = %q, want storage.StorageObjectResult", result.PayloadType)
	}
	if _, ok := result.Payload.(appstorage.StorageObjectResult); !ok {
		t.Fatalf("Payload = %T, want storage.StorageObjectResult", result.Payload)
	}
	return result
}

type recordingStorageService struct {
	getCalled     bool
	getRequest    appstorage.GetOwnStorageObjectRequest
	getResult     appstorage.StorageObjectResult
	getErr        error
	listCalled    bool
	listRequest   appstorage.ListOwnStorageObjectsRequest
	listResult    appstorage.StorageObjectResult
	listErr       error
	putCalled     bool
	putRequest    appstorage.PutOwnStorageObjectRequest
	putResult     appstorage.StorageObjectResult
	putErr        error
	deleteCalled  bool
	deleteRequest appstorage.DeleteOwnStorageObjectRequest
	deleteResult  appstorage.StorageObjectResult
	deleteErr     error
}

func (s *recordingStorageService) GetOwnStorageObject(_ context.Context, request appstorage.GetOwnStorageObjectRequest) (appstorage.StorageObjectResult, error) {
	s.getCalled = true
	s.getRequest = request
	if s.getErr != nil {
		return s.getResult, s.getErr
	}
	if s.getResult.Status != "" {
		return s.getResult, nil
	}
	return storageHandlerResult(appstorage.StorageObjectOperationStatusFound), nil
}

func (s *recordingStorageService) ListOwnStorageObjects(_ context.Context, request appstorage.ListOwnStorageObjectsRequest) (appstorage.StorageObjectResult, error) {
	s.listCalled = true
	s.listRequest = request
	if s.listErr != nil {
		return s.listResult, s.listErr
	}
	if s.listResult.Status != "" {
		return s.listResult, nil
	}
	return storageHandlerResult(appstorage.StorageObjectOperationStatusListed), nil
}

func (s *recordingStorageService) PutOwnStorageObject(_ context.Context, request appstorage.PutOwnStorageObjectRequest) (appstorage.StorageObjectResult, error) {
	s.putCalled = true
	s.putRequest = request
	if s.putErr != nil {
		return s.putResult, s.putErr
	}
	if s.putResult.Status != "" {
		return s.putResult, nil
	}
	return storageHandlerResult(appstorage.StorageObjectOperationStatusStored), nil
}

func (s *recordingStorageService) DeleteOwnStorageObject(_ context.Context, request appstorage.DeleteOwnStorageObjectRequest) (appstorage.StorageObjectResult, error) {
	s.deleteCalled = true
	s.deleteRequest = request
	if s.deleteErr != nil {
		return s.deleteResult, s.deleteErr
	}
	if s.deleteResult.Status != "" {
		return s.deleteResult, nil
	}
	return storageHandlerResult(appstorage.StorageObjectOperationStatusDeleted), nil
}

func storageHandlerResult(status appstorage.StorageObjectOperationStatus) appstorage.StorageObjectResult {
	object := storagemodule.StorageObject{
		ObjectID: "storage-object-1",
		Owner: storagemodule.StorageObjectOwner{
			Kind: storagemodule.OwnerKindPlayer,
			ID:   "player-1",
		},
		Identity: storagemodule.StorageObjectIdentity{
			Collection: "progress",
			Key:        "tutorial",
		},
		Value: storagemodule.StorageObjectValue{
			JSON: []byte(`{"level":4}`),
		},
		Version:   4,
		Status:    storagemodule.StorageObjectStatusActive,
		CreatedAt: time.Date(2026, 5, 22, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 5, 22, 10, 1, 0, 0, time.UTC),
	}
	if status == appstorage.StorageObjectOperationStatusDeleted {
		object.Status = storagemodule.StorageObjectStatusDeleted
	}
	result := appstorage.StorageObjectResult{
		Status:  status,
		Object:  object,
		Objects: []storagemodule.StorageObject{object},
		Version: object.Version,
	}
	if status == appstorage.StorageObjectOperationStatusListed {
		result.NextObjectKey = "tutorial-2"
	}
	return result
}
