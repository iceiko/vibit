package currency

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

const defaultRequestedBy = "currency_wallet_service"

type Operation string

const (
	OperationNewService                Operation = "NewService"
	OperationEnsurePlayerWallet        Operation = "EnsurePlayerWallet"
	OperationGetOwnWallet              Operation = "GetOwnWallet"
	OperationListOwnWalletBalances     Operation = "ListOwnWalletBalances"
	OperationGrantCurrency             Operation = "GrantCurrency"
	OperationSpendCurrency             Operation = "SpendCurrency"
	OperationListOwnWalletTransactions Operation = "ListOwnWalletTransactions"
)

type FailureClass string

const (
	FailureClassInvalidRequest        FailureClass = "invalid_request"
	FailureClassUnauthenticated       FailureClass = "unauthenticated"
	FailureClassNotFound              FailureClass = "not_found"
	FailureClassAlreadyExists         FailureClass = "already_exists"
	FailureClassWalletNotActive       FailureClass = "wallet_not_active"
	FailureClassInsufficientBalance   FailureClass = "insufficient_balance"
	FailureClassDuplicateTransaction  FailureClass = "duplicate_transaction"
	FailureClassVersionMismatch       FailureClass = "version_mismatch"
	FailureClassDependencyUnavailable FailureClass = "dependency_unavailable"
)

type PublicErrorCode string

const (
	PublicErrorCurrencyWalletInvalidRequest        PublicErrorCode = "CURRENCY_WALLET_INVALID_REQUEST"
	PublicErrorCurrencyWalletUnauthenticated       PublicErrorCode = "CURRENCY_WALLET_UNAUTHENTICATED"
	PublicErrorCurrencyWalletNotFound              PublicErrorCode = "CURRENCY_WALLET_NOT_FOUND"
	PublicErrorCurrencyWalletAlreadyExists         PublicErrorCode = "CURRENCY_WALLET_ALREADY_EXISTS"
	PublicErrorCurrencyWalletNotActive             PublicErrorCode = "CURRENCY_WALLET_NOT_ACTIVE"
	PublicErrorCurrencyWalletInsufficientBalance   PublicErrorCode = "CURRENCY_WALLET_INSUFFICIENT_BALANCE"
	PublicErrorCurrencyWalletDuplicateTransaction  PublicErrorCode = "CURRENCY_WALLET_DUPLICATE_TRANSACTION"
	PublicErrorCurrencyWalletVersionMismatch       PublicErrorCode = "CURRENCY_WALLET_VERSION_MISMATCH"
	PublicErrorCurrencyWalletUnavailable           PublicErrorCode = "CURRENCY_WALLET_UNAVAILABLE"
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
		return fmt.Sprintf("currency wallet service: %s", e.PublicCode)
	}
	return fmt.Sprintf("currency wallet service: %s: %s", e.Operation, e.PublicCode)
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

type WalletIDGenerator interface {
	GenerateCurrencyWalletID(context.Context) (string, error)
}

type TransactionIDGenerator interface {
	GenerateCurrencyWalletTransactionID(context.Context) (string, error)
}

type ServiceDependencies struct {
	UnitOfWorkRunner      UnitOfWorkRunner
	WalletIDGenerator     WalletIDGenerator
	TransactionIDGenerator TransactionIDGenerator
}

type Service struct {
	unitOfWorkRunner      UnitOfWorkRunner
	walletIDGenerator     WalletIDGenerator
	transactionIDGenerator TransactionIDGenerator
}

func NewService(dependencies ServiceDependencies) (Service, error) {
	if isNilInterface(dependencies.UnitOfWorkRunner) ||
		isNilInterface(dependencies.WalletIDGenerator) ||
		isNilInterface(dependencies.TransactionIDGenerator) {
		return Service{}, serviceFailure(OperationNewService, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errMissingDependency)
	}
	return Service{
		unitOfWorkRunner:      dependencies.UnitOfWorkRunner,
		walletIDGenerator:     dependencies.WalletIDGenerator,
		transactionIDGenerator: dependencies.TransactionIDGenerator,
	}, nil
}

type CurrencyWalletOperationStatus string

