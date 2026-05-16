package authentication

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	authenticationmodule "github.com/iceiko/vibit/runtime/internal/modules/authentication"
	playermodule "github.com/iceiko/vibit/runtime/internal/modules/player"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

func TestNewServiceRequiresUnitOfWorkRunner(t *testing.T) {
	_, err := NewService(ServiceDependencies{})
	assertServiceError(t, err, "NewService", FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
}

func TestNewServiceRejectsTypedNilUnitOfWorkRunner(t *testing.T) {
	var runner *recordingUnitOfWorkRunner
	_, err := NewService(ServiceDependencies{UnitOfWorkRunner: runner})
	assertServiceError(t, err, "NewService", FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
}

func TestNewServiceRequiresLoginDependencies(t *testing.T) {
	valid := validServiceDependencies(&recordingUnitOfWorkRunner{})
	cases := []struct {
		name   string
		mutate func(*ServiceDependencies)
	}{
		{name: "verifier key set", mutate: func(d *ServiceDependencies) { d.VerifierKeySet = VerifierKeySet{} }},
		{name: "access token random", mutate: func(d *ServiceDependencies) { d.AccessTokenRandom = nil }},
		{name: "clock", mutate: func(d *ServiceDependencies) { d.Clock = nil }},
		{name: "token record id generator", mutate: func(d *ServiceDependencies) { d.TokenRecordIDGenerator = nil }},
		{name: "access token lifetime", mutate: func(d *ServiceDependencies) { d.AccessTokenLifetime = 0 }},
		{name: "token audience", mutate: func(d *ServiceDependencies) { d.TokenAudience = " " }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dependencies := valid
			tc.mutate(&dependencies)
			_, err := NewService(dependencies)
			assertServiceError(t, err, "NewService", FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
		})
	}
}

func TestAuthenticateWithDeviceCredentialRejectsMissingProofWithoutUnitOfWork(t *testing.T) {
	runner := &recordingLoginRunner{unit: newFakeLoginUnitOfWork(nil, nil, nil)}
	service := mustNewService(t, runner)

	result, err := service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: " ",
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassMissingProof, PublicErrorAuthenticationProofMissing)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if result.Status != AuthenticationStatusRejected ||
		result.PublicErrorCode != PublicErrorAuthenticationProofMissing ||
		result.FailureClass != FailureClassMissingProof {
		t.Fatalf("result = %#v, want missing proof rejection", result)
	}
	if result.AccessToken != "" || result.TokenRecordID != "" || result.CredentialRecordID != "" || result.PlayerID != "" {
		t.Fatalf("missing proof result includes behavior-owned fields: %#v", result)
	}
}

func TestAuthenticateWithDeviceCredentialRejectsMalformedProofWithoutUnitOfWork(t *testing.T) {
	runner := &recordingLoginRunner{unit: newFakeLoginUnitOfWork(nil, nil, nil)}
	service := mustNewService(t, runner)
	proof := "not-a-valid-device-credential-proof"

	result, err := service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: proof,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassMalformedProof, PublicErrorAuthenticationProofMalformed)
	assertNoSecretLeak(t, err, proof)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if result.Status != AuthenticationStatusRejected ||
		result.PublicErrorCode != PublicErrorAuthenticationProofMalformed ||
		result.FailureClass != FailureClassMalformedProof {
		t.Fatalf("result = %#v, want malformed proof rejection", result)
	}
}

