package handle

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/configs/app/admincfg"
	"github.com/tempxla/stub2ch/internal/app/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params, sv *service.BoardService) {
		fmt.Fprintf(w, "OK")
		return
	}

	sv, _ := service.DefaultBoardService()

	router := httprouter.New()
	router.POST("/test/_admin/", injectService(sv)(authenticate(handleOK)))

	passphrase, err := ioutil.ReadFile("/tmp/pass_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	base64Sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	sid, err := sv.Admin.Login(string(passphrase), string(base64Sig))
	if err != nil {
		t.Errorf("setup failed. %v", err)
	}

	// Exercise
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/", nil)
	request.AddCookie(&http.Cookie{
		Name:  admincfg.LOGIN_COOKIE_NAME,
		Value: sid,
	})
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	if body := writer.Body.String(); body != "OK" {
		t.Errorf("%v", body)
	}
}

func TestAuthenticate_CookieMissing(t *testing.T) {

	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params, sv *service.BoardService) {
		fmt.Fprintf(w, "OK")
		return
	}

	sv, _ := service.DefaultBoardService()

	router := httprouter.New()
	router.POST("/test/_admin/", injectService(sv)(authenticate(handleOK)))

	// Exercise
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/", nil)
	// cookie is missing
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 403 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if body == "OK" {
		t.Error("body is OK.")
	}
}

func TestAuthenticate_WrongSession(t *testing.T) {
	// Setup
	handleOK := func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params, sv *service.BoardService) {
		fmt.Fprintf(w, "OK")
		return
	}

	sv, _ := service.DefaultBoardService()

	router := httprouter.New()
	router.POST("/test/_admin/", injectService(sv)(authenticate(handleOK)))

	// Exercise
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/", nil)
	request.AddCookie(&http.Cookie{
		Name:  admincfg.LOGIN_COOKIE_NAME,
		Value: "WRONG COOKIE VALUE",
	})
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 403 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if body == "OK" {
		t.Error("body is OK.")
	}
}

func TestHandleAdminLogin(t *testing.T) {
	// Setup
	passphrase, err := ioutil.ReadFile("/tmp/pass_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	base64Sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/login", nil)
	request.PostForm = make(map[string][]string)
	request.PostForm.Add(admincfg.LOGIN_PASSPHRASE_PARAM, string(passphrase))
	request.PostForm.Add(admincfg.LOGIN_SIGNATURE_PARAM, string(base64Sig))

	sv, _ := service.DefaultBoardService()
	router := NewBoardRouter(sv)

	// Exercise
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	if body := writer.Body.String(); !strings.Contains(body, "Admin Console") {
		t.Errorf("%v", body)
	}

	// Cookie TODO

}

func TestHandleAdminLogin_MissingPassphrase(t *testing.T) {
	// Setup

	base64Sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/login", nil)
	request.PostForm = make(map[string][]string)
	//request.PostForm.Add(admincfg.LOGIN_PASSPHRASE_PARAM, string(passphrase)) missing
	request.PostForm.Add(admincfg.LOGIN_SIGNATURE_PARAM, string(base64Sig))

	sv, _ := service.DefaultBoardService()
	router := NewBoardRouter(sv)

	// Exercise
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 403 {
		t.Errorf("Response code is %v", writer.Code)
	}
	if body := writer.Body.String(); strings.Contains(body, "Admin Console") {
		t.Errorf("%v", body)
	}

}

func TestHandleAdminLogin_MissingSignature(t *testing.T) {
	// Setup
	passphrase, err := ioutil.ReadFile("/tmp/pass_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/login", nil)
	request.PostForm = make(map[string][]string)
	request.PostForm.Add(admincfg.LOGIN_PASSPHRASE_PARAM, string(passphrase))
	//request.PostForm.Add(admincfg.LOGIN_SIGNATURE_PARAM, string(base64Sig)) missing

	sv, _ := service.DefaultBoardService()
	router := NewBoardRouter(sv)

	// Exercise
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 403 {
		t.Errorf("Response code is %v", writer.Code)
	}

	if body := writer.Body.String(); strings.Contains(body, "Admin Console") {
		t.Errorf("%v", body)
	}
}

func TestHandleAdminLogin_WrongSignature(t *testing.T) {
	// Setup
	passphrase, err := ioutil.ReadFile("/tmp/pass_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/login", nil)
	request.PostForm = make(map[string][]string)
	request.PostForm.Add(admincfg.LOGIN_PASSPHRASE_PARAM, string(passphrase))
	request.PostForm.Add(admincfg.LOGIN_SIGNATURE_PARAM, "wrong sig")

	sv, _ := service.DefaultBoardService()
	router := NewBoardRouter(sv)

	// Exercise
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 403 {
		t.Errorf("Response code is %v", writer.Code)
	}

	if body := writer.Body.String(); strings.Contains(body, "Admin Console") {
		t.Errorf("%v", body)
	}
}
