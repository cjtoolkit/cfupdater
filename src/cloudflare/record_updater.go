package cloudflare

import (
	"github.com/cjtoolkit/cfupdater/src/config"
	"github.com/cjtoolkit/cfupdater/src/iface"
	"github.com/cjtoolkit/cfupdater/src/network"
	"net/http"
	"time"
	"bytes"
	"encoding/json"
)

type RecordUpdater struct {
	client    iface.HttpClientInterface
	ip        *network.Ip
	dnsRecord *DnsRecord
	url       string
	httpRequest httpRequest
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
	}
}

func (rU RecordUpdater) RunUpdater() {
	address, updated := rU.ip.FetchIpAddress()
	if !updated {
		return
	}

	rU.dnsRecord.Content = address
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(*rU.dnsRecord)

	req := rU.httpRequest.newHttpRequest(updateDnsRecordMethod, rU.url)
	req.Body = buf

	rU.client.Do(req)
}