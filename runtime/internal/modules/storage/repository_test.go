package storage

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRepositoryInterfaceIsStorageNeutral(t *testing.T) {
	var _ Repository = recordingRepository{}
}

func TestOwnerKindAndObjectStatusAreClosedSets(t *testing.T) {
	if !OwnerKindPlayer.IsValid() {
		t.Fatalf("%q IsValid() = false, want true", OwnerKindPlayer)
	}
	if OwnerKind("global").IsValid() {
		t.Fatal(`OwnerKind("global").IsValid() = true, want false`)
	}

	for _, status := range []StorageObjectStatus{StorageObjectStatusActive, StorageObjectStatusDeleted} {
		if !status.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", status)
		}
	}
	if StorageObjectStatus("archived").IsValid() {
		t.Fatal(`StorageObjectStatus("archived").IsValid() = true, want false`)
	}
}

func TestNormalizeStorageObjectRecordTrimsAndCopiesValue(t *testing.T) {
	now := time.Date(2026, 5, 22, 9, 0, 0, 0, time.FixedZone("test", 8*60*60))
	deletedAt := now.Add(time.Minute)
	value := []byte(`{"level":3}`)

	record, err := NormalizeStorageObjectRecord(StorageObject{
		ObjectID:  " object-1 ",
		Owner:     StorageObjectOwner{Kind: OwnerKind(" player "), ID: " player-1 "},
		Identity:  StorageObjectIdentity{Collection: " loadouts ", Key: " primary "},
		Value:     StorageObjectValue{JSON: value},
		Version:   StorageObjectVersion(7),
		Status:    StorageObjectStatusDeleted,
		CreatedAt: now,
		UpdatedAt: now.Add(2 * time.Minute),
		DeletedAt: &deletedAt,
	})
	if err != nil {
		t.Fatalf("NormalizeStorageObjectRecord() error = %v, want nil", err)
	}

	if record.ObjectID != "object-1" ||
		record.Owner.Kind != OwnerKindPlayer ||
		record.Owner.ID != "player-1" ||
		record.Identity.Collection != "loadouts" ||
		record.Identity.Key != "primary" ||
		record.Status != StorageObjectStatusDeleted {
		t.Fatalf("normalized record = %#v, want trimmed fields", record)
	}
	if record.CreatedAt.Location() != time.UTC ||
		record.UpdatedAt.Location() != time.UTC ||
		record.DeletedAt == nil ||
		record.DeletedAt.Location() != time.UTC {
		t.Fatalf("normalized record times = %#v, want UTC times", record)
	}

	value[0] = '['
	if string(record.Value.JSON) != `{"level":3}` {
		t.Fatalf("record.Value.JSON = %q, want copied JSON value", string(record.Value.JSON))
	}
}

