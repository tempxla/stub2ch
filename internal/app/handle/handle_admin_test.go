package handle

import (
	// "github.com/tempxla/stub2ch/configs/app/admincfg"
	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	// "net/url"
	// "strings"
	"testing"
)

func TestAuthenticate(t *testing.T) {

}

func TestAuthenticate_CookieMissing(t *testing.T) {

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/_admin/", nil)

	router := NewBoardRouter(nil)
	router.ServeHTTP(writer, request)

	if writer.Code != 403 {
		t.Errorf("code: %v ", writer.Code)
	}
}

func TestAuthenticate_WrongPass(t *testing.T) {

}
