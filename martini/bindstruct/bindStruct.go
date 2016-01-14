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
	"github.com/tsaikd/KDGoLib/logutil"
	"github.com/tsaikd/KDGoLib/structutil/reflectstruct"
	"github.com/tsaikd/KDGoLib/structutil/trimstructspace"
)

// BindConfig config for bind
type BindConfig struct {
	DisableBindQuery   bool
	DisableBindParams  bool
	DisableTrimSpace   bool
	DisableGoValidator bool

	// Maximum amount of memory to use when parsing a multipart form.
	// Set this to whatever value you prefer; default is 20 MB.
	MaxMemory int64
}

// NewBindConfig create default BindConfig
func NewBindConfig() BindConfig {
	return BindConfig{
		MaxMemory: 20 * 1024 * 1024,
	}
}

var (
	// DefaultBindConfig default BindConfig
	DefaultBindConfig = NewBindConfig()
)

// BindStruct martini handler for bind struct with default config
func BindStruct(obj interface{}) martini.Handler {
	return BindStructWithConfig(obj, DefaultBindConfig)
}

// BindStructWithConfig martini handler for bind struct
func BindStructWithConfig(obj interface{}, config BindConfig) martini.Handler {
	return func(context martini.Context, req *http.Request, params martini.Params) (err error) {
		// handle return error because martini will ignore
		defer func() {
			if err != nil {
				logutil.GetLevelLogger(context).Errorln(err)
			}
		}()

		bindStruct := reflect.New(reflect.TypeOf(obj))
		defer func() {
			context.Map(bindStruct.Elem().Interface())
		}()

		if !config.DisableBindQuery {
			if err = bindQuery(bindStruct, req); err != nil {
				return
			}
		}

		if !config.DisableBindParams {
			if err = bindParams(bindStruct, params); err != nil {
				return
			}
		}

		contentType := req.Header.Get("Content-Type")
		if contentType != "" {
			if strings.Contains(contentType, "json") {
				if err = bindJSONBody(bindStruct, req); err != nil {
					return
				}
			} else if strings.Contains(contentType, "multipart/form-data") {
				if err = bindMultipartForm(bindStruct, config, req); err != nil {
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

		return
	}
}

// EnsureBindStruct ensure martini context mapped struct
func EnsureBindStruct(ctx martini.Context, obj interface{}) reflect.Value {
	objType := ctx.Get(reflect.TypeOf(obj))
	if objType.IsValid() {
		return objType
	}

	ctx.Invoke(BindStruct(obj))
	objType = ctx.Get(reflect.TypeOf(obj))
	return objType
}

func bindQuery(bindStruct reflect.Value, req *http.Request) (err error) {
	return reflectstruct.ReflectStruct(bindStruct.Interface(), req.URL.Query())
}

func bindParams(bindStruct reflect.Value, params martini.Params) (err error) {
	return reflectstruct.ReflectStruct(bindStruct.Interface(), params)
}

func bindJSONBody(bindStruct reflect.Value, req *http.Request) (err error) {
	if req.Body != nil {
		defer req.Body.Close()
		var value interface{}
		err = json.NewDecoder(req.Body).Decode(&value)
		if err != nil && err != io.EOF {
			return
		}
		return reflectstruct.ReflectStruct(bindStruct.Interface(), value)
	}
	return
}

func bindMultipartForm(bindStruct reflect.Value, config BindConfig, req *http.Request) (err error) {
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
