package encryption

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
)

// ed25519KeyGenerator implements KeyGenerator for Ed25519
type ed25519KeyGenerator struct{}

func (e *ed25519KeyGenerator) GenerateKeys() (string, string, error) {
	// Generate Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	// Encode private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: privateKey.Seed(),
	})

	// Encode public key to PEM format
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "ED25519 PUBLIC KEY",
		Bytes: publicKey,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}
