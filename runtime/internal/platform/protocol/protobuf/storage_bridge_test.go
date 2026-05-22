package protobuf

import (
	"errors"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
	storagev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/storage/v1"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
	"google.golang.org/protobuf/proto"
)

func TestRouteRequestWithStoragePayloadMapsOwnObjectRequests(t *testing.T) {
	expectedVersion := int64(7)

	tests := []struct {
		name    string
		request app.RouteRequest
		assert  func(t *testing.T, payload any)
	}{
		{
			name: "get",
			request: app.RouteRequest{
				Route: appstorage.GetOwnStorageObjectRoute(),
				Payload: &storagev1.GetOwnStorageObjectRequest{
					Collection: "progress",
					Key:        "tutorial",
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appstorage.GetOwnStorageObjectRequest)
				if !ok {
					t.Fatalf("payload = %T, want GetOwnStorageObjectRequest", payload)
				}
				if got.Collection != "progress" || got.Key != "tutorial" {
					t.Fatalf("payload = %#v, want mapped get request", got)
				}
			},
		},
		{
			name: "list",
			request: app.RouteRequest{
				Route: appstorage.ListOwnStorageObjectsRoute(),
				Payload: &storagev1.ListOwnStorageObjectsRequest{
					Collection: "progress",
					Limit:      25,
					AfterKey:   "tutorial",
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appstorage.ListOwnStorageObjectsRequest)
				if !ok {
					t.Fatalf("payload = %T, want ListOwnStorageObjectsRequest", payload)
				}
				if got.Collection != "progress" || got.Limit != 25 || got.AfterObjectKey != "tutorial" {
					t.Fatalf("payload = %#v, want mapped list request", got)
				}
			},
		},
		{
			name: "put",
			request: app.RouteRequest{
				Route: appstorage.PutOwnStorageObjectRoute(),
				Payload: &storagev1.PutOwnStorageObjectRequest{
					Collection:      "progress",
					Key:             "tutorial",
					ValueJson:       `{"level":4}`,
					ExpectedVersion: &expectedVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appstorage.PutOwnStorageObjectRequest)
				if !ok {
					t.Fatalf("payload = %T, want PutOwnStorageObjectRequest", payload)
				}
				if got.Collection != "progress" ||
					got.Key != "tutorial" ||
					string(got.ValueJSON) != `{"level":4}` ||
					got.ExpectedVersion == nil ||
					*got.ExpectedVersion != storagemodule.StorageObjectVersion(expectedVersion) {
					t.Fatalf("payload = %#v, want mapped put request", got)
				}
			},
		},
		{
			name: "delete",
			request: app.RouteRequest{
				Route: appstorage.DeleteOwnStorageObjectRoute(),
				Payload: &storagev1.DeleteOwnStorageObjectRequest{
					Collection:      "progress",
					Key:             "tutorial",
					ExpectedVersion: &expectedVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appstorage.DeleteOwnStorageObjectRequest)
				if !ok {
					t.Fatalf("payload = %T, want DeleteOwnStorageObjectRequest", payload)
				}
				if got.Collection != "progress" ||
					got.Key != "tutorial" ||
					got.ExpectedVersion == nil ||
					*got.ExpectedVersion != storagemodule.StorageObjectVersion(expectedVersion) {
					t.Fatalf("payload = %#v, want mapped delete request", got)
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

func TestRouteRequestWithStoragePayloadPreservesMissingExpectedVersion(t *testing.T) {
	mapped, err := RouteRequestWithDomainPayload(app.RouteRequest{
		Route: appstorage.PutOwnStorageObjectRoute(),
		Payload: &storagev1.PutOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
			ValueJson:  `{"level":4}`,
		},
	})
	if err != nil {
		t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
	}
	payload, ok := mapped.Payload.(appstorage.PutOwnStorageObjectRequest)
	if !ok {
		t.Fatalf("Payload = %T, want PutOwnStorageObjectRequest", mapped.Payload)
	}
	if payload.ExpectedVersion != nil {
		t.Fatalf("ExpectedVersion = %#v, want nil when optional field is absent", payload.ExpectedVersion)
	}
}

func TestRouteRequestWithStoragePayloadRejectsWrongPayload(t *testing.T) {
	_, err := RouteRequestWithDomainPayload(app.RouteRequest{
		Route:   appstorage.GetOwnStorageObjectRoute(),
		Payload: &storagev1.ListOwnStorageObjectsRequest{},
	})
	if err == nil {
		t.Fatal("RouteRequestWithDomainPayload() error = nil, want bridge error")
	}
	var bridgeErr *PayloadBridgeError
	if !errors.As(err, &bridgeErr) {
		t.Fatalf("error = %T %v, want *PayloadBridgeError", err, err)
	}
}

func TestProtoPayloadFromStorageResultMapsResponses(t *testing.T) {
	createdAt := time.Date(2026, 5, 22, 8, 30, 0, 123, time.FixedZone("test", 8*60*60))
	updatedAt := createdAt.Add(time.Minute)
	first := storageBridgeObject("progress", "tutorial", []byte(`{"level":4}`), 7, createdAt, updatedAt)
	second := storageBridgeObject("progress", "boss", []byte(`{"clear":true}`), 8, createdAt.Add(time.Hour), updatedAt.Add(time.Hour))

	getPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appstorage.GetOwnStorageObjectRoute(),
		Payload: appstorage.StorageObjectResult{
			Status:  appstorage.StorageObjectOperationStatusFound,
			Object:  first,
			Version: first.Version,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(get) error = %v, want nil", err)
	}
	getResponse, ok := getPayload.(*storagev1.GetOwnStorageObjectResponse)
	if !ok {
		t.Fatalf("get payload = %T, want GetOwnStorageObjectResponse", getPayload)
	}
	assertProtoStorageObject(t, getResponse.GetObject(), first)

	listPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appstorage.ListOwnStorageObjectsRoute(),
		Payload: appstorage.StorageObjectResult{
			Status:        appstorage.StorageObjectOperationStatusListed,
			Objects:       []storagemodule.StorageObject{first, second},
			NextObjectKey: "next-key",
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(list) error = %v, want nil", err)
	}
	listResponse, ok := listPayload.(*storagev1.ListOwnStorageObjectsResponse)
	if !ok {
		t.Fatalf("list payload = %T, want ListOwnStorageObjectsResponse", listPayload)
	}
	if len(listResponse.GetObjects()) != 2 || listResponse.GetNextKey() != "next-key" {
		t.Fatalf("list response = %#v, want two objects and next key", listResponse)
	}
	assertProtoStorageObject(t, listResponse.GetObjects()[1], second)

	putPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appstorage.PutOwnStorageObjectRoute(),
		Payload: &appstorage.StorageObjectResult{
			Status:  appstorage.StorageObjectOperationStatusStored,
			Object:  first,
			Version: first.Version,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(put) error = %v, want nil", err)
	}
	putResponse, ok := putPayload.(*storagev1.PutOwnStorageObjectResponse)
	if !ok {
		t.Fatalf("put payload = %T, want PutOwnStorageObjectResponse", putPayload)
	}
	if putResponse.GetVersion() != int64(first.Version) {
		t.Fatalf("put version = %d, want %d", putResponse.GetVersion(), first.Version)
	}
	assertProtoStorageObject(t, putResponse.GetObject(), first)

	deletePayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appstorage.DeleteOwnStorageObjectRoute(),
		Payload: appstorage.StorageObjectResult{
			Status:  appstorage.StorageObjectOperationStatusDeleted,
			Object:  second,
			Version: second.Version,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(delete) error = %v, want nil", err)
	}
	deleteResponse, ok := deletePayload.(*storagev1.DeleteOwnStorageObjectResponse)
	if !ok {
		t.Fatalf("delete payload = %T, want DeleteOwnStorageObjectResponse", deletePayload)
	}
	if !deleteResponse.GetDeleted() || deleteResponse.GetVersion() != int64(second.Version) {
		t.Fatalf("delete response = %#v, want deleted true and version", deleteResponse)
	}
}

func TestStorageProtoPayloadShapeOmitsOwnerAndTransportFields(t *testing.T) {
	messages := []proto.Message{
		&storagev1.StorageObject{},
		&storagev1.GetOwnStorageObjectRequest{},
		&storagev1.ListOwnStorageObjectsRequest{},
		&storagev1.PutOwnStorageObjectRequest{},
		&storagev1.DeleteOwnStorageObjectRequest{},
	}
	for _, message := range messages {
		descriptor := message.ProtoReflect().Descriptor()
		fields := descriptor.Fields()
		fieldNames := map[string]bool{}
		for i := 0; i < fields.Len(); i++ {
			fieldNames[string(fields.Get(i).Name())] = true
		}
		for _, forbidden := range []string{
			"owner_id",
			"player_id",
			"session_id",
			"access_token",
			"credential_proof",
			"lookup_digest",
			"verifier_digest",
			"s3_bucket",
		} {
			if fieldNames[forbidden] {
				t.Fatalf("%s has forbidden field %q", descriptor.FullName(), forbidden)
			}
		}
	}
}

func storageBridgeObject(collection, key string, value []byte, version storagemodule.StorageObjectVersion, createdAt, updatedAt time.Time) storagemodule.StorageObject {
	return storagemodule.StorageObject{
		ObjectID: "storage-object-1",
		Owner: storagemodule.StorageObjectOwner{
			Kind: storagemodule.OwnerKindPlayer,
			ID:   "player-1",
		},
		Identity: storagemodule.StorageObjectIdentity{
			Collection: collection,
			Key:        key,
		},
		Value: storagemodule.StorageObjectValue{
			JSON: append([]byte(nil), value...),
		},
		Version:   version,
		Status:    storagemodule.StorageObjectStatusActive,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func assertProtoStorageObject(t *testing.T, got *storagev1.StorageObject, want storagemodule.StorageObject) {
	t.Helper()
	if got == nil {
		t.Fatal("StorageObject = nil, want mapped object")
	}
	if got.GetCollection() != want.Identity.Collection ||
		got.GetKey() != want.Identity.Key ||
		got.GetValueJson() != string(want.Value.JSON) ||
		got.GetVersion() != int64(want.Version) ||
		got.GetCreatedAt() != want.CreatedAt.UTC().Format(time.RFC3339Nano) ||
		got.GetUpdatedAt() != want.UpdatedAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("StorageObject = %#v, want mapped object %#v", got, want)
	}
}
