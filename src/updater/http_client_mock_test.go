package updater

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

/*
Implement:
	httpClientInterface
*/
type httpClientMock struct {
	Error      func(name string)
	FnDo       func(req *http.Request) (resp *http.Response, err error)
	FnGet      func(url string) (resp *http.Response, err error)
	FnHead     func(url string) (resp *http.Response, err error)
	FnPost     func(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
	FnPostForm func(url string, data url.Values) (resp *http.Response, err error)
}

func newHttpClientMock(_error func(name string)) *httpClientMock {
	hCM := &httpClientMock{
		Error: _error,
	}
	hCM.reset()
	return hCM
}

func (hCM *httpClientMock) reset() {
	hCM.FnDo = func(req *http.Request) (resp *http.Response, err error) {
		hCM.Error("Do")
		return
	}

	hCM.FnGet = func(url string) (resp *http.Response, err error) {
		hCM.Error("Get")
		return
	}

	hCM.FnHead = func(url string) (resp *http.Response, err error) {
		hCM.Error("Head")
		return
	}

	hCM.FnPost = func(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
		hCM.Error("Post")
		return
	}

	hCM.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
		hCM.Error("PostForm")
		return
	}
}

func (hCM *httpClientMock) Do(req *http.Request) (resp *http.Response, err error) {
	return hCM.FnDo(req)
}

func (hCM *httpClientMock) Get(url string) (resp *http.Response, err error) {
	return hCM.FnGet(url)
}

func (hCM *httpClientMock) Head(url string) (resp *http.Response, err error) {
	return hCM.FnHead(url)
}

func (hCM *httpClientMock) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return hCM.FnPost(url, bodyType, body)
}

func (hCM *httpClientMock) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return hCM.FnPostForm(url, data)
}

/*
Implement:
	ReadCloser in "io"
*/
type readCloser struct {
	io.Reader
}

func newReadCloser(str string) readCloser {
	return readCloser{strings.NewReader(str)}
}

func (rc readCloser) Close() error {
	return nil
}
