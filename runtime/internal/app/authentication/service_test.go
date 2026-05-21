package authentication

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	sessionmodule "github.com/iceiko/vibit/runtime/internal/app/session"
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
		name       string
		publicCode PublicErrorCode
		mutate     func(*ServiceDependencies)
	}{
		{name: "verifier key set", publicCode: PublicErrorAuthenticationCredentialUnavailable, mutate: func(d *ServiceDependencies) { d.VerifierKeySet = VerifierKeySet{} }},
		{name: "access token random", publicCode: PublicErrorAuthenticationCredentialUnavailable, mutate: func(d *ServiceDependencies) { d.AccessTokenRandom = nil }},
		{name: "device credential random", publicCode: PublicErrorAuthenticationOnboardingUnavailable, mutate: func(d *ServiceDependencies) { d.DeviceCredentialRandom = nil }},
		{name: "clock", publicCode: PublicErrorAuthenticationCredentialUnavailable, mutate: func(d *ServiceDependencies) { d.Clock = nil }},
		{name: "token record id generator", publicCode: PublicErrorAuthenticationCredentialUnavailable, mutate: func(d *ServiceDependencies) { d.TokenRecordIDGenerator = nil }},
		{name: "session id generator", publicCode: PublicErrorAuthenticationCredentialUnavailable, mutate: func(d *ServiceDependencies) { d.SessionIDGenerator = nil }},
		{name: "player id generator", publicCode: PublicErrorAuthenticationOnboardingUnavailable, mutate: func(d *ServiceDependencies) { d.PlayerIDGenerator = nil }},
		{name: "player account event id generator", publicCode: PublicErrorAuthenticationOnboardingUnavailable, mutate: func(d *ServiceDependencies) { d.PlayerAccountEventIDGenerator = nil }},
		{name: "credential record id generator", publicCode: PublicErrorAuthenticationOnboardingUnavailable, mutate: func(d *ServiceDependencies) { d.CredentialRecordIDGenerator = nil }},
		{name: "access token lifetime", publicCode: PublicErrorAuthenticationCredentialUnavailable, mutate: func(d *ServiceDependencies) { d.AccessTokenLifetime = 0 }},
		{name: "token audience", publicCode: PublicErrorAuthenticationCredentialUnavailable, mutate: func(d *ServiceDependencies) { d.TokenAudience = " " }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dependencies := valid
			tc.mutate(&dependencies)
			_, err := NewService(dependencies)
			assertServiceError(t, err, "NewService", FailureClassDependencyUnavailable, tc.publicCode)
		})
	}
}

func TestOnboardLocalPlayerWithDeviceCredentialRejectsInvalidRequestWithoutUnitOfWork(t *testing.T) {
	runner := &recordingLoginRunner{unit: newFakeLoginUnitOfWork(nil, nil, nil, nil)}
	service := mustNewService(t, runner)

	result, err := service.OnboardLocalPlayerWithDeviceCredential(context.Background(), LocalOnboardingDeviceCredentialIssuanceRequest{
		DisplayName: " ",
		RequestedBy: "local-operator",
	})

	assertServiceError(t, err, OperationOnboardLocalPlayerWithDeviceCredential, FailureClassInvalidRequest, PublicErrorAuthenticationOnboardingInvalid)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	assertRejectedLocalOnboarding(t, result, PublicErrorAuthenticationOnboardingInvalid, FailureClassInvalidRequest)
}

func TestOnboardLocalPlayerWithDeviceCredentialHandlesGenerationFailureWithoutUnitOfWork(t *testing.T) {
	runner := &recordingLoginRunner{unit: newFakeLoginUnitOfWork(nil, nil, nil, nil)}
	dependencies := validServiceDependencies(runner)
	dependencies.DeviceCredentialRandom = strings.NewReader("short")
	service, err := NewService(dependencies)
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}

	result, err := service.OnboardLocalPlayerWithDeviceCredential(context.Background(), LocalOnboardingDeviceCredentialIssuanceRequest{
		DisplayName: "Player One",
	})

	assertServiceError(t, err, OperationOnboardLocalPlayerWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationOnboardingUnavailable)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	assertRejectedLocalOnboarding(t, result, PublicErrorAuthenticationOnboardingUnavailable, FailureClassDependencyUnavailable)
}

func TestOnboardLocalPlayerWithDeviceCredentialCreatesPlayerAndCredentialAfterCommit(t *testing.T) {
	fixture := newOnboardingFixture(t)

	result, err := fixture.service.OnboardLocalPlayerWithDeviceCredential(context.Background(), LocalOnboardingDeviceCredentialIssuanceRequest{
		DisplayName: " Player One ",
		RequestedBy: " local-dev ",
	})
	if err != nil {
		t.Fatalf("OnboardLocalPlayerWithDeviceCredential() error = %v, want nil", err)
	}

	wantCredentialText := base64.RawURLEncoding.EncodeToString(fixture.deviceCredentialMaterial)
	if result.Status != LocalOnboardingDeviceCredentialIssuanceStatusCreated ||
		result.PublicErrorCode != "" ||
		result.FailureClass != "" ||
		result.PlayerID != "player-1" ||
		result.CredentialRecordID != "credential-1" ||
		result.DeviceCredential != wantCredentialText ||
		!result.CreatedAt.Equal(fixture.clock.now) {
		t.Fatalf("result = %#v, want created local onboarding result", result)
	}

	assertEvents(t, fixture.events, []string{
		"read-device-credential-random",
		"begin",
		"new-player-account-repository",
		"new-authentication-repository",
		"generate-player-id",
		"generate-player-account-event-id",
		"generate-credential-record-id",
		"clock-now",
		"create-player-account",
		"store-credential",
		"commit",
	})

	if fixture.playerRepository.createCalls != 1 {
		t.Fatalf("CreatePlayerAccount calls = %d, want 1", fixture.playerRepository.createCalls)
	}
	playerMutation := fixture.playerRepository.lastCreateMutation
	if playerMutation.EventID != "player-event-1" ||
		playerMutation.PlayerID != "player-1" ||
		playerMutation.DisplayName != "Player One" ||
		playerMutation.AccountState != playermodule.AccountStateActive ||
		!playerMutation.OccurredAt.Equal(fixture.clock.now) ||
		playerMutation.RequestedBy != "local-dev" {
		t.Fatalf("CreatePlayerAccount mutation = %#v, want active local player account", playerMutation)
	}

	wantLookup, err := ComputeCredentialLookupDigest(fixture.keySet, fixture.deviceCredentialMaterial)
	if err != nil {
		t.Fatalf("ComputeCredentialLookupDigest() error = %v, want nil", err)
	}
	wantVerifier, err := ComputeCredentialVerifierDigest(fixture.keySet, fixture.deviceCredentialMaterial)
	if err != nil {
		t.Fatalf("ComputeCredentialVerifierDigest() error = %v, want nil", err)
	}
	if fixture.authenticationRepository.storeCredentialCalls != 1 {
		t.Fatalf("StoreCredential calls = %d, want 1", fixture.authenticationRepository.storeCredentialCalls)
	}
	stored := fixture.authenticationRepository.lastStoreCredentialMutation
	if stored.CredentialRecordID != "credential-1" ||
		stored.PlayerID != "player-1" ||
		stored.CredentialKind != credentialKindDeviceCredentialLogin ||
		stored.VerifierAlgorithm != verifierAlgorithmHMACSHA256V1 ||
		stored.VerifierVersion != verifierVersionV1 ||
		stored.VerifierKeyID != fixture.keySet.KeySetID() ||
		len(stored.ClientInstanceIDDigest) != 0 ||
		!stored.OccurredAt.Equal(fixture.clock.now) ||
		stored.RequestedBy != "local-dev" {
		t.Fatalf("StoreCredential mutation = %#v, want digest-only local credential record", stored)
	}
	assertBytesEqual(t, stored.CredentialLookupDigest, wantLookup.Bytes())
	assertBytesEqual(t, stored.CredentialVerifierDigest, wantVerifier.Bytes())
	if equalBytes(stored.CredentialLookupDigest, fixture.deviceCredentialMaterial) ||
		equalBytes(stored.CredentialVerifierDigest, fixture.deviceCredentialMaterial) ||
		strings.Contains(fmt.Sprintf("%#v", stored), wantCredentialText) {
		t.Fatalf("StoreCredential mutation leaks raw device credential material: %#v", stored)
	}
	if strings.Contains(fmtLocalOnboardingResult(result), string(fixture.deviceCredentialMaterial)) {
		t.Fatalf("local onboarding result leaks raw credential bytes: %#v", result)
	}
}

