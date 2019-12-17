package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/tempxla/stub2ch/configs/app/config"
	"io/ioutil"
	"testing"
)

func TestVerifyPKCS1v15(t *testing.T) {

	base64PK := config.RSA_PUBLIC

	var sha256Digest [sha256.Size]byte
	sha256byte, err := hex.DecodeString(config.ADMIN_PASSPHRASE_DIGEST)
	if err != nil {
		t.Error(err)
	}
	copy(sha256Digest[:], sha256byte)

	base64sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	err = VerifyPKCS1v15(base64PK, sha256Digest, string(base64sig))

	if err != nil {
		t.Error(err)
		t.Errorf("base64PK: %s", base64PK)
		t.Errorf("sha256Digest: %x", sha256Digest)
		t.Errorf("base64sig: %s", string(base64sig))
	}

}
