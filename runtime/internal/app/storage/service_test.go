package storage

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

func TestNewServiceRequiresDependencies(t *testing.T) {
	_, err := NewService(ServiceDependencies{})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable)

	var runner *recordingUnitOfWorkRunner
	_, err = NewService(ServiceDependencies{
		UnitOfWorkRunner: runner,
		ObjectIDGenerator: staticObjectIDGenerator{
			value: "storage-object-1",
		},
	})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable)

	_, err = NewService(ServiceDependencies{
		UnitOfWorkRunner: &recordingUnitOfWorkRunner{},
	})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable)
}

func TestGetOwnStorageObjectRejectsMetadataOnlyIdentityBeforeUnitOfWork(t *testing.T) {
	repository := &fakeStorageRepository{}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})

	result, err := service.GetOwnStorageObject(context.Background(), GetOwnStorageObjectRequest{
		Identity: app.MetadataOnlyIdentityFromSession(app.Session{
			PlayerID:  "player-1",
			SessionID: "metadata-session",
		}),
		Collection: "progress",
		Key:        "tutorial",
	})

	assertServiceError(t, err, OperationGetOwnStorageObject, FailureClassForbidden, PublicErrorStorageObjectForbidden)
	if result.Status != StorageObjectOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorStorageObjectForbidden ||
		result.FailureClass != FailureClassForbidden {
		t.Fatalf("result = %#v, want forbidden rejection", result)
	}
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if repository.totalCalls() != 0 {
		t.Fatalf("repository calls = %d, want 0", repository.totalCalls())
	}
}

func TestGetOwnStorageObjectRejectsMismatchedActorIdentityBeforeUnitOfWork(t *testing.T) {
	repository := &fakeStorageRepository{}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})
	identity := validatedPlayerIdentity("player-1")
	identity.ActorID = "player-2"

	result, err := service.GetOwnStorageObject(context.Background(), GetOwnStorageObjectRequest{
		Identity:    identity,
		Collection:  "progress",
		Key:         "tutorial",
		ClientOwner: "player-2",
	})

	assertServiceError(t, err, OperationGetOwnStorageObject, FailureClassForbidden, PublicErrorStorageObjectForbidden)
	if result.Status != StorageObjectOperationStatusRejected {
		t.Fatalf("result = %#v, want rejected", result)
	}
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if repository.totalCalls() != 0 {
		t.Fatalf("repository calls = %d, want 0", repository.totalCalls())
	}
}

func TestGetOwnStorageObjectDerivesOwnerFromValidatedIdentity(t *testing.T) {
	events := []string{}
	repository := &fakeStorageRepository{
		events:    &events,
		getResult: activeStorageObject("player-1", "progress", "tutorial", []byte(`{"level":3}`), 4),
	}
	runner := &recordingUnitOfWorkRunner{
		unit:   &fakeStorageUnitOfWork{repository: repository, events: &events},
		events: &events,
	}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1", events: &events})

	result, err := service.GetOwnStorageObject(context.Background(), GetOwnStorageObjectRequest{
		Identity:    validatedPlayerIdentity("player-1"),
		Collection:  " progress ",
		Key:         " tutorial ",
		ClientOwner: "player-2",
	})
	if err != nil {
		t.Fatalf("GetOwnStorageObject() error = %v, want nil", err)
	}
	if result.Status != StorageObjectOperationStatusFound ||
		result.Object.Owner.ID != "player-1" ||
		result.Object.Identity.Collection != "progress" ||
		result.Object.Identity.Key != "tutorial" ||
		result.Version != 4 {
		t.Fatalf("result = %#v, want found player-owned object", result)
	}
	if repository.getCalls != 1 {
		t.Fatalf("GetStorageObject calls = %d, want 1", repository.getCalls)
	}
	if repository.lastGetInput.Owner != (storagemodule.StorageObjectOwner{Kind: storagemodule.OwnerKindPlayer, ID: "player-1"}) {
		t.Fatalf("owner = %#v, want validated player owner", repository.lastGetInput.Owner)
	}
	if repository.lastGetInput.Identity.Collection != "progress" || repository.lastGetInput.Identity.Key != "tutorial" {
		t.Fatalf("identity = %#v, want normalized collection/key", repository.lastGetInput.Identity)
	}
	assertEvents(t, events, []string{"begin", "new-storage-object-repository", "get-storage-object", "commit"})
}

