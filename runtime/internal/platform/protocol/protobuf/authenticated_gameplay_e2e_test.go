package protobuf

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
	"github.com/iceiko/vibit/runtime/internal/app/bootstrap"
	appconnection "github.com/iceiko/vibit/runtime/internal/app/connection"
	apppresence "github.com/iceiko/vibit/runtime/internal/app/presence"
	sessionmodule "github.com/iceiko/vibit/runtime/internal/app/session"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	presencev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/presence/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	storagev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/storage/v1"
	authenticationmodule "github.com/iceiko/vibit/runtime/internal/modules/authentication"
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
	runner := e2eUnitOfWorkRunner{unit: &e2eUnitOfWork{
		authenticationRepository: authenticationRepository,
		playerRepository:         playerRepository,
		sessionRepository:        sessionRepository,
		storageRepository:        storageRepository,
	}}
	service, err := appauth.NewService(appauth.ServiceDependencies{
		UnitOfWorkRunner:              runner,
		VerifierKeySet:                keySet,
		AccessTokenRandom:             bytes.NewReader(bytesWithIncrementingSeed(70, appauth.RawSecretMaterialBytes)),
		DeviceCredentialRandom:        bytes.NewReader(bytesWithIncrementingSeed(20, appauth.RawSecretMaterialBytes)),
		Clock:                         clock,
		TokenRecordIDGenerator:        e2eTokenRecordIDGenerator{},
		SessionIDGenerator:            e2eSessionIDGenerator{},
		PlayerIDGenerator:             e2ePlayerIDGenerator{},
		PlayerAccountEventIDGenerator: e2ePlayerAccountEventIDGenerator{},
		CredentialRecordIDGenerator:   e2eCredentialRecordIDGenerator{},
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

type e2eTokenRecordIDGenerator struct{}

func (e2eTokenRecordIDGenerator) GenerateTokenRecordID(context.Context) (string, error) {
	return "token-e2e-1", nil
}

type e2eSessionIDGenerator struct{}

func (e2eSessionIDGenerator) GenerateSessionID(context.Context) (string, error) {
	return "runtime-session-e2e-1", nil
}

type e2ePlayerIDGenerator struct{}

func (e2ePlayerIDGenerator) GeneratePlayerID(context.Context) (string, error) {
	return "player-e2e-1", nil
}

type e2ePlayerAccountEventIDGenerator struct{}

func (e2ePlayerAccountEventIDGenerator) GeneratePlayerAccountEventID(context.Context) (string, error) {
	return "player-event-e2e-1", nil
}

type e2eCredentialRecordIDGenerator struct{}

func (e2eCredentialRecordIDGenerator) GenerateCredentialRecordID(context.Context) (string, error) {
	return "credential-e2e-1", nil
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
