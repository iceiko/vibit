package storage

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

const (
	MaxStorageObjectValueJSONBytes = 16 * 1024
	defaultRequestedBy             = "storage_object_service"
)

type Operation string

const (
	OperationNewService             Operation = "NewService"
	OperationGetOwnStorageObject    Operation = "GetOwnStorageObject"
	OperationListOwnStorageObjects  Operation = "ListOwnStorageObjects"
	OperationPutOwnStorageObject    Operation = "PutOwnStorageObject"
	OperationDeleteOwnStorageObject Operation = "DeleteOwnStorageObject"
)

type FailureClass string

const (
	FailureClassInvalidRequest        FailureClass = "invalid_request"
	FailureClassForbidden             FailureClass = "forbidden"
	FailureClassNotFound              FailureClass = "not_found"
	FailureClassAlreadyExists         FailureClass = "already_exists"
	FailureClassVersionMismatch       FailureClass = "version_mismatch"
	FailureClassDependencyUnavailable FailureClass = "dependency_unavailable"
)

type PublicErrorCode string

const (
	PublicErrorStorageObjectInvalidRequest  PublicErrorCode = "STORAGE_OBJECT_INVALID_REQUEST"
	PublicErrorStorageObjectNotFound        PublicErrorCode = "STORAGE_OBJECT_NOT_FOUND"
	PublicErrorStorageObjectAlreadyExists   PublicErrorCode = "STORAGE_OBJECT_ALREADY_EXISTS"
	PublicErrorStorageObjectVersionMismatch PublicErrorCode = "STORAGE_OBJECT_VERSION_MISMATCH"
	PublicErrorStorageObjectUnavailable     PublicErrorCode = "STORAGE_OBJECT_UNAVAILABLE"
	PublicErrorStorageObjectForbidden       PublicErrorCode = "STORAGE_OBJECT_FORBIDDEN"
)

type ServiceError struct {
	Operation  Operation
	Class      FailureClass
	PublicCode PublicErrorCode
	Err        error
}

func (e *ServiceError) Error() string {
	if e == nil {
		return ""
	}
	if e.Operation == "" {
		return fmt.Sprintf("storage object service: %s", e.PublicCode)
	}
	return fmt.Sprintf("storage object service: %s: %s", e.Operation, e.PublicCode)
}

func (e *ServiceError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *ServiceError) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

type UnitOfWorkRunner interface {
	WithinUnitOfWork(context.Context, func(context.Context, tx.UnitOfWork) error) error
}

type ObjectIDGenerator interface {
	GenerateStorageObjectID(context.Context) (string, error)
}

type ServiceDependencies struct {
	UnitOfWorkRunner  UnitOfWorkRunner
	ObjectIDGenerator ObjectIDGenerator
}

type Service struct {
	unitOfWorkRunner  UnitOfWorkRunner
	objectIDGenerator ObjectIDGenerator
}

func NewService(dependencies ServiceDependencies) (Service, error) {
	if isNilInterface(dependencies.UnitOfWorkRunner) || isNilInterface(dependencies.ObjectIDGenerator) {
		return Service{}, serviceFailure(OperationNewService, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, errMissingDependency)
	}
	return Service{
		unitOfWorkRunner:  dependencies.UnitOfWorkRunner,
		objectIDGenerator: dependencies.ObjectIDGenerator,
	}, nil
}

type StorageObjectOperationStatus string

const (
	StorageObjectOperationStatusRejected StorageObjectOperationStatus = "rejected"
	StorageObjectOperationStatusFound    StorageObjectOperationStatus = "found"
	StorageObjectOperationStatusListed   StorageObjectOperationStatus = "listed"
	StorageObjectOperationStatusStored   StorageObjectOperationStatus = "stored"
	StorageObjectOperationStatusDeleted  StorageObjectOperationStatus = "deleted"
)

type GetOwnStorageObjectRequest struct {
	Identity    app.RequestIdentity
	Collection  string
	Key         string
	ClientOwner string
}

type ListOwnStorageObjectsRequest struct {
	Identity       app.RequestIdentity
	Collection     string
	Limit          int
	AfterObjectKey string
}

type PutOwnStorageObjectRequest struct {
	Identity        app.RequestIdentity
	Collection      string
	Key             string
	ValueJSON       []byte
	ExpectedVersion *storagemodule.StorageObjectVersion
	ClientOwner     string
}

