package cloudflare

import (
	"encoding/json"
	"fmt"
	"github.com/cjtoolkit/cfupdater/src/config"
	"github.com/cjtoolkit/cfupdater/src/iface"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type DnsRecordsGetters struct {
	zone        string
	name        string
	client      iface.HttpClientInterface
	httpRequest httpRequestInterface
	log         iface.LoggerInterface
}

func NewDnsRecordsGetters() DnsRecordsGetters {
	config := config.GetConfig()

	return DnsRecordsGetters{
		zone: config.Zone,
		name: config.Name,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		log: log.New(os.Stdout, "DRG:", log.LstdFlags),
	}
}

func (d DnsRecordsGetters) GetRecords() (records []DnsRecord, err error) {
	var resp *http.Response
	resp, err = d.client.Do(
		d.httpRequest.newHttpRequest(listZonesMethod, listZones+"?"+(url.Values{
			"name":   {d.zone},
			"status": {"active"},
		}).Encode()),
	)
	if nil != err {
		return
	}

	zones := &zoneBase{}
	err = json.NewDecoder(resp.Body).Decode(zones)
	if nil != err {
		return
	} else if !zones.Success {
		err = fmt.Errorf("Zone has not been found. ")
		return
	}

	d.log.Println(fmt.Sprintf("Zone Id: %s ", zones.Result[0].Id))

	dnsUrl := urlSearchReplace(listDnsRecords, map[string]string{
		zoneIdentifier: zones.Result[0].Id,
	})

	for _, _type := range []string{"A", "AAAA"} {
		resp, err = d.client.Do(
			d.httpRequest.newHttpRequest(listDnsRecordsMethod, dnsUrl+"?"+(url.Values{
				"type": {_type},
				"name": {d.name},
			}).Encode()),
		)
		if nil == err {
			dnss := &dnsRecordBase{}
			err = json.NewDecoder(resp.Body).Decode(dnss)
			if nil == err && 0 != len(dnss.Result) {
				records = append(records, dnss.Result...)
				d.log.Println(fmt.Sprintf("'%s' record has been found, Id: %s ", _type, dnss.Result[0].Id))
			}
		}
	}

	if 0 == len(records) {
		err = fmt.Errorf("Could not find dns records for '%s' ", d.name)
	}

	return
}
