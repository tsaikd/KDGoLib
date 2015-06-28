package httputil

import (
	"crypto/tls"
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

func PingIgnoreCertificate(surl string) (err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(surl)
	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return &ErrorPingFailed{surl, resp.StatusCode}
	}
	return
}
