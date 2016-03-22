package errorJson

import (
	"reflect"

	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/tsaikd/KDGoLib/errutil"
)

// ReturnError define api function call return error struct
type ReturnError struct {
	Status int
	Body   []byte
	Error  error
}

// AddError add error to ReturnError
func (t *ReturnError) AddError(err error) {
	if err == nil {
		return
	}
	t.Error = errutil.NewErrorsSkip(1, err, t.Error)
	if t.Status == 0 {
		t.Status = 404
	}
}

type responseError struct {
	Status int `json:"status,omitempty"`
	*errutil.ErrorJSON
}

// RenderErrorJSON render error in json format
func RenderErrorJSON(render render.Render, status int, err error, errs ...error) {
	errConcat := append([]error{err}, errs...)
	reserr := responseError{
		Status:    status,
		ErrorJSON: errutil.NewJSON(errutil.NewErrorsSkip(1, errConcat...)),
	}
	render.JSON(status, reserr)
}

// BindMartini bind return error handler to martini instance
func BindMartini(m *martini.Martini) {
	if rv := m.Get(inject.InterfaceOf((*render.Render)(nil))); !rv.IsValid() {
		m.Use(render.Renderer())
	}

	m.Map(returnErrorHandler())
	m.Use(returnErrorProvider())

	return
}

func returnErrorProvider() martini.Handler {
	return func(ctx martini.Context) {
		returnError := &ReturnError{}
		ctx.Map(returnError)
	}
}

func returnErrorHandler() martini.ReturnHandler {
	return func(ctx martini.Context, vals []reflect.Value) {
		if ctx.Written() {
			return
		}

		returnError := ctx.Get(reflect.TypeOf(&ReturnError{})).Interface().(*ReturnError)
		render := ctx.Get(inject.InterfaceOf((*render.Render)(nil))).Interface().(render.Render)

		var responseVal reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			returnError.Status = int(vals[0].Int())
			responseVal = vals[1]
		} else if len(vals) > 0 {
			responseVal = vals[0]
		}

		if responseVal.Interface() == nil {
			return
		}
		if isError(responseVal) {
			returnError.AddError(responseVal.Interface().(error))
		}

		if returnError.Error != nil {
			RenderErrorJSON(render, returnError.Status, returnError.Error)
			return
		}

		if canDeref(responseVal) {
			responseVal = responseVal.Elem()
		}
		if isByteSlice(responseVal) {
			returnError.Body = responseVal.Bytes()
		} else {
			returnError.Body = []byte(responseVal.String())
		}
		render.Data(returnError.Status, returnError.Body)
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
