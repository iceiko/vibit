package currency

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

func TestNewServiceRequiresDependencies(t *testing.T) {
	_, err := NewService(ServiceDependencies{})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable)

	var runner *recordingUnitOfWorkRunner
	_, err = NewService(ServiceDependencies{
		UnitOfWorkRunner:      runner,
		WalletIDGenerator:     staticWalletIDGenerator{value: "wallet-1"},
		TransactionIDGenerator: staticTransactionIDGenerator{value: "tx-1"},
	})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable)

	_, err = NewService(ServiceDependencies{
		UnitOfWorkRunner:  &recordingUnitOfWorkRunner{},
		WalletIDGenerator: staticWalletIDGenerator{value: "wallet-1"},
	})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable)
}

func TestGetOwnWalletRejectsMetadataOnlyIdentityBeforeUnitOfWork(t *testing.T) {
	repository := &fakeCurrencyWalletRepository{}
	runner := &recordingUnitOfWorkRunner{unit: &fakeCurrencyUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticWalletIDGenerator{value: "wallet-1"}, staticTransactionIDGenerator{value: "tx-1"})

	result, err := service.GetOwnWallet(context.Background(), GetOwnWalletRequest{
		Identity: app.MetadataOnlyIdentityFromSession(app.Session{
			PlayerID:  "player-1",
			SessionID: "metadata-session",
		}),
	})

	assertServiceError(t, err, OperationGetOwnWallet, FailureClassUnauthenticated, PublicErrorCurrencyWalletUnauthenticated)
	if result.Status != CurrencyWalletOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorCurrencyWalletUnauthenticated ||
		result.FailureClass != FailureClassUnauthenticated {
		t.Fatalf("result = %#v, want unauthenticated rejection", result)
	}
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if repository.totalCalls() != 0 {
		t.Fatalf("repository calls = %d, want 0", repository.totalCalls())
	}
}

func TestEnsurePlayerWalletCreatesWhenMissingWithGeneratedID(t *testing.T) {
	events := []string{}
	repository := &fakeCurrencyWalletRepository{
		events:      &events,
		getOwnerErr: currencyRepositoryError(currencymodule.CurrencyWalletConflictWalletNotFound, "owner_id=player-1"),
		createResult: activeCurrencyWallet("wallet-1", "player-1",
			withWalletVersion(currencymodule.InitialCurrencyWalletVersion)),
	}
	runner := &recordingUnitOfWorkRunner{
		unit:   &fakeCurrencyUnitOfWork{repository: repository, events: &events},
		events: &events,
	}
	service := mustNewService(t, runner, staticWalletIDGenerator{value: "wallet-1", events: &events}, staticTransactionIDGenerator{value: "tx-1"})

	result, err := service.EnsurePlayerWallet(context.Background(), EnsurePlayerWalletRequest{
		Identity: validatedPlayerIdentity("player-1"),
	})
	if err != nil {
		t.Fatalf("EnsurePlayerWallet() error = %v, want nil", err)
	}
	if result.Status != CurrencyWalletOperationStatusEnsured ||
		result.Wallet.WalletID != "wallet-1" ||
		result.Wallet.Owner.ID != "player-1" ||
		result.Wallet.WalletVersion != currencymodule.InitialCurrencyWalletVersion {
		t.Fatalf("result = %#v, want ensured player wallet", result)
	}
	if repository.getOwnerCalls != 1 || repository.createCalls != 1 {
		t.Fatalf("getOwner/create calls = %d/%d, want 1/1", repository.getOwnerCalls, repository.createCalls)
	}
	if repository.lastCreateInput.WalletID != "wallet-1" ||
		repository.lastCreateInput.Owner != (currencymodule.CurrencyWalletOwner{Kind: currencymodule.CurrencyWalletOwnerKindPlayer, ID: "player-1"}) ||
		repository.lastCreateInput.InitialState != currencymodule.CurrencyWalletLifecycleActive ||
		repository.lastCreateInput.InitialVersion != currencymodule.InitialCurrencyWalletVersion ||
		repository.lastCreateInput.RequestedBy != defaultRequestedBy {
		t.Fatalf("CreateCurrencyWallet input = %#v, want server-derived player wallet", repository.lastCreateInput)
	}
	assertEvents(t, events, []string{"begin", "new-currency-wallet-repository", "get-currency-wallet-for-owner", "generate-wallet-id", "create-currency-wallet", "commit"})
}

