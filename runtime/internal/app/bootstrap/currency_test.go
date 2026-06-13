package bootstrap

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appcurrency "github.com/iceiko/vibit/runtime/internal/app/currency"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
)

func TestCurrencyRouteHandlersRegisterWalletRoutesAndPassValidatedIdentity(t *testing.T) {
	dispatcher := app.NewDispatcher()
	service := &recordingCurrencyService{}
	handlers := CurrencyRouteHandlers{
		Service: service,
		GrantPolicy: StaticCurrencyGrantPolicy{
			Allowed:       true,
			SystemActorID: "currency-route-local-proof",
		},
	}

	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}

	identity := app.ValidatedPlayerIdentity("player-1", app.Session{
		ConnectionID:    "connection-1",
		SessionID:       "session-1",
		PlayerID:        "player-1",
		ConnectionEpoch: 7,
	})
	walletVersion := currencymodule.CurrencyWalletVersion(3)
	balanceVersion := currencymodule.CurrencyBalanceVersion(5)
	afterTime := time.Date(2026, 5, 26, 9, 30, 0, 0, time.UTC)

	dispatchCurrencyRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-ensure",
		Route:     appcurrency.EnsurePlayerWalletRoute(),
		Identity:  identity,
		Payload:   appcurrency.EnsurePlayerWalletRequest{},
	})
	if !service.ensureCalled || service.ensureRequest.Identity != identity {
		t.Fatalf("ensure request = %#v, want validated identity", service.ensureRequest)
	}

	dispatchCurrencyRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-get",
		Route:     appcurrency.GetOwnWalletRoute(),
		Identity:  identity,
		Payload:   appcurrency.GetOwnWalletRequest{},
	})
	if !service.getCalled || service.getRequest.Identity != identity {
		t.Fatalf("get request = %#v, want validated identity", service.getRequest)
	}

	dispatchCurrencyRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-balances",
		Route:     appcurrency.ListOwnWalletBalancesRoute(),
		Identity:  identity,
		Payload: appcurrency.ListOwnWalletBalancesRequest{
			Limit:             25,
			AfterCurrencyCode: "GEMS",
		},
	})
	if !service.listBalancesCalled || service.listBalancesRequest.Identity != identity ||
		service.listBalancesRequest.Limit != 25 ||
		service.listBalancesRequest.AfterCurrencyCode != currencymodule.CurrencyCode("GEMS") {
		t.Fatalf("list balances request = %#v, want validated identity and pagination", service.listBalancesRequest)
	}

	dispatchCurrencyRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-grant",
		Route:     appcurrency.GrantCurrencyRoute(),
		Identity:  identity,
		Payload: appcurrency.GrantCurrencyRequest{
			CurrencyCode:           "GEMS",
			Amount:                 50,
			IdempotencyKey:         "grant-key",
			IdempotencyScope:       "local-proof",
			ReasonCode:             "test_grant",
			ExternalReference:      "ref-1",
			MetadataJSON:           []byte(`{"source":"test"}`),
			ExpectedWalletVersion:  &walletVersion,
			ExpectedBalanceVersion: &balanceVersion,
		},
	})
	if !service.grantCalled || service.grantRequest.Identity != identity ||
		service.grantRequest.SystemActorID != "currency-route-local-proof" ||
		service.grantRequest.CurrencyCode != currencymodule.CurrencyCode("GEMS") ||
		service.grantRequest.Amount != 50 ||
		service.grantRequest.ExpectedWalletVersion == nil ||
		*service.grantRequest.ExpectedWalletVersion != walletVersion ||
		service.grantRequest.ExpectedBalanceVersion == nil ||
		*service.grantRequest.ExpectedBalanceVersion != balanceVersion {
		t.Fatalf("grant request = %#v, want validated identity, system actor, amount, and expected versions", service.grantRequest)
	}

	dispatchCurrencyRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-spend",
		Route:     appcurrency.SpendCurrencyRoute(),
		Identity:  identity,
		Payload: appcurrency.SpendCurrencyRequest{
			CurrencyCode:           "GEMS",
			Amount:                 20,
			IdempotencyKey:         "spend-key",
			IdempotencyScope:       "local-proof",
			ReasonCode:             "test_spend",
			ExternalReference:      "ref-2",
			MetadataJSON:           []byte(`{"sink":"test"}`),
			ExpectedWalletVersion:  &walletVersion,
			ExpectedBalanceVersion: &balanceVersion,
		},
	})
	if !service.spendCalled || service.spendRequest.Identity != identity ||
		service.spendRequest.Amount != 20 ||
		service.spendRequest.ExpectedWalletVersion == nil ||
		*service.spendRequest.ExpectedWalletVersion != walletVersion ||
		service.spendRequest.ExpectedBalanceVersion == nil ||
		*service.spendRequest.ExpectedBalanceVersion != balanceVersion {
		t.Fatalf("spend request = %#v, want validated identity and expected versions", service.spendRequest)
	}

	dispatchCurrencyRoute(t, dispatcher, app.RouteRequest{
		RequestID: "request-transactions",
		Route:     appcurrency.ListOwnWalletTransactionsRoute(),
		Identity:  identity,
		Payload: appcurrency.ListOwnWalletTransactionsRequest{
			CurrencyCode:         "GEMS",
			Limit:                10,
			AfterTransactionID:   "txn-1",
			AfterTransactionTime: afterTime,
		},
	})
	if !service.listTransactionsCalled || service.listTransactionsRequest.Identity != identity ||
		service.listTransactionsRequest.CurrencyCode != currencymodule.CurrencyCode("GEMS") ||
		service.listTransactionsRequest.Limit != 10 ||
		service.listTransactionsRequest.AfterTransactionID != currencymodule.CurrencyWalletTransactionID("txn-1") ||
		!service.listTransactionsRequest.AfterTransactionTime.Equal(afterTime) {
		t.Fatalf("list transactions request = %#v, want validated identity and pagination", service.listTransactionsRequest)
	}
}

