package testutil

import (
	"github.com/tempxla/stub2ch/configs/app/setting"
)

type SettingStub struct{}

func (_ *SettingStub) BBS_TITLE() string            { return "VIP＠スタブ" }
func (_ *SettingStub) BBS_NONAME_NAME() string      { return "⊂二二二( ^ω^)二⊃" }
func (_ *SettingStub) BBS_UNICODE() string          { return "pass" }
func (_ *SettingStub) BBS_SUBJECT_COUNT() int       { return 128 }
func (_ *SettingStub) BBS_NAME_COUNT() int          { return 96 }
func (_ *SettingStub) BBS_MAIL_COUNT() int          { return 32 }
func (_ *SettingStub) BBS_MESSAGE_COUNT() int       { return 4096 }
func (_ *SettingStub) BBS_THREAD_TATESUGI() int     { return 8 }
func (_ *SettingStub) BBS_SLIP() string             { return "verbose" }
func (_ *SettingStub) BBS_DISP_IP() string          { return "" }
func (_ *SettingStub) BBS_FORCE_ID() string         { return "checked" }
func (_ *SettingStub) BBS_NO_ID() string            { return "" }
func (_ *SettingStub) BBS_JP_CHECK() string         { return "" }
func (_ *SettingStub) BBS_4WORLD() string           { return "" }
func (_ *SettingStub) BBS_YMD_WEEKS() string        { return "" }
func (_ *SettingStub) BBS_ARR() string              { return "" }
func (_ *SettingStub) BBS_SOKO() string             { return "ononon" }
func (_ *SettingStub) BBS_DISP_MSEC() int           { return 3 }
func (_ *SettingStub) STUB_WRITE_ENTITY_LIMIT() int { return 5000 }
func (_ *SettingStub) STUB_THREAD_COUNT() int       { return 500 }
func (_ *SettingStub) STUB_MESSAGE_COUNT() int      { return 1000 }
func (_ *SettingStub) STUB_DAT_CAPACITY() int       { return 500 * 1024 }

func NewSettingStub() setting.BBS {
	return &SettingStub{}
}
