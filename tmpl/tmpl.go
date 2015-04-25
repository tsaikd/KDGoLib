package tmpl

import (
	"bytes"
	"text/template"

	"github.com/tsaikd/KDGoLib/must"
)

func TemplateOnce(text string, data interface{}) (res string, err error) {
	buffer := bytes.NewBuffer(nil)
	tmpl := template.New("")
	if _, err = tmpl.Parse(text); err != nil {
		return
	}
	if err = tmpl.Execute(buffer, data); err != nil {
		return
	}
	res = buffer.String()
	return
}

func MustTemplateOnce(text string, data interface{}) string {
	return must.MustString(TemplateOnce(text, data))
}
