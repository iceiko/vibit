package protobuf

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
	"github.com/iceiko/vibit/runtime/internal/app/bootstrap"
	appconnection "github.com/iceiko/vibit/runtime/internal/app/connection"
	appcurrency "github.com/iceiko/vibit/runtime/internal/app/currency"
	appfriends "github.com/iceiko/vibit/runtime/internal/app/friends"
	apppresence "github.com/iceiko/vibit/runtime/internal/app/presence"
	sessionmodule "github.com/iceiko/vibit/runtime/internal/app/session"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
	currencyv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/currency/v1"
	friendsv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/friends/v1"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	presencev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/presence/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	storagev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/storage/v1"
	authenticationmodule "github.com/iceiko/vibit/runtime/internal/modules/authentication"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
	friendsmodule "github.com/iceiko/vibit/runtime/internal/modules/friends"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	playermodule "github.com/iceiko/vibit/runtime/internal/modules/player"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
	"google.golang.org/protobuf/proto"
)

func TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout(t *testing.T) {
	ctx := context.Background()
	fixture := newAuthenticatedGameplayE2EFixture(t)

	onboarding, err := fixture.service.OnboardLocalPlayerWithDeviceCredential(ctx, appauth.LocalOnboardingDeviceCredentialIssuanceRequest{
		DisplayName: "Alpha Player",
		RequestedBy: "local-e2e-test",
	})
	if err != nil {
		t.Fatalf("OnboardLocalPlayerWithDeviceCredential() error = %v, want nil", err)
	}
	if onboarding.Status != appauth.LocalOnboardingDeviceCredentialIssuanceStatusCreated ||
		onboarding.PlayerID != "player-e2e-1" ||
		onboarding.CredentialRecordID != "credential-e2e-1" ||
		onboarding.DeviceCredential == "" {
		t.Fatalf("onboarding result = %#v, want local player and one-time device credential", onboarding)
	}

	loginEnvelope := fixture.handleFrame(t, &frameStep{
		route:           app.AuthenticateWithDeviceCredentialRoute(),
		requestID:       "login-1",
		target:          app.Target{Scope: app.TargetScopePlayer, ID: onboarding.PlayerID},
		session:         app.Session{ConnectionID: "ws-e2e-1", ConnectionEpoch: 1},
		connectionID:    "ws-e2e-1",
		connectionEpoch: 1,
		payload: &authenticationv1.AuthenticateWithDeviceCredentialRequest{
			CredentialProof:       onboarding.DeviceCredential,
			RequestedPlayerId:     onboarding.PlayerID,
			ClientInstanceId:      "client-e2e-1",
			AccountCreationIntent: authenticationv1.AccountCreationIntent_ACCOUNT_CREATION_INTENT_AUTHENTICATE_EXISTING_ONLY,
		},
	})
	login := mustDecodePayloadAs[*authenticationv1.AuthenticateWithDeviceCredentialResponse](t, loginEnvelope)
	if login.GetAuthenticationStatus() != string(appauth.AuthenticationStatusAuthenticated) ||
		login.GetPlayerId() != onboarding.PlayerID ||
		login.GetAccessToken() == "" ||
		login.GetTokenRecordId() != "token-e2e-1" ||
		login.GetTokenType() != string(appauth.TokenTypeOpaqueAccess) {
		t.Fatalf("login response = %#v, want authenticated access token for onboarded player", login)
	}
	if loginEnvelope.GetSession().GetPlayerId() != onboarding.PlayerID ||
		loginEnvelope.GetSession().GetSessionId() != "runtime-session-e2e-1" ||
		loginEnvelope.GetSession().GetConnectionId() != "ws-e2e-1" {
		t.Fatalf("login response session = %#v, want login-created runtime session carrier", loginEnvelope.GetSession())
	}

	accessToken := login.GetAccessToken()
	if strings.Contains(login.String(), onboarding.DeviceCredential) {
		t.Fatalf("login response leaks one-time device credential: %v", login)
	}
	if fixture.authenticationRepository.containsRawSecret(onboarding.DeviceCredential, accessToken) {
		t.Fatal("authentication repository stored raw credential or access-token text")
	}

	if _, err := fixture.connectionRegistry.RegisterOpenConnection(ctx, appconnection.OpenConnection{
		ConnectionID:    "ws-e2e-1",
		ConnectionEpoch: 1,
		OpenedAt:        fixture.clock.Now().Add(time.Second),
	}); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	bindEnvelope := fixture.handleFrame(t, &frameStep{
		route:           app.BindConnectionRoute(),
		requestID:       "bind-1",
		target:          app.Target{Scope: app.TargetScopeSystem, ID: "runtime"},
		session:         app.Session{ConnectionID: "ws-e2e-1", SessionID: loginEnvelope.GetSession().GetSessionId(), PlayerID: onboarding.PlayerID, ConnectionEpoch: 1},
		connectionID:    "ws-e2e-1",
		connectionEpoch: 1,
		payload: &authenticationv1.BindConnectionRequest{
			AccessToken:      accessToken,
			ClientInstanceId: "client-e2e-1",
		},
	})
	binding := mustDecodePayloadAs[*authenticationv1.BindConnectionResponse](t, bindEnvelope)
	if binding.GetBindingStatus() != authenticationv1.ConnectionBindingStatus_CONNECTION_BINDING_STATUS_BOUND ||
		binding.GetPlayerId() != onboarding.PlayerID ||
		binding.GetConnectionId() != "ws-e2e-1" ||
		binding.GetConnectionEpoch() != 1 ||
		binding.GetSessionValidated() {
		t.Fatalf("BindConnection response = %#v, want bound player connection without session validation", binding)
	}
	if strings.Contains(binding.String(), accessToken) {
		t.Fatalf("BindConnection response leaks access token: %v", binding)
	}

	grantEnvelope := fixture.handleFrame(t, &frameStep{
		route:           inventory.GrantItemRoute(),
		requestID:       "grant-1",
		target:          app.Target{Scope: app.TargetScopePlayer, ID: onboarding.PlayerID},
		session:         app.Session{ConnectionID: "ws-e2e-1", SessionID: loginEnvelope.GetSession().GetSessionId(), PlayerID: onboarding.PlayerID, ConnectionEpoch: 1},
		connectionID:    "ws-e2e-1",
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     accessToken,
		payload: &inventoryv1.GrantItemRequest{
			PlayerId:    onboarding.PlayerID,
			ItemId:      "item-e2e-sword",
			Quantity:    2,
			Reason:      "e2e-proof",
			RequestedBy: onboarding.PlayerID,
		},
	})
	grant := mustDecodePayloadAs[*inventoryv1.GrantItemResponse](t, grantEnvelope)
	if grant.GetPlayerId() != onboarding.PlayerID ||
		grant.GetItemId() != "item-e2e-sword" ||
		grant.GetQuantity() != 2 ||
		grant.GetNewQuantity() != 2 ||
		grant.GetEvent() != inventory.EventItemGranted {
		t.Fatalf("GrantItem response = %#v, want protected inventory mutation", grant)
	}

	inventoryEnvelope := fixture.handleFrame(t, &frameStep{
		route:           inventory.GetInventoryRoute(),
		requestID:       "inventory-1",
		target:          app.Target{Scope: app.TargetScopePlayer, ID: onboarding.PlayerID},
		session:         app.Session{ConnectionID: "ws-e2e-1", SessionID: loginEnvelope.GetSession().GetSessionId(), PlayerID: onboarding.PlayerID, ConnectionEpoch: 1},
		connectionID:    "ws-e2e-1",
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     accessToken,
		payload: &inventoryv1.GetInventoryRequest{
			PlayerId:    onboarding.PlayerID,
			RequestedBy: onboarding.PlayerID,
		},
	})
	inventoryResponse := mustDecodePayloadAs[*inventoryv1.GetInventoryResponse](t, inventoryEnvelope)
	if inventoryResponse.GetPlayerId() != onboarding.PlayerID ||
		len(inventoryResponse.GetItems()) != 1 ||
		inventoryResponse.GetItems()[0].GetItemId() != "item-e2e-sword" ||
		inventoryResponse.GetItems()[0].GetQuantity() != 2 {
		t.Fatalf("GetInventory response = %#v, want protected inventory read after grant", inventoryResponse)
	}

	presenceEnvelope := fixture.handleFrame(t, &frameStep{
		route:           apppresence.GetPlayerPresenceRoute(),
		requestID:       "presence-1",
		target:          app.Target{Scope: app.TargetScopePlayer, ID: onboarding.PlayerID},
		session:         app.Session{ConnectionID: "ws-e2e-1", SessionID: loginEnvelope.GetSession().GetSessionId(), PlayerID: onboarding.PlayerID, ConnectionEpoch: 1},
		connectionID:    "ws-e2e-1",
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     accessToken,
		payload: &presencev1.GetPlayerPresenceRequest{
			PlayerId: onboarding.PlayerID,
		},
	})
	presence := mustDecodePayloadAs[*presencev1.GetPlayerPresenceResponse](t, presenceEnvelope)
	if presence.GetPlayerId() != onboarding.PlayerID ||
		presence.GetPresenceStatus() != presencev1.PresenceStatus_PRESENCE_STATUS_ONLINE ||
		presence.GetConnectionCount() != 1 ||
		len(presence.GetActiveConnections()) != 1 ||
		presence.GetActiveConnections()[0].GetConnectionId() != "ws-e2e-1" ||
		presence.GetActiveConnections()[0].GetConnectionEpoch() != 1 {
		t.Fatalf("presence response = %#v, want bound authenticated player online", presence)
	}

	logoutEnvelope := fixture.handleFrame(t, &frameStep{
		route:           app.LogoutAccessTokenRoute(),
		requestID:       "logout-1",
		target:          app.Target{Scope: app.TargetScopePlayer, ID: onboarding.PlayerID},
		session:         app.Session{ConnectionID: "ws-e2e-1", SessionID: loginEnvelope.GetSession().GetSessionId(), PlayerID: onboarding.PlayerID, ConnectionEpoch: 1},
		connectionID:    "ws-e2e-1",
		connectionEpoch: 1,
		payload: &authenticationv1.LogoutAccessTokenRequest{
			AccessToken:  accessToken,
			LogoutReason: "e2e_logout",
		},
	})
	logout := mustDecodePayloadAs[*authenticationv1.LogoutAccessTokenResponse](t, logoutEnvelope)
	if logout.GetLogoutStatus() != string(appauth.LogoutStatusRevoked) ||
		!logout.GetRevoked() ||
		logout.GetLogoutScope() != "presented_access_token" ||
		logout.GetTokenRecordId() != "token-e2e-1" {
		t.Fatalf("logout response = %#v, want presented token revocation", logout)
	}
	if strings.Contains(logout.String(), accessToken) {
		t.Fatalf("logout response leaks access token: %v", logout)
	}

	afterLogout := fixture.handleFrame(t, &frameStep{
		route:           inventory.GetInventoryRoute(),
		requestID:       "inventory-after-logout-1",
		target:          app.Target{Scope: app.TargetScopePlayer, ID: onboarding.PlayerID},
		session:         app.Session{ConnectionID: "ws-e2e-1", SessionID: loginEnvelope.GetSession().GetSessionId(), PlayerID: onboarding.PlayerID, ConnectionEpoch: 1},
		connectionID:    "ws-e2e-1",
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     accessToken,
		payload: &inventoryv1.GetInventoryRequest{
			PlayerId:    onboarding.PlayerID,
			RequestedBy: onboarding.PlayerID,
		},
	})
	assertErrorEnvelope(t, afterLogout, app.ErrorCodeAuthenticationTokenInvalid)
	assertNoFrameErrorSecretLeak(t, afterLogout, accessToken, onboarding.DeviceCredential)
}

func TestStorageObjectsProtocolRouteLocalAlphaFlow(t *testing.T) {
	ctx := context.Background()
	fixture := newAuthenticatedGameplayE2EFixture(t)
	player := fixture.authenticateAndBindLocalPlayer(t, ctx, "ws-storage-e2e-1", "client-storage-e2e-1")
	session := app.Session{
		ConnectionID:    player.connectionID,
		SessionID:       player.sessionID,
		PlayerID:        player.playerID,
		ConnectionEpoch: 1,
	}
	target := app.Target{Scope: app.TargetScopePlayer, ID: player.playerID}
	valueJSON := `{"checkpoint":"gate","level":4}`

	putEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appstorage.PutOwnStorageObjectRoute(),
		requestID:       "storage-put-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &storagev1.PutOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
			ValueJson:  valueJSON,
		},
	})
	put := mustDecodePayloadAs[*storagev1.PutOwnStorageObjectResponse](t, putEnvelope)
	if put.GetVersion() != int64(storagemodule.InitialStorageObjectVersion) ||
		put.GetObject().GetCollection() != "progress" ||
		put.GetObject().GetKey() != "tutorial" ||
		put.GetObject().GetValueJson() != valueJSON ||
		put.GetObject().GetVersion() != int64(storagemodule.InitialStorageObjectVersion) {
		t.Fatalf("PutOwnStorageObject response = %#v, want stored progress/tutorial object", put)
	}

	getEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appstorage.GetOwnStorageObjectRoute(),
		requestID:       "storage-get-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &storagev1.GetOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
		},
	})
	get := mustDecodePayloadAs[*storagev1.GetOwnStorageObjectResponse](t, getEnvelope)
	if get.GetObject().GetCollection() != "progress" ||
		get.GetObject().GetKey() != "tutorial" ||
		get.GetObject().GetValueJson() != valueJSON ||
		get.GetObject().GetVersion() != put.GetVersion() {
		t.Fatalf("GetOwnStorageObject response = %#v, want stored object", get)
	}

	listEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appstorage.ListOwnStorageObjectsRoute(),
		requestID:       "storage-list-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &storagev1.ListOwnStorageObjectsRequest{
			Collection: "progress",
			Limit:      10,
		},
	})
	list := mustDecodePayloadAs[*storagev1.ListOwnStorageObjectsResponse](t, listEnvelope)
	if len(list.GetObjects()) != 1 ||
		list.GetObjects()[0].GetCollection() != "progress" ||
		list.GetObjects()[0].GetKey() != "tutorial" ||
		list.GetObjects()[0].GetValueJson() != valueJSON ||
		list.GetNextKey() != "" {
		t.Fatalf("ListOwnStorageObjects response = %#v, want one progress object and no next key", list)
	}

	expectedVersion := put.GetVersion()
	deleteEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appstorage.DeleteOwnStorageObjectRoute(),
		requestID:       "storage-delete-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &storagev1.DeleteOwnStorageObjectRequest{
			Collection:      "progress",
			Key:             "tutorial",
			ExpectedVersion: &expectedVersion,
		},
	})
	deleted := mustDecodePayloadAs[*storagev1.DeleteOwnStorageObjectResponse](t, deleteEnvelope)
	if !deleted.GetDeleted() || deleted.GetVersion() != expectedVersion+1 {
		t.Fatalf("DeleteOwnStorageObject response = %#v, want deleted version %d", deleted, expectedVersion+1)
	}

	afterDelete := fixture.handleFrame(t, &frameStep{
		route:           appstorage.GetOwnStorageObjectRoute(),
		requestID:       "storage-get-after-delete-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &storagev1.GetOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
		},
	})
	assertErrorEnvelope(t, afterDelete, app.ErrorCode(appstorage.PublicErrorStorageObjectNotFound))
	assertNoFrameErrorSecretLeak(t, afterDelete, player.accessToken, player.deviceCredential, valueJSON)
}

