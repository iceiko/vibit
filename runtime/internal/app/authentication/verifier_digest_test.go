package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
)

func TestCanonicalDigestInputIsDeterministic(t *testing.T) {
	raw := bytesWithSeed(11, RawSecretMaterialBytes)

	first, err := canonicalDigestInput(PurposeLabelCredentialLookup, raw)
	if err != nil {
		t.Fatalf("first canonicalDigestInput() error = %v, want nil", err)
	}
	second, err := canonicalDigestInput(PurposeLabelCredentialLookup, raw)
	if err != nil {
		t.Fatalf("second canonicalDigestInput() error = %v, want nil", err)
	}

	assertBytesEqual(t, first, second)
}

func TestCanonicalDigestInputUsesVersionHeader(t *testing.T) {
	input := mustCanonicalDigestInput(t, PurposeLabelCredentialLookup, bytesWithSeed(12, RawSecretMaterialBytes))
	header := []byte(CanonicalInputVersion)

	if len(input) < len(header) {
		t.Fatalf("len(input) = %d, want at least %d", len(input), len(header))
	}
	assertBytesEqual(t, input[:len(header)], header)
}

func TestCanonicalDigestInputNullSeparatorPresent(t *testing.T) {
	input := mustCanonicalDigestInput(t, PurposeLabelCredentialLookup, bytesWithSeed(13, RawSecretMaterialBytes))

	if got := input[len(CanonicalInputVersion)]; got != 0x00 {
		t.Fatalf("separator = %#x, want 0x00", got)
	}
}

func TestCanonicalDigestInputLengthPrefixesPurposeLabel(t *testing.T) {
	purpose := PurposeLabelCredentialVerifier
	input := mustCanonicalDigestInput(t, purpose, bytesWithSeed(14, RawSecretMaterialBytes))
	offset := len(CanonicalInputVersion) + 1

	if got := binary.BigEndian.Uint16(input[offset : offset+2]); got != uint16(len(purpose)) {
		t.Fatalf("purpose length prefix = %d, want %d", got, len(purpose))
	}
	assertBytesEqual(t, input[offset+2:offset+2+len(purpose)], []byte(purpose))
}

func TestCanonicalDigestInputLengthPrefixesRawMaterial(t *testing.T) {
	purpose := PurposeLabelTokenLookup
	raw := bytesWithSeed(15, RawSecretMaterialBytes)
	input := mustCanonicalDigestInput(t, purpose, raw)
	offset := len(CanonicalInputVersion) + 1 + 2 + len(purpose)

	if got := binary.BigEndian.Uint16(input[offset : offset+2]); got != uint16(len(raw)) {
		t.Fatalf("raw material length prefix = %d, want %d", got, len(raw))
	}
	assertBytesEqual(t, input[offset+2:], raw)
}

func TestLookupAndVerifierPurposeLabelsDiffer(t *testing.T) {
	if PurposeLabelCredentialLookup == PurposeLabelCredentialVerifier {
		t.Fatal("credential lookup and verifier purpose labels must differ")
	}
	if PurposeLabelTokenLookup == PurposeLabelTokenVerifier {
		t.Fatal("token lookup and verifier purpose labels must differ")
	}
}

func TestCredentialAndTokenPurposeLabelsDiffer(t *testing.T) {
	if PurposeLabelCredentialLookup == PurposeLabelTokenLookup {
		t.Fatal("credential and token lookup purpose labels must differ")
	}
	if PurposeLabelCredentialVerifier == PurposeLabelTokenVerifier {
		t.Fatal("credential and token verifier purpose labels must differ")
	}
}

func TestDigestOutputIs32Bytes(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	raw := bytesWithSeed(16, RawSecretMaterialBytes)

	digest, err := ComputeCredentialLookupDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeCredentialLookupDigest() error = %v, want nil", err)
	}

	if len(digest.Bytes()) != VerifierDigestBytes {
		t.Fatalf("len(Bytes()) = %d, want %d", len(digest.Bytes()), VerifierDigestBytes)
	}
	if len(digest.Bytes()) != sha256.Size {
		t.Fatalf("len(Bytes()) = %d, want sha256.Size %d", len(digest.Bytes()), sha256.Size)
	}
}

