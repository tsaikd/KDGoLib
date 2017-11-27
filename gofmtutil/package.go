package gofmtutil

import (
	"go/scanner"
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
	"golang.org/x/tools/imports"
)

// errors
var (
	ErrGoImportsFailed  = errutil.NewFactory("goimports failed")
	ErrGoImportsFailed1 = errutil.NewFactory("%v\ngoimports failed")
)

// GoImports return formated go source
func GoImports(data []byte) (fmted []byte, err error) {
	options := imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}
	if fmted, err = imports.Process("", data, &options); err != nil {
		var extraInfos string
		switch errval := err.(type) {
		case scanner.ErrorList:
			lines := strings.Split(string(data), "\n")
			if len(errval) > 0 {
				errelem := errval[0]
				minline := errelem.Pos.Line - 4
				if minline < 0 {
					minline = 0
				}
				maxline := errelem.Pos.Line + 5
				if maxline >= len(lines) {
					maxline = len(lines)
				}
				extraInfos = strings.Join(lines[minline:maxline], "\n")
			}
		}
		if len(extraInfos) > 0 {
			return fmted, ErrGoImportsFailed1.New(err, extraInfos)
		}
		return fmted, ErrGoImportsFailed.New(err)
	}
	return
}
