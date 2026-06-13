package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/modules/currency"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type CurrencyWalletRepository struct {
	executor Executor
}

var _ currency.Repository = (*CurrencyWalletRepository)(nil)

func NewCurrencyWalletRepositoryForUnitOfWork(executor Executor) *CurrencyWalletRepository {
	return &CurrencyWalletRepository{executor: executor}
}

func (r *CurrencyWalletRepository) CreateCurrencyWallet(ctx context.Context, input currency.CreateCurrencyWalletInput) (currency.CurrencyWallet, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return currency.CurrencyWallet{}, err
	}

	normalized, err := currency.NormalizeCreateCurrencyWalletInput(input)
	if err != nil {
		return currency.CurrencyWallet{}, currencyWalletInvalidInputError("create")
	}

	record, err := scanCurrencyWalletRow(executor.QueryRow(
		ctx,
		insertCurrencyWalletSQL,
		string(normalized.WalletID),
		string(normalized.Owner.Kind),
		normalized.Owner.ID,
		string(normalized.InitialState),
		int64(normalized.InitialVersion),
	))
	if err != nil {
		return currency.CurrencyWallet{}, mapCurrencyWalletPostgresError("create", err, nil)
	}
	return record, nil
}

func (r *CurrencyWalletRepository) GetCurrencyWallet(ctx context.Context, input currency.GetCurrencyWalletInput) (currency.CurrencyWallet, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return currency.CurrencyWallet{}, err
	}

	normalized, err := currency.NormalizeGetCurrencyWalletInput(input)
	if err != nil {
		return currency.CurrencyWallet{}, currencyWalletInvalidInputError("get")
	}

	record, err := scanCurrencyWalletRow(executor.QueryRow(ctx, getCurrencyWalletSQL, string(normalized.WalletID)))
	if err != nil {
		return currency.CurrencyWallet{}, mapCurrencyWalletPostgresError("get", err, nil)
	}
	return record, nil
}

func (r *CurrencyWalletRepository) GetCurrencyWalletForOwner(ctx context.Context, input currency.GetCurrencyWalletForOwnerInput) (currency.CurrencyWallet, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return currency.CurrencyWallet{}, err
	}

	normalized, err := currency.NormalizeGetCurrencyWalletForOwnerInput(input)
	if err != nil {
		return currency.CurrencyWallet{}, currencyWalletInvalidInputError("get_for_owner")
	}

	record, err := scanCurrencyWalletRow(executor.QueryRow(
		ctx,
		getCurrencyWalletForOwnerSQL,
		string(normalized.Owner.Kind),
		normalized.Owner.ID,
	))
	if err != nil {
		return currency.CurrencyWallet{}, mapCurrencyWalletPostgresError("get_for_owner", err, nil)
	}
	return record, nil
}

func (r *CurrencyWalletRepository) ListCurrencyWalletBalances(ctx context.Context, input currency.ListCurrencyWalletBalancesInput) (currency.ListCurrencyWalletBalancesResult, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return currency.ListCurrencyWalletBalancesResult{}, err
	}

	normalized, err := currency.NormalizeListCurrencyWalletBalancesInput(input)
	if err != nil {
		return currency.ListCurrencyWalletBalancesResult{}, currencyWalletInvalidInputError("list_balances")
	}

	rows, err := executor.Query(
		ctx,
		listCurrencyWalletBalancesSQL,
		string(normalized.WalletID),
		string(normalized.AfterCurrencyCode),
		int32(normalized.Limit+1),
	)
	if err != nil {
		return currency.ListCurrencyWalletBalancesResult{}, mapCurrencyWalletPostgresError("list_balances", err, nil)
	}
	defer rows.Close()

	balances := make([]currency.CurrencyWalletBalance, 0, normalized.Limit)
	var nextCurrencyCode currency.CurrencyCode
	for rows.Next() {
		record, err := scanCurrencyWalletBalanceScanner(rows)
		if err != nil {
			return currency.ListCurrencyWalletBalancesResult{}, mapCurrencyWalletPostgresError("list_balances", err, nil)
		}
		if len(balances) >= normalized.Limit {
			nextCurrencyCode = record.CurrencyCode
			continue
		}
		balances = append(balances, record)
	}
	if err := rows.Err(); err != nil {
		return currency.ListCurrencyWalletBalancesResult{}, mapCurrencyWalletPostgresError("list_balances", err, nil)
	}
	return currency.NormalizeListCurrencyWalletBalancesResult(currency.ListCurrencyWalletBalancesResult{
		Balances:         balances,
		NextCurrencyCode: nextCurrencyCode,
	})
}

