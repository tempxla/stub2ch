package handle

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/configs/app/bbscfg"
	"github.com/tempxla/stub2ch/internal/app/service"
	"github.com/tempxla/stub2ch/internal/app/types/errors"
	"github.com/tempxla/stub2ch/internal/app/util"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	param_error_format = "bad parameter '%s' is: %v"
	top_load_delay     = 10 // seconds
	top_subject_limit  = 10
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
	setting := bbscfg.GetSetting(boardName)
	if setting == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	name, ok := requireName(w, r, setting)
	if !ok {
		return
	}
	mail, ok := requireMail(w, r, setting)
	if !ok {
		return
	}
	message, ok := requireMessage(w, r, setting)
	if !ok {
		return
	}
	_, ok = requireReferer(w, r, boardName)
	if !ok {
		return
	}
	ipAddr := getIP(r)

	// クッキー確認
	if executeWriteDatConfirmTmpl(w, r, ipAddr,
		boardName, name, mail, message, sv.StartedAt(), "", threadKey) {
		return
	}
	// 書き込み
	id := sv.ComputeId(ipAddr, boardName)
	resnum, err := sv.WriteDat(setting, boardName, threadKey, name, mail, id, message)
	if err != nil {
		// 存在しない or dat落ち or 1001 or 容量オーバー
		executeWriteDatNotFoundTmpl(w, r, boardName, threadKey, sv.StartedAt())
		return
	}
	// 書き込み完了
	logPrintWriteDone(boardName, threadKey, resnum, id, ipAddr)

	executeWriteDoneTmpl(w, r, boardName, threadKey, id, resnum, sv.StartedAt())
}

func executeWriteDoneTmpl(w http.ResponseWriter, r *http.Request,
	boardName, threadKey, id string, resnum int, startedAt time.Time) {

	setContentTypeHtmlSjis(w)

	w.Header().Add("x-Resnum", strconv.Itoa(resnum))
	//                         01234567890123456789
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

	setContentTypeHtmlSjis(w)

	// //hebi.5ch.net/test/read.cgi/news4vip/1575543566/
	view := fmt.Sprintf("//%s/test/read.cgi/%s/%s/", r.Host, boardName, threadKey)
	if err := writeDatNotFoundTmpl.Execute(w, view); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Returns false if Cookie Found.
func executeWriteDatConfirmTmpl(w http.ResponseWriter, r *http.Request, ipAddr string,
	boardName, name, mail, message string, startedAt time.Time, title, threadKey string) bool {

	if c, err := r.Cookie("PON"); err == nil && c.Value == ipAddr {
		if c, err := r.Cookie("yuki"); err == nil && c.Value == "akari" {
			// Cookie Found. Need not to forward Confirm page.
			return false
		}
	}

	setContentTypeHtmlSjis(w)

	// Domain属性を指定しないCookieは、Cookieを発行したホストのみに送信される
	expires := startedAt.Add(time.Duration(7*24) * time.Hour).UTC().Format(http.TimeFormat)
	w.Header().Add("Set-Cookie", fmt.Sprintf("PON=%s; expires=%s; path=/", ipAddr, expires))
	w.Header().Add("Set-Cookie", fmt.Sprintf("yuki=akari; expires=%s; path=/", expires))

	// Body
	view := map[string]string{
		"Title":     util.UTF8toSJISString(title),
		"Name":      util.UTF8toSJISString(name),
		"Mail":      util.UTF8toSJISString(mail),
		"Message":   util.UTF8toSJISString(message),
		"BoardName": boardName,
		"Time":      strconv.FormatInt(startedAt.Unix(), 10),
	}
	if threadKey != "" {
		view["ThreadKey"] = threadKey
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
	setting := bbscfg.GetSetting(boardName)
	if setting == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	title, ok := requireTitle(w, r, setting)
	if !ok {
		return
	}
	_, ok = requireTime(w, r)
	if !ok {
		return
	}
	name, ok := requireName(w, r, setting)
	if !ok {
		return
	}
	mail, ok := requireMail(w, r, setting)
	if !ok {
		return
	}
	message, ok := requireMessage(w, r, setting)
	if !ok {
		return
	}
	_, ok = requireReferer(w, r, boardName)
	if !ok {
		return
	}
	ipAddr := getIP(r)

	// クッキー確認
	if executeWriteDatConfirmTmpl(w, r, ipAddr,
		boardName, name, mail, message, sv.StartedAt(), title, "") {
		return
	}
	// スレ立て
	id := sv.ComputeId(ipAddr, boardName)
	threadKey, err := sv.CreateThread(setting, boardName, name, mail, id, message, title)
	if err != nil {
		// スレ立て失敗
		executeCreateThreadErrorTmpl(w, r, sv.StartedAt())
		return
	}
	// 書き込み完了
	logPrintWriteDone(boardName, threadKey, 1, id, ipAddr)

	executeWriteDoneTmpl(w, r, boardName, threadKey, id, 1, sv.StartedAt())
}

func logPrintWriteDone(boardName, threadKey string, resnum int, id, ipAddr string) {
	log.Printf("[WRITE DONE] /%s/%s/%d id:%s ip:%s ", boardName, threadKey, resnum, id, ipAddr)
}

func executeCreateThreadErrorTmpl(w http.ResponseWriter, r *http.Request, startedAt time.Time) {

	setContentTypeHtmlSjis(w)

	if err := createThreadErrorTmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleDat() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
		board := ps.ByName("board")
		threadKey := strings.Replace(ps.ByName("dat"), ".dat", "", 1)
		dat, lastModifiedTime, err := sv.MakeDat(board, threadKey)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				log.Printf("ERROR: handleDat. %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		sjisDat := util.UTF8toSJIS(dat)
		lastModified := lastModifiedTime.UTC().Format(http.TimeFormat)

		// 差分取得判定
		ifModifiedSince := r.Header.Get("If-Modified-Since")
		// 差分取得でない
		if ifModifiedSince == "" {
			setContentTypePlainSjis(w)
			w.Header().Add("Last-Modified", lastModified)
			fmt.Fprintf(w, string(sjisDat))
			return
		}
		// 更新されていない
		if ifModifiedSince == lastModified {
			w.WriteHeader(http.StatusNotModified) // 304
			return
		}
		// 差分取得
		rangeBytes, err := parseDatRange(r.Header.Get("Range"))
		if err != nil {
			// リクエストがおかしい
			http.Error(w, "Need Range ?", http.StatusBadRequest) // 400
		} else if rangeBytes > len(sjisDat) {
			// あぼーん有り
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable) // 416
		} else {
			// 差分DAT
			setContentTypePlainSjis(w)
			w.Header().Add("Last-Modified", lastModified)
			w.WriteHeader(http.StatusPartialContent) // 206
			fmt.Fprintf(w, string(sjisDat[rangeBytes:]))
		}
	}
}

func parseDatRange(rangeHeader string) (int, error) {
	//形式: bytes=3050-
	start := 6                  // bytes=^3050-
	end := len(rangeHeader) - 1 // bytes=3050^-
	if strings.Index(rangeHeader, "bytes=") != 0 {
		return -1, fmt.Errorf("parse error: %v", rangeHeader)
	}
	if strings.LastIndex(rangeHeader, "-") != end {
		return -1, fmt.Errorf("parse error: %v", rangeHeader)
	}
	return strconv.Atoi(rangeHeader[start:end])
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

		setContentTypePlainSjis(w)
		fmt.Fprintf(w, string(util.UTF8toSJIS(subjectTxt)))
	}
}

