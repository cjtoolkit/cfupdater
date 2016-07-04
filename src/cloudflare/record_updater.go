package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cjtoolkit/cfupdater/src/config"
	"github.com/cjtoolkit/cfupdater/src/iface"
	"github.com/cjtoolkit/cfupdater/src/network"
	"log"
	"net/http"
	"time"
	"os"
)

type bufCloser struct {
	*bytes.Buffer
}

func (b bufCloser) Close() error {
	return nil
}

type RecordUpdater struct {
	client      iface.HttpClientInterface
	ip          *network.Ip
	dnsRecord   *DnsRecord
	url         string
	httpRequest httpRequestInterface
	log         iface.LoggerInterface
}

var ipUrlsLookup = map[string]string{
	"A":    network.Ipv4_url,
	"AAAA": network.Ipv6_url,
}

func NewRecordUpdater(dnsRecord DnsRecord) RecordUpdater {
	config := config.GetConfig()

	return RecordUpdater{
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		ip:        network.NewIp(ipUrlsLookup[dnsRecord.Type]),
		dnsRecord: &dnsRecord,
		url: urlSearchReplace(updateDnsRecord, map[string]string{
			zoneIdentifier: dnsRecord.ZoneId,
			identifier:     dnsRecord.Id,
		}),
		httpRequest: httpRequest{},
		log: log.New(os.Stdout, "RU: ", log.LstdFlags),
	}
}

func (rU RecordUpdater) RunUpdater() {
	address, updated := rU.ip.FetchIpAddress()
	if !updated {
		return
	}

	rU.log.Println(fmt.Sprintf("Submitting updated IP address '%s'", address))

	rU.dnsRecord.Content = address
	buf := &bytes.Buffer{}
	defer buf.Reset()

	json.NewEncoder(buf).Encode(*rU.dnsRecord)

	req := rU.httpRequest.newHttpRequest(updateDnsRecordMethod, rU.url)
	req.Body = bufCloser{buf}

	rU.client.Do(req)
}