const (
	CurrencyWalletOperationStatusRejected CurrencyWalletOperationStatus = "rejected"
	CurrencyWalletOperationStatusEnsured  CurrencyWalletOperationStatus = "ensured"
	CurrencyWalletOperationStatusFound    CurrencyWalletOperationStatus = "found"
	CurrencyWalletOperationStatusListed   CurrencyWalletOperationStatus = "listed"
	CurrencyWalletOperationStatusGranted  CurrencyWalletOperationStatus = "granted"
	CurrencyWalletOperationStatusSpent    CurrencyWalletOperationStatus = "spent"
)

type EnsurePlayerWalletRequest struct {
	Identity app.RequestIdentity
}

type GetOwnWalletRequest struct {
	Identity app.RequestIdentity
}

type ListOwnWalletBalancesRequest struct {
	Identity          app.RequestIdentity
	Limit             int
	AfterCurrencyCode currencymodule.CurrencyCode
}

type GrantCurrencyRequest struct {
	Identity               app.RequestIdentity
	CurrencyCode           currencymodule.CurrencyCode
	Amount                 currencymodule.CurrencyAmount
	IdempotencyKey         currencymodule.CurrencyWalletIdempotencyKey
	IdempotencyScope       currencymodule.CurrencyWalletIdempotencyScope
	SystemActorID          string
	ReasonCode             string
	ExternalReference      string
	MetadataJSON           []byte
	ExpectedWalletVersion  *currencymodule.CurrencyWalletVersion
	ExpectedBalanceVersion *currencymodule.CurrencyBalanceVersion
}

type SpendCurrencyRequest struct {
	Identity               app.RequestIdentity
	CurrencyCode           currencymodule.CurrencyCode
	Amount                 currencymodule.CurrencyAmount
	IdempotencyKey         currencymodule.CurrencyWalletIdempotencyKey
	IdempotencyScope       currencymodule.CurrencyWalletIdempotencyScope
	ReasonCode             string
	ExternalReference      string
	MetadataJSON           []byte
	ExpectedWalletVersion  *currencymodule.CurrencyWalletVersion
	ExpectedBalanceVersion *currencymodule.CurrencyBalanceVersion
}

type ListOwnWalletTransactionsRequest struct {
	Identity             app.RequestIdentity
	CurrencyCode         currencymodule.CurrencyCode
	Limit                int
	AfterTransactionID   currencymodule.CurrencyWalletTransactionID
	AfterTransactionTime time.Time
}

type CurrencyWalletResult struct {
	Status                  CurrencyWalletOperationStatus
	PublicErrorCode         PublicErrorCode
	FailureClass            FailureClass
	Wallet                  currencymodule.CurrencyWallet
	Balance                 currencymodule.CurrencyWalletBalance
	Balances                []currencymodule.CurrencyWalletBalance
	Transaction             currencymodule.CurrencyWalletTransaction
	Transactions            []currencymodule.CurrencyWalletTransaction
	NextCurrencyCode        currencymodule.CurrencyCode
	NextTransactionID       currencymodule.CurrencyWalletTransactionID
	NextTransactionCreateAt time.Time
}

func (s Service) EnsurePlayerWallet(ctx context.Context, request EnsurePlayerWalletRequest) (CurrencyWalletResult, error) {
	owner, err := ownerFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorCurrencyWalletUnauthenticated, FailureClassUnauthenticated),
			serviceFailure(OperationEnsurePlayerWallet, FailureClassUnauthenticated, PublicErrorCurrencyWalletUnauthenticated, err)
	}

	var committedResult CurrencyWalletResult
	var failedResult CurrencyWalletResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := currencyRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationEnsurePlayerWallet, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, err)
		}
		wallet, err := repository.GetCurrencyWalletForOwner(runCtx, currencymodule.GetCurrencyWalletForOwnerInput{Owner: owner})
		if err == nil {
			normalized, normalizeErr := currencymodule.NormalizeCurrencyWalletRecord(wallet)
			if normalizeErr != nil {
				failedResult = rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable)
				return serviceFailure(OperationEnsurePlayerWallet, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errWalletRecordInvalid)
			}
			committedResult = CurrencyWalletResult{Status: CurrencyWalletOperationStatusEnsured, Wallet: normalized}
			return nil
		}
		if repositoryConflictClass(err) != currencymodule.CurrencyWalletConflictWalletNotFound {
			failedResult, err = mapRepositoryFailure(OperationEnsurePlayerWallet, err)
			return err
		}
		walletID, err := s.generatedWalletID(runCtx)
		if err != nil {
			failedResult = rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationEnsurePlayerWallet, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, err)
		}
		input, err := currencymodule.NormalizeCreateCurrencyWalletInput(currencymodule.CreateCurrencyWalletInput{
			WalletID:       currencymodule.CurrencyWalletID(walletID),
			Owner:          owner,
			InitialState:   currencymodule.CurrencyWalletLifecycleActive,
			InitialVersion: currencymodule.InitialCurrencyWalletVersion,
			RequestedBy:    defaultRequestedBy,
		})
		if err != nil {
			failedResult = rejectedResult(PublicErrorCurrencyWalletInvalidRequest, FailureClassInvalidRequest)
			return serviceFailure(OperationEnsurePlayerWallet, FailureClassInvalidRequest, PublicErrorCurrencyWalletInvalidRequest, errInvalidRequest)
		}
		wallet, err = repository.CreateCurrencyWallet(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationEnsurePlayerWallet, err)
			return err
		}
		normalized, err := currencymodule.NormalizeCurrencyWalletRecord(wallet)
		if err != nil {
			failedResult = rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationEnsurePlayerWallet, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errWalletRecordInvalid)
		}
		committedResult = CurrencyWalletResult{Status: CurrencyWalletOperationStatusEnsured, Wallet: normalized}
		return nil
	})
	return resultAfterUnitOfWork(OperationEnsurePlayerWallet, committedResult, failedResult, err)
}

