package currency

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
	ModuleName = "currency"

	MaxCurrencyCodeLength                   = 64
	MaxCurrencyWalletIdempotencyKeyLength   = 256
	MaxCurrencyWalletIdempotencyScopeLength = 128
	MaxCurrencyWalletReasonCodeLength       = 128
	MaxCurrencyWalletExternalRefLength      = 256

	DefaultListCurrencyWalletBalancesLimit     = 100
	MaxListCurrencyWalletBalancesLimit         = 500
	DefaultListCurrencyWalletTransactionsLimit = 100
	MaxListCurrencyWalletTransactionsLimit     = 500

	InitialCurrencyWalletVersion  CurrencyWalletVersion  = 1
	InitialCurrencyBalanceVersion CurrencyBalanceVersion = 1
)

type CurrencyWalletID string

type CurrencyWalletOwnerKind string

const (
	CurrencyWalletOwnerKindPlayer CurrencyWalletOwnerKind = "player"
)

func (k CurrencyWalletOwnerKind) IsValid() bool {
	return k == CurrencyWalletOwnerKindPlayer
}

type CurrencyWalletOwner struct {
	Kind CurrencyWalletOwnerKind
	ID   string
}

type CurrencyWalletLifecycleState string

const (
	CurrencyWalletLifecycleActive    CurrencyWalletLifecycleState = "active"
	CurrencyWalletLifecycleSuspended CurrencyWalletLifecycleState = "suspended"
	CurrencyWalletLifecycleClosed    CurrencyWalletLifecycleState = "closed"
)

func (s CurrencyWalletLifecycleState) IsValid() bool {
	switch s {
	case CurrencyWalletLifecycleActive,
		CurrencyWalletLifecycleSuspended,
		CurrencyWalletLifecycleClosed:
		return true
	default:
		return false
	}
}

type CurrencyWalletVersion int64

func (v CurrencyWalletVersion) IsValid() bool {
	return v > 0
}

type CurrencyWallet struct {
	WalletID       CurrencyWalletID
	Owner          CurrencyWalletOwner
	LifecycleState CurrencyWalletLifecycleState
	WalletVersion  CurrencyWalletVersion
	CreatedAt      time.Time
	UpdatedAt      time.Time
	StateChangedAt time.Time
	SuspendedAt    *time.Time
	ClosedAt       *time.Time
}

type CurrencyCode string

type CurrencyAmount int64

type CurrencyBalanceVersion int64

func (v CurrencyBalanceVersion) IsValid() bool {
	return v > 0
}

