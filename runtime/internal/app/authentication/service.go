package authentication

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/session"
	authenticationmodule "github.com/iceiko/vibit/runtime/internal/modules/authentication"
	playermodule "github.com/iceiko/vibit/runtime/internal/modules/player"
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
	FailureClassWrongCredentialKind    FailureClass = "wrong_credential_kind"
	FailureClassWrongVerifierAlgorithm FailureClass = "wrong_verifier_algorithm"
	FailureClassUnknownVerifierKeyID   FailureClass = "unknown_verifier_key_id"
	FailureClassUnsupportedVersion     FailureClass = "unsupported_verifier_version"
	FailureClassVerifierMismatch       FailureClass = "verifier_digest_mismatch"
	FailureClassCredentialNotActive    FailureClass = "credential_not_active"
	FailureClassPlayerAccountNotActive FailureClass = "player_account_not_active"
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

const (
	credentialKindDeviceCredentialLogin = "device_credential_login"
	tokenKindAccessToken                = "access_token"
	verifierAlgorithmHMACSHA256V1       = "vibit_hmac_sha256_v1"
	verifierVersionV1                   = 1
	defaultRequestedBy                  = "authentication_service"
	logoutReasonPresentedAccessToken    = "logout_presented_access_token"
)

var (
	errMissingDeviceCredentialProof   = errors.New("authentication service: missing device credential proof")
	errMalformedDeviceCredentialProof = errors.New("authentication service: malformed device credential proof")
	errMissingAuthenticationUOW       = errors.New("authentication service: authentication unit-of-work capability is required")
	errMissingDependency              = errors.New("authentication service: dependency is required")
	errInvalidDependency              = errors.New("authentication service: dependency is invalid")
	errMalformedRuntimeSessionID      = errors.New("authentication service: generated runtime session id is malformed")
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

type Clock interface {
	Now() time.Time
}

type TokenRecordIDGenerator interface {
	GenerateTokenRecordID(context.Context) (string, error)
}

type SessionIDGenerator interface {
	GenerateSessionID(context.Context) (string, error)
}

type ServiceDependencies struct {
	UnitOfWorkRunner       UnitOfWorkRunner
	VerifierKeySet         VerifierKeySet
	AccessTokenRandom      io.Reader
	Clock                  Clock
	TokenRecordIDGenerator TokenRecordIDGenerator
	SessionIDGenerator     SessionIDGenerator
	AccessTokenLifetime    time.Duration
	TokenAudience          string
}

type Service struct {
	unitOfWorkRunner       UnitOfWorkRunner
	verifierKeySet         VerifierKeySet
	accessTokenRandom      io.Reader
	clock                  Clock
	tokenRecordIDGenerator TokenRecordIDGenerator
	sessionIDGenerator     SessionIDGenerator
	accessTokenLifetime    time.Duration
	tokenAudience          string
}

func NewService(dependencies ServiceDependencies) (Service, error) {
	if isNilUnitOfWorkRunner(dependencies.UnitOfWorkRunner) {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
		}
	}
	if isZeroVerifierKeySet(dependencies.VerifierKeySet) {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
			Err:        errMissingDependency,
		}
	}
	if isNilInterface(dependencies.AccessTokenRandom) {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
			Err:        errMissingDependency,
		}
	}
	if isNilInterface(dependencies.Clock) {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
			Err:        errMissingDependency,
		}
	}
	if isNilInterface(dependencies.TokenRecordIDGenerator) {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
			Err:        errMissingDependency,
		}
	}
	if isNilInterface(dependencies.SessionIDGenerator) {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
			Err:        errMissingDependency,
		}
	}
	if dependencies.AccessTokenLifetime <= 0 {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
			Err:        errInvalidDependency,
		}
	}
	tokenAudience := strings.TrimSpace(dependencies.TokenAudience)
	if tokenAudience == "" {
		return Service{}, &ServiceError{
			Operation:  "NewService",
			Class:      FailureClassDependencyUnavailable,
			PublicCode: PublicErrorAuthenticationCredentialUnavailable,
			Err:        errMissingDependency,
		}
	}
	return Service{
		unitOfWorkRunner:       dependencies.UnitOfWorkRunner,
		verifierKeySet:         dependencies.VerifierKeySet,
		accessTokenRandom:      dependencies.AccessTokenRandom,
		clock:                  dependencies.Clock,
		tokenRecordIDGenerator: dependencies.TokenRecordIDGenerator,
		sessionIDGenerator:     dependencies.SessionIDGenerator,
		accessTokenLifetime:    dependencies.AccessTokenLifetime,
		tokenAudience:          tokenAudience,
	}, nil
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
	SessionID          string
	SessionCreated     bool
	SessionExpiresAt   time.Time
	TokenRecordID      string
	CredentialRecordID string
}