func TestListOwnStorageObjectsUsesValidatedOwnerAndBoundsPagination(t *testing.T) {
	sourceValue := []byte(`{"slot":1}`)
	repository := &fakeStorageRepository{
		listResult: storagemodule.ListStorageObjectsResult{
			Objects: []storagemodule.StorageObject{
				activeStorageObject("player-1", "loadout", "primary", sourceValue, 2),
			},
			NextObjectKey: "secondary",
		},
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})

	result, err := service.ListOwnStorageObjects(context.Background(), ListOwnStorageObjectsRequest{
		Identity:       validatedPlayerIdentity("player-1"),
		Collection:     " loadout ",
		Limit:          0,
		AfterObjectKey: " ",
	})
	if err != nil {
		t.Fatalf("ListOwnStorageObjects() error = %v, want nil", err)
	}
	if result.Status != StorageObjectOperationStatusListed ||
		len(result.Objects) != 1 ||
		result.Objects[0].Owner.ID != "player-1" ||
		result.NextObjectKey != "secondary" {
		t.Fatalf("result = %#v, want listed owner-scoped objects", result)
	}
	if repository.lastListInput.Owner.ID != "player-1" ||
		repository.lastListInput.Collection != "loadout" ||
		repository.lastListInput.Limit != storagemodule.DefaultListStorageObjectsLimit ||
		repository.lastListInput.AfterObjectKey != "" {
		t.Fatalf("ListStorageObjects input = %#v, want normalized owner collection and default limit", repository.lastListInput)
	}
	sourceValue[0] = '['
	if string(result.Objects[0].Value.JSON) != `{"slot":1}` {
		t.Fatalf("result value = %q, want copied value", string(result.Objects[0].Value.JSON))
	}
}

func TestListOwnStorageObjectsRejectsInvalidLimitBeforeUnitOfWork(t *testing.T) {
	repository := &fakeStorageRepository{}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})

	result, err := service.ListOwnStorageObjects(context.Background(), ListOwnStorageObjectsRequest{
		Identity:   validatedPlayerIdentity("player-1"),
		Collection: "loadout",
		Limit:      storagemodule.MaxListStorageObjectsLimit + 1,
	})

	assertServiceError(t, err, OperationListOwnStorageObjects, FailureClassInvalidRequest, PublicErrorStorageObjectInvalidRequest)
	if result.Status != StorageObjectOperationStatusRejected {
		t.Fatalf("result = %#v, want rejected", result)
	}
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if repository.totalCalls() != 0 {
		t.Fatalf("repository calls = %d, want 0", repository.totalCalls())
	}
}

func TestPutOwnStorageObjectCreatesObjectWithServerDerivedOwner(t *testing.T) {
	events := []string{}
	repository := &fakeStorageRepository{
		events:       &events,
		createResult: activeStorageObject("player-1", "progress", "tutorial", []byte(`{"level":3}`), 1),
	}
	generator := staticObjectIDGenerator{value: "storage-object-1", events: &events}
	runner := &recordingUnitOfWorkRunner{
		unit:   &fakeStorageUnitOfWork{repository: repository, events: &events},
		events: &events,
	}
	service := mustNewService(t, runner, generator)
	value := []byte(`{"level":3}`)

	result, err := service.PutOwnStorageObject(context.Background(), PutOwnStorageObjectRequest{
		Identity:    validatedPlayerIdentity("player-1"),
		Collection:  " progress ",
		Key:         " tutorial ",
		ValueJSON:   value,
		ClientOwner: "player-2",
	})
	if err != nil {
		t.Fatalf("PutOwnStorageObject() error = %v, want nil", err)
	}
	if result.Status != StorageObjectOperationStatusStored ||
		result.Object.Owner.ID != "player-1" ||
		result.Version != storagemodule.InitialStorageObjectVersion {
		t.Fatalf("result = %#v, want stored player object", result)
	}
	if repository.createCalls != 1 || repository.updateCalls != 0 {
		t.Fatalf("create/update calls = %d/%d, want 1/0", repository.createCalls, repository.updateCalls)
	}
	created := repository.lastCreateInput
	if created.ObjectID != "storage-object-1" ||
		created.Owner.ID != "player-1" ||
		created.Identity.Collection != "progress" ||
		created.Identity.Key != "tutorial" ||
		created.RequestedBy != defaultRequestedBy {
		t.Fatalf("CreateStorageObject input = %#v, want server-derived create input", created)
	}
	value[0] = '['
	if string(created.Value.JSON) != `{"level":3}` {
		t.Fatalf("CreateStorageObject value = %q, want copied request value", string(created.Value.JSON))
	}
	assertEvents(t, events, []string{"generate-storage-object-id", "begin", "new-storage-object-repository", "create-storage-object", "commit"})
}