type DeleteOwnStorageObjectRequest struct {
	Identity        app.RequestIdentity
	Collection      string
	Key             string
	ExpectedVersion *storagemodule.StorageObjectVersion
}

type StorageObjectResult struct {
	Status          StorageObjectOperationStatus
	PublicErrorCode PublicErrorCode
	FailureClass    FailureClass
	Object          storagemodule.StorageObject
	Objects         []storagemodule.StorageObject
	NextObjectKey   string
	Version         storagemodule.StorageObjectVersion
}

func (s Service) GetOwnStorageObject(ctx context.Context, request GetOwnStorageObjectRequest) (StorageObjectResult, error) {
	owner, err := ownerFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectForbidden, FailureClassForbidden),
			serviceFailure(OperationGetOwnStorageObject, FailureClassForbidden, PublicErrorStorageObjectForbidden, err)
	}
	input, err := storagemodule.NormalizeGetStorageObjectInput(storagemodule.GetStorageObjectInput{
		Owner: owner,
		Identity: storagemodule.StorageObjectIdentity{
			Collection: request.Collection,
			Key:        request.Key,
		},
	})
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationGetOwnStorageObject, FailureClassInvalidRequest, PublicErrorStorageObjectInvalidRequest, errInvalidRequest)
	}

	var committedResult StorageObjectResult
	var failedResult StorageObjectResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := storageRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationGetOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
		}
		record, err := repository.GetStorageObject(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationGetOwnStorageObject, err)
			return err
		}
		record, err = storagemodule.NormalizeStorageObjectRecord(record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationGetOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, errStorageObjectRecordInvalid)
		}
		committedResult = StorageObjectResult{
			Status:  StorageObjectOperationStatusFound,
			Object:  record,
			Version: record.Version,
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationGetOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
	}
	return committedResult, nil
}

func (s Service) ListOwnStorageObjects(ctx context.Context, request ListOwnStorageObjectsRequest) (StorageObjectResult, error) {
	owner, err := ownerFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectForbidden, FailureClassForbidden),
			serviceFailure(OperationListOwnStorageObjects, FailureClassForbidden, PublicErrorStorageObjectForbidden, err)
	}
	input, err := storagemodule.NormalizeListStorageObjectsInput(storagemodule.ListStorageObjectsInput{
		Owner:          owner,
		Collection:     request.Collection,
		Limit:          request.Limit,
		AfterObjectKey: strings.TrimSpace(request.AfterObjectKey),
	})
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationListOwnStorageObjects, FailureClassInvalidRequest, PublicErrorStorageObjectInvalidRequest, errInvalidRequest)
	}

	var committedResult StorageObjectResult
	var failedResult StorageObjectResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := storageRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationListOwnStorageObjects, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
		}
		result, err := repository.ListStorageObjects(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationListOwnStorageObjects, err)
			return err
		}
		objects := make([]storagemodule.StorageObject, 0, len(result.Objects))
		for _, object := range result.Objects {
			normalized, err := storagemodule.NormalizeStorageObjectRecord(object)
			if err != nil {
				failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
				return serviceFailure(OperationListOwnStorageObjects, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, errStorageObjectRecordInvalid)
			}
			objects = append(objects, normalized)
		}
		committedResult = StorageObjectResult{
			Status:        StorageObjectOperationStatusListed,
			Objects:       objects,
			NextObjectKey: strings.TrimSpace(result.NextObjectKey),
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationListOwnStorageObjects, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
	}
	return committedResult, nil
}

func (s Service) PutOwnStorageObject(ctx context.Context, request PutOwnStorageObjectRequest) (StorageObjectResult, error) {
	owner, err := ownerFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectForbidden, FailureClassForbidden),
			serviceFailure(OperationPutOwnStorageObject, FailureClassForbidden, PublicErrorStorageObjectForbidden, err)
	}
	if len(request.ValueJSON) > MaxStorageObjectValueJSONBytes {
		return rejectedResult(PublicErrorStorageObjectInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationPutOwnStorageObject, FailureClassInvalidRequest, PublicErrorStorageObjectInvalidRequest, errInvalidRequest)
	}
	updateInput, err := storagemodule.NormalizeUpdateStorageObjectInput(storagemodule.UpdateStorageObjectInput{
		Owner: owner,
		Identity: storagemodule.StorageObjectIdentity{
			Collection: request.Collection,
			Key:        request.Key,
		},
		Value: storagemodule.StorageObjectValue{
			JSON: request.ValueJSON,
		},
		ExpectedVersion: request.ExpectedVersion,
		RequestedBy:     defaultRequestedBy,
	})
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationPutOwnStorageObject, FailureClassInvalidRequest, PublicErrorStorageObjectInvalidRequest, errInvalidRequest)
	}

	if request.ExpectedVersion != nil {
		return s.updateOwnStorageObject(ctx, updateInput)
	}
	return s.createOrReplaceOwnStorageObject(ctx, updateInput)
}

