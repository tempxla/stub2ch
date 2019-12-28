package handle

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/internal/app/service"
	. "github.com/tempxla/stub2ch/internal/app/types"
	"github.com/tempxla/stub2ch/internal/app/util"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	param_error_format = "bad parameter '%s' is: %v"
)

func handleBbsCgi() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

		submit, err := process(requireOne(r, "submit"), sjisToUtf8String)
		if err != nil {
			http.Error(w, fmt.Sprintf(param_error_format, "submit", err), http.StatusBadRequest)
			return
		}

		switch submit {
		case "書き込む":
			// レスを書き込む
			handleWriteDat(w, r, sv)
		case "新規スレッド作成":
			// スレッドを立てる
			handleCreateThread(w, r, sv)
		case "上記全てを承諾して書き込む":
			if _, ok := r.PostForm["key"]; ok {
				handleWriteDat(w, r, sv)
			} else {
				handleCreateThread(w, r, sv)
			}
		default:
			fmt.Fprint(w, util.UTF8toSJISString("SJISで書いてね？"))
		}
	}
}

func handleWriteDat(w http.ResponseWriter, r *http.Request, sv *service.BoardService) {
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
	ipAddr := r.RemoteAddr
	if i := strings.LastIndexByte(ipAddr, ':'); i != -1 {
		ipAddr = ipAddr[:i]
	}
	id := sv.ComputeId(ipAddr, boardName)
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
	expires := startedAt.Add(time.Duration(7*24) * time.Hour).UTC().Format(http.TimeFormat)
	ipAddr := r.RemoteAddr
	if i := strings.LastIndexByte(ipAddr, ':'); i != -1 {
		ipAddr = ipAddr[:i]
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("PON=%s; expires=%s; path=/", ipAddr, expires))
	w.Header().Add("Set-Cookie", fmt.Sprintf("yuki=akari; expires=%s; path=/", expires))

	// Body
	view := map[string]string{
		"Title":     util.UTF8toSJISString(FromMaybe(title, "")),
		"Name":      util.UTF8toSJISString(name),
		"Mail":      util.UTF8toSJISString(mail),
		"Message":   util.UTF8toSJISString(message),
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

func handleCreateThread(w http.ResponseWriter, r *http.Request, sv *service.BoardService) {
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
	ipAddr := r.RemoteAddr
	if i := strings.LastIndexByte(ipAddr, ':'); i != -1 {
		ipAddr = ipAddr[:i]
	}
	id := sv.ComputeId(ipAddr, boardName)
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

func handleDat() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
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

func handleSubjectTxt() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
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
