package service

import (
	"../config"
	E "../entity"
	"bytes"
	"cloud.google.com/go/datastore"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"
)

const (
	datFormatN = "\n" + config.DAT_FORMAT
)

// Dependency injection for Board
type BoardService struct {
	repo BoardRepository
}

type BoardRepository interface {
	GetBoard(key *datastore.Key, entity *E.BoardEntity) (err error)
	PutBoard(key *datastore.Key, entity *E.BoardEntity) (err error)
	GetDat(key *datastore.Key, entity *E.DatEntity) (err error)
	PutDat(key *datastore.Key, entity *E.DatEntity) (err error)
}

func NewBoardService(repo BoardRepository) *BoardService {
	return &BoardService{
		repo: repo,
	}
}

// データストアからエンティティを取得しdatを返す
func (sv *BoardService) MakeDat(boardName string, threadKey string) (_ string, err error) {
	// Creates a Key instance.
	key := datastore.NameKey("Dat", threadKey,
		datastore.NameKey("Board", boardName, nil))

	// Gets a Board
	e := new(E.DatEntity)
	if err = sv.repo.GetDat(key, e); err != nil {
		return
	}

	return string(e.Dat), nil
}

// データストアからエンティティを取得しsubject.txtとして返す
func (sv *BoardService) MakeSubjectTxt(boardName string) (_ string, err error) {
	// Creates a Key instance.
	key := datastore.NameKey("Board", boardName, nil)

	// Gets a Board
	e := new(E.BoardEntity)
	if err = sv.repo.GetBoard(key, e); err != nil {
		return
	}

	buf := new(bytes.Buffer)
	for i, s := range e.Subjects {
		if i > 0 {
			fmt.Fprintf(buf, "\n")
		}
		fmt.Fprintf(buf, "%s.dat<>%s \t (%d)", s.ThreadKey, s.ThreadTitle, s.MessageCount)
	}

	return buf.String(), nil
}

// Creates a Thread
func (sv *BoardService) CreateNewThread(boardName string,
	name string, mail string, now time.Time, id string, message string, title string) (err error) {

	// Gets a Board entity
	boardKey := datastore.NameKey("Board", boardName, nil)
	board := &E.BoardEntity{}
	if err = sv.repo.GetBoard(boardKey, board); err != nil {
		return
	}

	// Adds to Subject
	threadKey := strconv.FormatInt(now.Unix(), 10)
	subject := E.Subject{
		ThreadKey:    threadKey,
		ThreadTitle:  title,
		MessageCount: 1,
		LastModified: now,
	}
	board.Subjects = append(board.Subjects, subject)

	if err = sv.repo.PutBoard(boardKey, board); err != nil {
		return
	}

	// Create dat
	datKey := datastore.NameKey("Dat", threadKey, boardKey)
	dat := createDat(name, mail, now, id, message, title)
	if err = sv.repo.PutDat(datKey, dat); err != nil {
		return
	}
	return nil
}

func (sv *BoardService) WriteDat(boardName string, threadKey string,
	name string, mail string, now time.Time, id string, message string) (err error) {

	// Creates a Key instance.
	boardKey := datastore.NameKey("Board", boardName, nil)
	datKey := datastore.NameKey("Dat", threadKey, boardKey)

	// Get Entities
	dat := new(E.DatEntity)
	if err = sv.repo.GetDat(datKey, dat); err != nil {
		return
	}
	board := new(E.BoardEntity)
	if err = sv.repo.GetBoard(boardKey, board); err != nil {
		return
	}

	// 書き込み
	appendDat(dat, name, mail, now, id, message)

	// subject.txtの更新
	err = updateSubjectsWhenWriteDat(board, threadKey, mail, now)
	if err != nil {
		return
	}

	// Push Entities
	if err = sv.repo.PutDat(datKey, dat); err != nil {
		return
	}
	if err = sv.repo.PutBoard(boardKey, board); err != nil {
		return
	}
	return nil
}

func updateSubjectsWhenWriteDat(board *E.BoardEntity,
	threadKey string, mail string, now time.Time) error {

	sbjLen := len(board.Subjects)
	pos := -1
	for i, sbj := range board.Subjects {
		if sbj.ThreadKey == threadKey {
			pos = i
			break
		}
	}
	if pos == -1 {
		return fmt.Errorf("fail update subjects. len:%v key:%v", sbjLen, threadKey)
	}
	board.Subjects[pos].MessageCount++
	board.Subjects[pos].LastModified = now

	// (´∀`∩)↑age↑
	if sbjLen > 1 && mail != "sage" {
		sbj := board.Subjects[pos]
		// 切り出し
		board.Subjects = append(board.Subjects[:pos], board.Subjects[pos+1:]...)
		// 先頭に追加
		board.Subjects, board.Subjects[0] =
			append(board.Subjects[0:1], board.Subjects[0:]...), sbj
	}
	return nil
}

// create dat. line: 1
func createDat(name string, mail string, date time.Time, id string, message string, title string) *E.DatEntity {
	dat := &E.DatEntity{}
	writeDat(dat, config.DAT_FORMAT, name, mail, date, id, message, title)
	return dat
}

// append dat. line: 2..
func appendDat(dat *E.DatEntity,
	name string, mail string, date time.Time, id string, message string) {

	writeDat(dat, datFormatN, name, mail, date, id, message, "")
}

func writeDat(dat *E.DatEntity, format string,
	name string, mail string, date time.Time, id string, message string, title string) {

	wr := bytes.NewBuffer(dat.Dat)
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	// 2行目以降はスレタイは無し
	fmt.Fprintf(wr, format,
		html.EscapeString(name),               // 名前
		html.EscapeString(mail),               // メール
		date.Format(config.DAT_DATE_LAYOUT),   // 年月日
		config.WEEK_DAYS_JP[date.Weekday()],   // 曜
		date.Format(config.DAT_TIME_LAYOUT),   // 時分秒
		id,                                    // ID
		escapeDat(html.EscapeString(message)), // 本文
		html.EscapeString(title))              // スレタイ

	dat.Dat = wr.Bytes()
}

func escapeDat(str string) string {
	return strings.ReplaceAll(str, "\n", "<br>")
}
