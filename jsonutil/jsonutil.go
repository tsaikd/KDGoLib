package jsonutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func UnmarshalHttpResponse(resp *http.Response, v interface{}) (data []byte, err error) {
	defer resp.Body.Close()
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	err = json.Unmarshal(data, v)
	return
}
