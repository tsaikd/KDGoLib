package errutil

import (
	"io"
	"os"
	"strings"
)

// Trace error stack, output to default ErrorFormatter, panic if output error
func Trace(err error) {
	TraceSkip(err, 1)
}

// TraceWrap trace err and wrap with wraperr only if err != nil
func TraceWrap(err error, wraperr error) {
	if err != nil {
		errs := NewErrorsSkip(1, wraperr, err)
		TraceSkip(errs, 1)
	}
}

// TraceSkip error stack, output to default ErrorFormatter, skip n function calls, panic if output error
func TraceSkip(err error, skip int) {
	var errtext string
	var errfmt error
	if errtext, errfmt = defaultTraceFormatter.FormatSkip(err, skip+1); errfmt != nil {
		panic(errfmt)
	}
	if errtext == "" {
		return
	}
	if !strings.HasSuffix(errtext, "\n") {
		errtext += "\n"
	}
	if _, errfmt = defaultTraceOutput.Write([]byte(errtext)); errfmt != nil {
		panic(errfmt)
	}
}

var defaultFormatter = ErrorFormatter(&ConsoleFormatter{
	Seperator: "; ",
})
var defaultTraceFormatter = TraceFormatter(&ConsoleFormatter{
	Seperator:  "; ",
	TimeFormat: "2006-01-02 15:04:05 ",
	LongFile:   true,
	Line:       true,
})
var defaultTraceOutput = io.Writer(os.Stderr)

// SetDefaultFormatter set default ErrorFormatter
func SetDefaultFormatter(formatter ErrorFormatter) {
	defaultFormatter = formatter
}

// SetDefaultTraceFormatter set default ErrorFormatter
func SetDefaultTraceFormatter(formatter TraceFormatter) {
	defaultTraceFormatter = formatter
}

// SetDefaultTraceOutput set default error trace output
func SetDefaultTraceOutput(writer io.Writer) {
	defaultTraceOutput = writer
}
