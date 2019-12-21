package handle

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// トップページ表示
func TestHandleIndex(t *testing.T) {
	// Setup
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	// Exercise
	router := NewBoardRouter(nil)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if !strings.Contains(body, "やあ （´・ω・｀)") {
		t.Errorf("body is %v", body)
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

func TestHandleTestDir(t *testing.T) {
	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "OK")
		return
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test/ok", nil)

	// Exercise
	router := httprouter.New()
	router.GET("/:board/ok", handleTestDir(handleOK))
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

func TestHandleTestDir_NotTestDir(t *testing.T) {
	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "OK")
		return
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test1/ok", nil)

	// Exercise
	router := httprouter.New()
	router.GET("/:board/ok", handleTestDir(handleOK))
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if body == "OK" {
		t.Errorf("body is %v", body)
	}
}

func TestHandleParseForm(t *testing.T) {
	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		okParam := r.PostForm["OK_PARAM"]
		if len(okParam) > 0 {
			fmt.Fprintf(w, okParam[0]) // OK
		}
		return
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.PostForm = make(map[string][]string)
	request.PostForm.Add("OK_PARAM", "OK")

	// Exercise
	router := httprouter.New()
	router.GET("/", handleParseForm(handleOK))
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