func (s Service) DeleteOwnStorageObject(ctx context.Context, request DeleteOwnStorageObjectRequest) (StorageObjectResult, error) {
	owner, err := ownerFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectForbidden, FailureClassForbidden),
			serviceFailure(OperationDeleteOwnStorageObject, FailureClassForbidden, PublicErrorStorageObjectForbidden, err)
	}
	input, err := storagemodule.NormalizeDeleteStorageObjectInput(storagemodule.DeleteStorageObjectInput{
		Owner: owner,
		Identity: storagemodule.StorageObjectIdentity{
			Collection: request.Collection,
			Key:        request.Key,
		},
		ExpectedVersion: request.ExpectedVersion,
		RequestedBy:     defaultRequestedBy,
	})
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationDeleteOwnStorageObject, FailureClassInvalidRequest, PublicErrorStorageObjectInvalidRequest, errInvalidRequest)
	}

	var committedResult StorageObjectResult
	var failedResult StorageObjectResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := storageRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationDeleteOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
		}
		record, err := repository.DeleteStorageObject(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationDeleteOwnStorageObject, err)
			return err
		}
		record, err = storagemodule.NormalizeStorageObjectRecord(record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationDeleteOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, errStorageObjectRecordInvalid)
		}
		committedResult = StorageObjectResult{
			Status:  StorageObjectOperationStatusDeleted,
			Object:  record,
			Version: record.Version,
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationDeleteOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
	}
	return committedResult, nil
}

func (s Service) updateOwnStorageObject(ctx context.Context, input storagemodule.UpdateStorageObjectInput) (StorageObjectResult, error) {
	var committedResult StorageObjectResult
	var failedResult StorageObjectResult
	err := s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := storageRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationPutOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
		}
		record, err := repository.UpdateStorageObject(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationPutOwnStorageObject, err)
			return err
		}
		record, err = storagemodule.NormalizeStorageObjectRecord(record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationPutOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, errStorageObjectRecordInvalid)
		}
		committedResult = StorageObjectResult{
			Status:  StorageObjectOperationStatusStored,
			Object:  record,
			Version: record.Version,
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationPutOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
	}
	return committedResult, nil
}

func (s Service) createOrReplaceOwnStorageObject(ctx context.Context, input storagemodule.UpdateStorageObjectInput) (StorageObjectResult, error) {
	objectID, err := s.generatedStorageObjectID(ctx)
	if err != nil {
		return rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationPutOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
	}

	var committedResult StorageObjectResult
	var failedResult StorageObjectResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := storageRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationPutOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
		}
		record, err := repository.CreateStorageObject(runCtx, storagemodule.CreateStorageObjectInput{
			ObjectID:    objectID,
			Owner:       input.Owner,
			Identity:    input.Identity,
			Value:       input.Value,
			RequestedBy: defaultRequestedBy,
		})
		if err != nil {
			if repositoryConflictClass(err) != storagemodule.StorageObjectConflictAlreadyExists {
				failedResult, err = mapRepositoryFailure(OperationPutOwnStorageObject, err)
				return err
			}
			record, err = repository.UpdateStorageObject(runCtx, input)
			if err != nil {
				failedResult, err = mapRepositoryFailure(OperationPutOwnStorageObject, err)
				return err
			}
		}
		record, err = storagemodule.NormalizeStorageObjectRecord(record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationPutOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, errStorageObjectRecordInvalid)
		}
		committedResult = StorageObjectResult{
			Status:  StorageObjectOperationStatusStored,
			Object:  record,
			Version: record.Version,
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedResult(PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationPutOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable, err)
	}
	return committedResult, nil
}

