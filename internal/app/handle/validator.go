package handle

import (
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/setting"
	"github.com/tempxla/stub2ch/internal/app/util"
	"html"
	"math/rand"
	"net/http"
	"strings"
	"time"
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
		if len(s) > max {
			err = fmt.Errorf(" len %s > %d", s, max)
		} else {
			str = s
		}
		return
	}
}

func maxByte(max int) func(string) (string, error) {
	return func(s string) (str string, err error) {
		if len([]byte(s)) > max {
			err = fmt.Errorf(" byte len %s > %d", s, max)
		} else {
			str = s
		}
		return
	}
}

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
		trip := util.ComputeTrip(util.UTF8toSJISString(s[idx+1:]))
		return fmt.Sprintf("%s </b>◆%s <b>",
			strings.ReplaceAll(html.EscapeString(s[:idx]), "◆", "◇"), trip), nil
	} else {
		return strings.ReplaceAll(html.EscapeString(s), "◆", "◇"), nil
	}
}

var omikuji_list = [...]string{"豚", "おっさん", "犬", "にゃあ", "男の娘", "神", "女神"}

func omikuji(s string) (string, error) {
	var idx int
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(1000)
	if x < 10 {
		idx = rand.Intn(len(omikuji_list))
	} else if x < 100 {
		idx = rand.Intn(len(omikuji_list) - 2)
	} else {
		idx = 0
	}

	return strings.ReplaceAll(s, "!omikuji", fmt.Sprintf("</b>【%s】<b>", omikuji_list[idx])), nil
}

func dama(s string) (string, error) {
	var money int
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(1000)
	if x < 10 {
		money = rand.Intn(100 * 1000)
	} else if x < 100 {
		money = rand.Intn(10 * 1000)
	} else {
		money = rand.Intn(1 * 1000)
	}
	return strings.ReplaceAll(s, "!dama", fmt.Sprintf("</b>【%d円】<b>", money*1000)), nil
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
		maxLen(10), between("0", "zzzzzzzzzz"))

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "bbs", err), http.StatusBadRequest)
		return "", false
	}
	return boardName, true
}

func requireThreadKey(w http.ResponseWriter, r *http.Request) (string, bool) {
	threadKey, err := process(requireOne(r, "key"),
		maxLen(10), between("0000000000", "9999999999"))

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "key", err), http.StatusBadRequest)
		return "", false
	}
	return threadKey, true
}

func requireTime(w http.ResponseWriter, r *http.Request) (string, bool) {
	t, err := process(requireOne(r, "time"), notEmpty)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "time", err), http.StatusBadRequest)
		return "", false
	}
	return t, true
}

func requireName(w http.ResponseWriter, r *http.Request, setting setting.BBS) (string, bool) {
	name, err := process(requireOne(r, "FROM"),
		maxByte(setting.BBS_NAME_COUNT()),
		sjisToUtf8String,
		trip, // 制御文字とかどうなるんやろ＞トリップ
		dama,
		omikuji,
		delBadChar,
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

func requireMail(w http.ResponseWriter, r *http.Request, setting setting.BBS) (string, bool) {
	mail, err := process(requireOne(r, "mail"),
		maxByte(setting.BBS_MAIL_COUNT()), sjisToUtf8String, delBadChar)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "mail", err), http.StatusBadRequest)
		return "", false
	}
	return mail, true
}

func requireMessage(w http.ResponseWriter, r *http.Request, setting setting.BBS) (string, bool) {
	message, err := process(requireOne(r, "MESSAGE"),
		maxByte(setting.BBS_MESSAGE_COUNT()), sjisToUtf8String, delBadChar, notBlank)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "MESSAGE", err), http.StatusBadRequest)
		return "", false
	}
	return message, true
}

func requireTitle(w http.ResponseWriter, r *http.Request, setting setting.BBS) (string, bool) {
	title, err := process(requireOne(r, "subject"),
		maxByte(setting.BBS_SUBJECT_COUNT()), sjisToUtf8String, delBadChar, notBlank)

	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "subject", err), http.StatusBadRequest)
		return "", false
	}
	return title, true
}
