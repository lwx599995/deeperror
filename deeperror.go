package deeperror

import (
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strings"
)

type deeperror struct {
	err      error
	filename string
	line     int
}

func getPosition(callDepth int) (filename string, line int) {
	pc := make([]uintptr, 1)
	n := runtime.Callers(callDepth+2, pc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(pc).Next()
	return frame.File, frame.Line
}

func (e *deeperror) setPosition(callDepth int) {
	e.filename, e.line = getPosition(callDepth + 1)
}

func (e *deeperror) markPosition() {
	const format = "\n--> %s:%d\t%v"
	if len(e.filename) > 0 {
		dir, file := path.Split(e.filename)
		// record base filename and parent directory only
		e.filename = path.Join(path.Base(dir), file)
	}
	e.err = fmt.Errorf(format, e.filename, e.line, e.err)
}

func (e *deeperror) getDeepestError() (errMsg string) {
	const (
		pattern = `-->.*\.go:[0-9]*`
		newline = "\n"
		tab     = "\t"
	)
	if e.err == nil {
		return
	}
	errMsg = e.err.Error()
	lines := strings.Split(errMsg, newline)
	if len(lines) == 0 {
		return
	}
	lastline := lines[len(lines)-1]
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return
	}
	lastline = strings.TrimPrefix(lastline, reg.FindString(lastline))
	errMsg = strings.TrimPrefix(lastline, tab)
	return
}

// WithContext adds position and annotation to an error. However, the error type will missing.
func WithContext(err error, annotation string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	e := &deeperror{
		err: err,
	}
	e.setPosition(1)
	if len(args) == 0 {
		e.err = fmt.Errorf("%s %v", annotation, e.err)
	} else {
		e.err = fmt.Errorf(fmt.Sprintf(annotation, args...)+" %v", e.err)
	}
	e.markPosition()
	return e.err
}

// WithPosition adds position to an error. However, the error type will missing.
func WithPosition(err error) error {
	if err == nil {
		return nil
	}
	e := &deeperror{
		err: err,
	}
	e.setPosition(1)
	e.markPosition()
	return e.err
}

// DeepestError returns the deepest error string.
// If require origin error string, call WithPosition(<origin error>) when annotate the origin error.
func DeepestError(err error) (errstr string) {
	if err == nil {
		return
	}
	e := &deeperror{
		err: err,
	}
	return e.getDeepestError()
}