func TestCurrencyRouteHandlersRejectMalformedPayloadBeforeService(t *testing.T) {
	service := &recordingCurrencyService{}
	handlers := CurrencyRouteHandlers{Service: service}

	result, err := handlers.HandleSpendCurrencyRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appcurrency.SpendCurrencyRoute(),
		Payload:   "not a currency request",
	})
	if err == nil {
		t.Fatal("HandleSpendCurrencyRoute() error = nil, want invalid request error")
	}
	if service.spendCalled {
		t.Fatal("currency service was called for malformed payload")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appcurrency.PublicErrorCurrencyWalletInvalidRequest) {
		t.Fatalf("result error = %#v, want currency invalid request", result.Error)
	}
}

func TestCurrencyRouteHandlersDoNotGrantWithoutServerSideGrantPolicy(t *testing.T) {
	service := &recordingCurrencyService{}
	handlers := CurrencyRouteHandlers{Service: service}

	result, err := handlers.HandleGrantCurrencyRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appcurrency.GrantCurrencyRoute(),
		Identity:  app.ValidatedPlayerIdentity("player-1", app.Session{PlayerID: "player-1"}),
		Payload: appcurrency.GrantCurrencyRequest{
			CurrencyCode:     "GEMS",
			Amount:           50,
			IdempotencyKey:   "grant-key",
			IdempotencyScope: "local-proof",
			MetadataJSON:     []byte(`{"secret":true}`),
		},
	})
	if err == nil {
		t.Fatal("HandleGrantCurrencyRoute() error = nil, want grant policy error")
	}
	if service.grantCalled {
		t.Fatal("currency grant service was called without server-side grant policy")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appcurrency.PublicErrorCurrencyWalletInvalidRequest) {
		t.Fatalf("result error = %#v, want currency invalid request", result.Error)
	}
	errorText := result.Error.Error()
	for _, leak := range []string{"grant-key", "metadata_json", "secret", "player-1", "GEMS"} {
		if strings.Contains(errorText, leak) {
			t.Fatalf("error %q leaked private grant detail %q", errorText, leak)
		}
	}
}

