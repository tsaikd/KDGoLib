package env

import (
	"os"
)

func GetString(key string, defvalue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defvalue
	}
	return value
}
