package handle

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/internal/app/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

// トップページ表示
func TestHandleIndex(t *testing.T) {
	// Setup
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if err := indexTmpl.Execute(writer, nil); err != nil {
		t.Errorf("Error executing template: %v", err)
	}
}

func TestProtect(t *testing.T) {
	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "OK")
		return
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	// Exercise
	router := httprouter.New()
	router.GET("/", protect(true)(handleOK))
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 503 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if body == "OK" {
		t.Error("body is OK.")
	}
}

func TestProtect_NotProtect(t *testing.T) {
	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "OK")
		return
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	// Exercise
	router := httprouter.New()
	router.GET("/", protect(false)(handleOK))
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if body != "OK" {
		t.Errorf("body is %v", body)
	}
}
