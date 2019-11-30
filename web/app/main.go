package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	indexTmpl = template.Must(
		template.ParseFiles(filepath.Join("..", "template", "index.html")),
	)
)

func main() {
	router := httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/:board/bbs.cgi", handleBbsCgi)
	router.GET("/:board/subject.txt", handleSubjectTxt)
	router.GET("/:board/setting.txt", handleSettingTxt)
	router.GET("/:board/dat/:dat", handleDat)

	// Serve static files out of the public directory.
	// By configuring a static handler in app.yaml, App Engine serves all the
	// static content itself. As a result, the following two lines are in
	// effect for development only.
	public := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	http.Handle("/static/", public)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

// handleIndex uses a template to create an index.html.
func handleIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := indexTmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleBbsCgi(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	board := ps.ByName("board")
	if board != "test" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "bbs.cgi, %s!\n", board)
}

func handleSubjectTxt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "subject.txt, %s!\n", ps.ByName("board"))
}
func handleSettingTxt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "setting.txt, %s!\n", ps.ByName("board"))
}
func handleDat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "dat, %s %s!\n", ps.ByName("board"), ps.ByName("dat"))
}