func (s Service) AuthenticateWithDeviceCredential(ctx context.Context, request DeviceCredentialAuthenticationRequest) (AuthenticationResult, error) {
	rawProof, err := decodeDeviceCredentialProof(request.CredentialProof)
	if err != nil {
		if errors.Is(err, errMissingDeviceCredentialProof) {
			return rejectedAuthenticationResult(PublicErrorAuthenticationProofMissing, FailureClassMissingProof),
				serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassMissingProof, PublicErrorAuthenticationProofMissing, err)
		}
		return rejectedAuthenticationResult(PublicErrorAuthenticationProofMalformed, FailureClassMalformedProof),
			serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassMalformedProof, PublicErrorAuthenticationProofMalformed, err)
	}
	defer zeroBytes(rawProof)

	credentialLookupDigest, err := ComputeCredentialLookupDigest(s.verifierKeySet, rawProof)
	if err != nil {
		return rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
	}

	var committedResult AuthenticationResult
	var failedResult AuthenticationResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repositories, ok := unit.(authenticationLoginUnitOfWork)
		if !ok {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, errMissingAuthenticationUOW)
		}

		authenticationRepository, err := repositories.NewAuthenticationRepository()
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}
		playerRepository, err := repositories.NewPlayerAccountRepository()
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}

		credential, err := authenticationRepository.FindCredentialByLookupDigest(runCtx, credentialLookupDigest.Bytes())
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialInvalid, FailureClassLookupMiss)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassLookupMiss, PublicErrorAuthenticationCredentialInvalid, err)
		}
		if failureClass, ok := s.rejectCredentialRecord(credential); ok {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialInvalid, failureClass)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, failureClass, PublicErrorAuthenticationCredentialInvalid, nil)
		}

		credentialVerifierDigest, err := ComputeCredentialVerifierDigest(s.verifierKeySet, rawProof)
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialInvalid, FailureClassVerifierMismatch)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassVerifierMismatch, PublicErrorAuthenticationCredentialInvalid, err)
		}
		comparison, err := CompareCredentialVerifierDigest(credentialVerifierDigest, credential.CredentialVerifierDigest)
		if err != nil || !comparison.Matched() {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialInvalid, FailureClassVerifierMismatch)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassVerifierMismatch, PublicErrorAuthenticationCredentialInvalid, err)
		}

		account, err := playerRepository.GetPlayerAccount(runCtx, credential.PlayerID)
		if err != nil || account.AccountState != playermodule.AccountStateActive {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialInvalid, FailureClassPlayerAccountNotActive)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassPlayerAccountNotActive, PublicErrorAuthenticationCredentialInvalid, err)
		}

		accessTokenMaterial, err := GenerateAccessTokenMaterial(s.accessTokenRandom)
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}
		rawAccessToken := accessTokenMaterial.RawBytes()
		defer zeroBytes(rawAccessToken)

		tokenLookupDigest, err := ComputeTokenLookupDigest(s.verifierKeySet, rawAccessToken)
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}
		tokenVerifierDigest, err := ComputeTokenVerifierDigest(s.verifierKeySet, rawAccessToken)
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}

		tokenRecordID, err := s.tokenRecordIDGenerator.GenerateTokenRecordID(runCtx)
		tokenRecordID = strings.TrimSpace(tokenRecordID)
		if err != nil || tokenRecordID == "" {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}

		issuedAt := s.clock.Now().UTC()
		if issuedAt.IsZero() {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, errInvalidDependency)
		}
		expiresAt := issuedAt.Add(s.accessTokenLifetime).UTC()

		tokenRecord, err := authenticationRepository.StoreToken(runCtx, authenticationmodule.StoreTokenMutation{
			TokenRecordID:       tokenRecordID,
			PlayerID:            credential.PlayerID,
			CredentialRecordID:  credential.CredentialRecordID,
			TokenKind:           tokenKindAccessToken,
			ActorKind:           string(app.ActorKindPlayer),
			TokenLookupDigest:   tokenLookupDigest.Bytes(),
			TokenVerifierDigest: tokenVerifierDigest.Bytes(),
			VerifierAlgorithm:   verifierAlgorithmHMACSHA256V1,
			VerifierVersion:     verifierVersionV1,
			VerifierKeyID:       s.verifierKeySet.KeySetID(),
			Audience:            s.tokenAudience,
			IssuedAt:            issuedAt,
			ExpiresAt:           expiresAt,
			RequestedBy:         defaultRequestedBy,
		})
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}

		sessionRepository, err := repositories.NewSessionRepository()
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}

		sessionID, err := s.generatedRuntimeSessionID(runCtx)
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}

		runtimeSession, err := sessionRepository.CreateRuntimeSession(runCtx, session.CreateRuntimeSessionMutation{
			SessionID:           sessionID,
			ActorKind:           session.ActorKindPlayer,
			ActorID:             credential.PlayerID,
			PlayerID:            credential.PlayerID,
			SessionStatus:       session.SessionStatusActive,
			IssuedAt:            issuedAt,
			ExpiresAt:           expiresAt,
			LastSeenAt:          issuedAt,
			AccessTokenRecordID: tokenRecord.TokenRecordID,
			RequestedBy:         defaultRequestedBy,
		})
		if err != nil {
			failedResult = rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
		}

		committedResult = AuthenticationResult{
			Status:             AuthenticationStatusAuthenticated,
			ActorKind:          app.ActorKindPlayer,
			PlayerID:           credential.PlayerID,
			AccessToken:        accessTokenMaterial.Text(),
			TokenType:          TokenTypeOpaqueAccess,
			IssuedAt:           issuedAt,
			ExpiresAt:          expiresAt,
			SessionID:          runtimeSession.SessionID,
			SessionCreated:     true,
			SessionExpiresAt:   runtimeSession.ExpiresAt,
			TokenRecordID:      tokenRecord.TokenRecordID,
			CredentialRecordID: credential.CredentialRecordID,
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedAuthenticationResult(PublicErrorAuthenticationCredentialUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationAuthenticateWithDeviceCredential, FailureClassDependencyUnavailable, PublicErrorAuthenticationCredentialUnavailable, err)
	}
	return committedResult, nil
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
	ProofStatusExpired      ProofStatus = "expired"
	ProofStatusRevoked      ProofStatus = "revoked"
	ProofStatusUnavailable  ProofStatus = "unavailable"
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

