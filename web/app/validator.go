package main

import (
	"fmt"
	"net/http"
	"strings"
)

func requireParam(r *http.Request, name string) (str string, err error) {
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

func require(r *http.Request, name string) func() (string, error) {
	return func() (string, error) {
		return requireParam(r, name)
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

func betweenStr(a, b string) func(string) (string, error) {
	return func(s string) (str string, err error) {
		if s < a || b < s {
			err = fmt.Errorf("%s < %s or %s < %s", s, a, b, s)
		} else {
			str = s
		}
		return
	}
}

func processParam(src func() (string, error),
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
