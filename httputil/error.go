package httputil

import (
	"fmt"
)

type ErrorPingFailed struct {
	Url        string
	StatusCode int
}

func (t *ErrorPingFailed) Error() string {
	return fmt.Sprintf("ping return unexpect status code: %v", t.StatusCode)
}
