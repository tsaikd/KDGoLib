package tmpl

import (
	"bytes"
	"text/template"

	"github.com/tsaikd/KDGoLib/must"
)

func TemplateString(tmpl *template.Template, data interface{}) (res string, err error) {
	buffer := bytes.NewBuffer(nil)
	if err = tmpl.Execute(buffer, data); err != nil {
		return
	}
	res = buffer.String()
	return
}

func TemplateOnce(text string, data interface{}) (res string, err error) {
	tmpl := template.New("")
	if _, err = tmpl.Parse(text); err != nil {
		return
	}
	if res, err = TemplateString(tmpl, data); err != nil {
		return
	}
	return
}

func MustTemplateOnce(text string, data interface{}) string {
	return must.MustString(TemplateOnce(text, data))
}
