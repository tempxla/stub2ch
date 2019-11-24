package main

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html"
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
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	datFormat1 = "%s<>%s<>%s(%s) %s ID:%s<> %s <>%s"
	datFormatN = "\n" + datFormat1
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

// create dat. line: 1
func createDat(name string, mail string, date time.Time, id string, message string, title string) []byte {
	return writeDat([]byte{}, datFormat1, name, mail, date, id, message, title)
}

// append dat. line: 2..
func appendDat(dat []byte,
	name string, mail string, date time.Time, id string, message string) []byte {
	return writeDat(dat, datFormatN, name, mail, date, id, message, "")
}

func writeDat(dat []byte, format string,
	name string, mail string, date time.Time, id string, message string, title string) []byte {

	wr := bytes.NewBuffer(dat)
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	// 2行目以降はスレタイは無し
	fmt.Fprintf(wr, format,
		html.EscapeString(name),               // 名前
		html.EscapeString(mail),               // メール
		date.Format(datDateLayout),            // 年月日
		weekdaysJp[date.Weekday()],            // 曜
		date.Format(datTimeLayout),            // 時分秒
		id,                                    // ID
		escapeDat(html.EscapeString(message)), // 本文
		html.EscapeString(title))              // スレタイ

	return wr.Bytes()
}

func escapeDat(str string) string {
	return strings.ReplaceAll(str, "\n", "<br>")
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