func TestListOwnWalletBalancesUsesServerDerivedOwnerAndDefaultsLimit(t *testing.T) {
	sourceBalance := activeCurrencyWalletBalance("wallet-1", "Gems", 25, 3)
	repository := &fakeCurrencyWalletRepository{
		getOwnerResult: activeCurrencyWallet("wallet-1", "player-1"),
		listBalancesResult: currencymodule.ListCurrencyWalletBalancesResult{
			Balances:         []currencymodule.CurrencyWalletBalance{sourceBalance},
			NextCurrencyCode: "Gold",
		},
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeCurrencyUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticWalletIDGenerator{value: "wallet-1"}, staticTransactionIDGenerator{value: "tx-1"})

	result, err := service.ListOwnWalletBalances(context.Background(), ListOwnWalletBalancesRequest{
		Identity:          validatedPlayerIdentity("player-1"),
		Limit:             0,
		AfterCurrencyCode: " ",
	})
	if err != nil {
		t.Fatalf("ListOwnWalletBalances() error = %v, want nil", err)
	}
	if result.Status != CurrencyWalletOperationStatusListed ||
		len(result.Balances) != 1 ||
		result.Balances[0].CurrencyCode != "Gems" ||
		result.NextCurrencyCode != "Gold" {
		t.Fatalf("result = %#v, want listed balances", result)
	}
	if repository.lastListBalancesInput.WalletID != "wallet-1" ||
		repository.lastListBalancesInput.Limit != currencymodule.DefaultListCurrencyWalletBalancesLimit ||
		repository.lastListBalancesInput.AfterCurrencyCode != "" {
		t.Fatalf("ListCurrencyWalletBalances input = %#v, want wallet id and default limit", repository.lastListBalancesInput)
	}
}

func TestGrantCurrencyUsesSystemActorAndCopiesMetadata(t *testing.T) {
	events := []string{}
	metadata := []byte(`{"source":"daily-login"}`)
	transaction := activeCurrencyWalletTransaction("tx-1", "wallet-1", currencymodule.CurrencyWalletTransactionGrant, "Gems", 10, 35,
		withTransactionActor(currencymodule.CurrencyWalletActorSystem, "economy"),
		withTransactionMetadata(metadata),
	)
	repository := &fakeCurrencyWalletRepository{
		events:            &events,
		getOwnerResult:    activeCurrencyWallet("wallet-1", "player-1"),
		recordGrantResult: transaction,
	}
	runner := &recordingUnitOfWorkRunner{
		unit:   &fakeCurrencyUnitOfWork{repository: repository, events: &events},
		events: &events,
	}
	service := mustNewService(t, runner, staticWalletIDGenerator{value: "wallet-1"}, staticTransactionIDGenerator{value: "tx-1", events: &events})

	result, err := service.GrantCurrency(context.Background(), GrantCurrencyRequest{
		Identity:          validatedPlayerIdentity("player-1"),
		CurrencyCode:      " Gems ",
		Amount:            10,
		IdempotencyKey:    " daily-1 ",
		IdempotencyScope:  " daily-login ",
		SystemActorID:     " economy ",
		ReasonCode:        " daily_login ",
		ExternalReference: " season-1 ",
		MetadataJSON:      metadata,
	})
	if err != nil {
		t.Fatalf("GrantCurrency() error = %v, want nil", err)
	}
	if result.Status != CurrencyWalletOperationStatusGranted ||
		result.Transaction.TransactionID != "tx-1" ||
		result.Transaction.Actor.Kind != currencymodule.CurrencyWalletActorSystem ||
		result.Transaction.Actor.ID != "economy" {
		t.Fatalf("result = %#v, want granted system transaction", result)
	}
	input := repository.lastRecordGrantInput
	if input.WalletID != "wallet-1" ||
		input.TransactionID != "tx-1" ||
		input.CurrencyCode != "Gems" ||
		input.Amount != 10 ||
		input.IdempotencyKey != "daily-1" ||
		input.IdempotencyScope != "daily-login" ||
		input.Actor.Kind != currencymodule.CurrencyWalletActorSystem ||
		input.Actor.ID != "economy" ||
		input.ReasonCode != "daily_login" ||
		input.ExternalReference != "season-1" {
		t.Fatalf("RecordCurrencyGrant input = %#v, want normalized grant input", input)
	}
	metadata[0] = '['
	if string(input.MetadataJSON) != `{"source":"daily-login"}` {
		t.Fatalf("grant metadata = %q, want copied request metadata", string(input.MetadataJSON))
	}
	assertEvents(t, events, []string{"generate-transaction-id", "begin", "new-currency-wallet-repository", "get-currency-wallet-for-owner", "record-currency-grant", "commit"})
}

