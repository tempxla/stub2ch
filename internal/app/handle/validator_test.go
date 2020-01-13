package handle

import (
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/util"
	"html"
	"net/http"
	"testing"
)

func TestRequireOne(t *testing.T) {

	tests := []struct {
		param     map[string][]string
		arg, want string
		err       error
	}{
		{map[string][]string{"p1": {"v1"}}, "p1", "v1", nil},
		{map[string][]string{"p1": {"v1"}}, "p2", "", fmt.Errorf("missing")},
		{map[string][]string{"p1": {}}, "p1", "", fmt.Errorf("empty")},
		{map[string][]string{"p1": {"v1", "v2"}}, "p1", "", fmt.Errorf("too many")},
	}

	for i, tt := range tests {
		r := &http.Request{
			PostForm: tt.param,
		}
		value, err := requireOne(r, tt.arg)()
		if value != tt.want ||
			(err == nil && tt.err != nil || err != nil && tt.err == nil) {
			t.Errorf("%d: requireOne(r, %v) = (%v, %v), want: (%v, %v)",
				i, tt.arg, value, err, tt.want, tt.err)
		}
	}
}

func TestNotEmpty(t *testing.T) {

	tests := []struct {
		arg, want string
		err       error
	}{
		{"s1", "s1", nil},
		{"", "", fmt.Errorf("0 byte")},
	}

	for i, tt := range tests {
		value, err := notEmpty(tt.arg)
		if value != tt.want || (err == nil && tt.err != nil || err != nil && tt.err == nil) {
			t.Errorf("%d: notEmpty(%v) = (%v, %v), want: (%v, %v)",
				i, tt.arg, value, err, tt.want, tt.err)
		}
	}
}

func TestNotBlank(t *testing.T) {

	tests := []struct {
		arg, want string
		err       error
	}{
		{" s1 ", " s1 ", nil},
		{" 　\n\r\t\v", "", fmt.Errorf("blank")},
	}

	for i, tt := range tests {
		value, err := notBlank(tt.arg)
		if value != tt.want || (err == nil && tt.err != nil || err != nil && tt.err == nil) {
			t.Errorf("%d: notBlank(%v) = (%v, %v), want: (%v, %v)",
				i, tt.arg, value, err, tt.want, tt.err)
		}
	}
}

func TestBetween(t *testing.T) {

	tests := []struct {
		min, max  string
		arg, want string
		err       error
	}{
		{"bbb", "ddd", "bbb", "bbb", nil},
		{"bbb", "ddd", "ccc", "ccc", nil},
		{"bbb", "ddd", "ddd", "ddd", nil},
		{"bbb", "ddd", "bbbb", "bbbb", nil},
		{"bbb", "ddd", "eee", "", fmt.Errorf("between")},
		{"bbb", "ddd", "aaa", "", fmt.Errorf("between")},
	}

	for i, tt := range tests {
		value, err := between(tt.min, tt.max)(tt.arg)
		if value != tt.want || (err == nil && tt.err != nil || err != nil && tt.err == nil) {
			t.Errorf("%d: between(%v, %v)(%v) = (%v, %v), want: (%v, %v)",
				i, tt.min, tt.max, tt.arg, value, err, tt.want, tt.err)
		}
	}
}

func TestMaxLen(t *testing.T) {

	tests := []struct {
		max       int
		arg, want string
		err       error
	}{
		{4, "1234", "1234", nil},
		{4, "１２３４", "１２３４", nil},
		{4, "12345", "", fmt.Errorf("maxLen")},
		{4, "１２３４５", "", fmt.Errorf("maxLen")},
	}

	for i, tt := range tests {
		value, err := maxLen(tt.max)(tt.arg)
		if value != tt.want || (err == nil && tt.err != nil || err != nil && tt.err == nil) {
			t.Errorf("%d: maxLen(%v)(%v) = (%v, %v), want: (%v, %v)",
				i, tt.max, tt.arg, value, err, tt.want, tt.err)
		}
	}
}