const tokenIssuedAtClockTolerance time.Duration = 0

func (s Service) ValidateAccessToken(ctx context.Context, request AccessTokenValidationRequest) (AccessTokenValidationResult, error) {
	rawToken, err := decodeAccessTokenProof(request.AccessToken)
	if err != nil {
		if errors.Is(err, errMissingAccessTokenProof) {
			return rejectedAccessTokenValidationResult(ProofStatusMissing, PublicErrorAuthenticationTokenMissing, FailureClassMissingProof),
				serviceFailure(OperationValidateAccessToken, FailureClassMissingProof, PublicErrorAuthenticationTokenMissing, err)
		}
		return rejectedAccessTokenValidationResult(ProofStatusMalformed, PublicErrorAuthenticationTokenMalformed, FailureClassMalformedProof),
			serviceFailure(OperationValidateAccessToken, FailureClassMalformedProof, PublicErrorAuthenticationTokenMalformed, err)
	}
	defer zeroBytes(rawToken)

	tokenLookupDigest, err := ComputeTokenLookupDigest(s.verifierKeySet, rawToken)
	if err != nil {
		return rejectedAccessTokenValidationResult(ProofStatusUnavailable, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationValidateAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
	}

	var committedResult AccessTokenValidationResult
	var failedResult AccessTokenValidationResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repositories, ok := unit.(authenticationValidationUnitOfWork)
		if !ok {
			failedResult = rejectedAccessTokenValidationResult(ProofStatusUnavailable, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationValidateAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, errMissingAuthenticationUOW)
		}

		authenticationRepository, err := repositories.NewAuthenticationRepository()
		if err != nil {
			failedResult = rejectedAccessTokenValidationResult(ProofStatusUnavailable, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationValidateAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
		}
		playerRepository, err := repositories.NewPlayerAccountRepository()
		if err != nil {
			failedResult = rejectedAccessTokenValidationResult(ProofStatusUnavailable, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationValidateAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
		}

		token, err := authenticationRepository.FindTokenByLookupDigest(runCtx, tokenLookupDigest.Bytes())
		if err != nil {
			failedResult = rejectedAccessTokenValidationResult(ProofStatusInvalid, PublicErrorAuthenticationTokenInvalid, FailureClassLookupMiss)
			return serviceFailure(OperationValidateAccessToken, FailureClassLookupMiss, PublicErrorAuthenticationTokenInvalid, err)
		}
		if failureClass, ok := s.rejectTokenRecord(token); ok {
			failedResult = rejectedAccessTokenValidationResult(validationProofStatusForTokenFailure(failureClass), PublicErrorAuthenticationTokenInvalid, failureClass)
			return serviceFailure(OperationValidateAccessToken, failureClass, PublicErrorAuthenticationTokenInvalid, nil)
		}

		tokenVerifierDigest, err := ComputeTokenVerifierDigest(s.verifierKeySet, rawToken)
		if err != nil {
			failedResult = rejectedAccessTokenValidationResult(ProofStatusInvalid, PublicErrorAuthenticationTokenInvalid, FailureClassVerifierMismatch)
			return serviceFailure(OperationValidateAccessToken, FailureClassVerifierMismatch, PublicErrorAuthenticationTokenInvalid, err)
		}
		comparison, err := CompareTokenVerifierDigest(tokenVerifierDigest, token.TokenVerifierDigest)
		if err != nil || !comparison.Matched() {
			failedResult = rejectedAccessTokenValidationResult(ProofStatusInvalid, PublicErrorAuthenticationTokenInvalid, FailureClassVerifierMismatch)
			return serviceFailure(OperationValidateAccessToken, FailureClassVerifierMismatch, PublicErrorAuthenticationTokenInvalid, err)
		}

		account, err := playerRepository.GetPlayerAccount(runCtx, token.PlayerID)
		if err != nil || account.AccountState != playermodule.AccountStateActive {
			failedResult = rejectedAccessTokenValidationResult(ProofStatusInvalid, PublicErrorAuthenticationTokenInvalid, FailureClassPlayerAccountNotActive)
			return serviceFailure(OperationValidateAccessToken, FailureClassPlayerAccountNotActive, PublicErrorAuthenticationTokenInvalid, err)
		}

		identity := app.ValidatedPlayerIdentity(token.PlayerID, app.Session{
			ConnectionID:    request.ConnectionID,
			ConnectionEpoch: request.ConnectionEpoch,
			PlayerID:        token.PlayerID,
		})
		identity.SessionValidated = false
		identity.PlayerIDValidated = true
		identity.ActorKind = app.ActorKindPlayer
		identity.ActorID = token.PlayerID
		identity.PlayerID = token.PlayerID

		committedResult = AccessTokenValidationResult{
			Status:             ValidationStatusValidated,
			ProofStatus:        ProofStatusValid,
			PublicErrorCode:    "",
			FailureClass:       "",
			Identity:           identity,
			ActorKind:          app.ActorKindPlayer,
			ActorID:            token.PlayerID,
			PlayerID:           token.PlayerID,
			PlayerIDValidated:  true,
			SessionValidated:   false,
			TokenRecordID:      token.TokenRecordID,
			CredentialRecordID: token.CredentialRecordID,
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedAccessTokenValidationResult(ProofStatusUnavailable, PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationValidateAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
	}
	return committedResult, nil
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

func (s Service) LogoutAccessToken(ctx context.Context, request LogoutAccessTokenRequest) (LogoutAccessTokenResult, error) {
	rawToken, err := decodeAccessTokenProof(request.AccessToken)
	if err != nil {
		if errors.Is(err, errMissingAccessTokenProof) {
			return rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenMissing, FailureClassMissingProof),
				serviceFailure(OperationLogoutAccessToken, FailureClassMissingProof, PublicErrorAuthenticationTokenMissing, err)
		}
		return rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenMalformed, FailureClassMalformedProof),
			serviceFailure(OperationLogoutAccessToken, FailureClassMalformedProof, PublicErrorAuthenticationTokenMalformed, err)
	}
	defer zeroBytes(rawToken)

	tokenLookupDigest, err := ComputeTokenLookupDigest(s.verifierKeySet, rawToken)
	if err != nil {
		return rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationLogoutAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
	}

	var committedResult LogoutAccessTokenResult
	var failedResult LogoutAccessTokenResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repositories, ok := unit.(authenticationLogoutUnitOfWork)
		if !ok {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationLogoutAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, errMissingAuthenticationUOW)
		}

		authenticationRepository, err := repositories.NewAuthenticationRepository()
		if err != nil {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationLogoutAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
		}

		token, err := authenticationRepository.FindTokenByLookupDigest(runCtx, tokenLookupDigest.Bytes())
		if err != nil {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenInvalid, FailureClassLookupMiss)
			return serviceFailure(OperationLogoutAccessToken, FailureClassLookupMiss, PublicErrorAuthenticationTokenInvalid, err)
		}
		if failureClass, ok := s.rejectTokenRecord(token); ok {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenInvalid, failureClass)
			return serviceFailure(OperationLogoutAccessToken, failureClass, PublicErrorAuthenticationTokenInvalid, nil)
		}

		tokenVerifierDigest, err := ComputeTokenVerifierDigest(s.verifierKeySet, rawToken)
		if err != nil {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenInvalid, FailureClassVerifierMismatch)
			return serviceFailure(OperationLogoutAccessToken, FailureClassVerifierMismatch, PublicErrorAuthenticationTokenInvalid, err)
		}
		comparison, err := CompareTokenVerifierDigest(tokenVerifierDigest, token.TokenVerifierDigest)
		if err != nil || !comparison.Matched() {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenInvalid, FailureClassVerifierMismatch)
			return serviceFailure(OperationLogoutAccessToken, FailureClassVerifierMismatch, PublicErrorAuthenticationTokenInvalid, err)
		}

		revokedAt := s.clock.Now().UTC()
		if revokedAt.IsZero() {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationLogoutAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, errInvalidDependency)
		}

		err = authenticationRepository.RevokeToken(runCtx, authenticationmodule.RevokeTokenMutation{
			TokenRecordID: token.TokenRecordID,
			RevokedAt:     revokedAt,
			RevokedReason: logoutReasonPresentedAccessToken,
			RequestedBy:   defaultRequestedBy,
		})
		if err != nil {
			failedResult = rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationLogoutAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
		}

		committedResult = LogoutAccessTokenResult{
			Status:        LogoutStatusRevoked,
			Revoked:       true,
			LogoutScope:   LogoutScopeToken,
			TokenRecordID: token.TokenRecordID,
			RevokedAt:     revokedAt,
		}
		return nil
	})
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedLogoutAccessTokenResult(PublicErrorAuthenticationTokenUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationLogoutAccessToken, FailureClassDependencyUnavailable, PublicErrorAuthenticationTokenUnavailable, err)
	}
	return committedResult, nil
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