func (r *CurrencyWalletRepository) RecordCurrencyGrant(ctx context.Context, input currency.RecordCurrencyGrantInput) (currency.CurrencyWalletTransaction, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return currency.CurrencyWalletTransaction{}, err
	}

	normalized, err := currency.NormalizeRecordCurrencyGrantInput(input)
	if err != nil {
		return currency.CurrencyWalletTransaction{}, currencyWalletInvalidInputError("record_grant")
	}

	record, err := scanCurrencyWalletTransactionRow(executor.QueryRow(
		ctx,
		recordCurrencyGrantSQL,
		string(normalized.TransactionID),
		string(normalized.WalletID),
		string(normalized.CurrencyCode),
		string(currency.CurrencyWalletTransactionGrant),
		int64(normalized.Amount),
		string(normalized.IdempotencyKey),
		string(normalized.IdempotencyScope),
		string(normalized.Actor.Kind),
		nullableText(normalized.Actor.ID),
		nullableText(normalized.ReasonCode),
		nullableText(normalized.ExternalReference),
		normalized.MetadataJSON,
		int64(normalized.Amount),
		nullableWalletVersion(normalized.ExpectedWalletVersion),
		nullableBalanceVersion(normalized.ExpectedBalanceVersion),
	))
	if err != nil {
		return currency.CurrencyWalletTransaction{}, mapCurrencyWalletPostgresError("record_grant", err, &mutationConflictHint{kind: currency.CurrencyWalletTransactionGrant})
	}
	return record, nil
}

func (r *CurrencyWalletRepository) RecordCurrencySpend(ctx context.Context, input currency.RecordCurrencySpendInput) (currency.CurrencyWalletTransaction, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return currency.CurrencyWalletTransaction{}, err
	}

	normalized, err := currency.NormalizeRecordCurrencySpendInput(input)
	if err != nil {
		return currency.CurrencyWalletTransaction{}, currencyWalletInvalidInputError("record_spend")
	}

	record, err := scanCurrencyWalletTransactionRow(executor.QueryRow(
		ctx,
		recordCurrencySpendSQL,
		string(normalized.TransactionID),
		string(normalized.WalletID),
		string(normalized.CurrencyCode),
		string(currency.CurrencyWalletTransactionSpend),
		-int64(normalized.Amount),
		string(normalized.IdempotencyKey),
		string(normalized.IdempotencyScope),
		string(normalized.Actor.Kind),
		nullableText(normalized.Actor.ID),
		nullableText(normalized.ReasonCode),
		nullableText(normalized.ExternalReference),
		normalized.MetadataJSON,
		int64(normalized.Amount),
		nullableWalletVersion(normalized.ExpectedWalletVersion),
		nullableBalanceVersion(normalized.ExpectedBalanceVersion),
	))
	if err != nil {
		return currency.CurrencyWalletTransaction{}, mapCurrencyWalletPostgresError("record_spend", err, &mutationConflictHint{kind: currency.CurrencyWalletTransactionSpend})
	}
	return record, nil
}

func (r *CurrencyWalletRepository) ListCurrencyWalletTransactions(ctx context.Context, input currency.ListCurrencyWalletTransactionsInput) (currency.ListCurrencyWalletTransactionsResult, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return currency.ListCurrencyWalletTransactionsResult{}, err
	}

	normalized, err := currency.NormalizeListCurrencyWalletTransactionsInput(input)
	if err != nil {
		return currency.ListCurrencyWalletTransactionsResult{}, currencyWalletInvalidInputError("list_transactions")
	}

	rows, err := executor.Query(
		ctx,
		listCurrencyWalletTransactionsSQL,
		string(normalized.WalletID),
		string(normalized.CurrencyCode),
		normalized.AfterTransactionTime,
		string(normalized.AfterTransactionID),
		int32(normalized.Limit+1),
	)
	if err != nil {
		return currency.ListCurrencyWalletTransactionsResult{}, mapCurrencyWalletPostgresError("list_transactions", err, nil)
	}
	defer rows.Close()

	transactions := make([]currency.CurrencyWalletTransaction, 0, normalized.Limit)
	var nextTransactionID currency.CurrencyWalletTransactionID
	var nextTransactionCreateAt = normalized.AfterTransactionTime
	for rows.Next() {
		record, err := scanCurrencyWalletTransactionScanner(rows)
		if err != nil {
			return currency.ListCurrencyWalletTransactionsResult{}, mapCurrencyWalletPostgresError("list_transactions", err, nil)
		}
		if len(transactions) >= normalized.Limit {
			nextTransactionID = record.TransactionID
			nextTransactionCreateAt = record.CreatedAt
			continue
		}
		transactions = append(transactions, record)
	}
	if err := rows.Err(); err != nil {
		return currency.ListCurrencyWalletTransactionsResult{}, mapCurrencyWalletPostgresError("list_transactions", err, nil)
	}
	return currency.NormalizeListCurrencyWalletTransactionsResult(currency.ListCurrencyWalletTransactionsResult{
		Transactions:            transactions,
		NextTransactionID:       nextTransactionID,
		NextTransactionCreateAt: nextTransactionCreateAt,
	})
}