func (s Service) GetOwnWallet(ctx context.Context, request GetOwnWalletRequest) (CurrencyWalletResult, error) {
	return s.withOwnerWallet(ctx, OperationGetOwnWallet, request.Identity, func(_ context.Context, _ currencymodule.Repository, wallet currencymodule.CurrencyWallet) (CurrencyWalletResult, error) {
		return CurrencyWalletResult{Status: CurrencyWalletOperationStatusFound, Wallet: wallet}, nil
	})
}

func (s Service) ListOwnWalletBalances(ctx context.Context, request ListOwnWalletBalancesRequest) (CurrencyWalletResult, error) {
	return s.withOwnerWallet(ctx, OperationListOwnWalletBalances, request.Identity, func(runCtx context.Context, repository currencymodule.Repository, wallet currencymodule.CurrencyWallet) (CurrencyWalletResult, error) {
		input, err := currencymodule.NormalizeListCurrencyWalletBalancesInput(currencymodule.ListCurrencyWalletBalancesInput{
			WalletID:          wallet.WalletID,
			Limit:             request.Limit,
			AfterCurrencyCode: currencymodule.CurrencyCode(strings.TrimSpace(string(request.AfterCurrencyCode))),
		})
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletInvalidRequest, FailureClassInvalidRequest),
				serviceFailure(OperationListOwnWalletBalances, FailureClassInvalidRequest, PublicErrorCurrencyWalletInvalidRequest, errInvalidRequest)
		}
		result, err := repository.ListCurrencyWalletBalances(runCtx, input)
		if err != nil {
			return mapRepositoryFailure(OperationListOwnWalletBalances, err)
		}
		normalized, err := currencymodule.NormalizeListCurrencyWalletBalancesResult(result)
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable),
				serviceFailure(OperationListOwnWalletBalances, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errBalanceRecordInvalid)
		}
		return CurrencyWalletResult{
			Status:           CurrencyWalletOperationStatusListed,
			Balances:         normalized.Balances,
			NextCurrencyCode: normalized.NextCurrencyCode,
		}, nil
	})
}