func TestCurrencyRouteHandlersMapServiceErrorsWithoutPrivateDetailLeak(t *testing.T) {
	secretDetail := "internal sql currency_wallet_transactions wallet-1 player-1 idempotency-key access-token metadata_json"
	service := &recordingCurrencyService{
		spendResult: appcurrency.CurrencyWalletResult{
			Status:          appcurrency.CurrencyWalletOperationStatusRejected,
			PublicErrorCode: appcurrency.PublicErrorCurrencyWalletVersionMismatch,
			FailureClass:    appcurrency.FailureClassVersionMismatch,
		},
		spendErr: &appcurrency.ServiceError{
			Operation:  appcurrency.OperationSpendCurrency,
			Class:      appcurrency.FailureClassVersionMismatch,
			PublicCode: appcurrency.PublicErrorCurrencyWalletVersionMismatch,
			Err:        errors.New(secretDetail),
		},
	}
	handlers := CurrencyRouteHandlers{Service: service}

	result, err := handlers.HandleSpendCurrencyRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appcurrency.SpendCurrencyRoute(),
		Payload: appcurrency.SpendCurrencyRequest{
			CurrencyCode:     "GEMS",
			Amount:           20,
			IdempotencyKey:   "idempotency-key",
			IdempotencyScope: "local-proof",
			MetadataJSON:     []byte(`{"secret":true}`),
		},
	})
	if err == nil {
		t.Fatal("HandleSpendCurrencyRoute() error = nil, want public currency error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appcurrency.PublicErrorCurrencyWalletVersionMismatch) {
		t.Fatalf("result error = %#v, want version mismatch code", result.Error)
	}
	errorText := result.Error.Error()
	for _, leak := range []string{secretDetail, "currency_wallet_transactions", "wallet-1", "player-1", "idempotency-key", "access-token", "metadata_json", "secret", "sql"} {
		if strings.Contains(errorText, leak) {
			t.Fatalf("error %q leaked private detail %q", errorText, leak)
		}
	}
	if result.Payload != nil {
		t.Fatalf("result Payload = %#v, want nil on currency failure", result.Payload)
	}
}

func TestCurrencyRouteHandlersRequireService(t *testing.T) {
	handlers := CurrencyRouteHandlers{}

	result, err := handlers.HandleGetOwnWalletRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     appcurrency.GetOwnWalletRoute(),
		Payload:   appcurrency.GetOwnWalletRequest{},
	})
	if err == nil {
		t.Fatal("HandleGetOwnWalletRoute() error = nil, want unavailable error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appcurrency.PublicErrorCurrencyWalletUnavailable) {
		t.Fatalf("result error = %#v, want currency unavailable", result.Error)
	}
}

func dispatchCurrencyRoute(t *testing.T, dispatcher *app.Dispatcher, request app.RouteRequest) app.ApplicationResult {
	t.Helper()
	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch(%s) error = %v, want nil", app.RenderRouteKey(request.Route), err)
	}
	if result.PayloadType != "currency.CurrencyWalletResult" {
		t.Fatalf("PayloadType = %q, want currency.CurrencyWalletResult", result.PayloadType)
	}
	if _, ok := result.Payload.(appcurrency.CurrencyWalletResult); !ok {
		t.Fatalf("Payload = %T, want currency.CurrencyWalletResult", result.Payload)
	}
	return result
}

type recordingCurrencyService struct {
	ensureCalled            bool
	ensureRequest           appcurrency.EnsurePlayerWalletRequest
	ensureResult            appcurrency.CurrencyWalletResult
	ensureErr               error
	getCalled               bool
	getRequest              appcurrency.GetOwnWalletRequest
	getResult               appcurrency.CurrencyWalletResult
	getErr                  error
	listBalancesCalled      bool
	listBalancesRequest     appcurrency.ListOwnWalletBalancesRequest
	listBalancesResult      appcurrency.CurrencyWalletResult
	listBalancesErr         error
	grantCalled             bool
	grantRequest            appcurrency.GrantCurrencyRequest
	grantResult             appcurrency.CurrencyWalletResult
	grantErr                error
	spendCalled             bool
	spendRequest            appcurrency.SpendCurrencyRequest
	spendResult             appcurrency.CurrencyWalletResult
	spendErr                error
	listTransactionsCalled  bool
	listTransactionsRequest appcurrency.ListOwnWalletTransactionsRequest
	listTransactionsResult  appcurrency.CurrencyWalletResult
	listTransactionsErr     error
}

func (s *recordingCurrencyService) EnsurePlayerWallet(_ context.Context, request appcurrency.EnsurePlayerWalletRequest) (appcurrency.CurrencyWalletResult, error) {
	s.ensureCalled = true
	s.ensureRequest = request
	return s.resultOrDefault(s.ensureResult, appcurrency.CurrencyWalletOperationStatusEnsured), s.ensureErr
}

func (s *recordingCurrencyService) GetOwnWallet(_ context.Context, request appcurrency.GetOwnWalletRequest) (appcurrency.CurrencyWalletResult, error) {
	s.getCalled = true
	s.getRequest = request
	return s.resultOrDefault(s.getResult, appcurrency.CurrencyWalletOperationStatusFound), s.getErr
}

