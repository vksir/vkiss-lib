package errutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrNotFound = errors.New("not found")
	ErrBusy     = errors.New("busy")
)

func wrap(err error) error {
	var pcs [1]uintptr
	// skip [runtime.Callers, this function]
	runtime.Callers(3, pcs[:])
	f, _ := runtime.CallersFrames(pcs[:]).Next()
	dir, file := filepath.Split(f.File)
	return fmt.Errorf("%s/%s:%d %w", filepath.Base(dir), file, f.Line, err)
}

func Wrap(err error) error {
	return wrap(err)
}

func WrapF(format string, a ...any) error {
	return wrap(fmt.Errorf(format, a...))
}

func WrapPath(op string, path string, err error) error {
	return wrap(&os.PathError{Op: op, Path: path, Err: err})
}

func WrapNotFound(dst string) error {
	return wrap(fmt.Errorf("%w: %s", ErrNotFound, dst))
}

func Check(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error exit:", wrap(err))
		os.Exit(1)
	}
}
