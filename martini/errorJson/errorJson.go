package errorJson

import (
	"github.com/martini-contrib/render"
)

type ResponseError struct {
	Status int      `json:"status,omitempty"`
	Error  string   `json:"error,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

func RenderErrorJSON(render render.Render, status int, err error, errs ...error) {
	var (
		reserr ResponseError
	)
	reserr.Status = status
	reserr.Error = err.Error()
	reserr.Errors = append(reserr.Errors, err.Error())
	for _, extraErr := range errs {
		reserr.Errors = append(reserr.Errors, extraErr.Error())
	}
	render.JSON(status, reserr)
}
