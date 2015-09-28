package bindstruct

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-martini/martini"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/injectutil"
	"github.com/tsaikd/KDGoLib/structutil/reflectstruct"
	"github.com/tsaikd/KDGoLib/structutil/trimstructspace"
)

type BindConfig struct {
	DisableBindQuery   bool
	DisableBindParams  bool
	DisableTrimSpace   bool
	DisableGoValidator bool

	// Maximum amount of memory to use when parsing a multipart form.
	// Set this to whatever value you prefer; default is 20 MB.
	MaxMemory int64
}

func NewBindConfig() BindConfig {
	return BindConfig{
		MaxMemory: 20 * 1024 * 1024,
	}
}

var (
	DefaultBindConfig = NewBindConfig()
)

func BindStruct(obj interface{}) martini.Handler {
	return BindStructWithConfig(obj, DefaultBindConfig)
}

func BindStructWithConfig(obj interface{}, config BindConfig) martini.Handler {
	return func(context martini.Context, req *http.Request) (err error) {
		bindStruct := reflect.New(reflect.TypeOf(obj))

		if !config.DisableBindQuery {
			if _, err = injectutil.Invoke(context, bindQuery(bindStruct)); err != nil {
				return
			}
		}

		if !config.DisableBindParams {
			if _, err = injectutil.Invoke(context, bindParams(bindStruct)); err != nil {
				return
			}
		}

		contentType := req.Header.Get("Content-Type")
		if contentType != "" {
			if strings.Contains(contentType, "json") {
				if _, err = injectutil.Invoke(context, bindJsonBody(bindStruct)); err != nil {
					return
				}
			} else if strings.Contains(contentType, "multipart/form-data") {
				if _, err = injectutil.Invoke(context, bindMultipartForm(bindStruct, config)); err != nil {
					return
				}
			}
		}

		if !config.DisableTrimSpace {
			if err = trimstructspace.TrimStructSpace(bindStruct.Interface()); err != nil {
				return
			}
		}

		if !config.DisableGoValidator {
			ok, err := govalidator.ValidateStruct(bindStruct.Interface())
			if err != nil {
				return err
			}
			if !ok {
				return errutil.New("invalid parameters")
			}
		}

		context.Map(bindStruct.Elem().Interface())
		return
	}
}

func bindQuery(bindStruct reflect.Value) martini.Handler {
	return func(req *http.Request) (err error) {
		return reflectstruct.ReflectStruct(bindStruct.Interface(), req.URL.Query())
	}
}

func bindParams(bindStruct reflect.Value) martini.Handler {
	return func(params martini.Params) (err error) {
		return reflectstruct.ReflectStruct(bindStruct.Interface(), params)
	}
}

func bindJsonBody(bindStruct reflect.Value) martini.Handler {
	return func(context martini.Context, req *http.Request) (err error) {
		if req.Body != nil {
			defer req.Body.Close()
			err = json.NewDecoder(req.Body).Decode(bindStruct.Interface())
			if err != nil && err != io.EOF {
				return
			}
		}
		return
	}
}

func bindMultipartForm(bindStruct reflect.Value, config BindConfig) martini.Handler {
	return func(context martini.Context, req *http.Request) (err error) {
		if err = req.ParseMultipartForm(config.MaxMemory); err != nil {
			return
		}
		if err = reflectstruct.ReflectStruct(bindStruct.Interface(), req.MultipartForm.Value); err != nil {
			return
		}
		if err = reflectstruct.ReflectStruct(bindStruct.Interface(), req.MultipartForm.File); err != nil {
			return
		}
		return
	}
}
