package bbscfg

type News4vip struct{}

func (_ *News4vip) BBS_TITLE() string            { return "VIP＠スタブ" }
func (_ *News4vip) BBS_NONAME_NAME() string      { return "⊂二二二( ^ω^)二⊃" }
func (_ *News4vip) BBS_UNICODE() string          { return "pass" }
func (_ *News4vip) BBS_SUBJECT_COUNT() int       { return 128 }
func (_ *News4vip) BBS_NAME_COUNT() int          { return 96 }
func (_ *News4vip) BBS_MAIL_COUNT() int          { return 32 }
func (_ *News4vip) BBS_MESSAGE_COUNT() int       { return 4096 }
func (_ *News4vip) BBS_THREAD_TATESUGI() int     { return 8 }
func (_ *News4vip) BBS_SLIP() string             { return "verbose" }
func (_ *News4vip) BBS_DISP_IP() string          { return "" }
func (_ *News4vip) BBS_FORCE_ID() string         { return "checked" }
func (_ *News4vip) BBS_NO_ID() string            { return "" }
func (_ *News4vip) BBS_JP_CHECK() string         { return "" }
func (_ *News4vip) BBS_4WORLD() string           { return "" }
func (_ *News4vip) BBS_YMD_WEEKS() string        { return "" }
func (_ *News4vip) BBS_ARR() string              { return "" }
func (_ *News4vip) BBS_SOKO() string             { return "ononon" }
func (_ *News4vip) BBS_DISP_MSEC() int           { return 3 }
func (_ *News4vip) STUB_WRITE_ENTITY_LIMIT() int { return 4000 } // 5000までいける
func (_ *News4vip) STUB_THREAD_COUNT() int       { return 500 }
func (_ *News4vip) STUB_MESSAGE_COUNT() int      { return 1000 }
func (_ *News4vip) STUB_DAT_CAPACITY() int       { return 500 * 1024 }
