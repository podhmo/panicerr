# panicerr

convert panic to error

```go
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
```