func TestPutOwnStorageObjectReplacesExistingObjectWithoutExpectedVersion(t *testing.T) {
	repository := &fakeStorageRepository{
		createErr:    storageRepositoryError(storagemodule.StorageObjectConflictAlreadyExists, `duplicate value_json {"secret":true}`),
		updateResult: activeStorageObject("player-1", "progress", "tutorial", []byte(`{"level":4}`), 2),
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-2"})

	result, err := service.PutOwnStorageObject(context.Background(), PutOwnStorageObjectRequest{
		Identity:   validatedPlayerIdentity("player-1"),
		Collection: "progress",
		Key:        "tutorial",
		ValueJSON:  []byte(`{"level":4}`),
	})
	if err != nil {
		t.Fatalf("PutOwnStorageObject() error = %v, want nil", err)
	}
	if result.Status != StorageObjectOperationStatusStored || result.Version != 2 {
		t.Fatalf("result = %#v, want replacement result", result)
	}
	if repository.createCalls != 1 || repository.updateCalls != 1 {
		t.Fatalf("create/update calls = %d/%d, want 1/1", repository.createCalls, repository.updateCalls)
	}
	if repository.lastUpdateInput.ExpectedVersion != nil {
		t.Fatalf("ExpectedVersion = %#v, want nil replacement without expected version", repository.lastUpdateInput.ExpectedVersion)
	}
}

func TestPutOwnStorageObjectRejectsOversizedValueBeforeUnitOfWork(t *testing.T) {
	repository := &fakeStorageRepository{}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})
	value := []byte(`{"payload":"` + strings.Repeat("x", MaxStorageObjectValueJSONBytes) + `"}`)

	result, err := service.PutOwnStorageObject(context.Background(), PutOwnStorageObjectRequest{
		Identity:   validatedPlayerIdentity("player-1"),
		Collection: "progress",
		Key:        "tutorial",
		ValueJSON:  value,
	})

	assertServiceError(t, err, OperationPutOwnStorageObject, FailureClassInvalidRequest, PublicErrorStorageObjectInvalidRequest)
	if result.Status != StorageObjectOperationStatusRejected {
		t.Fatalf("result = %#v, want rejected", result)
	}
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if repository.totalCalls() != 0 {
		t.Fatalf("repository calls = %d, want 0", repository.totalCalls())
	}
}

func TestPutOwnStorageObjectMapsVersionMismatchAndRedactsDetails(t *testing.T) {
	rawDetail := `UPDATE storage_objects SET value_json={"secret":true}`
	expected := storagemodule.StorageObjectVersion(7)
	repository := &fakeStorageRepository{
		updateErr: storageRepositoryError(storagemodule.StorageObjectConflictVersionMismatch, rawDetail),
	}
	generator := &recordingObjectIDGenerator{value: "storage-object-1"}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, generator)

	result, err := service.PutOwnStorageObject(context.Background(), PutOwnStorageObjectRequest{
		Identity:        validatedPlayerIdentity("player-1"),
		Collection:      "progress",
		Key:             "tutorial",
		ValueJSON:       []byte(`{"level":4}`),
		ExpectedVersion: &expected,
	})

	assertServiceError(t, err, OperationPutOwnStorageObject, FailureClassVersionMismatch, PublicErrorStorageObjectVersionMismatch)
	assertNoLeak(t, err, rawDetail)
	if result.Status != StorageObjectOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorStorageObjectVersionMismatch {
		t.Fatalf("result = %#v, want version mismatch rejection", result)
	}
	if generator.calls != 0 {
		t.Fatalf("object id generator calls = %d, want 0 for expected-version update", generator.calls)
	}
	if repository.createCalls != 0 || repository.updateCalls != 1 {
		t.Fatalf("create/update calls = %d/%d, want 0/1", repository.createCalls, repository.updateCalls)
	}
	if repository.lastUpdateInput.ExpectedVersion == nil || *repository.lastUpdateInput.ExpectedVersion != expected {
		t.Fatalf("ExpectedVersion = %#v, want %d", repository.lastUpdateInput.ExpectedVersion, expected)
	}
}