type authenticationLoginUnitOfWork interface {
	NewAuthenticationRepository() (authenticationmodule.Repository, error)
	NewPlayerAccountRepository() (playermodule.Repository, error)
	NewSessionRepository() (session.Repository, error)
}

type authenticationValidationUnitOfWork interface {
	NewAuthenticationRepository() (authenticationmodule.Repository, error)
	NewPlayerAccountRepository() (playermodule.Repository, error)
}

type authenticationLogoutUnitOfWork interface {
	NewAuthenticationRepository() (authenticationmodule.Repository, error)
}

func decodeDeviceCredentialProof(proof string) ([]byte, error) {
	if strings.TrimSpace(proof) == "" {
		return nil, errMissingDeviceCredentialProof
	}
	if strings.TrimSpace(proof) != proof {
		return nil, errMalformedDeviceCredentialProof
	}
	if len(proof) != EncodedSecretMaterialChars {
		return nil, errMalformedDeviceCredentialProof
	}
	raw, err := base64.RawURLEncoding.DecodeString(proof)
	if err != nil {
		return nil, errMalformedDeviceCredentialProof
	}
	if len(raw) != RawSecretMaterialBytes {
		return nil, errMalformedDeviceCredentialProof
	}
	return raw, nil
}

var errMissingAccessTokenProof = errors.New("authentication service: missing access token proof")
var errMalformedAccessTokenProof = errors.New("authentication service: malformed access token proof")