func TestCredentialLookupDigestUsesCredentialLookupKey(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	raw := bytesWithSeed(17, RawSecretMaterialBytes)

	digest, err := ComputeCredentialLookupDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeCredentialLookupDigest() error = %v, want nil", err)
	}

	assertDigest(t, digest, DigestClassCredentialLookup, expectedHMAC(t, keySet.CredentialLookupKey(), PurposeLabelCredentialLookup, raw))
}

func TestCredentialVerifierDigestUsesCredentialVerifierKey(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	raw := bytesWithSeed(18, RawSecretMaterialBytes)

	digest, err := ComputeCredentialVerifierDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeCredentialVerifierDigest() error = %v, want nil", err)
	}

	assertDigest(t, digest, DigestClassCredentialVerifier, expectedHMAC(t, keySet.CredentialVerifierKey(), PurposeLabelCredentialVerifier, raw))
}

func TestTokenLookupDigestUsesTokenLookupKey(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	raw := bytesWithSeed(19, RawSecretMaterialBytes)

	digest, err := ComputeTokenLookupDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeTokenLookupDigest() error = %v, want nil", err)
	}

	assertDigest(t, digest, DigestClassTokenLookup, expectedHMAC(t, keySet.TokenLookupKey(), PurposeLabelTokenLookup, raw))
}

func TestTokenVerifierDigestUsesTokenVerifierKey(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	raw := bytesWithSeed(20, RawSecretMaterialBytes)

	digest, err := ComputeTokenVerifierDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeTokenVerifierDigest() error = %v, want nil", err)
	}

	assertDigest(t, digest, DigestClassTokenVerifier, expectedHMAC(t, keySet.TokenVerifierKey(), PurposeLabelTokenVerifier, raw))
}

func TestDifferentKeysProduceDifferentDigests(t *testing.T) {
	firstKeySet := mustVerifierKeySet(t)
	secondConfig := validVerifierKeySetConfig()
	secondConfig.CredentialLookupKey = bytesWithSeed(21, MinVerifierKeyBytes)
	secondKeySet, err := NewVerifierKeySet(secondConfig)
	if err != nil {
		t.Fatalf("NewVerifierKeySet() error = %v, want nil", err)
	}
	raw := bytesWithSeed(22, RawSecretMaterialBytes)

	first, err := ComputeCredentialLookupDigest(firstKeySet, raw)
	if err != nil {
		t.Fatalf("first ComputeCredentialLookupDigest() error = %v, want nil", err)
	}
	second, err := ComputeCredentialLookupDigest(secondKeySet, raw)
	if err != nil {
		t.Fatalf("second ComputeCredentialLookupDigest() error = %v, want nil", err)
	}

	if equalBytes(first.Bytes(), second.Bytes()) {
		t.Fatal("different keys produced the same digest")
	}
}

func TestDifferentRawMaterialProducesDifferentDigests(t *testing.T) {
	keySet := mustVerifierKeySet(t)

	first, err := ComputeTokenLookupDigest(keySet, bytesWithSeed(23, RawSecretMaterialBytes))
	if err != nil {
		t.Fatalf("first ComputeTokenLookupDigest() error = %v, want nil", err)
	}
	second, err := ComputeTokenLookupDigest(keySet, bytesWithSeed(24, RawSecretMaterialBytes))
	if err != nil {
		t.Fatalf("second ComputeTokenLookupDigest() error = %v, want nil", err)
	}

	if equalBytes(first.Bytes(), second.Bytes()) {
		t.Fatal("different raw material produced the same digest")
	}
}

func TestDigestBytesAreCopiedOnReturn(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	digest, err := ComputeTokenVerifierDigest(keySet, bytesWithSeed(25, RawSecretMaterialBytes))
	if err != nil {
		t.Fatalf("ComputeTokenVerifierDigest() error = %v, want nil", err)
	}

	first := digest.Bytes()
	first[0] = first[0] + 1
	second := digest.Bytes()
	if second[0] == first[0] {
		t.Fatal("Bytes() returned mutable internal digest bytes")
	}
}

