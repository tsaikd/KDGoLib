package webutil

import (
	"io/ioutil"
	"net/http"

	"github.com/tsaikd/KDGoLib/errutil"
)

// ReadResponse read http response body and close resource
func ReadResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// NewResponseError consume response body and create error object if status code is abnormal
func NewResponseError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	body, err := ReadResponse(resp)
	if err != nil {
		return err
	}

	return errutil.New(body)
}