func (s *recordingCurrencyService) ListOwnWalletBalances(_ context.Context, request appcurrency.ListOwnWalletBalancesRequest) (appcurrency.CurrencyWalletResult, error) {
	s.listBalancesCalled = true
	s.listBalancesRequest = request
	return s.resultOrDefault(s.listBalancesResult, appcurrency.CurrencyWalletOperationStatusListed), s.listBalancesErr
}

func (s *recordingCurrencyService) GrantCurrency(_ context.Context, request appcurrency.GrantCurrencyRequest) (appcurrency.CurrencyWalletResult, error) {
	s.grantCalled = true
	s.grantRequest = request
	return s.resultOrDefault(s.grantResult, appcurrency.CurrencyWalletOperationStatusGranted), s.grantErr
}

func (s *recordingCurrencyService) SpendCurrency(_ context.Context, request appcurrency.SpendCurrencyRequest) (appcurrency.CurrencyWalletResult, error) {
	s.spendCalled = true
	s.spendRequest = request
	return s.resultOrDefault(s.spendResult, appcurrency.CurrencyWalletOperationStatusSpent), s.spendErr
}

func (s *recordingCurrencyService) ListOwnWalletTransactions(_ context.Context, request appcurrency.ListOwnWalletTransactionsRequest) (appcurrency.CurrencyWalletResult, error) {
	s.listTransactionsCalled = true
	s.listTransactionsRequest = request
	return s.resultOrDefault(s.listTransactionsResult, appcurrency.CurrencyWalletOperationStatusListed), s.listTransactionsErr
}

func (s *recordingCurrencyService) resultOrDefault(result appcurrency.CurrencyWalletResult, status appcurrency.CurrencyWalletOperationStatus) appcurrency.CurrencyWalletResult {
	if result.Status != "" {
		return result
	}
	wallet, balance, transaction := currencyHandlerRecords()
	result = appcurrency.CurrencyWalletResult{
		Status:                  status,
		Wallet:                  wallet,
		Balance:                 balance,
		Balances:                []currencymodule.CurrencyWalletBalance{balance},
		Transaction:             transaction,
		Transactions:            []currencymodule.CurrencyWalletTransaction{transaction},
		NextCurrencyCode:        "GOLD",
		NextTransactionID:       "txn-next",
		NextTransactionCreateAt: transaction.CreatedAt.Add(time.Second),
	}
	return result
}

func currencyHandlerRecords() (currencymodule.CurrencyWallet, currencymodule.CurrencyWalletBalance, currencymodule.CurrencyWalletTransaction) {
	now := time.Date(2026, 5, 26, 9, 0, 0, 0, time.UTC)
	wallet := currencymodule.CurrencyWallet{
		WalletID:       "wallet-1",
		Owner:          currencymodule.CurrencyWalletOwner{Kind: currencymodule.CurrencyWalletOwnerKindPlayer, ID: "player-1"},
		LifecycleState: currencymodule.CurrencyWalletLifecycleActive,
		WalletVersion:  3,
		CreatedAt:      now,
		UpdatedAt:      now.Add(time.Minute),
		StateChangedAt: now,
	}
	balance := currencymodule.CurrencyWalletBalance{
		WalletID:       wallet.WalletID,
		CurrencyCode:   "GEMS",
		BalanceAmount:  100,
		BalanceVersion: 5,
		CreatedAt:      now,
		UpdatedAt:      now.Add(time.Minute),
	}
	transaction := currencymodule.CurrencyWalletTransaction{
		TransactionID:     "txn-1",
		WalletID:          wallet.WalletID,
		CurrencyCode:      "GEMS",
		TransactionKind:   currencymodule.CurrencyWalletTransactionSpend,
		AmountDelta:       -20,
		BalanceAfter:      80,
		IdempotencyKey:    "spend-key",
		IdempotencyScope:  "local-proof",
		Actor:             currencymodule.CurrencyWalletActor{Kind: currencymodule.CurrencyWalletActorPlayer, ID: "player-1"},
		ReasonCode:        "test_spend",
		ExternalReference: "ref-2",
		MetadataJSON:      []byte(`{"sink":"test"}`),
		CreatedAt:         now.Add(2 * time.Minute),
	}
	return wallet, balance, transaction
}