func TestDeleteOwnStorageObjectChecksExpectedVersion(t *testing.T) {
	expected := storagemodule.StorageObjectVersion(4)
	deletedAt := fixedStorageObjectTime().Add(time.Minute)
	repository := &fakeStorageRepository{
		deleteResult: deletedStorageObject("player-1", "progress", "tutorial", []byte(`{"level":4}`), 5, deletedAt),
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})

	result, err := service.DeleteOwnStorageObject(context.Background(), DeleteOwnStorageObjectRequest{
		Identity:        validatedPlayerIdentity("player-1"),
		Collection:      " progress ",
		Key:             " tutorial ",
		ExpectedVersion: &expected,
	})
	if err != nil {
		t.Fatalf("DeleteOwnStorageObject() error = %v, want nil", err)
	}
	if result.Status != StorageObjectOperationStatusDeleted ||
		result.Object.Status != storagemodule.StorageObjectStatusDeleted ||
		result.Version != 5 {
		t.Fatalf("result = %#v, want deleted object", result)
	}
	if repository.deleteCalls != 1 {
		t.Fatalf("DeleteStorageObject calls = %d, want 1", repository.deleteCalls)
	}
	if repository.lastDeleteInput.Owner.ID != "player-1" ||
		repository.lastDeleteInput.Identity.Collection != "progress" ||
		repository.lastDeleteInput.Identity.Key != "tutorial" ||
		repository.lastDeleteInput.ExpectedVersion == nil ||
		*repository.lastDeleteInput.ExpectedVersion != expected ||
		repository.lastDeleteInput.RequestedBy != defaultRequestedBy {
		t.Fatalf("DeleteStorageObject input = %#v, want normalized expected-version delete", repository.lastDeleteInput)
	}
}

func TestGetOwnStorageObjectMapsNotFoundWithoutExistenceLeak(t *testing.T) {
	rawDetail := `owner_id=other-player value_json={"secret":true}`
	repository := &fakeStorageRepository{
		getErr: storageRepositoryError(storagemodule.StorageObjectConflictOwnerScopeMismatch, rawDetail),
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeStorageUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})

	result, err := service.GetOwnStorageObject(context.Background(), GetOwnStorageObjectRequest{
		Identity:   validatedPlayerIdentity("player-1"),
		Collection: "progress",
		Key:        "tutorial",
	})

	assertServiceError(t, err, OperationGetOwnStorageObject, FailureClassNotFound, PublicErrorStorageObjectNotFound)
	assertNoLeak(t, err, rawDetail)
	if result.Status != StorageObjectOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorStorageObjectNotFound {
		t.Fatalf("result = %#v, want not-found rejection", result)
	}
}

func TestGetOwnStorageObjectRequiresUnitOfWorkRepositoryCapability(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{unit: tx.NoopUnitOfWork{}}
	service := mustNewService(t, runner, staticObjectIDGenerator{value: "storage-object-1"})

	result, err := service.GetOwnStorageObject(context.Background(), GetOwnStorageObjectRequest{
		Identity:   validatedPlayerIdentity("player-1"),
		Collection: "progress",
		Key:        "tutorial",
	})

	assertServiceError(t, err, OperationGetOwnStorageObject, FailureClassDependencyUnavailable, PublicErrorStorageObjectUnavailable)
	if result.Status != StorageObjectOperationStatusRejected {
		t.Fatalf("result = %#v, want rejected", result)
	}
}

func mustNewService(t *testing.T, runner UnitOfWorkRunner, generator ObjectIDGenerator) Service {
	t.Helper()
	service, err := NewService(ServiceDependencies{
		UnitOfWorkRunner:  runner,
		ObjectIDGenerator: generator,
	})
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}
	return service
}

