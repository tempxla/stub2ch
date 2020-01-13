package handle

import (
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/bbscfg"
	"github.com/tempxla/stub2ch/internal/app/util"
	"html"
	"net/http"
	"strings"
	"unicode/utf8"
)

func requireOne(r *http.Request, name string) func() (string, error) {
	return func() (str string, err error) {
		if param, ok := r.PostForm[name]; !ok {
			err = fmt.Errorf("missing")
		} else if len(param) == 0 {
			err = fmt.Errorf("empty")
		} else if len(param) != 1 {
			err = fmt.Errorf("too many")
		} else {
			str = param[0]
		}
		return
	}
}

func notEmpty(s string) (str string, err error) {
	if s == "" {
		err = fmt.Errorf("0 byte")
	} else {
		str = s
	}
	return
}

func notBlank(s string) (str string, err error) {
	if strings.TrimSpace(s) == "" {
		err = fmt.Errorf("blank")
	} else {
		str = s
	}
	return
}

func between(a, b string) func(string) (string, error) {
	return func(s string) (str string, err error) {
		if s < a || b < s {
			err = fmt.Errorf("%s < %s or %s < %s", s, a, b, s)
		} else {
			str = s
		}
		return
	}
}

func maxLen(max int) func(string) (string, error) {
	return func(s string) (str string, err error) {
		if utf8.RuneCountInString(s) > max {
			err = fmt.Errorf(" len %s > %d", s, max)
		} else {
			str = s
		}
		return
	}
}

func maxByte(max int) func(string) (string, error) {
	return func(s string) (str string, err error) {
		if len(s) > max {
			err = fmt.Errorf(" byte len %s > %d", s, max)
		} else {
			str = s
		}
		return
	}
}

// \t, \n 以外の制御文字を削除する
// \t, \nは後の工程でなんとかする。
func delBadChar(s string) (string, error) {
	f := func(r rune) rune {
		if '\u0000' <= r && r < '\u0009' { // NUL <= r < HT
			return ' '
		}
		if '\u000A' < r && r <= '\u001F' { // LF < r < US
			return ' '
		}
		if '\u007F' == r { // r == DEL
			return ' '
		}
		return r
	}
	return strings.Map(f, s), nil
}

func htmlUnescapeString(s string) (string, error) {
	return html.UnescapeString(s), nil
}

func sjisToUtf8String(s string) (string, error) {
	return util.SJIStoUTF8String(s), nil
}

// トリップの関係でhtml.EscapeStringはここでやる
func trip(s string) (string, error) {
	idx := strings.IndexRune(s, '#')
	if idx != -1 {
		// トリップじゃい
		trip := util.ComputeTrip(util.UTF8toSJISString(s[idx+1:]))
		name := strings.ReplaceAll(html.EscapeString(s[:idx]), "◆", "◇")
		return fmt.Sprintf("%s </b>◆%s <b>", name, trip), nil
	} else {
		// トリップ無し
		name := strings.ReplaceAll(html.EscapeString(s), "◆", "◇")
		return name, nil
	}
}

func trimWhitespace(s string) (string, error) {
	return strings.Trim(s, "\t\n "), nil
}

func process(src func() (string, error),
	funcs ...func(string) (string, error)) (s string, e error) {

	s, e = src()
	if e != nil {
		return
	}

	for _, f := range funcs {
		s, e = f(s)
		if e != nil {
			return
		}
	}
	return
}

func requireBoardName(w http.ResponseWriter, r *http.Request) (string, bool) {
	boardName, err := process(requireOne(r, "bbs"),
		maxLen(10),
		between("0", "zzzzzzzzzz"),
	)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "bbs", err), http.StatusBadRequest)
		return "", false
	}
	return boardName, true
}

func requireThreadKey(w http.ResponseWriter, r *http.Request) (string, bool) {
	threadKey, err := process(requireOne(r, "key"),
		maxLen(10),
		between("0000000000", "9999999999"),
	)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "key", err), http.StatusBadRequest)
		return "", false
	}
	return threadKey, true
}

func requireTime(w http.ResponseWriter, r *http.Request) (string, bool) {
	t, err := process(requireOne(r, "time"),
		notEmpty,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "time", err), http.StatusBadRequest)
		return "", false
	}
	return t, true
}

func requireName(w http.ResponseWriter, r *http.Request, setting bbscfg.Setting) (string, bool) {
	name, err := process(requireOne(r, "FROM"),
		maxByte(setting.BBS_NAME_COUNT()),
		sjisToUtf8String,
		trip, // 制御文字とかどうなるんやろ＞トリップ
		delBadChar,
		trimWhitespace,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "FROM", err), http.StatusBadRequest)
		return "", false
	}
	if name == "" {
		name = setting.BBS_NONAME_NAME()
	}
	return name, true
}

func requireMail(w http.ResponseWriter, r *http.Request, setting bbscfg.Setting) (string, bool) {
	mail, err := process(requireOne(r, "mail"),
		maxByte(setting.BBS_MAIL_COUNT()),
		sjisToUtf8String,
		delBadChar,
		trimWhitespace,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "mail", err), http.StatusBadRequest)
		return "", false
	}
	return mail, true
}

func requireMessage(w http.ResponseWriter, r *http.Request, setting bbscfg.Setting) (string, bool) {
	message, err := process(requireOne(r, "MESSAGE"),
		maxByte(setting.BBS_MESSAGE_COUNT()),
		sjisToUtf8String,
		delBadChar,
		notBlank,
		trimWhitespace,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "MESSAGE", err), http.StatusBadRequest)
		return "", false
	}
	return message, true
}

func requireTitle(w http.ResponseWriter, r *http.Request, setting bbscfg.Setting) (string, bool) {
	title, err := process(requireOne(r, "subject"),
		maxByte(setting.BBS_SUBJECT_COUNT()),
		sjisToUtf8String,
		delBadChar,
		trimWhitespace,
		notBlank,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "subject", err), http.StatusBadRequest)
		return "", false
	}
	return title, true
}

func requireReferer(w http.ResponseWriter, r *http.Request, boardName string) (string, bool) {
	ref := r.Referer()
	if !strings.Contains(ref, r.Host) || !strings.Contains(ref, boardName) {
		http.Error(w, fmt.Sprintf(param_error_format, "referer", "BAD"), http.StatusBadRequest)
		return "", false
	}
	return ref, true
}
