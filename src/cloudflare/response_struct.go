package cloudflare

type zoneBase struct {
	Success bool   `json:"success"`
	Result  []zone `json:"result"`
}

type zone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type dnsRecordBase struct {
	Success bool        `json:"success"`
	Result  []DnsRecord `json:"result"`
}

type DnsRecord struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	Content    string          `json:"content"`
	Proxiable  bool            `json:"proxiable"`
	Proxied    bool            `json:"proxied"`
	Ttl        int64           `json:"ttl"`
	Locked     bool            `json:"locked"`
	ZoneId     string          `json:"zone_id"`
	ZoneName   string          `json:"zone_name"`
	CreatedOn  string          `json:"created_on"`
	ModifiedOn string          `json:"modified_on"`
	Data       map[string]bool `json:"data"`
}
