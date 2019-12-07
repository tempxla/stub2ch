package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequireParam_ok(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{"v1"},
		},
	}

	// Exercise
	s, err := requireParam(r, "p1")

	// Verify
	if err != nil {
		t.Errorf("err %v", err)
	}
	if s != "v1" {
		t.Errorf("value: %v", s)
	}
}

func TestRequireParam_missing(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{"v1"},
		},
	}

	// Exercise
	_, err := requireParam(r, "p2")

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
}

func TestRequireParam_empty(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{},
		},
	}

	// Exercise
	_, err := requireParam(r, "p1")

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
}

func TestRequireParam_many(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{"v1", "v2"},
		},
	}

	// Exercise
	_, err := requireParam(r, "p1")

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
}

func TestRequire(t *testing.T) {
	// Setup
	r := &http.Request{
		PostForm: map[string][]string{
			"p1": []string{"v1"},
		},
	}

	// Exercise
	s, err := require(r, "p1")()

	// Verify
	if err != nil {
		t.Errorf("err %v", err)
	}
	if s != "v1" {
		t.Errorf("value: %v", s)
	}
}

func TestNotEmpty(t *testing.T) {
	if _, err := notEmpty(""); err == nil {
		t.Error("err is nil")
	}
	s, err := notEmpty("s1")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if s != "s1" {
		t.Errorf("value: %v", s)
	}
}

func TestNotBlank(t *testing.T) {
	if _, err := notBlank(" 　\n\r\t\v"); err == nil {
		t.Error("err is nil")
	}
	s, err := notBlank(" s1 ")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if s != " s1 " {
		t.Errorf("value: %v", s)
	}
}

func TestBetweenStr(t *testing.T) {
	if _, err := betweenStr("bbb", "ddd")("eee"); err == nil {
		t.Error("err is nil")
	}
	if _, err := betweenStr("bbb", "ddd")("aaa"); err == nil {
		t.Error("err is nil")
	}
	xs := [...]string{"bbb", "ccc", "ddd", "bbbb"}
	for _, x := range xs {
		s, err := betweenStr("bbb", "ddd")(x)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		if s != x {
			t.Errorf("value: %v", s)
		}
	}
}

func TestProcessParam(t *testing.T) {
	// case 1
	f := func() (string, error) { return "s", fmt.Errorf("errrrrr") }

	if _, err := processParam(f); err == nil {
		t.Error("case 1: err is nil")
	}

	// case 2
	g := func() (string, error) { return "s", nil }

	s, err := processParam(g)
	if err != nil {
		t.Errorf("case 2 err: %v", err)
	}
	if s != "s" {
		t.Errorf("case 2 value: %v", s)
	}

	// case 3
	h1 := func(s string) (string, error) { return s + "t", nil }
	h2 := func(s string) (string, error) { return s + "u", fmt.Errorf("errhhhh") }
	if _, err := processParam(g, h1, h2); err == nil {
		t.Error("case 3: err is nil")
	}

	// case 4
	h3 := func(s string) (string, error) { return s + "u", nil }
	s, err = processParam(g, h1, h3)
	if err != nil {
		t.Errorf("case 4 err: %v", err)
	}
	if s != "stu" {
		t.Errorf("case 4 value: %v", s)
	}
}
