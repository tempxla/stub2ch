package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
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