type CurrencyWalletBalance struct {
	WalletID       CurrencyWalletID
	CurrencyCode   CurrencyCode
	BalanceAmount  CurrencyAmount
	BalanceVersion CurrencyBalanceVersion
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CurrencyWalletTransactionID string

type CurrencyWalletTransactionKind string

const (
	CurrencyWalletTransactionGrant CurrencyWalletTransactionKind = "grant"
	CurrencyWalletTransactionSpend CurrencyWalletTransactionKind = "spend"
)

func (k CurrencyWalletTransactionKind) IsValid() bool {
	switch k {
	case CurrencyWalletTransactionGrant, CurrencyWalletTransactionSpend:
		return true
	default:
		return false
	}
}

type CurrencyWalletActorKind string

const (
	CurrencyWalletActorSystem    CurrencyWalletActorKind = "system"
	CurrencyWalletActorPlayer    CurrencyWalletActorKind = "player"
	CurrencyWalletActorOperation CurrencyWalletActorKind = "operation"
)

func (k CurrencyWalletActorKind) IsValid() bool {
	switch k {
	case CurrencyWalletActorSystem,
		CurrencyWalletActorPlayer,
		CurrencyWalletActorOperation:
		return true
	default:
		return false
	}
}

type CurrencyWalletActor struct {
	Kind CurrencyWalletActorKind
	ID   string
}

type CurrencyWalletIdempotencyKey string

type CurrencyWalletIdempotencyScope string

type CurrencyWalletTransaction struct {
	TransactionID     CurrencyWalletTransactionID
	WalletID          CurrencyWalletID
	CurrencyCode      CurrencyCode
	TransactionKind   CurrencyWalletTransactionKind
	AmountDelta       CurrencyAmount
	BalanceAfter      CurrencyAmount
	IdempotencyKey    CurrencyWalletIdempotencyKey
	IdempotencyScope  CurrencyWalletIdempotencyScope
	Actor             CurrencyWalletActor
	ReasonCode        string
	ExternalReference string
	MetadataJSON      []byte
	CreatedAt         time.Time
}

type CreateCurrencyWalletInput struct {
	WalletID       CurrencyWalletID
	Owner          CurrencyWalletOwner
	InitialState   CurrencyWalletLifecycleState
	InitialVersion CurrencyWalletVersion
	RequestedBy    string
}

type GetCurrencyWalletInput struct {
	WalletID CurrencyWalletID
}

type GetCurrencyWalletForOwnerInput struct {
	Owner CurrencyWalletOwner
}

type ListCurrencyWalletBalancesInput struct {
	WalletID          CurrencyWalletID
	Limit             int
	AfterCurrencyCode CurrencyCode
}

type ListCurrencyWalletBalancesResult struct {
	Balances         []CurrencyWalletBalance
	NextCurrencyCode CurrencyCode
}

type RecordCurrencyGrantInput struct {
	WalletID               CurrencyWalletID
	TransactionID          CurrencyWalletTransactionID
	CurrencyCode           CurrencyCode
	Amount                 CurrencyAmount
	IdempotencyKey         CurrencyWalletIdempotencyKey
	IdempotencyScope       CurrencyWalletIdempotencyScope
	Actor                  CurrencyWalletActor
	ReasonCode             string
	ExternalReference      string
	MetadataJSON           []byte
	ExpectedWalletVersion  *CurrencyWalletVersion
	ExpectedBalanceVersion *CurrencyBalanceVersion
}

type RecordCurrencySpendInput struct {
	WalletID               CurrencyWalletID
	TransactionID          CurrencyWalletTransactionID
	CurrencyCode           CurrencyCode
	Amount                 CurrencyAmount
	IdempotencyKey         CurrencyWalletIdempotencyKey
	IdempotencyScope       CurrencyWalletIdempotencyScope
	Actor                  CurrencyWalletActor
	ReasonCode             string
	ExternalReference      string
	MetadataJSON           []byte
	ExpectedWalletVersion  *CurrencyWalletVersion
	ExpectedBalanceVersion *CurrencyBalanceVersion
}

type ListCurrencyWalletTransactionsInput struct {
	WalletID             CurrencyWalletID
	CurrencyCode         CurrencyCode
	Limit                int
	AfterTransactionID   CurrencyWalletTransactionID
	AfterTransactionTime time.Time
}

type ListCurrencyWalletTransactionsResult struct {
	Transactions            []CurrencyWalletTransaction
	NextTransactionID       CurrencyWalletTransactionID
	NextTransactionCreateAt time.Time
}

type Repository interface {
	CreateCurrencyWallet(context.Context, CreateCurrencyWalletInput) (CurrencyWallet, error)
	GetCurrencyWallet(context.Context, GetCurrencyWalletInput) (CurrencyWallet, error)
	GetCurrencyWalletForOwner(context.Context, GetCurrencyWalletForOwnerInput) (CurrencyWallet, error)
	ListCurrencyWalletBalances(context.Context, ListCurrencyWalletBalancesInput) (ListCurrencyWalletBalancesResult, error)
	RecordCurrencyGrant(context.Context, RecordCurrencyGrantInput) (CurrencyWalletTransaction, error)
	RecordCurrencySpend(context.Context, RecordCurrencySpendInput) (CurrencyWalletTransaction, error)
	ListCurrencyWalletTransactions(context.Context, ListCurrencyWalletTransactionsInput) (ListCurrencyWalletTransactionsResult, error)
}

type CurrencyWalletConflictClass string

const (
	CurrencyWalletConflictWalletNotFound                  CurrencyWalletConflictClass = "wallet_not_found"
	CurrencyWalletConflictWalletAlreadyExists             CurrencyWalletConflictClass = "wallet_already_exists"
	CurrencyWalletConflictWalletOwnerMismatch             CurrencyWalletConflictClass = "wallet_owner_mismatch"
	CurrencyWalletConflictWalletNotActive                 CurrencyWalletConflictClass = "wallet_not_active"
	CurrencyWalletConflictBalanceNotFound                 CurrencyWalletConflictClass = "currency_balance_not_found"
	CurrencyWalletConflictInvalidCurrencyCode             CurrencyWalletConflictClass = "invalid_currency_code"
	CurrencyWalletConflictInvalidCurrencyAmount           CurrencyWalletConflictClass = "invalid_currency_amount"
	CurrencyWalletConflictInsufficientBalance             CurrencyWalletConflictClass = "insufficient_balance"
	CurrencyWalletConflictDuplicateTransaction            CurrencyWalletConflictClass = "duplicate_transaction"
	CurrencyWalletConflictConflictingDuplicateTransaction CurrencyWalletConflictClass = "conflicting_duplicate_transaction"
	CurrencyWalletConflictInvalidIdempotencyKey           CurrencyWalletConflictClass = "invalid_idempotency_key"
	CurrencyWalletConflictInvalidTransactionKind          CurrencyWalletConflictClass = "invalid_transaction_kind"
	CurrencyWalletConflictVersionMismatch                 CurrencyWalletConflictClass = "version_mismatch"
	CurrencyWalletConflictStaleWalletVersion              CurrencyWalletConflictClass = "stale_wallet_version"
	CurrencyWalletConflictStaleBalanceVersion             CurrencyWalletConflictClass = "stale_balance_version"
	CurrencyWalletConflictStorageUnavailable              CurrencyWalletConflictClass = "storage_unavailable"
)

func (c CurrencyWalletConflictClass) IsValid() bool {
	switch c {
	case CurrencyWalletConflictWalletNotFound,
		CurrencyWalletConflictWalletAlreadyExists,
		CurrencyWalletConflictWalletOwnerMismatch,
		CurrencyWalletConflictWalletNotActive,
		CurrencyWalletConflictBalanceNotFound,
		CurrencyWalletConflictInvalidCurrencyCode,
		CurrencyWalletConflictInvalidCurrencyAmount,
		CurrencyWalletConflictInsufficientBalance,
		CurrencyWalletConflictDuplicateTransaction,
		CurrencyWalletConflictConflictingDuplicateTransaction,
		CurrencyWalletConflictInvalidIdempotencyKey,
		CurrencyWalletConflictInvalidTransactionKind,
		CurrencyWalletConflictVersionMismatch,
		CurrencyWalletConflictStaleWalletVersion,
		CurrencyWalletConflictStaleBalanceVersion,
		CurrencyWalletConflictStorageUnavailable:
		return true
	default:
		return false
	}
}

type CurrencyWalletConflict struct {
	Class          CurrencyWalletConflictClass
	Retryable      bool
	RedactedReason string
}

func (c CurrencyWalletConflict) Error() string {
	reason := strings.TrimSpace(c.RedactedReason)
	if reason == "" {
		reason = string(c.Class)
	}
	if reason == "" {
		return "currency wallet conflict"
	}
	return fmt.Sprintf("currency wallet conflict: %s", reason)
}

func (c CurrencyWalletConflict) Is(target error) bool {
	targetConflict, ok := target.(CurrencyWalletConflict)
	if !ok {
		return false
	}
	return c.Class == targetConflict.Class && c.Class != ""
}

var (
	ErrCurrencyWalletInvalidInput = errors.New("currency wallet repository: invalid input")
	ErrCurrencyWalletConflict     = errors.New("currency wallet repository: conflict")
	ErrCurrencyWalletUnavailable  = errors.New("currency wallet repository: storage unavailable")
)

type CurrencyWalletRepositoryError struct {
	Kind           error
	Conflict       CurrencyWalletConflict
	Operation      string
	RedactedReason string
	Err            error
}

func (e *CurrencyWalletRepositoryError) Error() string {
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
		return fmt.Sprintf("currency wallet repository %s failed", operation)
	}
	return fmt.Sprintf("currency wallet repository %s failed: %s", operation, reason)
}

