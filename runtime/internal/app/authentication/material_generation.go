package authentication

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const (
	RawSecretMaterialBytes     = 32
	EncodedSecretMaterialChars = 43
)

type MaterialKind string

const (
	MaterialKindDeviceCredential MaterialKind = "device_credential"
	MaterialKindAccessToken      MaterialKind = "access_token"
)

var (
	ErrMissingRandomSource      = errors.New("authentication material generation: missing random source")
	ErrRandomReadFailed         = errors.New("authentication material generation: random read failed")
	ErrInvalidGeneratedMaterial = errors.New("authentication material generation: invalid generated material")
)

type MaterialGenerationError struct {
	Kind         error
	MaterialKind MaterialKind
	Err          error
}

func (e *MaterialGenerationError) Error() string {
	if e == nil {
		return ""
	}
	if e.MaterialKind == "" {
		return e.Kind.Error()
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.MaterialKind)
}

func (e *MaterialGenerationError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *MaterialGenerationError) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.Kind, target) || errors.Is(e.Err, target)
}

type GeneratedSecretMaterial struct {
	kind MaterialKind
	raw  []byte
	text string
}

func (m GeneratedSecretMaterial) Kind() MaterialKind {
	return m.kind
}

func (m GeneratedSecretMaterial) RawBytes() []byte {
	return cloneBytes(m.raw)
}

func (m GeneratedSecretMaterial) Text() string {
	return m.text
}

func GenerateDeviceCredentialMaterial(random io.Reader) (GeneratedSecretMaterial, error) {
	return generateSecretMaterial(MaterialKindDeviceCredential, random)
}

func GenerateAccessTokenMaterial(random io.Reader) (GeneratedSecretMaterial, error) {
	return generateSecretMaterial(MaterialKindAccessToken, random)
}

func generateSecretMaterial(kind MaterialKind, random io.Reader) (GeneratedSecretMaterial, error) {
	if random == nil {
		return GeneratedSecretMaterial{}, &MaterialGenerationError{
			Kind:         ErrMissingRandomSource,
			MaterialKind: kind,
		}
	}

	raw := make([]byte, RawSecretMaterialBytes)
	if _, err := io.ReadFull(random, raw); err != nil {
		return GeneratedSecretMaterial{}, &MaterialGenerationError{
			Kind:         ErrRandomReadFailed,
			MaterialKind: kind,
			Err:          err,
		}
	}

	if isWeakRepeatedKey(raw) {
		return GeneratedSecretMaterial{}, &MaterialGenerationError{
			Kind:         ErrInvalidGeneratedMaterial,
			MaterialKind: kind,
		}
	}

	return GeneratedSecretMaterial{
		kind: kind,
		raw:  cloneBytes(raw),
		text: base64.RawURLEncoding.EncodeToString(raw),
	}, nil
}
