package handle

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/service"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var (
	indexTmpl             = template.Must(template.ParseFiles(filepath.Join("web", "template", "index.html")))
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
		handleTestDir(
			parseForm(
				injectService(sv)(
					handleAdminLogin()))))
	router.POST("/:board/_admin/",
		handleTestDir(
			parseForm(
				injectService(sv)(
					authenticate(
						handleAdmin())))))
	// 掲示板
	router.POST("/:board/bbs.cgi",
		protect(config.KEEP_OUT)(
			handleTestDir(
				parseForm(
					injectService(sv)(
						handleBbsCgi())))))
	router.GET("/:board/subject.txt",
		protect(config.KEEP_OUT)(
			injectService(sv)(
				handleSubjectTxt())))
	router.GET("/:board/dat/:dat",
		protect(config.KEEP_OUT)(
			injectService(sv)(
				handleDat())))

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

func parseForm(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r.ParseForm()
		h(w, r, ps)
	}
}
