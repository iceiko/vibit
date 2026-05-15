package authentication

import (
	"errors"
	"strings"
	"testing"
)

func TestNewVerifierKeySetAcceptsCompleteDistinctKeys(t *testing.T) {
	config := validVerifierKeySetConfig()

	keySet, err := NewVerifierKeySet(config)
	if err != nil {
		t.Fatalf("NewVerifierKeySet() error = %v, want nil", err)
	}

	if keySet.KeySetID() != "key-set-1" {
		t.Fatalf("KeySetID() = %q, want key-set-1", keySet.KeySetID())
	}
	assertBytesEqual(t, keySet.CredentialLookupKey(), config.CredentialLookupKey)
	assertBytesEqual(t, keySet.CredentialVerifierKey(), config.CredentialVerifierKey)
	assertBytesEqual(t, keySet.TokenLookupKey(), config.TokenLookupKey)
	assertBytesEqual(t, keySet.TokenVerifierKey(), config.TokenVerifierKey)
}

func TestNewVerifierKeySetTrimsKeySetID(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.KeySetID = "  key-set-1  "

	keySet, err := NewVerifierKeySet(config)
	if err != nil {
		t.Fatalf("NewVerifierKeySet() error = %v, want nil", err)
	}
	if keySet.KeySetID() != "key-set-1" {
		t.Fatalf("KeySetID() = %q, want key-set-1", keySet.KeySetID())
	}
}

func TestNewVerifierKeySetCopiesInputKeyMaterial(t *testing.T) {
	config := validVerifierKeySetConfig()

	keySet, err := NewVerifierKeySet(config)
	if err != nil {
		t.Fatalf("NewVerifierKeySet() error = %v, want nil", err)
	}

	config.CredentialLookupKey[0] = 99
	config.CredentialVerifierKey[0] = 99
	config.TokenLookupKey[0] = 99
	config.TokenVerifierKey[0] = 99

	if keySet.CredentialLookupKey()[0] == 99 ||
		keySet.CredentialVerifierKey()[0] == 99 ||
		keySet.TokenLookupKey()[0] == 99 ||
		keySet.TokenVerifierKey()[0] == 99 {
		t.Fatal("validated key set changed after mutating caller input")
	}
}

func TestVerifierKeySetAccessorsDoNotExposeMutableInternalSlices(t *testing.T) {
	keySet, err := NewVerifierKeySet(validVerifierKeySetConfig())
	if err != nil {
		t.Fatalf("NewVerifierKeySet() error = %v, want nil", err)
	}

	credentialLookup := keySet.CredentialLookupKey()
	credentialVerifier := keySet.CredentialVerifierKey()
	tokenLookup := keySet.TokenLookupKey()
	tokenVerifier := keySet.TokenVerifierKey()

	credentialLookup[0] = 99
	credentialVerifier[0] = 99
	tokenLookup[0] = 99
	tokenVerifier[0] = 99

	if keySet.CredentialLookupKey()[0] == 99 ||
		keySet.CredentialVerifierKey()[0] == 99 ||
		keySet.TokenLookupKey()[0] == 99 ||
		keySet.TokenVerifierKey()[0] == 99 {
		t.Fatal("accessor returned mutable internal slice")
	}
}

func TestNewVerifierKeySetMissingKeySetIDFailsClosed(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.KeySetID = "  "

	_, err := NewVerifierKeySet(config)
	assertErrorIs(t, err, ErrMissingKeySetID)
}

func TestNewVerifierKeySetMissingEachLogicalKeyFailsClosed(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*VerifierKeySetConfig)
		purpose KeyPurpose
	}{
		{name: "credential lookup", mutate: func(c *VerifierKeySetConfig) { c.CredentialLookupKey = nil }, purpose: KeyPurposeCredentialLookup},
		{name: "credential verifier", mutate: func(c *VerifierKeySetConfig) { c.CredentialVerifierKey = nil }, purpose: KeyPurposeCredentialVerifier},
		{name: "token lookup", mutate: func(c *VerifierKeySetConfig) { c.TokenLookupKey = nil }, purpose: KeyPurposeTokenLookup},
		{name: "token verifier", mutate: func(c *VerifierKeySetConfig) { c.TokenVerifierKey = nil }, purpose: KeyPurposeTokenVerifier},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := validVerifierKeySetConfig()
			tc.mutate(&config)

			_, err := NewVerifierKeySet(config)
			assertConfigError(t, err, ErrMissingKey, tc.purpose)
		})
	}
}

