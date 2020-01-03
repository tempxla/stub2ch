package handle

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/service"
	"github.com/tempxla/stub2ch/internal/app/util"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	user_agent = "Monazilla/1.00"
)

var (
	indexTmpl             = template.Must(template.ParseFiles(filepath.Join("web", "template", "index.html")))
	topTmpl               = template.Must(template.ParseFiles(filepath.Join("web", "template", "top.html")))
	datTmpl               = template.Must(template.ParseFiles(filepath.Join("web", "template", "dat.html")))
	writeDatConfirmTmpl   = template.Must(template.ParseFiles(filepath.Join("web", "template", "write_dat_confirm.html")))
	writeDatNotFoundTmpl  = template.Must(template.ParseFiles(filepath.Join("web", "template", "write_dat_not_found.html")))
	writeDatDoneTmpl      = template.Must(template.ParseFiles(filepath.Join("web", "template", "write_dat_done.html")))
	createThreadErrorTmpl = template.Must(template.ParseFiles(filepath.Join("web", "template", "create_thread_error.html")))
	adminIndexTmpl        = template.Must(template.ParseFiles(filepath.Join("web", "template", "admin", "index.html")))
)

type ServiceHandle func(http.ResponseWriter, *http.Request, httprouter.Params, *service.BoardService)

// HTTP routing
func NewBoardRouter(sv *service.BoardService) *httprouter.Router {
	router := httprouter.New()

	// トップ
	router.GET("/", handleIndex())

	// 管理ページ
	router.POST("/:board/_admin/login",
		protect(config.KEEP_OUT)(
			handleTestDir(
				handleParseForm(
					injectService(sv)(
						handleAdminLogin())))))
	router.POST("/:board/_admin/logout",
		protect(config.KEEP_OUT)(
			handleTestDir(
				handleParseForm(
					injectService(sv)(
						authenticate(
							handleAdminLogout()))))))
	router.POST("/:board/_admin/func/:fp1/:fp2",
		protect(config.KEEP_OUT)(
			handleTestDir(
				handleParseForm(
					injectService(sv)(
						authenticate(
							handleAdmin()))))))

	// 掲示板
	router.GET("/:board/",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				injectService(sv)(
					handleTop()))))
	router.GET("/:board/read.cgi/:boardName/:threadKey/",
		protect(config.KEEP_OUT)(
			handleTestDir(
				handleUserAgent(
					injectService(sv)(
						handleReadCgi())))))
	router.GET("/:board/SETTING.TXT",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				handleSettingTxt())))
	router.GET("/:board/head.txt",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				handleHeadTxt())))
	router.GET("/:board/subject.txt",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				injectService(sv)(
					handleSubjectTxt()))))
	router.GET("/:board/dat/:dat",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				injectService(sv)(
					handleDat()))))
	router.POST("/:board/bbs.cgi",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				handleTestDir(
					handleParseForm(
						injectService(sv)(
							handleBbsCgi()))))))

	// Jsonデモ
	router.POST("/:board/subject.json",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				handleParseForm(
					injectService(sv)(
						handlePrecure(
							handleSubjectJson()))))))
	router.POST("/:board/json/:dat",
		protect(config.KEEP_OUT)(
			handleUserAgent(
				handleParseForm(
					injectService(sv)(
						handlePrecure(
							handleDatJson()))))))

	// 静的ファイル
	// GAEの設定はapp.yamlなので、これは開発用
	// The path must end with "/*filepath"
	router.ServeFiles("/:board/_static/*filepath", http.Dir("web/static"))

	return router
}

func injectService(sv *service.BoardService) func(ServiceHandle) httprouter.Handle {
	return func(sh ServiceHandle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

			var boardService *service.BoardService
			var err error

			if sv != nil {
				boardService = sv
			} else {
				boardService, err = service.DefaultBoardService()
				if err != nil {
					http.Error(w, fmt.Sprintf("%v", err), http.StatusServiceUnavailable) // 503
					return
				}
			}

			// Injection
			sh(w, r, ps, boardService)
		}
	}
}

// トップページ
func handleIndex() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := indexTmpl.Execute(w, nil); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func protect(keepOut bool) func(httprouter.Handle) httprouter.Handle {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			if keepOut {
				http.Error(w, "KEEP OUT", http.StatusServiceUnavailable) // 503
				return
			}
			h(w, r, ps)
		}
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

func handleParseForm(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r.ParseForm()
		h(w, r, ps)
	}
}

func handleUserAgent(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// UAはあればよい感じ
		if len(r.UserAgent()) < len(user_agent) {
			setContentTypePlainSjis(w)
			http.Error(w, util.UTF8toSJISString("m9(^Д^)"), http.StatusBadRequest)
			return
		}

		h(w, r, ps)
	}
}

func setContentTypeHtmlSjis(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=Shift_JIS")
}

func setContentTypePlainSjis(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=Shift_JIS")
}

func getIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		comma := strings.Index(ipProxy, ",")
		if comma != -1 {
			// X-Forwarded-For: client1, proxy1, proxy2
			return ipProxy[:comma]
		}
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
