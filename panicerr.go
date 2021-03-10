package panicerr

import (
	"fmt"
	"runtime/debug"
)

var ErrFormat = "%s %s"

type Err struct {
	prefix string
	inner  error
	stack  []byte
}

func (e *Err) Error() string {
	return fmt.Sprintf(ErrFormat, e.prefix, e.inner.Error())
}
func (e *Err) Stack() string {
	// return (*(*string)(unsafe.Pointer(&e.stack))
	return string(e.stack)
}
func (e *Err) Unwrap() error {
	return e.inner
}

func Recoverer(prefix string, err *error) {
	r := recover()
	if r == nil {
		return
	}
	switch r := r.(type) {
	case error:
		*err = &Err{prefix: prefix, inner: r, stack: debug.Stack()}
	default:
		*err = &Err{prefix: prefix, inner: fmt.Errorf("%+v", r), stack: debug.Stack()}
	}
}