func TestOnboardLocalPlayerWithDeviceCredentialDoesNotReturnCredentialOnFailures(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*onboardingFixture)
	}{
		{name: "player repository unavailable", mutate: func(f *onboardingFixture) {
			f.runner.unit.playerErr = errors.New("player repository unavailable")
		}},
		{name: "authentication repository unavailable", mutate: func(f *onboardingFixture) {
			f.runner.unit.authenticationErr = errors.New("authentication repository unavailable")
		}},
		{name: "player id generation fails", mutate: func(f *onboardingFixture) {
			f.playerIDGenerator.err = errors.New("player id unavailable")
		}},
		{name: "player account creation fails", mutate: func(f *onboardingFixture) {
			f.playerRepository.createErr = errors.New("player create conflict")
		}},
		{name: "credential storage fails", mutate: func(f *onboardingFixture) {
			f.authenticationRepository.storeCredentialErr = errors.New("credential store conflict")
		}},
		{name: "commit fails", mutate: func(f *onboardingFixture) {
			f.runner.commitErr = errors.New("commit unavailable")
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newOnboardingFixture(t)
			tc.mutate(&fixture)

			result, err := fixture.service.OnboardLocalPlayerWithDeviceCredential(context.Background(), LocalOnboardingDeviceCredentialIssuanceRequest{
				DisplayName: "Player One",
			})

			assertServiceError(t, err, OperationOnboardLocalPlayerWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationOnboardingUnavailable)
			assertNoSecretLeak(t, err, base64.RawURLEncoding.EncodeToString(fixture.deviceCredentialMaterial))
			assertRejectedLocalOnboarding(t, result, PublicErrorAuthenticationOnboardingUnavailable, FailureClassDependencyUnavailable)
		})
	}
}

func TestOnboardLocalPlayerWithDeviceCredentialDoesNotIssueTokenOrSession(t *testing.T) {
	fixture := newOnboardingFixture(t)

	result, err := fixture.service.OnboardLocalPlayerWithDeviceCredential(context.Background(), LocalOnboardingDeviceCredentialIssuanceRequest{
		DisplayName: "Player One",
	})
	if err != nil {
		t.Fatalf("OnboardLocalPlayerWithDeviceCredential() error = %v, want nil", err)
	}
	if result.DeviceCredential == "" {
		t.Fatal("DeviceCredential is empty, want one-time credential text")
	}
	if fixture.tokenRandom.bytesRead != 0 ||
		fixture.authenticationRepository.storeTokenCalls != 0 ||
		fixture.sessionRepository.createCalls != 0 ||
		fixture.sessionIDGenerator.calls != 0 {
		t.Fatalf("onboarding touched token/session capabilities: tokenBytes=%d storeToken=%d sessionCreate=%d sessionIDs=%d",
			fixture.tokenRandom.bytesRead,
			fixture.authenticationRepository.storeTokenCalls,
			fixture.sessionRepository.createCalls,
			fixture.sessionIDGenerator.calls)
	}
}

func TestAuthenticateWithDeviceCredentialRejectsMissingProofWithoutUnitOfWork(t *testing.T) {
	runner := &recordingLoginRunner{unit: newFakeLoginUnitOfWork(nil, nil, nil, nil)}
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
	runner := &recordingLoginRunner{unit: newFakeLoginUnitOfWork(nil, nil, nil, nil)}
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
		result.SessionID != "runtime-session-1" ||
		!result.SessionCreated ||
		!result.IssuedAt.Equal(fixture.clock.now) ||
		!result.ExpiresAt.Equal(fixture.clock.now.Add(time.Hour)) ||
		!result.SessionExpiresAt.Equal(fixture.clock.now.Add(time.Hour)) {
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
		"new-session-repository",
		"generate-session-id",
		"create-runtime-session",
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

	if fixture.sessionRepository.createCalls != 1 {
		t.Fatalf("CreateRuntimeSession calls = %d, want 1", fixture.sessionRepository.createCalls)
	}
	if fixture.playerRepository.createCalls != 0 {
		t.Fatalf("login CreatePlayerAccount calls = %d, want 0", fixture.playerRepository.createCalls)
	}
	sessionMutation := fixture.sessionRepository.lastCreateMutation
	if sessionMutation.SessionID != "runtime-session-1" ||
		sessionMutation.ActorKind != sessionmodule.ActorKindPlayer ||
		sessionMutation.ActorID != "player-1" ||
		sessionMutation.PlayerID != "player-1" ||
		sessionMutation.SessionStatus != sessionmodule.SessionStatusActive ||
		!sessionMutation.IssuedAt.Equal(fixture.clock.now) ||
		!sessionMutation.ExpiresAt.Equal(fixture.clock.now.Add(time.Hour)) ||
		!sessionMutation.LastSeenAt.Equal(fixture.clock.now) ||
		sessionMutation.AccessTokenRecordID != "token-record-1" ||
		sessionMutation.RequestedBy != defaultRequestedBy {
		t.Fatalf("CreateRuntimeSession mutation = %#v, want active player runtime session linked to token record", sessionMutation)
	}
	if strings.Contains(fmt.Sprintf("%#v", sessionMutation), base64.RawURLEncoding.EncodeToString(fixture.tokenMaterial)) ||
		strings.Contains(fmt.Sprintf("%#v", sessionMutation), base64.RawURLEncoding.EncodeToString(fixture.proofMaterial)) {
		t.Fatalf("CreateRuntimeSession mutation leaks raw proof material: %#v", sessionMutation)
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
	if fixture.sessionRepository.createCalls != 0 {
		t.Fatalf("store failure created runtime session %d times, want 0", fixture.sessionRepository.createCalls)
	}
}

func TestAuthenticateWithDeviceCredentialDoesNotReturnTokenWhenSessionRepositoryUnavailable(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.runner.unit.sessionErr = errors.New("session repository unavailable")

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
	assertNoSecretLeak(t, err, base64.RawURLEncoding.EncodeToString(fixture.tokenMaterial))
	if result.AccessToken != "" ||
		result.TokenRecordID != "" ||
		result.SessionID != "" ||
		result.SessionCreated ||
		result.Status != AuthenticationStatusRejected {
		t.Fatalf("session repository failure result = %#v, want rejected result without token or session", result)
	}
	if fixture.sessionIDGenerator.calls != 0 || fixture.sessionRepository.createCalls != 0 {
		t.Fatalf("session repository failure advanced session creation: generator=%d create=%d", fixture.sessionIDGenerator.calls, fixture.sessionRepository.createCalls)
	}
}

func TestAuthenticateWithDeviceCredentialDoesNotReturnTokenWhenSessionIDGenerationFails(t *testing.T) {
	cases := []struct {
		name string
		id   string
		err  error
	}{
		{name: "generator error", err: errors.New("session id entropy unavailable")},
		{name: "blank id", id: " "},
		{name: "whitespace padded id", id: " runtime-session-1 "},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newLoginFixture(t)
			fixture.sessionIDGenerator.id = tc.id
			fixture.sessionIDGenerator.err = tc.err

			result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
				CredentialProof: fixture.proofText,
			})

			assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
			assertNoSecretLeak(t, err, base64.RawURLEncoding.EncodeToString(fixture.tokenMaterial))
			if result.AccessToken != "" ||
				result.TokenRecordID != "" ||
				result.SessionID != "" ||
				result.SessionCreated ||
				result.Status != AuthenticationStatusRejected {
				t.Fatalf("session id failure result = %#v, want rejected result without token or session", result)
			}
			if fixture.sessionRepository.createCalls != 0 {
				t.Fatalf("session id failure created runtime session %d times, want 0", fixture.sessionRepository.createCalls)
			}
		})
	}
}

