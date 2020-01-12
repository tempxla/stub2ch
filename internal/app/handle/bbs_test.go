package handle

import (
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/service"
	"github.com/tempxla/stub2ch/internal/app/service/repository"
	"github.com/tempxla/stub2ch/internal/app/util"
	"github.com/tempxla/stub2ch/tools/app/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// bbs.cgi がない
func TestHandleBbsCgi_404(t *testing.T) {
	// Setup
	var repo repository.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test1/bbs.cgi", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

// bbs.cgi
func TestHandleBbsCgi_MissingSubmit(t *testing.T) {
	// Setup
	var repo repository.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 400 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleBbsCgi_WrongSubmit(t *testing.T) {
	// Setup
	var repo repository.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.PostForm = make(map[string][]string)
	request.PostForm.Add("submit", "カキカキ")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	txt := writer.Body.String()
	if txt != util.UTF8toSJISString("SJISで書いてね？") {
		t.Errorf("actual: %v", txt)
	}
}

func TestHandleBbsCgi_writeDatOK(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "1234567890",
			ThreadTitle:  "XXXX",
			MessageCount: 1,
			LastModified: time.Now().Add(time.Duration(-1) * time.Hour),
			Dat:          "1行目",
		},
	},
	)
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("書き込む")},
		"bbs":     []string{"news4vip"},
		"key":     []string{"1234567890"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>書きこみました。</title>") {
		t.Errorf("NOT write_dat_done.html : %v", body)
	}
}

func TestHandleBbsCgi_writeDatOKthroughConfirm(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "1234567890",
			ThreadTitle:  "XXXX",
			MessageCount: 1,
			LastModified: time.Now().Add(time.Duration(-1) * time.Hour),
			Dat:          "1行目",
		},
	},
	)
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("上記全てを承諾して書き込む")},
		"bbs":     []string{"news4vip"},
		"key":     []string{"1234567890"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>書きこみました。</title>") {
		t.Errorf("NOT write_dat_done.html : %v", body)
	}
}

func TestHandleBbsCgi_createThreadOK(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{},
	},
	)
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("新規スレッド作成")},
		"bbs":     []string{"news4vip"},
		"subject": []string{"AAAA"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>書きこみました。</title>") {
		t.Errorf("NOT write_dat_done.html : %v", body)
	}
}

func TestHandleBbsCgi_createThreadOKthroughConfirm(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{},
	},
	)
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("上記全てを承諾して書き込む")},
		"bbs":     []string{"news4vip"},
		"subject": []string{"AAAA"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>書きこみました。</title>") {
		t.Errorf("NOT write_dat_done.html : %v", body)
	}
}

// パラメータ不備
// 本当は「ERROR: 送られてきたデータが壊れています」ページが返されると思う
func TestWriteDat_400(t *testing.T) {
	// Setup
	var repo repository.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	params := []map[string]string{
		// not 400
		// map[string]string{
		// 	"bbs":     "news4vip",
		// 	"key":     "1234567890",
		// 	"time":    "1",
		// 	"FROM":    "xxxx",
		// 	"mail":    "yyyy",
		// 	"MESSAGE": "aaaa",
		// },
		{
			"bbs":     "12345678901", // too long
			"key":     "1234567890",
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"key":     "12345678901", // too long
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"key":     "1234567890",
			"time":    "", // empty
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":  "news4vip",
			"key":  "1234567890",
			"time": "1",
			// "FROM":    "xxxx", // missing
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":  "news4vip",
			"key":  "1234567890",
			"time": "1",
			"FROM": "xxxx",
			// "mail":    "yyyy", // missing
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"key":     "1234567890",
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": " ", // balnk
		},
	}

	// Exercise
	for i, param := range params {
		// request
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
		request.PostForm = make(map[string][]string)
		for k, v := range param {
			request.PostForm.Add(k, v)
		}
		handleWriteDat(writer, request, sv)

		// Verify
		if writer.Code != 400 {
			t.Errorf("case %d . Response code is %v", i, writer.Code)
		}
	}
}