func decodeAccessTokenProof(proof string) ([]byte, error) {
	if strings.TrimSpace(proof) == "" {
		return nil, errMissingAccessTokenProof
	}
	if strings.TrimSpace(proof) != proof {
		return nil, errMalformedAccessTokenProof
	}
	if len(proof) != EncodedSecretMaterialChars {
		return nil, errMalformedAccessTokenProof
	}
	raw, err := base64.RawURLEncoding.DecodeString(proof)
	if err != nil {
		return nil, errMalformedAccessTokenProof
	}
	if len(raw) != RawSecretMaterialBytes {
		return nil, errMalformedAccessTokenProof
	}
	return raw, nil
}

func rejectedAccessTokenValidationResult(proofStatus ProofStatus, publicCode PublicErrorCode, failureClass FailureClass) AccessTokenValidationResult {
	return AccessTokenValidationResult{
		Status:          ValidationStatusRejected,
		ProofStatus:     proofStatus,
		PublicErrorCode: publicCode,
		FailureClass:    failureClass,
		Identity: app.RequestIdentity{
			Status: app.IdentityValidationUnknown,
		},
	}
}

func validationProofStatusForTokenFailure(failureClass FailureClass) ProofStatus {
	switch failureClass {
	case FailureClassTokenExpired:
		return ProofStatusExpired
	case FailureClassTokenRevoked:
		return ProofStatusRevoked
	default:
		return ProofStatusInvalid
	}
}