func handleSettingTxt() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		board := ps.ByName("board")
		stng := bbscfg.GetSetting(board)
		if stng == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		settingTxt := bbscfg.MakeSettingTxt(stng)
		setContentTypePlainSjis(w)
		fmt.Fprintf(w, string(util.UTF8toSJIS(settingTxt)))
	}
}

func handleHeadTxt() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		board := ps.ByName("board")
		headTxt := bbscfg.MakeHeadTxt(board)
		if headTxt == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		setContentTypePlainSjis(w)
		fmt.Fprintf(w, string(util.UTF8toSJIS(headTxt)))
	}
}

func handleBoard() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
		board := ps.ByName("board")
		stng := bbscfg.GetSetting(board)
		if stng == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		view := &struct {
			BoardName  string
			BoardTitle string
			Precure    int64
		}{
			board,
			util.UTF8toSJISString(stng.BBS_TITLE()),
			sv.StartedAt().Unix() - top_load_delay, // 初回ロードで待たせないため時間をずらしておく
		}

		setContentTypeHtmlSjis(w)

		if err := boardTmpl.Execute(w, view); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func handleSubjectJson() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

		board := ps.ByName("board")
		stng := bbscfg.GetSetting(board)
		if stng == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		subjectJson, err := sv.MakeSubjectJson(board, top_subject_limit)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				log.Printf("ERROR: handleSubjectJson. %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		fmt.Fprintf(w, string(subjectJson))
	}
}

func handleReadCgi() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
		boardName := ps.ByName("boardName")
		stng := bbscfg.GetSetting(boardName)
		if stng == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		threadKey := ps.ByName("threadKey")

		view := &struct {
			BoardName  string
			ThreadKey  string
			BoardTitle string
			Precure    int64
		}{
			boardName,
			threadKey,
			util.UTF8toSJISString(stng.BBS_TITLE()),
			sv.StartedAt().Unix() - top_load_delay, // 初回ロードで待たせないため時間をずらしておく
		}

		setContentTypeHtmlSjis(w)

		if err := datTmpl.Execute(w, view); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func handleDatJson() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

		board := ps.ByName("board")
		threadKey := strings.Replace(ps.ByName("dat"), ".json", "", 1)

		// 差分取得判定
		ifModifiedSince := r.Header.Get("If-Modified-Since")
		min := 1
		max := 11 // 暫定 10 までとする メッセージのため11個返す
		if ifModifiedSince != "" {
			// 差分取得ですよ
			rangeInt, err := strconv.Atoi(r.Header.Get("Range"))
			if err != nil {
				// リクエストがおかしい
				http.Error(w, "Need Range ?", http.StatusBadRequest) // 400
				return
			}
			min = rangeInt + 1
		}
		datJson, err := sv.MakeDatJson(board, threadKey, ifModifiedSince, min, max)
		if err != nil {
			if err == errors.NOT_MODIFIED {
				// 更新されていない
				w.WriteHeader(http.StatusNotModified) // 304
			} else if err == datastore.ErrNoSuchEntity {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				log.Printf("ERROR: handleDatJson. %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// 送信
		if ifModifiedSince != "" {
			// 差分取得ですよ
			w.WriteHeader(http.StatusPartialContent) // 206
		}
		fmt.Fprintf(w, string(datJson))
	}
}

func handlePrecure(sh ServiceHandle) ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

		precure, err := requireOne(r, "precure")()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest) // 400
			return
		}
		t, err := strconv.ParseInt(precure, 10, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest) // 400
			return
		}

		if sv.StartedAt().Unix() < t+top_load_delay {
			http.Error(w, "(；´Д｀)焦らないで...", http.StatusServiceUnavailable) // 503
			return
		}

		sh(w, r, ps, sv)
	}
}