func TestWriteDat_CookieMissing(t *testing.T) {

	// Setup
	var repo repository.BoardRepository
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4vip",
		"key":     "1234567890",
		"time":    "1",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)

	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	request.PostForm = make(map[string][]string)
	for k, v := range param {
		request.PostForm.Add(k, v)
	}
	handleWriteDat(writer, request, sv)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// header
	cookieCount := 0
	for k, vs := range writer.HeaderMap {
		if k == "Set-Cookie" {
			for _, v := range vs {
				if strings.HasPrefix(v, "PON=") {
					cookieCount++
				} else if strings.HasPrefix(v, "yuki=akari") {
					cookieCount++
				}
			}
		}
	}
	if cookieCount != 2 {
		t.Errorf("header: %v", writer.HeaderMap)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>■ 書き込み確認 ■</title>") {
		t.Errorf("NOT write_dat_confirm.html: %v", body)
	}
}

func TestWriteDat_NotFound(t *testing.T) {
	// Setup
	repo := testutil.EmptyBoardStub()
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4vip",
		"key":     "1234567890",
		"time":    "1",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	request.PostForm = make(map[string][]string)
	for k, v := range param {
		request.PostForm.Add(k, v)
	}
	handleWriteDat(writer, request, sv)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "ERROR: 該当するスレッドがありません。") {
		t.Errorf("NOT write_dat_not_found.html : %v", body)
	}
}

func TestWriteDat_Done(t *testing.T) {
	// Setup
	dat := "1行目\n2行目"
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "1234567890",
			ThreadTitle:  "XXXX",
			MessageCount: 2,
			LastModified: time.Now().Add(time.Duration(-1) * time.Hour),
			Dat:          dat,
		},
	},
	)
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4vip",
		"key":     "1234567890",
		"time":    "1",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	request.PostForm = make(map[string][]string)
	for k, v := range param {
		request.PostForm.Add(k, v)
	}
	handleWriteDat(writer, request, sv)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>書きこみました。</title>") {
		t.Errorf("NOT write_dat_done.html : %v", body)
	}
}

// パラメータ不備
// 本当は「ERROR: 送られてきたデータが壊れています」ページが返されると思う
func TestCreateThread_400(t *testing.T) {
	// Setup
	var repo repository.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	params := []map[string]string{
		// not 400
		// {
		// 	"bbs":     "news4vip",
		// 	"subject": "AAAAA",
		// 	"time":    "1",
		// 	"FROM":    "xxxx",
		// 	"mail":    "yyyy",
		// 	"MESSAGE": "aaaa",
		// },
		{
			"bbs":     "12345678901", // too long
			"subject": "AAAAA",
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"subject": " ", // blank
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"subject": "AAAAA",
			"time":    "", // empty
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"subject": "AAAAA",
			"time":    "1",
			// "FROM":    "xxxx", // missing
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"subject": "AAAAA",
			"time":    "1",
			"FROM":    "xxxx",
			// "mail":    "yyyy", // missing
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4vip",
			"subject": "AAAAA",
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": " ", // balnk
		},
	}

	// Exercise
	for i, param := range params {
		// request
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
		request.PostForm = make(map[string][]string)
		for k, v := range param {
			request.PostForm.Add(k, v)
		}
		handleCreateThread(writer, request, sv)

		// Verify
		if writer.Code != 400 {
			t.Errorf("case %d . Response code is %v", i, writer.Code)
		}
	}
}

func TestCreateThread_CookieMissing(t *testing.T) {

	// Setup
	var repo repository.BoardRepository
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4vip",
		"time":    "1",
		"subject": "AAAAA",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)

	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	request.PostForm = make(map[string][]string)
	for k, v := range param {
		request.PostForm.Add(k, v)
	}
	handleCreateThread(writer, request, sv)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// header
	cookieCount := 0
	for k, vs := range writer.HeaderMap {
		if k == "Set-Cookie" {
			for _, v := range vs {
				if strings.HasPrefix(v, "PON=") {
					cookieCount++
				} else if strings.HasPrefix(v, "yuki=akari") {
					cookieCount++
				}
			}
		}
	}
	if cookieCount != 2 {
		t.Errorf("header: %v", writer.HeaderMap)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>■ 書き込み確認 ■</title>") {
		t.Errorf("NOT write_dat_confirm.html: %v", body)
	}
}

func TestCreateThread_NotFound(t *testing.T) {
	// Setup
	repo := testutil.EmptyBoardStub()
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4vip",
		"time":    "1",
		"subject": "AAAAA",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	request.PostForm = make(map[string][]string)
	for k, v := range param {
		request.PostForm.Add(k, v)
	}
	handleCreateThread(writer, request, sv)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "ERROR: XXXXXXX") {
		t.Errorf("NOT create_thread_error.html : %v", body)
	}
}

