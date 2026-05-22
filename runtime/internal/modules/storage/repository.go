package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ModuleName = "storage"

	MaxCollectionLength            = 128
	MaxObjectKeyLength             = 256
	DefaultListStorageObjectsLimit = 100
	MaxListStorageObjectsLimit     = 500

	InitialStorageObjectVersion StorageObjectVersion = 1
)

type OwnerKind string

const (
	OwnerKindPlayer OwnerKind = "player"
)

func (k OwnerKind) IsValid() bool {
	return k == OwnerKindPlayer
}

type StorageObjectStatus string

const (
	StorageObjectStatusActive  StorageObjectStatus = "active"
	StorageObjectStatusDeleted StorageObjectStatus = "deleted"
)

func (s StorageObjectStatus) IsValid() bool {
	switch s {
	case StorageObjectStatusActive, StorageObjectStatusDeleted:
		return true
	default:
		return false
	}
}

type StorageObjectVersion int64

func (v StorageObjectVersion) IsValid() bool {
	return v > 0
}

type StorageObjectOwner struct {
	Kind OwnerKind
	ID   string
}

type StorageObjectIdentity struct {
	Collection string
	Key        string
}

type StorageObjectValue struct {
	JSON []byte
}

type StorageObject struct {
	ObjectID  string
	Owner     StorageObjectOwner
	Identity  StorageObjectIdentity
	Value     StorageObjectValue
	Version   StorageObjectVersion
	Status    StorageObjectStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CreateStorageObjectInput struct {
	ObjectID       string
	Owner          StorageObjectOwner
	Identity       StorageObjectIdentity
	Value          StorageObjectValue
	InitialVersion StorageObjectVersion
	Status         StorageObjectStatus
	RequestedBy    string
}

type GetStorageObjectInput struct {
	Owner    StorageObjectOwner
	Identity StorageObjectIdentity
}

type ListStorageObjectsInput struct {
	Owner          StorageObjectOwner
	Collection     string
	Limit          int
	AfterObjectKey string
}

type ListStorageObjectsResult struct {
	Objects       []StorageObject
	NextObjectKey string
}

type UpdateStorageObjectInput struct {
	Owner           StorageObjectOwner
	Identity        StorageObjectIdentity
	Value           StorageObjectValue
	ExpectedVersion *StorageObjectVersion
	RequestedBy     string
}

type DeleteStorageObjectInput struct {
	Owner           StorageObjectOwner
	Identity        StorageObjectIdentity
	ExpectedVersion *StorageObjectVersion
	RequestedBy     string
}

type Repository interface {
	CreateStorageObject(context.Context, CreateStorageObjectInput) (StorageObject, error)
	GetStorageObject(context.Context, GetStorageObjectInput) (StorageObject, error)
	ListStorageObjects(context.Context, ListStorageObjectsInput) (ListStorageObjectsResult, error)
	UpdateStorageObject(context.Context, UpdateStorageObjectInput) (StorageObject, error)
	DeleteStorageObject(context.Context, DeleteStorageObjectInput) (StorageObject, error)
}

type StorageObjectConflictClass string

const (
	StorageObjectConflictAlreadyExists          StorageObjectConflictClass = "object_already_exists"
	StorageObjectConflictNotFound               StorageObjectConflictClass = "object_not_found"
	StorageObjectConflictVersionMismatch        StorageObjectConflictClass = "version_mismatch"
	StorageObjectConflictInvalidExpectedVersion StorageObjectConflictClass = "invalid_expected_version"
	StorageObjectConflictOwnerScopeMismatch     StorageObjectConflictClass = "owner_scope_mismatch"
	StorageObjectConflictDeletedObject          StorageObjectConflictClass = "deleted_object"
	StorageObjectConflictStorageUnavailable     StorageObjectConflictClass = "storage_unavailable"
)

func (c StorageObjectConflictClass) IsValid() bool {
	switch c {
	case StorageObjectConflictAlreadyExists,
		StorageObjectConflictNotFound,
		StorageObjectConflictVersionMismatch,
		StorageObjectConflictInvalidExpectedVersion,
		StorageObjectConflictOwnerScopeMismatch,
		StorageObjectConflictDeletedObject,
		StorageObjectConflictStorageUnavailable:
		return true
	default:
		return false
	}
}

type StorageObjectConflict struct {
	Class          StorageObjectConflictClass
	Expected       StorageObjectVersion
	Actual         StorageObjectVersion
	Retryable      bool
	RedactedReason string
}

func (c StorageObjectConflict) Error() string {
	reason := strings.TrimSpace(c.RedactedReason)
	if reason == "" {
		reason = string(c.Class)
	}
	if reason == "" {
		return "storage object conflict"
	}
	return fmt.Sprintf("storage object conflict: %s", reason)
}

func (c StorageObjectConflict) Is(target error) bool {
	targetConflict, ok := target.(StorageObjectConflict)
	if !ok {
		return false
	}
	return c.Class == targetConflict.Class && c.Class != ""
}

var (
	ErrStorageObjectInvalidInput = errors.New("storage object repository: invalid input")
	ErrStorageObjectConflict     = errors.New("storage object repository: conflict")
	ErrStorageUnavailable        = errors.New("storage object repository: storage unavailable")
)

type StorageObjectRepositoryError struct {
	Kind           error
	Conflict       StorageObjectConflict
	Operation      string
	RedactedReason string
	Err            error
}

func (e *StorageObjectRepositoryError) Error() string {
	if e == nil {
		return ""
	}
	reason := strings.TrimSpace(e.RedactedReason)
	if reason == "" && e.Conflict.Class != "" {
		reason = string(e.Conflict.Class)
	}
	if reason == "" && e.Kind != nil {
		reason = e.Kind.Error()
	}
	operation := strings.TrimSpace(e.Operation)
	if operation == "" {
		operation = "operation"
	}
	if reason == "" {
		return fmt.Sprintf("storage object repository %s failed", operation)
	}
	return fmt.Sprintf("storage object repository %s failed: %s", operation, reason)
}

func (e *StorageObjectRepositoryError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *StorageObjectRepositoryError) Is(target error) bool {
	if e == nil {
		return false
	}
	if errors.Is(e.Kind, target) || errors.Is(e.Err, target) {
		return true
	}
	if e.Conflict.Class != "" && errors.Is(e.Conflict, target) {
		return true
	}
	return target == ErrStorageObjectConflict && e.Conflict.Class != ""
}

func NormalizeStorageObjectRecord(record StorageObject) (StorageObject, error) {
	var err error
	record.ObjectID, err = normalizeRequired("object_id", record.ObjectID)
	if err != nil {
		return StorageObject{}, err
	}
	record.Owner, err = NormalizeStorageObjectOwner(record.Owner)
	if err != nil {
		return StorageObject{}, err
	}
	record.Identity, err = NormalizeStorageObjectIdentity(record.Identity)
	if err != nil {
		return StorageObject{}, err
	}
	record.Value, err = NormalizeStorageObjectValue(record.Value)
	if err != nil {
		return StorageObject{}, err
	}
	if !record.Version.IsValid() {
		return StorageObject{}, errors.New("storage object: version must be positive")
	}
	if !record.Status.IsValid() {
		return StorageObject{}, errors.New("storage object: status is invalid")
	}
	record.CreatedAt, err = normalizeRequiredTime("created_at", record.CreatedAt)
	if err != nil {
		return StorageObject{}, err
	}
	record.UpdatedAt, err = normalizeRequiredTime("updated_at", record.UpdatedAt)
	if err != nil {
		return StorageObject{}, err
	}
	if record.UpdatedAt.Before(record.CreatedAt) {
		return StorageObject{}, errors.New("storage object: updated_at must not be before created_at")
	}
	if record.DeletedAt != nil {
		deletedAt := record.DeletedAt.UTC()
		if deletedAt.Before(record.CreatedAt) {
			return StorageObject{}, errors.New("storage object: deleted_at must not be before created_at")
		}
		if record.Status != StorageObjectStatusDeleted {
			return StorageObject{}, errors.New("storage object: deleted_at requires deleted status")
		}
		record.DeletedAt = &deletedAt
	}
	if record.Status == StorageObjectStatusDeleted && record.DeletedAt == nil {
		return StorageObject{}, errors.New("storage object: deleted status requires deleted_at")
	}
	return record, nil
}

func NormalizeCreateStorageObjectInput(input CreateStorageObjectInput) (CreateStorageObjectInput, error) {
	var err error
	input.ObjectID, err = normalizeRequired("object_id", input.ObjectID)
	if err != nil {
		return CreateStorageObjectInput{}, err
	}
	input.Owner, err = NormalizeStorageObjectOwner(input.Owner)
	if err != nil {
		return CreateStorageObjectInput{}, err
	}
	input.Identity, err = NormalizeStorageObjectIdentity(input.Identity)
	if err != nil {
		return CreateStorageObjectInput{}, err
	}
	input.Value, err = NormalizeStorageObjectValue(input.Value)
	if err != nil {
		return CreateStorageObjectInput{}, err
	}
	if input.InitialVersion == 0 {
		input.InitialVersion = InitialStorageObjectVersion
	}
	if !input.InitialVersion.IsValid() {
		return CreateStorageObjectInput{}, errors.New("storage object: initial_version must be positive")
	}
	if input.Status == "" {
		input.Status = StorageObjectStatusActive
	}
	if input.Status != StorageObjectStatusActive {
		return CreateStorageObjectInput{}, errors.New("storage object: created objects must start active")
	}
	input.RequestedBy, err = normalizeRequired("requested_by", input.RequestedBy)
	if err != nil {
		return CreateStorageObjectInput{}, err
	}
	return input, nil
}

func NormalizeGetStorageObjectInput(input GetStorageObjectInput) (GetStorageObjectInput, error) {
	var err error
	input.Owner, err = NormalizeStorageObjectOwner(input.Owner)
	if err != nil {
		return GetStorageObjectInput{}, err
	}
	input.Identity, err = NormalizeStorageObjectIdentity(input.Identity)
	if err != nil {
		return GetStorageObjectInput{}, err
	}
	return input, nil
}

func NormalizeListStorageObjectsInput(input ListStorageObjectsInput) (ListStorageObjectsInput, error) {
	var err error
	input.Owner, err = NormalizeStorageObjectOwner(input.Owner)
	if err != nil {
		return ListStorageObjectsInput{}, err
	}
	input.Collection, err = normalizeBoundedIdentifier("collection", input.Collection, MaxCollectionLength)
	if err != nil {
		return ListStorageObjectsInput{}, err
	}
	if input.Limit == 0 {
		input.Limit = DefaultListStorageObjectsLimit
	}
	if input.Limit < 0 || input.Limit > MaxListStorageObjectsLimit {
		return ListStorageObjectsInput{}, errors.New("storage object: list limit is invalid")
	}
	if strings.TrimSpace(input.AfterObjectKey) != "" {
		input.AfterObjectKey, err = normalizeBoundedIdentifier("after_object_key", input.AfterObjectKey, MaxObjectKeyLength)
		if err != nil {
			return ListStorageObjectsInput{}, err
		}
	}
	return input, nil
}

func NormalizeUpdateStorageObjectInput(input UpdateStorageObjectInput) (UpdateStorageObjectInput, error) {
	var err error
	input.Owner, err = NormalizeStorageObjectOwner(input.Owner)
	if err != nil {
		return UpdateStorageObjectInput{}, err
	}
	input.Identity, err = NormalizeStorageObjectIdentity(input.Identity)
	if err != nil {
		return UpdateStorageObjectInput{}, err
	}
	input.Value, err = NormalizeStorageObjectValue(input.Value)
	if err != nil {
		return UpdateStorageObjectInput{}, err
	}
	if input.ExpectedVersion != nil && !input.ExpectedVersion.IsValid() {
		return UpdateStorageObjectInput{}, errors.New("storage object: expected_version must be positive")
	}
	input.RequestedBy, err = normalizeRequired("requested_by", input.RequestedBy)
	if err != nil {
		return UpdateStorageObjectInput{}, err
	}
	return input, nil
}

func NormalizeDeleteStorageObjectInput(input DeleteStorageObjectInput) (DeleteStorageObjectInput, error) {
	var err error
	input.Owner, err = NormalizeStorageObjectOwner(input.Owner)
	if err != nil {
		return DeleteStorageObjectInput{}, err
	}
	input.Identity, err = NormalizeStorageObjectIdentity(input.Identity)
	if err != nil {
		return DeleteStorageObjectInput{}, err
	}
	if input.ExpectedVersion != nil && !input.ExpectedVersion.IsValid() {
		return DeleteStorageObjectInput{}, errors.New("storage object: expected_version must be positive")
	}
	input.RequestedBy, err = normalizeRequired("requested_by", input.RequestedBy)
	if err != nil {
		return DeleteStorageObjectInput{}, err
	}
	return input, nil
}

func NormalizeStorageObjectOwner(owner StorageObjectOwner) (StorageObjectOwner, error) {
	owner.Kind = OwnerKind(strings.TrimSpace(string(owner.Kind)))
	if !owner.Kind.IsValid() {
		return StorageObjectOwner{}, errors.New("storage object: owner_kind is invalid")
	}
	var err error
	owner.ID, err = normalizeRequired("owner_id", owner.ID)
	if err != nil {
		return StorageObjectOwner{}, err
	}
	return owner, nil
}

func NormalizeStorageObjectIdentity(identity StorageObjectIdentity) (StorageObjectIdentity, error) {
	var err error
	identity.Collection, err = normalizeBoundedIdentifier("collection", identity.Collection, MaxCollectionLength)
	if err != nil {
		return StorageObjectIdentity{}, err
	}
	identity.Key, err = normalizeBoundedIdentifier("object_key", identity.Key, MaxObjectKeyLength)
	if err != nil {
		return StorageObjectIdentity{}, err
	}
	return identity, nil
}

func NormalizeStorageObjectValue(value StorageObjectValue) (StorageObjectValue, error) {
	if len(value.JSON) == 0 {
		return StorageObjectValue{}, errors.New("storage object: value_json is required")
	}
	trimmed := bytes.TrimSpace(value.JSON)
	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &decoded); err != nil || decoded == nil {
		return StorageObjectValue{}, errors.New("storage object: value_json must be a JSON object")
	}
	value.JSON = cloneBytes(trimmed)
	return value, nil
}

func normalizeRequired(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("storage object: %s is required", name)
	}
	return value, nil
}

func normalizeBoundedIdentifier(name string, value string, maxLength int) (string, error) {
	value, err := normalizeRequired(name, value)
	if err != nil {
		return "", err
	}
	if len(value) > maxLength {
		return "", fmt.Errorf("storage object: %s is too long", name)
	}
	return value, nil
}

func normalizeRequiredTime(name string, value time.Time) (time.Time, error) {
	if value.IsZero() {
		return time.Time{}, fmt.Errorf("storage object: %s is required", name)
	}
	return value.UTC(), nil
}

func cloneBytes(source []byte) []byte {
	if len(source) == 0 {
		return nil
	}
	clone := make([]byte, len(source))
	copy(clone, source)
	return clone
}