func (e *CurrencyWalletRepositoryError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *CurrencyWalletRepositoryError) Is(target error) bool {
	if e == nil {
		return false
	}
	if errors.Is(e.Kind, target) || errors.Is(e.Err, target) {
		return true
	}
	if e.Conflict.Class != "" && errors.Is(e.Conflict, target) {
		return true
	}
	return target == ErrCurrencyWalletConflict && e.Conflict.Class != ""
}

func NormalizeCurrencyWalletRecord(record CurrencyWallet) (CurrencyWallet, error) {
	var err error
	record.WalletID, err = NormalizeCurrencyWalletID(record.WalletID)
	if err != nil {
		return CurrencyWallet{}, err
	}
	record.Owner, err = NormalizeCurrencyWalletOwner(record.Owner)
	if err != nil {
		return CurrencyWallet{}, err
	}
	record.LifecycleState = CurrencyWalletLifecycleState(strings.TrimSpace(string(record.LifecycleState)))
	if !record.LifecycleState.IsValid() {
		return CurrencyWallet{}, errors.New("currency wallet: lifecycle_state is invalid")
	}
	if !record.WalletVersion.IsValid() {
		return CurrencyWallet{}, errors.New("currency wallet: wallet_version must be positive")
	}
	record.CreatedAt, err = normalizeRequiredTime("created_at", record.CreatedAt)
	if err != nil {
		return CurrencyWallet{}, err
	}
	record.UpdatedAt, err = normalizeRequiredTime("updated_at", record.UpdatedAt)
	if err != nil {
		return CurrencyWallet{}, err
	}
	record.StateChangedAt, err = normalizeRequiredTime("state_changed_at", record.StateChangedAt)
	if err != nil {
		return CurrencyWallet{}, err
	}
	if record.UpdatedAt.Before(record.CreatedAt) {
		return CurrencyWallet{}, errors.New("currency wallet: updated_at must not be before created_at")
	}
	if record.StateChangedAt.Before(record.CreatedAt) {
		return CurrencyWallet{}, errors.New("currency wallet: state_changed_at must not be before created_at")
	}
	record.SuspendedAt, err = normalizeOptionalStateTime("suspended_at", record.SuspendedAt, record.CreatedAt)
	if err != nil {
		return CurrencyWallet{}, err
	}
	record.ClosedAt, err = normalizeOptionalStateTime("closed_at", record.ClosedAt, record.CreatedAt)
	if err != nil {
		return CurrencyWallet{}, err
	}
	if record.SuspendedAt != nil && record.LifecycleState != CurrencyWalletLifecycleSuspended {
		return CurrencyWallet{}, errors.New("currency wallet: suspended_at requires suspended state")
	}
	if record.ClosedAt != nil && record.LifecycleState != CurrencyWalletLifecycleClosed {
		return CurrencyWallet{}, errors.New("currency wallet: closed_at requires closed state")
	}
	if record.LifecycleState == CurrencyWalletLifecycleSuspended && record.SuspendedAt == nil {
		return CurrencyWallet{}, errors.New("currency wallet: suspended state requires suspended_at")
	}
	if record.LifecycleState == CurrencyWalletLifecycleClosed && record.ClosedAt == nil {
		return CurrencyWallet{}, errors.New("currency wallet: closed state requires closed_at")
	}
	return record, nil
}