func TestCreateThread_Done(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{},
	},
	)
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4vip",
		"time":    "1",
		"subject": "AAAAA",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: request.RemoteAddr})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})
	request.Header.Add("Referer", "http://"+request.Host+"/news4vip/")

	// Exercise
	request.PostForm = make(map[string][]string)
	for k, v := range param {
		request.PostForm.Add(k, v)
	}
	handleCreateThread(writer, request, sv)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	// body
	body := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if !strings.Contains(body, "<title>書きこみました。</title>") {
		t.Errorf("NOT write_dat_done.html : %v", body)
	}
}

func TestHandleDat_200(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey: "123",
			Dat:       "1行目\n2行目",
		},
	})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4vip/dat/123.dat", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	txt := string(util.SJIStoUTF8(writer.Body.Bytes()))
	if txt != "1行目\n2行目" {
		t.Errorf("dat actual: %v", txt)
	}
}

func TestHandleDat_404(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey: "123",
			Dat:       "1行目\n2行目",
		},
	})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4vip/dat/999.dat", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleDat_IfModified_304(t *testing.T) {
	// Setup
	now := time.Now()
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "123",
			Dat:          "1行目\n2行目\n",
			LastModified: now,
		},
	})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4vip/dat/123.dat", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.Header.Add("If-Modified-Since", now.UTC().Format(http.TimeFormat))

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 304 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleDat_IfModified_416(t *testing.T) {
	// Setup
	now := time.Now()
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "123",
			Dat:          "1行目\n2行目\n",
			LastModified: now.Add(time.Duration(-1 * 24 * time.Hour)),
		},
	})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4vip/dat/123.dat", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.Header.Add("If-Modified-Since", now.UTC().Format(http.TimeFormat))
	request.Header.Add("Range", fmt.Sprintf("bytes=%d-", len(util.UTF8toSJISString("1行目\n2行目\n"))+1))

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 416 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleDat_IfModified_206(t *testing.T) {
	// Setup
	now := time.Now()
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "123",
			Dat:          "1行目\n2行目\n",
			LastModified: now.Add(time.Duration(-1 * 24 * time.Hour)),
		},
	})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4vip/dat/123.dat", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.Header.Add("If-Modified-Since", now.UTC().Format(http.TimeFormat))
	request.Header.Add("Range", fmt.Sprintf("bytes=%d-", len(util.UTF8toSJISString("1行目\n"))))

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 206 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := util.SJIStoUTF8String(writer.Body.String())
	if body != "2行目\n" {
		t.Errorf("body: %v", body)
	}
}

func TestHandleDat_IfModified_Err(t *testing.T) {
	// Setup
	now := time.Now()
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "123",
			Dat:          "1行目\n2行目\n",
			LastModified: now.Add(time.Duration(-1 * 24 * time.Hour)),
		},
	})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4vip/dat/123.dat", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")
	request.Header.Add("If-Modified-Since", now.UTC().Format(http.TimeFormat))
	request.Header.Add("Range", "bytes=1000+")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 400 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestParseDatRange(t *testing.T) {
	not := func(x bool) bool { return !x }
	id := func(x bool) bool { return x }
	tests := []struct {
		arg      string
		cond     func(bool) bool
		expected int
	}{
		{"bytes=3050-", id, 3050},
		{"-bytes=3050-", not, 3050},
		{"bytes=3050", not, 3050},
		{"bytes=b3050-", not, 3050},
	}

	for i, tt := range tests {
		a, err := parseDatRange(tt.arg)
		if tt.cond(a != tt.expected) && tt.cond(err != nil) {
			t.Errorf("case %d: %v, %v", i, a, err)
		}
	}
}

func TestHandleSubjectTxt_200(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{
		{
			ThreadKey:    "222",
			ThreadTitle:  "YYY",
			MessageCount: 200,
		},
		{
			ThreadKey:    "111",
			ThreadTitle:  "XXX",
			MessageCount: 100,
		},
		{
			ThreadKey:    "333",
			ThreadTitle:  "ZZZ",
			MessageCount: 300,
		},
	})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4vip/subject.txt", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	txt := writer.Body.String()
	if txt != "222.dat<>YYY \t (200)\n111.dat<>XXX \t (100)\n333.dat<>ZZZ \t (300)\n" {
		t.Errorf("subject.txt actual: %v", txt)
	}
}

func TestHandleSubjectTxt_404(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4vip", []testutil.ThreadStub{})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test2/subject.txt", nil)
	request.Header.Add("User-Agent", "Monazilla/1.00")

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
