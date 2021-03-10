package panicerr

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

	if expected, actual := err.Error(), "ng hmm"; expected != actual {
		t.Errorf("expected message is %q, but actual is %q", expected, actual)
	}
}

func TestVerboseFormat(t *testing.T) {
	err := func() (err error) {
		defer Recoverer("err", &err)
		panic("content")
	}()

	t.Run("s", func(t *testing.T) {
		if expected, actual := "!! err content", fmt.Sprintf("!! %s", err); expected != actual {
			t.Errorf("expected message is %q, but actual is %q", expected, actual)
		}
	})
	t.Run("v", func(t *testing.T) {
		if expected, actual := "!! err content", fmt.Sprintf("!! %v", err); expected != actual {
			t.Errorf("expected message is %q, but actual is %q", expected, actual)
		}
	})
	t.Run("+v", func(t *testing.T) {
		if expected, actual := "!! err content\ngoroutine ", fmt.Sprintf("!! %+v", err); !strings.HasPrefix(actual, expected) {
			t.Errorf("expected message is %q, but actual is %q", expected, actual)
		}
	})

	w := &wrap{inner: err, tag: "@@@@"}
	t.Run("wrap s", func(t *testing.T) {
		if expected, actual := "!!  err content", fmt.Sprintf("!! %s", w); expected != actual {
			t.Errorf("expected message is %q, but actual is %q", expected, actual)
		}
	})
	t.Run("wrap v", func(t *testing.T) {
		// This is go's securty feature? ()
		if expected, actual := "!! ", fmt.Sprintf("!! %v", w); expected != actual {
			t.Errorf("expected message is %q, but actual is %q", expected, actual)
		}
	})
	t.Run("wrap +v", func(t *testing.T) {
		if expected, actual := "!! @@@@ err content\ngoroutine ", fmt.Sprintf("!! %+v", w); !strings.HasPrefix(actual, expected) {
			t.Errorf("expected message is %q, but actual is %q", expected, actual)
		}

	})
}

type wrap struct {
	tag   string
	inner error
}

func (w wrap) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s %+v", w.tag, w.inner)
			return
		}
		fallthrough
	default:
	case 's':
		fmt.Fprintf(s, "%s", w.inner)
	case 'q':
		fmt.Fprintf(s, "%q", w.inner)
	}
}