func TestAuthenticateWithDeviceCredentialStoresTokenDigestsAndReturnsTokenAfterCommit(t *testing.T) {
	fixture := newLoginFixture(t)
	service := fixture.service

	result, err := service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof:       fixture.proofText,
		RequestedPlayerID:     "ignored-by-service",
		ClientInstanceID:      "client-1",
		AccountCreationIntent: AccountCreationIntentCreate,
	})
	if err != nil {
		t.Fatalf("AuthenticateWithDeviceCredential() error = %v, want nil", err)
	}

	wantTokenText := base64.RawURLEncoding.EncodeToString(fixture.tokenMaterial)
	if result.Status != AuthenticationStatusAuthenticated ||
		result.ActorKind != app.ActorKindPlayer ||
		result.PlayerID != "player-1" ||
		result.AccessToken != wantTokenText ||
		result.TokenType != TokenTypeOpaqueAccess ||
		result.TokenRecordID != "token-record-1" ||
		result.CredentialRecordID != "credential-1" ||
		!result.IssuedAt.Equal(fixture.clock.now) ||
		!result.ExpiresAt.Equal(fixture.clock.now.Add(time.Hour)) {
		t.Fatalf("result = %#v, want authenticated token result", result)
	}

	assertEvents(t, fixture.events, []string{
		"begin",
		"new-authentication-repository",
		"new-player-account-repository",
		"find-credential",
		"get-player-account",
		"read-access-token-random",
		"generate-token-record-id",
		"clock-now",
		"store-token",
		"commit",
	})

	wantCredentialLookup, err := ComputeCredentialLookupDigest(fixture.keySet, fixture.proofMaterial)
	if err != nil {
		t.Fatalf("ComputeCredentialLookupDigest() error = %v, want nil", err)
	}
	assertBytesEqual(t, fixture.authenticationRepository.lastCredentialLookupDigest, wantCredentialLookup.Bytes())

	wantTokenLookup, err := ComputeTokenLookupDigest(fixture.keySet, fixture.tokenMaterial)
	if err != nil {
		t.Fatalf("ComputeTokenLookupDigest() error = %v, want nil", err)
	}
	wantTokenVerifier, err := ComputeTokenVerifierDigest(fixture.keySet, fixture.tokenMaterial)
	if err != nil {
		t.Fatalf("ComputeTokenVerifierDigest() error = %v, want nil", err)
	}
	stored := fixture.authenticationRepository.lastStoreTokenMutation
	if stored.TokenRecordID != "token-record-1" ||
		stored.PlayerID != "player-1" ||
		stored.CredentialRecordID != "credential-1" ||
		stored.TokenKind != tokenKindAccessToken ||
		stored.ActorKind != string(app.ActorKindPlayer) ||
		stored.VerifierAlgorithm != verifierAlgorithmHMACSHA256V1 ||
		stored.VerifierVersion != verifierVersionV1 ||
		stored.VerifierKeyID != fixture.keySet.KeySetID() ||
		stored.Audience != "gameplay" ||
		!stored.IssuedAt.Equal(fixture.clock.now) ||
		!stored.ExpiresAt.Equal(fixture.clock.now.Add(time.Hour)) ||
		stored.RequestedBy != defaultRequestedBy {
		t.Fatalf("StoreToken mutation = %#v, want configured digest-only access token", stored)
	}
	assertBytesEqual(t, stored.TokenLookupDigest, wantTokenLookup.Bytes())
	assertBytesEqual(t, stored.TokenVerifierDigest, wantTokenVerifier.Bytes())
	if equalBytes(stored.TokenLookupDigest, fixture.tokenMaterial) || equalBytes(stored.TokenVerifierDigest, fixture.tokenMaterial) {
		t.Fatal("StoreToken mutation received raw access-token material")
	}
}

func TestAuthenticateWithDeviceCredentialCollapsesLookupMissToInvalidCredential(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.authenticationRepository.findCredentialErr = errors.New("lookup miss")

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassLookupMiss, PublicErrorAuthenticationCredentialInvalid)
	assertNoSecretLeak(t, err, fixture.proofText)
	if result.Status != AuthenticationStatusRejected || result.PublicErrorCode != PublicErrorAuthenticationCredentialInvalid {
		t.Fatalf("result = %#v, want invalid credential rejection", result)
	}
	if result.AccessToken != "" || fixture.authenticationRepository.storeTokenCalls != 0 || fixture.tokenRandom.bytesRead != 0 {
		t.Fatalf("lookup miss generated or stored token: result=%#v storeCalls=%d tokenBytes=%d", result, fixture.authenticationRepository.storeTokenCalls, fixture.tokenRandom.bytesRead)
	}
}

