package env

import (
	"os"
	"strconv"
)

// like os.Getenv and provide default value
func GetString(key string, defvalue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defvalue
	}
	return value
}

func GetBool(key string, defvalue bool) (value bool) {
	envValue := os.Getenv(key)
	switch envValue {
	case "true", "True", "TRUE", "yes", "Yes", "YES", "1":
		value = true
	case "false", "False", "FALSE", "no", "No", "NO", "0":
		value = false
	default:
		value = defvalue
	}
	return value
}

func GetInt(key string, defvalue int) (value int) {
	var (
		err error
	)
	envValue := os.Getenv(key)
	if value, err = strconv.Atoi(envValue); err != nil {
		value = defvalue
	}
	return
}