func TestCurrencyWalletProtocolRouteLocalAlphaFlow(t *testing.T) {
	ctx := context.Background()
	fixture := newAuthenticatedGameplayE2EFixture(t)
	player := fixture.authenticateAndBindLocalPlayer(t, ctx, "ws-currency-e2e-1", "client-currency-e2e-1")
	session := app.Session{
		ConnectionID:    player.connectionID,
		SessionID:       player.sessionID,
		PlayerID:        player.playerID,
		ConnectionEpoch: 1,
	}
	target := app.Target{Scope: app.TargetScopePlayer, ID: player.playerID}
	metadataJSON := `{"source":"e2e","reward":"daily"}`

	ensureEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appcurrency.EnsurePlayerWalletRoute(),
		requestID:       "currency-ensure-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload:         &currencyv1.EnsurePlayerWalletRequest{},
	})
	ensure := mustDecodePayloadAs[*currencyv1.EnsurePlayerWalletResponse](t, ensureEnvelope)
	if ensure.GetStatus() != string(appcurrency.CurrencyWalletOperationStatusEnsured) ||
		ensure.GetWallet().GetWalletId() != "currency-wallet-e2e-1" ||
		ensure.GetWallet().GetOwnerKind() != string(currencymodule.CurrencyWalletOwnerKindPlayer) ||
		ensure.GetWallet().GetLifecycleState() != string(currencymodule.CurrencyWalletLifecycleActive) ||
		ensure.GetWallet().GetWalletVersion() != int64(currencymodule.InitialCurrencyWalletVersion) {
		t.Fatalf("EnsurePlayerWallet response = %#v, want ensured active player wallet", ensure)
	}

	getEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appcurrency.GetOwnWalletRoute(),
		requestID:       "currency-get-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload:         &currencyv1.GetOwnWalletRequest{},
	})
	get := mustDecodePayloadAs[*currencyv1.GetOwnWalletResponse](t, getEnvelope)
	if get.GetStatus() != string(appcurrency.CurrencyWalletOperationStatusFound) ||
		get.GetWallet().GetWalletId() != ensure.GetWallet().GetWalletId() ||
		get.GetWallet().GetWalletVersion() != ensure.GetWallet().GetWalletVersion() {
		t.Fatalf("GetOwnWallet response = %#v, want ensured wallet", get)
	}

	grantEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appcurrency.GrantCurrencyRoute(),
		requestID:       "currency-grant-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &currencyv1.GrantCurrencyRequest{
			CurrencyCode:      "GEMS",
			Amount:            125,
			IdempotencyKey:    "grant-e2e-1",
			IdempotencyScope:  "local-e2e",
			ReasonCode:        "daily_reward",
			ExternalReference: "reward-e2e-1",
			MetadataJson:      metadataJSON,
		},
	})
	grant := mustDecodePayloadAs[*currencyv1.GrantCurrencyResponse](t, grantEnvelope)
	if grant.GetStatus() != string(appcurrency.CurrencyWalletOperationStatusGranted) ||
		grant.GetTransaction().GetTransactionId() != "currency-transaction-e2e-1" ||
		grant.GetTransaction().GetCurrencyCode() != "GEMS" ||
		grant.GetTransaction().GetTransactionKind() != string(currencymodule.CurrencyWalletTransactionGrant) ||
		grant.GetTransaction().GetAmountDelta() != 125 ||
		grant.GetTransaction().GetBalanceAfter() != 125 ||
		grant.GetTransaction().GetActorKind() != string(currencymodule.CurrencyWalletActorSystem) ||
		grant.GetTransaction().GetMetadataJson() != metadataJSON {
		t.Fatalf("GrantCurrency response = %#v, want server-authorized grant transaction", grant)
	}

	spendEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appcurrency.SpendCurrencyRoute(),
		requestID:       "currency-spend-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &currencyv1.SpendCurrencyRequest{
			CurrencyCode:      "GEMS",
			Amount:            40,
			IdempotencyKey:    "spend-e2e-1",
			IdempotencyScope:  "local-e2e",
			ReasonCode:        "upgrade",
			ExternalReference: "upgrade-e2e-1",
			MetadataJson:      `{"source":"e2e","sink":"upgrade"}`,
		},
	})
	spend := mustDecodePayloadAs[*currencyv1.SpendCurrencyResponse](t, spendEnvelope)
	if spend.GetStatus() != string(appcurrency.CurrencyWalletOperationStatusSpent) ||
		spend.GetTransaction().GetTransactionId() != "currency-transaction-e2e-2" ||
		spend.GetTransaction().GetCurrencyCode() != "GEMS" ||
		spend.GetTransaction().GetTransactionKind() != string(currencymodule.CurrencyWalletTransactionSpend) ||
		spend.GetTransaction().GetAmountDelta() != -40 ||
		spend.GetTransaction().GetBalanceAfter() != 85 ||
		spend.GetTransaction().GetActorKind() != string(currencymodule.CurrencyWalletActorPlayer) {
		t.Fatalf("SpendCurrency response = %#v, want player spend transaction", spend)
	}

	balancesEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appcurrency.ListOwnWalletBalancesRoute(),
		requestID:       "currency-balances-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &currencyv1.ListOwnWalletBalancesRequest{
			Limit: 10,
		},
	})
	balances := mustDecodePayloadAs[*currencyv1.ListOwnWalletBalancesResponse](t, balancesEnvelope)
	if balances.GetStatus() != string(appcurrency.CurrencyWalletOperationStatusListed) ||
		len(balances.GetBalances()) != 1 ||
		balances.GetBalances()[0].GetCurrencyCode() != "GEMS" ||
		balances.GetBalances()[0].GetBalanceAmount() != 85 ||
		balances.GetBalances()[0].GetBalanceVersion() != 2 ||
		balances.GetNextCurrencyCode() != "" {
		t.Fatalf("ListOwnWalletBalances response = %#v, want one updated GEMS balance", balances)
	}

	transactionsEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appcurrency.ListOwnWalletTransactionsRoute(),
		requestID:       "currency-transactions-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &currencyv1.ListOwnWalletTransactionsRequest{
			CurrencyCode: "GEMS",
			Limit:        10,
		},
	})
	transactions := mustDecodePayloadAs[*currencyv1.ListOwnWalletTransactionsResponse](t, transactionsEnvelope)
	if transactions.GetStatus() != string(appcurrency.CurrencyWalletOperationStatusListed) ||
		len(transactions.GetTransactions()) != 2 ||
		transactions.GetTransactions()[0].GetTransactionId() != grant.GetTransaction().GetTransactionId() ||
		transactions.GetTransactions()[1].GetTransactionId() != spend.GetTransaction().GetTransactionId() ||
		transactions.GetNextTransactionId() != "" ||
		transactions.GetNextTransactionTime() != "" {
		t.Fatalf("ListOwnWalletTransactions response = %#v, want grant and spend transactions in order", transactions)
	}

	afterLogout := fixture.handleFrame(t, &frameStep{
		route:           app.LogoutAccessTokenRoute(),
		requestID:       "currency-logout-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		payload: &authenticationv1.LogoutAccessTokenRequest{
			AccessToken:  player.accessToken,
			LogoutReason: "currency_e2e_logout",
		},
	})
	logout := mustDecodePayloadAs[*authenticationv1.LogoutAccessTokenResponse](t, afterLogout)
	if logout.GetLogoutStatus() != string(appauth.LogoutStatusRevoked) || !logout.GetRevoked() {
		t.Fatalf("LogoutAccessToken response = %#v, want revoked currency route token", logout)
	}

	protectedAfterLogout := fixture.handleFrame(t, &frameStep{
		route:           appcurrency.GetOwnWalletRoute(),
		requestID:       "currency-get-after-logout-1",
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload:         &currencyv1.GetOwnWalletRequest{},
	})
	assertErrorEnvelope(t, protectedAfterLogout, app.ErrorCodeAuthenticationTokenInvalid)
	assertNoFrameErrorSecretLeak(t, protectedAfterLogout, player.accessToken, player.deviceCredential, metadataJSON)
}

func TestFriendsRelationshipProtocolRouteLocalAlphaFlow(t *testing.T) {
	ctx := context.Background()
	fixture := newAuthenticatedGameplayE2EFixture(t)
	playerA := fixture.authenticateAndBindLocalPlayer(t, ctx, "ws-friends-a-e2e-1", "client-friends-a-e2e-1")
	playerB := fixture.authenticateAndBindLocalPlayer(t, ctx, "ws-friends-b-e2e-1", "client-friends-b-e2e-1")

	sessionA := app.Session{
		ConnectionID:    playerA.connectionID,
		SessionID:       playerA.sessionID,
		PlayerID:        playerA.playerID,
		ConnectionEpoch: 1,
	}
	targetA := app.Target{Scope: app.TargetScopePlayer, ID: playerA.playerID}
	sessionB := app.Session{
		ConnectionID:    playerB.connectionID,
		SessionID:       playerB.sessionID,
		PlayerID:        playerB.playerID,
		ConnectionEpoch: 1,
	}
	targetB := app.Target{Scope: app.TargetScopePlayer, ID: playerB.playerID}

	sendEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.SendFriendRequestRoute(),
		requestID:       "friends-send-1",
		target:          targetA,
		session:         sessionA,
		connectionID:    playerA.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerA.accessToken,
		payload: &friendsv1.SendFriendRequestRequest{
			TargetPlayerId: playerB.playerID,
		},
	})
	send := mustDecodePayloadAs[*friendsv1.SendFriendRequestResponse](t, sendEnvelope)
	if send.GetStatus() != string(appfriends.FriendRelationshipOperationStatusSent) ||
		send.GetVersion() != int64(friendsmodule.InitialFriendRelationshipVersion) ||
		send.GetRelationship().GetRequestedByPlayerId() != playerA.playerID ||
		send.GetRelationship().GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusOutgoingRequestPending) {
		t.Fatalf("SendFriendRequest response = %#v, want outgoing pending relationship", send)
	}

	statusForBEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.GetFriendRelationshipStatusRoute(),
		requestID:       "friends-status-b-pending-1",
		target:          targetB,
		session:         sessionB,
		connectionID:    playerB.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerB.accessToken,
		payload: &friendsv1.GetFriendRelationshipStatusRequest{
			TargetPlayerId: playerA.playerID,
		},
	})
	statusForB := mustDecodePayloadAs[*friendsv1.GetFriendRelationshipStatusResponse](t, statusForBEnvelope)
	if statusForB.GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusIncomingRequestPending) ||
		statusForB.GetVersion() != send.GetVersion() {
		t.Fatalf("GetFriendRelationshipStatus response = %#v, want incoming pending for target player", statusForB)
	}

	acceptVersion := send.GetVersion()
	acceptEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.AcceptFriendRequestRoute(),
		requestID:       "friends-accept-1",
		target:          targetB,
		session:         sessionB,
		connectionID:    playerB.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerB.accessToken,
		payload: &friendsv1.AcceptFriendRequestRequest{
			TargetPlayerId:  playerA.playerID,
			ExpectedVersion: &acceptVersion,
		},
	})
	accept := mustDecodePayloadAs[*friendsv1.AcceptFriendRequestResponse](t, acceptEnvelope)
	if accept.GetStatus() != string(appfriends.FriendRelationshipOperationStatusAccepted) ||
		accept.GetRelationship().GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusFriends) ||
		accept.GetVersion() != send.GetVersion()+1 {
		t.Fatalf("AcceptFriendRequest response = %#v, want accepted friendship", accept)
	}

	listEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.ListFriendRelationshipsRoute(),
		requestID:       "friends-list-1",
		target:          targetA,
		session:         sessionA,
		connectionID:    playerA.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerA.accessToken,
		payload: &friendsv1.ListFriendRelationshipsRequest{
			Status: string(friendsmodule.FriendRelationshipStatusFriends),
			Limit:  10,
		},
	})
	list := mustDecodePayloadAs[*friendsv1.ListFriendRelationshipsResponse](t, listEnvelope)
	if list.GetPage() == nil ||
		len(list.GetPage().GetRelationships()) != 1 ||
		list.GetPage().GetRelationships()[0].GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusFriends) ||
		list.GetPage().GetNextPairToken() != "" {
		t.Fatalf("ListFriendRelationships response = %#v, want one accepted friendship", list)
	}

	removeVersion := accept.GetVersion()
	removeEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.RemoveFriendRoute(),
		requestID:       "friends-remove-1",
		target:          targetA,
		session:         sessionA,
		connectionID:    playerA.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerA.accessToken,
		payload: &friendsv1.RemoveFriendRequest{
			TargetPlayerId:  playerB.playerID,
			ExpectedVersion: &removeVersion,
		},
	})
	remove := mustDecodePayloadAs[*friendsv1.RemoveFriendResponse](t, removeEnvelope)
	if remove.GetStatus() != string(appfriends.FriendRelationshipOperationStatusRemoved) ||
		remove.GetRelationship().GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusRemoved) ||
		remove.GetVersion() != accept.GetVersion()+1 {
		t.Fatalf("RemoveFriend response = %#v, want removed relationship", remove)
	}

	blockVersion := remove.GetVersion()
	blockEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.BlockPlayerRoute(),
		requestID:       "friends-block-1",
		target:          targetA,
		session:         sessionA,
		connectionID:    playerA.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerA.accessToken,
		payload: &friendsv1.BlockPlayerRequest{
			TargetPlayerId:  playerB.playerID,
			ExpectedVersion: &blockVersion,
		},
	})
	block := mustDecodePayloadAs[*friendsv1.BlockPlayerResponse](t, blockEnvelope)
	if block.GetStatus() != string(appfriends.FriendRelationshipOperationStatusBlocked) ||
		block.GetRelationship().GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusBlockedByActor) ||
		block.GetVersion() != remove.GetVersion()+1 {
		t.Fatalf("BlockPlayer response = %#v, want actor-side blocked relationship", block)
	}

	unblockVersion := block.GetVersion()
	unblockEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.UnblockPlayerRoute(),
		requestID:       "friends-unblock-1",
		target:          targetA,
		session:         sessionA,
		connectionID:    playerA.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerA.accessToken,
		payload: &friendsv1.UnblockPlayerRequest{
			TargetPlayerId:  playerB.playerID,
			ExpectedVersion: &unblockVersion,
		},
	})
	unblock := mustDecodePayloadAs[*friendsv1.UnblockPlayerResponse](t, unblockEnvelope)
	if unblock.GetStatus() != string(appfriends.FriendRelationshipOperationStatusUnblocked) ||
		unblock.GetRelationship().GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusRemoved) ||
		unblock.GetVersion() != block.GetVersion()+1 {
		t.Fatalf("UnblockPlayer response = %#v, want removed relationship after unblock", unblock)
	}

	resendEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.SendFriendRequestRoute(),
		requestID:       "friends-resend-1",
		target:          targetA,
		session:         sessionA,
		connectionID:    playerA.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerA.accessToken,
		payload: &friendsv1.SendFriendRequestRequest{
			TargetPlayerId: playerB.playerID,
		},
	})
	resend := mustDecodePayloadAs[*friendsv1.SendFriendRequestResponse](t, resendEnvelope)
	if resend.GetStatus() != string(appfriends.FriendRelationshipOperationStatusSent) ||
		resend.GetRelationship().GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusOutgoingRequestPending) ||
		resend.GetVersion() != unblock.GetVersion()+1 {
		t.Fatalf("second SendFriendRequest response = %#v, want new pending request after unblock", resend)
	}

	rejectVersion := resend.GetVersion()
	rejectEnvelope := fixture.handleFrame(t, &frameStep{
		route:           appfriends.RejectFriendRequestRoute(),
		requestID:       "friends-reject-1",
		target:          targetB,
		session:         sessionB,
		connectionID:    playerB.connectionID,
		connectionEpoch: 1,
		authenticated:   true,
		accessToken:     playerB.accessToken,
		payload: &friendsv1.RejectFriendRequestRequest{
			TargetPlayerId:  playerA.playerID,
			ExpectedVersion: &rejectVersion,
		},
	})
	reject := mustDecodePayloadAs[*friendsv1.RejectFriendRequestResponse](t, rejectEnvelope)
	if reject.GetStatus() != string(appfriends.FriendRelationshipOperationStatusRequestRejected) ||
		reject.GetRelationship().GetPublicStatus() != string(appfriends.FriendRelationshipPublicStatusRejected) ||
		reject.GetVersion() != resend.GetVersion()+1 {
		t.Fatalf("RejectFriendRequest response = %#v, want rejected relationship", reject)
	}

	assertNoFrameErrorSecretLeak(t, rejectEnvelope, playerA.accessToken, playerA.deviceCredential, playerB.accessToken, playerB.deviceCredential)
}

func TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation(t *testing.T) {
	ctx := context.Background()

	t.Run("transport close", func(t *testing.T) {
		fixture := newAuthenticatedGameplayE2EFixture(t)
		player := fixture.authenticateAndBindLocalPlayer(t, ctx, "ws-presence-close-e2e-1", "client-presence-close-e2e-1")
		session := app.Session{
			ConnectionID:    player.connectionID,
			SessionID:       player.sessionID,
			PlayerID:        player.playerID,
			ConnectionEpoch: 1,
		}
		target := app.Target{Scope: app.TargetScopePlayer, ID: player.playerID}

		online := fixture.queryPlayerPresence(t, "presence-close-online-1", target, session, player, 1)
		assertPresenceOnline(t, online, player.playerID, player.connectionID, 1)

		if _, err := fixture.connectionRegistry.MarkConnectionClosed(ctx, appconnection.MarkClosed{
			ConnectionID:     appconnection.ConnectionID(player.connectionID),
			ConnectionEpoch:  1,
			ClosedAt:         fixture.clock.Now().Add(2 * time.Second),
			CloseReasonClass: "transport_closed",
		}); err != nil {
			t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
		}

		offline := fixture.queryPlayerPresence(t, "presence-close-offline-1", target, session, player, 1)
		assertPresenceOffline(t, offline, player.playerID)
		assertNoPresenceSecretLeak(t, offline, player.accessToken, player.deviceCredential)
	})

	t.Run("policy invalidation", func(t *testing.T) {
		fixture := newAuthenticatedGameplayE2EFixture(t)
		player := fixture.authenticateAndBindLocalPlayer(t, ctx, "ws-presence-invalidated-e2e-1", "client-presence-invalidated-e2e-1")
		session := app.Session{
			ConnectionID:    player.connectionID,
			SessionID:       player.sessionID,
			PlayerID:        player.playerID,
			ConnectionEpoch: 1,
		}
		target := app.Target{Scope: app.TargetScopePlayer, ID: player.playerID}

		online := fixture.queryPlayerPresence(t, "presence-invalidated-online-1", target, session, player, 1)
		assertPresenceOnline(t, online, player.playerID, player.connectionID, 1)

		closeResult, err := appconnection.NewClosePolicy(
			fixture.connectionRegistry,
			appconnection.WithClosePolicyClock(fixture.clock),
		).RequestClose(ctx, appconnection.CloseConnectionsCommand{
			Target:           appconnection.TargetConnection(appconnection.ConnectionID(player.connectionID), 1),
			ReasonClass:      appconnection.CloseReasonTokenRevoked,
			PublicVisibility: appconnection.ClosePublicVisibilityGenericReauthRequired,
			Retryability:     appconnection.CloseRetryabilityNotRetryable,
		})
		if err != nil {
			t.Fatalf("RequestClose() error = %v, want nil", err)
		}
		if len(closeResult.Intents) != 1 || closeResult.Intents[0].Outcome != appconnection.CloseOutcomeInvalidated {
			t.Fatalf("closeResult = %#v, want one invalidated close intent", closeResult)
		}

		offline := fixture.queryPlayerPresence(t, "presence-invalidated-offline-1", target, session, player, 1)
		assertPresenceOffline(t, offline, player.playerID)
		assertNoPresenceSecretLeak(t, offline, player.accessToken, player.deviceCredential)
	})
}

func TestAuthenticatedGameplayFailurePathsLocalAlphaFlow(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		handle       func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope
		wantError    app.ErrorCode
		extraNoLeaks []string
		expireToken  bool
		logoutBefore bool
	}{
		{
			name:      "protected inventory without authenticated wrapper",
			wantError: app.ErrorCodeAuthenticationTokenMissing,
			handle: func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope {
				t.Helper()
				return fixture.handleFrame(t, &frameStep{
					route:           inventory.GetInventoryRoute(),
					requestID:       "failure-missing-wrapper-1",
					target:          target,
					session:         session,
					connectionID:    player.connectionID,
					connectionEpoch: 1,
					payload: &inventoryv1.GetInventoryRequest{
						PlayerId:    player.playerID,
						RequestedBy: player.playerID,
					},
				})
			},
		},
		{
			name:      "protected inventory with malformed authenticated wrapper",
			wantError: app.ErrorCodeAuthenticationTokenMalformed,
			handle: func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope {
				t.Helper()
				return fixture.handleMalformedAuthenticatedInventoryFrame(t, "failure-malformed-wrapper-1", target, session, player)
			},
		},
		{
			name:      "protected inventory with malformed access token text",
			wantError: app.ErrorCodeAuthenticationTokenMalformed,
			handle: func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope {
				t.Helper()
				return fixture.handleFrame(t, &frameStep{
					route:           inventory.GetInventoryRoute(),
					requestID:       "failure-malformed-token-1",
					target:          target,
					session:         session,
					connectionID:    player.connectionID,
					connectionEpoch: 1,
					authenticated:   true,
					accessToken:     "not-a-valid-access-token",
					payload: &inventoryv1.GetInventoryRequest{
						PlayerId:    player.playerID,
						RequestedBy: player.playerID,
					},
				})
			},
			extraNoLeaks: []string{"not-a-valid-access-token"},
		},
		{
			name:      "protected inventory with unknown well-formed access token",
			wantError: app.ErrorCodeAuthenticationTokenInvalid,
			handle: func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope {
				t.Helper()
				unknownAccessToken := base64.RawURLEncoding.EncodeToString(bytesWithIncrementingSeed(120, appauth.RawSecretMaterialBytes))
				return fixture.handleFrame(t, &frameStep{
					route:           inventory.GetInventoryRoute(),
					requestID:       "failure-unknown-token-1",
					target:          target,
					session:         session,
					connectionID:    player.connectionID,
					connectionEpoch: 1,
					authenticated:   true,
					accessToken:     unknownAccessToken,
					payload: &inventoryv1.GetInventoryRequest{
						PlayerId:    player.playerID,
						RequestedBy: player.playerID,
					},
				})
			},
			extraNoLeaks: []string{base64.RawURLEncoding.EncodeToString(bytesWithIncrementingSeed(120, appauth.RawSecretMaterialBytes))},
		},
		{
			name:        "protected inventory with expired access token",
			wantError:   app.ErrorCodeAuthenticationTokenInvalid,
			expireToken: true,
			handle: func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope {
				t.Helper()
				return fixture.handleFrame(t, &frameStep{
					route:           inventory.GetInventoryRoute(),
					requestID:       "failure-expired-token-1",
					target:          target,
					session:         session,
					connectionID:    player.connectionID,
					connectionEpoch: 1,
					authenticated:   true,
					accessToken:     player.accessToken,
					payload: &inventoryv1.GetInventoryRequest{
						PlayerId:    player.playerID,
						RequestedBy: player.playerID,
					},
				})
			},
		},
		{
			name:         "protected inventory with revoked access token after logout",
			wantError:    app.ErrorCodeAuthenticationTokenInvalid,
			logoutBefore: true,
			handle: func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope {
				t.Helper()
				return fixture.handleFrame(t, &frameStep{
					route:           inventory.GetInventoryRoute(),
					requestID:       "failure-revoked-token-1",
					target:          target,
					session:         session,
					connectionID:    player.connectionID,
					connectionEpoch: 1,
					authenticated:   true,
					accessToken:     player.accessToken,
					payload: &inventoryv1.GetInventoryRequest{
						PlayerId:    player.playerID,
						RequestedBy: player.playerID,
					},
				})
			},
		},
		{
			name:      "protected presence with missing authenticated wrapper",
			wantError: app.ErrorCodeAuthenticationTokenMissing,
			handle: func(t *testing.T, fixture authenticatedGameplayE2EFixture, player e2eAuthenticatedPlayer, session app.Session, target app.Target) *protocolv1.Envelope {
				t.Helper()
				return fixture.handleFrame(t, &frameStep{
					route:           apppresence.GetPlayerPresenceRoute(),
					requestID:       "failure-presence-missing-wrapper-1",
					target:          target,
					session:         session,
					connectionID:    player.connectionID,
					connectionEpoch: 1,
					payload: &presencev1.GetPlayerPresenceRequest{
						PlayerId: player.playerID,
					},
				})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newAuthenticatedGameplayE2EFixture(t)
			player := fixture.authenticateAndBindLocalPlayer(t, ctx, "ws-failure-e2e-1", "client-failure-e2e-1")
			session := app.Session{
				ConnectionID:    player.connectionID,
				SessionID:       player.sessionID,
				PlayerID:        player.playerID,
				ConnectionEpoch: 1,
			}
			target := app.Target{Scope: app.TargetScopePlayer, ID: player.playerID}

			if tc.expireToken {
				fixture.authenticationRepository.expireTokenForE2E(t, "token-e2e-1", fixture.clock.Now().Add(-time.Second))
			}
			if tc.logoutBefore {
				logoutEnvelope := fixture.handleFrame(t, &frameStep{
					route:           app.LogoutAccessTokenRoute(),
					requestID:       "failure-logout-before-protected-request-1",
					target:          target,
					session:         session,
					connectionID:    player.connectionID,
					connectionEpoch: 1,
					payload: &authenticationv1.LogoutAccessTokenRequest{
						AccessToken:  player.accessToken,
						LogoutReason: "failure_path_e2e_logout",
					},
				})
				logout := mustDecodePayloadAs[*authenticationv1.LogoutAccessTokenResponse](t, logoutEnvelope)
				if logout.GetLogoutStatus() != string(appauth.LogoutStatusRevoked) || !logout.GetRevoked() {
					t.Fatalf("logout response = %#v, want revoked before protected failure request", logout)
				}
				assertNoFrameErrorSecretLeak(t, logoutEnvelope, player.accessToken, player.deviceCredential)
			}

			response := tc.handle(t, fixture, player, session, target)
			assertErrorEnvelope(t, response, tc.wantError)
			noLeaks := append([]string{player.accessToken, player.deviceCredential}, tc.extraNoLeaks...)
			assertNoFrameErrorSecretLeak(t, response, noLeaks...)
		})
	}
}

type authenticatedGameplayE2EFixture struct {
	clock                    e2eClock
	service                  appauth.Service
	handler                  FrameHandler
	connectionRegistry       *appconnection.InMemoryRegistry
	authenticationRepository *e2eAuthenticationRepository
}