func TestAuthenticateWithDeviceCredentialRejectsInactiveWrongKindOrWrongVerifierPosture(t *testing.T) {
	cases := []struct {
		name         string
		mutate       func(*authenticationmodule.CredentialRecord)
		failureClass FailureClass
	}{
		{name: "wrong kind", mutate: func(r *authenticationmodule.CredentialRecord) { r.CredentialKind = "password" }, failureClass: FailureClassWrongCredentialKind},
		{name: "inactive credential", mutate: func(r *authenticationmodule.CredentialRecord) {
			r.CredentialStatus = authenticationmodule.CredentialStatusRevoked
		}, failureClass: FailureClassCredentialNotActive},
		{name: "wrong algorithm", mutate: func(r *authenticationmodule.CredentialRecord) { r.VerifierAlgorithm = "other" }, failureClass: FailureClassWrongVerifierAlgorithm},
		{name: "wrong version", mutate: func(r *authenticationmodule.CredentialRecord) { r.VerifierVersion = 2 }, failureClass: FailureClassUnsupportedVersion},
		{name: "wrong key id", mutate: func(r *authenticationmodule.CredentialRecord) { r.VerifierKeyID = "other-key-set" }, failureClass: FailureClassUnknownVerifierKeyID},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newLoginFixture(t)
			tc.mutate(&fixture.authenticationRepository.credential)

			result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
				CredentialProof: fixture.proofText,
			})

			assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, tc.failureClass, PublicErrorAuthenticationCredentialInvalid)
			if result.Status != AuthenticationStatusRejected || result.PublicErrorCode != PublicErrorAuthenticationCredentialInvalid {
				t.Fatalf("result = %#v, want invalid credential rejection", result)
			}
			if fixture.tokenRandom.bytesRead != 0 || fixture.authenticationRepository.storeTokenCalls != 0 {
				t.Fatalf("rejected credential generated or stored token: tokenBytes=%d storeCalls=%d", fixture.tokenRandom.bytesRead, fixture.authenticationRepository.storeTokenCalls)
			}
		})
	}
}

func TestAuthenticateWithDeviceCredentialComparesVerifierBeforeTokenGeneration(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.authenticationRepository.credential.CredentialVerifierDigest = bytesWithSeed(250, VerifierDigestBytes)

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassVerifierMismatch, PublicErrorAuthenticationCredentialInvalid)
	if result.Status != AuthenticationStatusRejected || result.PublicErrorCode != PublicErrorAuthenticationCredentialInvalid {
		t.Fatalf("result = %#v, want verifier mismatch rejection", result)
	}
	if fixture.playerRepository.getCalls != 0 || fixture.tokenRandom.bytesRead != 0 || fixture.authenticationRepository.storeTokenCalls != 0 {
		t.Fatalf("verifier mismatch advanced too far: getPlayer=%d tokenBytes=%d storeCalls=%d", fixture.playerRepository.getCalls, fixture.tokenRandom.bytesRead, fixture.authenticationRepository.storeTokenCalls)
	}
}

func TestAuthenticateWithDeviceCredentialRequiresActivePlayerAccount(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.playerRepository.account.AccountState = playermodule.AccountStateDisabled

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassPlayerAccountNotActive, PublicErrorAuthenticationCredentialInvalid)
	if result.Status != AuthenticationStatusRejected || result.PublicErrorCode != PublicErrorAuthenticationCredentialInvalid {
		t.Fatalf("result = %#v, want invalid credential rejection", result)
	}
	if fixture.tokenRandom.bytesRead != 0 || fixture.authenticationRepository.storeTokenCalls != 0 {
		t.Fatalf("inactive player generated or stored token: tokenBytes=%d storeCalls=%d", fixture.tokenRandom.bytesRead, fixture.authenticationRepository.storeTokenCalls)
	}
}

func TestAuthenticateWithDeviceCredentialDoesNotReturnTokenWhenStoreFails(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.authenticationRepository.storeTokenErr = errors.New("store token unavailable")

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
	assertNoSecretLeak(t, err, fixture.proofText)
	assertNoSecretLeak(t, err, base64.RawURLEncoding.EncodeToString(fixture.tokenMaterial))
	if result.AccessToken != "" || result.TokenRecordID != "" || result.Status != AuthenticationStatusRejected {
		t.Fatalf("store failure result = %#v, want rejected result without token", result)
	}
}

func TestAuthenticateWithDeviceCredentialDoesNotReturnTokenWhenCommitFails(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.runner.commitErr = errors.New("commit unavailable")

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
	if result.AccessToken != "" || result.TokenRecordID != "" || result.Status != AuthenticationStatusRejected {
		t.Fatalf("commit failure result = %#v, want rejected result without token", result)
	}
}

