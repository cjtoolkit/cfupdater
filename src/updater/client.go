package updater

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	URL_IPV4 = "https://icanhazip.com/"
	URL_IPV6 = "https://ipv6.icanhazip.com/"
)

/*
Implements:
	clientInterface
*/
type client struct {
	sync.Mutex

	httpClient httpClientInterface
	url        string
	data       Data
	logger     loggerInterface
}

func (c *client) getObjects() (ipv4, ipv6 *Object) {
	resp, err := c.httpClient.PostForm(c.url, url.Values{
		"a":     {"rec_load_all"},
		"tkn":   {c.data.Tkn},
		"email": {c.data.Email},
		"z":     {c.data.Z},
	})

	if err != nil {
		c.logger.Println("API Timeout")
		return
	}

	rec := &recloadall{}

	err = json.NewDecoder(resp.Body).Decode(rec)
	resp.Body.Close()

	if err != nil {
		c.logger.Println("JSON Decoder Failed")
		return
	}

	if rec.Result != "success" {
		str := ""
		if rec.Msg != nil {
			str = *rec.Msg
		}
		c.logger.Println("API Request Failed:", str)
		return
	}

	for _, object := range rec.Response.Record.Objects {
		if object.Type == "A" && object.Name == c.data.Name {
			ipv4 = object
			break
		}
	}

	for _, object := range rec.Response.Record.Objects {
		if object.Type == "AAAA" && object.Name == c.data.Name {
			ipv6 = object
			break
		}
	}
	return
}

func (c *client) GetObjects() (ipv4, ipv6 *Object) {
	for attemp := 0; attemp < 3; attemp++ {
		ipv4, ipv6 = c.getObjects()
		if ipv4 != nil || ipv6 != nil {
			break
		}
	}
	return
}

func (c *client) getUrlAndType(ob *Object) (url, _type string) {
	url = URL_IPV4
	_type = "ipv4"
	if ob.Type == "AAAA" {
		url = URL_IPV6
		_type = "ipv6"
	}
	return
}

func (c client) runOn(
	ob *Object,
	ipurl string,
	iptype string,
	address *string,
) {
	resp, err := c.httpClient.Get(ipurl)
	if err != nil {
		c.logger.Println(ipurl, ": Timed out")
		return
	}

	b, _ := ioutil.ReadAll(resp.Body)
	respaddress := strings.TrimSpace(string(b))
	b = nil

	if net.ParseIP(respaddress) == nil {
		c.logger.Println(respaddress, "is not a valid IP Address")
		return
	}

	if respaddress == *address {
		return
	}

	resp, err = c.httpClient.PostForm(c.url, url.Values{
		"a":            {"rec_edit"},
		"z":            {c.data.Z},
		"type":         {ob.Type},
		"id":           {ob.Id},
		"name":         {ob.Name},
		"content":      {respaddress},
		"ttl":          {ob.Ttl},
		"service_mode": {ob.ServiceMode},
		"email":        {c.data.Email},
		"tkn":          {c.data.Tkn},
	})

	if err != nil {
		c.logger.Println(c.url, ": Timed out (", iptype, ")")
		return
	}

	_editRes := &editRes{}

	err = json.NewDecoder(resp.Body).Decode(_editRes)
	if err != nil {
		c.logger.Println("CfUpdater json decoder failed (", iptype, ")")
		return
	}

	if _editRes.Result != "success" {
		c.logger.Println("CfUpdater failed to update IP address (", iptype, ")")
		return
	}

	*address = respaddress
}

func (c *client) RunOn(ob *Object) {
	c.Lock()
	cs := *c
	ipurl, iptype := c.getUrlAndType(ob)
	address := strings.TrimSpace(ob.Content)
	c.Unlock()

	for {
		cs.runOn(ob, ipurl, iptype, &address)
		time.Sleep(time.Duration(cs.data.Hour) * time.Hour)
	}
}
