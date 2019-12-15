package handle

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/memcache"
	"github.com/tempxla/stub2ch/internal/app/service"
	. "github.com/tempxla/stub2ch/internal/app/types"
	"github.com/tempxla/stub2ch/internal/app/util"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	param_error_format = "bad parameter '%s' is: %v"
)

var (
	indexTmpl             = template.Must(template.ParseFiles(filepath.Join("web", "template", "index.html")))
	writeDatConfirmTmpl   = template.Must(template.ParseFiles(filepath.Join("web", "template", "write_dat_confirm.html")))
	writeDatNotFoundTmpl  = template.Must(template.ParseFiles(filepath.Join("web", "template", "write_dat_not_found.html")))
	writeDatDoneTmpl      = template.Must(template.ParseFiles(filepath.Join("web", "template", "write_dat_done.html")))
	createThreadErrorTmpl = template.Must(template.ParseFiles(filepath.Join("web", "template", "create_thread_error.html")))
	adminIndexTmpl        = template.Must(template.ParseFiles(filepath.Join("web", "template", "admin", "index.html")))
)

// HTTP routing
func NewBoardRouter(sv *service.BoardService) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/:board/_admin/:mode", handleTestDir(serviceChain(sv, authenticate, handleAdmin)))
	router.POST("/:board/bbs.cgi", protect(handleTestDir(serviceChain(sv, handleBbsCgi))))
	router.GET("/:board/subject.txt", protect(serviceChain(sv, handleSubjectTxt)))
	router.GET("/:board/dat/:dat", protect(serviceChain(sv, handleDat)))
	return router
}

func serviceChain(sv *service.BoardService, handles ...func(*service.BoardService) httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var boardService *service.BoardService
		if sv != nil {
			boardService = sv
		} else {
			ctx := context.Background()

			// Creates a client.
			client, err := datastore.NewClient(ctx, config.PROJECT_ID)
			if err != nil {
				log.Fatalf("Failed to create client: %v", err)
			}

			repo := &service.BoardStore{
				Context: ctx,
				Client:  client,
			}
			sysEnv := &service.SysEnv{
				StartedTime: time.Now(),
			}
			mem := &memcache.AlterMemcache{
				Context: ctx,
				Client:  client,
			}

			boardService = service.NewBoardService(service.RepoConf(repo), service.EnvConf(sysEnv), service.MemConf(mem))
		}

		// Injection
		for _, h := range handles {
			h(boardService)(w, r, ps)
		}
	}
}

func protect(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if config.KEEP_OUT {
			http.Error(w, "KEEP OUT", http.StatusServiceUnavailable) // 503
			return
		}
		h(w, r, ps)
	}
}

func authenticate(sv *service.BoardService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	}
}

func handleAdmin(sv *service.BoardService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	}
}

// トップページ表示
func handleIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := indexTmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleTestDir(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		board := ps.ByName("board")
		if board != "test" {
			http.NotFound(w, r)
			return
		}
		h(w, r, ps)
	}
}

func handleBbsCgi(sv *service.BoardService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		r.ParseForm()

		submit, err := process(requireOne(r, "submit"), url.QueryUnescape)
		if err != nil {
			http.Error(w, fmt.Sprintf(param_error_format, "submit", err), http.StatusBadRequest)
			return
		}

		switch submit {
		case "書き込む":
			// レスを書き込む
			handleWriteDat(sv, w, r)
		case "新規スレッド作成":
			// スレッドを立てる
			handleCreateThread(sv, w, r)
		case "上記全てを承諾して書き込む":
			if _, ok := r.PostForm["key"]; ok {
				handleWriteDat(sv, w, r)
			} else {
				handleCreateThread(sv, w, r)
			}
		default:
			fmt.Fprint(w, "SJISで書いてね？")
		}
	}
}

func handleWriteDat(sv *service.BoardService, w http.ResponseWriter, r *http.Request) {
	boardName, ok := requireBoardName(w, r)
	if !ok {
		return
	}
	threadKey, ok := requireThreadKey(w, r)
	if !ok {
		return
	}
	_, ok = requireTime(w, r)
	if !ok {
		return
	}
	name, ok := requireName(w, r)
	if !ok {
		return
	}
	mail, ok := requireMail(w, r)
	if !ok {
		return
	}
	message, ok := requireMessage(w, r)
	if !ok {
		return
	}
	// クッキー確認
	if executeWriteDatConfirmTmpl(w, r,
		boardName, name, mail, message, sv.StartedAt(), Nothing(), Just(threadKey)) {
		return
	}
	// 書き込み
	id := sv.ComputeId(r.RemoteAddr, boardName)
	resnum, err := sv.WriteDat(boardName, threadKey, name, mail, id, message)
	if err != nil {
		// 存在しない or dat落ち or 1001 or 容量オーバー
		executeWriteDatNotFoundTmpl(w, r, boardName, threadKey, sv.StartedAt())
		return
	}
	// 書き込み完了
	executeWriteDoneTmpl(w, r, boardName, threadKey, id, resnum, sv.StartedAt())
}

