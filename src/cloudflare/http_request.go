package cloudflare

import (
	"github.com/cjtoolkit/cfupdater/src/config"
	"net/http"
	"net/url"
)

const (
	authEmail   = "X-Auth-Email"
	authKey     = "X-Auth-Key"
	contentType = "Content-Type"

	contentTypeValue = "application/json"
)

type httpRequest struct{}

func (_ httpRequest) newHttpRequest(method, rawUrl string) *http.Request {
	parsedUrl, err := url.Parse(rawUrl)
	if nil != err {
		panic(err)
	}

	config := config.GetConfig()

	return &http.Request{
		Method: method,
		URL:    parsedUrl,
		Header: &http.Header{
			authEmail:   {config.Email},
			authKey:     {config.ApiKey},
			contentType: {contentTypeValue},
		},
	}
}

type httpRequestInterface interface {
	newHttpRequest(method, url string) *http.Request
}
