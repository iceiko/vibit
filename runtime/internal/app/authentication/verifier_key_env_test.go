package authentication

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

func TestLoadVerifierKeySetFromEnvironmentAcceptsBase64URLUnpaddedKeySet(t *testing.T) {
	config := validVerifierKeySetConfig()
	env := environmentFromConfig(config, encodeRawURL)

	keySet, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))
	if err != nil {
		t.Fatalf("LoadVerifierKeySetFromEnvironment() error = %v, want nil", err)
	}

	assertLoadedKeySet(t, keySet, config)
}

func TestLoadVerifierKeySetFromEnvironmentAcceptsStandardBase64PaddedKeySet(t *testing.T) {
	config := validVerifierKeySetConfig()
	env := environmentFromConfig(config, base64.StdEncoding.EncodeToString)

	keySet, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))
	if err != nil {
		t.Fatalf("LoadVerifierKeySetFromEnvironment() error = %v, want nil", err)
	}

	assertLoadedKeySet(t, keySet, config)
}

func TestLoadVerifierKeySetFromEnvironmentRequiresLookup(t *testing.T) {
	_, err := LoadVerifierKeySetFromEnvironment(nil)

	assertEnvironmentError(t, err, ErrMissingEnvironmentLookup, "", "")
}

func TestLoadVerifierKeySetFromEnvironmentMissingEachVariableFailsClosed(t *testing.T) {
	config := validVerifierKeySetConfig()
	env := environmentFromConfig(config, encodeRawURL)

	for _, variable := range verifierKeyEnvironmentVariableNames() {
		t.Run(variable, func(t *testing.T) {
			mutated := cloneStringMap(env)
			delete(mutated, variable)

			_, err := LoadVerifierKeySetFromEnvironment(mapLookup(mutated))
			assertEnvironmentError(t, err, ErrMissingEnvironmentVariable, variable, "")
		})
	}
}

func TestLoadVerifierKeySetFromEnvironmentInvalidEachEncodedKeyFailsClosed(t *testing.T) {
	config := validVerifierKeySetConfig()
	env := environmentFromConfig(config, encodeRawURL)

	cases := []struct {
		variable string
		purpose  KeyPurpose
	}{
		{variable: EnvCredentialLookupKey, purpose: KeyPurposeCredentialLookup},
		{variable: EnvCredentialVerifierKey, purpose: KeyPurposeCredentialVerifier},
		{variable: EnvTokenLookupKey, purpose: KeyPurposeTokenLookup},
		{variable: EnvTokenVerifierKey, purpose: KeyPurposeTokenVerifier},
	}

	for _, tc := range cases {
		t.Run(tc.variable, func(t *testing.T) {
			mutated := cloneStringMap(env)
			mutated[tc.variable] = "not valid base64 text!"

			_, err := LoadVerifierKeySetFromEnvironment(mapLookup(mutated))
			assertEnvironmentError(t, err, ErrInvalidEnvironmentKeyText, tc.variable, tc.purpose)
		})
	}
}

func TestLoadVerifierKeySetFromEnvironmentDecodedShortKeyFailsThroughValidator(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.CredentialLookupKey = bytesWithSeed(9, MinVerifierKeyBytes-1)
	env := environmentFromConfig(config, encodeRawURL)

	_, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))

	assertEnvironmentError(t, err, ErrInvalidEnvironmentKeySet, EnvCredentialLookupKey, KeyPurposeCredentialLookup)
	assertErrorIs(t, err, ErrKeyTooShort)
}

func TestLoadVerifierKeySetFromEnvironmentDuplicateDecodedKeysFailThroughValidator(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.TokenVerifierKey = cloneBytes(config.CredentialLookupKey)
	env := environmentFromConfig(config, encodeRawURL)

	_, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))

	assertEnvironmentError(t, err, ErrInvalidEnvironmentKeySet, EnvTokenVerifierKey, KeyPurposeTokenVerifier)
	assertErrorIs(t, err, ErrDuplicateKey)
}

func TestLoadVerifierKeySetFromEnvironmentAllZeroKeyFailsThroughValidator(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.CredentialVerifierKey = make([]byte, MinVerifierKeyBytes)
	env := environmentFromConfig(config, encodeRawURL)

	_, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))

	assertEnvironmentError(t, err, ErrInvalidEnvironmentKeySet, EnvCredentialVerifierKey, KeyPurposeCredentialVerifier)
	assertErrorIs(t, err, ErrWeakRepeatedKey)
}

func TestLoadVerifierKeySetFromEnvironmentRepeatedSingleByteKeyFailsThroughValidator(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.TokenLookupKey = repeatedBytes(8, MinVerifierKeyBytes)
	env := environmentFromConfig(config, encodeRawURL)

	_, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))

	assertEnvironmentError(t, err, ErrInvalidEnvironmentKeySet, EnvTokenLookupKey, KeyPurposeTokenLookup)
	assertErrorIs(t, err, ErrWeakRepeatedKey)
}