func (r *CurrencyWalletRepository) requireExecutor() (Executor, error) {
	if r == nil || r.executor == nil {
		return nil, errors.New("postgres currency wallet: unit-of-work executor is required")
	}
	return r.executor, nil
}

func scanCurrencyWalletRow(row pgx.Row) (currency.CurrencyWallet, error) {
	return scanCurrencyWalletScanner(row)
}

func scanCurrencyWalletScanner(row scanner) (currency.CurrencyWallet, error) {
	var record currency.CurrencyWallet
	var walletID string
	var ownerKind string
	var ownerID string
	var lifecycleState string
	var walletVersion int64
	var suspendedAt pgtype.Timestamptz
	var closedAt pgtype.Timestamptz

	if err := row.Scan(
		&walletID,
		&ownerKind,
		&ownerID,
		&lifecycleState,
		&walletVersion,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.StateChangedAt,
		&suspendedAt,
		&closedAt,
	); err != nil {
		return currency.CurrencyWallet{}, err
	}

	record.WalletID = currency.CurrencyWalletID(walletID)
	record.Owner.Kind = currency.CurrencyWalletOwnerKind(strings.TrimSpace(ownerKind))
	record.Owner.ID = ownerID
	record.LifecycleState = currency.CurrencyWalletLifecycleState(strings.TrimSpace(lifecycleState))
	record.WalletVersion = currency.CurrencyWalletVersion(walletVersion)
	record.SuspendedAt = nullableTimestamptzUTC(suspendedAt)
	record.ClosedAt = nullableTimestamptzUTC(closedAt)

	normalized, err := currency.NormalizeCurrencyWalletRecord(record)
	if err != nil {
		return currency.CurrencyWallet{}, fmt.Errorf("%w: row shape", currency.ErrCurrencyWalletUnavailable)
	}
	return normalized, nil
}

func scanCurrencyWalletBalanceScanner(row scanner) (currency.CurrencyWalletBalance, error) {
	var record currency.CurrencyWalletBalance
	var walletID string
	var currencyCode string
	var balanceAmount int64
	var balanceVersion int64

	if err := row.Scan(
		&walletID,
		&currencyCode,
		&balanceAmount,
		&balanceVersion,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return currency.CurrencyWalletBalance{}, err
	}

	record.WalletID = currency.CurrencyWalletID(walletID)
	record.CurrencyCode = currency.CurrencyCode(currencyCode)
	record.BalanceAmount = currency.CurrencyAmount(balanceAmount)
	record.BalanceVersion = currency.CurrencyBalanceVersion(balanceVersion)
	normalized, err := currency.NormalizeCurrencyWalletBalanceRecord(record)
	if err != nil {
		return currency.CurrencyWalletBalance{}, fmt.Errorf("%w: row shape", currency.ErrCurrencyWalletUnavailable)
	}
	return normalized, nil
}

func scanCurrencyWalletTransactionRow(row pgx.Row) (currency.CurrencyWalletTransaction, error) {
	return scanCurrencyWalletTransactionScanner(row)
}

