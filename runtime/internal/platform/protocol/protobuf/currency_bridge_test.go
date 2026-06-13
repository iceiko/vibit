package protobuf

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appcurrency "github.com/iceiko/vibit/runtime/internal/app/currency"
	currencyv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/currency/v1"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
	"google.golang.org/protobuf/proto"
)

func TestRouteRequestWithCurrencyPayloadMapsWalletRequests(t *testing.T) {
	walletVersion := int64(7)
	balanceVersion := int64(9)
	afterTime := time.Date(2026, 5, 26, 9, 30, 1, 234, time.UTC)

	tests := []struct {
		name    string
		request app.RouteRequest
		assert  func(t *testing.T, payload any)
	}{
		{
			name: "ensure wallet",
			request: app.RouteRequest{
				Route:   appcurrency.EnsurePlayerWalletRoute(),
				Payload: &currencyv1.EnsurePlayerWalletRequest{},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				if _, ok := payload.(appcurrency.EnsurePlayerWalletRequest); !ok {
					t.Fatalf("payload = %T, want EnsurePlayerWalletRequest", payload)
				}
			},
		},
		{
			name: "get own wallet",
			request: app.RouteRequest{
				Route:   appcurrency.GetOwnWalletRoute(),
				Payload: &currencyv1.GetOwnWalletRequest{},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				if _, ok := payload.(appcurrency.GetOwnWalletRequest); !ok {
					t.Fatalf("payload = %T, want GetOwnWalletRequest", payload)
				}
			},
		},
		{
			name: "list balances",
			request: app.RouteRequest{
				Route: appcurrency.ListOwnWalletBalancesRoute(),
				Payload: &currencyv1.ListOwnWalletBalancesRequest{
					Limit:             25,
					AfterCurrencyCode: "GEMS",
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appcurrency.ListOwnWalletBalancesRequest)
				if !ok {
					t.Fatalf("payload = %T, want ListOwnWalletBalancesRequest", payload)
				}
				if got.Limit != 25 || got.AfterCurrencyCode != currencymodule.CurrencyCode("GEMS") {
					t.Fatalf("payload = %#v, want balance pagination request", got)
				}
			},
		},
		{
			name: "grant",
			request: app.RouteRequest{
				Route: appcurrency.GrantCurrencyRoute(),
				Payload: &currencyv1.GrantCurrencyRequest{
					CurrencyCode:           "GEMS",
					Amount:                 50,
					IdempotencyKey:         "grant-key",
					IdempotencyScope:       "local-proof",
					ReasonCode:             "test_grant",
					ExternalReference:      "ref-1",
					MetadataJson:           `{"source":"test"}`,
					ExpectedWalletVersion:  &walletVersion,
					ExpectedBalanceVersion: &balanceVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appcurrency.GrantCurrencyRequest)
				if !ok {
					t.Fatalf("payload = %T, want GrantCurrencyRequest", payload)
				}
				if got.CurrencyCode != currencymodule.CurrencyCode("GEMS") ||
					got.Amount != 50 ||
					got.IdempotencyKey != currencymodule.CurrencyWalletIdempotencyKey("grant-key") ||
					got.IdempotencyScope != currencymodule.CurrencyWalletIdempotencyScope("local-proof") ||
					string(got.MetadataJSON) != `{"source":"test"}` ||
					got.ExpectedWalletVersion == nil ||
					*got.ExpectedWalletVersion != currencymodule.CurrencyWalletVersion(walletVersion) ||
					got.ExpectedBalanceVersion == nil ||
					*got.ExpectedBalanceVersion != currencymodule.CurrencyBalanceVersion(balanceVersion) {
					t.Fatalf("payload = %#v, want mapped grant request", got)
				}
				if got.SystemActorID != "" {
					t.Fatalf("SystemActorID = %q, want bridge to leave server-authoritative actor to handler/service", got.SystemActorID)
				}
			},
		},
		{
			name: "spend",
			request: app.RouteRequest{
				Route: appcurrency.SpendCurrencyRoute(),
				Payload: &currencyv1.SpendCurrencyRequest{
					CurrencyCode:           "GEMS",
					Amount:                 20,
					IdempotencyKey:         "spend-key",
					IdempotencyScope:       "local-proof",
					ReasonCode:             "test_spend",
					ExternalReference:      "ref-2",
					MetadataJson:           `{"sink":"test"}`,
					ExpectedWalletVersion:  &walletVersion,
					ExpectedBalanceVersion: &balanceVersion,
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appcurrency.SpendCurrencyRequest)
				if !ok {
					t.Fatalf("payload = %T, want SpendCurrencyRequest", payload)
				}
				if got.CurrencyCode != currencymodule.CurrencyCode("GEMS") ||
					got.Amount != 20 ||
					got.IdempotencyKey != currencymodule.CurrencyWalletIdempotencyKey("spend-key") ||
					string(got.MetadataJSON) != `{"sink":"test"}` ||
					got.ExpectedWalletVersion == nil ||
					*got.ExpectedWalletVersion != currencymodule.CurrencyWalletVersion(walletVersion) ||
					got.ExpectedBalanceVersion == nil ||
					*got.ExpectedBalanceVersion != currencymodule.CurrencyBalanceVersion(balanceVersion) {
					t.Fatalf("payload = %#v, want mapped spend request", got)
				}
			},
		},
		{
			name: "list transactions",
			request: app.RouteRequest{
				Route: appcurrency.ListOwnWalletTransactionsRoute(),
				Payload: &currencyv1.ListOwnWalletTransactionsRequest{
					CurrencyCode:         "GEMS",
					Limit:                10,
					AfterTransactionId:   "txn-1",
					AfterTransactionTime: afterTime.Format(time.RFC3339Nano),
				},
			},
			assert: func(t *testing.T, payload any) {
				t.Helper()
				got, ok := payload.(appcurrency.ListOwnWalletTransactionsRequest)
				if !ok {
					t.Fatalf("payload = %T, want ListOwnWalletTransactionsRequest", payload)
				}
				if got.CurrencyCode != currencymodule.CurrencyCode("GEMS") ||
					got.Limit != 10 ||
					got.AfterTransactionID != currencymodule.CurrencyWalletTransactionID("txn-1") ||
					!got.AfterTransactionTime.Equal(afterTime) {
					t.Fatalf("payload = %#v, want mapped transaction pagination request", got)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mapped, err := RouteRequestWithDomainPayload(tc.request)
			if err != nil {
				t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
			}
			tc.assert(t, mapped.Payload)
		})
	}
}

func TestRouteRequestWithCurrencyPayloadPreservesMissingExpectedVersions(t *testing.T) {
	mapped, err := RouteRequestWithDomainPayload(app.RouteRequest{
		Route: appcurrency.SpendCurrencyRoute(),
		Payload: &currencyv1.SpendCurrencyRequest{
			CurrencyCode:     "GEMS",
			Amount:           20,
			IdempotencyKey:   "spend-key",
			IdempotencyScope: "local-proof",
		},
	})
	if err != nil {
		t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
	}
	payload, ok := mapped.Payload.(appcurrency.SpendCurrencyRequest)
	if !ok {
		t.Fatalf("Payload = %T, want SpendCurrencyRequest", mapped.Payload)
	}
	if payload.ExpectedWalletVersion != nil || payload.ExpectedBalanceVersion != nil {
		t.Fatalf("expected versions = %#v/%#v, want nil when optional fields are absent", payload.ExpectedWalletVersion, payload.ExpectedBalanceVersion)
	}
}

func TestRouteRequestWithCurrencyPayloadRejectsMalformedTimeAndWrongPayload(t *testing.T) {
	_, err := RouteRequestWithDomainPayload(app.RouteRequest{
		Route: appcurrency.ListOwnWalletTransactionsRoute(),
		Payload: &currencyv1.ListOwnWalletTransactionsRequest{
			AfterTransactionTime: "not-rfc3339",
		},
	})
	if err == nil {
		t.Fatal("RouteRequestWithDomainPayload(malformed time) error = nil, want bridge error")
	}
	assertCurrencyBridgeError(t, err)

	_, err = RouteRequestWithDomainPayload(app.RouteRequest{
		Route:   appcurrency.GetOwnWalletRoute(),
		Payload: &currencyv1.ListOwnWalletBalancesRequest{},
	})
	if err == nil {
		t.Fatal("RouteRequestWithDomainPayload(wrong payload) error = nil, want bridge error")
	}
	assertCurrencyBridgeError(t, err)
}

func TestProtoPayloadFromCurrencyResultMapsResponses(t *testing.T) {
	wallet, balance, grant, spend := currencyBridgeRecords()

	ensurePayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appcurrency.EnsurePlayerWalletRoute(),
		Payload: appcurrency.CurrencyWalletResult{
			Status: appcurrency.CurrencyWalletOperationStatusEnsured,
			Wallet: wallet,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(ensure) error = %v, want nil", err)
	}
	ensure, ok := ensurePayload.(*currencyv1.EnsurePlayerWalletResponse)
	if !ok {
		t.Fatalf("ensure payload = %T, want EnsurePlayerWalletResponse", ensurePayload)
	}
	assertProtoCurrencyWallet(t, ensure.GetWallet(), wallet)
	if ensure.GetStatus() != string(appcurrency.CurrencyWalletOperationStatusEnsured) {
		t.Fatalf("ensure status = %q, want ensured", ensure.GetStatus())
	}

	getPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appcurrency.GetOwnWalletRoute(),
		Payload: &appcurrency.CurrencyWalletResult{
			Status: appcurrency.CurrencyWalletOperationStatusFound,
			Wallet: wallet,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(get) error = %v, want nil", err)
	}
	get, ok := getPayload.(*currencyv1.GetOwnWalletResponse)
	if !ok {
		t.Fatalf("get payload = %T, want GetOwnWalletResponse", getPayload)
	}
	assertProtoCurrencyWallet(t, get.GetWallet(), wallet)

	listBalancesPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appcurrency.ListOwnWalletBalancesRoute(),
		Payload: appcurrency.CurrencyWalletResult{
			Status:           appcurrency.CurrencyWalletOperationStatusListed,
			Balances:         []currencymodule.CurrencyWalletBalance{balance},
			NextCurrencyCode: "GOLD",
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(list balances) error = %v, want nil", err)
	}
	listBalances, ok := listBalancesPayload.(*currencyv1.ListOwnWalletBalancesResponse)
	if !ok {
		t.Fatalf("list balances payload = %T, want ListOwnWalletBalancesResponse", listBalancesPayload)
	}
	if len(listBalances.GetBalances()) != 1 || listBalances.GetNextCurrencyCode() != "GOLD" {
		t.Fatalf("list balances response = %#v, want one balance and next currency", listBalances)
	}
	assertProtoCurrencyBalance(t, listBalances.GetBalances()[0], balance)

	grantPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appcurrency.GrantCurrencyRoute(),
		Payload: appcurrency.CurrencyWalletResult{
			Status:      appcurrency.CurrencyWalletOperationStatusGranted,
			Transaction: grant,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(grant) error = %v, want nil", err)
	}
	grantResponse, ok := grantPayload.(*currencyv1.GrantCurrencyResponse)
	if !ok {
		t.Fatalf("grant payload = %T, want GrantCurrencyResponse", grantPayload)
	}
	assertProtoCurrencyTransaction(t, grantResponse.GetTransaction(), grant)

	spendPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appcurrency.SpendCurrencyRoute(),
		Payload: appcurrency.CurrencyWalletResult{
			Status:      appcurrency.CurrencyWalletOperationStatusSpent,
			Transaction: spend,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(spend) error = %v, want nil", err)
	}
	spendResponse, ok := spendPayload.(*currencyv1.SpendCurrencyResponse)
	if !ok {
		t.Fatalf("spend payload = %T, want SpendCurrencyResponse", spendPayload)
	}
	assertProtoCurrencyTransaction(t, spendResponse.GetTransaction(), spend)

	listTransactionsPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: appcurrency.ListOwnWalletTransactionsRoute(),
		Payload: appcurrency.CurrencyWalletResult{
			Status:                  appcurrency.CurrencyWalletOperationStatusListed,
			Transactions:            []currencymodule.CurrencyWalletTransaction{grant, spend},
			NextTransactionID:       "txn-next",
			NextTransactionCreateAt: spend.CreatedAt.Add(time.Second),
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult(list transactions) error = %v, want nil", err)
	}
	listTransactions, ok := listTransactionsPayload.(*currencyv1.ListOwnWalletTransactionsResponse)
	if !ok {
		t.Fatalf("list transactions payload = %T, want ListOwnWalletTransactionsResponse", listTransactionsPayload)
	}
	if len(listTransactions.GetTransactions()) != 2 ||
		listTransactions.GetNextTransactionId() != "txn-next" ||
		listTransactions.GetNextTransactionTime() != spend.CreatedAt.Add(time.Second).UTC().Format(time.RFC3339Nano) {
		t.Fatalf("list transactions response = %#v, want transactions and next cursor", listTransactions)
	}
	assertProtoCurrencyTransaction(t, listTransactions.GetTransactions()[1], spend)
}

func TestCurrencyProtoPayloadShapeOmitsOwnerWalletActorAndCredentialProofFields(t *testing.T) {
	messages := []proto.Message{
		&currencyv1.CurrencyWallet{},
		&currencyv1.CurrencyWalletBalance{},
		&currencyv1.CurrencyWalletTransaction{},
		&currencyv1.EnsurePlayerWalletRequest{},
		&currencyv1.GetOwnWalletRequest{},
		&currencyv1.ListOwnWalletBalancesRequest{},
		&currencyv1.GrantCurrencyRequest{},
		&currencyv1.SpendCurrencyRequest{},
		&currencyv1.ListOwnWalletTransactionsRequest{},
	}
	for _, message := range messages {
		descriptor := message.ProtoReflect().Descriptor()
		fields := descriptor.Fields()
		fieldNames := map[string]bool{}
		for i := 0; i < fields.Len(); i++ {
			fieldNames[string(fields.Get(i).Name())] = true
		}
		for _, forbidden := range []string{
			"owner_id",
			"player_id",
			"session_id",
			"access_token",
			"credential_proof",
			"lookup_digest",
			"verifier_digest",
			"actor_id",
			"wallet_lookup_id",
			"payment_provider_payload",
			"catalog_id",
			"reward_id",
			"inventory_item_id",
		} {
			if fieldNames[forbidden] {
				t.Fatalf("%s has forbidden field %q", descriptor.FullName(), forbidden)
			}
		}
	}
}

func assertCurrencyBridgeError(t *testing.T, err error) {
	t.Helper()
	var bridgeErr *PayloadBridgeError
	if !errors.As(err, &bridgeErr) {
		t.Fatalf("error = %T %v, want *PayloadBridgeError", err, err)
	}
	for _, leak := range []string{"metadata_json", "access-token", "wallet-1", "player-1", "idempotency"} {
		if leak != "" && strings.Contains(err.Error(), leak) {
			t.Fatalf("bridge error %q leaks private detail %q", err.Error(), leak)
		}
	}
}

func currencyBridgeRecords() (currencymodule.CurrencyWallet, currencymodule.CurrencyWalletBalance, currencymodule.CurrencyWalletTransaction, currencymodule.CurrencyWalletTransaction) {
	createdAt := time.Date(2026, 5, 26, 9, 0, 0, 123, time.FixedZone("test", 8*60*60))
	wallet := currencymodule.CurrencyWallet{
		WalletID:       "wallet-1",
		Owner:          currencymodule.CurrencyWalletOwner{Kind: currencymodule.CurrencyWalletOwnerKindPlayer, ID: "player-1"},
		LifecycleState: currencymodule.CurrencyWalletLifecycleActive,
		WalletVersion:  3,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt.Add(time.Minute),
		StateChangedAt: createdAt.Add(2 * time.Minute),
	}
	balance := currencymodule.CurrencyWalletBalance{
		WalletID:       wallet.WalletID,
		CurrencyCode:   "GEMS",
		BalanceAmount:  100,
		BalanceVersion: 5,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt.Add(time.Minute),
	}
	grant := currencymodule.CurrencyWalletTransaction{
		TransactionID:     "txn-grant-1",
		WalletID:          wallet.WalletID,
		CurrencyCode:      "GEMS",
		TransactionKind:   currencymodule.CurrencyWalletTransactionGrant,
		AmountDelta:       50,
		BalanceAfter:      150,
		IdempotencyKey:    "grant-key",
		IdempotencyScope:  "local-proof",
		Actor:             currencymodule.CurrencyWalletActor{Kind: currencymodule.CurrencyWalletActorSystem, ID: "currency-route"},
		ReasonCode:        "test_grant",
		ExternalReference: "ref-1",
		MetadataJSON:      []byte(`{"source":"test"}`),
		CreatedAt:         createdAt.Add(3 * time.Minute),
	}
	spend := currencymodule.CurrencyWalletTransaction{
		TransactionID:     "txn-spend-1",
		WalletID:          wallet.WalletID,
		CurrencyCode:      "GEMS",
		TransactionKind:   currencymodule.CurrencyWalletTransactionSpend,
		AmountDelta:       -20,
		BalanceAfter:      130,
		IdempotencyKey:    "spend-key",
		IdempotencyScope:  "local-proof",
		Actor:             currencymodule.CurrencyWalletActor{Kind: currencymodule.CurrencyWalletActorPlayer, ID: "player-1"},
		ReasonCode:        "test_spend",
		ExternalReference: "ref-2",
		MetadataJSON:      []byte(`{"sink":"test"}`),
		CreatedAt:         createdAt.Add(4 * time.Minute),
	}
	return wallet, balance, grant, spend
}

func assertProtoCurrencyWallet(t *testing.T, got *currencyv1.CurrencyWallet, want currencymodule.CurrencyWallet) {
	t.Helper()
	if got == nil {
		t.Fatal("CurrencyWallet = nil, want mapped wallet")
	}
	if got.GetWalletId() != string(want.WalletID) ||
		got.GetOwnerKind() != string(want.Owner.Kind) ||
		got.GetLifecycleState() != string(want.LifecycleState) ||
		got.GetWalletVersion() != int64(want.WalletVersion) ||
		got.GetCreatedAt() != want.CreatedAt.UTC().Format(time.RFC3339Nano) ||
		got.GetUpdatedAt() != want.UpdatedAt.UTC().Format(time.RFC3339Nano) ||
		got.GetStateChangedAt() != want.StateChangedAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("CurrencyWallet = %#v, want mapped wallet %#v", got, want)
	}
}