func NormalizeCurrencyWalletBalanceRecord(record CurrencyWalletBalance) (CurrencyWalletBalance, error) {
	var err error
	record.WalletID, err = NormalizeCurrencyWalletID(record.WalletID)
	if err != nil {
		return CurrencyWalletBalance{}, err
	}
	record.CurrencyCode, err = NormalizeCurrencyCode(record.CurrencyCode)
	if err != nil {
		return CurrencyWalletBalance{}, err
	}
	if record.BalanceAmount < 0 {
		return CurrencyWalletBalance{}, errors.New("currency wallet: balance_amount must be non-negative")
	}
	if !record.BalanceVersion.IsValid() {
		return CurrencyWalletBalance{}, errors.New("currency wallet: balance_version must be positive")
	}
	record.CreatedAt, err = normalizeRequiredTime("created_at", record.CreatedAt)
	if err != nil {
		return CurrencyWalletBalance{}, err
	}
	record.UpdatedAt, err = normalizeRequiredTime("updated_at", record.UpdatedAt)
	if err != nil {
		return CurrencyWalletBalance{}, err
	}
	if record.UpdatedAt.Before(record.CreatedAt) {
		return CurrencyWalletBalance{}, errors.New("currency wallet: updated_at must not be before created_at")
	}
	return record, nil
}

