package errutil

import "bytes"

// NewConsoleFormatter create JSONErrorFormatter instance
func NewConsoleFormatter(seperator string) *ConsoleFormatter {
	return &ConsoleFormatter{
		seperator: seperator,
	}
}

// ConsoleFormatter used to format error object in console readable
type ConsoleFormatter struct {
	seperator string
}

// Format error object
func (t *ConsoleFormatter) Format(errin error) (errtext string, err error) {
	return t.FormatSkip(errin, 1)
}

// FormatSkip trace error line and format object
func (t *ConsoleFormatter) FormatSkip(errin error, skip int) (errtext string, err error) {
	if t.seperator == "" {
		return getErrorText(errin), nil
	}

	errobj := castErrorObject(nil, skip+1, errin)
	if errobj == nil {
		return "", nil
	}

	buffer := &bytes.Buffer{}

	if walkerr := WalkErrors(errobj, func(errloop ErrorObject) (stop bool, walkerr error) {
		if buffer.Len() > 0 {
			if _, errio := buffer.WriteString(t.seperator); errio != nil {
				return true, errio
			}
		}
		if _, errio := buffer.WriteString(getErrorText(errloop)); errio != nil {
			return true, errio
		}
		return false, nil
	}); walkerr != nil {
		return buffer.String(), walkerr
	}

	return buffer.String(), nil
}