func TestAuthenticateWithDeviceCredentialDoesNotReturnTokenWhenSessionCreationFails(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.sessionRepository.createErr = errors.New("session create conflict for runtime-session-1")

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
	assertNoSecretLeak(t, err, base64.RawURLEncoding.EncodeToString(fixture.tokenMaterial))
	assertNoSecretLeak(t, err, "runtime-session-1")
	if result.AccessToken != "" ||
		result.TokenRecordID != "" ||
		result.SessionID != "" ||
		result.SessionCreated ||
		result.Status != AuthenticationStatusRejected {
		t.Fatalf("session creation failure result = %#v, want rejected result without token or session", result)
	}
	if fixture.sessionRepository.createCalls != 1 {
		t.Fatalf("session creation calls = %d, want 1", fixture.sessionRepository.createCalls)
	}
}

func TestAuthenticateWithDeviceCredentialDoesNotReturnTokenWhenCommitFails(t *testing.T) {
	fixture := newLoginFixture(t)
	fixture.runner.commitErr = errors.New("commit unavailable")

	result, err := fixture.service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof: fixture.proofText,
	})

	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable)
	if result.AccessToken != "" ||
		result.TokenRecordID != "" ||
		result.SessionID != "" ||
		result.SessionCreated ||
		result.Status != AuthenticationStatusRejected {
		t.Fatalf("commit failure result = %#v, want rejected result without token or session", result)
	}
}

func TestValidateAccessTokenRejectsMissingProofWithoutUnitOfWork(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)

	result, err := service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: " ",
		Route: app.RouteKey{
			Kind:   app.MessageKindCommand,
			Module: "inventory",
			Name:   "GrantItem",
		},
		ConnectionID:    "connection-1",
		ConnectionEpoch: 7,
	})

	assertServiceError(t, err, OperationValidateAccessToken, FailureClassMissingProof, PublicErrorAuthenticationTokenMissing)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	assertRejectedAccessTokenValidation(t, result, ProofStatusMissing, PublicErrorAuthenticationTokenMissing, FailureClassMissingProof)
}

func TestValidateAccessTokenRejectsMalformedProofWithoutUnitOfWork(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)
	proof := "not-a-valid-access-token-proof"

	result, err := service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: proof,
	})

	assertServiceError(t, err, OperationValidateAccessToken, FailureClassMalformedProof, PublicErrorAuthenticationTokenMalformed)
	assertNoSecretLeak(t, err, proof)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	assertRejectedAccessTokenValidation(t, result, ProofStatusMalformed, PublicErrorAuthenticationTokenMalformed, FailureClassMalformedProof)
}

func TestValidateAccessTokenReturnsValidatedIdentityAfterCommit(t *testing.T) {
	fixture := newValidationFixture(t)

	result, err := fixture.service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: fixture.tokenText,
		Route: app.RouteKey{
			Kind:   app.MessageKindCommand,
			Module: "inventory",
			Name:   "GrantItem",
		},
		ConnectionID:    "connection-1",
		ConnectionEpoch: 7,
	})
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v, want nil", err)
	}

	if result.Status != ValidationStatusValidated ||
		result.ProofStatus != ProofStatusValid ||
		result.ActorKind != app.ActorKindPlayer ||
		result.ActorID != "player-1" ||
		result.PlayerID != "player-1" ||
		!result.PlayerIDValidated ||
		result.SessionValidated ||
		result.TokenRecordID != "token-record-1" ||
		result.CredentialRecordID != "credential-1" {
		t.Fatalf("result = %#v, want validated player identity without validated session", result)
	}
	if result.Identity.Status != app.IdentityValidationValidated ||
		result.Identity.ActorKind != app.ActorKindPlayer ||
		result.Identity.ActorID != "player-1" ||
		result.Identity.PlayerID != "player-1" ||
		result.Identity.ConnectionID != "connection-1" ||
		result.Identity.ConnectionEpoch != 7 ||
		!result.Identity.PlayerIDValidated ||
		result.Identity.SessionValidated {
		t.Fatalf("identity = %#v, want validated player identity with SessionValidated false", result.Identity)
	}

	assertEvents(t, fixture.events, []string{
		"begin",
		"new-authentication-repository",
		"new-player-account-repository",
		"find-token",
		"clock-now",
		"get-player-account",
		"commit",
	})

	wantLookup, err := ComputeTokenLookupDigest(fixture.keySet, fixture.tokenMaterial)
	if err != nil {
		t.Fatalf("ComputeTokenLookupDigest() error = %v, want nil", err)
	}
	assertBytesEqual(t, fixture.authenticationRepository.lastTokenLookupDigest, wantLookup.Bytes())
	if fixture.authenticationRepository.storeTokenCalls != 0 || fixture.tokenRandom.bytesRead != 0 {
		t.Fatalf("validation generated or stored token: storeCalls=%d tokenBytes=%d", fixture.authenticationRepository.storeTokenCalls, fixture.tokenRandom.bytesRead)
	}
	if strings.Contains(fmtValidationResult(result), fixture.tokenText) {
		t.Fatalf("validation result leaks raw access token: %#v", result)
	}
}

func TestValidateAccessTokenCollapsesLookupMissToInvalidToken(t *testing.T) {
	fixture := newValidationFixture(t)
	fixture.authenticationRepository.findTokenErr = errors.New("lookup miss")

	result, err := fixture.service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: fixture.tokenText,
	})

	assertServiceError(t, err, OperationValidateAccessToken, FailureClassLookupMiss, PublicErrorAuthenticationTokenInvalid)
	assertNoSecretLeak(t, err, fixture.tokenText)
	assertRejectedAccessTokenValidation(t, result, ProofStatusInvalid, PublicErrorAuthenticationTokenInvalid, FailureClassLookupMiss)
	if fixture.playerRepository.getCalls != 0 {
		t.Fatalf("lookup miss called player repository %d times, want 0", fixture.playerRepository.getCalls)
	}
}

