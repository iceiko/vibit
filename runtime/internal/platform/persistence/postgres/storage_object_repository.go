package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/modules/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type StorageObjectRepository struct {
	executor Executor
}

var _ storage.Repository = (*StorageObjectRepository)(nil)

func NewStorageObjectRepositoryForUnitOfWork(executor Executor) *StorageObjectRepository {
	return &StorageObjectRepository{executor: executor}
}

func (r *StorageObjectRepository) CreateStorageObject(ctx context.Context, input storage.CreateStorageObjectInput) (storage.StorageObject, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return storage.StorageObject{}, err
	}

	normalized, err := storage.NormalizeCreateStorageObjectInput(input)
	if err != nil {
		return storage.StorageObject{}, storageObjectRepositoryError("create", storage.ErrStorageObjectInvalidInput, storage.StorageObjectConflictInvalidExpectedVersion, false)
	}

	record, err := scanStorageObjectRow(executor.QueryRow(
		ctx,
		insertStorageObjectSQL,
		normalized.ObjectID,
		string(normalized.Owner.Kind),
		normalized.Owner.ID,
		normalized.Identity.Collection,
		normalized.Identity.Key,
		normalized.Value.JSON,
		int64(normalized.InitialVersion),
	))
	if err != nil {
		return storage.StorageObject{}, mapStorageObjectPostgresError("create", err, nil)
	}
	return record, nil
}

func (r *StorageObjectRepository) GetStorageObject(ctx context.Context, input storage.GetStorageObjectInput) (storage.StorageObject, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return storage.StorageObject{}, err
	}

	normalized, err := storage.NormalizeGetStorageObjectInput(input)
	if err != nil {
		return storage.StorageObject{}, storageObjectRepositoryError("get", storage.ErrStorageObjectInvalidInput, storage.StorageObjectConflictOwnerScopeMismatch, false)
	}

	record, err := scanStorageObjectRow(executor.QueryRow(
		ctx,
		getStorageObjectSQL,
		string(normalized.Owner.Kind),
		normalized.Owner.ID,
		normalized.Identity.Collection,
		normalized.Identity.Key,
	))
	if err != nil {
		return storage.StorageObject{}, mapStorageObjectPostgresError("get", err, nil)
	}
	return record, nil
}

func (r *StorageObjectRepository) ListStorageObjects(ctx context.Context, input storage.ListStorageObjectsInput) (storage.ListStorageObjectsResult, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return storage.ListStorageObjectsResult{}, err
	}

	normalized, err := storage.NormalizeListStorageObjectsInput(input)
	if err != nil {
		return storage.ListStorageObjectsResult{}, storageObjectRepositoryError("list", storage.ErrStorageObjectInvalidInput, storage.StorageObjectConflictOwnerScopeMismatch, false)
	}

	rows, err := executor.Query(
		ctx,
		listStorageObjectsSQL,
		string(normalized.Owner.Kind),
		normalized.Owner.ID,
		normalized.Collection,
		normalized.AfterObjectKey,
		int32(normalized.Limit+1),
	)
	if err != nil {
		return storage.ListStorageObjectsResult{}, mapStorageObjectPostgresError("list", err, nil)
	}
	defer rows.Close()

	objects := make([]storage.StorageObject, 0, normalized.Limit)
	var nextObjectKey string
	for rows.Next() {
		record, err := scanStorageObjectScanner(rows)
		if err != nil {
			return storage.ListStorageObjectsResult{}, mapStorageObjectPostgresError("list", err, nil)
		}
		if len(objects) >= normalized.Limit {
			nextObjectKey = record.Identity.Key
			continue
		}
		objects = append(objects, record)
	}
	if err := rows.Err(); err != nil {
		return storage.ListStorageObjectsResult{}, mapStorageObjectPostgresError("list", err, nil)
	}
	return storage.ListStorageObjectsResult{Objects: objects, NextObjectKey: nextObjectKey}, nil
}