func scanCurrencyWalletTransactionScanner(row scanner) (currency.CurrencyWalletTransaction, error) {
	var record currency.CurrencyWalletTransaction
	var transactionID string
	var walletID string
	var currencyCode string
	var transactionKind string
	var amountDelta int64
	var balanceAfter int64
	var idempotencyKey string
	var idempotencyScope string
	var actorKind string
	var actorID pgtype.Text
	var reasonCode pgtype.Text
	var externalReference pgtype.Text
	var metadataJSON []byte

	if err := row.Scan(
		&transactionID,
		&walletID,
		&currencyCode,
		&transactionKind,
		&amountDelta,
		&balanceAfter,
		&idempotencyKey,
		&idempotencyScope,
		&actorKind,
		&actorID,
		&reasonCode,
		&externalReference,
		&metadataJSON,
		&record.CreatedAt,
	); err != nil {
		return currency.CurrencyWalletTransaction{}, err
	}

	record.TransactionID = currency.CurrencyWalletTransactionID(transactionID)
	record.WalletID = currency.CurrencyWalletID(walletID)
	record.CurrencyCode = currency.CurrencyCode(currencyCode)
	record.TransactionKind = currency.CurrencyWalletTransactionKind(strings.TrimSpace(transactionKind))
	record.AmountDelta = currency.CurrencyAmount(amountDelta)
	record.BalanceAfter = currency.CurrencyAmount(balanceAfter)
	record.IdempotencyKey = currency.CurrencyWalletIdempotencyKey(idempotencyKey)
	record.IdempotencyScope = currency.CurrencyWalletIdempotencyScope(idempotencyScope)
	record.Actor = currency.CurrencyWalletActor{
		Kind: currency.CurrencyWalletActorKind(strings.TrimSpace(actorKind)),
		ID:   nullableTextValue(actorID),
	}
	record.ReasonCode = nullableTextValue(reasonCode)
	record.ExternalReference = nullableTextValue(externalReference)
	record.MetadataJSON = cloneBytes(metadataJSON)

	normalized, err := currency.NormalizeCurrencyWalletTransactionRecord(record)
	if err != nil {
		return currency.CurrencyWalletTransaction{}, fmt.Errorf("%w: row shape", currency.ErrCurrencyWalletUnavailable)
	}
	return normalized, nil
}

type mutationConflictHint struct {
	kind currency.CurrencyWalletTransactionKind
}

func mapCurrencyWalletPostgresError(operation string, err error, hint *mutationConflictHint) error {
	if errors.Is(err, pgx.ErrNoRows) {
		if hint != nil && hint.kind == currency.CurrencyWalletTransactionSpend {
			return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletConflict, currency.CurrencyWalletConflictInsufficientBalance, false)
		}
		if hint != nil {
			return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletConflict, currency.CurrencyWalletConflictWalletNotFound, false)
		}
		return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletConflict, currency.CurrencyWalletConflictWalletNotFound, false)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			if pgErr.ConstraintName == "currency_wallet_transactions_idempotency_uq" {
				return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletConflict, currency.CurrencyWalletConflictDuplicateTransaction, false)
			}
			return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletConflict, currency.CurrencyWalletConflictWalletAlreadyExists, false)
		case "23502", "23503", "23514":
			return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletUnavailable, currency.CurrencyWalletConflictStorageUnavailable, true)
		default:
			return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletUnavailable, currency.CurrencyWalletConflictStorageUnavailable, true)
		}
	}

	if errors.Is(err, currency.ErrCurrencyWalletUnavailable) {
		return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletUnavailable, currency.CurrencyWalletConflictStorageUnavailable, true)
	}
	return currencyWalletRepositoryError(operation, currency.ErrCurrencyWalletUnavailable, currency.CurrencyWalletConflictStorageUnavailable, true)
}

func currencyWalletInvalidInputError(operation string) error {
	return &currency.CurrencyWalletRepositoryError{
		Kind: currency.ErrCurrencyWalletInvalidInput,
		Conflict: currency.CurrencyWalletConflict{
			Class:          currency.CurrencyWalletConflictStorageUnavailable,
			Retryable:      false,
			RedactedReason: "invalid_input",
		},
		Operation:      operation,
		RedactedReason: "invalid_input",
	}
}

func currencyWalletRepositoryError(operation string, kind error, class currency.CurrencyWalletConflictClass, retryable bool) error {
	reason := string(class)
	if reason == "" && kind != nil {
		reason = kind.Error()
	}
	return &currency.CurrencyWalletRepositoryError{
		Kind: kind,
		Conflict: currency.CurrencyWalletConflict{
			Class:          class,
			Retryable:      retryable,
			RedactedReason: reason,
		},
		Operation:      operation,
		RedactedReason: reason,
	}
}

func nullableWalletVersion(version *currency.CurrencyWalletVersion) any {
	if version == nil {
		return nil
	}
	return int64(*version)
}

func nullableBalanceVersion(version *currency.CurrencyBalanceVersion) any {
	if version == nil {
		return nil
	}
	return int64(*version)
}

const currencyWalletColumns = `
wallet_id,
owner_kind,
owner_id,
lifecycle_state,
wallet_version,
created_at,
updated_at,
state_changed_at,
suspended_at,
closed_at`

const currencyWalletBalanceColumns = `
wallet_id,
currency_code,
balance_amount,
balance_version,
created_at,
updated_at`