func TestValidateAccessTokenFailsClosedWithoutRepositoryCallOrIdentity(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)
	proof := "raw-access-token-proof"

	result, err := service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: proof,
		Route: app.RouteKey{
			Kind:   app.MessageKindCommand,
			Module: "inventory",
			Name:   "GrantItem",
		},
		ConnectionID:    "connection-1",
		ConnectionEpoch: 7,
	})
	assertServiceError(t, err, OperationValidateAccessToken, FailureClassNotImplemented, PublicErrorAuthenticationNotImplemented)
	assertNoSecretLeak(t, err, proof)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if result.Status != ValidationStatusNotImplemented {
		t.Fatalf("Status = %q, want %q", result.Status, ValidationStatusNotImplemented)
	}
	if result.ProofStatus != ProofStatusNotEvaluated {
		t.Fatalf("ProofStatus = %q, want %q", result.ProofStatus, ProofStatusNotEvaluated)
	}
	if result.PublicErrorCode != PublicErrorAuthenticationNotImplemented {
		t.Fatalf("PublicErrorCode = %q, want %q", result.PublicErrorCode, PublicErrorAuthenticationNotImplemented)
	}
	if result.Identity.Status != app.IdentityValidationUnknown {
		t.Fatalf("Identity.Status = %q, want %q", result.Identity.Status, app.IdentityValidationUnknown)
	}
	if result.PlayerIDValidated || result.SessionValidated || result.PlayerID != "" || result.ActorID != "" || result.TokenRecordID != "" {
		t.Fatalf("validation result includes proven identity fields: %#v", result)
	}
}

func TestLogoutAccessTokenFailsClosedWithoutRepositoryCall(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)
	proof := "raw-access-token-proof"

	result, err := service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
		AccessToken:  proof,
		LogoutReason: "client_logout",
	})
	assertServiceError(t, err, OperationLogoutAccessToken, FailureClassNotImplemented, PublicErrorAuthenticationNotImplemented)
	assertNoSecretLeak(t, err, proof)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if result.Status != LogoutStatusNotImplemented {
		t.Fatalf("Status = %q, want %q", result.Status, LogoutStatusNotImplemented)
	}
	if result.Revoked || result.TokenRecordID != "" {
		t.Fatalf("logout result includes revocation fields: %#v", result)
	}
}

func TestRefreshAccessTokenFailsClosedAsUnsupportedWithoutRepositoryCall(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)
	proof := "raw-access-token-proof"

	result, err := service.RefreshAccessToken(context.Background(), RefreshAccessTokenRequest{
		AccessToken: proof,
	})
	assertServiceError(t, err, OperationRefreshAccessToken, FailureClassRefreshNotSupported, PublicErrorAuthenticationRefreshNotSupported)
	assertNoSecretLeak(t, err, proof)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if result.Status != RefreshStatusUnsupported {
		t.Fatalf("Status = %q, want %q", result.Status, RefreshStatusUnsupported)
	}
	if result.AccessToken != "" || result.TokenRecordID != "" {
		t.Fatalf("refresh result includes token issuance fields: %#v", result)
	}
}

func TestServiceErrorRedactsInternalDetails(t *testing.T) {
	err := &ServiceError{
		Operation:  OperationAuthenticateWithDeviceCredential,
		Class:      FailureClassVerifierMismatch,
		PublicCode: PublicErrorAuthenticationCredentialInvalid,
		Err:        errors.New("internal verifier mismatch detail"),
	}

	text := err.Error()
	for _, forbidden := range []string{
		"internal verifier mismatch detail",
		string(FailureClassVerifierMismatch),
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("ServiceError() = %q, want redacted from %q", text, forbidden)
		}
	}
	if !strings.Contains(text, string(PublicErrorAuthenticationCredentialInvalid)) {
		t.Fatalf("ServiceError() = %q, want public code", text)
	}
}

func TestServiceSkeletonExposesReservedVocabulary(t *testing.T) {
	operations := []Operation{
		OperationAuthenticateWithDeviceCredential,
		OperationValidateAccessToken,
		OperationLogoutAccessToken,
		OperationRefreshAccessToken,
	}
	for _, operation := range operations {
		if operation == "" {
			t.Fatal("operation vocabulary must not be empty")
		}
	}

	publicCodes := []PublicErrorCode{
		PublicErrorAuthenticationProofMissing,
		PublicErrorAuthenticationProofMalformed,
		PublicErrorAuthenticationCredentialInvalid,
		PublicErrorAuthenticationTokenMissing,
		PublicErrorAuthenticationTokenMalformed,
		PublicErrorAuthenticationTokenInvalid,
		PublicErrorAuthenticationCredentialUnavailable,
		PublicErrorAuthenticationTokenUnavailable,
		PublicErrorAuthenticationRefreshNotSupported,
		PublicErrorAuthenticationNotImplemented,
	}
	for _, code := range publicCodes {
		if code == "" {
			t.Fatal("public error code vocabulary must not be empty")
		}
	}
}

