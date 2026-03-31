package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Err error type with the boolean property 'Fatal', also supports multiple errors.
type Err struct {
	Fatal  bool // Fatal severity level, 'true' -> fatal error, should stop further processing
	Errors []error
}

func (e *Err) Error() string {
	switch len(e.Errors) {
	case 0:
		if e.Fatal {
			return "Fatal error"
		}
		return "Error"
	case 1:
		msg := e.Errors[0].Error()
		ftl := ""
		if e.Fatal {
			if msg == "" {
				return "Fatal error"
			}
			if msg[0] != '[' {
				ftl = "[FATAL] "
			} else {
				ftl = "[FATAL]"
			}
		}
		return fmt.Sprintf("%v%v", ftl, msg)
	default:
		var sb strings.Builder
		ftl := ""
		if e.Fatal {
			ftl = "[FATAL] "
		}
		fmt.Fprintf(&sb, "%verrors:\n - %v", ftl, e.Errors[0])
		for _, err := range e.Errors[1:] {
			fmt.Fprintf(&sb, "\n - %v", err)
		}
		return sb.String()
	}
}

func New(fatal bool, errs ...string) (r *Err) {
	r = &Err{
		Fatal: fatal,
	}
	for _, err := range errs {
		r.Errors = append(r.Errors, errors.New(err))
	}
	return
}

func Count(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(*Err); ok {
		return len(e.Errors)
	}
	return 1
}

func IsFatal(err error) bool {
	if e, ok := err.(*Err); ok {
		return e.Fatal // for error type Err, check the e.Fatal flag
	}
	return true
}

func Fatal(err string) *Err {
	return &Err{
		Fatal:  true,
		Errors: []error{errors.New(err)},
	}
}

func Fatalf(format string, a ...any) *Err {
	return &Err{
		Fatal:  true,
		Errors: []error{fmt.Errorf(format, a...)},
	}
}

func NonFatal(err string) *Err {
	return &Err{
		Fatal:  false,
		Errors: []error{errors.New(err)},
	}
}

func NonFatalf(format string, a ...any) *Err {
	return &Err{
		Fatal:  false,
		Errors: []error{fmt.Errorf(format, a...)},
	}
}

// Append append new errors in 'errs' to 'base'.
// - if 'base' is of type *Err, 'Fatal' in 'base' will be used.
// - if 'base' is not of type *Err:
//   - 'base' will first be converted to *Err with 'Fatal' default to false.
//   - if 'errs' is not empty and the first error in 'errs' is of type *Err, the 'Fatal' value in the first error will be used.
func Append(base error, errs ...error) (r *Err) {
	var isErr bool
	if base != nil {
		if e, ok := base.(*Err); !ok {
			r = &Err{ // default is non-fatal
				Errors: []error{base},
			}
		} else {
			isErr = true
			r = e
		}
	}

	if len(errs) > 0 {
		if r == nil {
			r = &Err{} // default is non-fatal
		}
		if !isErr { // if is *Err, err.Fatal follows the value in base
			if e, ok := errs[0].(*Err); ok {
				r.Fatal = e.Fatal
			}
		}
		// err.Errors = append(err.Errors, errs...)
		for _, err := range errs {
			if e, ok := err.(*Err); ok {
				r.Errors = append(r.Errors, e.Errors...)
			} else {
				r.Errors = append(r.Errors, err)
			}
		}
	}
	return
}

func Appendf(base error, format string, a ...any) *Err {
	return Append(base, fmt.Errorf(format, a...))
}
