package main

import (
	E "./entity"
	"./service"
	"./testutil"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleBbsCgi_404(t *testing.T) {

	router := httprouter.New()
	router.GET("/:board/bbs.cgi", handleBbsCgi)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test1/bbs.cgi", nil)
	router.ServeHTTP(writer, request)

	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleBbsCgi_200(t *testing.T) {

	router := httprouter.New()
	router.GET("/:board/bbs.cgi", handleBbsCgi)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test/bbs.cgi", nil)
	router.ServeHTTP(writer, request)

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
	sv := service.NewBoardService(repo)

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
	sv := service.NewBoardService(repo)

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
					ThreadKey:    "111",
					ThreadTitle:  "XXX",
					MessageCount: 100,
					LastFloat:    time.Now().Add(time.Duration(2) * time.Hour),
				},
				E.Subject{
					ThreadKey:    "222",
					ThreadTitle:  "YYY",
					MessageCount: 200,
					LastFloat:    time.Now().Add(time.Duration(3) * time.Hour),
				},
				E.Subject{
					ThreadKey:    "333",
					ThreadTitle:  "ZZZ",
					MessageCount: 300,
					LastFloat:    time.Now().Add(time.Duration(1) * time.Hour),
				},
			}},
		},
	}
	sv := service.NewBoardService(repo)

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
				E.Subject{
					ThreadKey:    "111",
					ThreadTitle:  "XXX",
					MessageCount: 100,
					LastFloat:    time.Now().Add(time.Duration(2) * time.Hour),
				},
				E.Subject{
					ThreadKey:    "222",
					ThreadTitle:  "YYY",
					MessageCount: 200,
					LastFloat:    time.Now().Add(time.Duration(3) * time.Hour),
				},
				E.Subject{
					ThreadKey:    "333",
					ThreadTitle:  "ZZZ",
					MessageCount: 300,
					LastFloat:    time.Now().Add(time.Duration(1) * time.Hour),
				},
			}},
		},
	}
	sv := service.NewBoardService(repo)

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
