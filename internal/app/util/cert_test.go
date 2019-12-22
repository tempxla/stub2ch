package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/tempxla/stub2ch/configs/app/admincfg"
	"io/ioutil"
	"testing"
)

func TestVerifyPKCS1v15(t *testing.T) {

	base64Pub := admincfg.RSA_PUBLIC

	var sha256Digest [sha256.Size]byte
	sha256byte, err := hex.DecodeString(admincfg.LOGIN_PASSPHRASE_DIGEST)
	if err != nil {
		t.Error(err)
	}
	copy(sha256Digest[:], sha256byte)

	base64Sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	err = VerifyPKCS1v15(base64Pub, sha256Digest, string(base64Sig))

	if err != nil {
		t.Error(err)
		t.Errorf("base64Pub: %s", base64Pub)
		t.Errorf("sha256Digest: %x", sha256Digest)
		t.Errorf("base64sig: %s", string(base64Sig))
	}
}

func TestVerifyPKCS1v15_Bad(t *testing.T) {

	// public key
	base64Pub := admincfg.RSA_PUBLIC
	base64Pub_not_base64 := "NOT BASE64"
	base64Pub_not_Pub := "44GG44KT44GT"

	// digest
	var sha256Digest [sha256.Size]byte
	sha256byte, err := hex.DecodeString(admincfg.LOGIN_PASSPHRASE_DIGEST)
	if err != nil {
		t.Error(err)
	}
	copy(sha256Digest[:], sha256byte)
	var sha256Digest_bad [sha256.Size]byte
	copy(sha256Digest_bad[:], []byte("AA"))

	// signature
	base64Sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}
	base64sig_not_base64 := []byte("NOT BASE64")

	// test cases
	tests := []struct {
		pub string
		dig [sha256.Size]byte
		sig []byte
	}{
		{base64Pub_not_base64, sha256Digest, base64Sig},
		{base64Pub_not_Pub, sha256Digest, base64Sig},
		{base64Pub, sha256Digest, base64sig_not_base64},
		{base64Pub, sha256Digest_bad, base64Sig},
	}

	// Verify
	for _, ts := range tests {
		err = VerifyPKCS1v15(ts.pub, ts.dig, string(ts.sig))
		if err == nil {
			t.Error("err is nil")
		}
	}
}
