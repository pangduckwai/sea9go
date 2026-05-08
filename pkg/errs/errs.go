package errs

import (
	"errors"
	"fmt"
	"strings"
)

// Err error type with the boolean property 'fatal', supports multiple errors.
type Err struct {
	fatal  bool // Fatal severity level, 'true' -> fatal error, should stop further processing
	errors []error
	labels []string
}

func (e *Err) Error() (str string) {
	lgth := len(e.errors)
	if lgth < 1 {
		if e.fatal {
			str = "error"
		} else {
			str = "non-fatal error"
		}
	} else {
		var bdr strings.Builder
		var msg string
		i0 := 1
		if e.errors[0] != nil && e.errors[0].Error() != "" {
			msg = e.errors[0].Error()
		} else {
			// e.errors[0] is nil or empty string
			if lgth == 2 {
				msg = e.errors[1].Error()
				i0 = 2
			} else if lgth > 2 {
				msg = "errors"
			}
		}
		if !e.fatal {
			if msg == "" {
				fmt.Fprint(&bdr, "non-fatal error")
			} else if msg[0] != '[' {
				fmt.Fprintf(&bdr, "[non-fatal] %v", msg)
			} else {
				fmt.Fprintf(&bdr, "[non-fatal]%v", msg)
			}
		} else {
			if msg == "" {
				fmt.Fprint(&bdr, "error")
			} else {
				fmt.Fprint(&bdr, msg)
			}
		}

		spr := ":\n"
		for _, err := range e.errors[i0:] {
			if _, ok := err.(*Err); ok {
				fmt.Fprintf(&bdr, "%v = %v", spr, err)
			} else {
				fmt.Fprintf(&bdr, "%v - %v", spr, err)
			}
			spr = "\n"
		}

		str = bdr.String()
	}

	if len(e.labels) > 0 {
		var lbl strings.Builder
		for _, msg := range e.labels {
			i0, i1 := 0, len(msg)
			if strings.HasPrefix(msg, "[") {
				i0++
			}
			if strings.HasSuffix(msg, "]") {
				i1--
			}
			if i1 > i0 {
				fmt.Fprintf(&lbl, "[%v]", msg[i0:i1])
			}
		}
		if strings.HasPrefix(str, "[") {
			str = fmt.Sprintf("%v%v", lbl.String(), str)
		} else {
			str = fmt.Sprintf("%v %v", lbl.String(), str)
		}
	}

	return
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
		sub := 0
		if len(e.errors) > 0 && e.errors[0] == nil {
			sub = 1
		}
		return len(e.errors) - sub
	}
	return 1
}

func IsFatal(err error) bool {
	if e, ok := err.(*Err); ok {
		return e.fatal // for error type Err, check the e.fatal flag
	}
	return true // other error types are considered fatal
}

func Fatal(msg string) *Err {
	if msg != "" {
		return &Err{
			fatal:  true,
			errors: []error{errors.New(msg)},
		}
	}
	return &Err{
		fatal: true,
	}
}

func Fatalf(format string, a ...any) *Err {
	if format != "" {
		return &Err{
			fatal:  true,
			errors: []error{fmt.Errorf(format, a...)},
		}
	}
	return &Err{
		fatal: true,
	}
}

func NonFatal(msg string) *Err {
	if msg != "" {
		return &Err{
			fatal:  false,
			errors: []error{errors.New(msg)},
		}
	}
	return &Err{
		fatal: false,
	}
}

func NonFatalf(format string, a ...any) *Err {
	if format != "" {
		return &Err{
			fatal:  false,
			errors: []error{fmt.Errorf(format, a...)},
		}
	}
	return &Err{
		fatal: false,
	}
}

// Append append new errors in 'errs' to 'base'.
// - if 'base' is not of type *Err:
//   - 'base' will first be converted to *Err with 'fatal' default to true
//   - if 'errs' is not empty, the 'fatal' value of the first non-empty error in 'errs' of type *Err will be used
//
// - if 'base' is of type *Err and is not nil:
//   - 'fatal' in 'base' will be used
//   - if error list of 'base' is empty, a nil will be added to the error list of the result.
func Append(base error, errs ...error) (r *Err) {
	if len(errs) > 0 {
		var set, ok bool
		if r, ok = base.(*Err); !ok {
			if base != nil && base.Error() != "" {
				r = &Err{
					fatal:  true, // non *Err instances are considered fatal
					errors: []error{base},
				}
			}
		} else if r != nil { // should check for nil after type assertion
			set = true // 'fatal' value of the result is determined (use the one in 'base')
			if len(r.errors) < 1 {
				// base is of type *Err but the error list is empty
				r.errors = append(r.errors, nil)
			}
		}

		if r == nil { // that is, if 'base' is an empty error
			r = &Err{
				fatal: true, // default to fatal since system errors (non *Err instances) are considered fatal.
			}
		}
		for _, err := range errs {
			if e, ok := err.(*Err); ok {
				if e != nil && len(e.errors) > 0 {
					if !set {
						r.fatal = e.fatal // use the 'fatal' value of the first non-empty *Err in 'errs'
						set = true        // 'fatal' value of the result is determined
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
	if len(r.errors) < 1 {
		r = nil // to avoid the case Append(base, nil) returns a non-nil *Err with no error in it
	}
	return
}

// Appendf append a new error created by fmt.Errorf().
func Appendf(base error, format string, a ...any) *Err {
	return Append(base, fmt.Errorf(format, a...))
}

func Wrap(e error, labels ...string) (r *Err) {
	var ok bool
	if r, ok = e.(*Err); !ok || r == nil {
		if e != nil && e.Error() != "" {
			r = &Err{
				fatal:  true,
				errors: []error{e},
			}
		} else {
			r = &Err{
				fatal: true,
			}
		}
	}
	r.labels = append(r.labels, labels...)
	return
}
