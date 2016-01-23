package updater

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/url"
	"testing"
)

func TestClient(t *testing.T) {
	// start let

	shouldNotBeCalledFn := func(name string) {
		t.Errorf("'%s' should not be called!", name)
	}

	logger := newLoggerMock()
	httpClient := newHttpClientMock(shouldNotBeCalledFn)

	_client := &client{
		httpClient: httpClient,
		logger:     logger,
		url:        API_URL,
		data: Data{
			Name: "test",
		},
	}

	reset := func() {
		logger.Buf.Reset()
		httpClient.reset()
	}

	// end let

	Convey("'getUrlAndType' should return ipv4 url and type on 'A' record", t, func() {
		url, _type := _client.getUrlAndType(&Object{Type: "A"})
		So(url, ShouldEqual, URL_IPV4)
		So(_type, ShouldEqual, "ipv4")
	})

	reset()

	Convey("'getUrlAndType' should return ipv6 url and type on 'AAAA' record", t, func() {
		url, _type := _client.getUrlAndType(&Object{Type: "AAAA"})
		So(url, ShouldEqual, URL_IPV6)
		So(_type, ShouldEqual, "ipv6")
	})

	reset()

	Convey("In 'getObjects' we got a timeout from 'httpClient'", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, "rec_load_all")
			err = fmt.Errorf("I have timedout, sorry")
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4, ShouldBeNil)
		So(ipv6, ShouldBeNil)

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("API Timeout"))
	})

	reset()

	Convey("In 'getObjects' we got a response from 'httpClient' but json decoder failed", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, "rec_load_all")
			resp = &http.Response{Body: newReadCloser(shoddyJson)}
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4, ShouldBeNil)
		So(ipv6, ShouldBeNil)

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("JSON Decoder Failed"))
	})

	reset()

	Convey("In 'getObjects' we got a response from 'httpClient', json decoder went well, but the result was a fail", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, "rec_load_all")
			resp = &http.Response{Body: newReadCloser(recLoadAllFailer)}
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4, ShouldBeNil)
		So(ipv6, ShouldBeNil)

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("API Request Failed:", "I am error"))
	})

	reset()

	Convey("In 'getObjects' we got a response from 'httpClient', json decoder went well, the result were great, we manage to obtain both 'A' and 'AAAA' records", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, "rec_load_all")
			resp = &http.Response{Body: newReadCloser(recLoadAllSuccess)}
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4.Type, ShouldEqual, "A")
		So(ipv6.Type, ShouldEqual, "AAAA")
	})

	reset()

	Convey("In 'getObjects' we got a response from 'httpClient', json decoder went well, the result were great, we manage to obtain both 'A' record", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, "rec_load_all")
			resp = &http.Response{Body: newReadCloser(recLoadAllSuccessIpV4)}
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4.Type, ShouldEqual, "A")
		So(ipv6, ShouldBeNil)
	})

	reset()

	Convey("In 'getObjects' we got a response from 'httpClient', json decoder went well, the result were great, we manage to obtain both 'AAAA' record", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, "rec_load_all")
			resp = &http.Response{Body: newReadCloser(recLoadAllSuccessIpV6)}
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4, ShouldBeNil)
		So(ipv6.Type, ShouldEqual, "AAAA")
	})

	reset()
}
