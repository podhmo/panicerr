# panicerr

convert panic to error

```go
func TestSimpleWithGoroutine(t *testing.T) {
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() (err error) {
		defer Recoverer(&err)
		t.Log("ok")
		return nil
	})
	g.Go(func() (err error) {
		defer Recoverer(&err)
		t.Log("ng")
		panic("hmm")
	})
	err := g.Wait()
	if err == nil {
		t.Errorf("expected error, but nil, something wrong")
	}
}
```
