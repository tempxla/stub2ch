package handle

import (
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/util"
	"html"
	"net/http"
	"net/http/httptest"
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

func TestTrimWhitespace(t *testing.T) {
	tests := []struct {
		arg, want string
	}{
		{"\t\n abc \t\n ", "abc"},
	}
	for _, tt := range tests {
		value, err := trimWhitespace(tt.arg)
		if value != tt.want || err != nil {
			t.Errorf("trimWhitespace(%v) = (%v, %v), want: %v", tt.arg, value, err, tt.want)
		}
	}
}

func TestProcess(t *testing.T) {
	src := func() (string, error) { return "", nil }
	// sample func
	errorFunc1 := func(s string) (string, error) { return "", fmt.Errorf("errrrrr") }
	successFuncS := func(s string) (string, error) { return "s", nil }
	successFuncT := func(s string) (string, error) { return s + "t", nil }
	errorFunc2 := func(s string) (string, error) { return "", fmt.Errorf("errhhhh") }

	tests := []struct {
		funcs []func(string) (string, error)
		want  string
		err   error
	}{
		{[]func(string) (string, error){successFuncS}, "s", nil},
		{[]func(string) (string, error){successFuncS, successFuncT}, "st", nil},
		{[]func(string) (string, error){errorFunc1}, "", fmt.Errorf("errrrrr")},
		{[]func(string) (string, error){successFuncS, errorFunc2}, "", fmt.Errorf("errhhhh")},
	}

	for i, tt := range tests {
		value, err := process(src, tt.funcs...)
		if value != tt.want || (err == nil && tt.err != nil || err != nil && tt.err == nil) {
			t.Errorf("%d: process(src, funcs) = (%v, %v), want: %v", i, value, err, tt.want)
		}
	}
}

func TestRequireBoardName(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s, ok := requireBoardName(w, r)
		if !ok {
			return
		}
		fmt.Fprint(w, s)
	})

	tests := []struct {
		param string
		code  int
		body  string
	}{
		{"news4test", 200, "news4test"},
		{"01234567890", 400, ""},
	}

	for i, tt := range tests {
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/", nil)
		request.PostForm = map[string][]string{"bbs": {tt.param}}
		mux.ServeHTTP(writer, request)

		if writer.Code != tt.code {
			t.Errorf("%d: wirte.Code = %d, want: %d", i, writer.Code, tt.code)
		}
		if tt.code == 200 && writer.Body.String() != tt.body {
			t.Errorf("%d: writer.Body.String() = %s, want: %s", i, writer.Body.String(), tt.body)
		}
	}
}
