package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

var rsaBitSizeParams = []int{2048, 4096}

// rsaKeyGenerator implements KeyGenerator for RSA
type rsaKeyGenerator struct {
	Bits int
}

func (r *rsaKeyGenerator) GenerateKeys() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, r.Bits)
	if err != nil {
		return "", "", err
	}

	// Encode private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Encode public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}