func TestValidateAccessTokenRejectsTokenPostureFailures(t *testing.T) {
	cases := []struct {
		name         string
		mutate       func(*authenticationmodule.TokenRecord, fixedClock)
		failureClass FailureClass
		proofStatus  ProofStatus
	}{
		{name: "wrong kind", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.TokenKind = "refresh_token" }, failureClass: FailureClassTokenNotActive, proofStatus: ProofStatusInvalid},
		{name: "expired status", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) {
			r.TokenStatus = authenticationmodule.TokenStatusExpired
		}, failureClass: FailureClassTokenExpired, proofStatus: ProofStatusExpired},
		{name: "revoked status", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) {
			r.TokenStatus = authenticationmodule.TokenStatusRevoked
		}, failureClass: FailureClassTokenRevoked, proofStatus: ProofStatusRevoked},
		{name: "wrong algorithm", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.VerifierAlgorithm = "other" }, failureClass: FailureClassWrongVerifierAlgorithm, proofStatus: ProofStatusInvalid},
		{name: "wrong version", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.VerifierVersion = 2 }, failureClass: FailureClassUnsupportedVersion, proofStatus: ProofStatusInvalid},
		{name: "wrong key id", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.VerifierKeyID = "other-key-set" }, failureClass: FailureClassUnknownVerifierKeyID, proofStatus: ProofStatusInvalid},
		{name: "wrong audience", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.Audience = "admin" }, failureClass: FailureClassTokenNotActive, proofStatus: ProofStatusInvalid},
		{name: "future issued", mutate: func(r *authenticationmodule.TokenRecord, clock fixedClock) {
			r.IssuedAt = clock.now.Add(time.Second)
		}, failureClass: FailureClassTokenNotActive, proofStatus: ProofStatusInvalid},
		{name: "expired by clock", mutate: func(r *authenticationmodule.TokenRecord, clock fixedClock) {
			r.ExpiresAt = clock.now.Add(-time.Second)
		}, failureClass: FailureClassTokenExpired, proofStatus: ProofStatusExpired},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newValidationFixture(t)
			tc.mutate(&fixture.authenticationRepository.token, fixture.clock)

			result, err := fixture.service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
				AccessToken: fixture.tokenText,
			})

			assertServiceError(t, err, OperationValidateAccessToken, tc.failureClass, PublicErrorAuthenticationTokenInvalid)
			assertRejectedAccessTokenValidation(t, result, tc.proofStatus, PublicErrorAuthenticationTokenInvalid, tc.failureClass)
			if fixture.playerRepository.getCalls != 0 {
				t.Fatalf("token posture failure called player repository %d times, want 0", fixture.playerRepository.getCalls)
			}
		})
	}
}

func TestValidateAccessTokenComparesVerifierBeforePlayerLookup(t *testing.T) {
	fixture := newValidationFixture(t)
	fixture.authenticationRepository.token.TokenVerifierDigest = bytesWithSeed(251, VerifierDigestBytes)

	result, err := fixture.service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: fixture.tokenText,
	})

	assertServiceError(t, err, OperationValidateAccessToken, FailureClassVerifierMismatch, PublicErrorAuthenticationTokenInvalid)
	assertNoSecretLeak(t, err, fixture.tokenText)
	assertRejectedAccessTokenValidation(t, result, ProofStatusInvalid, PublicErrorAuthenticationTokenInvalid, FailureClassVerifierMismatch)
	if fixture.playerRepository.getCalls != 0 {
		t.Fatalf("verifier mismatch called player repository %d times, want 0", fixture.playerRepository.getCalls)
	}
}

func TestValidateAccessTokenRequiresActivePlayerAccount(t *testing.T) {
	fixture := newValidationFixture(t)
	fixture.playerRepository.account.AccountState = playermodule.AccountStateDisabled

	result, err := fixture.service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: fixture.tokenText,
	})

	assertServiceError(t, err, OperationValidateAccessToken, FailureClassPlayerAccountNotActive, PublicErrorAuthenticationTokenInvalid)
	assertRejectedAccessTokenValidation(t, result, ProofStatusInvalid, PublicErrorAuthenticationTokenInvalid, FailureClassPlayerAccountNotActive)
}

func TestValidateAccessTokenMapsDependencyFailuresToStoreUnavailable(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*validationFixture)
	}{
		{name: "missing unit-of-work capability", mutate: func(f *validationFixture) {
			f.runner.unit = nil
		}},
		{name: "authentication repository unavailable", mutate: func(f *validationFixture) {
			f.runner.unit.authenticationErr = errors.New("authentication repository unavailable")
		}},
		{name: "player repository unavailable", mutate: func(f *validationFixture) {
			f.runner.unit.playerErr = errors.New("player repository unavailable")
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newValidationFixture(t)
			tc.mutate(&fixture)

			result, err := fixture.service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
				AccessToken: fixture.tokenText,
			})

			assertServiceError(t, err, OperationValidateAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable)
			assertRejectedAccessTokenValidation(t, result, ProofStatusUnavailable, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
		})
	}
}

func TestValidateAccessTokenDoesNotReturnIdentityWhenCommitFails(t *testing.T) {
	fixture := newValidationFixture(t)
	fixture.runner.commitErr = errors.New("commit unavailable")

	result, err := fixture.service.ValidateAccessToken(context.Background(), AccessTokenValidationRequest{
		AccessToken: fixture.tokenText,
	})

	assertServiceError(t, err, OperationValidateAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable)
	assertNoSecretLeak(t, err, fixture.tokenText)
	assertRejectedAccessTokenValidation(t, result, ProofStatusUnavailable, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
}

func TestLogoutAccessTokenRejectsMissingProofWithoutUnitOfWork(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)

	result, err := service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
		AccessToken:  " ",
		LogoutReason: "client_logout",
	})
	assertServiceError(t, err, OperationLogoutAccessToken, FailureClassMissingProof, PublicErrorAuthenticationTokenMissing)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	assertRejectedLogoutAccessToken(t, result, PublicErrorAuthenticationTokenMissing, FailureClassMissingProof)
}

func TestLogoutAccessTokenRejectsMalformedProofWithoutUnitOfWork(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)
	proof := "raw-access-token-proof"

	result, err := service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
		AccessToken:  proof,
		LogoutReason: "client_logout",
	})
	assertServiceError(t, err, OperationLogoutAccessToken, FailureClassMalformedProof, PublicErrorAuthenticationTokenMalformed)
	assertNoSecretLeak(t, err, proof)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	assertRejectedLogoutAccessToken(t, result, PublicErrorAuthenticationTokenMalformed, FailureClassMalformedProof)
}

func TestLogoutAccessTokenRevokesPresentedTokenAfterVerifierComparisonAndCommit(t *testing.T) {
	fixture := newValidationFixture(t)

	result, err := fixture.service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
		AccessToken:  fixture.tokenText,
		LogoutReason: "client_logout",
	})
	if err != nil {
		t.Fatalf("LogoutAccessToken() error = %v, want nil", err)
	}

	if result.Status != LogoutStatusRevoked ||
		result.PublicErrorCode != "" ||
		result.FailureClass != "" ||
		!result.Revoked ||
		result.LogoutScope != LogoutScopeToken ||
		result.TokenRecordID != "token-record-1" ||
		!result.RevokedAt.Equal(fixture.clock.now) {
		t.Fatalf("result = %#v, want revoked presented-token result", result)
	}
	assertEvents(t, fixture.events, []string{
		"begin",
		"new-authentication-repository",
		"find-token",
		"clock-now",
		"clock-now",
		"revoke-token",
		"commit",
	})

	wantLookup, err := ComputeTokenLookupDigest(fixture.keySet, fixture.tokenMaterial)
	if err != nil {
		t.Fatalf("ComputeTokenLookupDigest() error = %v, want nil", err)
	}
	assertBytesEqual(t, fixture.authenticationRepository.lastTokenLookupDigest, wantLookup.Bytes())
	if fixture.authenticationRepository.revokeTokenCalls != 1 {
		t.Fatalf("RevokeToken calls = %d, want 1", fixture.authenticationRepository.revokeTokenCalls)
	}
	revoked := fixture.authenticationRepository.lastRevokeTokenMutation
	if revoked.TokenRecordID != "token-record-1" ||
		!revoked.RevokedAt.Equal(fixture.clock.now) ||
		revoked.RevokedReason != logoutReasonPresentedAccessToken ||
		revoked.CleanupAfter != nil ||
		revoked.RequestedBy != defaultRequestedBy {
		t.Fatalf("RevokeToken mutation = %#v, want presented access-token revocation", revoked)
	}
	if fixture.playerRepository.getCalls != 0 ||
		fixture.sessionRepository.createCalls != 0 ||
		fixture.sessionRepository.revokeCalls != 0 ||
		fixture.tokenRandom.bytesRead != 0 ||
		fixture.authenticationRepository.storeTokenCalls != 0 {
		t.Fatalf("logout touched deferred capabilities: playerGet=%d sessionCreate=%d sessionRevoke=%d tokenBytes=%d storeToken=%d",
			fixture.playerRepository.getCalls,
			fixture.sessionRepository.createCalls,
			fixture.sessionRepository.revokeCalls,
			fixture.tokenRandom.bytesRead,
			fixture.authenticationRepository.storeTokenCalls)
	}
	if strings.Contains(fmtLogoutResult(result), fixture.tokenText) {
		t.Fatalf("logout result leaks raw access token: %#v", result)
	}
}