func rejectedAuthenticationResult(publicCode PublicErrorCode, failureClass FailureClass) AuthenticationResult {
	return AuthenticationResult{
		Status:          AuthenticationStatusRejected,
		PublicErrorCode: publicCode,
		FailureClass:    failureClass,
	}
}

func rejectedLogoutAccessTokenResult(publicCode PublicErrorCode, failureClass FailureClass) LogoutAccessTokenResult {
	return LogoutAccessTokenResult{
		Status:          LogoutStatusRejected,
		PublicErrorCode: publicCode,
		FailureClass:    failureClass,
		LogoutScope:     LogoutScopeUnspecified,
	}
}

func (s Service) generatedRuntimeSessionID(ctx context.Context) (string, error) {
	sessionID, err := s.sessionIDGenerator.GenerateSessionID(ctx)
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(sessionID)
	if trimmed == "" || trimmed != sessionID {
		return "", errMalformedRuntimeSessionID
	}
	return trimmed, nil
}

func serviceFailure(operation Operation, class FailureClass, publicCode PublicErrorCode, err error) error {
	return &ServiceError{
		Operation:  operation,
		Class:      class,
		PublicCode: publicCode,
		Err:        err,
	}
}

func (s Service) rejectCredentialRecord(record authenticationmodule.CredentialRecord) (FailureClass, bool) {
	if record.CredentialKind != credentialKindDeviceCredentialLogin {
		return FailureClassWrongCredentialKind, true
	}
	if record.CredentialStatus != authenticationmodule.CredentialStatusActive {
		return FailureClassCredentialNotActive, true
	}
	if record.VerifierAlgorithm != verifierAlgorithmHMACSHA256V1 {
		return FailureClassWrongVerifierAlgorithm, true
	}
	if record.VerifierVersion != verifierVersionV1 {
		return FailureClassUnsupportedVersion, true
	}
	if record.VerifierKeyID != s.verifierKeySet.KeySetID() {
		return FailureClassUnknownVerifierKeyID, true
	}
	return "", false
}

