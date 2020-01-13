package handle

import (
	"fmt"
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
		if value != tt.want {
			t.Errorf("%d: requireOne(r, %v) = %v, want: %v", i, tt.arg, value, tt.want)
		}
		if tt.err == nil && err != nil || tt.err != nil && err == nil {
			t.Errorf("%d: requireOne(r, %v) = %v \n"+
				"tt.err = %v, err = %v",
				i, tt.arg, value,
				tt.err, err)
		}
	}
}

func TestNotEmpty(t *testing.T) {
	// error case
	if _, err := notEmpty(""); err == nil {
		t.Error("err is nil")
	}
	// not error
	s, err := notEmpty("s1")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if s != "s1" {
		t.Errorf("value: %v", s)
	}
}

func TestNotBlank(t *testing.T) {
	// error case
	if _, err := notBlank(" 　\n\r\t\v"); err == nil {
		t.Error("err is nil")
	}
	// not error
	s, err := notBlank(" s1 ")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if s != " s1 " {
		t.Errorf("value: %v", s)
	}
}

func TestBetween(t *testing.T) {
	// error case
	if _, err := between("bbb", "ddd")("eee"); err == nil {
		t.Error("err is nil")
	}
	if _, err := between("bbb", "ddd")("aaa"); err == nil {
		t.Error("err is nil")
	}
	// not error
	xs := [...]string{"bbb", "ccc", "ddd", "bbbb"}
	for _, x := range xs {
		s, err := between("bbb", "ddd")(x)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		if s != x {
			t.Errorf("value: %v", s)
		}
	}
}

func TestMaxLen(t *testing.T) {
	// error case
	if _, err := maxLen(4)("12345"); err == nil {
		t.Error("err is nil")
	}
	// not error
	s, err := maxLen(4)("1234")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if s != "1234" {
		t.Errorf("value: %v", s)
	}
}

func TestDelBadChar(t *testing.T) {
	rs := []rune{'a', '\n', 'b', '\t', 'c', '[', '\u0000', '\u000B', '\u001F', '\u007F', ']'}

	act, _ := delBadChar(string(rs))
	exp := "a\nb\tc[    ]"

	if act != exp {
		t.Errorf("\nact: %#v\nexp: %#v", act, exp)
	}
}

func TestTrip(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"名無し", "名無し"},
		{"名無し◆ABC◆XYZ", "名無し◇ABC◇XYZ"},
		{"名無し##9CA39C423D4881A6..", "名無し </b>◆moussy./hk <b>"},
		{"名無し<>", "名無し&lt;&gt;"},
	}
	for _, tt := range tests {
		actual, _ := trip(tt.name)
		if actual != tt.expected {
			t.Errorf("%v : %v", tt.name, actual)
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