func TestLogoutAccessTokenCollapsesLookupMissToInvalidToken(t *testing.T) {
	fixture := newValidationFixture(t)
	fixture.authenticationRepository.findTokenErr = errors.New("lookup miss")

	result, err := fixture.service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
		AccessToken: fixture.tokenText,
	})

	assertServiceError(t, err, OperationLogoutAccessToken, FailureClassLookupMiss, PublicErrorAuthenticationTokenInvalid)
	assertNoSecretLeak(t, err, fixture.tokenText)
	assertRejectedLogoutAccessToken(t, result, PublicErrorAuthenticationTokenInvalid, FailureClassLookupMiss)
	if fixture.authenticationRepository.revokeTokenCalls != 0 || fixture.playerRepository.getCalls != 0 {
		t.Fatalf("lookup miss advanced too far: revoke=%d playerGet=%d", fixture.authenticationRepository.revokeTokenCalls, fixture.playerRepository.getCalls)
	}
}

func TestLogoutAccessTokenRejectsTokenPostureFailures(t *testing.T) {
	cases := []struct {
		name         string
		mutate       func(*authenticationmodule.TokenRecord, fixedClock)
		failureClass FailureClass
	}{
		{name: "wrong kind", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.TokenKind = "refresh_token" }, failureClass: FailureClassTokenNotActive},
		{name: "expired status", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) {
			r.TokenStatus = authenticationmodule.TokenStatusExpired
		}, failureClass: FailureClassTokenExpired},
		{name: "revoked status", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) {
			r.TokenStatus = authenticationmodule.TokenStatusRevoked
		}, failureClass: FailureClassTokenRevoked},
		{name: "wrong algorithm", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.VerifierAlgorithm = "other" }, failureClass: FailureClassWrongVerifierAlgorithm},
		{name: "wrong version", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.VerifierVersion = 2 }, failureClass: FailureClassUnsupportedVersion},
		{name: "wrong key id", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.VerifierKeyID = "other-key-set" }, failureClass: FailureClassUnknownVerifierKeyID},
		{name: "wrong audience", mutate: func(r *authenticationmodule.TokenRecord, _ fixedClock) { r.Audience = "admin" }, failureClass: FailureClassTokenNotActive},
		{name: "future issued", mutate: func(r *authenticationmodule.TokenRecord, clock fixedClock) {
			r.IssuedAt = clock.now.Add(time.Second)
		}, failureClass: FailureClassTokenNotActive},
		{name: "expired by clock", mutate: func(r *authenticationmodule.TokenRecord, clock fixedClock) {
			r.ExpiresAt = clock.now.Add(-time.Second)
		}, failureClass: FailureClassTokenExpired},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newValidationFixture(t)
			tc.mutate(&fixture.authenticationRepository.token, fixture.clock)

			result, err := fixture.service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
				AccessToken: fixture.tokenText,
			})

			assertServiceError(t, err, OperationLogoutAccessToken, tc.failureClass, PublicErrorAuthenticationTokenInvalid)
			assertRejectedLogoutAccessToken(t, result, PublicErrorAuthenticationTokenInvalid, tc.failureClass)
			if fixture.authenticationRepository.revokeTokenCalls != 0 || fixture.playerRepository.getCalls != 0 {
				t.Fatalf("token posture failure advanced too far: revoke=%d playerGet=%d", fixture.authenticationRepository.revokeTokenCalls, fixture.playerRepository.getCalls)
			}
		})
	}
}

func TestLogoutAccessTokenComparesVerifierBeforeRevocation(t *testing.T) {
	fixture := newValidationFixture(t)
	fixture.authenticationRepository.token.TokenVerifierDigest = bytesWithSeed(252, VerifierDigestBytes)

	result, err := fixture.service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
		AccessToken: fixture.tokenText,
	})

	assertServiceError(t, err, OperationLogoutAccessToken, FailureClassVerifierMismatch, PublicErrorAuthenticationTokenInvalid)
	assertNoSecretLeak(t, err, fixture.tokenText)
	assertRejectedLogoutAccessToken(t, result, PublicErrorAuthenticationTokenInvalid, FailureClassVerifierMismatch)
	if fixture.authenticationRepository.revokeTokenCalls != 0 || fixture.playerRepository.getCalls != 0 {
		t.Fatalf("verifier mismatch advanced too far: revoke=%d playerGet=%d", fixture.authenticationRepository.revokeTokenCalls, fixture.playerRepository.getCalls)
	}
}

