package handle

import (
	_ "github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/internal/app/service"
	"github.com/tempxla/stub2ch/internal/app/util"
	"github.com/tempxla/stub2ch/tools/app/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