func ownerFromValidatedIdentity(identity app.RequestIdentity) (storagemodule.StorageObjectOwner, error) {
	playerID := strings.TrimSpace(identity.PlayerID)
	actorID := strings.TrimSpace(identity.ActorID)
	if identity.Status != app.IdentityValidationValidated ||
		identity.ActorKind != app.ActorKindPlayer ||
		!identity.PlayerIDValidated ||
		playerID == "" ||
		actorID == "" ||
		playerID != actorID {
		return storagemodule.StorageObjectOwner{}, errForbiddenIdentity
	}
	return storagemodule.NormalizeStorageObjectOwner(storagemodule.StorageObjectOwner{
		Kind: storagemodule.OwnerKindPlayer,
		ID:   playerID,
	})
}

type storageUnitOfWork interface {
	NewStorageObjectRepository() (storagemodule.Repository, error)
}

func storageRepositoryFromUnitOfWork(unit tx.UnitOfWork) (storagemodule.Repository, error) {
	repositories, ok := unit.(storageUnitOfWork)
	if !ok {
		return nil, errMissingStorageUnitOfWork
	}
	repository, err := repositories.NewStorageObjectRepository()
	if err != nil {
		return nil, err
	}
	if isNilInterface(repository) {
		return nil, errMissingRepository
	}
	return repository, nil
}

func mapRepositoryFailure(operation Operation, err error) (StorageObjectResult, error) {
	publicCode, class := publicFailureForRepositoryError(err)
	return rejectedResult(publicCode, class), serviceFailure(operation, class, publicCode, nil)
}

func publicFailureForRepositoryError(err error) (PublicErrorCode, FailureClass) {
	switch repositoryConflictClass(err) {
	case storagemodule.StorageObjectConflictNotFound,
		storagemodule.StorageObjectConflictOwnerScopeMismatch,
		storagemodule.StorageObjectConflictDeletedObject:
		return PublicErrorStorageObjectNotFound, FailureClassNotFound
	case storagemodule.StorageObjectConflictAlreadyExists:
		return PublicErrorStorageObjectAlreadyExists, FailureClassAlreadyExists
	case storagemodule.StorageObjectConflictVersionMismatch:
		return PublicErrorStorageObjectVersionMismatch, FailureClassVersionMismatch
	case storagemodule.StorageObjectConflictInvalidExpectedVersion:
		return PublicErrorStorageObjectInvalidRequest, FailureClassInvalidRequest
	case storagemodule.StorageObjectConflictStorageUnavailable:
		return PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable
	default:
		if errors.Is(err, storagemodule.ErrStorageObjectInvalidInput) {
			return PublicErrorStorageObjectInvalidRequest, FailureClassInvalidRequest
		}
		if errors.Is(err, storagemodule.ErrStorageUnavailable) {
			return PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable
		}
		return PublicErrorStorageObjectUnavailable, FailureClassDependencyUnavailable
	}
}

func repositoryConflictClass(err error) storagemodule.StorageObjectConflictClass {
	var repositoryErr *storagemodule.StorageObjectRepositoryError
	if errors.As(err, &repositoryErr) {
		return repositoryErr.Conflict.Class
	}
	var conflict storagemodule.StorageObjectConflict
	if errors.As(err, &conflict) {
		return conflict.Class
	}
	return ""
}

func (s Service) generatedStorageObjectID(ctx context.Context) (string, error) {
	objectID, err := s.objectIDGenerator.GenerateStorageObjectID(ctx)
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(objectID)
	if trimmed == "" || trimmed != objectID {
		return "", errMalformedObjectID
	}
	return trimmed, nil
}

func rejectedResult(publicCode PublicErrorCode, class FailureClass) StorageObjectResult {
	return StorageObjectResult{
		Status:          StorageObjectOperationStatusRejected,
		PublicErrorCode: publicCode,
		FailureClass:    class,
	}
}

func serviceFailure(operation Operation, class FailureClass, publicCode PublicErrorCode, err error) error {
	return &ServiceError{
		Operation:  operation,
		Class:      class,
		PublicCode: publicCode,
		Err:        err,
	}
}

func isNilInterface(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

var (
	errMissingDependency          = errors.New("storage object service: dependency is required")
	errForbiddenIdentity          = errors.New("storage object service: validated player identity is required")
	errInvalidRequest             = errors.New("storage object service: invalid request")
	errMissingStorageUnitOfWork   = errors.New("storage object service: storage unit-of-work capability is required")
	errMissingRepository          = errors.New("storage object service: storage repository is required")
	errStorageObjectRecordInvalid = errors.New("storage object service: storage object record is invalid")
	errMalformedObjectID          = errors.New("storage object service: generated object id is malformed")
)
