package authentication

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestCredentialVerifierDigestMatchReturnsMatched(t *testing.T) {
	digest := mustCredentialVerifierDigest(t, 31)

	result, err := CompareCredentialVerifierDigest(digest, digest.Bytes())
	if err != nil {
		t.Fatalf("CompareCredentialVerifierDigest() error = %v, want nil", err)
	}

	assertComparisonResult(t, result, DigestClassCredentialVerifier, true)
}

func TestCredentialVerifierDigestMismatchReturnsNotMatched(t *testing.T) {
	first := mustCredentialVerifierDigest(t, 32)
	second := mustCredentialVerifierDigest(t, 33)

	result, err := CompareCredentialVerifierDigest(first, second.Bytes())
	if err != nil {
		t.Fatalf("CompareCredentialVerifierDigest() error = %v, want nil", err)
	}

	assertComparisonResult(t, result, DigestClassCredentialVerifier, false)
}

func TestTokenVerifierDigestMatchReturnsMatched(t *testing.T) {
	digest := mustTokenVerifierDigest(t, 34)

	result, err := CompareTokenVerifierDigest(digest, digest.Bytes())
	if err != nil {
		t.Fatalf("CompareTokenVerifierDigest() error = %v, want nil", err)
	}

	assertComparisonResult(t, result, DigestClassTokenVerifier, true)
}

func TestTokenVerifierDigestMismatchReturnsNotMatched(t *testing.T) {
	first := mustTokenVerifierDigest(t, 35)
	second := mustTokenVerifierDigest(t, 36)

	result, err := CompareTokenVerifierDigest(first, second.Bytes())
	if err != nil {
		t.Fatalf("CompareTokenVerifierDigest() error = %v, want nil", err)
	}

	assertComparisonResult(t, result, DigestClassTokenVerifier, false)
}

func TestVerifierComparisonUsesCryptoHmacEqual(t *testing.T) {
	source, err := os.ReadFile("verifier_comparison.go")
	if err != nil {
		t.Fatalf("ReadFile(verifier_comparison.go) error = %v, want nil", err)
	}
	text := string(source)

	if !strings.Contains(text, "\"crypto/hmac\"") {
		t.Fatal("verifier_comparison.go does not import crypto/hmac")
	}
	if !strings.Contains(text, "hmac.Equal(") {
		t.Fatal("verifier_comparison.go does not use hmac.Equal")
	}
	for _, forbidden := range []string{
		"bytes.Equal(",
		"reflect.DeepEqual(",
		"subtle.ConstantTimeCompare(",
		"string(computed",
		"string(stored",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("verifier_comparison.go contains forbidden comparison marker %q", forbidden)
		}
	}
}

func TestVerifierComparisonRejectsLookupDigestClasses(t *testing.T) {
	credentialLookup := mustCredentialLookupDigest(t, 37)
	tokenLookup := mustTokenLookupDigest(t, 38)

	_, err := CompareCredentialVerifierDigest(credentialLookup, credentialLookup.Bytes())
	assertComparisonError(t, err, ErrWrongComputedDigestClass, DigestClassCredentialLookup)

	_, err = CompareTokenVerifierDigest(tokenLookup, tokenLookup.Bytes())
	assertComparisonError(t, err, ErrWrongComputedDigestClass, DigestClassTokenLookup)
}

func TestVerifierComparisonRejectsWrongVerifierDigestClass(t *testing.T) {
	credential := mustCredentialVerifierDigest(t, 39)
	token := mustTokenVerifierDigest(t, 40)

	_, err := CompareCredentialVerifierDigest(token, token.Bytes())
	assertComparisonError(t, err, ErrWrongComputedDigestClass, DigestClassTokenVerifier)

	_, err = CompareTokenVerifierDigest(credential, credential.Bytes())
	assertComparisonError(t, err, ErrWrongComputedDigestClass, DigestClassCredentialVerifier)
}

func TestVerifierComparisonRejectsMissingComputedDigest(t *testing.T) {
	_, err := CompareCredentialVerifierDigest(ComputedDigest{}, bytesWithSeed(41, VerifierDigestBytes))
	assertComparisonError(t, err, ErrMissingComputedDigest, DigestClassCredentialVerifier)
}

func TestVerifierComparisonRejectsMissingStoredDigest(t *testing.T) {
	digest := mustTokenVerifierDigest(t, 42)

	_, err := CompareTokenVerifierDigest(digest, nil)
	assertComparisonError(t, err, ErrMissingStoredDigest, DigestClassTokenVerifier)
}

func TestVerifierComparisonRejectsMalformedStoredDigestLength(t *testing.T) {
	digest := mustCredentialVerifierDigest(t, 43)

	_, err := CompareCredentialVerifierDigest(digest, bytesWithSeed(44, VerifierDigestBytes-1))
	assertComparisonError(t, err, ErrMalformedStoredDigest, DigestClassCredentialVerifier)
}

func TestVerifierComparisonRejectsMalformedComputedDigestLength(t *testing.T) {
	digest := ComputedDigest{
		class: DigestClassTokenVerifier,
		bytes: bytesWithSeed(45, VerifierDigestBytes-1),
	}

	_, err := CompareTokenVerifierDigest(digest, bytesWithSeed(46, VerifierDigestBytes))
	assertComparisonError(t, err, ErrInvalidVerifierComparison, DigestClassTokenVerifier)
}