func NormalizeCurrencyWalletTransactionRecord(record CurrencyWalletTransaction) (CurrencyWalletTransaction, error) {
	var err error
	record.TransactionID, err = NormalizeCurrencyWalletTransactionID(record.TransactionID)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.WalletID, err = NormalizeCurrencyWalletID(record.WalletID)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.CurrencyCode, err = NormalizeCurrencyCode(record.CurrencyCode)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.TransactionKind = CurrencyWalletTransactionKind(strings.TrimSpace(string(record.TransactionKind)))
	if !record.TransactionKind.IsValid() {
		return CurrencyWalletTransaction{}, errors.New("currency wallet: transaction_kind is invalid")
	}
	if err := validateTransactionDelta(record.TransactionKind, record.AmountDelta); err != nil {
		return CurrencyWalletTransaction{}, err
	}
	if record.BalanceAfter < 0 {
		return CurrencyWalletTransaction{}, errors.New("currency wallet: balance_after must be non-negative")
	}
	record.IdempotencyKey, err = NormalizeCurrencyWalletIdempotencyKey(record.IdempotencyKey)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.IdempotencyScope, err = NormalizeCurrencyWalletIdempotencyScope(record.IdempotencyScope)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.Actor, err = NormalizeCurrencyWalletActor(record.Actor)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.ReasonCode, err = normalizeOptionalBounded("reason_code", record.ReasonCode, MaxCurrencyWalletReasonCodeLength)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.ExternalReference, err = normalizeOptionalBounded("external_reference", record.ExternalReference, MaxCurrencyWalletExternalRefLength)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.MetadataJSON, err = NormalizeCurrencyWalletMetadataJSON(record.MetadataJSON)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	record.CreatedAt, err = normalizeRequiredTime("created_at", record.CreatedAt)
	if err != nil {
		return CurrencyWalletTransaction{}, err
	}
	return record, nil
}

func NormalizeCreateCurrencyWalletInput(input CreateCurrencyWalletInput) (CreateCurrencyWalletInput, error) {
	var err error
	input.WalletID, err = NormalizeCurrencyWalletID(input.WalletID)
	if err != nil {
		return CreateCurrencyWalletInput{}, err
	}
	input.Owner, err = NormalizeCurrencyWalletOwner(input.Owner)
	if err != nil {
		return CreateCurrencyWalletInput{}, err
	}
	if strings.TrimSpace(string(input.InitialState)) == "" {
		input.InitialState = CurrencyWalletLifecycleActive
	} else {
		input.InitialState = CurrencyWalletLifecycleState(strings.TrimSpace(string(input.InitialState)))
	}
	if input.InitialState != CurrencyWalletLifecycleActive {
		return CreateCurrencyWalletInput{}, errors.New("currency wallet: created wallets must start active")
	}
	if input.InitialVersion == 0 {
		input.InitialVersion = InitialCurrencyWalletVersion
	}
	if !input.InitialVersion.IsValid() {
		return CreateCurrencyWalletInput{}, errors.New("currency wallet: initial_version must be positive")
	}
	input.RequestedBy, err = normalizeRequired("requested_by", input.RequestedBy)
	if err != nil {
		return CreateCurrencyWalletInput{}, err
	}
	return input, nil
}

func NormalizeGetCurrencyWalletInput(input GetCurrencyWalletInput) (GetCurrencyWalletInput, error) {
	var err error
	input.WalletID, err = NormalizeCurrencyWalletID(input.WalletID)
	if err != nil {
		return GetCurrencyWalletInput{}, err
	}
	return input, nil
}

func NormalizeGetCurrencyWalletForOwnerInput(input GetCurrencyWalletForOwnerInput) (GetCurrencyWalletForOwnerInput, error) {
	var err error
	input.Owner, err = NormalizeCurrencyWalletOwner(input.Owner)
	if err != nil {
		return GetCurrencyWalletForOwnerInput{}, err
	}
	return input, nil
}

func NormalizeListCurrencyWalletBalancesInput(input ListCurrencyWalletBalancesInput) (ListCurrencyWalletBalancesInput, error) {
	var err error
	input.WalletID, err = NormalizeCurrencyWalletID(input.WalletID)
	if err != nil {
		return ListCurrencyWalletBalancesInput{}, err
	}
	if input.Limit == 0 {
		input.Limit = DefaultListCurrencyWalletBalancesLimit
	}
	if input.Limit < 0 || input.Limit > MaxListCurrencyWalletBalancesLimit {
		return ListCurrencyWalletBalancesInput{}, errors.New("currency wallet: balance list limit is invalid")
	}
	if strings.TrimSpace(string(input.AfterCurrencyCode)) != "" {
		input.AfterCurrencyCode, err = NormalizeCurrencyCode(input.AfterCurrencyCode)
		if err != nil {
			return ListCurrencyWalletBalancesInput{}, err
		}
	}
	return input, nil
}