func validatedPlayerIdentity(playerID string) app.RequestIdentity {
	return app.ValidatedPlayerIdentity(playerID, app.Session{
		ConnectionID:    "connection-1",
		SessionID:       "session-1",
		PlayerID:        playerID,
		ConnectionEpoch: 7,
	})
}

func activeStorageObject(ownerID, collection, key string, value []byte, version storagemodule.StorageObjectVersion) storagemodule.StorageObject {
	if version == 0 {
		version = storagemodule.InitialStorageObjectVersion
	}
	createdAt := fixedStorageObjectTime()
	return storagemodule.StorageObject{
		ObjectID: "storage-object-1",
		Owner: storagemodule.StorageObjectOwner{
			Kind: storagemodule.OwnerKindPlayer,
			ID:   ownerID,
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
		UpdatedAt: createdAt.Add(time.Second),
	}
}

func deletedStorageObject(ownerID, collection, key string, value []byte, version storagemodule.StorageObjectVersion, deletedAt time.Time) storagemodule.StorageObject {
	record := activeStorageObject(ownerID, collection, key, value, version)
	record.Status = storagemodule.StorageObjectStatusDeleted
	record.UpdatedAt = deletedAt.UTC()
	record.DeletedAt = &record.UpdatedAt
	return record
}

func fixedStorageObjectTime() time.Time {
	return time.Date(2026, 5, 22, 12, 0, 0, 0, time.UTC)
}

func storageRepositoryError(class storagemodule.StorageObjectConflictClass, rawDetail string) error {
	kind := storagemodule.ErrStorageObjectConflict
	if class == storagemodule.StorageObjectConflictStorageUnavailable {
		kind = storagemodule.ErrStorageUnavailable
	}
	return &storagemodule.StorageObjectRepositoryError{
		Kind: kind,
		Conflict: storagemodule.StorageObjectConflict{
			Class:          class,
			Retryable:      class == storagemodule.StorageObjectConflictStorageUnavailable,
			RedactedReason: string(class),
		},
		Operation:      "test",
		RedactedReason: string(class),
		Err:            errors.New(rawDetail),
	}
}

func assertServiceError(t *testing.T, err error, operation Operation, class FailureClass, publicCode PublicErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want service error %s/%s", class, publicCode)
	}
	var serviceErr *ServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("error = %T %v, want *ServiceError", err, err)
	}
	if serviceErr.Operation != operation ||
		serviceErr.Class != class ||
		serviceErr.PublicCode != publicCode {
		t.Fatalf("service error = %#v, want operation=%s class=%s public=%s", serviceErr, operation, class, publicCode)
	}
}

func assertNoLeak(t *testing.T, err error, forbidden string) {
	t.Helper()
	if err == nil {
		return
	}
	if forbidden != "" && strings.Contains(err.Error(), forbidden) {
		t.Fatalf("error %q leaks forbidden detail %q", err.Error(), forbidden)
	}
	if strings.Contains(err.Error(), "value_json") ||
		strings.Contains(err.Error(), "secret") ||
		strings.Contains(err.Error(), "other-player") {
		t.Fatalf("error %q leaks storage detail", err.Error())
	}
}

func assertEvents(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("events = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("events = %#v, want %#v", got, want)
		}
	}
}

type recordingUnitOfWorkRunner struct {
	unit      tx.UnitOfWork
	events    *[]string
	calls     int
	commitErr error
}

func (r *recordingUnitOfWorkRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	r.calls += 1
	if ctx == nil {
		ctx = context.Background()
	}
	unit := r.unit
	if unit == nil {
		unit = tx.NoopUnitOfWork{}
	}
	appendEvent(r.events, "begin")
	if err := fn(ctx, unit); err != nil {
		appendEvent(r.events, "rollback")
		return err
	}
	if r.commitErr != nil {
		appendEvent(r.events, "rollback")
		return r.commitErr
	}
	appendEvent(r.events, "commit")
	return nil
}

type fakeStorageUnitOfWork struct {
	ctx             context.Context
	repository      storagemodule.Repository
	repositoryErr   error
	events          *[]string
	repositoryCalls int
}