type loginFixture struct {
	keySet                   VerifierKeySet
	proofMaterial            []byte
	proofText                string
	tokenMaterial            []byte
	events                   *[]string
	runner                   *recordingLoginRunner
	authenticationRepository *fakeAuthenticationRepository
	playerRepository         *fakePlayerRepository
	tokenRandom              *recordingReader
	clock                    fixedClock
	service                  Service
}

func newLoginFixture(t *testing.T) loginFixture {
	t.Helper()
	keySet := mustVerifierKeySet(t)
	proofMaterial := bytesWithSeed(141, RawSecretMaterialBytes)
	proofText := base64.RawURLEncoding.EncodeToString(proofMaterial)
	credentialVerifierDigest, err := ComputeCredentialVerifierDigest(keySet, proofMaterial)
	if err != nil {
		t.Fatalf("ComputeCredentialVerifierDigest() error = %v, want nil", err)
	}

	events := []string{}
	authRepo := &fakeAuthenticationRepository{
		events: &events,
		credential: authenticationmodule.CredentialRecord{
			CredentialRecordID:       "credential-1",
			PlayerID:                 "player-1",
			CredentialKind:           credentialKindDeviceCredentialLogin,
			CredentialStatus:         authenticationmodule.CredentialStatusActive,
			CredentialVerifierDigest: credentialVerifierDigest.Bytes(),
			VerifierAlgorithm:        verifierAlgorithmHMACSHA256V1,
			VerifierVersion:          verifierVersionV1,
			VerifierKeyID:            keySet.KeySetID(),
		},
	}
	playerRepo := &fakePlayerRepository{
		events: &events,
		account: playermodule.Account{
			PlayerID:     "player-1",
			DisplayName:  "Player One",
			AccountState: playermodule.AccountStateActive,
		},
	}
	unit := newFakeLoginUnitOfWork(&events, authRepo, playerRepo)
	runner := &recordingLoginRunner{events: &events, unit: unit}
	tokenMaterial := bytesWithSeed(151, RawSecretMaterialBytes)
	tokenRandom := &recordingReader{events: &events, data: tokenMaterial, event: "read-access-token-random"}
	clock := fixedClock{events: &events, now: time.Date(2026, 5, 16, 1, 2, 3, 0, time.UTC)}
	tokenIDGenerator := &recordingTokenRecordIDGenerator{events: &events, id: "token-record-1"}
	service, err := NewService(ServiceDependencies{
		UnitOfWorkRunner:       runner,
		VerifierKeySet:         keySet,
		AccessTokenRandom:      tokenRandom,
		Clock:                  clock,
		TokenRecordIDGenerator: tokenIDGenerator,
		AccessTokenLifetime:    time.Hour,
		TokenAudience:          "gameplay",
	})
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}

	return loginFixture{
		keySet:                   keySet,
		proofMaterial:            proofMaterial,
		proofText:                proofText,
		tokenMaterial:            tokenMaterial,
		events:                   &events,
		runner:                   runner,
		authenticationRepository: authRepo,
		playerRepository:         playerRepo,
		tokenRandom:              tokenRandom,
		clock:                    clock,
		service:                  service,
	}
}

func validServiceDependencies(runner UnitOfWorkRunner) ServiceDependencies {
	return ServiceDependencies{
		UnitOfWorkRunner:       runner,
		VerifierKeySet:         mustVerifierKeySetForHelper(),
		AccessTokenRandom:      strings.NewReader(string(bytesWithSeed(201, RawSecretMaterialBytes))),
		Clock:                  fixedClock{now: time.Date(2026, 5, 16, 1, 2, 3, 0, time.UTC)},
		TokenRecordIDGenerator: staticTokenRecordIDGenerator("token-record-1"),
		AccessTokenLifetime:    time.Hour,
		TokenAudience:          "gameplay",
	}
}

func mustVerifierKeySetForHelper() VerifierKeySet {
	keySet, err := NewVerifierKeySet(validVerifierKeySetConfig())
	if err != nil {
		panic(err)
	}
	return keySet
}

