package setting

type BBS interface {
	_A() string
	BBS_TITLE() string
	BBS_TITLE_ORIG() string
	BBS_TITLE_PICTURE() string
	BBS_TITLE_COLOR() string
	BBS_TITLE_LINK() string
	BBS_BG_COLOR() string
	BBS_BG_PICTURE() string
	BBS_NONAME_NAME() string
	BBS_MAKETHREAD_COLOR() string
	BBS_MENU_COLOR() string
	BBS_THREAD_COLOR() string
	BBS_TEXT_COLOR() string
	BBS_NAME_COLOR() string
	BBS_LINK_COLOR() string
	BBS_ALINK_COLOR() string
	BBS_VLINK_COLOR() string
	BBS_THREAD_NUMBER() int
	BBS_CONTENTS_NUMBER() int
	BBS_LINE_NUMBER() int
	BBS_MAX_MENU_THREAD() int
	BBS_SUBJECT_COLOR() string
	BBS_UNICODE() string
	BBS_NAMECOOKIE_CHECK() string
	BBS_MAILCOOKIE_CHECK() string
	BBS_SUBJECT_COUNT() int
	BBS_NAME_COUNT() int
	BBS_MAIL_COUNT() int
	BBS_MESSAGE_COUNT() int
	BBS_THREAD_TATESUGI() int
	NANASHI_CHECK() string
	BBS_PROXY_CHECK() string
	BBS_OVERSEA_PROXY() string
	BBS_RAWIP_CHECK() string
	BBS_SLIP() string
	BBS_DISP_IP() string
	BBS_FORCE_ID() string
	BBS_BE_ID() string
	BBS_BE_TYPE2() string
	BBS_NO_ID() string
	BBS_JP_CHECK() string
	BBS_4WORLD() string
	BBS_YMD_WEEKS() string
	BBS_BBX_PASS() string
	BBS_SOKO() string
	BBS_USE_VIPQ2() int
	BBS_DISP_MSEC() int
}

func GetSetting(boardName string) BBS {
	switch boardName {
	case "news4vip":
		return &News4vip{}
	case "poverty":
		return &Poverty{}
	default:
		return &Stub2ch{}
	}
}