func (s Service) GrantCurrency(ctx context.Context, request GrantCurrencyRequest) (CurrencyWalletResult, error) {
	transactionID, err := s.generatedTransactionID(ctx)
	if err != nil {
		return rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationGrantCurrency, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, err)
	}
	return s.withOwnerWallet(ctx, OperationGrantCurrency, request.Identity, func(runCtx context.Context, repository currencymodule.Repository, wallet currencymodule.CurrencyWallet) (CurrencyWalletResult, error) {
		input, err := currencymodule.NormalizeRecordCurrencyGrantInput(currencymodule.RecordCurrencyGrantInput{
			WalletID:               wallet.WalletID,
			TransactionID:          currencymodule.CurrencyWalletTransactionID(transactionID),
			CurrencyCode:           request.CurrencyCode,
			Amount:                 request.Amount,
			IdempotencyKey:         request.IdempotencyKey,
			IdempotencyScope:       request.IdempotencyScope,
			Actor:                  currencymodule.CurrencyWalletActor{Kind: currencymodule.CurrencyWalletActorSystem, ID: request.SystemActorID},
			ReasonCode:             request.ReasonCode,
			ExternalReference:      request.ExternalReference,
			MetadataJSON:           request.MetadataJSON,
			ExpectedWalletVersion:  request.ExpectedWalletVersion,
			ExpectedBalanceVersion: request.ExpectedBalanceVersion,
		})
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletInvalidRequest, FailureClassInvalidRequest),
				serviceFailure(OperationGrantCurrency, FailureClassInvalidRequest, PublicErrorCurrencyWalletInvalidRequest, errInvalidRequest)
		}
		transaction, err := repository.RecordCurrencyGrant(runCtx, input)
		if err != nil {
			return mapRepositoryFailure(OperationGrantCurrency, err)
		}
		normalized, err := currencymodule.NormalizeCurrencyWalletTransactionRecord(transaction)
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable),
				serviceFailure(OperationGrantCurrency, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errTransactionRecordInvalid)
		}
		return CurrencyWalletResult{Status: CurrencyWalletOperationStatusGranted, Transaction: normalized}, nil
	})
}

func (s Service) SpendCurrency(ctx context.Context, request SpendCurrencyRequest) (CurrencyWalletResult, error) {
	transactionID, err := s.generatedTransactionID(ctx)
	if err != nil {
		return rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationSpendCurrency, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, err)
	}
	return s.withOwnerWallet(ctx, OperationSpendCurrency, request.Identity, func(runCtx context.Context, repository currencymodule.Repository, wallet currencymodule.CurrencyWallet) (CurrencyWalletResult, error) {
		ownerID := strings.TrimSpace(wallet.Owner.ID)
		input, err := currencymodule.NormalizeRecordCurrencySpendInput(currencymodule.RecordCurrencySpendInput{
			WalletID:               wallet.WalletID,
			TransactionID:          currencymodule.CurrencyWalletTransactionID(transactionID),
			CurrencyCode:           request.CurrencyCode,
			Amount:                 request.Amount,
			IdempotencyKey:         request.IdempotencyKey,
			IdempotencyScope:       request.IdempotencyScope,
			Actor:                  currencymodule.CurrencyWalletActor{Kind: currencymodule.CurrencyWalletActorPlayer, ID: ownerID},
			ReasonCode:             request.ReasonCode,
			ExternalReference:      request.ExternalReference,
			MetadataJSON:           request.MetadataJSON,
			ExpectedWalletVersion:  request.ExpectedWalletVersion,
			ExpectedBalanceVersion: request.ExpectedBalanceVersion,
		})
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletInvalidRequest, FailureClassInvalidRequest),
				serviceFailure(OperationSpendCurrency, FailureClassInvalidRequest, PublicErrorCurrencyWalletInvalidRequest, errInvalidRequest)
		}
		transaction, err := repository.RecordCurrencySpend(runCtx, input)
		if err != nil {
			return mapRepositoryFailure(OperationSpendCurrency, err)
		}
		normalized, err := currencymodule.NormalizeCurrencyWalletTransactionRecord(transaction)
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable),
				serviceFailure(OperationSpendCurrency, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errTransactionRecordInvalid)
		}
		return CurrencyWalletResult{Status: CurrencyWalletOperationStatusSpent, Transaction: normalized}, nil
	})
}

func (s Service) ListOwnWalletTransactions(ctx context.Context, request ListOwnWalletTransactionsRequest) (CurrencyWalletResult, error) {
	return s.withOwnerWallet(ctx, OperationListOwnWalletTransactions, request.Identity, func(runCtx context.Context, repository currencymodule.Repository, wallet currencymodule.CurrencyWallet) (CurrencyWalletResult, error) {
		input, err := currencymodule.NormalizeListCurrencyWalletTransactionsInput(currencymodule.ListCurrencyWalletTransactionsInput{
			WalletID:             wallet.WalletID,
			CurrencyCode:         currencymodule.CurrencyCode(strings.TrimSpace(string(request.CurrencyCode))),
			Limit:                request.Limit,
			AfterTransactionID:   currencymodule.CurrencyWalletTransactionID(strings.TrimSpace(string(request.AfterTransactionID))),
			AfterTransactionTime: request.AfterTransactionTime,
		})
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletInvalidRequest, FailureClassInvalidRequest),
				serviceFailure(OperationListOwnWalletTransactions, FailureClassInvalidRequest, PublicErrorCurrencyWalletInvalidRequest, errInvalidRequest)
		}
		result, err := repository.ListCurrencyWalletTransactions(runCtx, input)
		if err != nil {
			return mapRepositoryFailure(OperationListOwnWalletTransactions, err)
		}
		normalized, err := currencymodule.NormalizeListCurrencyWalletTransactionsResult(result)
		if err != nil {
			return rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable),
				serviceFailure(OperationListOwnWalletTransactions, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errTransactionRecordInvalid)
		}
		return CurrencyWalletResult{
			Status:                  CurrencyWalletOperationStatusListed,
			Transactions:            normalized.Transactions,
			NextTransactionID:       normalized.NextTransactionID,
			NextTransactionCreateAt: normalized.NextTransactionCreateAt,
		}, nil
	})
}