func newAuthenticatedGameplayE2EFixture(t *testing.T) authenticatedGameplayE2EFixture {
	t.Helper()

	clock := e2eClock{now: time.Date(2026, 5, 21, 9, 0, 0, 0, time.UTC)}
	keySet := mustE2EVerifierKeySet(t)
	authenticationRepository := newE2EAuthenticationRepository()
	playerRepository := newE2EPlayerRepository()
	sessionRepository := newE2ESessionRepository()
	storageRepository := newE2EStorageRepository(clock)
	friendsRepository := newE2EFriendRelationshipRepository(clock)
	currencyRepository := newE2ECurrencyWalletRepository(clock)
	runner := e2eUnitOfWorkRunner{unit: &e2eUnitOfWork{
		authenticationRepository: authenticationRepository,
		playerRepository:         playerRepository,
		sessionRepository:        sessionRepository,
		storageRepository:        storageRepository,
		friendsRepository:        friendsRepository,
		currencyRepository:       currencyRepository,
	}}
	service, err := appauth.NewService(appauth.ServiceDependencies{
		UnitOfWorkRunner:              runner,
		VerifierKeySet:                keySet,
		AccessTokenRandom:             bytes.NewReader(bytesWithIncrementingSeed(70, appauth.RawSecretMaterialBytes*16)),
		DeviceCredentialRandom:        bytes.NewReader(bytesWithIncrementingSeed(20, appauth.RawSecretMaterialBytes*16)),
		Clock:                         clock,
		TokenRecordIDGenerator:        &e2eTokenRecordIDGenerator{},
		SessionIDGenerator:            &e2eSessionIDGenerator{},
		PlayerIDGenerator:             &e2ePlayerIDGenerator{},
		PlayerAccountEventIDGenerator: &e2ePlayerAccountEventIDGenerator{},
		CredentialRecordIDGenerator:   &e2eCredentialRecordIDGenerator{},
		AccessTokenLifetime:           time.Hour,
		TokenAudience:                 "gameplay-e2e",
	})
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}

	inventoryRepository := inventory.NewMemoryRepository()
	dispatcher, err := bootstrap.NewInventoryDispatcher(bootstrap.InventoryOptions{
		Repositories:     bootstrap.StaticInventoryRepositoryProvider{Repository: inventoryRepository},
		PermissionPolicy: inventory.StaticPermissionPolicy{GrantAllowed: true, ReadAllowed: true},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 16},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "e2e-inventory-event"},
		Clock:            clock,
	})
	if err != nil {
		t.Fatalf("NewInventoryDispatcher() error = %v, want nil", err)
	}
	if err := (bootstrap.AuthenticationRouteHandlers{Service: service}).RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("Register authentication routes error = %v, want nil", err)
	}

	connectionRegistry := appconnection.NewInMemoryRegistry(clock)
	if err := (apppresence.Handlers{Registry: connectionRegistry}).RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("Register presence routes error = %v, want nil", err)
	}

	storageService, err := appstorage.NewService(appstorage.ServiceDependencies{
		UnitOfWorkRunner:  runner,
		ObjectIDGenerator: e2eStorageObjectIDGenerator{},
	})
	if err != nil {
		t.Fatalf("storage NewService() error = %v, want nil", err)
	}
	if err := (bootstrap.StorageRouteHandlers{Service: storageService}).RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("Register storage routes error = %v, want nil", err)
	}
	friendsService, err := appfriends.NewService(appfriends.ServiceDependencies{
		UnitOfWorkRunner:        runner,
		RelationshipIDGenerator: &e2eFriendRelationshipIDGenerator{},
	})
	if err != nil {
		t.Fatalf("friends NewService() error = %v, want nil", err)
	}
	if err := (bootstrap.FriendsRouteHandlers{Service: friendsService}).RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("Register friends routes error = %v, want nil", err)
	}
	currencyService, err := appcurrency.NewService(appcurrency.ServiceDependencies{
		UnitOfWorkRunner:       runner,
		WalletIDGenerator:      &e2eCurrencyWalletIDGenerator{},
		TransactionIDGenerator: &e2eCurrencyTransactionIDGenerator{},
	})
	if err != nil {
		t.Fatalf("currency NewService() error = %v, want nil", err)
	}
	if err := (bootstrap.CurrencyRouteHandlers{
		Service: currencyService,
		GrantPolicy: bootstrap.StaticCurrencyGrantPolicy{
			Allowed:       true,
			SystemActorID: "currency-route-local-proof",
		},
	}).RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("Register currency routes error = %v, want nil", err)
	}

	routeValidator := appauth.NewRouteAccessTokenValidator(service)
	handler := FrameHandler{
		Dispatcher: app.SessionValidatingDispatcher{
			Dispatcher: dispatcher,
			Validator:  app.MetadataOnlySessionValidator{},
		},
		RouteProtector: app.NewRouteProtector(routeValidator),
		ConnectionBinder: e2eRegistryConnectionBinder{
			binder:   app.NewConnectionBinder(routeValidator),
			registry: connectionRegistry,
		},
	}

	return authenticatedGameplayE2EFixture{
		clock:                    clock,
		service:                  service,
		handler:                  handler,
		connectionRegistry:       connectionRegistry,
		authenticationRepository: authenticationRepository,
	}
}

func (f authenticatedGameplayE2EFixture) handleFrame(t *testing.T, step *frameStep) *protocolv1.Envelope {
	t.Helper()
	if step == nil {
		t.Fatal("frame step is nil")
	}

	payload := step.payload
	if step.authenticated {
		innerBytes, err := proto.Marshal(step.payload)
		if err != nil {
			t.Fatalf("proto.Marshal(inner payload) error = %v, want nil", err)
		}
		payload = &authenticationv1.AuthenticatedRequest{
			AccessToken:      step.accessToken,
			InnerPayloadType: PayloadType(step.payload),
			InnerPayload:     innerBytes,
		}
	}

	envelope := mustBuildEnvelope(t, step.route, step.requestID, step.target, step.session, payload)
	encoded, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal(envelope) error = %v, want nil", err)
	}

	responses, err := f.handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID:    step.connectionID,
		ConnectionEpoch: step.connectionEpoch,
		Payload:         encoded,
	})
	if err != nil {
		t.Fatalf("HandleFrame(%s) error = %v, want nil", step.requestID, err)
	}
	return mustUnmarshalSingleResponse(t, responses)
}

func (f authenticatedGameplayE2EFixture) handleMalformedAuthenticatedInventoryFrame(t *testing.T, requestID string, target app.Target, session app.Session, player e2eAuthenticatedPlayer) *protocolv1.Envelope {
	t.Helper()

	envelope := mustBuildEnvelope(
		t,
		inventory.GetInventoryRoute(),
		requestID,
		target,
		session,
		&authenticationv1.AuthenticatedRequest{
			AccessToken:      player.accessToken,
			InnerPayloadType: PayloadType(&inventoryv1.GetInventoryRequest{}),
		},
	)
	encoded, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal(malformed authenticated envelope) error = %v, want nil", err)
	}
	responses, err := f.handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID:    player.connectionID,
		ConnectionEpoch: session.ConnectionEpoch,
		Payload:         encoded,
	})
	if err != nil {
		t.Fatalf("HandleFrame(%s) error = %v, want nil", requestID, err)
	}
	return mustUnmarshalSingleResponse(t, responses)
}

func (f authenticatedGameplayE2EFixture) queryPlayerPresence(t *testing.T, requestID string, target app.Target, session app.Session, player e2eAuthenticatedPlayer, epoch uint64) *presencev1.GetPlayerPresenceResponse {
	t.Helper()

	envelope := f.handleFrame(t, &frameStep{
		route:           apppresence.GetPlayerPresenceRoute(),
		requestID:       requestID,
		target:          target,
		session:         session,
		connectionID:    player.connectionID,
		connectionEpoch: epoch,
		authenticated:   true,
		accessToken:     player.accessToken,
		payload: &presencev1.GetPlayerPresenceRequest{
			PlayerId: player.playerID,
		},
	})
	return mustDecodePayloadAs[*presencev1.GetPlayerPresenceResponse](t, envelope)
}

type frameStep struct {
	route           app.RouteKey
	requestID       string
	target          app.Target
	session         app.Session
	connectionID    string
	connectionEpoch uint64
	authenticated   bool
	accessToken     string
	payload         proto.Message
}

func mustDecodePayloadAs[T proto.Message](t *testing.T, envelope *protocolv1.Envelope) T {
	t.Helper()
	if envelope.GetKind() == protocolv1.MessageKind_MESSAGE_KIND_ERROR {
		t.Fatalf("response is error envelope: %#v", envelope.GetError())
	}
	payload, err := DecodePayload(envelope.GetPayloadType(), envelope.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(%q) error = %v, want nil", envelope.GetPayloadType(), err)
	}
	typed, ok := payload.(T)
	if !ok {
		t.Fatalf("payload = %T, want requested response type", payload)
	}
	return typed
}

type e2eAuthenticatedPlayer struct {
	playerID         string
	sessionID        string
	connectionID     string
	accessToken      string
	deviceCredential string
}

func (f authenticatedGameplayE2EFixture) authenticateAndBindLocalPlayer(t *testing.T, ctx context.Context, connectionID string, clientInstanceID string) e2eAuthenticatedPlayer {
	t.Helper()

	onboarding, err := f.service.OnboardLocalPlayerWithDeviceCredential(ctx, appauth.LocalOnboardingDeviceCredentialIssuanceRequest{
		DisplayName: "Storage Alpha Player",
		RequestedBy: "local-storage-e2e-test",
	})
	if err != nil {
		t.Fatalf("OnboardLocalPlayerWithDeviceCredential() error = %v, want nil", err)
	}

	loginEnvelope := f.handleFrame(t, &frameStep{
		route:           app.AuthenticateWithDeviceCredentialRoute(),
		requestID:       "storage-login-1",
		target:          app.Target{Scope: app.TargetScopePlayer, ID: onboarding.PlayerID},
		session:         app.Session{ConnectionID: connectionID, ConnectionEpoch: 1},
		connectionID:    connectionID,
		connectionEpoch: 1,
		payload: &authenticationv1.AuthenticateWithDeviceCredentialRequest{
			CredentialProof:       onboarding.DeviceCredential,
			RequestedPlayerId:     onboarding.PlayerID,
			ClientInstanceId:      clientInstanceID,
			AccountCreationIntent: authenticationv1.AccountCreationIntent_ACCOUNT_CREATION_INTENT_AUTHENTICATE_EXISTING_ONLY,
		},
	})
	login := mustDecodePayloadAs[*authenticationv1.AuthenticateWithDeviceCredentialResponse](t, loginEnvelope)
	if login.GetAuthenticationStatus() != string(appauth.AuthenticationStatusAuthenticated) ||
		login.GetPlayerId() != onboarding.PlayerID ||
		login.GetAccessToken() == "" {
		t.Fatalf("storage login response = %#v, want authenticated player", login)
	}

	if _, err := f.connectionRegistry.RegisterOpenConnection(ctx, appconnection.OpenConnection{
		ConnectionID:    appconnection.ConnectionID(connectionID),
		ConnectionEpoch: 1,
		OpenedAt:        f.clock.Now().Add(time.Second),
	}); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	bindEnvelope := f.handleFrame(t, &frameStep{
		route:           app.BindConnectionRoute(),
		requestID:       "storage-bind-1",
		target:          app.Target{Scope: app.TargetScopeSystem, ID: "runtime"},
		session:         app.Session{ConnectionID: connectionID, SessionID: loginEnvelope.GetSession().GetSessionId(), PlayerID: onboarding.PlayerID, ConnectionEpoch: 1},
		connectionID:    connectionID,
		connectionEpoch: 1,
		payload: &authenticationv1.BindConnectionRequest{
			AccessToken:      login.GetAccessToken(),
			ClientInstanceId: clientInstanceID,
		},
	})
	binding := mustDecodePayloadAs[*authenticationv1.BindConnectionResponse](t, bindEnvelope)
	if binding.GetBindingStatus() != authenticationv1.ConnectionBindingStatus_CONNECTION_BINDING_STATUS_BOUND ||
		binding.GetPlayerId() != onboarding.PlayerID ||
		binding.GetConnectionId() != connectionID ||
		binding.GetConnectionEpoch() != 1 {
		t.Fatalf("storage BindConnection response = %#v, want bound storage proof connection", binding)
	}

	return e2eAuthenticatedPlayer{
		playerID:         onboarding.PlayerID,
		sessionID:        loginEnvelope.GetSession().GetSessionId(),
		connectionID:     connectionID,
		accessToken:      login.GetAccessToken(),
		deviceCredential: onboarding.DeviceCredential,
	}
}

func assertPresenceOnline(t *testing.T, response *presencev1.GetPlayerPresenceResponse, playerID string, connectionID string, epoch uint64) {
	t.Helper()

	if response.GetPlayerId() != playerID ||
		response.GetPresenceStatus() != presencev1.PresenceStatus_PRESENCE_STATUS_ONLINE ||
		response.GetConnectionCount() != 1 ||
		len(response.GetActiveConnections()) != 1 ||
		response.GetLastSeenAt() == "" {
		t.Fatalf("presence response = %#v, want online self-presence", response)
	}
	active := response.GetActiveConnections()[0]
	if active.GetConnectionId() != connectionID ||
		active.GetConnectionEpoch() != epoch ||
		active.GetLastSeenAt() == "" ||
		active.GetBoundAt() == "" {
		t.Fatalf("active presence connection = %#v, want bound connection metadata", active)
	}
}

func assertPresenceOffline(t *testing.T, response *presencev1.GetPlayerPresenceResponse, playerID string) {
	t.Helper()

	if response.GetPlayerId() != playerID ||
		response.GetPresenceStatus() != presencev1.PresenceStatus_PRESENCE_STATUS_OFFLINE ||
		response.GetConnectionCount() != 0 ||
		len(response.GetActiveConnections()) != 0 ||
		len(response.GetRuntimeSessionIds()) != 0 ||
		response.GetLastSeenAt() != "" ||
		response.GetObservedAt() == "" {
		t.Fatalf("presence response = %#v, want offline without active connection metadata", response)
	}
}

func assertNoPresenceSecretLeak(t *testing.T, response *presencev1.GetPlayerPresenceResponse, secrets ...string) {
	t.Helper()

	text := response.String()
	for _, secret := range secrets {
		if secret != "" && strings.Contains(text, secret) {
			t.Fatalf("presence response leaks secret %q: %v", secret, response)
		}
	}
}

func mustE2EVerifierKeySet(t *testing.T) appauth.VerifierKeySet {
	t.Helper()
	keySet, err := appauth.NewVerifierKeySet(appauth.VerifierKeySetConfig{
		KeySetID:              "e2e-key-set",
		CredentialLookupKey:   bytesWithIncrementingSeed(1, appauth.MinVerifierKeyBytes),
		CredentialVerifierKey: bytesWithIncrementingSeed(2, appauth.MinVerifierKeyBytes),
		TokenLookupKey:        bytesWithIncrementingSeed(3, appauth.MinVerifierKeyBytes),
		TokenVerifierKey:      bytesWithIncrementingSeed(4, appauth.MinVerifierKeyBytes),
	})
	if err != nil {
		t.Fatalf("NewVerifierKeySet() error = %v, want nil", err)
	}
	return keySet
}

func bytesWithIncrementingSeed(seed byte, length int) []byte {
	out := make([]byte, length)
	for i := range out {
		out[i] = seed + byte(i)
	}
	return out
}

type e2eClock struct {
	now time.Time
}

func (c e2eClock) Now() time.Time {
	return c.now
}

type e2eTokenRecordIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2eTokenRecordIDGenerator) GenerateTokenRecordID(context.Context) (string, error) {
	return g.nextID("token-e2e")
}

type e2eSessionIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2eSessionIDGenerator) GenerateSessionID(context.Context) (string, error) {
	return g.nextID("runtime-session-e2e")
}

type e2ePlayerIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2ePlayerIDGenerator) GeneratePlayerID(context.Context) (string, error) {
	return g.nextID("player-e2e")
}

type e2ePlayerAccountEventIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2ePlayerAccountEventIDGenerator) GeneratePlayerAccountEventID(context.Context) (string, error) {
	return g.nextID("player-event-e2e")
}

type e2eCredentialRecordIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2eCredentialRecordIDGenerator) GenerateCredentialRecordID(context.Context) (string, error) {
	return g.nextID("credential-e2e")
}