func TestEmptyRawMaterialFailsClosed(t *testing.T) {
	_, err := ComputeCredentialLookupDigest(mustVerifierKeySet(t), nil)
	assertDigestError(t, err, ErrMissingRawMaterial, DigestClassCredentialLookup)
}

func TestMissingVerifierKeySetFailsClosed(t *testing.T) {
	_, err := ComputeTokenVerifierDigest(VerifierKeySet{}, bytesWithSeed(26, RawSecretMaterialBytes))
	assertDigestError(t, err, ErrMissingVerifierKeySet, DigestClassTokenVerifier)
}

func TestVerifierDigestErrorsDoNotIncludeDigestOrKeyMaterial(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	raw := bytesWithSeed(27, RawSecretMaterialBytes)
	digest, err := ComputeCredentialVerifierDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeCredentialVerifierDigest() error = %v, want nil", err)
	}

	_, err = ComputeCredentialVerifierDigest(VerifierKeySet{}, raw)
	if err == nil {
		t.Fatal("ComputeCredentialVerifierDigest() error = nil, want error")
	}

	message := err.Error()
	for _, allowed := range []string{ErrMissingVerifierKeySet.Error(), string(DigestClassCredentialVerifier)} {
		if !strings.Contains(message, allowed) {
			t.Fatalf("error %q does not include allowed marker %q", message, allowed)
		}
	}
	for _, forbidden := range [][]byte{
		raw,
		keySet.CredentialVerifierKey(),
		digest.Bytes(),
	} {
		text := string(forbidden)
		if text != "" && strings.Contains(message, text) {
			t.Fatalf("error %q contains forbidden material %q", message, text)
		}
	}
}

func TestVerifierDigestHelperDoesNotCompareVerifiers(t *testing.T) {
	keySet := mustVerifierKeySet(t)
	raw := bytesWithSeed(28, RawSecretMaterialBytes)

	lookup, err := ComputeCredentialLookupDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeCredentialLookupDigest() error = %v, want nil", err)
	}
	verifier, err := ComputeCredentialVerifierDigest(keySet, raw)
	if err != nil {
		t.Fatalf("ComputeCredentialVerifierDigest() error = %v, want nil", err)
	}

	if lookup.Class() == verifier.Class() {
		t.Fatal("lookup and verifier digest classes must differ")
	}
	if len(lookup.Bytes()) != VerifierDigestBytes || len(verifier.Bytes()) != VerifierDigestBytes {
		t.Fatal("digest helpers should only compute digest bytes")
	}
}

func mustVerifierKeySet(t *testing.T) VerifierKeySet {
	t.Helper()
	keySet, err := NewVerifierKeySet(validVerifierKeySetConfig())
	if err != nil {
		t.Fatalf("NewVerifierKeySet() error = %v, want nil", err)
	}
	return keySet
}

func mustCanonicalDigestInput(t *testing.T, purposeLabel string, rawMaterial []byte) []byte {
	t.Helper()
	input, err := canonicalDigestInput(purposeLabel, rawMaterial)
	if err != nil {
		t.Fatalf("canonicalDigestInput() error = %v, want nil", err)
	}
	return input
}

func expectedHMAC(t *testing.T, key []byte, purposeLabel string, rawMaterial []byte) []byte {
	t.Helper()
	input := mustCanonicalDigestInput(t, purposeLabel, rawMaterial)
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(input)
	return mac.Sum(nil)
}

func assertDigest(t *testing.T, digest ComputedDigest, class DigestClass, want []byte) {
	t.Helper()
	if digest.Class() != class {
		t.Fatalf("Class() = %q, want %q", digest.Class(), class)
	}
	assertBytesEqual(t, digest.Bytes(), want)
}

func assertDigestError(t *testing.T, err error, target error, class DigestClass) {
	t.Helper()
	assertErrorIs(t, err, target)
	var digestErr *VerifierDigestError
	if !errors.As(err, &digestErr) {
		t.Fatalf("error = %T, want *VerifierDigestError", err)
	}
	if digestErr.Class != class {
		t.Fatalf("Class = %q, want %q", digestErr.Class, class)
	}
}
