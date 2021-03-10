package panicerr

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"
)

var (
	ErrFormat        = "%s %s"
	ErrVerboseFormat = "%s %+v\n%s"
)

type Err struct {
	prefix string
	inner  error
	stack  []byte
}

func (e *Err) Error() string {
	return fmt.Sprintf(ErrFormat, e.prefix, e.inner.Error())
}
func (e *Err) Stack() string {
	// // return (*(*string)(unsafe.Pointer(&e.stack))

	// include:	goroutine 23 [running]:
	// trim:  :	runtime/debug.Stack(0x114b92c, 0x3, 0xc00002e520)
	// trim:  :	         /opt/local/lib/go/src/runtime/debug/stack.go:24 +0x9f
	// trim:  :	 github.com/podhmo/panicerr.Recoverer(0x114b983, 0x3, 0xc000091f30)
	// trim:  :	         ~/ghq/github.com/podhmo/panicerr/panicerr.go:70 +0xca
	// trim:  :	 panic(0x1121fa0, 0x1175090)
	// trim:  :	         /opt/local/lib/go/src/runtime/panic.go:969 +0x175
	// include:	github.com/podhmo/panicerr.TestVerboseFormat.func1(0x0, 0x0)
	// include:	         ~/ghq/github.com/podhmo/panicerr/panicerr_test.go:68 +0x87
	// include:	github.com/podhmo/panicerr.TestVerboseFormat(0xc00009ac00)
	// ...
	lines := strings.SplitN(string(e.stack), "\n", 9)
	return lines[0] + lines[8]
}
func (e *Err) Unwrap() error {
	return e.inner
}

func (e *Err) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, ErrVerboseFormat, e.prefix, e.inner, e.Stack())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
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