func TestLoadVerifierKeySetFromEnvironmentEmptyKeySetIDFailsThroughValidator(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.KeySetID = "   "
	env := environmentFromConfig(config, encodeRawURL)

	_, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))

	assertEnvironmentError(t, err, ErrInvalidEnvironmentKeySet, EnvVerifierKeySetID, "")
	assertErrorIs(t, err, ErrMissingKeySetID)
}

func TestLoadVerifierKeySetFromEnvironmentErrorsDoNotIncludeSecretValues(t *testing.T) {
	config := validVerifierKeySetConfig()
	config.KeySetID = "sensitive-key-set-id"
	env := environmentFromConfig(config, encodeRawURL)
	env[EnvCredentialLookupKey] = "not valid base64 text!"

	_, err := LoadVerifierKeySetFromEnvironment(mapLookup(env))
	if err == nil {
		t.Fatal("LoadVerifierKeySetFromEnvironment() error = nil, want error")
	}

	message := err.Error()
	for _, allowed := range []string{ErrInvalidEnvironmentKeyText.Error(), EnvCredentialLookupKey} {
		if !strings.Contains(message, allowed) {
			t.Fatalf("error %q does not include safe context %q", message, allowed)
		}
	}
	for _, forbidden := range []string{
		"sensitive-key-set-id",
		env[EnvCredentialLookupKey],
		env[EnvCredentialVerifierKey],
		env[EnvTokenLookupKey],
		env[EnvTokenVerifierKey],
		string(config.CredentialLookupKey),
		string(config.CredentialVerifierKey),
		string(config.TokenLookupKey),
		string(config.TokenVerifierKey),
	} {
		if forbidden != "" && strings.Contains(message, forbidden) {
			t.Fatalf("error %q contains forbidden secret material %q", message, forbidden)
		}
	}
}

func TestLoadVerifierKeySetFromProcessEnvironmentReadsRequiredEnvironmentOnly(t *testing.T) {
	config := validVerifierKeySetConfig()
	env := environmentFromConfig(config, encodeRawURL)
	for name, value := range env {
		t.Setenv(name, value)
	}

	keySet, err := LoadVerifierKeySetFromProcessEnvironment()
	if err != nil {
		t.Fatalf("LoadVerifierKeySetFromProcessEnvironment() error = %v, want nil", err)
	}

	assertLoadedKeySet(t, keySet, config)
}

func environmentFromConfig(config VerifierKeySetConfig, encode func([]byte) string) map[string]string {
	return map[string]string{
		EnvVerifierKeySetID:      config.KeySetID,
		EnvCredentialLookupKey:   encode(config.CredentialLookupKey),
		EnvCredentialVerifierKey: encode(config.CredentialVerifierKey),
		EnvTokenLookupKey:        encode(config.TokenLookupKey),
		EnvTokenVerifierKey:      encode(config.TokenVerifierKey),
	}
}

func verifierKeyEnvironmentVariableNames() []string {
	return []string{
		EnvVerifierKeySetID,
		EnvCredentialLookupKey,
		EnvCredentialVerifierKey,
		EnvTokenLookupKey,
		EnvTokenVerifierKey,
	}
}

func mapLookup(values map[string]string) EnvLookup {
	return func(name string) (string, bool) {
		value, ok := values[name]
		return value, ok
	}
}

func cloneStringMap(values map[string]string) map[string]string {
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func encodeRawURL(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}

func assertLoadedKeySet(t *testing.T, keySet VerifierKeySet, config VerifierKeySetConfig) {
	t.Helper()
	if keySet.KeySetID() != strings.TrimSpace(config.KeySetID) {
		t.Fatalf("KeySetID() = %q, want %q", keySet.KeySetID(), strings.TrimSpace(config.KeySetID))
	}
	assertBytesEqual(t, keySet.CredentialLookupKey(), config.CredentialLookupKey)
	assertBytesEqual(t, keySet.CredentialVerifierKey(), config.CredentialVerifierKey)
	assertBytesEqual(t, keySet.TokenLookupKey(), config.TokenLookupKey)
	assertBytesEqual(t, keySet.TokenVerifierKey(), config.TokenVerifierKey)
}

func assertEnvironmentError(t *testing.T, err error, target error, variable string, purpose KeyPurpose) {
	t.Helper()
	assertErrorIs(t, err, target)

	var envErr *VerifierKeyEnvironmentError
	if !errors.As(err, &envErr) {
		t.Fatalf("error = %T, want *VerifierKeyEnvironmentError", err)
	}
	if envErr.Variable != variable {
		t.Fatalf("Variable = %q, want %q", envErr.Variable, variable)
	}
	if envErr.Purpose != purpose {
		t.Fatalf("Purpose = %q, want %q", envErr.Purpose, purpose)
	}
}
