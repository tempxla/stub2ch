package bbscfg

import (
	"bytes"
	"fmt"
)

var (
	settings = map[string]Setting{
		"news4vip": &News4vip{},
		"poverty":  &Poverty{},
	}
)

type Setting interface {
	BBS_TITLE() string
	BBS_NONAME_NAME() string
	BBS_UNICODE() string
	BBS_SUBJECT_COUNT() int
	BBS_NAME_COUNT() int
	BBS_MAIL_COUNT() int
	BBS_MESSAGE_COUNT() int
	BBS_THREAD_TATESUGI() int
	BBS_SLIP() string
	BBS_DISP_IP() string
	BBS_FORCE_ID() string
	BBS_NO_ID() string
	BBS_JP_CHECK() string
	BBS_4WORLD() string
	BBS_YMD_WEEKS() string
	BBS_ARR() string
	BBS_SOKO() string
	BBS_DISP_MSEC() int
	STUB_WRITE_ENTITY_LIMIT() int // 一日の書き込み数制限
	STUB_THREAD_COUNT() int       // 許容スレッド数
	STUB_MESSAGE_COUNT() int      // 許容レス数
	STUB_DAT_CAPACITY() int       // 許容バイト数
}

func GetSetting(boardName string) Setting {
	return settings[boardName]
}

func MakeSettingTxt(setting Setting) []byte {
	sb := &bytes.Buffer{}
	fmt.Fprintf(sb, "BBS_TITLE=%s\n", setting.BBS_TITLE())
	fmt.Fprintf(sb, "BBS_NONAME_NAME=%s\n", setting.BBS_NONAME_NAME())
	fmt.Fprintf(sb, "BBS_UNICODE=%s\n", setting.BBS_UNICODE())
	fmt.Fprintf(sb, "BBS_SUBJECT_COUNT=%d\n", setting.BBS_SUBJECT_COUNT())
	fmt.Fprintf(sb, "BBS_NAME_COUNT=%d\n", setting.BBS_NAME_COUNT())
	fmt.Fprintf(sb, "BBS_MAIL_COUNT=%d\n", setting.BBS_MAIL_COUNT())
	fmt.Fprintf(sb, "BBS_MESSAGE_COUNT=%d\n", setting.BBS_MESSAGE_COUNT())
	fmt.Fprintf(sb, "BBS_THREAD_TATESUGI=%d\n", setting.BBS_THREAD_TATESUGI())
	fmt.Fprintf(sb, "BBS_SLIP=%s\n", setting.BBS_SLIP())
	fmt.Fprintf(sb, "BBS_DISP_IP=%s\n", setting.BBS_DISP_IP())
	fmt.Fprintf(sb, "BBS_FORCE_ID=%s\n", setting.BBS_FORCE_ID())
	fmt.Fprintf(sb, "BBS_NO_ID=%s\n", setting.BBS_NO_ID())
	fmt.Fprintf(sb, "BBS_JP_CHECK=%s\n", setting.BBS_JP_CHECK())
	fmt.Fprintf(sb, "BBS_4WORLD=%s\n", setting.BBS_4WORLD())
	fmt.Fprintf(sb, "BBS_YMD_WEEKS=%s\n", setting.BBS_YMD_WEEKS())
	fmt.Fprintf(sb, "BBS_ARR=%s\n", setting.BBS_ARR())
	fmt.Fprintf(sb, "BBS_SOKO=%s\n", setting.BBS_SOKO())
	fmt.Fprintf(sb, "BBS_DISP_MSEC=%d\n", setting.BBS_DISP_MSEC())

	return sb.Bytes()
}
