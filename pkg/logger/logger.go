package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const fRM_LOGF_MILLI = "2006-01-02T15:04:05.000"

const mASK_ERR uint8 = 1  // 0000 0001
const mASK_EXIT uint8 = 2 // 0000 0010

// _log write log to console
// arguments:
// - ctrl: flow control: x0 - write to stdout; x1 - write to stderr; 1x write to stderr and exit(1)
func _log(
	ctrl uint8, // flow control: 00 - write to stdout; 01 - write to stderr; 1x - exit(1)
	frm string, // format string
	a ...any, // values to plug into the format string
) {
	if ctrl&mASK_ERR > 0 {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%v %v", time.Now().Format(fRM_LOGF_MILLI), frm), a...)
	} else {
		fmt.Printf(fmt.Sprintf("%v %v", time.Now().Format(fRM_LOGF_MILLI), frm), a...)
	}
	if ctrl&mASK_EXIT > 0 {
		os.Exit(1)
	}
}

// Init initialize base loggers.
// - log   : base logger to stdout
// - err   : base logger to stderr
// - fatal : base logger to stderr which exit after writing the message
func Init() (
	log, err, fatal func(string, ...any),
) {
	log = func(frm string, a ...any) {
		_log(0, frm, a...)
	}
	err = func(frm string, a ...any) {
		_log(1, frm, a...)
	}
	fatal = func(frm string, a ...any) {
		_log(3, frm, a...)
	}
	return
}

// AddPrefix build a logger with fixed prefix.
func AddPrefix(
	log func(string, ...any),
	pfx ...string,
) func(string, ...any) {
	var str strings.Builder
	for _, p := range pfx {
		fmt.Fprintf(&str, "[%v]", p)
	}
	return func(frm string, a ...any) {
		log(fmt.Sprintf("%v%v", str.String(), frm), a...)
	}
}

// AddLabels build a logger with fillable labels of specified paddings.
// - val[0, 2, 4...] indicates style: 0 - left justify, 1 - right justify, 2 - pad zero
// - val[1, 3, 5...] indicates the padding value
// NOTE if len of val is an odd number, the last element is ignored
func AddLabels(
	log func(string, ...any),
	val ...int, // (0 - left justify, 1 - right justify, 2 - pad zero, padding value)
) (
	fn func(string, ...any),
	cn int, // number of 'dangling' place holders
) {
	var str strings.Builder
	for i := 0; i < len(val); i += 2 {
		switch val[i] {
		case 0:
			cn++
			fmt.Fprintf(&str, "[%%-%vv]", val[i+1])
		case 1:
			cn++
			fmt.Fprintf(&str, "[%%%vv]", val[i+1])
		case 2:
			cn++
			fmt.Fprintf(&str, "[%%0%vv]", val[i+1])
		}
	}
	fn = func(frm string, a ...any) {
		log(fmt.Sprintf("%v%v", str.String(), frm), a...)
	}
	return
}

// Prefill build a logger obtained from AddLabels() with pre-filled labels.
func Prefill(
	log func(string, ...any),
	args ...any,
) func(string, ...any) {
	return func(frm string, a ...any) {
		log(frm, append(args, a...)...)
	}
}