func (r *StorageObjectRepository) UpdateStorageObject(ctx context.Context, input storage.UpdateStorageObjectInput) (storage.StorageObject, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return storage.StorageObject{}, err
	}

	normalized, err := storage.NormalizeUpdateStorageObjectInput(input)
	if err != nil {
		return storage.StorageObject{}, storageObjectRepositoryError("update", storage.ErrStorageObjectInvalidInput, storage.StorageObjectConflictInvalidExpectedVersion, false)
	}

	if normalized.ExpectedVersion != nil {
		record, err := scanStorageObjectRow(executor.QueryRow(
			ctx,
			updateStorageObjectWithExpectedVersionSQL,
			string(normalized.Owner.Kind),
			normalized.Owner.ID,
			normalized.Identity.Collection,
			normalized.Identity.Key,
			normalized.Value.JSON,
			int64(*normalized.ExpectedVersion),
		))
		if err != nil {
			return storage.StorageObject{}, mapStorageObjectPostgresError("update", err, normalized.ExpectedVersion)
		}
		return record, nil
	}

	record, err := scanStorageObjectRow(executor.QueryRow(
		ctx,
		updateStorageObjectSQL,
		string(normalized.Owner.Kind),
		normalized.Owner.ID,
		normalized.Identity.Collection,
		normalized.Identity.Key,
		normalized.Value.JSON,
	))
	if err != nil {
		return storage.StorageObject{}, mapStorageObjectPostgresError("update", err, nil)
	}
	return record, nil
}

func (r *StorageObjectRepository) DeleteStorageObject(ctx context.Context, input storage.DeleteStorageObjectInput) (storage.StorageObject, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return storage.StorageObject{}, err
	}

	normalized, err := storage.NormalizeDeleteStorageObjectInput(input)
	if err != nil {
		return storage.StorageObject{}, storageObjectRepositoryError("delete", storage.ErrStorageObjectInvalidInput, storage.StorageObjectConflictInvalidExpectedVersion, false)
	}

	if normalized.ExpectedVersion != nil {
		record, err := scanStorageObjectRow(executor.QueryRow(
			ctx,
			deleteStorageObjectWithExpectedVersionSQL,
			string(normalized.Owner.Kind),
			normalized.Owner.ID,
			normalized.Identity.Collection,
			normalized.Identity.Key,
			int64(*normalized.ExpectedVersion),
		))
		if err != nil {
			return storage.StorageObject{}, mapStorageObjectPostgresError("delete", err, normalized.ExpectedVersion)
		}
		return record, nil
	}

	record, err := scanStorageObjectRow(executor.QueryRow(
		ctx,
		deleteStorageObjectSQL,
		string(normalized.Owner.Kind),
		normalized.Owner.ID,
		normalized.Identity.Collection,
		normalized.Identity.Key,
	))
	if err != nil {
		return storage.StorageObject{}, mapStorageObjectPostgresError("delete", err, nil)
	}
	return record, nil
}

func (r *StorageObjectRepository) requireExecutor() (Executor, error) {
	if r == nil || r.executor == nil {
		return nil, errors.New("postgres storage object: unit-of-work executor is required")
	}
	return r.executor, nil
}

func scanStorageObjectRow(row pgx.Row) (storage.StorageObject, error) {
	return scanStorageObjectScanner(row)
}

func scanStorageObjectScanner(row scanner) (storage.StorageObject, error) {
	var record storage.StorageObject
	var ownerKind string
	var value []byte
	var version int64
	var deletedAt pgtype.Timestamptz

	if err := row.Scan(
		&record.ObjectID,
		&ownerKind,
		&record.Owner.ID,
		&record.Identity.Collection,
		&record.Identity.Key,
		&value,
		&version,
		&record.CreatedAt,
		&record.UpdatedAt,
		&deletedAt,
	); err != nil {
		return storage.StorageObject{}, err
	}

	record.Owner.Kind = storage.OwnerKind(strings.TrimSpace(ownerKind))
	record.Value = storage.StorageObjectValue{JSON: value}
	record.Version = storage.StorageObjectVersion(version)
	record.DeletedAt = nullableTimestamptzUTC(deletedAt)
	if record.DeletedAt != nil {
		record.Status = storage.StorageObjectStatusDeleted
	} else {
		record.Status = storage.StorageObjectStatusActive
	}

	normalized, err := storage.NormalizeStorageObjectRecord(record)
	if err != nil {
		return storage.StorageObject{}, fmt.Errorf("%w: row shape", storage.ErrStorageUnavailable)
	}
	return normalized, nil
}

