package util

import (
	"bytes"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

func UTF8toSJIS(utf8 []byte) []byte {
	buf := new(bytes.Buffer)
	w := transform.NewWriter(buf, japanese.ShiftJIS.NewEncoder())
	w.Write(utf8)
	return buf.Bytes()
}

func SJIStoUTF8(sjis []byte) []byte {
	r := bytes.NewReader(sjis)
	wrap := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	buf, _ := ioutil.ReadAll(wrap)
	return buf
}

func UTF8toSJISString(utf8 string) string {
	buf := new(bytes.Buffer)
	w := transform.NewWriter(buf, japanese.ShiftJIS.NewEncoder())
	w.Write([]byte(utf8))
	return buf.String()
}

func SJIStoUTF8String(sjis string) string {
	r := strings.NewReader(sjis)
	wrap := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	buf, _ := ioutil.ReadAll(wrap)
	return string(buf)
}
