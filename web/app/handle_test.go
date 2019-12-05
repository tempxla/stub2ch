package main

import (
	E "./entity"
	"./service"
	"./testutil"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleBbsCgi_404(t *testing.T) {

	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(repo, sysEnv)

	router := httprouter.New()
	router.GET("/:board/bbs.cgi", handleBbsCgi(sv))

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test1/bbs.cgi", nil)
	router.ServeHTTP(writer, request)

	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleBbsCgi_200(t *testing.T) {

	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(repo, sysEnv)

	router := httprouter.New()
	router.GET("/:board/bbs.cgi", handleBbsCgi(sv))

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test/bbs.cgi", nil)
	router.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

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

func TestHandleDat_200(t *testing.T) {
	// Setup
	repo := &testutil.BoardStub{
		DatMap: map[string]map[string]*E.DatEntity{
			"news4test": map[string]*E.DatEntity{
				"123": &E.DatEntity{
					Dat: []byte("1行目\n2行目"),
				},
			},
		},
	}
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(repo, env)

	router := httprouter.New()
	router.GET("/:board/dat/:dat", handleDat(sv))

	// Exercise
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test/dat/123.dat", nil)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	txt := writer.Body.String()
	if txt != "1行目\n2行目" {
		t.Errorf("dat actual: %v", txt)
	}
}

func TestHandleDat_400(t *testing.T) {
	// Setup
	repo := &testutil.BoardStub{
		DatMap: map[string]map[string]*E.DatEntity{
			"news4test": map[string]*E.DatEntity{
				"123": &E.DatEntity{
					Dat: []byte("1行目\n2行目"),
				},
			},
		},
	}
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(repo, env)

	router := httprouter.New()
	router.GET("/:board/dat/:dat", handleDat(sv))

	// Exercise
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test/dat/999.dat", nil)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleSubjectTxt_200(t *testing.T) {
	// Setup
	repo := &testutil.BoardStub{
		BoardMap: map[string]*E.BoardEntity{
			"news4test": &E.BoardEntity{Subjects: []E.Subject{
				E.Subject{
					ThreadKey:    "222",
					ThreadTitle:  "YYY",
					MessageCount: 200,
				},
				E.Subject{
					ThreadKey:    "111",
					ThreadTitle:  "XXX",
					MessageCount: 100,
				},
				E.Subject{
					ThreadKey:    "333",
					ThreadTitle:  "ZZZ",
					MessageCount: 300,
				},
			}},
		},
	}
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(repo, env)

	router := httprouter.New()
	router.GET("/:board/subject.txt", handleSubjectTxt(sv))

	// Exercise
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test/subject.txt", nil)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	txt := writer.Body.String()
	if txt != "222.dat<>YYY \t (200)\n111.dat<>XXX \t (100)\n333.dat<>ZZZ \t (300)" {
		t.Errorf("subject.txt actual: %v", txt)
	}
}

func TestHandleSubjectTxt_404(t *testing.T) {
	// Setup
	repo := &testutil.BoardStub{
		BoardMap: map[string]*E.BoardEntity{
			"news4test": &E.BoardEntity{Subjects: []E.Subject{
				E.Subject{},
			}},
		},
	}
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(repo, env)

	router := httprouter.New()
	router.GET("/:board/subject.txt", handleSubjectTxt(sv))

	// Exercise
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test2/subject.txt", nil)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