func executeWriteDoneTmpl(w http.ResponseWriter, r *http.Request,
	boardName, threadKey, id string, resnum int, startedAt time.Time) {

	w.Header().Add("Date", startedAt.UTC().Format(http.TimeFormat))
	w.Header().Add("Content-Type", "text/html; charset=Shift_JIS")
	w.Header().Add("x-Resnum", strconv.Itoa(resnum))
	//                              12345678901234567890
	mills := startedAt.Format("2006-01-02 15:04:05.000")[20:]
	w.Header().Add("x-PostDate", strconv.FormatInt(startedAt.Unix(), 10)+"."+mills)
	w.Header().Add("x-PosterID", id)
	// Body
	view := map[string]string{
		// //leia.2ch.net/test/read.cgi/poverty/1575541744/l50
		"URL": fmt.Sprintf("//%s/test/read.cgi/%s/%s/l50", r.Host, boardName, threadKey),
		"Sec": fmt.Sprintf("%f", time.Now().Sub(startedAt).Seconds()),
	}
	if err := writeDatDoneTmpl.Execute(w, view); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func executeWriteDatNotFoundTmpl(w http.ResponseWriter, r *http.Request,
	boardName, threadKey string, startedAt time.Time) {

	w.Header().Add("Content-Type", "text/html; charset=Shift_JIS")
	w.Header().Add("Date", startedAt.UTC().Format(http.TimeFormat))
	// //hebi.5ch.net/test/read.cgi/news4vip/1575543566/
	view := fmt.Sprintf("//%s/test/read.cgi/%s/%s/", r.Host, boardName, threadKey)
	if err := writeDatNotFoundTmpl.Execute(w, view); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Returns false if Cookie Found.
func executeWriteDatConfirmTmpl(w http.ResponseWriter, r *http.Request,
	boardName, name, mail, message string, startedAt time.Time, title, threadKey Maybe) bool {

	if c, err := r.Cookie("PON"); err == nil && c.Value != "" {
		if c, err := r.Cookie("yuki"); err == nil && c.Value == "akari" {
			// Cookie Found. Need not to forward Confirm page.
			return false
		}
	}
	w.Header().Add("Content-Type", "text/html; charset=Shift_JIS")
	w.Header().Add("Date", startedAt.UTC().Format(http.TimeFormat))
	// Domain属性を指定しないCookieは、Cookieを発行したホストのみに送信される
	expires := startedAt.Add(time.Duration(7*24) * time.Hour)
	w.Header().Add("Set-Cookie", fmt.Sprintf("PON=%s; expires=%s; path=/", r.RemoteAddr, expires))
	w.Header().Add("Set-Cookie", fmt.Sprintf("yuki=akari; expires=%s; path=/", expires))
	// Body
	view := map[string]string{
		"Title":     FromMaybe(title, ""),
		"Name":      name,
		"Mail":      mail,
		"Message":   message,
		"BoardName": boardName,
		"Time":      strconv.FormatInt(startedAt.Unix(), 10),
	}
	if IsJust(threadKey) {
		view["ThreadKey"] = FromJust(threadKey)
	}
	if err := writeDatConfirmTmpl.Execute(w, view); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	return true
}

func handleCreateThread(sv *service.BoardService, w http.ResponseWriter, r *http.Request) {
	boardName, ok := requireBoardName(w, r)
	if !ok {
		return
	}
	title, ok := requireTitle(w, r)
	if !ok {
		return
	}
	_, ok = requireTime(w, r)
	if !ok {
		return
	}
	name, ok := requireName(w, r)
	if !ok {
		return
	}
	mail, ok := requireMail(w, r)
	if !ok {
		return
	}
	message, ok := requireMessage(w, r)
	if !ok {
		return
	}
	// クッキー確認
	if executeWriteDatConfirmTmpl(w, r,
		boardName, name, mail, message, sv.StartedAt(), Just(title), Nothing()) {
		return
	}
	// スレ立て
	id := sv.ComputeId(r.RemoteAddr, boardName)
	threadKey, err := sv.CreateThread(boardName, name, mail, sv.StartedAt(), id, message, title)
	if err != nil {
		// スレ立て失敗
		executeCreateThreadErrorTmpl(w, r, sv.StartedAt())
		return
	}
	// 書き込み完了
	executeWriteDoneTmpl(w, r, boardName, threadKey, id, 1, sv.StartedAt())
}

func executeCreateThreadErrorTmpl(w http.ResponseWriter, r *http.Request, startedAt time.Time) {

	w.Header().Add("Content-Type", "text/html; charset=Shift_JIS")
	w.Header().Add("Date", startedAt.UTC().Format(http.TimeFormat))

	if err := createThreadErrorTmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleDat(sv *service.BoardService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		board := ps.ByName("board")
		threadKey := strings.Replace(ps.ByName("dat"), ".dat", "", 1)
		dat, err := sv.MakeDat(board, threadKey)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				log.Printf("ERROR: handleDat. %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		fmt.Fprintf(w, string(util.UTF8toSJIS(dat)))
	}
}

func handleSubjectTxt(sv *service.BoardService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		board := ps.ByName("board")
		subjectTxt, err := sv.MakeSubjectTxt(board)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				log.Printf("ERROR: handleSubjectTxt. %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		fmt.Fprintf(w, string(util.UTF8toSJIS(subjectTxt)))
	}
}