const currencyWalletTransactionColumns = `
transaction_id,
wallet_id,
currency_code,
transaction_kind,
amount_delta,
balance_after,
idempotency_key,
idempotency_scope,
actor_kind,
actor_id,
reason_code,
external_reference,
metadata_json,
created_at`

const insertCurrencyWalletSQL = `
INSERT INTO currency_wallets (
    wallet_id,
    owner_kind,
    owner_id,
    lifecycle_state,
    wallet_version,
    created_at,
    updated_at,
    state_changed_at
)
VALUES ($1, $2, $3, $4, $5, now(), now(), now())
RETURNING ` + currencyWalletColumns

const getCurrencyWalletSQL = `
SELECT ` + currencyWalletColumns + `
FROM currency_wallets
WHERE wallet_id = $1`

const getCurrencyWalletForOwnerSQL = `
SELECT ` + currencyWalletColumns + `
FROM currency_wallets
WHERE owner_kind = $1
  AND owner_id = $2`

const listCurrencyWalletBalancesSQL = `
SELECT ` + currencyWalletBalanceColumns + `
FROM currency_wallet_balances
WHERE wallet_id = $1
  AND currency_code > $2
ORDER BY currency_code
LIMIT $3`

const recordCurrencyGrantSQL = `
WITH wallet_match AS (
    SELECT wallet_id
    FROM currency_wallets
    WHERE wallet_id = $2
      AND lifecycle_state = 'active'
      AND ($14::bigint IS NULL OR wallet_version = $14)
), balance_upsert AS (
    INSERT INTO currency_wallet_balances (
        wallet_id,
        currency_code,
        balance_amount,
        balance_version,
        created_at,
        updated_at
    )
    SELECT wallet_id, $3, $13, 1, now(), now()
    FROM wallet_match
    ON CONFLICT (wallet_id, currency_code)
    DO UPDATE SET
        balance_amount = currency_wallet_balances.balance_amount + EXCLUDED.balance_amount,
        balance_version = currency_wallet_balances.balance_version + 1,
        updated_at = now()
    WHERE $15::bigint IS NULL OR currency_wallet_balances.balance_version = $15
    RETURNING wallet_id, currency_code, balance_amount
), transaction_insert AS (
    INSERT INTO currency_wallet_transactions (
        transaction_id,
        wallet_id,
        currency_code,
        transaction_kind,
        amount_delta,
        balance_after,
        idempotency_key,
        idempotency_scope,
        actor_kind,
        actor_id,
        reason_code,
        external_reference,
        metadata_json,
        created_at
    )
    SELECT $1, wallet_id, currency_code, $4, $5, balance_amount, $6, $7, $8, $9, $10, $11, $12, now()
    FROM balance_upsert
    RETURNING ` + currencyWalletTransactionColumns + `
)
SELECT ` + currencyWalletTransactionColumns + `
FROM transaction_insert`

const recordCurrencySpendSQL = `
WITH wallet_match AS (
    SELECT wallet_id
    FROM currency_wallets
    WHERE wallet_id = $2
      AND lifecycle_state = 'active'
      AND ($14::bigint IS NULL OR wallet_version = $14)
), balance_update AS (
    UPDATE currency_wallet_balances
    SET balance_amount = balance_amount - $13,
        balance_version = balance_version + 1,
        updated_at = now()
    WHERE wallet_id IN (SELECT wallet_id FROM wallet_match)
      AND currency_code = $3
      AND balance_amount >= $13
      AND ($15::bigint IS NULL OR balance_version = $15)
    RETURNING wallet_id, currency_code, balance_amount
), transaction_insert AS (
    INSERT INTO currency_wallet_transactions (
        transaction_id,
        wallet_id,
        currency_code,
        transaction_kind,
        amount_delta,
        balance_after,
        idempotency_key,
        idempotency_scope,
        actor_kind,
        actor_id,
        reason_code,
        external_reference,
        metadata_json,
        created_at
    )
    SELECT $1, wallet_id, currency_code, $4, $5, balance_amount, $6, $7, $8, $9, $10, $11, $12, now()
    FROM balance_update
    RETURNING ` + currencyWalletTransactionColumns + `
)
SELECT ` + currencyWalletTransactionColumns + `
FROM transaction_insert`

const listCurrencyWalletTransactionsSQL = `
SELECT ` + currencyWalletTransactionColumns + `
FROM currency_wallet_transactions
WHERE wallet_id = $1
  AND ($2::text = '' OR currency_code = $2)
  AND (created_at, transaction_id) > ($3, $4)
ORDER BY created_at, transaction_id
LIMIT $5`