type e2eFriendRelationshipIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2eFriendRelationshipIDGenerator) GenerateFriendRelationshipID(context.Context) (string, error) {
	return g.nextID("friend-relationship-e2e")
}

type e2eCurrencyWalletIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2eCurrencyWalletIDGenerator) GenerateCurrencyWalletID(context.Context) (string, error) {
	return g.nextID("currency-wallet-e2e")
}

type e2eCurrencyTransactionIDGenerator struct {
	mu   sync.Mutex
	next int
}

func (g *e2eCurrencyTransactionIDGenerator) GenerateCurrencyWalletTransactionID(context.Context) (string, error) {
	return g.nextID("currency-transaction-e2e")
}

func (g *e2eTokenRecordIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

func (g *e2eSessionIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

func (g *e2ePlayerIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

func (g *e2ePlayerAccountEventIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

func (g *e2eCredentialRecordIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

func (g *e2eFriendRelationshipIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

func (g *e2eCurrencyWalletIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

func (g *e2eCurrencyTransactionIDGenerator) nextID(prefix string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return fmt.Sprintf("%s-%d", prefix, g.next), nil
}

type e2eUnitOfWorkRunner struct {
	unit *e2eUnitOfWork
}

func (r e2eUnitOfWorkRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	if fn == nil {
		return errors.New("e2e: unit-of-work function is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if r.unit == nil {
		return fn(ctx, tx.NoopUnitOfWork{})
	}
	return fn(ctx, r.unit)
}

type e2eUnitOfWork struct {
	authenticationRepository authenticationmodule.Repository
	playerRepository         playermodule.Repository
	sessionRepository        sessionmodule.Repository
	storageRepository        storagemodule.Repository
	friendsRepository        friendsmodule.Repository
	currencyRepository       currencymodule.Repository
}

func (u *e2eUnitOfWork) Context() context.Context {
	return context.Background()
}

func (u *e2eUnitOfWork) NewAuthenticationRepository() (authenticationmodule.Repository, error) {
	if u == nil || u.authenticationRepository == nil {
		return nil, errors.New("e2e: authentication repository unavailable")
	}
	return u.authenticationRepository, nil
}

func (u *e2eUnitOfWork) NewPlayerAccountRepository() (playermodule.Repository, error) {
	if u == nil || u.playerRepository == nil {
		return nil, errors.New("e2e: player repository unavailable")
	}
	return u.playerRepository, nil
}

func (u *e2eUnitOfWork) NewSessionRepository() (sessionmodule.Repository, error) {
	if u == nil || u.sessionRepository == nil {
		return nil, errors.New("e2e: session repository unavailable")
	}
	return u.sessionRepository, nil
}

func (u *e2eUnitOfWork) NewStorageObjectRepository() (storagemodule.Repository, error) {
	if u == nil || u.storageRepository == nil {
		return nil, errors.New("e2e: storage repository unavailable")
	}
	return u.storageRepository, nil
}

func (u *e2eUnitOfWork) NewFriendRelationshipRepository() (friendsmodule.Repository, error) {
	if u == nil || u.friendsRepository == nil {
		return nil, errors.New("e2e: friends relationship repository unavailable")
	}
	return u.friendsRepository, nil
}

func (u *e2eUnitOfWork) NewCurrencyWalletRepository() (currencymodule.Repository, error) {
	if u == nil || u.currencyRepository == nil {
		return nil, errors.New("e2e: currency wallet repository unavailable")
	}
	return u.currencyRepository, nil
}

type e2eAuthenticationRepository struct {
	mu                    sync.Mutex
	credentialsByLookup   map[string]authenticationmodule.CredentialRecord
	tokensByLookup        map[string]authenticationmodule.TokenRecord
	tokenLookupByRecordID map[string]string
	credentialMutations   []authenticationmodule.StoreCredentialMutation
	tokenMutations        []authenticationmodule.StoreTokenMutation
}

func newE2EAuthenticationRepository() *e2eAuthenticationRepository {
	return &e2eAuthenticationRepository{
		credentialsByLookup:   make(map[string]authenticationmodule.CredentialRecord),
		tokensByLookup:        make(map[string]authenticationmodule.TokenRecord),
		tokenLookupByRecordID: make(map[string]string),
	}
}

func (r *e2eAuthenticationRepository) StoreCredential(_ context.Context, mutation authenticationmodule.StoreCredentialMutation) (authenticationmodule.CredentialRecord, error) {
	mutation, err := authenticationmodule.NormalizeStoreCredentialMutation(mutation)
	if err != nil {
		return authenticationmodule.CredentialRecord{}, err
	}

	record := authenticationmodule.CredentialRecord{
		CredentialRecordID:       mutation.CredentialRecordID,
		PlayerID:                 mutation.PlayerID,
		CredentialKind:           mutation.CredentialKind,
		CredentialStatus:         authenticationmodule.CredentialStatusActive,
		CredentialLookupDigest:   cloneBytesForE2E(mutation.CredentialLookupDigest),
		CredentialVerifierDigest: cloneBytesForE2E(mutation.CredentialVerifierDigest),
		VerifierAlgorithm:        mutation.VerifierAlgorithm,
		VerifierVersion:          mutation.VerifierVersion,
		VerifierKeyID:            mutation.VerifierKeyID,
		ClientInstanceIDDigest:   cloneBytesForE2E(mutation.ClientInstanceIDDigest),
		CreatedAt:                mutation.OccurredAt,
		UpdatedAt:                mutation.OccurredAt,
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.credentialMutations = append(r.credentialMutations, cloneStoreCredentialMutationForE2E(mutation))
	r.credentialsByLookup[digestKey(mutation.CredentialLookupDigest)] = record
	return cloneCredentialRecordForE2E(record), nil
}

func (r *e2eAuthenticationRepository) FindCredentialByLookupDigest(_ context.Context, digest []byte) (authenticationmodule.CredentialRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.credentialsByLookup[digestKey(digest)]
	if !ok {
		return authenticationmodule.CredentialRecord{}, errors.New("e2e: credential lookup miss")
	}
	return cloneCredentialRecordForE2E(record), nil
}

func (r *e2eAuthenticationRepository) StoreToken(_ context.Context, mutation authenticationmodule.StoreTokenMutation) (authenticationmodule.TokenRecord, error) {
	mutation, err := authenticationmodule.NormalizeStoreTokenMutation(mutation)
	if err != nil {
		return authenticationmodule.TokenRecord{}, err
	}

	record := authenticationmodule.TokenRecord{
		TokenRecordID:       mutation.TokenRecordID,
		TokenKind:           mutation.TokenKind,
		TokenStatus:         authenticationmodule.TokenStatusActive,
		ActorKind:           mutation.ActorKind,
		PlayerID:            mutation.PlayerID,
		CredentialRecordID:  mutation.CredentialRecordID,
		TokenLookupDigest:   cloneBytesForE2E(mutation.TokenLookupDigest),
		TokenVerifierDigest: cloneBytesForE2E(mutation.TokenVerifierDigest),
		VerifierAlgorithm:   mutation.VerifierAlgorithm,
		VerifierVersion:     mutation.VerifierVersion,
		VerifierKeyID:       mutation.VerifierKeyID,
		Audience:            mutation.Audience,
		IssuedAt:            mutation.IssuedAt,
		ExpiresAt:           mutation.ExpiresAt,
		CreatedAt:           mutation.IssuedAt,
		UpdatedAt:           mutation.IssuedAt,
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokenMutations = append(r.tokenMutations, cloneStoreTokenMutationForE2E(mutation))
	lookupKey := digestKey(mutation.TokenLookupDigest)
	r.tokensByLookup[lookupKey] = record
	r.tokenLookupByRecordID[record.TokenRecordID] = lookupKey
	return cloneTokenRecordForE2E(record), nil
}

func (r *e2eAuthenticationRepository) FindTokenByLookupDigest(_ context.Context, digest []byte) (authenticationmodule.TokenRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.tokensByLookup[digestKey(digest)]
	if !ok {
		return authenticationmodule.TokenRecord{}, errors.New("e2e: token lookup miss")
	}
	return cloneTokenRecordForE2E(record), nil
}

func (r *e2eAuthenticationRepository) RevokeCredential(context.Context, authenticationmodule.RevokeCredentialMutation) error {
	return errors.New("e2e: unexpected credential revocation")
}

func (r *e2eAuthenticationRepository) RevokeToken(_ context.Context, mutation authenticationmodule.RevokeTokenMutation) error {
	mutation, err := authenticationmodule.NormalizeRevokeTokenMutation(mutation)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	lookupKey, ok := r.tokenLookupByRecordID[mutation.TokenRecordID]
	if !ok {
		return errors.New("e2e: token revocation lookup miss")
	}
	record := r.tokensByLookup[lookupKey]
	revokedAt := mutation.RevokedAt
	record.TokenStatus = authenticationmodule.TokenStatusRevoked
	record.RevokedAt = &revokedAt
	record.RevokedReason = mutation.RevokedReason
	record.UpdatedAt = mutation.RevokedAt
	r.tokensByLookup[lookupKey] = record
	return nil
}

func (r *e2eAuthenticationRepository) expireTokenForE2E(t *testing.T, tokenRecordID string, expiresAt time.Time) {
	t.Helper()

	r.mu.Lock()
	defer r.mu.Unlock()
	lookupKey, ok := r.tokenLookupByRecordID[strings.TrimSpace(tokenRecordID)]
	if !ok {
		t.Fatalf("token record %q not found for e2e expiration", tokenRecordID)
	}
	record := r.tokensByLookup[lookupKey]
	record.ExpiresAt = expiresAt.UTC()
	record.UpdatedAt = expiresAt.UTC()
	r.tokensByLookup[lookupKey] = record
}

func (r *e2eAuthenticationRepository) ListTokensEligibleForCleanup(context.Context, authenticationmodule.TokenCleanupQuery) ([]authenticationmodule.TokenRecord, error) {
	return nil, errors.New("e2e: unexpected token cleanup query")
}

func (r *e2eAuthenticationRepository) containsRawSecret(secrets ...string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	text := strings.Builder{}
	for _, mutation := range r.credentialMutations {
		text.WriteString(base64.RawURLEncoding.EncodeToString(mutation.CredentialLookupDigest))
		text.WriteString(base64.RawURLEncoding.EncodeToString(mutation.CredentialVerifierDigest))
	}
	for _, mutation := range r.tokenMutations {
		text.WriteString(base64.RawURLEncoding.EncodeToString(mutation.TokenLookupDigest))
		text.WriteString(base64.RawURLEncoding.EncodeToString(mutation.TokenVerifierDigest))
	}
	joined := text.String()
	for _, secret := range secrets {
		if secret != "" && strings.Contains(joined, secret) {
			return true
		}
	}
	return false
}

type e2ePlayerRepository struct {
	mu       sync.Mutex
	accounts map[string]playermodule.Account
}

func newE2EPlayerRepository() *e2ePlayerRepository {
	return &e2ePlayerRepository{accounts: make(map[string]playermodule.Account)}
}

func (r *e2ePlayerRepository) CreatePlayerAccount(_ context.Context, mutation playermodule.CreatePlayerAccountMutation) (playermodule.Account, error) {
	mutation, err := playermodule.NormalizeCreatePlayerAccountMutation(mutation)
	if err != nil {
		return playermodule.Account{}, err
	}

	account := playermodule.Account{
		PlayerID:     mutation.PlayerID,
		DisplayName:  mutation.DisplayName,
		AccountState: playermodule.AccountStateActive,
		CreatedAt:    mutation.OccurredAt,
		UpdatedAt:    mutation.OccurredAt,
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[account.PlayerID] = account
	return account, nil
}

func (r *e2ePlayerRepository) GetPlayerAccount(_ context.Context, playerID string) (playermodule.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	account, ok := r.accounts[strings.TrimSpace(playerID)]
	if !ok {
		return playermodule.Account{}, errors.New("e2e: player lookup miss")
	}
	return account, nil
}

type e2eSessionRepository struct {
	mu       sync.Mutex
	sessions map[string]sessionmodule.RuntimeSession
}

func newE2ESessionRepository() *e2eSessionRepository {
	return &e2eSessionRepository{sessions: make(map[string]sessionmodule.RuntimeSession)}
}

func (r *e2eSessionRepository) CreateRuntimeSession(_ context.Context, mutation sessionmodule.CreateRuntimeSessionMutation) (sessionmodule.RuntimeSession, error) {
	mutation, err := sessionmodule.NormalizeCreateRuntimeSessionMutation(mutation)
	if err != nil {
		return sessionmodule.RuntimeSession{}, err
	}
	record := sessionmodule.RuntimeSession{
		SessionID:           mutation.SessionID,
		ActorKind:           mutation.ActorKind,
		ActorID:             mutation.ActorID,
		PlayerID:            mutation.PlayerID,
		SessionStatus:       sessionmodule.SessionStatusActive,
		IssuedAt:            mutation.IssuedAt,
		ExpiresAt:           mutation.ExpiresAt,
		LastSeenAt:          mutation.LastSeenAt,
		AccessTokenRecordID: mutation.AccessTokenRecordID,
		CreatedAt:           mutation.IssuedAt,
		UpdatedAt:           mutation.IssuedAt,
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[record.SessionID] = record
	return record, nil
}

func (r *e2eSessionRepository) GetRuntimeSession(_ context.Context, query sessionmodule.GetRuntimeSessionQuery) (sessionmodule.RuntimeSession, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.sessions[strings.TrimSpace(query.SessionID)]
	if !ok {
		return sessionmodule.RuntimeSession{}, errors.New("e2e: session lookup miss")
	}
	return record, nil
}

func (r *e2eSessionRepository) FindActiveSessionByID(_ context.Context, query sessionmodule.FindActiveSessionByIDQuery) (sessionmodule.RuntimeSession, error) {
	record, err := r.GetRuntimeSession(context.Background(), sessionmodule.GetRuntimeSessionQuery{SessionID: query.SessionID})
	if err != nil {
		return sessionmodule.RuntimeSession{}, err
	}
	if record.SessionStatus != sessionmodule.SessionStatusActive || !record.ExpiresAt.After(query.ObservedAt.UTC()) {
		return sessionmodule.RuntimeSession{}, errors.New("e2e: session is not active")
	}
	return record, nil
}

func (r *e2eSessionRepository) UpdateRuntimeSessionLastSeen(context.Context, sessionmodule.UpdateRuntimeSessionLastSeenMutation) (sessionmodule.RuntimeSession, error) {
	return sessionmodule.RuntimeSession{}, errors.New("e2e: unexpected session last-seen update")
}

func (r *e2eSessionRepository) MarkRuntimeSessionExpired(context.Context, sessionmodule.MarkRuntimeSessionExpiredMutation) (sessionmodule.RuntimeSession, error) {
	return sessionmodule.RuntimeSession{}, errors.New("e2e: unexpected session expiration")
}

func (r *e2eSessionRepository) RevokeRuntimeSession(context.Context, sessionmodule.RevokeRuntimeSessionMutation) (sessionmodule.RuntimeSession, error) {
	return sessionmodule.RuntimeSession{}, errors.New("e2e: unexpected session revocation")
}

func (r *e2eSessionRepository) ListActiveSessionsForPlayer(context.Context, sessionmodule.ListActiveSessionsForPlayerQuery) ([]sessionmodule.RuntimeSession, error) {
	return nil, errors.New("e2e: unexpected active session list")
}

type e2eStorageObjectIDGenerator struct{}

func (e2eStorageObjectIDGenerator) GenerateStorageObjectID(context.Context) (string, error) {
	return "storage-object-e2e-1", nil
}

type e2eStorageRepository struct {
	mu      sync.Mutex
	clock   e2eClock
	objects map[e2eStorageObjectKey]storagemodule.StorageObject
}

type e2eStorageObjectKey struct {
	ownerKind  storagemodule.OwnerKind
	ownerID    string
	collection string
	key        string
}

func newE2EStorageRepository(clock e2eClock) *e2eStorageRepository {
	return &e2eStorageRepository{
		clock:   clock,
		objects: make(map[e2eStorageObjectKey]storagemodule.StorageObject),
	}
}

func (r *e2eStorageRepository) CreateStorageObject(_ context.Context, input storagemodule.CreateStorageObjectInput) (storagemodule.StorageObject, error) {
	normalized, err := storagemodule.NormalizeCreateStorageObjectInput(input)
	if err != nil {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("create", storagemodule.StorageObjectConflictInvalidExpectedVersion, 0, 0)
	}
	key := e2eStorageKey(normalized.Owner, normalized.Identity)

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.objects[key]; ok {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("create", storagemodule.StorageObjectConflictAlreadyExists, 0, 0)
	}
	record := storagemodule.StorageObject{
		ObjectID:  normalized.ObjectID,
		Owner:     normalized.Owner,
		Identity:  normalized.Identity,
		Value:     normalized.Value,
		Version:   normalized.InitialVersion,
		Status:    storagemodule.StorageObjectStatusActive,
		CreatedAt: r.clock.Now(),
		UpdatedAt: r.clock.Now(),
	}
	record, err = storagemodule.NormalizeStorageObjectRecord(record)
	if err != nil {
		return storagemodule.StorageObject{}, err
	}
	r.objects[key] = cloneStorageObjectForE2E(record)
	return cloneStorageObjectForE2E(record), nil
}

func (r *e2eStorageRepository) GetStorageObject(_ context.Context, input storagemodule.GetStorageObjectInput) (storagemodule.StorageObject, error) {
	normalized, err := storagemodule.NormalizeGetStorageObjectInput(input)
	if err != nil {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("get", storagemodule.StorageObjectConflictOwnerScopeMismatch, 0, 0)
	}
	key := e2eStorageKey(normalized.Owner, normalized.Identity)

	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.objects[key]
	if !ok {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("get", storagemodule.StorageObjectConflictNotFound, 0, 0)
	}
	if record.Status != storagemodule.StorageObjectStatusActive {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("get", storagemodule.StorageObjectConflictDeletedObject, 0, 0)
	}
	return cloneStorageObjectForE2E(record), nil
}

func (r *e2eStorageRepository) ListStorageObjects(_ context.Context, input storagemodule.ListStorageObjectsInput) (storagemodule.ListStorageObjectsResult, error) {
	normalized, err := storagemodule.NormalizeListStorageObjectsInput(input)
	if err != nil {
		return storagemodule.ListStorageObjectsResult{}, e2eStorageRepositoryError("list", storagemodule.StorageObjectConflictOwnerScopeMismatch, 0, 0)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	objects := make([]storagemodule.StorageObject, 0, len(r.objects))
	for _, record := range r.objects {
		if record.Status != storagemodule.StorageObjectStatusActive ||
			record.Owner != normalized.Owner ||
			record.Identity.Collection != normalized.Collection ||
			record.Identity.Key <= normalized.AfterObjectKey {
			continue
		}
		objects = append(objects, cloneStorageObjectForE2E(record))
	}
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Identity.Key < objects[j].Identity.Key
	})

	var nextKey string
	if len(objects) > normalized.Limit {
		nextKey = objects[normalized.Limit].Identity.Key
		objects = objects[:normalized.Limit]
	}
	return storagemodule.ListStorageObjectsResult{
		Objects:       objects,
		NextObjectKey: nextKey,
	}, nil
}

func (r *e2eStorageRepository) UpdateStorageObject(_ context.Context, input storagemodule.UpdateStorageObjectInput) (storagemodule.StorageObject, error) {
	normalized, err := storagemodule.NormalizeUpdateStorageObjectInput(input)
	if err != nil {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("update", storagemodule.StorageObjectConflictInvalidExpectedVersion, 0, 0)
	}
	key := e2eStorageKey(normalized.Owner, normalized.Identity)

	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.objects[key]
	if !ok {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("update", storagemodule.StorageObjectConflictNotFound, 0, 0)
	}
	if record.Status != storagemodule.StorageObjectStatusActive {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("update", storagemodule.StorageObjectConflictDeletedObject, 0, 0)
	}
	if normalized.ExpectedVersion != nil && record.Version != *normalized.ExpectedVersion {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("update", storagemodule.StorageObjectConflictVersionMismatch, *normalized.ExpectedVersion, record.Version)
	}

	record.Value = normalized.Value
	record.Version++
	record.UpdatedAt = r.clock.Now()
	record, err = storagemodule.NormalizeStorageObjectRecord(record)
	if err != nil {
		return storagemodule.StorageObject{}, err
	}
	r.objects[key] = cloneStorageObjectForE2E(record)
	return cloneStorageObjectForE2E(record), nil
}

func (r *e2eStorageRepository) DeleteStorageObject(_ context.Context, input storagemodule.DeleteStorageObjectInput) (storagemodule.StorageObject, error) {
	normalized, err := storagemodule.NormalizeDeleteStorageObjectInput(input)
	if err != nil {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("delete", storagemodule.StorageObjectConflictInvalidExpectedVersion, 0, 0)
	}
	key := e2eStorageKey(normalized.Owner, normalized.Identity)

	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.objects[key]
	if !ok {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("delete", storagemodule.StorageObjectConflictNotFound, 0, 0)
	}
	if record.Status != storagemodule.StorageObjectStatusActive {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("delete", storagemodule.StorageObjectConflictDeletedObject, 0, 0)
	}
	if normalized.ExpectedVersion != nil && record.Version != *normalized.ExpectedVersion {
		return storagemodule.StorageObject{}, e2eStorageRepositoryError("delete", storagemodule.StorageObjectConflictVersionMismatch, *normalized.ExpectedVersion, record.Version)
	}

	deletedAt := r.clock.Now()
	record.Version++
	record.Status = storagemodule.StorageObjectStatusDeleted
	record.UpdatedAt = deletedAt
	record.DeletedAt = &deletedAt
	record, err = storagemodule.NormalizeStorageObjectRecord(record)
	if err != nil {
		return storagemodule.StorageObject{}, err
	}
	r.objects[key] = cloneStorageObjectForE2E(record)
	return cloneStorageObjectForE2E(record), nil
}

func e2eStorageKey(owner storagemodule.StorageObjectOwner, identity storagemodule.StorageObjectIdentity) e2eStorageObjectKey {
	return e2eStorageObjectKey{
		ownerKind:  owner.Kind,
		ownerID:    owner.ID,
		collection: identity.Collection,
		key:        identity.Key,
	}
}

func e2eStorageRepositoryError(operation string, class storagemodule.StorageObjectConflictClass, expected storagemodule.StorageObjectVersion, actual storagemodule.StorageObjectVersion) error {
	return &storagemodule.StorageObjectRepositoryError{
		Kind: storagemodule.ErrStorageObjectConflict,
		Conflict: storagemodule.StorageObjectConflict{
			Class:          class,
			Expected:       expected,
			Actual:         actual,
			Retryable:      class == storagemodule.StorageObjectConflictVersionMismatch,
			RedactedReason: string(class),
		},
		Operation:      operation,
		RedactedReason: string(class),
	}
}

type e2eFriendRelationshipRepository struct {
	mu            sync.Mutex
	clock         e2eClock
	relationships map[e2eFriendRelationshipKey]friendsmodule.FriendRelationship
}

type e2eFriendRelationshipKey struct {
	low  string
	high string
}

func newE2EFriendRelationshipRepository(clock e2eClock) *e2eFriendRelationshipRepository {
	return &e2eFriendRelationshipRepository{
		clock:         clock,
		relationships: make(map[e2eFriendRelationshipKey]friendsmodule.FriendRelationship),
	}
}

func (r *e2eFriendRelationshipRepository) CreateOrUpdateFriendRequest(_ context.Context, input friendsmodule.SendFriendRequestInput) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeSendFriendRequestInput(input)
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("create_or_update_request", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	pair, err := friendsmodule.NormalizeFriendRelationshipPair(friendsmodule.FriendRelationshipPair{
		PlayerLowID:  normalized.Actor.PlayerID,
		PlayerHighID: normalized.TargetPlayerID,
	})
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("create_or_update_request", friendsmodule.FriendRelationshipConflictSelfRelationshipForbidden, false)
	}
	key := e2eFriendRelationshipKeyForPair(pair)
	now := r.clock.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.relationships[key]
	if !exists {
		record := friendsmodule.FriendRelationship{
			RelationshipID:      normalized.RelationshipID,
			Pair:                pair,
			LifecycleState:      friendsmodule.FriendRelationshipLifecyclePending,
			RequestedByPlayerID: normalized.Actor.PlayerID,
			Version:             friendsmodule.InitialFriendRelationshipVersion,
			CreatedAt:           now,
			UpdatedAt:           now,
			StateChangedAt:      now,
		}
		return r.storeLocked(key, record)
	}
	if current.BlockState.BlockedByLowAt != nil || current.BlockState.BlockedByHighAt != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("create_or_update_request", friendsmodule.FriendRelationshipConflictBlockedRelationship, false)
	}
	switch current.LifecycleState {
	case friendsmodule.FriendRelationshipLifecyclePending:
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("create_or_update_request", friendsmodule.FriendRelationshipConflictDuplicatePendingRequest, false)
	case friendsmodule.FriendRelationshipLifecycleFriends:
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("create_or_update_request", friendsmodule.FriendRelationshipConflictAlreadyFriends, false)
	case friendsmodule.FriendRelationshipLifecycleRejected, friendsmodule.FriendRelationshipLifecycleRemoved:
		current.RelationshipID = normalized.RelationshipID
		current.LifecycleState = friendsmodule.FriendRelationshipLifecyclePending
		current.RequestedByPlayerID = normalized.Actor.PlayerID
		current.RespondedByPlayerID = ""
		current.RemovedByPlayerID = ""
		current.Version++
		current.UpdatedAt = now
		current.StateChangedAt = now
		current.RejectedAt = nil
		current.RemovedAt = nil
		return r.storeLocked(key, current)
	default:
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("create_or_update_request", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
}

func (r *e2eFriendRelationshipRepository) GetRelationshipByPair(_ context.Context, input friendsmodule.GetFriendRelationshipInput) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeGetFriendRelationshipInput(input)
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("get_by_pair", friendsmodule.FriendRelationshipConflictPairIdentity, false)
	}
	key := e2eFriendRelationshipKeyForPair(normalized.Pair)

	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.relationships[key]
	if !ok {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("get_by_pair", friendsmodule.FriendRelationshipConflictNotFound, false)
	}
	return cloneFriendRelationshipForE2E(record), nil
}

func (r *e2eFriendRelationshipRepository) ListRelationshipsForPlayer(_ context.Context, input friendsmodule.ListFriendRelationshipsInput) (friendsmodule.ListFriendRelationshipsResult, error) {
	normalized, err := friendsmodule.NormalizeListFriendRelationshipsInput(input)
	if err != nil {
		return friendsmodule.ListFriendRelationshipsResult{}, e2eFriendRelationshipRepositoryError("list_for_player", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	relationships := make([]friendsmodule.FriendRelationship, 0, len(r.relationships))
	for _, record := range r.relationships {
		if record.Pair.PlayerLowID != normalized.PlayerID && record.Pair.PlayerHighID != normalized.PlayerID {
			continue
		}
		if e2eFriendRelationshipPairToken(record.Pair) <= normalized.AfterPairToken {
			continue
		}
		if !e2eFriendRelationshipMatchesStatus(record, normalized.Status) {
			continue
		}
		relationships = append(relationships, cloneFriendRelationshipForE2E(record))
	}
	sort.Slice(relationships, func(i, j int) bool {
		left := e2eFriendRelationshipPairToken(relationships[i].Pair)
		right := e2eFriendRelationshipPairToken(relationships[j].Pair)
		return left < right
	})

	nextPairToken := ""
	if len(relationships) > normalized.Limit {
		nextPairToken = e2eFriendRelationshipPairToken(relationships[normalized.Limit].Pair)
		relationships = relationships[:normalized.Limit]
	}
	return friendsmodule.NormalizeListFriendRelationshipsResult(friendsmodule.ListFriendRelationshipsResult{
		Relationships: relationships,
		NextPairToken: nextPairToken,
	})
}

func (r *e2eFriendRelationshipRepository) AcceptFriendRequest(_ context.Context, input friendsmodule.AcceptFriendRequestInput) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeAcceptFriendRequestInput(input)
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("accept_request", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	return r.transitionPendingRelationship("accept_request", normalized.Pair, normalized.Actor.PlayerID, normalized.ExpectedVersion, func(record friendsmodule.FriendRelationship, now time.Time) friendsmodule.FriendRelationship {
		record.LifecycleState = friendsmodule.FriendRelationshipLifecycleFriends
		record.RespondedByPlayerID = normalized.Actor.PlayerID
		record.RemovedByPlayerID = ""
		record.Version++
		record.UpdatedAt = now
		record.StateChangedAt = now
		record.RejectedAt = nil
		record.RemovedAt = nil
		return record
	})
}

func (r *e2eFriendRelationshipRepository) RejectFriendRequest(_ context.Context, input friendsmodule.RejectFriendRequestInput) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeRejectFriendRequestInput(input)
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("reject_request", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	return r.transitionPendingRelationship("reject_request", normalized.Pair, normalized.Actor.PlayerID, normalized.ExpectedVersion, func(record friendsmodule.FriendRelationship, now time.Time) friendsmodule.FriendRelationship {
		record.LifecycleState = friendsmodule.FriendRelationshipLifecycleRejected
		record.RespondedByPlayerID = normalized.Actor.PlayerID
		record.RemovedByPlayerID = ""
		record.Version++
		record.UpdatedAt = now
		record.StateChangedAt = now
		record.RejectedAt = &now
		record.RemovedAt = nil
		return record
	})
}

func (r *e2eFriendRelationshipRepository) RemoveFriend(_ context.Context, input friendsmodule.RemoveFriendInput) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeRemoveFriendInput(input)
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("remove_friend", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	key := e2eFriendRelationshipKeyForPair(normalized.Pair)
	now := r.clock.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	current, err := r.relationshipForMutationLocked("remove_friend", key, normalized.ExpectedVersion)
	if err != nil {
		return friendsmodule.FriendRelationship{}, err
	}
	if current.LifecycleState != friendsmodule.FriendRelationshipLifecyclePending &&
		current.LifecycleState != friendsmodule.FriendRelationshipLifecycleFriends {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("remove_friend", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	if current.BlockState.BlockedByLowAt != nil || current.BlockState.BlockedByHighAt != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("remove_friend", friendsmodule.FriendRelationshipConflictBlockedRelationship, false)
	}

	current.LifecycleState = friendsmodule.FriendRelationshipLifecycleRemoved
	current.RemovedByPlayerID = normalized.Actor.PlayerID
	current.Version++
	current.UpdatedAt = now
	current.StateChangedAt = now
	current.RemovedAt = &now
	return r.storeLocked(key, current)
}

func (r *e2eFriendRelationshipRepository) SetPlayerBlock(_ context.Context, input friendsmodule.BlockPlayerInput) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeBlockPlayerInput(input)
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("set_player_block", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	pair, err := friendsmodule.NormalizeFriendRelationshipPair(friendsmodule.FriendRelationshipPair{
		PlayerLowID:  normalized.Actor.PlayerID,
		PlayerHighID: normalized.TargetPlayerID,
	})
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("set_player_block", friendsmodule.FriendRelationshipConflictSelfRelationshipForbidden, false)
	}
	key := e2eFriendRelationshipKeyForPair(pair)
	now := r.clock.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	current, err := r.relationshipForMutationLocked("set_player_block", key, normalized.ExpectedVersion)
	if err != nil {
		return friendsmodule.FriendRelationship{}, err
	}
	if normalized.Actor.PlayerID == current.Pair.PlayerLowID {
		current.BlockState.BlockedByLowAt = &now
	} else {
		current.BlockState.BlockedByHighAt = &now
	}
	if current.LifecycleState == friendsmodule.FriendRelationshipLifecyclePending ||
		current.LifecycleState == friendsmodule.FriendRelationshipLifecycleFriends {
		current.LifecycleState = friendsmodule.FriendRelationshipLifecycleRemoved
		current.RemovedByPlayerID = normalized.Actor.PlayerID
		if current.RemovedAt == nil {
			current.RemovedAt = &now
		}
	}
	current.Version++
	current.UpdatedAt = now
	current.StateChangedAt = now
	return r.storeLocked(key, current)
}

func (r *e2eFriendRelationshipRepository) ClearPlayerBlock(_ context.Context, input friendsmodule.UnblockPlayerInput) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeUnblockPlayerInput(input)
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("clear_player_block", friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	pair, err := friendsmodule.NormalizeFriendRelationshipPair(friendsmodule.FriendRelationshipPair{
		PlayerLowID:  normalized.Actor.PlayerID,
		PlayerHighID: normalized.TargetPlayerID,
	})
	if err != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError("clear_player_block", friendsmodule.FriendRelationshipConflictSelfRelationshipForbidden, false)
	}
	key := e2eFriendRelationshipKeyForPair(pair)
	now := r.clock.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	current, err := r.relationshipForMutationLocked("clear_player_block", key, normalized.ExpectedVersion)
	if err != nil {
		return friendsmodule.FriendRelationship{}, err
	}
	if normalized.Actor.PlayerID == current.Pair.PlayerLowID {
		current.BlockState.BlockedByLowAt = nil
	} else {
		current.BlockState.BlockedByHighAt = nil
	}
	current.Version++
	current.UpdatedAt = now
	return r.storeLocked(key, current)
}

func (r *e2eFriendRelationshipRepository) transitionPendingRelationship(operation string, pair friendsmodule.FriendRelationshipPair, actorPlayerID string, expectedVersion *friendsmodule.FriendRelationshipVersion, transition func(friendsmodule.FriendRelationship, time.Time) friendsmodule.FriendRelationship) (friendsmodule.FriendRelationship, error) {
	key := e2eFriendRelationshipKeyForPair(pair)
	now := r.clock.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	current, err := r.relationshipForMutationLocked(operation, key, expectedVersion)
	if err != nil {
		return friendsmodule.FriendRelationship{}, err
	}
	if current.LifecycleState != friendsmodule.FriendRelationshipLifecyclePending || current.RequestedByPlayerID == actorPlayerID {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError(operation, friendsmodule.FriendRelationshipConflictInvalidTransition, false)
	}
	if current.BlockState.BlockedByLowAt != nil || current.BlockState.BlockedByHighAt != nil {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError(operation, friendsmodule.FriendRelationshipConflictBlockedRelationship, false)
	}
	current = transition(current, now)
	return r.storeLocked(key, current)
}

func (r *e2eFriendRelationshipRepository) relationshipForMutationLocked(operation string, key e2eFriendRelationshipKey, expectedVersion *friendsmodule.FriendRelationshipVersion) (friendsmodule.FriendRelationship, error) {
	current, ok := r.relationships[key]
	if !ok {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError(operation, friendsmodule.FriendRelationshipConflictNotFound, false)
	}
	if expectedVersion != nil && current.Version != *expectedVersion {
		return friendsmodule.FriendRelationship{}, e2eFriendRelationshipRepositoryError(operation, friendsmodule.FriendRelationshipConflictStaleVersion, true)
	}
	return current, nil
}

func (r *e2eFriendRelationshipRepository) storeLocked(key e2eFriendRelationshipKey, record friendsmodule.FriendRelationship) (friendsmodule.FriendRelationship, error) {
	normalized, err := friendsmodule.NormalizeFriendRelationshipRecord(record)
	if err != nil {
		return friendsmodule.FriendRelationship{}, err
	}
	r.relationships[key] = cloneFriendRelationshipForE2E(normalized)
	return cloneFriendRelationshipForE2E(normalized), nil
}

func e2eFriendRelationshipMatchesStatus(record friendsmodule.FriendRelationship, status friendsmodule.FriendRelationshipStatus) bool {
	blocked := record.BlockState.BlockedByLowAt != nil || record.BlockState.BlockedByHighAt != nil
	switch status {
	case "":
		return true
	case friendsmodule.FriendRelationshipStatusPending:
		return record.LifecycleState == friendsmodule.FriendRelationshipLifecyclePending && !blocked
	case friendsmodule.FriendRelationshipStatusFriends:
		return record.LifecycleState == friendsmodule.FriendRelationshipLifecycleFriends && !blocked
	case friendsmodule.FriendRelationshipStatusBlocked:
		return blocked
	case friendsmodule.FriendRelationshipStatusEnded:
		return (record.LifecycleState == friendsmodule.FriendRelationshipLifecycleRemoved ||
			record.LifecycleState == friendsmodule.FriendRelationshipLifecycleRejected) && !blocked
	default:
		return false
	}
}

func e2eFriendRelationshipKeyForPair(pair friendsmodule.FriendRelationshipPair) e2eFriendRelationshipKey {
	return e2eFriendRelationshipKey{low: pair.PlayerLowID, high: pair.PlayerHighID}
}

func e2eFriendRelationshipPairToken(pair friendsmodule.FriendRelationshipPair) string {
	return pair.PlayerLowID + "|" + pair.PlayerHighID
}

func e2eFriendRelationshipRepositoryError(operation string, class friendsmodule.FriendRelationshipConflictClass, retryable bool) error {
	kind := friendsmodule.ErrFriendRelationshipConflict
	if class == friendsmodule.FriendRelationshipConflictStorageUnavailable {
		kind = friendsmodule.ErrFriendRelationshipUnavailable
	}
	return &friendsmodule.FriendRelationshipRepositoryError{
		Kind: kind,
		Conflict: friendsmodule.FriendRelationshipConflict{
			Class:          class,
			Retryable:      retryable,
			RedactedReason: string(class),
		},
		Operation:      operation,
		RedactedReason: string(class),
	}
}

type e2eCurrencyWalletRepository struct {
	mu                   sync.Mutex
	clock                e2eClock
	wallets              map[currencymodule.CurrencyWalletID]currencymodule.CurrencyWallet
	walletIDByOwner      map[e2eCurrencyWalletOwnerKey]currencymodule.CurrencyWalletID
	balances             map[e2eCurrencyWalletBalanceKey]currencymodule.CurrencyWalletBalance
	transactions         []currencymodule.CurrencyWalletTransaction
	transactionByID      map[currencymodule.CurrencyWalletTransactionID]currencymodule.CurrencyWalletTransaction
	transactionByIdemKey map[e2eCurrencyWalletIdempotencyKey]currencymodule.CurrencyWalletTransaction
}

type e2eCurrencyWalletOwnerKey struct {
	kind currencymodule.CurrencyWalletOwnerKind
	id   string
}

type e2eCurrencyWalletBalanceKey struct {
	walletID     currencymodule.CurrencyWalletID
	currencyCode currencymodule.CurrencyCode
}

type e2eCurrencyWalletIdempotencyKey struct {
	walletID currencymodule.CurrencyWalletID
	scope    currencymodule.CurrencyWalletIdempotencyScope
	key      currencymodule.CurrencyWalletIdempotencyKey
}

func newE2ECurrencyWalletRepository(clock e2eClock) *e2eCurrencyWalletRepository {
	return &e2eCurrencyWalletRepository{
		clock:                clock,
		wallets:              make(map[currencymodule.CurrencyWalletID]currencymodule.CurrencyWallet),
		walletIDByOwner:      make(map[e2eCurrencyWalletOwnerKey]currencymodule.CurrencyWalletID),
		balances:             make(map[e2eCurrencyWalletBalanceKey]currencymodule.CurrencyWalletBalance),
		transactionByID:      make(map[currencymodule.CurrencyWalletTransactionID]currencymodule.CurrencyWalletTransaction),
		transactionByIdemKey: make(map[e2eCurrencyWalletIdempotencyKey]currencymodule.CurrencyWalletTransaction),
	}
}

func (r *e2eCurrencyWalletRepository) CreateCurrencyWallet(_ context.Context, input currencymodule.CreateCurrencyWalletInput) (currencymodule.CurrencyWallet, error) {
	normalized, err := currencymodule.NormalizeCreateCurrencyWalletInput(input)
	if err != nil {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("create", currencymodule.CurrencyWalletConflictStorageUnavailable, false)
	}
	now := r.clock.Now().UTC()
	record := currencymodule.CurrencyWallet{
		WalletID:       normalized.WalletID,
		Owner:          normalized.Owner,
		LifecycleState: normalized.InitialState,
		WalletVersion:  normalized.InitialVersion,
		CreatedAt:      now,
		UpdatedAt:      now,
		StateChangedAt: now,
	}
	record, err = currencymodule.NormalizeCurrencyWalletRecord(record)
	if err != nil {
		return currencymodule.CurrencyWallet{}, err
	}
	ownerKey := e2eCurrencyOwnerKey(record.Owner)

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.wallets[record.WalletID]; ok {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("create", currencymodule.CurrencyWalletConflictWalletAlreadyExists, false)
	}
	if _, ok := r.walletIDByOwner[ownerKey]; ok {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("create", currencymodule.CurrencyWalletConflictWalletAlreadyExists, false)
	}
	r.wallets[record.WalletID] = cloneCurrencyWalletForE2E(record)
	r.walletIDByOwner[ownerKey] = record.WalletID
	return cloneCurrencyWalletForE2E(record), nil
}

func (r *e2eCurrencyWalletRepository) GetCurrencyWallet(_ context.Context, input currencymodule.GetCurrencyWalletInput) (currencymodule.CurrencyWallet, error) {
	normalized, err := currencymodule.NormalizeGetCurrencyWalletInput(input)
	if err != nil {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("get", currencymodule.CurrencyWalletConflictStorageUnavailable, false)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.wallets[normalized.WalletID]
	if !ok {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("get", currencymodule.CurrencyWalletConflictWalletNotFound, false)
	}
	return cloneCurrencyWalletForE2E(record), nil
}

func (r *e2eCurrencyWalletRepository) GetCurrencyWalletForOwner(_ context.Context, input currencymodule.GetCurrencyWalletForOwnerInput) (currencymodule.CurrencyWallet, error) {
	normalized, err := currencymodule.NormalizeGetCurrencyWalletForOwnerInput(input)
	if err != nil {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("get_for_owner", currencymodule.CurrencyWalletConflictStorageUnavailable, false)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	walletID, ok := r.walletIDByOwner[e2eCurrencyOwnerKey(normalized.Owner)]
	if !ok {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("get_for_owner", currencymodule.CurrencyWalletConflictWalletNotFound, false)
	}
	record, ok := r.wallets[walletID]
	if !ok {
		return currencymodule.CurrencyWallet{}, e2eCurrencyWalletRepositoryError("get_for_owner", currencymodule.CurrencyWalletConflictWalletNotFound, false)
	}
	return cloneCurrencyWalletForE2E(record), nil
}

func (r *e2eCurrencyWalletRepository) ListCurrencyWalletBalances(_ context.Context, input currencymodule.ListCurrencyWalletBalancesInput) (currencymodule.ListCurrencyWalletBalancesResult, error) {
	normalized, err := currencymodule.NormalizeListCurrencyWalletBalancesInput(input)
	if err != nil {
		return currencymodule.ListCurrencyWalletBalancesResult{}, e2eCurrencyWalletRepositoryError("list_balances", currencymodule.CurrencyWalletConflictStorageUnavailable, false)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.wallets[normalized.WalletID]; !ok {
		return currencymodule.ListCurrencyWalletBalancesResult{}, e2eCurrencyWalletRepositoryError("list_balances", currencymodule.CurrencyWalletConflictWalletNotFound, false)
	}
	balances := make([]currencymodule.CurrencyWalletBalance, 0, len(r.balances))
	for _, record := range r.balances {
		if record.WalletID != normalized.WalletID || record.CurrencyCode <= normalized.AfterCurrencyCode {
			continue
		}
		balances = append(balances, cloneCurrencyWalletBalanceForE2E(record))
	}
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].CurrencyCode < balances[j].CurrencyCode
	})
	var nextCurrencyCode currencymodule.CurrencyCode
	if len(balances) > normalized.Limit {
		nextCurrencyCode = balances[normalized.Limit].CurrencyCode
		balances = balances[:normalized.Limit]
	}
	return currencymodule.NormalizeListCurrencyWalletBalancesResult(currencymodule.ListCurrencyWalletBalancesResult{
		Balances:         balances,
		NextCurrencyCode: nextCurrencyCode,
	})
}

func (r *e2eCurrencyWalletRepository) RecordCurrencyGrant(_ context.Context, input currencymodule.RecordCurrencyGrantInput) (currencymodule.CurrencyWalletTransaction, error) {
	normalized, err := currencymodule.NormalizeRecordCurrencyGrantInput(input)
	if err != nil {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError("record_grant", currencymodule.CurrencyWalletConflictStorageUnavailable, false)
	}
	return r.recordCurrencyMutation("record_grant", normalized.WalletID, normalized.TransactionID, normalized.CurrencyCode, normalized.Amount, normalized.IdempotencyKey, normalized.IdempotencyScope, normalized.Actor, normalized.ReasonCode, normalized.ExternalReference, normalized.MetadataJSON, normalized.ExpectedWalletVersion, normalized.ExpectedBalanceVersion)
}

func (r *e2eCurrencyWalletRepository) RecordCurrencySpend(_ context.Context, input currencymodule.RecordCurrencySpendInput) (currencymodule.CurrencyWalletTransaction, error) {
	normalized, err := currencymodule.NormalizeRecordCurrencySpendInput(input)
	if err != nil {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError("record_spend", currencymodule.CurrencyWalletConflictStorageUnavailable, false)
	}
	return r.recordCurrencyMutation("record_spend", normalized.WalletID, normalized.TransactionID, normalized.CurrencyCode, -normalized.Amount, normalized.IdempotencyKey, normalized.IdempotencyScope, normalized.Actor, normalized.ReasonCode, normalized.ExternalReference, normalized.MetadataJSON, normalized.ExpectedWalletVersion, normalized.ExpectedBalanceVersion)
}

func (r *e2eCurrencyWalletRepository) ListCurrencyWalletTransactions(_ context.Context, input currencymodule.ListCurrencyWalletTransactionsInput) (currencymodule.ListCurrencyWalletTransactionsResult, error) {
	normalized, err := currencymodule.NormalizeListCurrencyWalletTransactionsInput(input)
	if err != nil {
		return currencymodule.ListCurrencyWalletTransactionsResult{}, e2eCurrencyWalletRepositoryError("list_transactions", currencymodule.CurrencyWalletConflictStorageUnavailable, false)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.wallets[normalized.WalletID]; !ok {
		return currencymodule.ListCurrencyWalletTransactionsResult{}, e2eCurrencyWalletRepositoryError("list_transactions", currencymodule.CurrencyWalletConflictWalletNotFound, false)
	}
	transactions := make([]currencymodule.CurrencyWalletTransaction, 0, len(r.transactions))
	for _, record := range r.transactions {
		if record.WalletID != normalized.WalletID {
			continue
		}
		if normalized.CurrencyCode != "" && record.CurrencyCode != normalized.CurrencyCode {
			continue
		}
		if !normalized.AfterTransactionTime.IsZero() {
			if record.CreatedAt.Before(normalized.AfterTransactionTime) {
				continue
			}
			if record.CreatedAt.Equal(normalized.AfterTransactionTime) && record.TransactionID <= normalized.AfterTransactionID {
				continue
			}
		} else if normalized.AfterTransactionID != "" && record.TransactionID <= normalized.AfterTransactionID {
			continue
		}
		transactions = append(transactions, cloneCurrencyWalletTransactionForE2E(record))
	}
	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].CreatedAt.Equal(transactions[j].CreatedAt) {
			return transactions[i].TransactionID < transactions[j].TransactionID
		}
		return transactions[i].CreatedAt.Before(transactions[j].CreatedAt)
	})
	var nextTransactionID currencymodule.CurrencyWalletTransactionID
	var nextTransactionCreateAt time.Time
	if len(transactions) > normalized.Limit {
		nextTransactionID = transactions[normalized.Limit].TransactionID
		nextTransactionCreateAt = transactions[normalized.Limit].CreatedAt
		transactions = transactions[:normalized.Limit]
	}
	return currencymodule.NormalizeListCurrencyWalletTransactionsResult(currencymodule.ListCurrencyWalletTransactionsResult{
		Transactions:            transactions,
		NextTransactionID:       nextTransactionID,
		NextTransactionCreateAt: nextTransactionCreateAt,
	})
}

func (r *e2eCurrencyWalletRepository) recordCurrencyMutation(
	operation string,
	walletID currencymodule.CurrencyWalletID,
	transactionID currencymodule.CurrencyWalletTransactionID,
	currencyCode currencymodule.CurrencyCode,
	amountDelta currencymodule.CurrencyAmount,
	idempotencyKey currencymodule.CurrencyWalletIdempotencyKey,
	idempotencyScope currencymodule.CurrencyWalletIdempotencyScope,
	actor currencymodule.CurrencyWalletActor,
	reasonCode string,
	externalReference string,
	metadataJSON []byte,
	expectedWalletVersion *currencymodule.CurrencyWalletVersion,
	expectedBalanceVersion *currencymodule.CurrencyBalanceVersion,
) (currencymodule.CurrencyWalletTransaction, error) {
	kind := currencymodule.CurrencyWalletTransactionGrant
	if amountDelta < 0 {
		kind = currencymodule.CurrencyWalletTransactionSpend
	}
	now := r.clock.Now().UTC().Add(time.Duration(len(r.transactions)+1) * time.Second)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.transactionByID[transactionID]; ok {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictDuplicateTransaction, false)
	}
	idemKey := e2eCurrencyWalletIdempotencyKey{
		walletID: walletID,
		scope:    idempotencyScope,
		key:      idempotencyKey,
	}
	if _, ok := r.transactionByIdemKey[idemKey]; ok {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictDuplicateTransaction, false)
	}
	wallet, ok := r.wallets[walletID]
	if !ok {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictWalletNotFound, false)
	}
	if wallet.LifecycleState != currencymodule.CurrencyWalletLifecycleActive {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictWalletNotActive, false)
	}
	if expectedWalletVersion != nil && wallet.WalletVersion != *expectedWalletVersion {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictStaleWalletVersion, true)
	}

	balanceKey := e2eCurrencyWalletBalanceKey{walletID: walletID, currencyCode: currencyCode}
	balance, exists := r.balances[balanceKey]
	if !exists {
		if amountDelta < 0 {
			return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictInsufficientBalance, false)
		}
		balance = currencymodule.CurrencyWalletBalance{
			WalletID:       walletID,
			CurrencyCode:   currencyCode,
			BalanceVersion: currencymodule.InitialCurrencyBalanceVersion,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	} else if expectedBalanceVersion != nil && balance.BalanceVersion != *expectedBalanceVersion {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictStaleBalanceVersion, true)
	}

	nextBalance := balance.BalanceAmount + amountDelta
	if nextBalance < 0 {
		return currencymodule.CurrencyWalletTransaction{}, e2eCurrencyWalletRepositoryError(operation, currencymodule.CurrencyWalletConflictInsufficientBalance, false)
	}
	balance.BalanceAmount = nextBalance
	if exists {
		balance.BalanceVersion++
	}
	balance.UpdatedAt = now
	balance, err := currencymodule.NormalizeCurrencyWalletBalanceRecord(balance)
	if err != nil {
		return currencymodule.CurrencyWalletTransaction{}, err
	}
	transaction, err := currencymodule.NormalizeCurrencyWalletTransactionRecord(currencymodule.CurrencyWalletTransaction{
		TransactionID:     transactionID,
		WalletID:          walletID,
		CurrencyCode:      currencyCode,
		TransactionKind:   kind,
		AmountDelta:       amountDelta,
		BalanceAfter:      nextBalance,
		IdempotencyKey:    idempotencyKey,
		IdempotencyScope:  idempotencyScope,
		Actor:             actor,
		ReasonCode:        reasonCode,
		ExternalReference: externalReference,
		MetadataJSON:      cloneBytesForE2E(metadataJSON),
		CreatedAt:         now,
	})
	if err != nil {
		return currencymodule.CurrencyWalletTransaction{}, err
	}

	r.balances[balanceKey] = cloneCurrencyWalletBalanceForE2E(balance)
	r.transactions = append(r.transactions, cloneCurrencyWalletTransactionForE2E(transaction))
	r.transactionByID[transaction.TransactionID] = cloneCurrencyWalletTransactionForE2E(transaction)
	r.transactionByIdemKey[idemKey] = cloneCurrencyWalletTransactionForE2E(transaction)
	return cloneCurrencyWalletTransactionForE2E(transaction), nil
}

func e2eCurrencyOwnerKey(owner currencymodule.CurrencyWalletOwner) e2eCurrencyWalletOwnerKey {
	return e2eCurrencyWalletOwnerKey{kind: owner.Kind, id: owner.ID}
}

func e2eCurrencyWalletRepositoryError(operation string, class currencymodule.CurrencyWalletConflictClass, retryable bool) error {
	kind := currencymodule.ErrCurrencyWalletConflict
	if class == currencymodule.CurrencyWalletConflictStorageUnavailable {
		kind = currencymodule.ErrCurrencyWalletUnavailable
	}
	return &currencymodule.CurrencyWalletRepositoryError{
		Kind: kind,
		Conflict: currencymodule.CurrencyWalletConflict{
			Class:          class,
			Retryable:      retryable,
			RedactedReason: string(class),
		},
		Operation:      operation,
		RedactedReason: string(class),
	}
}

type e2eRegistryConnectionBinder struct {
	binder   app.ConnectionBinder
	registry *appconnection.InMemoryRegistry
}

func (b e2eRegistryConnectionBinder) BindConnection(ctx context.Context, request app.ConnectionBindingRequest) (app.ConnectionBindingResult, error) {
	result, err := b.binder.BindConnection(ctx, request)
	if err != nil || !result.Bound || b.registry == nil {
		return result, err
	}

	_, err = b.registry.BindConnectionIdentity(ctx, appconnection.BindIdentity{
		ConnectionID:     appconnection.ConnectionID(result.ConnectionID),
		ConnectionEpoch:  appconnection.ConnectionEpoch(result.ConnectionEpoch),
		ActorKind:        appconnection.ActorKind(result.Identity.ActorKind),
		PlayerID:         appconnection.PlayerID(result.Identity.PlayerID),
		RuntimeSessionID: appconnection.RuntimeSessionID(result.Identity.SessionID),
		ValidatedAt:      result.BoundAt,
	})
	if err != nil {
		return app.ConnectionBindingResult{
				BindingStatus:   app.ConnectionBindingStatusRejected,
				PublicErrorCode: app.ErrorCodeConnectionBindingUnavailable,
			}, &app.ApplicationError{
				Code:    app.ErrorCodeConnectionBindingUnavailable,
				Message: "connection binding registry is unavailable",
				Route:   app.BindConnectionRoute(),
			}
	}
	return result, nil
}

func digestKey(digest []byte) string {
	return base64.RawURLEncoding.EncodeToString(digest)
}

func cloneBytesForE2E(value []byte) []byte {
	if value == nil {
		return nil
	}
	return append([]byte(nil), value...)
}

func cloneStoreCredentialMutationForE2E(mutation authenticationmodule.StoreCredentialMutation) authenticationmodule.StoreCredentialMutation {
	mutation.CredentialLookupDigest = cloneBytesForE2E(mutation.CredentialLookupDigest)
	mutation.CredentialVerifierDigest = cloneBytesForE2E(mutation.CredentialVerifierDigest)
	mutation.ClientInstanceIDDigest = cloneBytesForE2E(mutation.ClientInstanceIDDigest)
	return mutation
}

func cloneStoreTokenMutationForE2E(mutation authenticationmodule.StoreTokenMutation) authenticationmodule.StoreTokenMutation {
	mutation.TokenLookupDigest = cloneBytesForE2E(mutation.TokenLookupDigest)
	mutation.TokenVerifierDigest = cloneBytesForE2E(mutation.TokenVerifierDigest)
	return mutation
}

func cloneCredentialRecordForE2E(record authenticationmodule.CredentialRecord) authenticationmodule.CredentialRecord {
	record.CredentialLookupDigest = cloneBytesForE2E(record.CredentialLookupDigest)
	record.CredentialVerifierDigest = cloneBytesForE2E(record.CredentialVerifierDigest)
	record.ClientInstanceIDDigest = cloneBytesForE2E(record.ClientInstanceIDDigest)
	return record
}

func cloneTokenRecordForE2E(record authenticationmodule.TokenRecord) authenticationmodule.TokenRecord {
	record.TokenLookupDigest = cloneBytesForE2E(record.TokenLookupDigest)
	record.TokenVerifierDigest = cloneBytesForE2E(record.TokenVerifierDigest)
	return record
}

func cloneStorageObjectForE2E(record storagemodule.StorageObject) storagemodule.StorageObject {
	record.Value.JSON = cloneBytesForE2E(record.Value.JSON)
	if record.DeletedAt != nil {
		deletedAt := *record.DeletedAt
		record.DeletedAt = &deletedAt
	}
	return record
}

func cloneFriendRelationshipForE2E(record friendsmodule.FriendRelationship) friendsmodule.FriendRelationship {
	if record.BlockState.BlockedByLowAt != nil {
		blockedAt := *record.BlockState.BlockedByLowAt
		record.BlockState.BlockedByLowAt = &blockedAt
	}
	if record.BlockState.BlockedByHighAt != nil {
		blockedAt := *record.BlockState.BlockedByHighAt
		record.BlockState.BlockedByHighAt = &blockedAt
	}
	if record.RejectedAt != nil {
		rejectedAt := *record.RejectedAt
		record.RejectedAt = &rejectedAt
	}
	if record.RemovedAt != nil {
		removedAt := *record.RemovedAt
		record.RemovedAt = &removedAt
	}
	return record
}

func cloneCurrencyWalletForE2E(record currencymodule.CurrencyWallet) currencymodule.CurrencyWallet {
	if record.SuspendedAt != nil {
		suspendedAt := *record.SuspendedAt
		record.SuspendedAt = &suspendedAt
	}
	if record.ClosedAt != nil {
		closedAt := *record.ClosedAt
		record.ClosedAt = &closedAt
	}
	return record
}

func cloneCurrencyWalletBalanceForE2E(record currencymodule.CurrencyWalletBalance) currencymodule.CurrencyWalletBalance {
	return record
}

func cloneCurrencyWalletTransactionForE2E(record currencymodule.CurrencyWalletTransaction) currencymodule.CurrencyWalletTransaction {
	record.MetadataJSON = cloneBytesForE2E(record.MetadataJSON)
	return record
}
