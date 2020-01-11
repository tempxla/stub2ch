package bbscfg

type Poverty struct{}

func (_ *Poverty) BBS_TITLE() string            { return "嫌儲＠スタブ" }
func (_ *Poverty) BBS_NONAME_NAME() string      { return "（ヽ´ん`）" }
func (_ *Poverty) BBS_UNICODE() string          { return "pass" }
func (_ *Poverty) BBS_SUBJECT_COUNT() int       { return 128 }
func (_ *Poverty) BBS_NAME_COUNT() int          { return 96 }
func (_ *Poverty) BBS_MAIL_COUNT() int          { return 96 }
func (_ *Poverty) BBS_MESSAGE_COUNT() int       { return 4096 }
func (_ *Poverty) BBS_THREAD_TATESUGI() int     { return 8 }
func (_ *Poverty) BBS_SLIP() string             { return "vvvvv" }
func (_ *Poverty) BBS_DISP_IP() string          { return "" }
func (_ *Poverty) BBS_FORCE_ID() string         { return "checked" }
func (_ *Poverty) BBS_NO_ID() string            { return "" }
func (_ *Poverty) BBS_JP_CHECK() string         { return "" }
func (_ *Poverty) BBS_4WORLD() string           { return "" }
func (_ *Poverty) BBS_YMD_WEEKS() string        { return "" }
func (_ *Poverty) BBS_ARR() string              { return "" }
func (_ *Poverty) BBS_SOKO() string             { return "ononon" }
func (_ *Poverty) BBS_DISP_MSEC() int           { return 3 }
func (_ *Poverty) STUB_WRITE_ENTITY_LIMIT() int { return 4000 } // 5000までいける
func (_ *Poverty) STUB_THREAD_COUNT() int       { return 500 }
func (_ *Poverty) STUB_MESSAGE_COUNT() int      { return 1000 }
func (_ *Poverty) STUB_DAT_CAPACITY() int       { return 500 * 1024 }
