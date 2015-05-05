package must

import (
	"text/template"
)

func MustString(text string, err error) string {
	if err != nil {
		panic(err)
	}
	return text
}

func MustTemplate(tmpl *template.Template, err error) *template.Template {
	if err != nil {
		panic(err)
	}
	return tmpl
}
