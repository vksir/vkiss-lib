package errutil

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func wrap(err error) error {
	var pcs [1]uintptr
	// skip [runtime.Callers, this function]
	runtime.Callers(3, pcs[:])
	f, _ := runtime.CallersFrames(pcs[:]).Next()
	function := f.Function[strings.LastIndex(f.Function, "/")+1:]
	return fmt.Errorf("%s:%d > %w", function, f.Line, err)
}

func Wrap(err error) error {
	return wrap(err)
}

func WrapPathErr(op string, path string, err error) error {
	return wrap(&os.PathError{Op: op, Path: path, Err: err})
}
