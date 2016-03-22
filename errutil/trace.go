package errutil

import "os"

var defaultFormatter ErrorFormatter

// Trace error stack, output to default ErrorFormatter, panic if output error
func Trace(err error) {
	TraceSkip(err, 1)
}

// TraceSkip error stack, output to default ErrorFormatter, skip n function calls, panic if output error
func TraceSkip(err error, skip int) {
	if defaultFormatter == nil {
		defaultFormatter = NewJSONErrorFormatter(os.Stderr)
	}

	if traceFormatter, ok := defaultFormatter.(TraceFormatter); ok {
		errfmt := traceFormatter.FormatSkip(err, skip+1)
		if errfmt != nil {
			panic(errfmt)
		}
	} else {
		errfmt := defaultFormatter.Format(err)
		if errfmt != nil {
			panic(errfmt)
		}
	}
}

// SetDefaultFormatter set default ErrorFormatter
func SetDefaultFormatter(formatter ErrorFormatter) {
	defaultFormatter = formatter
}
