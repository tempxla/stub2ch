package setting

type Stub2ch struct{}

func (_ *Stub2ch) BBS_TITLE() string        { return "スタブ" }
func (_ *Stub2ch) BBS_NONAME_NAME() string  { return "名無しさん＠スタブ" }
func (_ *Stub2ch) BBS_UNICODE() string      { return "pass" }
func (_ *Stub2ch) BBS_SUBJECT_COUNT() int   { return 128 }
func (_ *Stub2ch) BBS_NAME_COUNT() int      { return 96 }
func (_ *Stub2ch) BBS_MAIL_COUNT() int      { return 32 }
func (_ *Stub2ch) BBS_MESSAGE_COUNT() int   { return 4096 }
func (_ *Stub2ch) BBS_THREAD_TATESUGI() int { return 8 }
func (_ *Stub2ch) BBS_SLIP() string         { return "verbose" }
func (_ *Stub2ch) BBS_DISP_IP() string      { return "" }
func (_ *Stub2ch) BBS_FORCE_ID() string     { return "checked" }
func (_ *Stub2ch) BBS_NO_ID() string        { return "" }
func (_ *Stub2ch) BBS_JP_CHECK() string     { return "" }
func (_ *Stub2ch) BBS_4WORLD() string       { return "" }
func (_ *Stub2ch) BBS_YMD_WEEKS() string    { return "" }
func (_ *Stub2ch) BBS_ARR() string          { return "" }
func (_ *Stub2ch) BBS_SOKO() string         { return "ononon" }
func (_ *Stub2ch) BBS_DISP_MSEC() int       { return 3 }
