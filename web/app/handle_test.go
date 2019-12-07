package main

import (
	E "./entity"
	"./service"
	"./testutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// トップページ表示
func TestHandleIndex(t *testing.T) {
	// Setup
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(repo, sysEnv)

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	// Exercise
	router := newBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if err := indexTmpl.Execute(writer, nil); err != nil {
		t.Errorf("Error executing template: %v", err)
	}
}

// bbs.cgi がない
func TestHandleBbsCgi_404(t *testing.T) {
	// Setup
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(repo, sysEnv)

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test1/bbs.cgi", nil)

	// Exercise
	router := newBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

// bbs.cgi へのGET
func TestHandleBbsCgi_200(t *testing.T) {
	// Setup
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(repo, sysEnv)

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)

	// Exercise
	router := newBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
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

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test/dat/123.dat", nil)

	// Exercise
	router := newBoardRouter(sv)
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

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test/dat/999.dat", nil)

	// Exercise
	router := newBoardRouter(sv)
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

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test/subject.txt", nil)

	// Exercise
	router := newBoardRouter(sv)
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

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test2/subject.txt", nil)

	// Exercise
	router := newBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
