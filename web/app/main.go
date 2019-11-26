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
	"sort"
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

func makeSubjectTxt(boardName string) (string, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return "", err
	}

	// Creates a Key instance.
	key := datastore.NameKey("Board", boardName, nil)

	// Gets a Board
	e := new(BoardEntity)
	if err := client.Get(ctx, key, e); err != nil {
		return "", err
	}

	// Sort
	sort.Sort(e.Subjects)

	buf := new(bytes.Buffer)
	for _, s := range e.Subjects {
		fmt.Fprintf(buf, "%s.dat<>%s \t (%d)", s.ThreadKey, s.ThreadTitle, s.MessageCount)
	}

	return buf.String(), nil
}

// Creates a Thread
func (sv *BoardService) createNewThread(boardName string,
	name string, mail string, now time.Time, id string, message string, title string) (err error) {

	// Gets a Board entity
	boardKey := datastore.NameKey("Board", boardName, nil)
	board := &BoardEntity{}
	if err = sv.GetBoard(boardKey, board); err != nil {
		return
	}

	// Adds to Subject
	threadKey := strconv.FormatInt(now.Unix(), 10)
	subject := Subject{
		ThreadKey:    threadKey,
		ThreadTitle:  title,
		MessageCount: 1,
		LastFloat:    now,
		LastModified: now,
	}
	board.Subjects = append(board.Subjects, subject)

	if err = sv.PutBoard(boardKey, board); err != nil {
		return
	}

	// Create dat
	datKey := datastore.NameKey("Dat", threadKey, boardKey)
	dat := createDat(name, mail, now, id, message, title)
	if err = sv.PutDat(datKey, dat); err != nil {
		return
	}
	return nil
}

// create dat. line: 1
func createDat(name string, mail string, date time.Time, id string, message string, title string) *DatEntity {
	dat := &DatEntity{}
	writeDat(dat, datFormat1, name, mail, date, id, message, title)
	return dat
}

// append dat. line: 2..
func appendDat(dat *DatEntity,
	name string, mail string, date time.Time, id string, message string) {

	writeDat(dat, datFormatN, name, mail, date, id, message, "")
}

func writeDat(dat *DatEntity, format string,
	name string, mail string, date time.Time, id string, message string, title string) {

	wr := bytes.NewBuffer(dat.Dat)
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

	dat.Dat = wr.Bytes()
}

func escapeDat(str string) string {
	return strings.ReplaceAll(str, "\n", "<br>")
}
