package util

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
)

func VerifyPKCS1v15(base64PK string, sha256Digest [sha256.Size]byte, base64sig string) error {

	pubKey, err := base64.StdEncoding.DecodeString(base64PK)
	if err != nil {
		return err
	}

	pub, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return err
	}

	signature, err := base64.StdEncoding.DecodeString(base64sig)
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, sha256Digest[:], signature)
	if err != nil {
		return err
	}

	return nil
}
