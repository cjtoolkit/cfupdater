package network

import (
	"github.com/cjtoolkit/cfupdater/src/iface"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
)

const (
	ipv4_url = "https://ipv4.icanhazip.com/"
	ipv6_url = "https://ipv6.icanhazip.com/"

	http_client_timeout time.Duration = 5 * time.Second
)

type Ip struct {
	currentAddress string
	client         iface.HttpClientInterface
	ipLookupUrl    string
}

func newIp(ipLookupUrl string) *Ip {
	return &Ip{
		client: &http.Client{
			Timeout: http_client_timeout,
		},
		ipLookupUrl: ipLookupUrl,
	}
}

func NewIpV4() *Ip {
	return newIp(ipv4_url)
}

func NewIpV6() *Ip {
	return newIp(ipv6_url)
}

func (ip *Ip) FetchIpAddress() (address string, updated bool) {
	currentAddress := ip.currentAddress

	resp, err := ip.client.Get(ip.ipLookupUrl)
	if nil != err {
		return
	}

	addressBytes, _ := ioutil.ReadAll(resp.Body)

	address = strings.TrimSpace(string(addressBytes))
	updated = address != currentAddress

	ip.currentAddress = address
	return
}
