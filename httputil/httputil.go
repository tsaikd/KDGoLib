package httputil

import (
	"io/ioutil"
	"net/http"
)

func ReadResponse(resp *http.Response) (body string, err error) {
	var (
		data []byte
	)

	defer resp.Body.Close()
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	body = string(data)
	return
}
