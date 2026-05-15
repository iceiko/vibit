package authentication

import (
	"crypto/hmac"
	"errors"
	"fmt"
)

var (
	ErrVerifierDigestMismatch    = errors.New("authentication verifier comparison: verifier digest mismatch")
	ErrMissingComputedDigest     = errors.New("authentication verifier comparison: missing computed digest")
	ErrWrongComputedDigestClass  = errors.New("authentication verifier comparison: wrong computed digest class")
	ErrMissingStoredDigest       = errors.New("authentication verifier comparison: missing stored digest")
	ErrMalformedStoredDigest     = errors.New("authentication verifier comparison: malformed stored digest")
	ErrInvalidVerifierComparison = errors.New("authentication verifier comparison: invalid comparison input")
)

type VerifierComparisonError struct {
	Kind  error
	Class DigestClass
}

func (e *VerifierComparisonError) Error() string {
	if e == nil {
		return ""
	}
	if e.Class == "" {
		return e.Kind.Error()
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Class)
}

func (e *VerifierComparisonError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Kind
}

type VerifierComparisonResult struct {
	class   DigestClass
	matched bool
}

func (r VerifierComparisonResult) Class() DigestClass {
	return r.class
}

func (r VerifierComparisonResult) Matched() bool {
	return r.matched
}

func CompareCredentialVerifierDigest(computed ComputedDigest, storedDigest []byte) (VerifierComparisonResult, error) {
	return compareVerifierDigest(DigestClassCredentialVerifier, computed, storedDigest)
}

func CompareTokenVerifierDigest(computed ComputedDigest, storedDigest []byte) (VerifierComparisonResult, error) {
	return compareVerifierDigest(DigestClassTokenVerifier, computed, storedDigest)
}

func compareVerifierDigest(expectedClass DigestClass, computed ComputedDigest, storedDigest []byte) (VerifierComparisonResult, error) {
	if expectedClass != DigestClassCredentialVerifier && expectedClass != DigestClassTokenVerifier {
		return VerifierComparisonResult{}, &VerifierComparisonError{Kind: ErrInvalidVerifierComparison, Class: expectedClass}
	}
	if computed.Class() == "" || len(computed.bytes) == 0 {
		return VerifierComparisonResult{}, &VerifierComparisonError{Kind: ErrMissingComputedDigest, Class: expectedClass}
	}
	if computed.Class() != expectedClass {
		return VerifierComparisonResult{}, &VerifierComparisonError{Kind: ErrWrongComputedDigestClass, Class: computed.Class()}
	}
	computedBytes := computed.Bytes()
	if len(computedBytes) != VerifierDigestBytes {
		return VerifierComparisonResult{}, &VerifierComparisonError{Kind: ErrInvalidVerifierComparison, Class: expectedClass}
	}
	if len(storedDigest) == 0 {
		return VerifierComparisonResult{}, &VerifierComparisonError{Kind: ErrMissingStoredDigest, Class: expectedClass}
	}
	if len(storedDigest) != VerifierDigestBytes {
		return VerifierComparisonResult{}, &VerifierComparisonError{Kind: ErrMalformedStoredDigest, Class: expectedClass}
	}

	return VerifierComparisonResult{
		class:   expectedClass,
		matched: hmac.Equal(computedBytes, storedDigest),
	}, nil
}