func (s Service) withOwnerWallet(ctx context.Context, operation Operation, identity app.RequestIdentity, fn func(context.Context, currencymodule.Repository, currencymodule.CurrencyWallet) (CurrencyWalletResult, error)) (CurrencyWalletResult, error) {
	owner, err := ownerFromValidatedIdentity(identity)
	if err != nil {
		return rejectedResult(PublicErrorCurrencyWalletUnauthenticated, FailureClassUnauthenticated),
			serviceFailure(operation, FailureClassUnauthenticated, PublicErrorCurrencyWalletUnauthenticated, err)
	}

	var committedResult CurrencyWalletResult
	var failedResult CurrencyWalletResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := currencyRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(operation, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, err)
		}
		wallet, err := repository.GetCurrencyWalletForOwner(runCtx, currencymodule.GetCurrencyWalletForOwnerInput{Owner: owner})
		if err != nil {
			failedResult, err = mapRepositoryFailure(operation, err)
			return err
		}
		normalized, err := currencymodule.NormalizeCurrencyWalletRecord(wallet)
		if err != nil {
			failedResult = rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(operation, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, errWalletRecordInvalid)
		}
		committedResult, err = fn(runCtx, repository, normalized)
		if err != nil {
			failedResult = committedResult
			return err
		}
		return nil
	})
	return resultAfterUnitOfWork(operation, committedResult, failedResult, err)
}

func ownerFromValidatedIdentity(identity app.RequestIdentity) (currencymodule.CurrencyWalletOwner, error) {
	playerID := strings.TrimSpace(identity.PlayerID)
	actorID := strings.TrimSpace(identity.ActorID)
	if identity.Status != app.IdentityValidationValidated ||
		identity.ActorKind != app.ActorKindPlayer ||
		!identity.PlayerIDValidated ||
		playerID == "" ||
		actorID == "" ||
		playerID != actorID {
		return currencymodule.CurrencyWalletOwner{}, errUnauthenticatedIdentity
	}
	return currencymodule.NormalizeCurrencyWalletOwner(currencymodule.CurrencyWalletOwner{
		Kind: currencymodule.CurrencyWalletOwnerKindPlayer,
		ID:   playerID,
	})
}

type currencyUnitOfWork interface {
	NewCurrencyWalletRepository() (currencymodule.Repository, error)
}

func currencyRepositoryFromUnitOfWork(unit tx.UnitOfWork) (currencymodule.Repository, error) {
	repositories, ok := unit.(currencyUnitOfWork)
	if !ok {
		return nil, errMissingCurrencyUnitOfWork
	}
	repository, err := repositories.NewCurrencyWalletRepository()
	if err != nil {
		return nil, err
	}
	if isNilInterface(repository) {
		return nil, errMissingRepository
	}
	return repository, nil
}

func mapRepositoryFailure(operation Operation, err error) (CurrencyWalletResult, error) {
	publicCode, class := publicFailureForRepositoryError(err)
	return rejectedResult(publicCode, class), serviceFailure(operation, class, publicCode, nil)
}

