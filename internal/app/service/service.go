package service

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/configs/app/secretcfg"
	. "github.com/tempxla/stub2ch/internal/app/types"
	"html"
	"strconv"
	"strings"
	"time"
)

const (
	dat_date_layout = "2006/01/02"
	dat_time_layout = "15:04:05.000"
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	dat_format = "%s<>%s<>%s(%s) %s ID:%s<> %s <>%s\n"
)

var (
	week_days_jp = [...]string{"日", "月", "火", "水", "木", "金", "土"}
)

// Dependency injection for Board
type BoardService struct {
	repo  BoardRepository
	env   BoardEnvironment
	Admin *AdminFunction
}

func DefaultBoardService() (*BoardService, error) {

	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	repo := &BoardStore{
		Context: ctx,
		Client:  client,
	}
	sysEnv := &SysEnv{
		StartedTime:   time.Now(),
		ComputeIdSalt: secretcfg.COMPUTE_ID_SALT,
	}
	mem := &AlterMemcache{
		Context: ctx,
		Client:  client,
	}
	adminRepo := &AdminBoardStore{
		repo: repo,
	}

	return NewBoardService(RepoConf(repo), EnvConf(sysEnv), AdminConf(adminRepo, mem)), nil
}

func NewBoardService(config ...func(*BoardService) *BoardService) *BoardService {
	sv := &BoardService{
		Admin: &AdminFunction{},
	}
	for _, conf := range config {
		sv = conf(sv)
	}
	return sv
}

func RepoConf(repo BoardRepository) func(*BoardService) *BoardService {
	return func(sv *BoardService) *BoardService {
		sv.repo = repo
		return sv
	}
}

func EnvConf(env BoardEnvironment) func(*BoardService) *BoardService {
	return func(sv *BoardService) *BoardService {
		sv.env = env
		return sv
	}
}

func AdminConf(repo AdminBoardRepository, mem BoardMemcache) func(*BoardService) *BoardService {
	return func(sv *BoardService) *BoardService {
		sv.Admin.repo = repo
		sv.Admin.mem = mem
		return sv
	}
}

// データストアからエンティティを取得しdatを返す
func (sv *BoardService) MakeDat(boardName string, threadKey string) (_ []byte, err error) {
	// Creates a Key instance.
	key := sv.repo.DatKey(threadKey, sv.repo.BoardKey(boardName))

	// Gets a Board
	e := new(DatEntity)
	if err = sv.repo.GetDat(key, e); err != nil {
		return
	}

	return e.Dat, nil
}

// データストアからエンティティを取得しsubject.txtとして返す
func (sv *BoardService) MakeSubjectTxt(boardName string) (_ []byte, err error) {
	// Creates a Key instance.
	key := sv.repo.BoardKey(boardName)

	// Gets a Board
	e := new(BoardEntity)
	if err = sv.repo.GetBoard(key, e); err != nil {
		return
	}

	buf := new(bytes.Buffer)
	for _, s := range e.Subjects {
		fmt.Fprintf(buf, "%s.dat<>%s \t (%d)\n", s.ThreadKey, s.ThreadTitle, s.MessageCount)
	}

	return buf.Bytes(), nil
}

// Creates a Thread
func (sv *BoardService) CreateThread(boardName string,
	name string, mail string, now time.Time, id string, message string,
	title string) (threadKey string, err error) {

	// New subject
	threadKey = strconv.FormatInt(now.Unix(), 10)
	subject := Subject{
		ThreadKey:    threadKey,
		ThreadTitle:  title,
		MessageCount: 1,
		LastModified: now,
	}

	// Create dat
	dat := createDat(name, mail, now, id, message, title)

	// Key
	boardKey := sv.repo.BoardKey(boardName)
	datKey := sv.repo.DatKey(threadKey, boardKey)

	// Start transaction
	err = sv.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		// Get And Check
		board := &BoardEntity{}
		if err := sv.repo.TxGetBoard(tx, boardKey, board); err != nil {
			return err
		}
		for _, sbj := range board.Subjects {
			if sbj.ThreadKey == datKey.Key.Name {
				return fmt.Errorf("thread key is duplicate")
			}
		}
		// Add
		board.Subjects = append(board.Subjects, subject)
		// Save
		if err := sv.repo.TxPutBoard(tx, boardKey, board); err != nil {
			return err
		}
		if err := sv.repo.TxPutDat(tx, datKey, dat); err != nil {
			return err
		}
		return nil
	})
	return
}

