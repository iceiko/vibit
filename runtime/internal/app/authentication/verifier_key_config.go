package authentication

import (
	"errors"
	"fmt"
	"strings"
)

const MinVerifierKeyBytes = 32

type KeyPurpose string

const (
	KeyPurposeCredentialLookup   KeyPurpose = "credential_lookup"
	KeyPurposeCredentialVerifier KeyPurpose = "credential_verifier"
	KeyPurposeTokenLookup        KeyPurpose = "token_lookup"
	KeyPurposeTokenVerifier      KeyPurpose = "token_verifier"
)

var (
	ErrMissingKeySetID = errors.New("authentication verifier key config: missing key set id")
	ErrMissingKey      = errors.New("authentication verifier key config: missing required key")
	ErrKeyTooShort     = errors.New("authentication verifier key config: key too short")
	ErrDuplicateKey    = errors.New("authentication verifier key config: duplicate logical key")
	ErrWeakRepeatedKey = errors.New("authentication verifier key config: weak repeated key")
)

type VerifierKeySetConfig struct {
	KeySetID              string
	CredentialLookupKey   []byte
	CredentialVerifierKey []byte
	TokenLookupKey        []byte
	TokenVerifierKey      []byte
}

type VerifierKeyConfigError struct {
	Kind    error
	Purpose KeyPurpose
}

func (e *VerifierKeyConfigError) Error() string {
	if e == nil {
		return ""
	}
	if e.Purpose == "" {
		return e.Kind.Error()
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Purpose)
}

func (e *VerifierKeyConfigError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Kind
}

type VerifierKeySet struct {
	keySetID              string
	credentialLookupKey   []byte
	credentialVerifierKey []byte
	tokenLookupKey        []byte
	tokenVerifierKey      []byte
}

func NewVerifierKeySet(config VerifierKeySetConfig) (VerifierKeySet, error) {
	keySetID := strings.TrimSpace(config.KeySetID)
	if keySetID == "" {
		return VerifierKeySet{}, &VerifierKeyConfigError{Kind: ErrMissingKeySetID}
	}

	keys := []struct {
		purpose KeyPurpose
		value   []byte
	}{
		{purpose: KeyPurposeCredentialLookup, value: config.CredentialLookupKey},
		{purpose: KeyPurposeCredentialVerifier, value: config.CredentialVerifierKey},
		{purpose: KeyPurposeTokenLookup, value: config.TokenLookupKey},
		{purpose: KeyPurposeTokenVerifier, value: config.TokenVerifierKey},
	}

	for _, key := range keys {
		if len(key.value) == 0 {
			return VerifierKeySet{}, &VerifierKeyConfigError{Kind: ErrMissingKey, Purpose: key.purpose}
		}
		if len(key.value) < MinVerifierKeyBytes {
			return VerifierKeySet{}, &VerifierKeyConfigError{Kind: ErrKeyTooShort, Purpose: key.purpose}
		}
		if isWeakRepeatedKey(key.value) {
			return VerifierKeySet{}, &VerifierKeyConfigError{Kind: ErrWeakRepeatedKey, Purpose: key.purpose}
		}
	}

	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if equalBytes(keys[i].value, keys[j].value) {
				return VerifierKeySet{}, &VerifierKeyConfigError{Kind: ErrDuplicateKey, Purpose: keys[j].purpose}
			}
		}
	}

	return VerifierKeySet{
		keySetID:              keySetID,
		credentialLookupKey:   cloneBytes(config.CredentialLookupKey),
		credentialVerifierKey: cloneBytes(config.CredentialVerifierKey),
		tokenLookupKey:        cloneBytes(config.TokenLookupKey),
		tokenVerifierKey:      cloneBytes(config.TokenVerifierKey),
	}, nil
}

func (s VerifierKeySet) KeySetID() string {
	return s.keySetID
}

func (s VerifierKeySet) CredentialLookupKey() []byte {
	return cloneBytes(s.credentialLookupKey)
}

func (s VerifierKeySet) CredentialVerifierKey() []byte {
	return cloneBytes(s.credentialVerifierKey)
}

func (s VerifierKeySet) TokenLookupKey() []byte {
	return cloneBytes(s.tokenLookupKey)
}

func (s VerifierKeySet) TokenVerifierKey() []byte {
	return cloneBytes(s.tokenVerifierKey)
}

func cloneBytes(value []byte) []byte {
	if value == nil {
		return nil
	}
	copyValue := make([]byte, len(value))
	copy(copyValue, value)
	return copyValue
}

func equalBytes(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func isWeakRepeatedKey(value []byte) bool {
	if len(value) == 0 {
		return false
	}
	first := value[0]
	for _, item := range value[1:] {
		if item != first {
			return false
		}
	}
	return true
}
