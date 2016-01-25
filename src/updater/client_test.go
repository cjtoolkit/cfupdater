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

	Convey("In 'getObjects' we got a timeout (an error) from 'httpClient'", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, REC_LOAD_ALL)
			err = fmt.Errorf("I have timedout, sorry")
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4, ShouldBeNil)
		So(ipv6, ShouldBeNil)

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("Http Client Error:", "I have timedout, sorry"))
	})

	reset()

	Convey("In 'getObjects' we got a response from 'httpClient' but json decoder failed", t, func() {
		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, REC_LOAD_ALL)
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
			So(data.Get("a"), ShouldEqual, REC_LOAD_ALL)
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
			So(data.Get("a"), ShouldEqual, REC_LOAD_ALL)
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
			So(data.Get("a"), ShouldEqual, REC_LOAD_ALL)
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
			So(data.Get("a"), ShouldEqual, REC_LOAD_ALL)
			resp = &http.Response{Body: newReadCloser(recLoadAllSuccessIpV6)}
			return
		}

		ipv4, ipv6 := _client.getObjects()
		So(ipv4, ShouldBeNil)
		So(ipv6.Type, ShouldEqual, "AAAA")
	})

	reset()

	Convey("Run updater on 'A' record, but timeout while requesting current IP address", t, func() {
		ipv4 := &Object{
			Content: "127.0.0.1",
			Type:    "A",
		}

		ipurl, iptype := _client.getUrlAndType(ipv4)

		address := ipv4.Content

		httpClient.FnGet = func(url string) (resp *http.Response, err error) {
			So(ipurl, ShouldEqual, URL_IPV4)
			err = fmt.Errorf("I have timed out")
			return
		}

		(*_client).runOn(ipv4, ipurl, iptype, &address)

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln(ipurl, ": Http Error:", "I have timed out"))
	})

	reset()

	Convey("Run updater on 'A' record, manage to get up-to-date IP address, but it's not a valid IP address", t, func() {
		ipv4 := &Object{
			Content: "127.0.0.1",
			Type:    "A",
		}

		ipurl, iptype := _client.getUrlAndType(ipv4)

		address := ipv4.Content

		httpClient.FnGet = func(url string) (resp *http.Response, err error) {
			So(ipurl, ShouldEqual, URL_IPV4)
			resp = &http.Response{Body: newReadCloser("192.168.1_0")}
			return
		}

		(*_client).runOn(ipv4, ipurl, iptype, &address)

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("192.168.1_0", "is not a valid IP Address"))
	})

	reset()

	Convey("Run updater on 'A' record, manage to get up-to-date IP address, but the current IP address is up-to-date, therefore don't need to do anything else", t, func() {
		ipv4 := &Object{
			Content: "127.0.0.1",
			Type:    "A",
		}

		ipurl, iptype := _client.getUrlAndType(ipv4)

		address := ipv4.Content

		httpClient.FnGet = func(url string) (resp *http.Response, err error) {
			So(ipurl, ShouldEqual, URL_IPV4)
			resp = &http.Response{Body: newReadCloser(" 127.0.0.1 ")}
			return
		}

		(*_client).runOn(ipv4, ipurl, iptype, &address)
	})

	reset()

	Convey("Run updater on 'A' record, manage to get up-to-date IP address, the current IP address is out-of-date, but timedout while trying to update it.", t, func() {
		ipv4 := &Object{
			Content: "127.0.0.1",
			Type:    "A",
		}

		ipurl, iptype := _client.getUrlAndType(ipv4)

		address := ipv4.Content

		httpClient.FnGet = func(url string) (resp *http.Response, err error) {
			So(ipurl, ShouldEqual, URL_IPV4)
			resp = &http.Response{Body: newReadCloser(" 192.168.1.1 ")}
			return
		}

		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, REC_EDIT)
			err = fmt.Errorf("Sorry, I have timed out :(")
			return
		}

		(*_client).runOn(ipv4, ipurl, iptype, &address)

		So(address, ShouldEqual, "127.0.0.1")
		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln(API_URL, iptype, ": Http Error:", "Sorry, I have timed out :("))
	})

	reset()

	Convey("Run updater on 'A' record, manage to get up-to-date IP address, the current IP address is out-of-date, got a response while updating IP, but JSON Decoder failed", t, func() {
		ipv4 := &Object{
			Content: "127.0.0.1",
			Type:    "A",
		}

		ipurl, iptype := _client.getUrlAndType(ipv4)

		address := ipv4.Content

		httpClient.FnGet = func(url string) (resp *http.Response, err error) {
			So(ipurl, ShouldEqual, URL_IPV4)
			resp = &http.Response{Body: newReadCloser(" 192.168.1.1 ")}
			return
		}

		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, REC_EDIT)
			resp = &http.Response{Body: newReadCloser(shoddyJson)}
			return
		}

		(*_client).runOn(ipv4, ipurl, iptype, &address)

		So(address, ShouldEqual, "127.0.0.1")
		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("CfUpdater json decoder failed (", iptype, ")"))
	})

	reset()

	Convey("Run updater on 'A' record, manage to get up-to-date IP address, the current IP address is out-of-date, got a response while updating IP, but the result was unsuccessful", t, func() {
		ipv4 := &Object{
			Content: "127.0.0.1",
			Type:    "A",
		}

		ipurl, iptype := _client.getUrlAndType(ipv4)

		address := ipv4.Content

		httpClient.FnGet = func(url string) (resp *http.Response, err error) {
			So(ipurl, ShouldEqual, URL_IPV4)
			resp = &http.Response{Body: newReadCloser(" 192.168.1.1 ")}
			return
		}

		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, REC_EDIT)
			resp = &http.Response{Body: newReadCloser(editResFailer)}
			return
		}

		(*_client).runOn(ipv4, ipurl, iptype, &address)

		So(address, ShouldEqual, "127.0.0.1")
		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("CfUpdater failed to update IP address (", iptype, ")"))
	})

	reset()

	Convey("Run updater on 'A' record, manage to get up-to-date IP address, the current IP address is out-of-date, got a response while updating IP, the result was successful", t, func() {
		ipv4 := &Object{
			Content: "127.0.0.1",
			Type:    "A",
		}

		ipurl, iptype := _client.getUrlAndType(ipv4)

		address := ipv4.Content

		httpClient.FnGet = func(url string) (resp *http.Response, err error) {
			So(ipurl, ShouldEqual, URL_IPV4)
			resp = &http.Response{Body: newReadCloser(" 192.168.1.1 ")}
			return
		}

		httpClient.FnPostForm = func(url string, data url.Values) (resp *http.Response, err error) {
			So(url, ShouldEqual, API_URL)
			So(data.Get("a"), ShouldEqual, REC_EDIT)
			resp = &http.Response{Body: newReadCloser(editResSuccess)}
			return
		}

		(*_client).runOn(ipv4, ipurl, iptype, &address)

		So(address, ShouldEqual, "192.168.1.1")
	})
}