func NormalizeListCurrencyWalletBalancesResult(result ListCurrencyWalletBalancesResult) (ListCurrencyWalletBalancesResult, error) {
	if len(result.Balances) > 0 {
		balances := make([]CurrencyWalletBalance, len(result.Balances))
		for i, balance := range result.Balances {
			normalized, err := NormalizeCurrencyWalletBalanceRecord(balance)
			if err != nil {
				return ListCurrencyWalletBalancesResult{}, err
			}
			balances[i] = normalized
		}
		result.Balances = balances
	}
	if strings.TrimSpace(string(result.NextCurrencyCode)) != "" {
		next, err := NormalizeCurrencyCode(result.NextCurrencyCode)
		if err != nil {
			return ListCurrencyWalletBalancesResult{}, err
		}
		result.NextCurrencyCode = next
	}
	return result, nil
}

func NormalizeRecordCurrencyGrantInput(input RecordCurrencyGrantInput) (RecordCurrencyGrantInput, error) {
	var err error
	input.WalletID, input.TransactionID, input.CurrencyCode, err = normalizeMutationIdentity(input.WalletID, input.TransactionID, input.CurrencyCode)
	if err != nil {
		return RecordCurrencyGrantInput{}, err
	}
	if input.Amount <= 0 {
		return RecordCurrencyGrantInput{}, errors.New("currency wallet: grant amount must be positive")
	}
	input.IdempotencyKey, input.IdempotencyScope, input.Actor, input.ReasonCode, input.ExternalReference, input.MetadataJSON, err = normalizeMutationMetadata(input.IdempotencyKey, input.IdempotencyScope, input.Actor, input.ReasonCode, input.ExternalReference, input.MetadataJSON)
	if err != nil {
		return RecordCurrencyGrantInput{}, err
	}
	input.ExpectedWalletVersion, err = copyValidWalletVersion(input.ExpectedWalletVersion)
	if err != nil {
		return RecordCurrencyGrantInput{}, err
	}
	input.ExpectedBalanceVersion, err = copyValidBalanceVersion(input.ExpectedBalanceVersion)
	if err != nil {
		return RecordCurrencyGrantInput{}, err
	}
	return input, nil
}

func NormalizeRecordCurrencySpendInput(input RecordCurrencySpendInput) (RecordCurrencySpendInput, error) {
	var err error
	input.WalletID, input.TransactionID, input.CurrencyCode, err = normalizeMutationIdentity(input.WalletID, input.TransactionID, input.CurrencyCode)
	if err != nil {
		return RecordCurrencySpendInput{}, err
	}
	if input.Amount <= 0 {
		return RecordCurrencySpendInput{}, errors.New("currency wallet: spend amount must be positive")
	}
	input.IdempotencyKey, input.IdempotencyScope, input.Actor, input.ReasonCode, input.ExternalReference, input.MetadataJSON, err = normalizeMutationMetadata(input.IdempotencyKey, input.IdempotencyScope, input.Actor, input.ReasonCode, input.ExternalReference, input.MetadataJSON)
	if err != nil {
		return RecordCurrencySpendInput{}, err
	}
	input.ExpectedWalletVersion, err = copyValidWalletVersion(input.ExpectedWalletVersion)
	if err != nil {
		return RecordCurrencySpendInput{}, err
	}
	input.ExpectedBalanceVersion, err = copyValidBalanceVersion(input.ExpectedBalanceVersion)
	if err != nil {
		return RecordCurrencySpendInput{}, err
	}
	return input, nil
}

func NormalizeListCurrencyWalletTransactionsInput(input ListCurrencyWalletTransactionsInput) (ListCurrencyWalletTransactionsInput, error) {
	var err error
	input.WalletID, err = NormalizeCurrencyWalletID(input.WalletID)
	if err != nil {
		return ListCurrencyWalletTransactionsInput{}, err
	}
	if strings.TrimSpace(string(input.CurrencyCode)) != "" {
		input.CurrencyCode, err = NormalizeCurrencyCode(input.CurrencyCode)
		if err != nil {
			return ListCurrencyWalletTransactionsInput{}, err
		}
	}
	if input.Limit == 0 {
		input.Limit = DefaultListCurrencyWalletTransactionsLimit
	}
	if input.Limit < 0 || input.Limit > MaxListCurrencyWalletTransactionsLimit {
		return ListCurrencyWalletTransactionsInput{}, errors.New("currency wallet: transaction list limit is invalid")
	}
	if strings.TrimSpace(string(input.AfterTransactionID)) != "" {
		input.AfterTransactionID, err = NormalizeCurrencyWalletTransactionID(input.AfterTransactionID)
		if err != nil {
			return ListCurrencyWalletTransactionsInput{}, err
		}
	}
	if !input.AfterTransactionTime.IsZero() {
		input.AfterTransactionTime = input.AfterTransactionTime.UTC()
	}
	return input, nil
}

