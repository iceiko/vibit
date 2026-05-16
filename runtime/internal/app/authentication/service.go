package authentication

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

type Operation string

const (
	OperationAuthenticateWithDeviceCredential Operation = "AuthenticateWithDeviceCredential"
	OperationValidateAccessToken              Operation = "ValidateAccessToken"
	OperationLogoutAccessToken                Operation = "LogoutAccessToken"
	OperationRefreshAccessToken               Operation = "RefreshAccessToken"
)

type FailureClass string

const (
	FailureClassNotImplemented         FailureClass = "not_implemented"
	FailureClassRefreshNotSupported    FailureClass = "refresh_not_supported"
	FailureClassMissingProof           FailureClass = "missing_proof"
	FailureClassMalformedProof         FailureClass = "malformed_proof"
	FailureClassLookupMiss             FailureClass = "lookup_miss"
	FailureClassWrongVerifierAlgorithm FailureClass = "wrong_verifier_algorithm"
	FailureClassUnknownVerifierKeyID   FailureClass = "unknown_verifier_key_id"
	FailureClassUnsupportedVersion     FailureClass = "unsupported_verifier_version"
	FailureClassVerifierMismatch       FailureClass = "verifier_digest_mismatch"
	FailureClassCredentialNotActive    FailureClass = "credential_not_active"
	FailureClassTokenNotActive         FailureClass = "token_not_active"
	FailureClassTokenExpired           FailureClass = "token_expired"
	FailureClassTokenRevoked           FailureClass = "token_revoked"
	FailureClassDependencyUnavailable  FailureClass = "repository_unavailable"
)

type PublicErrorCode string

const (
	PublicErrorAuthenticationProofMissing          PublicErrorCode = "AUTHENTICATION_PROOF_MISSING"
	PublicErrorAuthenticationProofMalformed        PublicErrorCode = "AUTHENTICATION_PROOF_MALFORMED"
	PublicErrorAuthenticationCredentialInvalid     PublicErrorCode = "AUTHENTICATION_CREDENTIAL_INVALID"
	PublicErrorAuthenticationTokenMissing          PublicErrorCode = "AUTHENTICATION_TOKEN_MISSING"
	PublicErrorAuthenticationTokenMalformed        PublicErrorCode = "AUTHENTICATION_TOKEN_MALFORMED"
	PublicErrorAuthenticationTokenInvalid          PublicErrorCode = "AUTHENTICATION_TOKEN_INVALID"
	PublicErrorAuthenticationCredentialUnavailable PublicErrorCode = "AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE"
	PublicErrorAuthenticationTokenUnavailable      PublicErrorCode = "AUTHENTICATION_TOKEN_STORE_UNAVAILABLE"
	PublicErrorAuthenticationRefreshNotSupported   PublicErrorCode = "AUTHENTICATION_REFRESH_NOT_SUPPORTED"
	PublicErrorAuthenticationNotImplemented        PublicErrorCode = "AUTHENTICATION_NOT_IMPLEMENTED"
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
		return fmt.Sprintf("authentication service: %s", e.PublicCode)
	}
	return fmt.Sprintf("authentication service: %s: %s", e.Operation, e.PublicCode)
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

type ServiceDependencies struct {
	UnitOfWorkRunner UnitOfWorkRunner
}

type Service struct {
	unitOfWorkRunner UnitOfWorkRunner
}

func NewService(dependencies ServiceDependencies) (Service, error) {
	if isNilUnitOfWorkRunner(dependencies.UnitOfWorkRunner) {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
		}
	}
	return Service{unitOfWorkRunner: dependencies.UnitOfWorkRunner}, nil
}

type AccountCreationIntent string

const (
	AccountCreationIntentUnspecified AccountCreationIntent = "unspecified"
	AccountCreationIntentCreate      AccountCreationIntent = "create"
	AccountCreationIntentReject      AccountCreationIntent = "reject"
)

type DeviceCredentialAuthenticationRequest struct {
	CredentialProof       string
	RequestedPlayerID     string
	ClientInstanceID      string
	AccountCreationIntent AccountCreationIntent
}

type AuthenticationStatus string

const (
	AuthenticationStatusNotImplemented AuthenticationStatus = "not_implemented"
	AuthenticationStatusRejected       AuthenticationStatus = "rejected"
	AuthenticationStatusAuthenticated  AuthenticationStatus = "authenticated"
)

type TokenType string

const (
	TokenTypeOpaqueAccess TokenType = "opaque_access_token"
)

type AuthenticationResult struct {
	Status             AuthenticationStatus
	PublicErrorCode    PublicErrorCode
	FailureClass       FailureClass
	ActorKind          app.ActorKind
	PlayerID           string
	AccessToken        string
	TokenType          TokenType
	IssuedAt           time.Time
	ExpiresAt          time.Time
	TokenRecordID      string
	CredentialRecordID string
}

