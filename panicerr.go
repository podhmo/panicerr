package panicerr

import (
	"fmt"
	"runtime/debug"
)

type Err struct {
	inner error
	stack []byte
}

func (e *Err) Error() string {
	return e.inner.Error()
}
func (e *Err) Stack() string {
	// return (*(*string)(unsafe.Pointer(&e.stack))
	return string(e.stack)
}
func (e *Err) Unwrap() error {
	return e.inner
}

func Recoverer(err *error) {
	r := recover()
	if r == nil {
		return
	}
	switch r := r.(type) {
	case error:
		*err = &Err{inner: r, stack: debug.Stack()}
	default:
		*err = &Err{inner: fmt.Errorf("%+v", r), stack: debug.Stack()}
	}
}
