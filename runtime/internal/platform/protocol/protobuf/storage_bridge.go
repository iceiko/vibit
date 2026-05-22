package protobuf

import (
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
	storagev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/storage/v1"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
	"google.golang.org/protobuf/proto"
)

func routeRequestWithStoragePayload(request app.RouteRequest) (app.RouteRequest, bool, error) {
	switch request.Route {
	case appstorage.GetOwnStorageObjectRoute():
		payload, ok := request.Payload.(*storagev1.GetOwnStorageObjectRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.storage.v1.GetOwnStorageObjectRequest")
		}
		request.Payload = appstorage.GetOwnStorageObjectRequest{
			Collection: payload.GetCollection(),
			Key:        payload.GetKey(),
		}
		return request, true, nil

	case appstorage.ListOwnStorageObjectsRoute():
		payload, ok := request.Payload.(*storagev1.ListOwnStorageObjectsRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.storage.v1.ListOwnStorageObjectsRequest")
		}
		request.Payload = appstorage.ListOwnStorageObjectsRequest{
			Collection:     payload.GetCollection(),
			Limit:          int(payload.GetLimit()),
			AfterObjectKey: payload.GetAfterKey(),
		}
		return request, true, nil

	case appstorage.PutOwnStorageObjectRoute():
		payload, ok := request.Payload.(*storagev1.PutOwnStorageObjectRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.storage.v1.PutOwnStorageObjectRequest")
		}
		request.Payload = appstorage.PutOwnStorageObjectRequest{
			Collection:      payload.GetCollection(),
			Key:             payload.GetKey(),
			ValueJSON:       []byte(payload.GetValueJson()),
			ExpectedVersion: storageVersionFromOptionalInt64(payload.ExpectedVersion),
		}
		return request, true, nil

	case appstorage.DeleteOwnStorageObjectRoute():
		payload, ok := request.Payload.(*storagev1.DeleteOwnStorageObjectRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.storage.v1.DeleteOwnStorageObjectRequest")
		}
		request.Payload = appstorage.DeleteOwnStorageObjectRequest{
			Collection:      payload.GetCollection(),
			Key:             payload.GetKey(),
			ExpectedVersion: storageVersionFromOptionalInt64(payload.ExpectedVersion),
		}
		return request, true, nil

	default:
		return request, false, nil
	}
}

func protoPayloadFromStorageRoute(route app.RouteKey, payload any) (proto.Message, bool, error) {
	switch route {
	case appstorage.GetOwnStorageObjectRoute():
		result, ok := storageResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be storage.StorageObjectResult")
		}
		return &storagev1.GetOwnStorageObjectResponse{
			Object: protoStorageObject(result.Object),
		}, true, nil

	case appstorage.ListOwnStorageObjectsRoute():
		result, ok := storageResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be storage.StorageObjectResult")
		}
		objects := make([]*storagev1.StorageObject, 0, len(result.Objects))
		for _, object := range result.Objects {
			objects = append(objects, protoStorageObject(object))
		}
		return &storagev1.ListOwnStorageObjectsResponse{
			Objects: objects,
			NextKey: result.NextObjectKey,
		}, true, nil

	case appstorage.PutOwnStorageObjectRoute():
		result, ok := storageResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be storage.StorageObjectResult")
		}
		return &storagev1.PutOwnStorageObjectResponse{
			Object:  protoStorageObject(result.Object),
			Version: int64(result.Version),
		}, true, nil

	case appstorage.DeleteOwnStorageObjectRoute():
		result, ok := storageResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be storage.StorageObjectResult")
		}
		return &storagev1.DeleteOwnStorageObjectResponse{
			Deleted: result.Status == appstorage.StorageObjectOperationStatusDeleted,
			Version: int64(result.Version),
		}, true, nil

	default:
		return nil, false, nil
	}
}

func storageResultPayload(payload any) (appstorage.StorageObjectResult, bool) {
	result, ok := payload.(appstorage.StorageObjectResult)
	if !ok {
		if pointerResult, pointerOK := payload.(*appstorage.StorageObjectResult); pointerOK && pointerResult != nil {
			result = *pointerResult
			ok = true
		}
	}
	return result, ok
}

func storageVersionFromOptionalInt64(value *int64) *storagemodule.StorageObjectVersion {
	if value == nil {
		return nil
	}
	version := storagemodule.StorageObjectVersion(*value)
	return &version
}

func protoStorageObject(object storagemodule.StorageObject) *storagev1.StorageObject {
	return &storagev1.StorageObject{
		Collection: object.Identity.Collection,
		Key:        object.Identity.Key,
		ValueJson:  string(object.Value.JSON),
		Version:    int64(object.Version),
		CreatedAt:  formatStorageTime(object.CreatedAt),
		UpdatedAt:  formatStorageTime(object.UpdatedAt),
	}
}

func formatStorageTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