func mustNewService(t *testing.T, runner UnitOfWorkRunner) Service {
	t.Helper()
	service, err := NewService(validServiceDependencies(runner))
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}
	return service
}

func assertServiceError(t *testing.T, err error, operation Operation, class FailureClass, publicCode PublicErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatal("error = nil, want ServiceError")
	}
	var serviceErr *ServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("error = %T, want *ServiceError", err)
	}
	if serviceErr.Operation != operation {
		t.Fatalf("Operation = %q, want %q", serviceErr.Operation, operation)
	}
	if serviceErr.Class != class {
		t.Fatalf("Class = %q, want %q", serviceErr.Class, class)
	}
	if serviceErr.PublicCode != publicCode {
		t.Fatalf("PublicCode = %q, want %q", serviceErr.PublicCode, publicCode)
	}
}

func assertNoSecretLeak(t *testing.T, err error, secret string) {
	t.Helper()
	if err == nil {
		t.Fatal("error = nil, want redacted error")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error %q leaks secret %q", err.Error(), secret)
	}
}

func assertEvents(t *testing.T, got *[]string, want []string) {
	t.Helper()
	if got == nil {
		t.Fatal("events = nil")
	}
	if len(*got) != len(want) {
		t.Fatalf("events = %#v, want %#v", *got, want)
	}
	for i := range want {
		if (*got)[i] != want[i] {
			t.Fatalf("events[%d] = %q, want %q; all events=%#v", i, (*got)[i], want[i], *got)
		}
	}
}

type recordingUnitOfWorkRunner struct {
	calls int
}

func (r *recordingUnitOfWorkRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	r.calls++
	if fn == nil {
		return nil
	}
	return fn(ctx, tx.NoopUnitOfWork{})
}

type recordingLoginRunner struct {
	calls     int
	events    *[]string
	unit      *fakeLoginUnitOfWork
	commitErr error
}

func (r *recordingLoginRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	r.calls++
	if r.events != nil {
		*r.events = append(*r.events, "begin")
	}
	if fn == nil {
		return nil
	}
	if r.unit == nil {
		r.unit = newFakeLoginUnitOfWork(r.events, nil, nil)
	}
	err := fn(ctx, r.unit)
	if err != nil {
		if r.events != nil {
			*r.events = append(*r.events, "rollback")
		}
		return err
	}
	if r.commitErr != nil {
		if r.events != nil {
			*r.events = append(*r.events, "commit")
		}
		return r.commitErr
	}
	if r.events != nil {
		*r.events = append(*r.events, "commit")
	}
	return nil
}

type fakeLoginUnitOfWork struct {
	ctx                      context.Context
	events                   *[]string
	authenticationRepository authenticationmodule.Repository
	playerRepository         playermodule.Repository
	authenticationErr        error
	playerErr                error
}

func newFakeLoginUnitOfWork(events *[]string, authenticationRepository authenticationmodule.Repository, playerRepository playermodule.Repository) *fakeLoginUnitOfWork {
	return &fakeLoginUnitOfWork{
		ctx:                      context.Background(),
		events:                   events,
		authenticationRepository: authenticationRepository,
		playerRepository:         playerRepository,
	}
}

func (u *fakeLoginUnitOfWork) Context() context.Context {
	if u.ctx == nil {
		return context.Background()
	}
	return u.ctx
}

func (u *fakeLoginUnitOfWork) NewAuthenticationRepository() (authenticationmodule.Repository, error) {
	if u.events != nil {
		*u.events = append(*u.events, "new-authentication-repository")
	}
	if u.authenticationErr != nil {
		return nil, u.authenticationErr
	}
	return u.authenticationRepository, nil
}

func (u *fakeLoginUnitOfWork) NewPlayerAccountRepository() (playermodule.Repository, error) {
	if u.events != nil {
		*u.events = append(*u.events, "new-player-account-repository")
	}
	if u.playerErr != nil {
		return nil, u.playerErr
	}
	return u.playerRepository, nil
}

type fakeAuthenticationRepository struct {
	events                     *[]string
	credential                 authenticationmodule.CredentialRecord
	findCredentialErr          error
	storeTokenErr              error
	lastCredentialLookupDigest []byte
	lastStoreTokenMutation     authenticationmodule.StoreTokenMutation
	findCredentialCalls        int
	storeTokenCalls            int
}

