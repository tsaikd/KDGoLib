package bindStruct

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/gin-gonic-util/apiutil"
	"github.com/tsaikd/KDGoLib/structutil/reflectstruct"
	"github.com/tsaikd/KDGoLib/structutil/trimstructspace"
	"github.com/tsaikd/govalidator"
)

// errors
var (
	ErrorInjectorNotFound = errutil.NewFactory("injector not found, bindStruct require gin-injector middleware")
)

// Config for binding
type Config struct {
	DisableBindQuery   bool
	DisableBindParams  bool
	DisableTrimSpace   bool
	DisableGoValidator bool

	// Maximum amount of memory to use when parsing a multipart form.
	// Set this to whatever value you prefer; default is 20 MB.
	MaxMemory int64
}

// NewConfig create default BindConfig
func NewConfig() Config {
	return Config{
		MaxMemory: 20 * 1024 * 1024,
	}
}

var (
	// DefaultConfig default Config for binding
	DefaultConfig = NewConfig()
)

// BindStruct martini handler for bind struct with default config
func BindStruct(obj interface{}) gin.HandlerFunc {
	return BindStructWithConfig(obj, DefaultConfig)
}

// BindStructWithConfig martini handler for bind struct
func BindStructWithConfig(obj interface{}, config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		inj, exist := apiutil.Get(c)
		if !exist {
			c.AbortWithError(http.StatusServiceUnavailable, ErrorInjectorNotFound.New(nil))
			return
		}

		bindStruct := reflect.New(reflect.TypeOf(obj))
		defer func() {
			inj.Map(bindStruct.Elem().Interface())
		}()

		req := c.Request

		if !config.DisableBindQuery {
			if err := bindQuery(bindStruct, req); err != nil {
				c.Error(err)
				return
			}
		}

		if !config.DisableBindParams {
			if err := bindParams(bindStruct, c.Params); err != nil {
				c.Error(err)
				return
			}
		}

		if err := c.Bind(bindStruct.Interface()); err != nil {
			c.Error(err)
			return
		}

		// contentType := req.Header.Get("Content-Type")
		// if contentType != "" {
		// 	if strings.Contains(contentType, "json") {
		// 		if err = bindJSONBody(bindStruct, req); err != nil {
		// 			return
		// 		}
		// 	} else if strings.Contains(contentType, "multipart/form-data") {
		// 		if err = bindMultipartForm(bindStruct, config, req); err != nil {
		// 			return
		// 		}
		// 	}
		// }
		//
		if !config.DisableTrimSpace {
			if err := trimstructspace.TrimStructSpace(bindStruct.Interface()); err != nil {
				c.Error(err)
				return
			}
		}

		if !config.DisableGoValidator {
			ok, err := govalidator.ValidateStruct(bindStruct.Interface())
			if err != nil {
				c.Error(err)
				return
			}
			if !ok {
				c.Error(errutil.New("invalid parameters"))
				return
			}
		}

		return
	}
}

// EnsureBindStruct ensure martini context mapped struct
// func EnsureBindStruct(ctx martini.Context, obj interface{}) reflect.Value {
// 	objType := ctx.Get(reflect.TypeOf(obj))
// 	if objType.IsValid() {
// 		return objType
// 	}
//
// 	ctx.Invoke(BindStruct(obj))
// 	objType = ctx.Get(reflect.TypeOf(obj))
// 	return objType
// }

func bindQuery(bindStruct reflect.Value, req *http.Request) (err error) {
	return reflectstruct.ReflectStruct(bindStruct.Interface(), req.URL.Query())
}

func bindParams(bindStruct reflect.Value, params gin.Params) (err error) {
	paramMap := map[string]interface{}{}
	for _, param := range params {
		paramMap[param.Key] = param.Value
	}
	return reflectstruct.ReflectStruct(bindStruct.Interface(), paramMap)
}

func bindJSONBody(bindStruct reflect.Value, req *http.Request) (err error) {
	if req.Body != nil {
		defer req.Body.Close()
		err = json.NewDecoder(req.Body).Decode(bindStruct.Interface())
		if err != nil && err != io.EOF {
			return
		}
	}
	return
}

func bindMultipartForm(bindStruct reflect.Value, config Config, req *http.Request) (err error) {
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