func TestNormalizeStorageObjectRecordRejectsInvalidShape(t *testing.T) {
	valid := validStorageObjectRecord()
	tests := []struct {
		name   string
		mutate func(*StorageObject)
	}{
		{name: "object_id", mutate: func(r *StorageObject) { r.ObjectID = " " }},
		{name: "owner_kind", mutate: func(r *StorageObject) { r.Owner.Kind = "global" }},
		{name: "owner_id", mutate: func(r *StorageObject) { r.Owner.ID = " " }},
		{name: "collection", mutate: func(r *StorageObject) { r.Identity.Collection = " " }},
		{name: "object_key", mutate: func(r *StorageObject) { r.Identity.Key = " " }},
		{name: "value_json_missing", mutate: func(r *StorageObject) { r.Value.JSON = nil }},
		{name: "value_json_not_object", mutate: func(r *StorageObject) { r.Value.JSON = []byte(`[]`) }},
		{name: "version", mutate: func(r *StorageObject) { r.Version = 0 }},
		{name: "status", mutate: func(r *StorageObject) { r.Status = "archived" }},
		{name: "created_at", mutate: func(r *StorageObject) { r.CreatedAt = time.Time{} }},
		{name: "updated_at", mutate: func(r *StorageObject) { r.UpdatedAt = time.Time{} }},
		{name: "updated_before_created", mutate: func(r *StorageObject) { r.UpdatedAt = r.CreatedAt.Add(-time.Second) }},
		{name: "deleted_status_without_deleted_at", mutate: func(r *StorageObject) { r.Status = StorageObjectStatusDeleted; r.DeletedAt = nil }},
		{name: "deleted_at_without_deleted_status", mutate: func(r *StorageObject) { at := r.CreatedAt; r.DeletedAt = &at }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := valid
			tt.mutate(&record)
			if _, err := NormalizeStorageObjectRecord(record); err == nil {
				t.Fatal("NormalizeStorageObjectRecord() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeCreateStorageObjectInputDefaultsActiveVersionAndCopiesValue(t *testing.T) {
	value := []byte(`{"xp":10}`)

	input, err := NormalizeCreateStorageObjectInput(CreateStorageObjectInput{
		ObjectID:    " object-1 ",
		Owner:       StorageObjectOwner{Kind: OwnerKind(" player "), ID: " player-1 "},
		Identity:    StorageObjectIdentity{Collection: " progress ", Key: " tutorial "},
		Value:       StorageObjectValue{JSON: value},
		RequestedBy: " storage_service ",
	})
	if err != nil {
		t.Fatalf("NormalizeCreateStorageObjectInput() error = %v, want nil", err)
	}

	if input.ObjectID != "object-1" ||
		input.Owner.Kind != OwnerKindPlayer ||
		input.Owner.ID != "player-1" ||
		input.Identity.Collection != "progress" ||
		input.Identity.Key != "tutorial" ||
		input.InitialVersion != InitialStorageObjectVersion ||
		input.Status != StorageObjectStatusActive ||
		input.RequestedBy != "storage_service" {
		t.Fatalf("normalized input = %#v, want trimmed active create input", input)
	}
	value[0] = '['
	if string(input.Value.JSON) != `{"xp":10}` {
		t.Fatalf("input.Value.JSON = %q, want copied JSON value", string(input.Value.JSON))
	}
}

func TestNormalizeCreateStorageObjectInputRejectsInvalidShape(t *testing.T) {
	valid := validCreateStorageObjectInput()
	tests := []struct {
		name   string
		mutate func(*CreateStorageObjectInput)
	}{
		{name: "object_id", mutate: func(i *CreateStorageObjectInput) { i.ObjectID = " " }},
		{name: "owner", mutate: func(i *CreateStorageObjectInput) { i.Owner.ID = " " }},
		{name: "identity", mutate: func(i *CreateStorageObjectInput) { i.Identity.Collection = strings.Repeat("a", MaxCollectionLength+1) }},
		{name: "value", mutate: func(i *CreateStorageObjectInput) { i.Value.JSON = []byte(`"scalar"`) }},
		{name: "initial_version", mutate: func(i *CreateStorageObjectInput) { i.InitialVersion = -1 }},
		{name: "status", mutate: func(i *CreateStorageObjectInput) { i.Status = StorageObjectStatusDeleted }},
		{name: "requested_by", mutate: func(i *CreateStorageObjectInput) { i.RequestedBy = " " }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := valid
			tt.mutate(&input)
			if _, err := NormalizeCreateStorageObjectInput(input); err == nil {
				t.Fatal("NormalizeCreateStorageObjectInput() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeGetListUpdateDeleteInputs(t *testing.T) {
	get, err := NormalizeGetStorageObjectInput(GetStorageObjectInput{
		Owner:    StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Identity: StorageObjectIdentity{Collection: " profile ", Key: " public "},
	})
	if err != nil {
		t.Fatalf("NormalizeGetStorageObjectInput() error = %v, want nil", err)
	}
	if get.Owner.Kind != OwnerKindPlayer || get.Owner.ID != "player-1" || get.Identity.Collection != "profile" || get.Identity.Key != "public" {
		t.Fatalf("get input = %#v, want trimmed owner identity", get)
	}

	list, err := NormalizeListStorageObjectsInput(ListStorageObjectsInput{
		Owner:          StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Collection:     " profile ",
		Limit:          0,
		AfterObjectKey: " public ",
	})
	if err != nil {
		t.Fatalf("NormalizeListStorageObjectsInput() error = %v, want nil", err)
	}
	if list.Collection != "profile" || list.Limit != DefaultListStorageObjectsLimit || list.AfterObjectKey != "public" {
		t.Fatalf("list input = %#v, want trimmed collection, default limit, and cursor", list)
	}

	expected := StorageObjectVersion(3)
	replacement := []byte(`{"class":"mage"}`)
	update, err := NormalizeUpdateStorageObjectInput(UpdateStorageObjectInput{
		Owner:           StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Identity:        StorageObjectIdentity{Collection: " loadouts ", Key: " primary "},
		Value:           StorageObjectValue{JSON: replacement},
		ExpectedVersion: &expected,
		RequestedBy:     " storage_service ",
	})
	if err != nil {
		t.Fatalf("NormalizeUpdateStorageObjectInput() error = %v, want nil", err)
	}
	if update.ExpectedVersion == nil || *update.ExpectedVersion != expected || update.RequestedBy != "storage_service" {
		t.Fatalf("update input = %#v, want expected version and requested_by", update)
	}
	replacement[0] = '['
	if string(update.Value.JSON) != `{"class":"mage"}` {
		t.Fatalf("update.Value.JSON = %q, want copied JSON value", string(update.Value.JSON))
	}

	deleteInput, err := NormalizeDeleteStorageObjectInput(DeleteStorageObjectInput{
		Owner:           StorageObjectOwner{Kind: " player ", ID: " player-1 "},
		Identity:        StorageObjectIdentity{Collection: " loadouts ", Key: " primary "},
		ExpectedVersion: &expected,
		RequestedBy:     " storage_service ",
	})
	if err != nil {
		t.Fatalf("NormalizeDeleteStorageObjectInput() error = %v, want nil", err)
	}
	if deleteInput.ExpectedVersion == nil || *deleteInput.ExpectedVersion != expected || deleteInput.RequestedBy != "storage_service" {
		t.Fatalf("delete input = %#v, want expected version and requested_by", deleteInput)
	}
}

func TestNormalizeListStorageObjectsInputRejectsInvalidPagination(t *testing.T) {
	valid := ListStorageObjectsInput{
		Owner:      StorageObjectOwner{Kind: OwnerKindPlayer, ID: "player-1"},
		Collection: "profile",
		Limit:      DefaultListStorageObjectsLimit,
	}
	tests := []struct {
		name   string
		mutate func(*ListStorageObjectsInput)
	}{
		{name: "missing_owner", mutate: func(i *ListStorageObjectsInput) { i.Owner.ID = " " }},
		{name: "missing_collection", mutate: func(i *ListStorageObjectsInput) { i.Collection = " " }},
		{name: "limit_negative", mutate: func(i *ListStorageObjectsInput) { i.Limit = -1 }},
		{name: "limit_too_large", mutate: func(i *ListStorageObjectsInput) { i.Limit = MaxListStorageObjectsLimit + 1 }},
		{name: "cursor_too_long", mutate: func(i *ListStorageObjectsInput) { i.AfterObjectKey = strings.Repeat("a", MaxObjectKeyLength+1) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := valid
			tt.mutate(&input)
			if _, err := NormalizeListStorageObjectsInput(input); err == nil {
				t.Fatal("NormalizeListStorageObjectsInput() error = nil, want rejection")
			}
		})
	}
}

func TestConflictAndRepositoryErrorsAreTypedAndRedacted(t *testing.T) {
	conflict := StorageObjectConflict{
		Class:          StorageObjectConflictVersionMismatch,
		Expected:       StorageObjectVersion(3),
		Actual:         StorageObjectVersion(4),
		Retryable:      true,
		RedactedReason: "version_mismatch",
	}
	if !conflict.Class.IsValid() {
		t.Fatalf("%q IsValid() = false, want true", conflict.Class)
	}
	if strings.Contains(conflict.Error(), "player-1") || strings.Contains(conflict.Error(), `{"`) {
		t.Fatalf("conflict error leaks storage material: %q", conflict.Error())
	}

	repositoryErr := &StorageObjectRepositoryError{
		Kind:           ErrStorageUnavailable,
		Conflict:       conflict,
		Operation:      "update",
		RedactedReason: "storage_unavailable",
		Err:            errors.New(`driver leaked {"secret":true}`),
	}
	if !errors.Is(repositoryErr, ErrStorageUnavailable) {
		t.Fatalf("errors.Is(repositoryErr, ErrStorageUnavailable) = false, want true")
	}
	if !errors.Is(repositoryErr, conflict) {
		t.Fatalf("errors.Is(repositoryErr, conflict) = false, want true")
	}
	if strings.Contains(repositoryErr.Error(), "secret") || strings.Contains(repositoryErr.Error(), "driver leaked") {
		t.Fatalf("repository error leaks wrapped detail: %q", repositoryErr.Error())
	}
}

func TestStorageObjectHasNoSecretTransportOrBlobFields(t *testing.T) {
	forbiddenFragments := []string{
		"Token",
		"Credential",
		"Verifier",
		"Digest",
		"Connection",
		"Subprotocol",
		"Remote",
		"Blob",
		"File",
		"Bucket",
		"S3",
	}
	objectType := reflect.TypeOf(StorageObject{})
	for i := 0; i < objectType.NumField(); i++ {
		fieldName := objectType.Field(i).Name
		for _, fragment := range forbiddenFragments {
			if strings.Contains(fieldName, fragment) {
				t.Fatalf("StorageObject field %q contains forbidden fragment %q", fieldName, fragment)
			}
		}
	}
}

func validStorageObjectRecord() StorageObject {
	now := time.Date(2026, 5, 22, 1, 2, 3, 0, time.UTC)
	return StorageObject{
		ObjectID:  "object-1",
		Owner:     StorageObjectOwner{Kind: OwnerKindPlayer, ID: "player-1"},
		Identity:  StorageObjectIdentity{Collection: "profile", Key: "public"},
		Value:     StorageObjectValue{JSON: []byte(`{"name":"p1"}`)},
		Version:   StorageObjectVersion(1),
		Status:    StorageObjectStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func validCreateStorageObjectInput() CreateStorageObjectInput {
	return CreateStorageObjectInput{
		ObjectID:    "object-1",
		Owner:       StorageObjectOwner{Kind: OwnerKindPlayer, ID: "player-1"},
		Identity:    StorageObjectIdentity{Collection: "profile", Key: "public"},
		Value:       StorageObjectValue{JSON: []byte(`{"name":"p1"}`)},
		RequestedBy: "storage_service",
	}
}

type recordingRepository struct{}

func (recordingRepository) CreateStorageObject(context.Context, CreateStorageObjectInput) (StorageObject, error) {
	return StorageObject{}, nil
}

func (recordingRepository) GetStorageObject(context.Context, GetStorageObjectInput) (StorageObject, error) {
	return StorageObject{}, nil
}

func (recordingRepository) ListStorageObjects(context.Context, ListStorageObjectsInput) (ListStorageObjectsResult, error) {
	return ListStorageObjectsResult{}, nil
}

func (recordingRepository) UpdateStorageObject(context.Context, UpdateStorageObjectInput) (StorageObject, error) {
	return StorageObject{}, nil
}

func (recordingRepository) DeleteStorageObject(context.Context, DeleteStorageObjectInput) (StorageObject, error) {
	return StorageObject{}, nil
}