func TestNewVerifierKeySetShortEachLogicalKeyFailsClosed(t *testing.T) {
	shortKey := bytesWithSeed(7, MinVerifierKeyBytes-1)
	cases := []struct {
		name    string
		mutate  func(*VerifierKeySetConfig)
		purpose KeyPurpose
	}{
		{name: "credential lookup", mutate: func(c *VerifierKeySetConfig) { c.CredentialLookupKey = shortKey }, purpose: KeyPurposeCredentialLookup},
		{name: "credential verifier", mutate: func(c *VerifierKeySetConfig) { c.CredentialVerifierKey = shortKey }, purpose: KeyPurposeCredentialVerifier},
		{name: "token lookup", mutate: func(c *VerifierKeySetConfig) { c.TokenLookupKey = shortKey }, purpose: KeyPurposeTokenLookup},
		{name: "token verifier", mutate: func(c *VerifierKeySetConfig) { c.TokenVerifierKey = shortKey }, purpose: KeyPurposeTokenVerifier},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := validVerifierKeySetConfig()
			tc.mutate(&config)

			_, err := NewVerifierKeySet(config)
			assertConfigError(t, err, ErrKeyTooShort, tc.purpose)
		})
	}
}

func TestNewVerifierKeySetDuplicateLogicalKeysFailClosed(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.TokenVerifierKey = cloneBytes(config.CredentialLookupKey)

	_, err := NewVerifierKeySet(config)
	assertConfigError(t, err, ErrDuplicateKey, KeyPurposeTokenVerifier)
}

func TestNewVerifierKeySetAllZeroKeyFailsClosed(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.CredentialLookupKey = make([]byte, MinVerifierKeyBytes)

	_, err := NewVerifierKeySet(config)
	assertConfigError(t, err, ErrWeakRepeatedKey, KeyPurposeCredentialLookup)
}

func TestNewVerifierKeySetRepeatedSingleByteKeyFailsClosed(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.TokenLookupKey = repeatedBytes(7, MinVerifierKeyBytes)

	_, err := NewVerifierKeySet(config)
	assertConfigError(t, err, ErrWeakRepeatedKey, KeyPurposeTokenLookup)
}

func TestNewVerifierKeySetErrorsDoNotIncludeSecretBytesOrFullKeySetID(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.KeySetID = "sensitive-key-set-id"
	config.CredentialVerifierKey = []byte("short-secret-fragment")

	_, err := NewVerifierKeySet(config)
	if err == nil {
		t.Fatal("NewVerifierKeySet() error = nil, want validation error")
	}

	message := err.Error()
	for _, forbidden := range []string{
		"sensitive-key-set-id",
		"short-secret-fragment",
		string(config.CredentialLookupKey),
		string(config.TokenLookupKey),
		string(config.TokenVerifierKey),
	} {
		if forbidden != "" && strings.Contains(message, forbidden) {
			t.Fatalf("error %q contains forbidden secret material %q", message, forbidden)
		}
	}
}

func validVerifierKeySetConfig() VerifierKeySetConfig {
	return VerifierKeySetConfig{
		KeySetID:              "key-set-1",
		CredentialLookupKey:   bytesWithSeed(1, MinVerifierKeyBytes),
		CredentialVerifierKey: bytesWithSeed(2, MinVerifierKeyBytes),
		TokenLookupKey:        bytesWithSeed(3, MinVerifierKeyBytes),
		TokenVerifierKey:      bytesWithSeed(4, MinVerifierKeyBytes),
	}
}

func bytesWithSeed(seed byte, length int) []byte {
	value := make([]byte, length)
	for i := range value {
		value[i] = seed + byte(i%11)
	}
	return value
}

func repeatedBytes(value byte, length int) []byte {
	items := make([]byte, length)
	for i := range items {
		items[i] = value
	}
	return items
}

func assertBytesEqual(t *testing.T, got, want []byte) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func assertErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("error = %v, want errors.Is(_, %v)", err, target)
	}
}

func assertConfigError(t *testing.T, err error, target error, purpose KeyPurpose) {
	t.Helper()
	assertErrorIs(t, err, target)
	var configErr *VerifierKeyConfigError
	if !errors.As(err, &configErr) {
		t.Fatalf("error = %T, want *VerifierKeyConfigError", err)
	}
	if configErr.Purpose != purpose {
		t.Fatalf("Purpose = %q, want %q", configErr.Purpose, purpose)
	}
}
