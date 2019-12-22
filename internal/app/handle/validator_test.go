package handle

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequireOne_ok(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{"v1"},
		},
	}

	// Exercise
	s, err := requireOne(r, "p1")()

	// Verify
	if err != nil {
		t.Errorf("err %v", err)
	}
	if s != "v1" {
		t.Errorf("value: %v", s)
	}
}

func TestRequireOne_missing(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{"v1"},
		},
	}

	// Exercise
	_, err := requireOne(r, "p2")()

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
}

func TestRequireOne_empty(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{},
		},
	}

	// Exercise
	_, err := requireOne(r, "p1")()

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
}

func TestRequireOne_many(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{"v1", "v2"},
		},
	}

	// Exercise
	_, err := requireOne(r, "p1")()

	// Verify
	if err == nil {
		t.Errorf("err is nil")
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