func assertProtoCurrencyBalance(t *testing.T, got *currencyv1.CurrencyWalletBalance, want currencymodule.CurrencyWalletBalance) {
	t.Helper()
	if got == nil {
		t.Fatal("CurrencyWalletBalance = nil, want mapped balance")
	}
	if got.GetCurrencyCode() != string(want.CurrencyCode) ||
		got.GetBalanceAmount() != int64(want.BalanceAmount) ||
		got.GetBalanceVersion() != int64(want.BalanceVersion) ||
		got.GetCreatedAt() != want.CreatedAt.UTC().Format(time.RFC3339Nano) ||
		got.GetUpdatedAt() != want.UpdatedAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("CurrencyWalletBalance = %#v, want mapped balance %#v", got, want)
	}
}

func assertProtoCurrencyTransaction(t *testing.T, got *currencyv1.CurrencyWalletTransaction, want currencymodule.CurrencyWalletTransaction) {
	t.Helper()
	if got == nil {
		t.Fatal("CurrencyWalletTransaction = nil, want mapped transaction")
	}
	if got.GetTransactionId() != string(want.TransactionID) ||
		got.GetCurrencyCode() != string(want.CurrencyCode) ||
		got.GetTransactionKind() != string(want.TransactionKind) ||
		got.GetAmountDelta() != int64(want.AmountDelta) ||
		got.GetBalanceAfter() != int64(want.BalanceAfter) ||
		got.GetActorKind() != string(want.Actor.Kind) ||
		got.GetReasonCode() != want.ReasonCode ||
		got.GetExternalReference() != want.ExternalReference ||
		got.GetMetadataJson() != string(want.MetadataJSON) ||
		got.GetCreatedAt() != want.CreatedAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("CurrencyWalletTransaction = %#v, want mapped transaction %#v", got, want)
	}
}
