package encryption

import (
	"errors"
	"slices"
)

type EncryptionAlgorithm string

const (
	EncryptionAlgorithmRSA     EncryptionAlgorithm = "RSA"
	EncryptionAlgorithmED25519 EncryptionAlgorithm = "ED25519"
)

// KeyGenerator defines the interface for key generation
type KeyGenerator interface {
	GenerateKeys() (privateKey string, publicKey string, err error)
}

// KeyGeneratorFactory creates a KeyGenerator based on the algorithm
func KeyGeneratorFactory(algo EncryptionAlgorithm, args ...interface{}) (KeyGenerator, error) {
	var bits int
	if len(args) > 0 {
		if b, ok := args[0].(int); ok {
			bits = b
		} else {
			return nil, errors.New("invalid argument type for bits")
		}
	}

	switch algo {
	case EncryptionAlgorithmRSA:
		if !slices.Contains(rsaBitSizeParams, bits) {
			return nil, errors.New("invalid bit size for RSA")
		}
		return &rsaKeyGenerator{Bits: bits}, nil
	case EncryptionAlgorithmED25519:
		return &ed25519KeyGenerator{}, nil
	default:
		return nil, errors.New("unsupported algorithm")
	}
}