func (sv *BoardService) WriteDat(boardName, threadKey,
	name, mail, id, message string) (resnum int, err error) {

	// Creates a Key instance.
	boardKey := sv.repo.BoardKey(boardName)
	datKey := sv.repo.DatKey(threadKey, boardKey)

	err = sv.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		// Get Entities
		dat := new(DatEntity)
		if err := sv.repo.TxGetDat(tx, datKey, dat); err != nil {
			return err
		}
		board := new(BoardEntity)
		if err := sv.repo.TxGetBoard(tx, boardKey, board); err != nil {
			return err
		}

		// 書き込み
		appendDat(dat, name, mail, sv.env.StartedAt(), id, message)

		// subject.txtの更新
		resnum, err = updateSubjectsWhenWriteDat(board, threadKey, mail, sv.env.StartedAt())
		if err != nil {
			return err
		}

		// Push Entities
		if err := sv.repo.TxPutDat(tx, datKey, dat); err != nil {
			return err
		}
		if err = sv.repo.TxPutBoard(tx, boardKey, board); err != nil {
			return err
		}
		return nil
	})
	return
}

func updateSubjectsWhenWriteDat(board *BoardEntity,
	threadKey string, mail string, now time.Time) (resnum int, err error) {

	sbjLen := len(board.Subjects)
	pos := -1
	for i, sbj := range board.Subjects {
		if sbj.ThreadKey == threadKey {
			pos = i
			break
		}
	}
	if pos == -1 {
		err = fmt.Errorf("fail update subjects. len:%v key:%v", sbjLen, threadKey)
		return
	}
	resnum = board.Subjects[pos].MessageCount + 1
	board.Subjects[pos].MessageCount = resnum
	board.Subjects[pos].LastModified = now

	// (´∀`∩)↑age↑
	if sbjLen > 1 && mail != "sage" {
		// 切り出し
		sbj := board.Subjects[pos]
		board.Subjects = append(board.Subjects[:pos], board.Subjects[pos+1:]...)
		// 先頭に追加
		board.Subjects, board.Subjects[0] =
			append(board.Subjects[0:1], board.Subjects[0:]...), sbj
	}
	return
}

// create dat. line: 1
func createDat(name string, mail string, date time.Time, id string, message string, title string) *DatEntity {
	dat := &DatEntity{}
	writeDat(dat, dat_format, name, mail, date, id, message, title)
	return dat
}

// append dat. line: 2..
func appendDat(dat *DatEntity,
	name string, mail string, date time.Time, id string, message string) {

	writeDat(dat, dat_format, name, mail, date, id, message, "")
}

func writeDat(dat *DatEntity, format string,
	name string, mail string, date time.Time, id string, message string, title string) {

	wr := bytes.NewBuffer(dat.Dat)
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	// 2行目以降はスレタイは無し
	fmt.Fprintf(wr, format,
		html.EscapeString(name),               // 名前
		html.EscapeString(mail),               // メール
		date.Format(dat_date_layout),          // 年月日
		week_days_jp[date.Weekday()],          // 曜
		date.Format(dat_time_layout),          // 時分秒
		id,                                    // ID
		escapeDat(html.EscapeString(message)), // 本文
		html.EscapeString(title))              // スレタイ

	dat.Dat = wr.Bytes()
}

func escapeDat(str string) string {
	return strings.ReplaceAll(str, "\n", "<br>")
}

func (sv *BoardService) ComputeId(ipAddr, boardName string) string {
	// http://age.s22.xrea.com/talk2ch/id.txt
	ipmd5 := fmt.Sprintf("%x", md5.Sum([]byte(ipAddr)))
	ipmd5 = ipmd5[len(ipmd5)-4:]
	ipmd5 += boardName
	// ipmd5 += strconv.Itoa(sv.env.Now().Day())
	ipmd5 += sv.env.StartedAt().Format("2006/01/02")
	ipmd5 += sv.env.SaltComputeId()

	// full := md5.Sum([]byte(ipmd5))
	full := sha256.Sum256([]byte(ipmd5))
	return string(base64.StdEncoding.EncodeToString(full[:])[0:8])
}

func (sv *BoardService) StartedAt() time.Time {
	return sv.env.StartedAt()
}
