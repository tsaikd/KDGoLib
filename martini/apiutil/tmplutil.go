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
	FuncName string
	FuncArgs RequestParams
	Path     string
	Params   RequestParams
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
		reflectStruct, ok := requestStruct.(reflect.Type)
		if !ok {
			reflectStruct = reflect.TypeOf(requestStruct)
		}
		for i := 0; i < reflectStruct.NumField(); i++ {
			field := reflectStruct.Field(i)

			if fieldTag := field.Tag.Get("apijs"); fieldTag == "inherit" {
				extends := AnalyzeRequestStruct("", field.Type, nil, nil)
				for _, paramType := range extends.Params {
					reqinfo.FuncArgs.Upsert(paramType.FieldName, &paramType.FieldType)
					reqinfo.Params.Upsert(paramType.FieldName, &paramType.FieldType)
				}
				continue
			}

			if fieldTag := field.Tag.Get("json"); fieldTag != "" {
				if tags := strings.Split(fieldTag, ","); len(tags) > 0 {
					jsonName := tags[0]
					if jsonName != "-" {
						reqinfo.FuncArgs.Upsert(jsonName, &field)
						reqinfo.Params.Upsert(jsonName, &field)
					}
				}
			}
		}
	}

	// analyze pattern for path
	if pattern != "" {
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

// GetFuncNameByPattern return function name by analyze pattern, check param in requestMap
func GetFuncNameByPattern(
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
	versionInfos := regexpAPIVersionPrefix.FindStringSubmatch(pattern)
	if len(versionInfos) > 1 {
		result += "V" + versionInfos[1]
	}

	return result
}

// GetPathByPattern return api path by analyze pattern, consume param in requestMap
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
