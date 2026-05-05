package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Err error type with the boolean property 'Fatal', also supports multiple errors.
type Err struct {
	fatal  bool // Fatal severity level, 'true' -> fatal error, should stop further processing
	errors []error
}

func (e *Err) Error() string {
	switch len(e.errors) {
	case 0:
		if e.fatal {
			return "Fatal error"
		}
		return "Error"
	case 1:
		msg := e.errors[0].Error()
		if e.fatal {
			if msg == "" {
				return "Fatal error"
			} else if msg[0] != '[' {
				return fmt.Sprintf("[FATAL] %v", msg)
			} else {
				return fmt.Sprintf("[FATAL]%v", msg)
			}
		}
		return msg
	default:
		var sb strings.Builder
		if e.fatal {
			fmt.Fprint(&sb, "[FATAL]:")
		} else {
			fmt.Fprint(&sb, "errors:")
		}
		for _, err := range e.errors {
			if _, ok := err.(*Err); ok {
				fmt.Fprintf(&sb, "\n %v.", err)
			} else {
				fmt.Fprintf(&sb, "\n %v", err)
			}
		}
		return sb.String()
	}
}

func New(fatal bool, errs ...string) (r *Err) {
	r = &Err{
		fatal: fatal,
	}
	for _, err := range errs {
		r.errors = append(r.errors, errors.New(err))
	}
	return
}

func Count(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(*Err); ok {
		return len(e.errors)
	}
	return 1
}

func IsFatal(err error) bool {
	if e, ok := err.(*Err); ok {
		return e.fatal // for error type Err, check the e.Fatal flag
	}
	return true
}

func Fatal(err string) *Err {
	return &Err{
		fatal:  true,
		errors: []error{errors.New(err)},
	}
}

func Fatalf(format string, a ...any) *Err {
	return &Err{
		fatal:  true,
		errors: []error{fmt.Errorf(format, a...)},
	}
}

func NonFatal(err string) *Err {
	return &Err{
		fatal:  false,
		errors: []error{errors.New(err)},
	}
}

func NonFatalf(format string, a ...any) *Err {
	return &Err{
		fatal:  false,
		errors: []error{fmt.Errorf(format, a...)},
	}
}

// Append append new errors in 'errs' to 'base'.
// - if 'base' is of type *Err, 'Fatal' in 'base' will be used.
// - if 'base' is not of type *Err:
//   - 'base' will first be converted to *Err with 'Fatal' default to true.
//   - if 'errs' is not empty, the 'Fatal' value of the first non-empty error in 'errs' of type *Err will be used.
func Append(base error, errs ...error) (r *Err) {
	var set bool
	if e, ok := base.(*Err); !ok {
		if base != nil && base.Error() != "" {
			r = &Err{ // default is fatal
				fatal:  true,
				errors: []error{base},
			}
		}
	} else if e != nil {
		set = true
		r = e
	}

	if len(errs) > 0 {
		if r == nil {
			r = &Err{
				fatal: true,
			} // default is fatal
		}
		for _, err := range errs {
			if e, ok := err.(*Err); ok {
				if e != nil && len(e.errors) > 0 {
					if !set {
						r.fatal = e.fatal
						set = true
					}
					r.errors = append(r.errors, e.errors...) // flatten nested *Err
				}
			} else {
				if err != nil && err.Error() != "" {
					r.errors = append(r.errors, err)
				}
			}
		}
	}
	return
}

// Appendf append a new error created by fmt.Errorf().
func Appendf(base error, format string, a ...any) *Err {
	return Append(base, fmt.Errorf(format, a...))
}
