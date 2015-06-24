package errorJson

import (
	"net/http"
	"reflect"

	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/tsaikd/KDGoLib/errutil"
)

type ResponseError struct {
	Status int `json:"status,omitempty"`
	*errutil.ErrorSlice
}

func newResponseError(status int, err error, errs ...error) (reserr ResponseError) {
	errs = append([]error{err}, errs...)
	reserr.Status = status
	reserr.ErrorSlice = errutil.NewErrorSlice(errs...)
	return
}

func RenderErrorJSON(render render.Render, status int, err error, errs ...error) {
	reserr := newResponseError(status, err, errs...)
	render.JSON(status, reserr)
}

func ReturnErrorProvider() martini.ReturnHandler {
	return func(ctx martini.Context, vals []reflect.Value) {
		rv := ctx.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
		res := rv.Interface().(http.ResponseWriter)
		var responseVal reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			status := vals[0].Int()
			if status > 0 {
				res.WriteHeader(int(status))
			}
			responseVal = vals[1]
		} else if len(vals) > 0 {
			responseVal = vals[0]
		}

		if responseVal.Interface() == nil {
			return
		}
		if isError(responseVal) {
			r := ctx.Get(inject.InterfaceOf((*render.Render)(nil))).Interface().(render.Render)
			err := responseVal.Interface().(error)
			RenderErrorJSON(r, 404, err)
			return
		}

		if canDeref(responseVal) {
			responseVal = responseVal.Elem()
		}
		if isByteSlice(responseVal) {
			res.Write(responseVal.Bytes())
		} else {
			res.Write([]byte(responseVal.String()))
		}
	}
}

func isError(val reflect.Value) bool {
	_, ok := val.Interface().(error)
	return ok
}

func isByteSlice(val reflect.Value) bool {
	return val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8
}

func canDeref(val reflect.Value) bool {
	return val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr
}
