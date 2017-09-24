package mf

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var onces = &sync.Map{}

// Deps runs the given functions as dependencies of the calling function.
// Dependencies must only be func() or func() error.  The function calling Deps
// is guaranteed that all dependent functions will be run exactly once when Deps
// returns.  Dependent functions may in turn declare their own dependencies
// using Deps. Each dependency is run in their own goroutines.
func Deps(fns ...interface{}) {
	for _, f := range fns {
		switch f.(type) {
		case func(), func() error:
			// ok
		default:
			panic(errors.Errorf("Invalid type for dependent function: %T. Dependencies must be func() or func() error", f))
		}
	}

	var errs []string
	mu := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	for _, f := range fns {
		fn := addDep(f)
		wg.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					mu.Lock()
					errs = append(errs, fmt.Sprint(err))
					mu.Unlock()
				}
				wg.Done()
			}()
			if err := fn.run(); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprint(err))
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if len(errs) > 0 {
		panic(strings.Join(errs, "\n"))
	}
}

func addDep(f interface{}) *onceFun {
	var fn func() error
	switch f := f.(type) {
	case func():
		fn = func() error { f(); return nil }
	case func() error:
		fn = f
	}

	n := name(f)
	of, _ := onces.LoadOrStore(n, &onceFun{
		fn: fn,
	})
	return of.(*onceFun)
}

func name(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type onceFun struct {
	once sync.Once
	fn   func() error
}

func (o *onceFun) run() error {
	var err error
	o.once.Do(func() {
		err = o.fn()
	})
	return err
}
