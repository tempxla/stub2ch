package main

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	projectID     = "stub2ch"
	datDateLayout = "2006/01/02"
	datTimeLayout = "15:04:05.000"
)

var (
	indexTmpl = template.Must(
		template.ParseFiles(filepath.Join("..", "template", "index.html")),
	)
	weekdaysJp = [...]string{"日", "月", "火", "水", "木", "金", "土"}
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

func createNewThread(boardName string, title string, name string, mail string, message string) {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get Board
	key := datastore.NameKey("Board", boardName, nil)
	board := new(BoardEntity)
	if err := client.Get(ctx, key, board); err != nil {
		log.Fatalf("Failed to get BoardEntity: %v", err)
	}

	// Add to subject
	now := time.Now()
	threadKey := strconv.FormatInt(now.Unix(), 10)
	subject := Subject{
		ThreadKey:    threadKey,
		ThreadTitle:  title,
		MessageCount: 1,
		LastFloat:    now, // if sage case ?
		LastModified: now,
	}
	board.Subjects = append(board.Subjects, subject)

	if _, err := client.Put(ctx, key, board); err != nil {
		log.Fatalf("Failed to save board: %v", err)
	}

	// // Create dat
	// ancestor := datastore.NameKey("Board", boardName, nil)
	// key = datastore.NameKey("Dat", threadKey, ancestor)
	// var dat []byte
	// if _, err := client.Put(ctx, key, dat); err != nil {
	// 	log.Fatalf("Failed to save dat: %v", err)
	// }
}

func createDat(
	name string, mail string, date time.Time,
	id string, message string, title string) []byte {

	d := date.Format(datDateLayout)
	t := date.Format(datTimeLayout)
	w := weekdaysJp[date.Weekday()]
	wr := bytes.NewBuffer([]byte{})
	fmt.Fprintf(wr, "%s<>%s<>%s(%s) %s ID:%s<> %s <>%s",
		name, mail, d, w, t, id, message, title)
	return wr.Bytes()
}

func appendDat(dat []byte,
	name string, mail string, date time.Time,
	id string, message string) {

	d := date.Format(datDateLayout)
	t := date.Format(datTimeLayout)
	w := weekdaysJp[date.Weekday()]
	wr := bytes.NewBuffer(dat)
	fmt.Fprintf(wr, "\n%s<>%s<>%s(%s) %s ID:%s<> %s <>",
		name, mail, d, w, t, id, message)
}

func escapeDat(str string) string {
	return strings.Replace(str, "\n", "<br>", -1) // ReplaceAll
}

// Kind=Board
// Key=BoardName
type BoardEntity struct {
	Subjects []Subject `datastore:",noindex"`
}

type Subject struct {
	ThreadKey    string
	ThreadTitle  string
	MessageCount int
	LastFloat    time.Time
	LastModified time.Time
}

// Kind=Dat
// Ancestor=Board
// Key=ThreadKey
type DatEntity []byte