func (s Service) rejectTokenRecord(record authenticationmodule.TokenRecord) (FailureClass, bool) {
	if record.TokenKind != tokenKindAccessToken {
		return FailureClassTokenNotActive, true
	}
	if record.TokenStatus != authenticationmodule.TokenStatusActive {
		switch record.TokenStatus {
		case authenticationmodule.TokenStatusExpired:
			return FailureClassTokenExpired, true
		case authenticationmodule.TokenStatusRevoked:
			return FailureClassTokenRevoked, true
		default:
			return FailureClassTokenNotActive, true
		}
	}
	if record.VerifierAlgorithm != verifierAlgorithmHMACSHA256V1 {
		return FailureClassWrongVerifierAlgorithm, true
	}
	if record.VerifierVersion != verifierVersionV1 {
		return FailureClassUnsupportedVersion, true
	}
	if record.VerifierKeyID != s.verifierKeySet.KeySetID() {
		return FailureClassUnknownVerifierKeyID, true
	}
	if strings.TrimSpace(record.Audience) != s.tokenAudience {
		return FailureClassTokenNotActive, true
	}
	now := s.clock.Now().UTC()
	if record.IssuedAt.IsZero() {
		return FailureClassTokenNotActive, true
	}
	if record.IssuedAt.After(now.Add(tokenIssuedAtClockTolerance)) {
		return FailureClassTokenNotActive, true
	}
	if !record.ExpiresAt.IsZero() && now.After(record.ExpiresAt) {
		return FailureClassTokenExpired, true
	}
	return "", false
}

func zeroBytes(value []byte) {
	for i := range value {
		value[i] = 0
	}
}

func isZeroVerifierKeySet(keySet VerifierKeySet) bool {
	return keySet.KeySetID() == "" ||
		len(keySet.CredentialLookupKey()) == 0 ||
		len(keySet.CredentialVerifierKey()) == 0 ||
		len(keySet.TokenLookupKey()) == 0 ||
		len(keySet.TokenVerifierKey()) == 0
}

func isNilUnitOfWorkRunner(runner UnitOfWorkRunner) bool {
	return isNilInterface(runner)
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
