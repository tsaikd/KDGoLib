package webutil

import (
	"crypto/tls"
	"net/http"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorPing = errutil.NewFactory("ping %q return unexpect status code %q")
)

// PingIgnoreCertificate ping url but ignore https certification check
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
		return ErrorPing.New(nil, surl, resp.StatusCode)
	}
	return
}