func TestLogoutAccessTokenMapsDependencyFailuresToStoreUnavailable(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*validationFixture)
	}{
		{name: "missing unit-of-work capability", mutate: func(f *validationFixture) {
			f.runner.unit = nil
		}},
		{name: "authentication repository unavailable", mutate: func(f *validationFixture) {
			f.runner.unit.authenticationErr = errors.New("authentication repository unavailable")
		}},
		{name: "revoke token unavailable", mutate: func(f *validationFixture) {
			f.authenticationRepository.revokeTokenErr = errors.New("token revocation unavailable")
		}},
		{name: "commit unavailable", mutate: func(f *validationFixture) {
			f.runner.commitErr = errors.New("commit unavailable")
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newValidationFixture(t)
			tc.mutate(&fixture)

			result, err := fixture.service.LogoutAccessToken(context.Background(), LogoutAccessTokenRequest{
				AccessToken: fixture.tokenText,
			})

			assertServiceError(t, err, OperationLogoutAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable)
			assertNoSecretLeak(t, err, fixture.tokenText)
			assertRejectedLogoutAccessToken(t, result, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			if result.Revoked || result.TokenRecordID != "" || !result.RevokedAt.IsZero() {
				t.Fatalf("dependency failure result claims revocation: %#v", result)
			}
		})
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
		OperationOnboardLocalPlayerWithDeviceCredential,
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
		PublicErrorAuthenticationOnboardingInvalid,
		PublicErrorAuthenticationOnboardingUnavailable,
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
	sessionIDGenerator       *recordingSessionIDGenerator
	sessionRepository        *fakeSessionRepository
	service                  Service
}

type onboardingFixture struct {
	keySet                   VerifierKeySet
	deviceCredentialMaterial []byte
	events                   *[]string
	runner                   *recordingLoginRunner
	authenticationRepository *fakeAuthenticationRepository
	playerRepository         *fakePlayerRepository
	sessionRepository        *fakeSessionRepository
	deviceCredentialRandom   *recordingReader
	tokenRandom              *recordingReader
	clock                    fixedClock
	playerIDGenerator        *recordingPlayerIDGenerator
	playerEventIDGenerator   *recordingPlayerAccountEventIDGenerator
	credentialIDGenerator    *recordingCredentialRecordIDGenerator
	sessionIDGenerator       *recordingSessionIDGenerator
	service                  Service
}

type validationFixture struct {
	keySet                   VerifierKeySet
	tokenMaterial            []byte
	tokenText                string
	events                   *[]string
	runner                   *recordingLoginRunner
	authenticationRepository *fakeAuthenticationRepository
	playerRepository         *fakePlayerRepository
	sessionRepository        *fakeSessionRepository
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
	sessionRepo := &fakeSessionRepository{events: &events}
	unit := newFakeLoginUnitOfWork(&events, authRepo, playerRepo, sessionRepo)
	runner := &recordingLoginRunner{events: &events, unit: unit}
	tokenMaterial := bytesWithSeed(151, RawSecretMaterialBytes)
	tokenRandom := &recordingReader{events: &events, data: tokenMaterial, event: "read-access-token-random"}
	deviceCredentialRandom := &recordingReader{events: &events, data: bytesWithSeed(152, RawSecretMaterialBytes), event: "read-device-credential-random"}
	clock := fixedClock{events: &events, now: time.Date(2026, 5, 16, 1, 2, 3, 0, time.UTC)}
	tokenIDGenerator := &recordingTokenRecordIDGenerator{events: &events, id: "token-record-1"}
	sessionIDGenerator := &recordingSessionIDGenerator{events: &events, id: "runtime-session-1"}
	service, err := NewService(ServiceDependencies{
		UnitOfWorkRunner:              runner,
		VerifierKeySet:                keySet,
		AccessTokenRandom:             tokenRandom,
		DeviceCredentialRandom:        deviceCredentialRandom,
		Clock:                         clock,
		TokenRecordIDGenerator:        tokenIDGenerator,
		SessionIDGenerator:            sessionIDGenerator,
		PlayerIDGenerator:             staticPlayerIDGenerator("player-1"),
		PlayerAccountEventIDGenerator: staticPlayerAccountEventIDGenerator("player-event-1"),
		CredentialRecordIDGenerator:   staticCredentialRecordIDGenerator("credential-1"),
		AccessTokenLifetime:           time.Hour,
		TokenAudience:                 "gameplay",
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
		sessionIDGenerator:       sessionIDGenerator,
		sessionRepository:        sessionRepo,
		service:                  service,
	}
}

func newOnboardingFixture(t *testing.T) onboardingFixture {
	t.Helper()
	keySet := mustVerifierKeySet(t)
	events := []string{}
	authRepo := &fakeAuthenticationRepository{events: &events}
	playerRepo := &fakePlayerRepository{events: &events}
	sessionRepo := &fakeSessionRepository{events: &events}
	unit := newFakeLoginUnitOfWork(&events, authRepo, playerRepo, sessionRepo)
	runner := &recordingLoginRunner{events: &events, unit: unit}
	deviceCredentialMaterial := bytesWithSeed(181, RawSecretMaterialBytes)
	deviceCredentialRandom := &recordingReader{events: &events, data: deviceCredentialMaterial, event: "read-device-credential-random"}
	tokenRandom := &recordingReader{events: &events, data: bytesWithSeed(182, RawSecretMaterialBytes), event: "read-access-token-random"}
	clock := fixedClock{events: &events, now: time.Date(2026, 5, 21, 1, 2, 3, 0, time.UTC)}
	playerIDGenerator := &recordingPlayerIDGenerator{events: &events, id: "player-1"}
	playerEventIDGenerator := &recordingPlayerAccountEventIDGenerator{events: &events, id: "player-event-1"}
	credentialIDGenerator := &recordingCredentialRecordIDGenerator{events: &events, id: "credential-1"}
	sessionIDGenerator := &recordingSessionIDGenerator{events: &events, id: "runtime-session-1"}
	service, err := NewService(ServiceDependencies{
		UnitOfWorkRunner:              runner,
		VerifierKeySet:                keySet,
		AccessTokenRandom:             tokenRandom,
		DeviceCredentialRandom:        deviceCredentialRandom,
		Clock:                         clock,
		TokenRecordIDGenerator:        &recordingTokenRecordIDGenerator{events: &events, id: "token-record-1"},
		SessionIDGenerator:            sessionIDGenerator,
		PlayerIDGenerator:             playerIDGenerator,
		PlayerAccountEventIDGenerator: playerEventIDGenerator,
		CredentialRecordIDGenerator:   credentialIDGenerator,
		AccessTokenLifetime:           time.Hour,
		TokenAudience:                 "gameplay",
	})
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}

	return onboardingFixture{
		keySet:                   keySet,
		deviceCredentialMaterial: deviceCredentialMaterial,
		events:                   &events,
		runner:                   runner,
		authenticationRepository: authRepo,
		playerRepository:         playerRepo,
		sessionRepository:        sessionRepo,
		deviceCredentialRandom:   deviceCredentialRandom,
		tokenRandom:              tokenRandom,
		clock:                    clock,
		playerIDGenerator:        playerIDGenerator,
		playerEventIDGenerator:   playerEventIDGenerator,
		credentialIDGenerator:    credentialIDGenerator,
		sessionIDGenerator:       sessionIDGenerator,
		service:                  service,
	}
}

func newValidationFixture(t *testing.T) validationFixture {
	t.Helper()
	keySet := mustVerifierKeySet(t)
	tokenMaterial := bytesWithSeed(161, RawSecretMaterialBytes)
	tokenText := base64.RawURLEncoding.EncodeToString(tokenMaterial)
	tokenVerifierDigest, err := ComputeTokenVerifierDigest(keySet, tokenMaterial)
	if err != nil {
		t.Fatalf("ComputeTokenVerifierDigest() error = %v, want nil", err)
	}

	events := []string{}
	clock := fixedClock{events: &events, now: time.Date(2026, 5, 16, 1, 2, 3, 0, time.UTC)}
	authRepo := &fakeAuthenticationRepository{
		events: &events,
		token: authenticationmodule.TokenRecord{
			TokenRecordID:       "token-record-1",
			TokenKind:           tokenKindAccessToken,
			TokenStatus:         authenticationmodule.TokenStatusActive,
			ActorKind:           string(app.ActorKindPlayer),
			PlayerID:            "player-1",
			CredentialRecordID:  "credential-1",
			TokenVerifierDigest: tokenVerifierDigest.Bytes(),
			VerifierAlgorithm:   verifierAlgorithmHMACSHA256V1,
			VerifierVersion:     verifierVersionV1,
			VerifierKeyID:       keySet.KeySetID(),
			Audience:            "gameplay",
			IssuedAt:            clock.now.Add(-time.Minute),
			ExpiresAt:           clock.now.Add(time.Hour),
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
	sessionRepo := &fakeSessionRepository{events: &events}
	unit := newFakeLoginUnitOfWork(&events, authRepo, playerRepo, sessionRepo)
	runner := &recordingLoginRunner{events: &events, unit: unit}
	tokenRandom := &recordingReader{events: &events, data: bytesWithSeed(171, RawSecretMaterialBytes), event: "read-access-token-random"}
	deviceCredentialRandom := &recordingReader{events: &events, data: bytesWithSeed(172, RawSecretMaterialBytes), event: "read-device-credential-random"}
	tokenIDGenerator := &recordingTokenRecordIDGenerator{events: &events, id: "token-record-1"}
	sessionIDGenerator := &recordingSessionIDGenerator{events: &events, id: "runtime-session-1"}
	service, err := NewService(ServiceDependencies{
		UnitOfWorkRunner:              runner,
		VerifierKeySet:                keySet,
		AccessTokenRandom:             tokenRandom,
		DeviceCredentialRandom:        deviceCredentialRandom,
		Clock:                         clock,
		TokenRecordIDGenerator:        tokenIDGenerator,
		SessionIDGenerator:            sessionIDGenerator,
		PlayerIDGenerator:             staticPlayerIDGenerator("player-1"),
		PlayerAccountEventIDGenerator: staticPlayerAccountEventIDGenerator("player-event-1"),
		CredentialRecordIDGenerator:   staticCredentialRecordIDGenerator("credential-1"),
		AccessTokenLifetime:           time.Hour,
		TokenAudience:                 "gameplay",
	})
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}

	return validationFixture{
		keySet:                   keySet,
		tokenMaterial:            tokenMaterial,
		tokenText:                tokenText,
		events:                   &events,
		runner:                   runner,
		authenticationRepository: authRepo,
		playerRepository:         playerRepo,
		sessionRepository:        sessionRepo,
		tokenRandom:              tokenRandom,
		clock:                    clock,
		service:                  service,
	}
}

func validServiceDependencies(runner UnitOfWorkRunner) ServiceDependencies {
	return ServiceDependencies{
		UnitOfWorkRunner:              runner,
		VerifierKeySet:                mustVerifierKeySetForHelper(),
		AccessTokenRandom:             strings.NewReader(string(bytesWithSeed(201, RawSecretMaterialBytes))),
		DeviceCredentialRandom:        strings.NewReader(string(bytesWithSeed(202, RawSecretMaterialBytes))),
		Clock:                         fixedClock{now: time.Date(2026, 5, 16, 1, 2, 3, 0, time.UTC)},
		TokenRecordIDGenerator:        staticTokenRecordIDGenerator("token-record-1"),
		SessionIDGenerator:            staticSessionIDGenerator("runtime-session-1"),
		PlayerIDGenerator:             staticPlayerIDGenerator("player-1"),
		PlayerAccountEventIDGenerator: staticPlayerAccountEventIDGenerator("player-event-1"),
		CredentialRecordIDGenerator:   staticCredentialRecordIDGenerator("credential-1"),
		AccessTokenLifetime:           time.Hour,
		TokenAudience:                 "gameplay",
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

func assertRejectedAccessTokenValidation(t *testing.T, result AccessTokenValidationResult, proofStatus ProofStatus, publicCode PublicErrorCode, failureClass FailureClass) {
	t.Helper()
	if result.Status != ValidationStatusRejected ||
		result.ProofStatus != proofStatus ||
		result.PublicErrorCode != publicCode ||
		result.FailureClass != failureClass {
		t.Fatalf("result = %#v, want rejected validation with proof=%q public=%q class=%q", result, proofStatus, publicCode, failureClass)
	}
	if result.Identity.Status != app.IdentityValidationUnknown ||
		result.PlayerIDValidated ||
		result.SessionValidated ||
		result.ActorKind != "" ||
		result.ActorID != "" ||
		result.PlayerID != "" ||
		result.TokenRecordID != "" ||
		result.CredentialRecordID != "" {
		t.Fatalf("rejected validation result includes proven identity fields: %#v", result)
	}
}

func assertRejectedLogoutAccessToken(t *testing.T, result LogoutAccessTokenResult, publicCode PublicErrorCode, failureClass FailureClass) {
	t.Helper()
	if result.Status != LogoutStatusRejected ||
		result.PublicErrorCode != publicCode ||
		result.FailureClass != failureClass ||
		result.Revoked ||
		result.LogoutScope != LogoutScopeUnspecified ||
		result.TokenRecordID != "" ||
		!result.RevokedAt.IsZero() {
		t.Fatalf("result = %#v, want rejected logout with public=%q class=%q", result, publicCode, failureClass)
	}
}

func assertRejectedLocalOnboarding(t *testing.T, result LocalOnboardingDeviceCredentialIssuanceResult, publicCode PublicErrorCode, failureClass FailureClass) {
	t.Helper()
	if result.Status != LocalOnboardingDeviceCredentialIssuanceStatusRejected ||
		result.PublicErrorCode != publicCode ||
		result.FailureClass != failureClass ||
		result.PlayerID != "" ||
		result.CredentialRecordID != "" ||
		result.DeviceCredential != "" ||
		!result.CreatedAt.IsZero() {
		t.Fatalf("result = %#v, want rejected local onboarding with public=%q class=%q", result, publicCode, failureClass)
	}
}

func fmtValidationResult(result AccessTokenValidationResult) string {
	return fmt.Sprintf("%#v", result)
}

func fmtLocalOnboardingResult(result LocalOnboardingDeviceCredentialIssuanceResult) string {
	return fmt.Sprintf("%#v", result)
}

func fmtLogoutResult(result LogoutAccessTokenResult) string {
	return fmt.Sprintf("%#v", result)
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
		return fn(ctx, tx.NoopUnitOfWork{})
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
	sessionRepository        sessionmodule.Repository
	authenticationErr        error
	playerErr                error
	sessionErr               error
}

func newFakeLoginUnitOfWork(events *[]string, authenticationRepository authenticationmodule.Repository, playerRepository playermodule.Repository, sessionRepository sessionmodule.Repository) *fakeLoginUnitOfWork {
	return &fakeLoginUnitOfWork{
		ctx:                      context.Background(),
		events:                   events,
		authenticationRepository: authenticationRepository,
		playerRepository:         playerRepository,
		sessionRepository:        sessionRepository,
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

func (u *fakeLoginUnitOfWork) NewSessionRepository() (sessionmodule.Repository, error) {
	if u.events != nil {
		*u.events = append(*u.events, "new-session-repository")
	}
	if u.sessionErr != nil {
		return nil, u.sessionErr
	}
	return u.sessionRepository, nil
}

type fakeAuthenticationRepository struct {
	events                      *[]string
	credential                  authenticationmodule.CredentialRecord
	token                       authenticationmodule.TokenRecord
	findCredentialErr           error
	storeCredentialErr          error
	findTokenErr                error
	storeTokenErr               error
	revokeTokenErr              error
	lastCredentialLookupDigest  []byte
	lastTokenLookupDigest       []byte
	lastStoreCredentialMutation authenticationmodule.StoreCredentialMutation
	lastStoreTokenMutation      authenticationmodule.StoreTokenMutation
	lastRevokeTokenMutation     authenticationmodule.RevokeTokenMutation
	findCredentialCalls         int
	storeCredentialCalls        int
	findTokenCalls              int
	storeTokenCalls             int
	revokeTokenCalls            int
}

func (r *fakeAuthenticationRepository) StoreCredential(_ context.Context, mutation authenticationmodule.StoreCredentialMutation) (authenticationmodule.CredentialRecord, error) {
	r.storeCredentialCalls++
	if r.events != nil {
		*r.events = append(*r.events, "store-credential")
	}
	r.lastStoreCredentialMutation = mutation
	if r.storeCredentialErr != nil {
		return authenticationmodule.CredentialRecord{}, r.storeCredentialErr
	}
	return authenticationmodule.CredentialRecord{
		CredentialRecordID: mutation.CredentialRecordID,
		PlayerID:           mutation.PlayerID,
		CredentialKind:     mutation.CredentialKind,
		CredentialStatus:   authenticationmodule.CredentialStatusActive,
		VerifierAlgorithm:  mutation.VerifierAlgorithm,
		VerifierVersion:    mutation.VerifierVersion,
		VerifierKeyID:      mutation.VerifierKeyID,
		CreatedAt:          mutation.OccurredAt,
		UpdatedAt:          mutation.OccurredAt,
	}, nil
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

func (r *fakeAuthenticationRepository) FindTokenByLookupDigest(_ context.Context, digest []byte) (authenticationmodule.TokenRecord, error) {
	r.findTokenCalls++
	if r.events != nil {
		*r.events = append(*r.events, "find-token")
	}
	r.lastTokenLookupDigest = cloneBytes(digest)
	if r.findTokenErr != nil {
		return authenticationmodule.TokenRecord{}, r.findTokenErr
	}
	return r.token, nil
}

func (r *fakeAuthenticationRepository) RevokeCredential(context.Context, authenticationmodule.RevokeCredentialMutation) error {
	return errors.New("unexpected RevokeCredential call")
}

func (r *fakeAuthenticationRepository) RevokeToken(_ context.Context, mutation authenticationmodule.RevokeTokenMutation) error {
	r.revokeTokenCalls++
	if r.events != nil {
		*r.events = append(*r.events, "revoke-token")
	}
	r.lastRevokeTokenMutation = mutation
	if r.revokeTokenErr != nil {
		return r.revokeTokenErr
	}
	return nil
}

func (r *fakeAuthenticationRepository) ListTokensEligibleForCleanup(context.Context, authenticationmodule.TokenCleanupQuery) ([]authenticationmodule.TokenRecord, error) {
	return nil, errors.New("unexpected ListTokensEligibleForCleanup call")
}

type fakePlayerRepository struct {
	events             *[]string
	account            playermodule.Account
	createErr          error
	getErr             error
	createCalls        int
	getCalls           int
	playerIDs          []string
	lastCreateMutation playermodule.CreatePlayerAccountMutation
}

type fakeSessionRepository struct {
	events             *[]string
	createErr          error
	createCalls        int
	revokeCalls        int
	lastCreateMutation sessionmodule.CreateRuntimeSessionMutation
	lastRevokeMutation sessionmodule.RevokeRuntimeSessionMutation
}

func (r *fakeSessionRepository) CreateRuntimeSession(_ context.Context, mutation sessionmodule.CreateRuntimeSessionMutation) (sessionmodule.RuntimeSession, error) {
	r.createCalls++
	if r.events != nil {
		*r.events = append(*r.events, "create-runtime-session")
	}
	r.lastCreateMutation = mutation
	if r.createErr != nil {
		return sessionmodule.RuntimeSession{}, r.createErr
	}
	return sessionmodule.RuntimeSession{
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
	}, nil
}

func (r *fakeSessionRepository) GetRuntimeSession(context.Context, sessionmodule.GetRuntimeSessionQuery) (sessionmodule.RuntimeSession, error) {
	return sessionmodule.RuntimeSession{}, errors.New("unexpected GetRuntimeSession call")
}

func (r *fakeSessionRepository) FindActiveSessionByID(context.Context, sessionmodule.FindActiveSessionByIDQuery) (sessionmodule.RuntimeSession, error) {
	return sessionmodule.RuntimeSession{}, errors.New("unexpected FindActiveSessionByID call")
}

func (r *fakeSessionRepository) UpdateRuntimeSessionLastSeen(context.Context, sessionmodule.UpdateRuntimeSessionLastSeenMutation) (sessionmodule.RuntimeSession, error) {
	return sessionmodule.RuntimeSession{}, errors.New("unexpected UpdateRuntimeSessionLastSeen call")
}

func (r *fakeSessionRepository) MarkRuntimeSessionExpired(context.Context, sessionmodule.MarkRuntimeSessionExpiredMutation) (sessionmodule.RuntimeSession, error) {
	return sessionmodule.RuntimeSession{}, errors.New("unexpected MarkRuntimeSessionExpired call")
}

func (r *fakeSessionRepository) RevokeRuntimeSession(context.Context, sessionmodule.RevokeRuntimeSessionMutation) (sessionmodule.RuntimeSession, error) {
	r.revokeCalls++
	return sessionmodule.RuntimeSession{}, errors.New("unexpected RevokeRuntimeSession call")
}

func (r *fakeSessionRepository) ListActiveSessionsForPlayer(context.Context, sessionmodule.ListActiveSessionsForPlayerQuery) ([]sessionmodule.RuntimeSession, error) {
	return nil, errors.New("unexpected ListActiveSessionsForPlayer call")
}

func (r *fakePlayerRepository) CreatePlayerAccount(_ context.Context, mutation playermodule.CreatePlayerAccountMutation) (playermodule.Account, error) {
	r.createCalls++
	if r.events != nil {
		*r.events = append(*r.events, "create-player-account")
	}
	r.lastCreateMutation = mutation
	if r.createErr != nil {
		return playermodule.Account{}, r.createErr
	}
	return playermodule.Account{
		PlayerID:     mutation.PlayerID,
		DisplayName:  mutation.DisplayName,
		AccountState: playermodule.AccountStateActive,
		CreatedAt:    mutation.OccurredAt,
		UpdatedAt:    mutation.OccurredAt,
	}, nil
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

type recordingPlayerIDGenerator struct {
	events *[]string
	id     string
	err    error
}

func (g *recordingPlayerIDGenerator) GeneratePlayerID(context.Context) (string, error) {
	if g.events != nil {
		*g.events = append(*g.events, "generate-player-id")
	}
	if g.err != nil {
		return "", g.err
	}
	return g.id, nil
}

type staticPlayerIDGenerator string

func (g staticPlayerIDGenerator) GeneratePlayerID(context.Context) (string, error) {
	return string(g), nil
}

type recordingPlayerAccountEventIDGenerator struct {
	events *[]string
	id     string
	err    error
}

func (g *recordingPlayerAccountEventIDGenerator) GeneratePlayerAccountEventID(context.Context) (string, error) {
	if g.events != nil {
		*g.events = append(*g.events, "generate-player-account-event-id")
	}
	if g.err != nil {
		return "", g.err
	}
	return g.id, nil
}

type staticPlayerAccountEventIDGenerator string

func (g staticPlayerAccountEventIDGenerator) GeneratePlayerAccountEventID(context.Context) (string, error) {
	return string(g), nil
}

type recordingCredentialRecordIDGenerator struct {
	events *[]string
	id     string
	err    error
}

func (g *recordingCredentialRecordIDGenerator) GenerateCredentialRecordID(context.Context) (string, error) {
	if g.events != nil {
		*g.events = append(*g.events, "generate-credential-record-id")
	}
	if g.err != nil {
		return "", g.err
	}
	return g.id, nil
}

type staticCredentialRecordIDGenerator string

func (g staticCredentialRecordIDGenerator) GenerateCredentialRecordID(context.Context) (string, error) {
	return string(g), nil
}

type recordingSessionIDGenerator struct {
	events *[]string
	id     string
	err    error
	calls  int
}

func (g *recordingSessionIDGenerator) GenerateSessionID(context.Context) (string, error) {
	g.calls++
	if g.events != nil {
		*g.events = append(*g.events, "generate-session-id")
	}
	if g.err != nil {
		return "", g.err
	}
	return g.id, nil
}

type staticSessionIDGenerator string

func (g staticSessionIDGenerator) GenerateSessionID(context.Context) (string, error) {
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
