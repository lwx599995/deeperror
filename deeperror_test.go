package deeperror_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lwx599995/deeperror"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

func TestWithContext(t *testing.T) {
	if err := a(); err != nil {
		err = deeperror.WithContext(err, "call %s failed", "a")
		fmt.Printf("%s\tERROR\tTestWithContext: %s\n", time.Now().Format(TIME_FORMAT), err)
	}
}

func TestWithPosition(t *testing.T) {
	if err := b(); err != nil {
		err = deeperror.WithPosition(err)
		fmt.Printf("%s\tERROR\tTestWithPosition: %s\n", time.Now().Format(TIME_FORMAT), err)
	}
}
func TestDeepestError(t *testing.T) {
	if err := b(); err != nil {
		str := deeperror.DeepestError(err)
		fmt.Printf("%s\tINFO\tTestDeepestError: %s\n", time.Now().Format(TIME_FORMAT), str)
	}
	if err := a(); err != nil {
		str := deeperror.DeepestError(err)
		fmt.Printf("%s\tINFO\tTestDeepestError: %s\n", time.Now().Format(TIME_FORMAT), str)
	}
}

func a() error {
	if err := aa(); err != nil {
		return deeperror.WithContext(err, "call aa failed")
	}
	return nil
}

func aa() error {
	if err := aaa(); err != nil {
		return deeperror.WithContext(err, "call aaa failed")
	}
	if err := aab(); err != nil {
		return deeperror.WithContext(err, "call aab failed")
	}
	return nil
}

func aaa() error {
	return nil
}

func aab() error {
	err := errors.New("origin error aab")
	if err != nil {
		return deeperror.WithContext(err, "new error ->")
	}
	return nil
}

func b() error {
	if err := bb(); err != nil {
		return deeperror.WithPosition(err)
	}
	return nil
}

func bb() error {
	if err := bbb(); err != nil {
		return deeperror.WithPosition(err)
	}
	if err := bbc(); err != nil {
		return deeperror.WithPosition(err)
	}
	return nil
}

func bbb() error {
	return nil
}

func bbc() error {
	err := errors.New("origin error bbc")
	if err != nil {
		return deeperror.WithPosition(err)
	}
	return nil
}