func TestVerifierComparisonResultDoesNotExposeDigestBytes(t *testing.T) {
	digest := mustCredentialVerifierDigest(t, 47)
	result, err := CompareCredentialVerifierDigest(digest, digest.Bytes())
	if err != nil {
		t.Fatalf("CompareCredentialVerifierDigest() error = %v, want nil", err)
	}

	rendered := result.Class() + DigestClass(result.MatchedStringForTestOnly())
	forbiddenDigestText := string(digest.Bytes())
	if forbiddenDigestText != "" && strings.Contains(string(rendered), forbiddenDigestText) {
		t.Fatalf("comparison result exposed digest bytes")
	}
}

func TestVerifierComparisonErrorsAreRedacted(t *testing.T) {
	computed := mustTokenVerifierDigest(t, 48)
	stored := bytesWithSeed(49, VerifierDigestBytes-1)

	_, err := CompareTokenVerifierDigest(computed, stored)
	if err == nil {
		t.Fatal("CompareTokenVerifierDigest() error = nil, want error")
	}

	message := err.Error()
	for _, allowed := range []string{ErrMalformedStoredDigest.Error(), string(DigestClassTokenVerifier)} {
		if !strings.Contains(message, allowed) {
			t.Fatalf("error %q does not include allowed marker %q", message, allowed)
		}
	}
	for _, forbidden := range [][]byte{
		computed.Bytes(),
		stored,
		bytesWithSeed(50, RawSecretMaterialBytes),
		mustVerifierKeySet(t).TokenVerifierKey(),
	} {
		text := string(forbidden)
		if text != "" && strings.Contains(message, text) {
			t.Fatalf("error %q contains forbidden material %q", message, text)
		}
	}
}

func TestVerifierComparisonDoesNotCompareRawMaterial(t *testing.T) {
	source, err := os.ReadFile("verifier_comparison.go")
	if err != nil {
		t.Fatalf("ReadFile(verifier_comparison.go) error = %v, want nil", err)
	}
	text := string(source)

	for _, required := range []string{
		"func CompareCredentialVerifierDigest(computed ComputedDigest, storedDigest []byte)",
		"func CompareTokenVerifierDigest(computed ComputedDigest, storedDigest []byte)",
	} {
		if !strings.Contains(text, required) {
			t.Fatalf("verifier_comparison.go is missing required signature %q", required)
		}
	}
	for _, forbidden := range []string{
		"GeneratedSecretMaterial",
		"RawBytes()",
		"rawMaterial",
		"rawCredential",
		"rawToken",
		"deviceCredential",
		"accessToken",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("verifier_comparison.go contains forbidden raw-material marker %q", forbidden)
		}
	}
}

func TestVerifierComparisonHelpersDoNotImplementAuthenticationServiceBehavior(t *testing.T) {
	source, err := os.ReadFile("verifier_comparison.go")
	if err != nil {
		t.Fatalf("ReadFile(verifier_comparison.go) error = %v, want nil", err)
	}
	text := string(source)

	for _, forbidden := range []string{
		"AuthenticateWithDeviceCredential",
		"ValidateAccessToken",
		"LogoutAccessToken",
		"RefreshAccessToken",
		"Repository",
		"Bearer",
		"Authorization",
		"WebSocket",
		"protobuf",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("verifier_comparison.go contains forbidden service marker %q", forbidden)
		}
	}
}

func mustCredentialLookupDigest(t *testing.T, seed byte) ComputedDigest {
	t.Helper()
	digest, err := ComputeCredentialLookupDigest(mustVerifierKeySet(t), bytesWithSeed(seed, RawSecretMaterialBytes))
	if err != nil {
		t.Fatalf("ComputeCredentialLookupDigest() error = %v, want nil", err)
	}
	return digest
}

func mustCredentialVerifierDigest(t *testing.T, seed byte) ComputedDigest {
	t.Helper()
	digest, err := ComputeCredentialVerifierDigest(mustVerifierKeySet(t), bytesWithSeed(seed, RawSecretMaterialBytes))
	if err != nil {
		t.Fatalf("ComputeCredentialVerifierDigest() error = %v, want nil", err)
	}
	return digest
}

func mustTokenLookupDigest(t *testing.T, seed byte) ComputedDigest {
	t.Helper()
	digest, err := ComputeTokenLookupDigest(mustVerifierKeySet(t), bytesWithSeed(seed, RawSecretMaterialBytes))
	if err != nil {
		t.Fatalf("ComputeTokenLookupDigest() error = %v, want nil", err)
	}
	return digest
}

func mustTokenVerifierDigest(t *testing.T, seed byte) ComputedDigest {
	t.Helper()
	digest, err := ComputeTokenVerifierDigest(mustVerifierKeySet(t), bytesWithSeed(seed, RawSecretMaterialBytes))
	if err != nil {
		t.Fatalf("ComputeTokenVerifierDigest() error = %v, want nil", err)
	}
	return digest
}

func assertComparisonResult(t *testing.T, result VerifierComparisonResult, class DigestClass, matched bool) {
	t.Helper()
	if result.Class() != class {
		t.Fatalf("Class() = %q, want %q", result.Class(), class)
	}
	if result.Matched() != matched {
		t.Fatalf("Matched() = %v, want %v", result.Matched(), matched)
	}
}

func assertComparisonError(t *testing.T, err error, target error, class DigestClass) {
	t.Helper()
	assertErrorIs(t, err, target)
	var comparisonErr *VerifierComparisonError
	if !errors.As(err, &comparisonErr) {
		t.Fatalf("error = %T, want *VerifierComparisonError", err)
	}
	if comparisonErr.Class != class {
		t.Fatalf("Class = %q, want %q", comparisonErr.Class, class)
	}
}

func (r VerifierComparisonResult) MatchedStringForTestOnly() string {
	if r.Matched() {
		return "matched"
	}
	return "not_matched"
}
