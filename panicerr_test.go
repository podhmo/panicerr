package panicerr

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestSimple(t *testing.T) {
	f := func() (err error) {
		defer Recoverer("hi", &err)
		panic("ho")
	}
	err := f()
	if err == nil {
		t.Errorf("expected error, but nil, something wrong")
	}

	if expected, actual := err.Error(), "hi ho"; expected != actual {
		t.Errorf("expected message is %q, but actual is %q", expected, actual)
	}
}
func TestSimplePanicWithError(t *testing.T) {
	this := fmt.Errorf("THIS")
	f := func() (err error) {
		defer Recoverer("", &err)
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
func TestSimpleWithGoroutine(t *testing.T) {
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() (err error) {
		defer Recoverer("ok", &err)
		t.Log("ok")
		return nil
	})
	g.Go(func() (err error) {
		defer Recoverer("ng", &err)
		t.Log("ng")
		panic("hmm")
	})
	err := g.Wait()
	if err == nil {
		t.Errorf("expected error, but nil, something wrong")
	}
}
