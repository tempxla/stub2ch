package main

import (
	"./service"
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	param_error_format = "bad parameter '%s' is: %v"
)

var (
	indexTmpl = template.Must(
		template.ParseFiles(filepath.Join("..", "template", "index.html")),
	)
)

// handleIndex uses a template to create an index.html.
func handleIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := indexTmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleBbsCgi(sv *service.BoardService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		board := ps.ByName("board")
		if board != "test" {
			http.NotFound(w, r)
			return
		}

		r.ParseForm()

		submit, err := processParam(require(r, "submit"), url.QueryUnescape)
		if err != nil {
			fmt.Fprint(w, param_error_format, "submit", err)
			return
		}

		switch submit {
		case "書き込む":
			// レスを書き込む
			handleWriteDat(sv, w, r)
		case "新規スレッド作成":
			// スレッドを立てる
		default:
			fmt.Fprint(w, "SJISで書いてね？")
		}
	}
}

func requireParam(r *http.Request, name string) (str string, err error) {
	if param, ok := r.PostForm[name]; !ok {
		err = fmt.Errorf("missing")
	} else if len(param) == 0 {
		err = fmt.Errorf("empty")
	} else if len(param) != 1 {
		err = fmt.Errorf("too many")
	} else {
		str = param[0]
	}
	return
}

func require(r *http.Request, name string) func() (string, error) {
	return func() (string, error) {
		return requireParam(r, name)
	}
}

func notEmpty(s string) (str string, err error) {
	if s == "" {
		err = fmt.Errorf("0 byte")
	} else {
		str = s
	}
	return
}

func notBlank(s string) (str string, err error) {
	if strings.TrimSpace(s) == "" {
		err = fmt.Errorf("blank")
	} else {
		str = s
	}
	return
}

func betweenStr(a, b string) func(string) (string, error) {
	return func(s string) (str string, err error) {
		if s < a || b < s {
			err = fmt.Errorf("%s < %s or %s < %s", s, a, b, s)
		} else {
			str = s
		}
		return
	}
}

func processParam(src func() (string, error),
	funcs ...func(string) (string, error)) (s string, e error) {

	s, e = src()
	if e != nil {
		return
	}

	for _, f := range funcs {
		s, e = f(s)
		if e != nil {
			return
		}
	}
	return
}

func computeId(ipAddr string) string {
	return ""
}

func handleWriteDat(sv *service.BoardService, w http.ResponseWriter, r *http.Request) {
	boardName, err := processParam(require(r, "bbs"), betweenStr("0", "zzzzzzzzzz"))
	if err != nil {
		fmt.Fprintf(w, param_error_format, "bbs", err)
		return
	}
	threadKey, err := processParam(require(r, "key"), betweenStr("0000000000", "9999999999"))
	if err != nil {
		fmt.Fprintf(w, param_error_format, "key", err)
		return
	}
	_, err = processParam(require(r, "time"), notEmpty)
	if err != nil {
		fmt.Fprintf(w, param_error_format, "time", err)
		return
	}
	name, err := processParam(require(r, "FROM"), url.QueryUnescape)
	if err != nil {
		fmt.Fprintf(w, param_error_format, "FROM", err)
		return
	}
	mail, err := processParam(require(r, "mail"), url.QueryUnescape)
	if err != nil {
		fmt.Fprintf(w, param_error_format, "mail", err)
		return
	}
	message, err := processParam(require(r, "MESSAGE"), url.QueryUnescape, notBlank)
	if err != nil {
		fmt.Fprintf(w, param_error_format, "MESSAGE", err)
		return
	}
	err = sv.WriteDat(boardName, threadKey, name, mail, computeId(r.RemoteAddr), message)
	if err != nil {
		// 存在しない
		// 1001
		// dat落ち
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
		fmt.Fprintf(w, dat)
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
		fmt.Fprintf(w, subjectTxt)
	}
}

func handleSettingTxt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "setting.txt, %s!\n", ps.ByName("board"))
}
