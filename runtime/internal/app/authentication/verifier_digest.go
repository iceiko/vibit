package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	CanonicalInputVersion = "vibit.auth.verifier.input.v1"
	VerifierDigestBytes   = sha256.Size
)

type DigestClass string

const (
	DigestClassCredentialLookup   DigestClass = "credential_lookup"
	DigestClassCredentialVerifier DigestClass = "credential_verifier"
	DigestClassTokenLookup        DigestClass = "token_lookup"
	DigestClassTokenVerifier      DigestClass = "token_verifier"
)

const (
	PurposeLabelCredentialLookup   = "vibit.credential.lookup.v1"
	PurposeLabelCredentialVerifier = "vibit.credential.verifier.v1"
	PurposeLabelTokenLookup        = "vibit.access_token.lookup.v1"
	PurposeLabelTokenVerifier      = "vibit.access_token.verifier.v1"
)

var (
	ErrMissingVerifierKeySet = errors.New("authentication verifier digest: missing verifier key set")
	ErrMissingRawMaterial    = errors.New("authentication verifier digest: missing raw material")
	ErrInvalidDigestInput    = errors.New("authentication verifier digest: invalid digest input")
)

type VerifierDigestError struct {
	Kind  error
	Class DigestClass
}

func (e *VerifierDigestError) Error() string {
	if e == nil {
		return ""
	}
	if e.Class == "" {
		return e.Kind.Error()
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Class)
}

func (e *VerifierDigestError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Kind
}

type ComputedDigest struct {
	class DigestClass
	bytes []byte
}

func (d ComputedDigest) Class() DigestClass {
	return d.class
}

func (d ComputedDigest) Bytes() []byte {
	return cloneBytes(d.bytes)
}

func ComputeCredentialLookupDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error) {
	return computeVerifierDigest(DigestClassCredentialLookup, PurposeLabelCredentialLookup, keySet.CredentialLookupKey(), rawMaterial)
}

func ComputeCredentialVerifierDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error) {
	return computeVerifierDigest(DigestClassCredentialVerifier, PurposeLabelCredentialVerifier, keySet.CredentialVerifierKey(), rawMaterial)
}

func ComputeTokenLookupDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error) {
	return computeVerifierDigest(DigestClassTokenLookup, PurposeLabelTokenLookup, keySet.TokenLookupKey(), rawMaterial)
}

func ComputeTokenVerifierDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error) {
	return computeVerifierDigest(DigestClassTokenVerifier, PurposeLabelTokenVerifier, keySet.TokenVerifierKey(), rawMaterial)
}

func computeVerifierDigest(class DigestClass, purposeLabel string, key []byte, rawMaterial []byte) (ComputedDigest, error) {
	if len(key) == 0 {
		return ComputedDigest{}, &VerifierDigestError{Kind: ErrMissingVerifierKeySet, Class: class}
	}
	if len(rawMaterial) == 0 {
		return ComputedDigest{}, &VerifierDigestError{Kind: ErrMissingRawMaterial, Class: class}
	}

	input, err := canonicalDigestInput(purposeLabel, rawMaterial)
	if err != nil {
		return ComputedDigest{}, &VerifierDigestError{Kind: ErrInvalidDigestInput, Class: class}
	}

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(input)
	return ComputedDigest{
		class: class,
		bytes: cloneBytes(mac.Sum(nil)),
	}, nil
}

func canonicalDigestInput(purposeLabel string, rawMaterial []byte) ([]byte, error) {
	if purposeLabel == "" {
		return nil, ErrInvalidDigestInput
	}
	if len(purposeLabel) > 0xffff || len(rawMaterial) > 0xffff {
		return nil, ErrInvalidDigestInput
	}

	input := make([]byte, 0, len(CanonicalInputVersion)+1+2+len(purposeLabel)+2+len(rawMaterial))
	input = append(input, CanonicalInputVersion...)
	input = append(input, 0x00)
	input = appendUint16(input, len(purposeLabel))
	input = append(input, purposeLabel...)
	input = appendUint16(input, len(rawMaterial))
	input = append(input, rawMaterial...)
	return input, nil
}

func appendUint16(target []byte, value int) []byte {
	var encoded [2]byte
	binary.BigEndian.PutUint16(encoded[:], uint16(value))
	return append(target, encoded[:]...)
}