func (s Service) AuthenticateWithDeviceCredential(_ context.Context, _ DeviceCredentialAuthenticationRequest) (AuthenticationResult, error) {
	_ = s.unitOfWorkRunner
	result := AuthenticationResult{
		Status:          AuthenticationStatusNotImplemented,
		PublicErrorCode: PublicErrorAuthenticationNotImplemented,
		FailureClass:    FailureClassNotImplemented,
	}
	return result, serviceNotImplemented(OperationAuthenticateWithDeviceCredential)
}

type AccessTokenValidationRequest struct {
	AccessToken     string
	Route           app.RouteKey
	ConnectionID    string
	ConnectionEpoch uint64
}

type ValidationStatus string

const (
	ValidationStatusNotImplemented ValidationStatus = "not_implemented"
	ValidationStatusRejected       ValidationStatus = "rejected"
	ValidationStatusValidated      ValidationStatus = "validated"
)

type ProofStatus string

const (
	ProofStatusNotEvaluated ProofStatus = "not_evaluated"
	ProofStatusMissing      ProofStatus = "missing"
	ProofStatusMalformed    ProofStatus = "malformed"
	ProofStatusInvalid      ProofStatus = "invalid"
	ProofStatusValid        ProofStatus = "valid"
)

type AccessTokenValidationResult struct {
	Status             ValidationStatus
	ProofStatus        ProofStatus
	PublicErrorCode    PublicErrorCode
	FailureClass       FailureClass
	Identity           app.RequestIdentity
	ActorKind          app.ActorKind
	ActorID            string
	PlayerID           string
	PlayerIDValidated  bool
	SessionValidated   bool
	TokenRecordID      string
	CredentialRecordID string
}

func (s Service) ValidateAccessToken(_ context.Context, _ AccessTokenValidationRequest) (AccessTokenValidationResult, error) {
	_ = s.unitOfWorkRunner
	result := AccessTokenValidationResult{
		Status:          ValidationStatusNotImplemented,
		ProofStatus:     ProofStatusNotEvaluated,
		PublicErrorCode: PublicErrorAuthenticationNotImplemented,
		FailureClass:    FailureClassNotImplemented,
		Identity: app.RequestIdentity{
			Status: app.IdentityValidationUnknown,
		},
	}
	return result, serviceNotImplemented(OperationValidateAccessToken)
}

type LogoutAccessTokenRequest struct {
	AccessToken  string
	LogoutReason string
}

type LogoutStatus string

const (
	LogoutStatusNotImplemented LogoutStatus = "not_implemented"
	LogoutStatusRejected       LogoutStatus = "rejected"
	LogoutStatusRevoked        LogoutStatus = "revoked"
)

type LogoutScope string

const (
	LogoutScopeUnspecified LogoutScope = "unspecified"
	LogoutScopeToken       LogoutScope = "token"
)

type LogoutAccessTokenResult struct {
	Status          LogoutStatus
	PublicErrorCode PublicErrorCode
	FailureClass    FailureClass
	Revoked         bool
	LogoutScope     LogoutScope
	TokenRecordID   string
	RevokedAt       time.Time
}

func (s Service) LogoutAccessToken(_ context.Context, _ LogoutAccessTokenRequest) (LogoutAccessTokenResult, error) {
	_ = s.unitOfWorkRunner
	result := LogoutAccessTokenResult{
		Status:          LogoutStatusNotImplemented,
		PublicErrorCode: PublicErrorAuthenticationNotImplemented,
		FailureClass:    FailureClassNotImplemented,
	}
	return result, serviceNotImplemented(OperationLogoutAccessToken)
}

type RefreshAccessTokenRequest struct {
	AccessToken string
}

type RefreshStatus string

const (
	RefreshStatusUnsupported RefreshStatus = "unsupported"
	RefreshStatusRefreshed   RefreshStatus = "refreshed"
)

type RefreshAccessTokenResult struct {
	Status          RefreshStatus
	PublicErrorCode PublicErrorCode
	FailureClass    FailureClass
	AccessToken     string
	TokenType       TokenType
	IssuedAt        time.Time
	ExpiresAt       time.Time
	TokenRecordID   string
}

func (s Service) RefreshAccessToken(_ context.Context, _ RefreshAccessTokenRequest) (RefreshAccessTokenResult, error) {
	_ = s.unitOfWorkRunner
	result := RefreshAccessTokenResult{
		Status:          RefreshStatusUnsupported,
		PublicErrorCode: PublicErrorAuthenticationRefreshNotSupported,
		FailureClass:    FailureClassRefreshNotSupported,
	}
	return result, &ServiceError{
		Operation:  OperationRefreshAccessToken,
		Class:      FailureClassRefreshNotSupported,
		PublicCode: PublicErrorAuthenticationRefreshNotSupported,
	}
}

func serviceNotImplemented(operation Operation) error {
	return &ServiceError{
		Operation:  operation,
		Class:      FailureClassNotImplemented,
		PublicCode: PublicErrorAuthenticationNotImplemented,
	}
}

func isNilUnitOfWorkRunner(runner UnitOfWorkRunner) bool {
	if runner == nil {
		return true
	}
	value := reflect.ValueOf(runner)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
