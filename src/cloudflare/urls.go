package cloudflare

import (
	"net/http"
	"strings"
)

const (
	baseUrl = "https://api.cloudflare.com/client/v4"

	zoneIdentifier = ":zone_identifier"
	identifier     = ":identifier"

	listZones       = baseUrl + "/zones"
	listDnsRecords  = listZones + "/" + zoneIdentifier + "/dns_records"
	updateDnsRecord = listDnsRecords + "/" + identifier

	listZonesMethod       = http.MethodGet
	listDnsRecordsMethod  = http.MethodGet
	updateDnsRecordMethod = http.MethodPut
)

func urlSearchReplace(rawUrl string, urlData map[string]string) string {
	for key, value := range urlData {
		rawUrl = strings.Replace(rawUrl, key, value, 1)
	}

	return rawUrl
}
