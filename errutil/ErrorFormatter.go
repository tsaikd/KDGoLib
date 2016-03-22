package errutil

import (
	"bytes"
	"encoding/json"
	"io"
)

// ErrorFormatter to format error
type ErrorFormatter interface {
	Format(error) error
}

// TraceFormatter to trace error occur line formatter
type TraceFormatter interface {
	ErrorFormatter
	FormatSkip(errin error, skip int) error
}

// NewJSONErrorFormatter create JSONErrorFormatter instance
func NewJSONErrorFormatter(writer io.Writer) *JSONErrorFormatter {
	return &JSONErrorFormatter{
		writer: writer,
	}
}

// JSONErrorFormatter used to format error to json
type JSONErrorFormatter struct {
	writer io.Writer
}

// Format error to json
func (t *JSONErrorFormatter) Format(errin error) (err error) {
	return t.FormatSkip(errin, 1)
}

// FormatSkip trace error line and format to json
func (t *JSONErrorFormatter) FormatSkip(errin error, skip int) (err error) {
	errjson := newJSON(skip+1, errin)
	if errjson == nil {
		return
	}

	return json.NewEncoder(t.writer).Encode(errjson)
}

// NewBufferErrorFormatter create BufferErrorFormatter instance
func NewBufferErrorFormatter() *BufferErrorFormatter {
	return &BufferErrorFormatter{
		JSONErrorFormatter: JSONErrorFormatter{
			writer: &bytes.Buffer{},
		},
	}
}

// BufferErrorFormatter buffer error output, used for testing
type BufferErrorFormatter struct {
	JSONErrorFormatter
}

func (t *BufferErrorFormatter) String() string {
	return t.writer.(*bytes.Buffer).String()
}
