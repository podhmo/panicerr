package panicerr

import (
	"errors"
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	f := func() (err error) {
		defer Recoverer(&err)
		panic("ho")
	}
	err := f()
	if err == nil {
		t.Errorf("expected error, but nil, something wrong")
	}
}
func TestSimplePanicWithError(t *testing.T) {
	this := fmt.Errorf("THIS")
	f := func() (err error) {
		defer Recoverer(&err)
		panic(this)
	}
	err := f()
	if err == nil {
		t.Errorf("expected error, but nil, something wrong")
	}
	if !errors.Is(err, this) {
		t.Errorf("expected error is %[1]T:%[1]v, but got %[2]T:%[2]v, unwrap is failed", this, err)
	}
}