func (r *fakeAuthenticationRepository) StoreCredential(context.Context, authenticationmodule.StoreCredentialMutation) (authenticationmodule.CredentialRecord, error) {
	return authenticationmodule.CredentialRecord{}, errors.New("unexpected StoreCredential call")
}

func (r *fakeAuthenticationRepository) FindCredentialByLookupDigest(_ context.Context, digest []byte) (authenticationmodule.CredentialRecord, error) {
	r.findCredentialCalls++
	if r.events != nil {
		*r.events = append(*r.events, "find-credential")
	}
	r.lastCredentialLookupDigest = cloneBytes(digest)
	if r.findCredentialErr != nil {
		return authenticationmodule.CredentialRecord{}, r.findCredentialErr
	}
	return r.credential, nil
}

func (r *fakeAuthenticationRepository) StoreToken(_ context.Context, mutation authenticationmodule.StoreTokenMutation) (authenticationmodule.TokenRecord, error) {
	r.storeTokenCalls++
	if r.events != nil {
		*r.events = append(*r.events, "store-token")
	}
	r.lastStoreTokenMutation = mutation
	if r.storeTokenErr != nil {
		return authenticationmodule.TokenRecord{}, r.storeTokenErr
	}
	return authenticationmodule.TokenRecord{
		TokenRecordID:      mutation.TokenRecordID,
		TokenKind:          mutation.TokenKind,
		TokenStatus:        authenticationmodule.TokenStatusActive,
		ActorKind:          mutation.ActorKind,
		PlayerID:           mutation.PlayerID,
		CredentialRecordID: mutation.CredentialRecordID,
	}, nil
}

func (r *fakeAuthenticationRepository) FindTokenByLookupDigest(context.Context, []byte) (authenticationmodule.TokenRecord, error) {
	return authenticationmodule.TokenRecord{}, errors.New("unexpected FindTokenByLookupDigest call")
}

func (r *fakeAuthenticationRepository) RevokeCredential(context.Context, authenticationmodule.RevokeCredentialMutation) error {
	return errors.New("unexpected RevokeCredential call")
}

func (r *fakeAuthenticationRepository) RevokeToken(context.Context, authenticationmodule.RevokeTokenMutation) error {
	return errors.New("unexpected RevokeToken call")
}

func (r *fakeAuthenticationRepository) ListTokensEligibleForCleanup(context.Context, authenticationmodule.TokenCleanupQuery) ([]authenticationmodule.TokenRecord, error) {
	return nil, errors.New("unexpected ListTokensEligibleForCleanup call")
}

type fakePlayerRepository struct {
	events    *[]string
	account   playermodule.Account
	getErr    error
	getCalls  int
	playerIDs []string
}

func (r *fakePlayerRepository) CreatePlayerAccount(context.Context, playermodule.CreatePlayerAccountMutation) (playermodule.Account, error) {
	return playermodule.Account{}, errors.New("unexpected CreatePlayerAccount call")
}

func (r *fakePlayerRepository) GetPlayerAccount(_ context.Context, playerID string) (playermodule.Account, error) {
	r.getCalls++
	r.playerIDs = append(r.playerIDs, playerID)
	if r.events != nil {
		*r.events = append(*r.events, "get-player-account")
	}
	if r.getErr != nil {
		return playermodule.Account{}, r.getErr
	}
	return r.account, nil
}

type recordingReader struct {
	events    *[]string
	data      []byte
	event     string
	bytesRead int
}

func (r *recordingReader) Read(target []byte) (int, error) {
	if r.events != nil && r.bytesRead == 0 {
		*r.events = append(*r.events, r.event)
	}
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := copy(target, r.data)
	r.data = r.data[n:]
	r.bytesRead += n
	return n, nil
}

type recordingTokenRecordIDGenerator struct {
	events *[]string
	id     string
	err    error
}

func (g *recordingTokenRecordIDGenerator) GenerateTokenRecordID(context.Context) (string, error) {
	if g.events != nil {
		*g.events = append(*g.events, "generate-token-record-id")
	}
	if g.err != nil {
		return "", g.err
	}
	return g.id, nil
}

type staticTokenRecordIDGenerator string

func (g staticTokenRecordIDGenerator) GenerateTokenRecordID(context.Context) (string, error) {
	return string(g), nil
}

type fixedClock struct {
	events *[]string
	now    time.Time
}

func (c fixedClock) Now() time.Time {
	if c.events != nil {
		*c.events = append(*c.events, "clock-now")
	}
	return c.now
}
