package service

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/configs/app/secretcfg"
	"github.com/tempxla/stub2ch/configs/app/setting"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
	jboard "github.com/tempxla/stub2ch/internal/app/types/json/board"
	"github.com/tempxla/stub2ch/internal/app/util"
	"html"
	"strconv"
	"strings"
	"time"
)

const (
	dat_date_layout = "2006/01/02"
	dat_time_layout = "15:04:05.000"
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	dat_format      = "%s<>%s<>%s(%s) %s ID:%s<> %s <>%s\n"
	dat_format_1001 = "%d<><>Over %d Thread<> このスレッドは%dを超えました。 <br> 新しいスレッドを立ててください。 <>\n"
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

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	repo := &BoardStore{
		Context: ctx,
		Client:  client,
	}
	sysEnv := &SysEnv{
		StartedTime:   time.Now().In(jst),
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
func (sv *BoardService) MakeDat(boardName, threadKey string) (_ []byte, _ time.Time, err error) {
	// Creates a Key instance.
	key := sv.repo.DatKey(threadKey, sv.repo.BoardKey(boardName))

	// Gets a Board
	dat := new(dat.Entity)
	if err = sv.repo.GetDat(key, dat); err != nil {
		return
	}

	return dat.Bytes, dat.LastModified, nil
}

// データストアからエンティティを取得しsubject.txtとして返す
func (sv *BoardService) MakeSubjectTxt(boardName string) (_ []byte, err error) {
	// Creates a Key instance.
	key := sv.repo.BoardKey(boardName)

	// Gets a Board
	e := new(board.Entity)
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
func (sv *BoardService) CreateThread(stng setting.BBS, boardName string,
	name string, mail string, now time.Time, id string, message string,
	title string) (threadKey string, err error) {

	// New subject
	threadKey = strconv.FormatInt(now.Unix(), 10)
	subject := board.Subject{
		ThreadKey:    threadKey,
		ThreadTitle:  escapeDat(html.EscapeString(title)),
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
		boardEntity := &board.Entity{}
		if err := sv.repo.TxGetBoard(tx, boardKey, boardEntity); err != nil {
			return err
		}
		for _, sbj := range boardEntity.Subjects {
			if sbj.ThreadKey == datKey.DSKey.Name {
				return fmt.Errorf("thread key is duplicate")
			}
		}

		// 制限チェキ
		if n := len(boardEntity.Subjects); n >= stng.STUB_THREAD_COUNT() {
			return fmt.Errorf("%d: これ以上スレ立てできません。。。", n)
		}
		if n := boardEntity.WriteCount; n >= stng.STUB_WRITE_ENTITY_LIMIT() {
			return fmt.Errorf("%d: 今日はこれ以上スレ立てできません。。。", n)
		}

		// 先頭に追加
		boardEntity.Subjects = append([]board.Subject{subject}, boardEntity.Subjects...)
		boardEntity.WriteCount++

		// Save
		if err := sv.repo.TxPutBoard(tx, boardKey, boardEntity); err != nil {
			return err
		}
		if err := sv.repo.TxPutDat(tx, datKey, dat); err != nil {
			return err
		}
		return nil
	})
	return
}

func (sv *BoardService) WriteDat(stng setting.BBS, boardName, threadKey,
	name, mail, id, message string) (resnum int, err error) {

	// Creates a Key instance.
	boardKey := sv.repo.BoardKey(boardName)
	datKey := sv.repo.DatKey(threadKey, boardKey)

	err = sv.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		// Get Entities
		dat := new(dat.Entity)
		if err := sv.repo.TxGetDat(tx, datKey, dat); err != nil {
			return err
		}
		board := new(board.Entity)
		if err := sv.repo.TxGetBoard(tx, boardKey, board); err != nil {
			return err
		}

		// 容量オーバー
		if len(util.UTF8toSJIS(dat.Bytes)) >= stng.STUB_DAT_CAPACITY() {
			return fmt.Errorf("容量超過: これ以上書き込めません。。。")
		}

		// 書き込み
		appendDat(dat, name, mail, sv.env.StartedAt(), id, message)

		// subject.txtの更新
		resnum, err = updateSubjectsWhenWriteDat(stng, board, threadKey, mail, sv.env.StartedAt())
		if err != nil {
			return err
		}

		// 1001カキコ
		if n := stng.STUB_MESSAGE_COUNT(); resnum == n {
			dat.Bytes = append(dat.Bytes, []byte(fmt.Sprintf(dat_format_1001, n+1, n, n))...)
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

func updateSubjectsWhenWriteDat(stng setting.BBS, board *board.Entity,
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

	// データストア制限チェック
	if board.WriteCount >= stng.STUB_WRITE_ENTITY_LIMIT() {
		err = fmt.Errorf("%d: 今日はこれ以上書き込めません。。。", board.WriteCount)
		return
	}

	resnum = board.Subjects[pos].MessageCount + 1

	// 1001チェキ
	maxMsgCnt := stng.STUB_MESSAGE_COUNT()
	if resnum > maxMsgCnt {
		err = fmt.Errorf("%d: これ以上書き込めません。。。", resnum)
		return
	}

	// エンティティ更新
	if resnum == maxMsgCnt {
		board.Subjects[pos].MessageCount = resnum + 1 // 1000だったら1001にしてしまう
	} else {
		board.Subjects[pos].MessageCount = resnum
	}
	board.Subjects[pos].LastModified = now
	board.WriteCount++

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
func createDat(name string, mail string, date time.Time, id string, message string, title string) *dat.Entity {
	dat := &dat.Entity{}
	writeDat(dat, dat_format, name, mail, date, id, message, title)
	return dat
}

// append dat. line: 2..
func appendDat(dat *dat.Entity,
	name string, mail string, date time.Time, id string, message string) {

	writeDat(dat, dat_format, name, mail, date, id, message, "")
}

func writeDat(dat *dat.Entity, format string,
	name string, mail string, date time.Time, id string, message string, title string) {

	wr := bytes.NewBuffer(dat.Bytes)
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	// 2行目以降はスレタイは無し
	fmt.Fprintf(wr, format,
		// 名前: トリップの関係でhtml.EscapeStringはトリップのところでやる
		escapeDat(name),
		escapeDat(html.EscapeString(mail)), // メール
		date.Format(dat_date_layout),       // 年月日
		week_days_jp[date.Weekday()],       // 曜
		date.Format(dat_time_layout),       // 時分秒
		id,                                 // ID
		escapeDatMessage(html.EscapeString(message)), // 本文
		escapeDat(html.EscapeString(title)),          // スレタイ
	)

	dat.Bytes = wr.Bytes()
	dat.LastModified = date
}

func escapeDatMessage(str string) string {
	return strings.ReplaceAll(str, "\n", "<br>")
}

func escapeDat(str string) string {
	s := strings.ReplaceAll(str, "\n", "")
	return s
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

// データストアからエンティティを取得しjsonとして返す
func (sv *BoardService) MakeSubjectJson(boardName string, limit int) (_ []byte, err error) {

	// Creates a Key instance.
	key := sv.repo.BoardKey(boardName)

	// Gets a Board
	e := new(board.Entity)
	if err = sv.repo.GetBoard(key, e); err != nil {
		return
	}

	jsonObj := &jboard.Object{
		Subjects: []jboard.Subject{},
		Precure:  sv.StartedAt().Unix(),
	}

	ln := len(e.Subjects)
	for i := 0; i < limit && i < ln; i++ {
		var sbj jboard.Subject
		sbj.ThreadKey = e.Subjects[i].ThreadKey
		sbj.ThreadTitle = e.Subjects[i].ThreadTitle
		sbj.MessageCount = e.Subjects[i].MessageCount
		sbj.LastModified = e.Subjects[i].LastModified.Format("2006/01/02 15:04:05")
		jsonObj.Subjects = append(jsonObj.Subjects, sbj)
	}

	return json.Marshal(jsonObj)
}
