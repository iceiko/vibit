package authentication

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	EnvVerifierKeySetID      = "VIBIT_AUTH_VERIFIER_KEY_SET_ID"
	EnvCredentialLookupKey   = "VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY"
	EnvCredentialVerifierKey = "VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY"
	EnvTokenLookupKey        = "VIBIT_AUTH_TOKEN_LOOKUP_KEY"
	EnvTokenVerifierKey      = "VIBIT_AUTH_TOKEN_VERIFIER_KEY"
)

var (
	ErrMissingEnvironmentLookup   = errors.New("authentication verifier key environment: missing lookup")
	ErrMissingEnvironmentVariable = errors.New("authentication verifier key environment: missing variable")
	ErrInvalidEnvironmentKeyText  = errors.New("authentication verifier key environment: invalid key encoding")
	ErrInvalidEnvironmentKeySet   = errors.New("authentication verifier key environment: invalid key set")
)

type EnvLookup func(name string) (string, bool)

type VerifierKeyEnvironmentError struct {
	Kind     error
	Variable string
	Purpose  KeyPurpose
	Err      error
}

func (e *VerifierKeyEnvironmentError) Error() string {
	if e == nil {
		return ""
	}
	parts := []string{e.Kind.Error()}
	if e.Variable != "" {
		parts = append(parts, e.Variable)
	}
	if e.Purpose != "" {
		parts = append(parts, string(e.Purpose))
	}
	return strings.Join(parts, ": ")
}

func (e *VerifierKeyEnvironmentError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *VerifierKeyEnvironmentError) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.Kind, target) || errors.Is(e.Err, target)
}

func LoadVerifierKeySetFromEnvironment(lookup EnvLookup) (VerifierKeySet, error) {
	if lookup == nil {
		return VerifierKeySet{}, &VerifierKeyEnvironmentError{Kind: ErrMissingEnvironmentLookup}
	}

	keySetID, ok := lookup(EnvVerifierKeySetID)
	if !ok {
		return VerifierKeySet{}, missingEnvironmentVariable(EnvVerifierKeySetID)
	}

	config := VerifierKeySetConfig{
		KeySetID: strings.TrimSpace(keySetID),
	}

	for _, variable := range verifierKeyEnvironmentVariables() {
		value, ok := lookup(variable.name)
		if !ok {
			return VerifierKeySet{}, missingEnvironmentVariable(variable.name)
		}
		decoded, err := decodeVerifierKeyEnvironmentValue(value)
		if err != nil {
			return VerifierKeySet{}, &VerifierKeyEnvironmentError{
				Kind:     ErrInvalidEnvironmentKeyText,
				Variable: variable.name,
				Purpose:  variable.purpose,
				Err:      err,
			}
		}
		variable.assign(&config, decoded)
	}

	keySet, err := NewVerifierKeySet(config)
	if err != nil {
		return VerifierKeySet{}, environmentValidatorError(err)
	}
	return keySet, nil
}

func LoadVerifierKeySetFromProcessEnvironment() (VerifierKeySet, error) {
	return LoadVerifierKeySetFromEnvironment(os.LookupEnv)
}

type verifierKeyEnvironmentVariable struct {
	name    string
	purpose KeyPurpose
	assign  func(*VerifierKeySetConfig, []byte)
}

func verifierKeyEnvironmentVariables() []verifierKeyEnvironmentVariable {
	return []verifierKeyEnvironmentVariable{
		{
			name:    EnvCredentialLookupKey,
			purpose: KeyPurposeCredentialLookup,
			assign: func(config *VerifierKeySetConfig, value []byte) {
				config.CredentialLookupKey = value
			},
		},
		{
			name:    EnvCredentialVerifierKey,
			purpose: KeyPurposeCredentialVerifier,
			assign: func(config *VerifierKeySetConfig, value []byte) {
				config.CredentialVerifierKey = value
			},
		},
		{
			name:    EnvTokenLookupKey,
			purpose: KeyPurposeTokenLookup,
			assign: func(config *VerifierKeySetConfig, value []byte) {
				config.TokenLookupKey = value
			},
		},
		{
			name:    EnvTokenVerifierKey,
			purpose: KeyPurposeTokenVerifier,
			assign: func(config *VerifierKeySetConfig, value []byte) {
				config.TokenVerifierKey = value
			},
		},
	}
}

func decodeVerifierKeyEnvironmentValue(value string) ([]byte, error) {
	text := strings.TrimSpace(value)
	if decoded, err := base64.RawURLEncoding.DecodeString(text); err == nil {
		return decoded, nil
	}
	if decoded, err := base64.StdEncoding.DecodeString(text); err == nil {
		return decoded, nil
	}
	return nil, fmt.Errorf("base64 decode failed")
}

func missingEnvironmentVariable(variable string) error {
	return &VerifierKeyEnvironmentError{
		Kind:     ErrMissingEnvironmentVariable,
		Variable: variable,
	}
}

func environmentValidatorError(err error) error {
	envErr := &VerifierKeyEnvironmentError{
		Kind: ErrInvalidEnvironmentKeySet,
		Err:  err,
	}

	if errors.Is(err, ErrMissingKeySetID) {
		envErr.Variable = EnvVerifierKeySetID
		return envErr
	}

	var configErr *VerifierKeyConfigError
	if errors.As(err, &configErr) {
		envErr.Purpose = configErr.Purpose
		envErr.Variable = environmentVariableForPurpose(configErr.Purpose)
	}

	return envErr
}

func environmentVariableForPurpose(purpose KeyPurpose) string {
	switch purpose {
	case KeyPurposeCredentialLookup:
		return EnvCredentialLookupKey
	case KeyPurposeCredentialVerifier:
		return EnvCredentialVerifierKey
	case KeyPurposeTokenLookup:
		return EnvTokenLookupKey
	case KeyPurposeTokenVerifier:
		return EnvTokenVerifierKey
	default:
		return ""
	}
}