func NormalizeListCurrencyWalletTransactionsResult(result ListCurrencyWalletTransactionsResult) (ListCurrencyWalletTransactionsResult, error) {
	if len(result.Transactions) > 0 {
		transactions := make([]CurrencyWalletTransaction, len(result.Transactions))
		for i, transaction := range result.Transactions {
			normalized, err := NormalizeCurrencyWalletTransactionRecord(transaction)
			if err != nil {
				return ListCurrencyWalletTransactionsResult{}, err
			}
			transactions[i] = normalized
		}
		result.Transactions = transactions
	}
	if strings.TrimSpace(string(result.NextTransactionID)) != "" {
		next, err := NormalizeCurrencyWalletTransactionID(result.NextTransactionID)
		if err != nil {
			return ListCurrencyWalletTransactionsResult{}, err
		}
		result.NextTransactionID = next
	}
	if !result.NextTransactionCreateAt.IsZero() {
		result.NextTransactionCreateAt = result.NextTransactionCreateAt.UTC()
	}
	return result, nil
}

func NormalizeCurrencyWalletID(id CurrencyWalletID) (CurrencyWalletID, error) {
	value, err := normalizeRequired("wallet_id", string(id))
	if err != nil {
		return "", err
	}
	return CurrencyWalletID(value), nil
}

func NormalizeCurrencyWalletTransactionID(id CurrencyWalletTransactionID) (CurrencyWalletTransactionID, error) {
	value, err := normalizeRequired("transaction_id", string(id))
	if err != nil {
		return "", err
	}
	return CurrencyWalletTransactionID(value), nil
}

func NormalizeCurrencyWalletOwner(owner CurrencyWalletOwner) (CurrencyWalletOwner, error) {
	owner.Kind = CurrencyWalletOwnerKind(strings.TrimSpace(string(owner.Kind)))
	if !owner.Kind.IsValid() {
		return CurrencyWalletOwner{}, errors.New("currency wallet: owner_kind is invalid")
	}
	var err error
	owner.ID, err = normalizeRequired("owner_id", owner.ID)
	if err != nil {
		return CurrencyWalletOwner{}, err
	}
	return owner, nil
}

func NormalizeCurrencyCode(code CurrencyCode) (CurrencyCode, error) {
	value, err := normalizeBounded("currency_code", string(code), MaxCurrencyCodeLength)
	if err != nil {
		return "", err
	}
	return CurrencyCode(value), nil
}

func NormalizeCurrencyWalletIdempotencyKey(key CurrencyWalletIdempotencyKey) (CurrencyWalletIdempotencyKey, error) {
	value, err := normalizeBounded("idempotency_key", string(key), MaxCurrencyWalletIdempotencyKeyLength)
	if err != nil {
		return "", err
	}
	return CurrencyWalletIdempotencyKey(value), nil
}

func NormalizeCurrencyWalletIdempotencyScope(scope CurrencyWalletIdempotencyScope) (CurrencyWalletIdempotencyScope, error) {
	value, err := normalizeBounded("idempotency_scope", string(scope), MaxCurrencyWalletIdempotencyScopeLength)
	if err != nil {
		return "", err
	}
	return CurrencyWalletIdempotencyScope(value), nil
}

func NormalizeCurrencyWalletActor(actor CurrencyWalletActor) (CurrencyWalletActor, error) {
	actor.Kind = CurrencyWalletActorKind(strings.TrimSpace(string(actor.Kind)))
	if !actor.Kind.IsValid() {
		return CurrencyWalletActor{}, errors.New("currency wallet: actor_kind is invalid")
	}
	actor.ID = strings.TrimSpace(actor.ID)
	if actor.Kind == CurrencyWalletActorPlayer && actor.ID == "" {
		return CurrencyWalletActor{}, errors.New("currency wallet: actor_id is required for player actors")
	}
	return actor, nil
}

func NormalizeCurrencyWalletMetadataJSON(metadata []byte) ([]byte, error) {
	if len(bytes.TrimSpace(metadata)) == 0 {
		return nil, nil
	}
	trimmed := bytes.TrimSpace(metadata)
	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &decoded); err != nil || decoded == nil {
		return nil, errors.New("currency wallet: metadata_json must be a JSON object")
	}
	return cloneBytes(trimmed), nil
}