func publicFailureForRepositoryError(err error) (PublicErrorCode, FailureClass) {
	switch repositoryConflictClass(err) {
	case currencymodule.CurrencyWalletConflictWalletNotFound,
		currencymodule.CurrencyWalletConflictBalanceNotFound,
		currencymodule.CurrencyWalletConflictWalletOwnerMismatch:
		return PublicErrorCurrencyWalletNotFound, FailureClassNotFound
	case currencymodule.CurrencyWalletConflictWalletAlreadyExists:
		return PublicErrorCurrencyWalletAlreadyExists, FailureClassAlreadyExists
	case currencymodule.CurrencyWalletConflictWalletNotActive:
		return PublicErrorCurrencyWalletNotActive, FailureClassWalletNotActive
	case currencymodule.CurrencyWalletConflictInvalidCurrencyCode,
		currencymodule.CurrencyWalletConflictInvalidCurrencyAmount,
		currencymodule.CurrencyWalletConflictInvalidIdempotencyKey,
		currencymodule.CurrencyWalletConflictInvalidTransactionKind:
		return PublicErrorCurrencyWalletInvalidRequest, FailureClassInvalidRequest
	case currencymodule.CurrencyWalletConflictInsufficientBalance:
		return PublicErrorCurrencyWalletInsufficientBalance, FailureClassInsufficientBalance
	case currencymodule.CurrencyWalletConflictDuplicateTransaction,
		currencymodule.CurrencyWalletConflictConflictingDuplicateTransaction:
		return PublicErrorCurrencyWalletDuplicateTransaction, FailureClassDuplicateTransaction
	case currencymodule.CurrencyWalletConflictVersionMismatch,
		currencymodule.CurrencyWalletConflictStaleWalletVersion,
		currencymodule.CurrencyWalletConflictStaleBalanceVersion:
		return PublicErrorCurrencyWalletVersionMismatch, FailureClassVersionMismatch
	case currencymodule.CurrencyWalletConflictStorageUnavailable:
		return PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable
	default:
		if errors.Is(err, currencymodule.ErrCurrencyWalletInvalidInput) {
			return PublicErrorCurrencyWalletInvalidRequest, FailureClassInvalidRequest
		}
		if errors.Is(err, currencymodule.ErrCurrencyWalletUnavailable) {
			return PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable
		}
		return PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable
	}
}

func repositoryConflictClass(err error) currencymodule.CurrencyWalletConflictClass {
	var repositoryErr *currencymodule.CurrencyWalletRepositoryError
	if errors.As(err, &repositoryErr) {
		return repositoryErr.Conflict.Class
	}
	var conflict currencymodule.CurrencyWalletConflict
	if errors.As(err, &conflict) {
		return conflict.Class
	}
	return ""
}

func (s Service) generatedWalletID(ctx context.Context) (string, error) {
	walletID, err := s.walletIDGenerator.GenerateCurrencyWalletID(ctx)
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(walletID)
	if trimmed == "" || trimmed != walletID {
		return "", errMalformedWalletID
	}
	return trimmed, nil
}

func (s Service) generatedTransactionID(ctx context.Context) (string, error) {
	transactionID, err := s.transactionIDGenerator.GenerateCurrencyWalletTransactionID(ctx)
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(transactionID)
	if trimmed == "" || trimmed != transactionID {
		return "", errMalformedTransactionID
	}
	return trimmed, nil
}

func resultAfterUnitOfWork(operation Operation, committedResult CurrencyWalletResult, failedResult CurrencyWalletResult, err error) (CurrencyWalletResult, error) {
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedResult(PublicErrorCurrencyWalletUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(operation, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable, err)
	}
	return committedResult, nil
}

func rejectedResult(publicCode PublicErrorCode, class FailureClass) CurrencyWalletResult {
	return CurrencyWalletResult{
		Status:          CurrencyWalletOperationStatusRejected,
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
	errMissingDependency          = errors.New("currency wallet service: dependency is required")
	errUnauthenticatedIdentity    = errors.New("currency wallet service: validated player identity is required")
	errInvalidRequest             = errors.New("currency wallet service: invalid request")
	errMissingCurrencyUnitOfWork  = errors.New("currency wallet service: currency unit-of-work capability is required")
	errMissingRepository          = errors.New("currency wallet service: currency wallet repository is required")
	errWalletRecordInvalid        = errors.New("currency wallet service: wallet record is invalid")
	errBalanceRecordInvalid       = errors.New("currency wallet service: balance record is invalid")
	errTransactionRecordInvalid   = errors.New("currency wallet service: transaction record is invalid")
	errMalformedWalletID          = errors.New("currency wallet service: generated wallet id is malformed")
	errMalformedTransactionID     = errors.New("currency wallet service: generated transaction id is malformed")
)
