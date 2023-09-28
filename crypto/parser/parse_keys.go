package parser

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParseRsaPrivateKeyFromPemStr(prvtKeyPEM []byte) (*rsa.PrivateKey, error) {
	// casting []byte(privtKeyPEM) - if args types got string
	block, _ := pem.Decode(prvtKeyPEM)
	if block == nil {
		return nil, errors.New("Failed to parse PRVT PEM block containing the key !")
	}

	// The pem file must start with ---- BEGIN RSA PRIVATE KEY
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func ParseRsaPublicKeyFromPemStr(pubKeyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		return nil, errors.New("Failed to parse PUB PEM block containing the key !")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}
