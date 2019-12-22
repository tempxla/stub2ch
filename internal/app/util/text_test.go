package util

import (
	"bytes"
	"testing"
)

func TestUTF8toSJIS(t *testing.T) {
	utf8 := "あいうえお"
	sjis := UTF8toSJIS([]byte(utf8))

	if !bytes.Equal(sjis, []byte{
		0x82, 0xa0, // あ
		0x82, 0xa2, // い
		0x82, 0xa4, // う
		0x82, 0xa6, // え
		0x82, 0xa8, // お
	}) {
		t.Errorf("%v", sjis)
	}
}

func TestSJIStoUTF8(t *testing.T) {
	sjis := []byte{
		0x82, 0xa0, // あ
		0x82, 0xa2, // い
		0x82, 0xa4, // う
		0x82, 0xa6, // え
		0x82, 0xa8, // お
	}

	utf8 := SJIStoUTF8(sjis)
	if string(utf8) != "あいうえお" {
		t.Errorf("%v", utf8)
	}
}
