package errors

import (
	"errors"
	"fmt"
)

// Err common error type able to specify fatal or not
type Err struct {
	Fatal bool // Fatal severity level, 'true' -> fatal error, should stop further processing
	Msg   error
}

func (e *Err) Error() string {
	ftl := ""
	if e.Fatal {
		ftl = "[FATAL]"
	}
	return fmt.Sprintf("%v%v", ftl, e.Msg)
}

func IsFatal(err error) bool {
	if e, ok := err.(*Err); ok {
		return e.Fatal // for error type Err, check the e.Fatal flag
	}
	return true
}

func Fatal(err string) *Err {
	return &Err{
		Fatal: true,
		Msg:   errors.New(err),
	}
}

func Fatalf(format string, a ...any) *Err {
	return &Err{
		Fatal: true,
		Msg:   fmt.Errorf(format, a...),
	}
}

func NonFatal(err string) *Err {
	return &Err{
		Fatal: false,
		Msg:   errors.New(err),
	}
}

func NonFatalf(format string, a ...any) *Err {
	return &Err{
		Fatal: false,
		Msg:   fmt.Errorf(format, a...),
	}
}

func Errors(errs ...error) (err error) {
	lgth := len(errs)
	if lgth > 1 {
		err = fmt.Errorf("errors:\n - %v", errs[0])
		for _, e := range errs[1:] {
			err = fmt.Errorf("%v\n - %v", err, e)
		}
	} else if lgth > 0 {
		err = fmt.Errorf("%v", errs[0])
	}
	return
}