func mapStorageObjectPostgresError(operation string, err error, expectedVersion *storage.StorageObjectVersion) error {
	if errors.Is(err, pgx.ErrNoRows) {
		if expectedVersion != nil {
			return storageObjectRepositoryError(operation, storage.ErrStorageObjectConflict, storage.StorageObjectConflictVersionMismatch, true)
		}
		return storageObjectRepositoryError(operation, storage.ErrStorageObjectConflict, storage.StorageObjectConflictNotFound, false)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return storageObjectRepositoryError(operation, storage.ErrStorageObjectConflict, storage.StorageObjectConflictAlreadyExists, false)
		case "23502", "23503", "23514":
			return storageObjectRepositoryError(operation, storage.ErrStorageUnavailable, storage.StorageObjectConflictStorageUnavailable, true)
		default:
			return storageObjectRepositoryError(operation, storage.ErrStorageUnavailable, storage.StorageObjectConflictStorageUnavailable, true)
		}
	}

	if errors.Is(err, storage.ErrStorageUnavailable) {
		return storageObjectRepositoryError(operation, storage.ErrStorageUnavailable, storage.StorageObjectConflictStorageUnavailable, true)
	}
	return storageObjectRepositoryError(operation, storage.ErrStorageUnavailable, storage.StorageObjectConflictStorageUnavailable, true)
}

func storageObjectRepositoryError(operation string, kind error, class storage.StorageObjectConflictClass, retryable bool) error {
	reason := string(class)
	if reason == "" && kind != nil {
		reason = kind.Error()
	}
	return &storage.StorageObjectRepositoryError{
		Kind: kind,
		Conflict: storage.StorageObjectConflict{
			Class:          class,
			Retryable:      retryable,
			RedactedReason: reason,
		},
		Operation:      operation,
		RedactedReason: reason,
	}
}

const storageObjectColumns = `
object_id,
owner_kind,
owner_id,
collection,
object_key,
value_json,
version,
created_at,
updated_at,
deleted_at`

const insertStorageObjectSQL = `
INSERT INTO storage_objects (
    object_id,
    owner_kind,
    owner_id,
    collection,
    object_key,
    value_json,
    version,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING ` + storageObjectColumns

const getStorageObjectSQL = `
SELECT ` + storageObjectColumns + `
FROM storage_objects
WHERE owner_kind = $1
  AND owner_id = $2
  AND collection = $3
  AND object_key = $4
  AND deleted_at IS NULL`

const listStorageObjectsSQL = `
SELECT ` + storageObjectColumns + `
FROM storage_objects
WHERE owner_kind = $1
  AND owner_id = $2
  AND collection = $3
  AND object_key > $4
  AND deleted_at IS NULL
ORDER BY object_key
LIMIT $5`

const updateStorageObjectSQL = `
UPDATE storage_objects
SET value_json = $5,
    version = version + 1,
    updated_at = now()
WHERE owner_kind = $1
  AND owner_id = $2
  AND collection = $3
  AND object_key = $4
  AND deleted_at IS NULL
RETURNING ` + storageObjectColumns

const updateStorageObjectWithExpectedVersionSQL = `
UPDATE storage_objects
SET value_json = $5,
    version = version + 1,
    updated_at = now()
WHERE owner_kind = $1
  AND owner_id = $2
  AND collection = $3
  AND object_key = $4
  AND version = $6
  AND deleted_at IS NULL
RETURNING ` + storageObjectColumns

const deleteStorageObjectSQL = `
UPDATE storage_objects
SET deleted_at = now(),
    version = version + 1,
    updated_at = now()
WHERE owner_kind = $1
  AND owner_id = $2
  AND collection = $3
  AND object_key = $4
  AND deleted_at IS NULL
RETURNING ` + storageObjectColumns

const deleteStorageObjectWithExpectedVersionSQL = `
UPDATE storage_objects
SET deleted_at = now(),
    version = version + 1,
    updated_at = now()
WHERE owner_kind = $1
  AND owner_id = $2
  AND collection = $3
  AND object_key = $4
  AND version = $5
  AND deleted_at IS NULL
RETURNING ` + storageObjectColumns