func TestMaxByte(t *testing.T) {

	tests := []struct {
		max       int
		arg, want string
		err       error
	}{
		{4, "1234", "1234", nil},
		{6, "１２", "１２", nil},
		{4, "12345", "", fmt.Errorf("maxByte")},
		{6, "１２3", "", fmt.Errorf("maxByte")},
		// SJIS
		{5, util.UTF8toSJISString("ABCDE"), util.UTF8toSJISString("ABCDE"), nil},
		{5, util.UTF8toSJISString("ｱｲｳｴｵ"), util.UTF8toSJISString("ｱｲｳｴｵ"), nil},
		{6, util.UTF8toSJISString("亜細亜"), util.UTF8toSJISString("亜細亜"), nil},
		{5, util.UTF8toSJISString("ABCDEA"), "", fmt.Errorf("maxByte")},
		{5, util.UTF8toSJISString("ｱｲｳｴｵA"), "", fmt.Errorf("maxByte")},
		{6, util.UTF8toSJISString("亜細亜A"), "", fmt.Errorf("maxByte")},
	}

	for i, tt := range tests {
		value, err := maxByte(tt.max)(tt.arg)
		if value != tt.want || (err == nil && tt.err != nil || err != nil && tt.err == nil) {
			t.Errorf("%d: maxByte(%v)(%v) = (%v, %v), want: (%v, %v)",
				i, tt.max, tt.arg, value, err, tt.want, tt.err)
		}
	}
}

func TestDelBadChar(t *testing.T) {
	rs := []rune{'a', '\n', 'b', '\t', 'c', '[', '\u0000', '\u000B', '\u001F', '\u007F', ']'}

	act, err := delBadChar(string(rs))
	exp := "a\nb\tc[    ]"

	if act != exp || err != nil {
		t.Errorf("\nact: %#v\nexp: %#v, err = %v", act, exp, err)
	}
}

func TestHtmlUnescapeString(t *testing.T) {
	tests := []struct {
		arg, want string
	}{
		{`<>"'&`, html.UnescapeString(`<>"'&`)},
	}
	for _, tt := range tests {
		value, err := htmlUnescapeString(tt.arg)
		if value != tt.want || err != nil {
			t.Errorf("htmlUnescapeString(%v) = (%v, %v), want: %v", tt.arg, value, err, tt.want)
		}
	}
}

func TestSjisToUtf8String(t *testing.T) {
	tests := []struct {
		arg, want string
	}{
		{util.UTF8toSJISString("あいうえお"), "あいうえお"},
	}
	for _, tt := range tests {
		value, err := sjisToUtf8String(tt.arg)
		if value != tt.want || err != nil {
			t.Errorf("sjisToUtf8String(%v) = (%v, %v), want: %v", tt.arg, value, err, tt.want)
		}
	}
}

func TestTrip(t *testing.T) {
	tests := []struct {
		name, expected string
	}{
		{"名無し", "名無し"},
		{"名無し◆ABC◆XYZ", "名無し◇ABC◇XYZ"},
		{"名無し##9CA39C423D4881A6..", "名無し </b>◆moussy./hk <b>"},
		{"名無し<>", "名無し&lt;&gt;"},
	}
	for _, tt := range tests {
		actual, _ := trip(tt.name)
		if actual != tt.expected {
			t.Errorf("trip(%v) = %v, want: %v", tt.name, actual, tt.expected)
		}
	}
}

func TestProcessParam(t *testing.T) {
	// case 1
	f := func() (string, error) { return "s", fmt.Errorf("errrrrr") }

	if _, err := process(f); err == nil {
		t.Error("case 1: err is nil")
	}

	// case 2
	g := func() (string, error) { return "s", nil }

	s, err := process(g)
	if err != nil {
		t.Errorf("case 2 err: %v", err)
	}
	if s != "s" {
		t.Errorf("case 2 value: %v", s)
	}

	// case 3
	h1 := func(s string) (string, error) { return s + "t", nil }
	h2 := func(s string) (string, error) { return s + "u", fmt.Errorf("errhhhh") }
	if _, err := process(g, h1, h2); err == nil {
		t.Error("case 3: err is nil")
	}

	// case 4
	h3 := func(s string) (string, error) { return s + "u", nil }
	s, err = process(g, h1, h3)
	if err != nil {
		t.Errorf("case 4 err: %v", err)
	}
	if s != "stu" {
		t.Errorf("case 4 value: %v", s)
	}
}
