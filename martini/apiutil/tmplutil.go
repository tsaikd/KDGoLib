package apiutil

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
)

var (
	regexpAPIVersionPrefix    = regexp.MustCompile(`^/(\d+)`)
	regexpAPISlashCamelCase   = regexp.MustCompile(`/\w`)
	regexpAPIParam            = regexp.MustCompile(`:[^/#?()\.\\]+`) // reference martini/router.go#210
	regexpAPIParam2           = regexp.MustCompile(`\(\?P\<([^/#?()\.\\]+)\>.*?\)`)
	regexpStartOptSpace       = regexp.MustCompile(`(?m)^\s*`)
	regexpStartSpaceEmptyLine = regexp.MustCompile(`(?m)^\s+$`)
	regexpContinuousNewLine   = regexp.MustCompile(`(?m)^\n\n+$`)
)

// RequestInfo contains request function name, params
type RequestInfo struct {
	APIName      string
	FuncName     string // FuncName := APIName with version suffix
	FuncArgs     RequestParams
	Path         string
	Params       RequestParams
	ExtraStructs ExtraStructs
}

// AnalyzeRequestStruct analyze request struct object, produce request struct map
// struct field with tag `apijs:"inherit"` will extend the field arguments
// repl used for custom function name generator
func AnalyzeRequestStruct(
	pattern string,
	requestStruct interface{},
	replFuncName func(param string) string,
	replPath func(param string) string,
) (reqinfo RequestInfo) {
	// analyze pattern for function name
	if pattern != "" {
		reqinfo.APIName = GetAPINameByPattern(pattern, nil)

		// fill FuncName and add param to FuncArgs, Params
		reqinfo.FuncName = GetFuncNameByPattern(pattern, func(param string) string {
			reqinfo.FuncArgs.Upsert(param, nil)
			reqinfo.Params.Upsert(param, nil)
			if replFuncName != nil {
				return replFuncName(param)
			}
			return ""
		})
	}

	// analyze request struct
	if requestStruct != nil {
		AnalyzeStructParams(requestStruct, &reqinfo.FuncArgs, &reqinfo.Params, &reqinfo.ExtraStructs)
	}

	// analyze pattern for path
	if pattern != "" {
		// fill Path and update param to FuncArgs, Params
		reqinfo.Path = GetPathByPattern(pattern, func(param string) string {
			result := param
			if replPath != nil {
				result = replPath(param)
			}
			if result == "" {
				reqinfo.FuncArgs.Upsert(param, nil)
				reqinfo.Params.Upsert(param, nil)
			} else {
				reqinfo.Params.Delete(param)
			}
			return result
		})
	}

	return
}

// GetAPINameByPattern return API name by analyze pattern
func GetAPINameByPattern(
	pattern string,
	repl func(string) string,
) string {
	result := pattern
	result = regexpAPIVersionPrefix.ReplaceAllString(result, "")

	result = regexpAPIParam.ReplaceAllStringFunc(result, func(str string) string {
		param := str[1:]
		return "(?P<" + param + ">[^/]+)"
	})

	result = regexpAPIParam2.ReplaceAllStringFunc(result, func(str string) string {
		matches := regexpAPIParam2.FindStringSubmatch(str)
		param := matches[1]
		if repl == nil {
			return ""
		}
		return repl(param)
	})

	result = regexpAPISlashCamelCase.ReplaceAllStringFunc(result, func(str string) string {
		return strings.ToUpper(str[1:])
	})

	// remove all slash
	result = strings.Replace(result, "/", "", -1)

	return result
}

// GetFuncNameByPattern return API name by analyze pattern
func GetFuncNameByPattern(
	pattern string,
	repl func(string) string,
) string {
	result := GetAPINameByPattern(pattern, repl)
	versionInfos := regexpAPIVersionPrefix.FindStringSubmatch(pattern)
	if len(versionInfos) > 1 {
		result += "V" + versionInfos[1]
	}
	return result
}

// AnalyzeStructParams analyze requestStruct fields and fill funcArgs/params
func AnalyzeStructParams(
	requestStruct interface{},
	funcArgs *RequestParams,
	params *RequestParams,
	extraStructs *ExtraStructs,
) {
	reflectStruct, ok := requestStruct.(reflect.Type)
	if !ok {
		reflectStruct = reflect.TypeOf(requestStruct)
	}
	for i := 0; i < reflectStruct.NumField(); i++ {
		field := reflectStruct.Field(i)

		if fieldTag := field.Tag.Get("apijs"); fieldTag == "inherit" {
			AnalyzeStructParams(field.Type, funcArgs, params, extraStructs)
			continue
		}

		if fieldTag := field.Tag.Get("json"); fieldTag != "" {
			if tags := strings.Split(fieldTag, ","); len(tags) > 0 {
				jsonName := tags[0]
				if jsonName == "-" {
					continue
				}

				funcArgs.Upsert(jsonName, &field)
				params.Upsert(jsonName, &field)

				structType := ensureReflectStruct(field.Type)
				switch structType.Kind() {
				case reflect.Struct:
					extraStruct := extraStructs.Upsert(jsonName, field.Type)
					AnalyzeStructParams(
						structType,
						&extraStruct.Params,
						nil,
						extraStructs,
					)
				}
			}
		}
	}
}

func ensureReflectStruct(reflectType reflect.Type) reflect.Type {
	switch reflectType.Kind() {
	case reflect.Slice:
		return reflectType.Elem()
	}
	return reflectType
}

// GetPathByPattern return api path by analyze pattern
func GetPathByPattern(
	pattern string,
	repl func(string) string,
) string {
	result := pattern

	result = regexpAPIParam.ReplaceAllStringFunc(result, func(str string) string {
		param := str[1:]
		return "(?P<" + param + ">[^/]+)"
	})

	result = regexpAPIParam2.ReplaceAllStringFunc(result, func(str string) string {
		matches := regexpAPIParam2.FindStringSubmatch(str)
		param := matches[1]
		if repl == nil {
			return param
		}
		return repl(param)
	})

	return result
}

// TailComma return "," if idx is the last of params
func TailComma(idx int, params RequestParams) string {
	if idx < len(params)-1 {
		return ","
	}
	return ""
}

// TrimReplaceStartSpace trim space and replace start space with repl per line
func TrimReplaceStartSpace(src string, repl string) string {
	result := strings.TrimSpace(src)
	if result == "" {
		return ""
	}
	return regexpStartOptSpace.ReplaceAllString(result, repl)
}

// TrimStartSpaceEmptyLine remove empty line with start space
func TrimStartSpaceEmptyLine(src string) (result string, err error) {
	lines := strings.Split(src, "\n")
	buffer := bytes.Buffer{}
	for _, line := range lines {
		if !regexpStartSpaceEmptyLine.MatchString(line) {
			if _, err = buffer.WriteString(line + "\n"); err != nil {
				return
			}
		}
	}
	return buffer.String(), nil
}

// ReplaceContinuousNewLine replace continuous newline
func ReplaceContinuousNewLine(src string, repl string) string {
	return regexpContinuousNewLine.ReplaceAllString(src, repl)
}
