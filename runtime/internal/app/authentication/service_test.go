package authentication

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
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

func TestAuthenticateWithDeviceCredentialFailsClosedWithoutRepositoryCall(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{}
	service := mustNewService(t, runner)
	proof := "raw-device-credential-proof"

	result, err := service.AuthenticateWithDeviceCredential(context.Background(), DeviceCredentialAuthenticationRequest{
		CredentialProof:       proof,
		RequestedPlayerID:     "player-1",
		ClientInstanceID:      "client-1",
		AccountCreationIntent: AccountCreationIntentCreate,
	})
	assertServiceError(t, err, OperationAuthenticateWithDeviceCredential, FailureClassNotImplemented, PublicErrorAuthenticationNotImplemented)
	assertNoSecretLeak(t, err, proof)
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if result.Status != AuthenticationStatusNotImplemented {
		t.Fatalf("Status = %q, want %q", result.Status, AuthenticationStatusNotImplemented)
	}
	if result.PublicErrorCode != PublicErrorAuthenticationNotImplemented {
		t.Fatalf("PublicErrorCode = %q, want %q", result.PublicErrorCode, PublicErrorAuthenticationNotImplemented)
	}
	if result.FailureClass != FailureClassNotImplemented {
		t.Fatalf("FailureClass = %q, want %q", result.FailureClass, FailureClassNotImplemented)
	}
	if result.AccessToken != "" || result.TokenRecordID != "" || result.CredentialRecordID != "" || result.PlayerID != "" {
		t.Fatalf("result includes behavior-owned fields: %#v", result)
	}
	if result.ActorKind != "" {
		t.Fatalf("ActorKind = %q, want empty until login executes", result.ActorKind)
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

func mustNewService(t *testing.T, runner UnitOfWorkRunner) Service {
	t.Helper()
	service, err := NewService(ServiceDependencies{UnitOfWorkRunner: runner})
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