func (u *fakeStorageUnitOfWork) Context() context.Context {
	if u.ctx == nil {
		return context.Background()
	}
	return u.ctx
}

func (u *fakeStorageUnitOfWork) NewStorageObjectRepository() (storagemodule.Repository, error) {
	u.repositoryCalls += 1
	appendEvent(u.events, "new-storage-object-repository")
	if u.repositoryErr != nil {
		return nil, u.repositoryErr
	}
	return u.repository, nil
}

type fakeStorageRepository struct {
	events *[]string

	createResult storagemodule.StorageObject
	createErr    error
	getResult    storagemodule.StorageObject
	getErr       error
	listResult   storagemodule.ListStorageObjectsResult
	listErr      error
	updateResult storagemodule.StorageObject
	updateErr    error
	deleteResult storagemodule.StorageObject
	deleteErr    error

	createCalls int
	getCalls    int
	listCalls   int
	updateCalls int
	deleteCalls int

	lastCreateInput storagemodule.CreateStorageObjectInput
	lastGetInput    storagemodule.GetStorageObjectInput
	lastListInput   storagemodule.ListStorageObjectsInput
	lastUpdateInput storagemodule.UpdateStorageObjectInput
	lastDeleteInput storagemodule.DeleteStorageObjectInput
}

func (r *fakeStorageRepository) CreateStorageObject(_ context.Context, input storagemodule.CreateStorageObjectInput) (storagemodule.StorageObject, error) {
	r.createCalls += 1
	r.lastCreateInput = input
	appendEvent(r.events, "create-storage-object")
	if r.createErr != nil {
		return storagemodule.StorageObject{}, r.createErr
	}
	return r.createResult, nil
}

func (r *fakeStorageRepository) GetStorageObject(_ context.Context, input storagemodule.GetStorageObjectInput) (storagemodule.StorageObject, error) {
	r.getCalls += 1
	r.lastGetInput = input
	appendEvent(r.events, "get-storage-object")
	if r.getErr != nil {
		return storagemodule.StorageObject{}, r.getErr
	}
	return r.getResult, nil
}

func (r *fakeStorageRepository) ListStorageObjects(_ context.Context, input storagemodule.ListStorageObjectsInput) (storagemodule.ListStorageObjectsResult, error) {
	r.listCalls += 1
	r.lastListInput = input
	appendEvent(r.events, "list-storage-objects")
	if r.listErr != nil {
		return storagemodule.ListStorageObjectsResult{}, r.listErr
	}
	return r.listResult, nil
}

func (r *fakeStorageRepository) UpdateStorageObject(_ context.Context, input storagemodule.UpdateStorageObjectInput) (storagemodule.StorageObject, error) {
	r.updateCalls += 1
	r.lastUpdateInput = input
	appendEvent(r.events, "update-storage-object")
	if r.updateErr != nil {
		return storagemodule.StorageObject{}, r.updateErr
	}
	return r.updateResult, nil
}

func (r *fakeStorageRepository) DeleteStorageObject(_ context.Context, input storagemodule.DeleteStorageObjectInput) (storagemodule.StorageObject, error) {
	r.deleteCalls += 1
	r.lastDeleteInput = input
	appendEvent(r.events, "delete-storage-object")
	if r.deleteErr != nil {
		return storagemodule.StorageObject{}, r.deleteErr
	}
	return r.deleteResult, nil
}

func (r *fakeStorageRepository) totalCalls() int {
	return r.createCalls + r.getCalls + r.listCalls + r.updateCalls + r.deleteCalls
}

type staticObjectIDGenerator struct {
	value  string
	events *[]string
}

func (g staticObjectIDGenerator) GenerateStorageObjectID(context.Context) (string, error) {
	appendEvent(g.events, "generate-storage-object-id")
	return g.value, nil
}

type recordingObjectIDGenerator struct {
	value string
	err   error
	calls int
}

func (g *recordingObjectIDGenerator) GenerateStorageObjectID(context.Context) (string, error) {
	g.calls += 1
	if g.err != nil {
		return "", g.err
	}
	return g.value, nil
}

func appendEvent(events *[]string, event string) {
	if events == nil {
		return
	}
	*events = append(*events, event)
}
