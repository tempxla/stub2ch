package errors

import (
	goerr "errors"
)

var (
	NOT_MODIFIED = goerr.New("Not Modified")
)