func TestSpendCurrencyUsesPlayerActorAndMapsInsufficientBalance(t *testing.T) {
	rawDetail := `wallet_id=wallet-1 balance_after=-7 metadata_json={"secret":true}`
	repository := &fakeCurrencyWalletRepository{
		getOwnerResult:  activeCurrencyWallet("wallet-1", "player-1"),
		recordSpendErr: currencyRepositoryError(currencymodule.CurrencyWalletConflictInsufficientBalance, rawDetail),
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeCurrencyUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticWalletIDGenerator{value: "wallet-1"}, staticTransactionIDGenerator{value: "tx-spend-1"})

	result, err := service.SpendCurrency(context.Background(), SpendCurrencyRequest{
		Identity:         validatedPlayerIdentity("player-1"),
		CurrencyCode:     "Gems",
		Amount:           40,
		IdempotencyKey:   "spend-1",
		IdempotencyScope: "shop",
	})

	assertServiceError(t, err, OperationSpendCurrency, FailureClassInsufficientBalance, PublicErrorCurrencyWalletInsufficientBalance)
	assertNoLeak(t, err, rawDetail)
	if result.Status != CurrencyWalletOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorCurrencyWalletInsufficientBalance {
		t.Fatalf("result = %#v, want insufficient-balance rejection", result)
	}
	if repository.recordSpendCalls != 1 {
		t.Fatalf("RecordCurrencySpend calls = %d, want 1", repository.recordSpendCalls)
	}
	input := repository.lastRecordSpendInput
	if input.Actor.Kind != currencymodule.CurrencyWalletActorPlayer ||
		input.Actor.ID != "player-1" ||
		input.Amount != 40 ||
		input.TransactionID != "tx-spend-1" {
		t.Fatalf("RecordCurrencySpend input = %#v, want player spend input", input)
	}
}

func TestListOwnWalletTransactionsUsesValidatedOwnerAndCopiesResults(t *testing.T) {
	createdAt := fixedCurrencyWalletTime().Add(5 * time.Minute)
	sourceMetadata := []byte(`{"match":"m-1"}`)
	transaction := activeCurrencyWalletTransaction("tx-1", "wallet-1", currencymodule.CurrencyWalletTransactionSpend, "Gems", -5, 20,
		withTransactionActor(currencymodule.CurrencyWalletActorPlayer, "player-1"),
		withTransactionCreatedAt(createdAt),
		withTransactionMetadata(sourceMetadata),
	)
	repository := &fakeCurrencyWalletRepository{
		getOwnerResult: activeCurrencyWallet("wallet-1", "player-1"),
		listTransactionsResult: currencymodule.ListCurrencyWalletTransactionsResult{
			Transactions:            []currencymodule.CurrencyWalletTransaction{transaction},
			NextTransactionID:       "tx-2",
			NextTransactionCreateAt: createdAt.Add(time.Minute),
		},
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeCurrencyUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticWalletIDGenerator{value: "wallet-1"}, staticTransactionIDGenerator{value: "tx-ignored"})

	result, err := service.ListOwnWalletTransactions(context.Background(), ListOwnWalletTransactionsRequest{
		Identity:             validatedPlayerIdentity("player-1"),
		CurrencyCode:         " Gems ",
		Limit:                0,
		AfterTransactionID:   " ",
		AfterTransactionTime: time.Time{},
	})
	if err != nil {
		t.Fatalf("ListOwnWalletTransactions() error = %v, want nil", err)
	}
	if result.Status != CurrencyWalletOperationStatusListed ||
		len(result.Transactions) != 1 ||
		result.Transactions[0].TransactionID != "tx-1" ||
		result.NextTransactionID != "tx-2" {
		t.Fatalf("result = %#v, want listed transactions", result)
	}
	if repository.lastListTransactionsInput.WalletID != "wallet-1" ||
		repository.lastListTransactionsInput.CurrencyCode != "Gems" ||
		repository.lastListTransactionsInput.Limit != currencymodule.DefaultListCurrencyWalletTransactionsLimit {
		t.Fatalf("ListCurrencyWalletTransactions input = %#v, want normalized wallet transaction list", repository.lastListTransactionsInput)
	}
	sourceMetadata[0] = '['
	if string(result.Transactions[0].MetadataJSON) != `{"match":"m-1"}` {
		t.Fatalf("transaction metadata = %q, want copied metadata", string(result.Transactions[0].MetadataJSON))
	}
}

