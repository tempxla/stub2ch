package handle

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	boardName, err := process(requireOne(r, "bbs"), maxLen(10), between("0", "zzzzzzzzzz"))
	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "bbs", err), http.StatusBadRequest)
		return "", false
	}
	return boardName, true
}

func requireThreadKey(w http.ResponseWriter, r *http.Request) (string, bool) {
	threadKey, err := process(requireOne(r, "key"), maxLen(10), between("0000000000", "9999999999"))
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

func requireName(w http.ResponseWriter, r *http.Request) (string, bool) {
	name, err := process(requireOne(r, "FROM"), url.QueryUnescape)
	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "FROM", err), http.StatusBadRequest)
		return "", false
	}
	return name, true
}

func requireMail(w http.ResponseWriter, r *http.Request) (string, bool) {
	mail, err := process(requireOne(r, "mail"), url.QueryUnescape)
	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "mail", err), http.StatusBadRequest)
		return "", false
	}
	return mail, true
}

func requireMessage(w http.ResponseWriter, r *http.Request) (string, bool) {
	message, err := process(requireOne(r, "MESSAGE"), url.QueryUnescape, notBlank)
	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "MESSAGE", err), http.StatusBadRequest)
		return "", false
	}
	return message, true
}

func requireTitle(w http.ResponseWriter, r *http.Request) (string, bool) {
	title, err := process(requireOne(r, "subject"), url.QueryUnescape, notBlank)
	if err != nil {
		http.Error(w, fmt.Sprintf(param_error_format, "subject", err), http.StatusBadRequest)
		return "", false
	}
	return title, true
}