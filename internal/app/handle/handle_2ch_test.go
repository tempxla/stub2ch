package handle

import (
	"github.com/tempxla/stub2ch/internal/app/service"
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
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test1/bbs.cgi", nil)

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
	var repo service.BoardRepository
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
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("書き込む")},
		"bbs":     []string{"news4test"},
		"key":     []string{"1234567890"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("上記全てを承諾して書き込む")},
		"bbs":     []string{"news4test"},
		"key":     []string{"1234567890"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("新規スレッド作成")},
		"bbs":     []string{"news4test"},
		"subject": []string{"AAAA"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
	request.PostForm = map[string][]string{
		"submit":  []string{util.UTF8toSJISString("上記全てを承諾して書き込む")},
		"bbs":     []string{"news4test"},
		"subject": []string{"AAAA"},
		"time":    []string{"1"},
		"FROM":    []string{"xxxx"},
		"mail":    []string{"sage"},
		"MESSAGE": []string{util.UTF8toSJISString("書き")},
	}
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	params := []map[string]string{
		// not 400
		// map[string]string{
		// 	"bbs":     "news4test",
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
			"bbs":     "news4test",
			"key":     "12345678901", // too long
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4test",
			"key":     "1234567890",
			"time":    "", // empty
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":  "news4test",
			"key":  "1234567890",
			"time": "1",
			// "FROM":    "xxxx", // missing
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":  "news4test",
			"key":  "1234567890",
			"time": "1",
			"FROM": "xxxx",
			// "mail":    "yyyy", // missing
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4test",
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
	var repo service.BoardRepository
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4test",
		"key":     "1234567890",
		"time":    "1",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)

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
		"bbs":     "news4test",
		"key":     "1234567890",
		"time":    "1",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
		"bbs":     "news4test",
		"key":     "1234567890",
		"time":    "1",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	var repo service.BoardRepository
	var sysEnv service.BoardEnvironment
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	params := []map[string]string{
		// not 400
		// {
		// 	"bbs":     "news4test",
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
			"bbs":     "news4test",
			"subject": " ", // blank
			"time":    "1",
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4test",
			"subject": "AAAAA",
			"time":    "", // empty
			"FROM":    "xxxx",
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4test",
			"subject": "AAAAA",
			"time":    "1",
			// "FROM":    "xxxx", // missing
			"mail":    "yyyy",
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4test",
			"subject": "AAAAA",
			"time":    "1",
			"FROM":    "xxxx",
			// "mail":    "yyyy", // missing
			"MESSAGE": "aaaa",
		},
		{
			"bbs":     "news4test",
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
	var repo service.BoardRepository
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4test",
		"time":    "1",
		"subject": "AAAAA",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)

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
		"bbs":     "news4test",
		"time":    "1",
		"subject": "AAAAA",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
		{},
	},
	)
	sysEnv := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv))

	param := map[string]string{
		"bbs":     "news4test",
		"time":    "1",
		"subject": "AAAAA",
		"FROM":    "xxxx",
		"mail":    "yyyy",
		"MESSAGE": "aaaa",
	}

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test/bbs.cgi", nil)
	request.AddCookie(&http.Cookie{Name: "PON", Value: "1.1.1.1"})
	request.AddCookie(&http.Cookie{Name: "yuki", Value: "akari"})

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
	request, _ := http.NewRequest("GET", "/news4test/dat/123.dat", nil)

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
	request, _ := http.NewRequest("GET", "/news4test/dat/999.dat", nil)

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleSubjectTxt_200(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
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
	request, _ := http.NewRequest("GET", "/news4test/subject.txt", nil)

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
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{})
	env := &service.SysEnv{
		StartedTime: time.Now(),
	}
	sv := service.NewBoardService(service.RepoConf(repo), service.EnvConf(env))

	// request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/news4test2/subject.txt", nil)

	// Exercise
	router := NewBoardRouter(sv)
	router.ServeHTTP(writer, request)

	// Verify
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