func TestGetOwnWalletRequiresUnitOfWorkCurrencyRepositoryCapability(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{unit: tx.NoopUnitOfWork{}}
	service := mustNewService(t, runner, staticWalletIDGenerator{value: "wallet-1"}, staticTransactionIDGenerator{value: "tx-1"})

	result, err := service.GetOwnWallet(context.Background(), GetOwnWalletRequest{
		Identity: validatedPlayerIdentity("player-1"),
	})

	assertServiceError(t, err, OperationGetOwnWallet, FailureClassDependencyUnavailable, PublicErrorCurrencyWalletUnavailable)
	if result.Status != CurrencyWalletOperationStatusRejected {
		t.Fatalf("result = %#v, want rejected", result)
	}
}

func mustNewService(t *testing.T, runner UnitOfWorkRunner, walletGenerator WalletIDGenerator, transactionGenerator TransactionIDGenerator) Service {
	t.Helper()
	service, err := NewService(ServiceDependencies{
		UnitOfWorkRunner:      runner,
		WalletIDGenerator:     walletGenerator,
		TransactionIDGenerator: transactionGenerator,
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

func activeCurrencyWallet(walletID string, ownerID string, options ...func(*currencymodule.CurrencyWallet)) currencymodule.CurrencyWallet {
	createdAt := fixedCurrencyWalletTime()
	record := currencymodule.CurrencyWallet{
		WalletID:       currencymodule.CurrencyWalletID(walletID),
		Owner:          currencymodule.CurrencyWalletOwner{Kind: currencymodule.CurrencyWalletOwnerKindPlayer, ID: ownerID},
		LifecycleState: currencymodule.CurrencyWalletLifecycleActive,
		WalletVersion:  currencymodule.InitialCurrencyWalletVersion,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt.Add(time.Second),
		StateChangedAt: createdAt,
	}
	for _, option := range options {
		option(&record)
	}
	return record
}

func withWalletVersion(version currencymodule.CurrencyWalletVersion) func(*currencymodule.CurrencyWallet) {
	return func(record *currencymodule.CurrencyWallet) {
		record.WalletVersion = version
	}
}

func activeCurrencyWalletBalance(walletID string, currencyCode string, amount currencymodule.CurrencyAmount, version currencymodule.CurrencyBalanceVersion) currencymodule.CurrencyWalletBalance {
	createdAt := fixedCurrencyWalletTime()
	return currencymodule.CurrencyWalletBalance{
		WalletID:       currencymodule.CurrencyWalletID(walletID),
		CurrencyCode:   currencymodule.CurrencyCode(currencyCode),
		BalanceAmount:  amount,
		BalanceVersion: version,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt.Add(time.Second),
	}
}

func activeCurrencyWalletTransaction(transactionID, walletID string, kind currencymodule.CurrencyWalletTransactionKind, currencyCode string, amountDelta currencymodule.CurrencyAmount, balanceAfter currencymodule.CurrencyAmount, options ...func(*currencymodule.CurrencyWalletTransaction)) currencymodule.CurrencyWalletTransaction {
	createdAt := fixedCurrencyWalletTime()
	record := currencymodule.CurrencyWalletTransaction{
		TransactionID:    currencymodule.CurrencyWalletTransactionID(transactionID),
		WalletID:         currencymodule.CurrencyWalletID(walletID),
		CurrencyCode:     currencymodule.CurrencyCode(currencyCode),
		TransactionKind:  kind,
		AmountDelta:      amountDelta,
		BalanceAfter:     balanceAfter,
		IdempotencyKey:   "idem-1",
		IdempotencyScope: "scope-1",
		Actor:            currencymodule.CurrencyWalletActor{Kind: currencymodule.CurrencyWalletActorSystem, ID: "system"},
		MetadataJSON:     []byte(`{"ok":true}`),
		CreatedAt:        createdAt,
	}
	for _, option := range options {
		option(&record)
	}
	return record
}

func withTransactionActor(kind currencymodule.CurrencyWalletActorKind, id string) func(*currencymodule.CurrencyWalletTransaction) {
	return func(record *currencymodule.CurrencyWalletTransaction) {
		record.Actor = currencymodule.CurrencyWalletActor{Kind: kind, ID: id}
	}
}

func withTransactionMetadata(metadata []byte) func(*currencymodule.CurrencyWalletTransaction) {
	return func(record *currencymodule.CurrencyWalletTransaction) {
		record.MetadataJSON = append([]byte(nil), metadata...)
	}
}

func withTransactionCreatedAt(createdAt time.Time) func(*currencymodule.CurrencyWalletTransaction) {
	return func(record *currencymodule.CurrencyWalletTransaction) {
		record.CreatedAt = createdAt
	}
}

func fixedCurrencyWalletTime() time.Time {
	return time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC)
}

func currencyRepositoryError(class currencymodule.CurrencyWalletConflictClass, rawDetail string) error {
	kind := currencymodule.ErrCurrencyWalletConflict
	if class == currencymodule.CurrencyWalletConflictStorageUnavailable {
		kind = currencymodule.ErrCurrencyWalletUnavailable
	}
	return &currencymodule.CurrencyWalletRepositoryError{
		Kind: kind,
		Conflict: currencymodule.CurrencyWalletConflict{
			Class:          class,
			RedactedReason: rawDetail,
		},
		Operation:      "fake",
		RedactedReason: rawDetail,
		Err:            errors.New(rawDetail),
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

type fakeCurrencyUnitOfWork struct {
	ctx             context.Context
	repository      currencymodule.Repository
	repositoryErr   error
	events          *[]string
	repositoryCalls int
}

func (u *fakeCurrencyUnitOfWork) Context() context.Context {
	if u.ctx == nil {
		return context.Background()
	}
	return u.ctx
}

func (u *fakeCurrencyUnitOfWork) NewCurrencyWalletRepository() (currencymodule.Repository, error) {
	u.repositoryCalls += 1
	appendEvent(u.events, "new-currency-wallet-repository")
	if u.repositoryErr != nil {
		return nil, u.repositoryErr
	}
	return u.repository, nil
}

type fakeCurrencyWalletRepository struct {
	events *[]string

	createResult           currencymodule.CurrencyWallet
	createErr              error
	getResult              currencymodule.CurrencyWallet
	getErr                 error
	getOwnerResult         currencymodule.CurrencyWallet
	getOwnerErr            error
	listBalancesResult     currencymodule.ListCurrencyWalletBalancesResult
	listBalancesErr        error
	recordGrantResult      currencymodule.CurrencyWalletTransaction
	recordGrantErr         error
	recordSpendResult      currencymodule.CurrencyWalletTransaction
	recordSpendErr         error
	listTransactionsResult currencymodule.ListCurrencyWalletTransactionsResult
	listTransactionsErr    error

	createCalls           int
	getCalls              int
	getOwnerCalls         int
	listBalancesCalls     int
	recordGrantCalls      int
	recordSpendCalls      int
	listTransactionsCalls int

	lastCreateInput           currencymodule.CreateCurrencyWalletInput
	lastGetInput              currencymodule.GetCurrencyWalletInput
	lastGetOwnerInput         currencymodule.GetCurrencyWalletForOwnerInput
	lastListBalancesInput     currencymodule.ListCurrencyWalletBalancesInput
	lastRecordGrantInput      currencymodule.RecordCurrencyGrantInput
	lastRecordSpendInput      currencymodule.RecordCurrencySpendInput
	lastListTransactionsInput currencymodule.ListCurrencyWalletTransactionsInput
}

func (r *fakeCurrencyWalletRepository) CreateCurrencyWallet(_ context.Context, input currencymodule.CreateCurrencyWalletInput) (currencymodule.CurrencyWallet, error) {
	r.createCalls += 1
	r.lastCreateInput = input
	appendEvent(r.events, "create-currency-wallet")
	if r.createErr != nil {
		return currencymodule.CurrencyWallet{}, r.createErr
	}
	return r.createResult, nil
}

func (r *fakeCurrencyWalletRepository) GetCurrencyWallet(_ context.Context, input currencymodule.GetCurrencyWalletInput) (currencymodule.CurrencyWallet, error) {
	r.getCalls += 1
	r.lastGetInput = input
	appendEvent(r.events, "get-currency-wallet")
	if r.getErr != nil {
		return currencymodule.CurrencyWallet{}, r.getErr
	}
	return r.getResult, nil
}

func (r *fakeCurrencyWalletRepository) GetCurrencyWalletForOwner(_ context.Context, input currencymodule.GetCurrencyWalletForOwnerInput) (currencymodule.CurrencyWallet, error) {
	r.getOwnerCalls += 1
	r.lastGetOwnerInput = input
	appendEvent(r.events, "get-currency-wallet-for-owner")
	if r.getOwnerErr != nil {
		return currencymodule.CurrencyWallet{}, r.getOwnerErr
	}
	return r.getOwnerResult, nil
}

func (r *fakeCurrencyWalletRepository) ListCurrencyWalletBalances(_ context.Context, input currencymodule.ListCurrencyWalletBalancesInput) (currencymodule.ListCurrencyWalletBalancesResult, error) {
	r.listBalancesCalls += 1
	r.lastListBalancesInput = input
	appendEvent(r.events, "list-currency-wallet-balances")
	if r.listBalancesErr != nil {
		return currencymodule.ListCurrencyWalletBalancesResult{}, r.listBalancesErr
	}
	return r.listBalancesResult, nil
}

func (r *fakeCurrencyWalletRepository) RecordCurrencyGrant(_ context.Context, input currencymodule.RecordCurrencyGrantInput) (currencymodule.CurrencyWalletTransaction, error) {
	r.recordGrantCalls += 1
	r.lastRecordGrantInput = input
	appendEvent(r.events, "record-currency-grant")
	if r.recordGrantErr != nil {
		return currencymodule.CurrencyWalletTransaction{}, r.recordGrantErr
	}
	return r.recordGrantResult, nil
}

func (r *fakeCurrencyWalletRepository) RecordCurrencySpend(_ context.Context, input currencymodule.RecordCurrencySpendInput) (currencymodule.CurrencyWalletTransaction, error) {
	r.recordSpendCalls += 1
	r.lastRecordSpendInput = input
	appendEvent(r.events, "record-currency-spend")
	if r.recordSpendErr != nil {
		return currencymodule.CurrencyWalletTransaction{}, r.recordSpendErr
	}
	return r.recordSpendResult, nil
}

func (r *fakeCurrencyWalletRepository) ListCurrencyWalletTransactions(_ context.Context, input currencymodule.ListCurrencyWalletTransactionsInput) (currencymodule.ListCurrencyWalletTransactionsResult, error) {
	r.listTransactionsCalls += 1
	r.lastListTransactionsInput = input
	appendEvent(r.events, "list-currency-wallet-transactions")
	if r.listTransactionsErr != nil {
		return currencymodule.ListCurrencyWalletTransactionsResult{}, r.listTransactionsErr
	}
	return r.listTransactionsResult, nil
}

func (r *fakeCurrencyWalletRepository) totalCalls() int {
	return r.createCalls + r.getCalls + r.getOwnerCalls + r.listBalancesCalls + r.recordGrantCalls + r.recordSpendCalls + r.listTransactionsCalls
}

type staticWalletIDGenerator struct {
	value  string
	events *[]string
}

func (g staticWalletIDGenerator) GenerateCurrencyWalletID(context.Context) (string, error) {
	appendEvent(g.events, "generate-wallet-id")
	return g.value, nil
}

type staticTransactionIDGenerator struct {
	value  string
	events *[]string
}

func (g staticTransactionIDGenerator) GenerateCurrencyWalletTransactionID(context.Context) (string, error) {
	appendEvent(g.events, "generate-transaction-id")
	return g.value, nil
}

func assertServiceError(t *testing.T, err error, operation Operation, class FailureClass, publicCode PublicErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want ServiceError %s", publicCode)
	}
	var serviceErr *ServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("error type = %T, want *ServiceError", err)
	}
	if serviceErr.Operation != operation || serviceErr.Class != class || serviceErr.PublicCode != publicCode {
		t.Fatalf("ServiceError = %#v, want operation=%s class=%s public=%s", serviceErr, operation, class, publicCode)
	}
}

func assertNoLeak(t *testing.T, err error, secret string) {
	t.Helper()
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error %q leaks secret detail %q", err.Error(), secret)
	}
}

func assertEvents(t *testing.T, got []string, want []string) {
	t.Helper()
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("events = %#v, want %#v", got, want)
	}
}

func appendEvent(events *[]string, event string) {
	if events != nil {
		*events = append(*events, event)
	}
}
