package authentication

import (
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestGenerateDeviceCredentialMaterialUses32RandomBytes(t *testing.T) {
	source := bytesWithSeed(21, RawSecretMaterialBytes)
	reader := &countingReader{data: source}

	material, err := GenerateDeviceCredentialMaterial(reader)
	if err != nil {
		t.Fatalf("GenerateDeviceCredentialMaterial() error = %v, want nil", err)
	}

	if reader.bytesRead != RawSecretMaterialBytes {
		t.Fatalf("bytesRead = %d, want %d", reader.bytesRead, RawSecretMaterialBytes)
	}
	assertBytesEqual(t, material.RawBytes(), source)
}

func TestGenerateAccessTokenMaterialUses32RandomBytes(t *testing.T) {
	source := bytesWithSeed(31, RawSecretMaterialBytes)
	reader := &countingReader{data: source}

	material, err := GenerateAccessTokenMaterial(reader)
	if err != nil {
		t.Fatalf("GenerateAccessTokenMaterial() error = %v, want nil", err)
	}

	if reader.bytesRead != RawSecretMaterialBytes {
		t.Fatalf("bytesRead = %d, want %d", reader.bytesRead, RawSecretMaterialBytes)
	}
	assertBytesEqual(t, material.RawBytes(), source)
}

func TestGeneratedMaterialTextIsBase64URLUnpadded(t *testing.T) {
	source := bytesWithSeed(41, RawSecretMaterialBytes)

	material, err := GenerateAccessTokenMaterial(strings.NewReader(string(source)))
	if err != nil {
		t.Fatalf("GenerateAccessTokenMaterial() error = %v, want nil", err)
	}

	text := material.Text()
	if strings.Contains(text, "=") {
		t.Fatalf("Text() = %q, want no padding", text)
	}
	if strings.ContainsAny(text, "+/") {
		t.Fatalf("Text() = %q, want URL-safe alphabet", text)
	}
	if text != base64.RawURLEncoding.EncodeToString(source) {
		t.Fatalf("Text() = %q, want RawURLEncoding text", text)
	}
}

func TestGeneratedMaterialTextLengthIs43Characters(t *testing.T) {
	material, err := GenerateDeviceCredentialMaterial(strings.NewReader(string(bytesWithSeed(51, RawSecretMaterialBytes))))
	if err != nil {
		t.Fatalf("GenerateDeviceCredentialMaterial() error = %v, want nil", err)
	}

	if len(material.Text()) != EncodedSecretMaterialChars {
		t.Fatalf("len(Text()) = %d, want %d", len(material.Text()), EncodedSecretMaterialChars)
	}
}

func TestGeneratedMaterialTextRoundTripsToRawBytes(t *testing.T) {
	source := bytesWithSeed(61, RawSecretMaterialBytes)

	material, err := GenerateAccessTokenMaterial(strings.NewReader(string(source)))
	if err != nil {
		t.Fatalf("GenerateAccessTokenMaterial() error = %v, want nil", err)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(material.Text())
	if err != nil {
		t.Fatalf("DecodeString(Text()) error = %v, want nil", err)
	}
	assertBytesEqual(t, decoded, source)
}

func TestGeneratedMaterialKindIsPreserved(t *testing.T) {
	device, err := GenerateDeviceCredentialMaterial(strings.NewReader(string(bytesWithSeed(71, RawSecretMaterialBytes))))
	if err != nil {
		t.Fatalf("GenerateDeviceCredentialMaterial() error = %v, want nil", err)
	}
	token, err := GenerateAccessTokenMaterial(strings.NewReader(string(bytesWithSeed(81, RawSecretMaterialBytes))))
	if err != nil {
		t.Fatalf("GenerateAccessTokenMaterial() error = %v, want nil", err)
	}

	if device.Kind() != MaterialKindDeviceCredential {
		t.Fatalf("device.Kind() = %q, want %q", device.Kind(), MaterialKindDeviceCredential)
	}
	if token.Kind() != MaterialKindAccessToken {
		t.Fatalf("token.Kind() = %q, want %q", token.Kind(), MaterialKindAccessToken)
	}
}

func TestGeneratedMaterialRawBytesAreCopiedOnReturn(t *testing.T) {
	source := bytesWithSeed(91, RawSecretMaterialBytes)

	material, err := GenerateDeviceCredentialMaterial(strings.NewReader(string(source)))
	if err != nil {
		t.Fatalf("GenerateDeviceCredentialMaterial() error = %v, want nil", err)
	}

	first := material.RawBytes()
	first[0] = 0
	second := material.RawBytes()
	if second[0] == 0 {
		t.Fatal("RawBytes() returned mutable internal material")
	}
	assertBytesEqual(t, second, source)
}

func TestGenerateMaterialNilRandomSourceFailsClosed(t *testing.T) {
	_, err := GenerateDeviceCredentialMaterial(nil)
	assertMaterialGenerationError(t, err, ErrMissingRandomSource, MaterialKindDeviceCredential)
}

func TestGenerateMaterialRandomSourceErrorFailsClosed(t *testing.T) {
	readerErr := errors.New("reader unavailable")

	_, err := GenerateAccessTokenMaterial(errorReader{err: readerErr})

	assertMaterialGenerationError(t, err, ErrRandomReadFailed, MaterialKindAccessToken)
	if !errors.Is(err, readerErr) {
		t.Fatalf("error = %v, want wrapped reader error", err)
	}
}

func TestGenerateMaterialShortRandomReadFailsClosed(t *testing.T) {
	_, err := GenerateDeviceCredentialMaterial(strings.NewReader(string(bytesWithSeed(101, RawSecretMaterialBytes-1))))
	assertMaterialGenerationError(t, err, ErrRandomReadFailed, MaterialKindDeviceCredential)
}

func TestGenerateMaterialAllZeroGeneratedMaterialFailsClosed(t *testing.T) {
	_, err := GenerateAccessTokenMaterial(strings.NewReader(string(make([]byte, RawSecretMaterialBytes))))
	assertMaterialGenerationError(t, err, ErrInvalidGeneratedMaterial, MaterialKindAccessToken)
}

func TestGenerateMaterialRepeatedSingleByteGeneratedMaterialFailsClosed(t *testing.T) {
	_, err := GenerateDeviceCredentialMaterial(strings.NewReader(string(repeatedBytes(9, RawSecretMaterialBytes))))
	assertMaterialGenerationError(t, err, ErrInvalidGeneratedMaterial, MaterialKindDeviceCredential)
}

func TestGeneratedMaterialsAreNotConstantWithProgressingSource(t *testing.T) {
	reader := strings.NewReader(string(append(
		bytesWithSeed(111, RawSecretMaterialBytes),
		bytesWithSeed(121, RawSecretMaterialBytes)...,
	)))

	first, err := GenerateAccessTokenMaterial(reader)
	if err != nil {
		t.Fatalf("first GenerateAccessTokenMaterial() error = %v, want nil", err)
	}
	second, err := GenerateAccessTokenMaterial(reader)
	if err != nil {
		t.Fatalf("second GenerateAccessTokenMaterial() error = %v, want nil", err)
	}

	if first.Text() == second.Text() {
		t.Fatal("generated materials are constant across progressing source reads")
	}
	if equalBytes(first.RawBytes(), second.RawBytes()) {
		t.Fatal("generated raw materials are constant across progressing source reads")
	}
}

func TestMaterialGenerationErrorsDoNotIncludeRawOrEncodedMaterial(t *testing.T) {
	source := []byte("short-secret-fragment")
	encoded := base64.RawURLEncoding.EncodeToString(source)

	_, err := GenerateAccessTokenMaterial(strings.NewReader(string(source)))
	if err == nil {
		t.Fatal("GenerateAccessTokenMaterial() error = nil, want error")
	}

	message := err.Error()
	for _, allowed := range []string{ErrRandomReadFailed.Error(), string(MaterialKindAccessToken)} {
		if !strings.Contains(message, allowed) {
			t.Fatalf("error %q does not include allowed marker %q", message, allowed)
		}
	}
	for _, forbidden := range []string{string(source), encoded} {
		if forbidden != "" && strings.Contains(message, forbidden) {
			t.Fatalf("error %q contains forbidden generated material %q", message, forbidden)
		}
	}
}

func TestMaterialGenerationHelperDoesNotComputeDigestsOrCompareVerifiers(t *testing.T) {
	material, err := GenerateDeviceCredentialMaterial(strings.NewReader(string(bytesWithSeed(131, RawSecretMaterialBytes))))
	if err != nil {
		t.Fatalf("GenerateDeviceCredentialMaterial() error = %v, want nil", err)
	}

	if len(material.RawBytes()) != RawSecretMaterialBytes {
		t.Fatalf("len(RawBytes()) = %d, want %d", len(material.RawBytes()), RawSecretMaterialBytes)
	}
	if material.Text() == "" {
		t.Fatal("Text() is empty")
	}
}

type countingReader struct {
	data      []byte
	bytesRead int
}

func (r *countingReader) Read(target []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := copy(target, r.data)
	r.data = r.data[n:]
	r.bytesRead += n
	return n, nil
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}

func assertMaterialGenerationError(t *testing.T, err error, target error, kind MaterialKind) {
	t.Helper()
	assertErrorIs(t, err, target)
	var materialErr *MaterialGenerationError
	if !errors.As(err, &materialErr) {
		t.Fatalf("error = %T, want *MaterialGenerationError", err)
	}
	if materialErr.MaterialKind != kind {
		t.Fatalf("MaterialKind = %q, want %q", materialErr.MaterialKind, kind)
	}
}