func normalizeMutationIdentity(walletID CurrencyWalletID, transactionID CurrencyWalletTransactionID, currencyCode CurrencyCode) (CurrencyWalletID, CurrencyWalletTransactionID, CurrencyCode, error) {
	normalizedWalletID, err := NormalizeCurrencyWalletID(walletID)
	if err != nil {
		return "", "", "", err
	}
	normalizedTransactionID, err := NormalizeCurrencyWalletTransactionID(transactionID)
	if err != nil {
		return "", "", "", err
	}
	normalizedCurrencyCode, err := NormalizeCurrencyCode(currencyCode)
	if err != nil {
		return "", "", "", err
	}
	return normalizedWalletID, normalizedTransactionID, normalizedCurrencyCode, nil
}

func normalizeMutationMetadata(key CurrencyWalletIdempotencyKey, scope CurrencyWalletIdempotencyScope, actor CurrencyWalletActor, reasonCode string, externalReference string, metadata []byte) (CurrencyWalletIdempotencyKey, CurrencyWalletIdempotencyScope, CurrencyWalletActor, string, string, []byte, error) {
	normalizedKey, err := NormalizeCurrencyWalletIdempotencyKey(key)
	if err != nil {
		return "", "", CurrencyWalletActor{}, "", "", nil, err
	}
	normalizedScope, err := NormalizeCurrencyWalletIdempotencyScope(scope)
	if err != nil {
		return "", "", CurrencyWalletActor{}, "", "", nil, err
	}
	normalizedActor, err := NormalizeCurrencyWalletActor(actor)
	if err != nil {
		return "", "", CurrencyWalletActor{}, "", "", nil, err
	}
	normalizedReason, err := normalizeOptionalBounded("reason_code", reasonCode, MaxCurrencyWalletReasonCodeLength)
	if err != nil {
		return "", "", CurrencyWalletActor{}, "", "", nil, err
	}
	normalizedExternal, err := normalizeOptionalBounded("external_reference", externalReference, MaxCurrencyWalletExternalRefLength)
	if err != nil {
		return "", "", CurrencyWalletActor{}, "", "", nil, err
	}
	normalizedMetadata, err := NormalizeCurrencyWalletMetadataJSON(metadata)
	if err != nil {
		return "", "", CurrencyWalletActor{}, "", "", nil, err
	}
	return normalizedKey, normalizedScope, normalizedActor, normalizedReason, normalizedExternal, normalizedMetadata, nil
}

func validateTransactionDelta(kind CurrencyWalletTransactionKind, amount CurrencyAmount) error {
	switch kind {
	case CurrencyWalletTransactionGrant:
		if amount <= 0 {
			return errors.New("currency wallet: grant amount_delta must be positive")
		}
	case CurrencyWalletTransactionSpend:
		if amount >= 0 {
			return errors.New("currency wallet: spend amount_delta must be negative")
		}
	default:
		return errors.New("currency wallet: transaction_kind is invalid")
	}
	return nil
}

func copyValidWalletVersion(version *CurrencyWalletVersion) (*CurrencyWalletVersion, error) {
	if version == nil {
		return nil, nil
	}
	if !version.IsValid() {
		return nil, errors.New("currency wallet: expected_wallet_version must be positive")
	}
	copied := *version
	return &copied, nil
}

func copyValidBalanceVersion(version *CurrencyBalanceVersion) (*CurrencyBalanceVersion, error) {
	if version == nil {
		return nil, nil
	}
	if !version.IsValid() {
		return nil, errors.New("currency wallet: expected_balance_version must be positive")
	}
	copied := *version
	return &copied, nil
}

func normalizeOptionalStateTime(name string, value *time.Time, createdAt time.Time) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	normalized := value.UTC()
	if normalized.Before(createdAt) {
		return nil, fmt.Errorf("currency wallet: %s must not be before created_at", name)
	}
	return &normalized, nil
}

func normalizeRequired(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("currency wallet: %s is required", name)
	}
	return value, nil
}

func normalizeBounded(name string, value string, maxLength int) (string, error) {
	value, err := normalizeRequired(name, value)
	if err != nil {
		return "", err
	}
	if len(value) > maxLength {
		return "", fmt.Errorf("currency wallet: %s is too long", name)
	}
	return value, nil
}

func normalizeOptionalBounded(name string, value string, maxLength int) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}
	return normalizeBounded(name, value, maxLength)
}

func normalizeRequiredTime(name string, value time.Time) (time.Time, error) {
	if value.IsZero() {
		return time.Time{}, fmt.Errorf("currency wallet: %s is required", name)
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
