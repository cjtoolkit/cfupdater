package updater

import (
	"encoding/json"
	"io/ioutil"
	"net"
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

	httpClient  httpClientInterface
	url         string
	data        Data
	logger      loggerInterface
	transformer transformer
}

func (c *client) getObjects() (ipv4, ipv6 *Object) {
	resp, err := c.httpClient.PostForm(c.url, c.transformer.getRecLoadAllValues(c.data))

	if err != nil {
		c.logger.Println("Http Client Error:", err.Error())
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

	ipv4 = c.filterObject(rec, "A")
	ipv6 = c.filterObject(rec, "AAAA")

	return
}

func (c *client) filterObject(rec *recloadall, recType string) (ob *Object) {
	for _, object := range rec.Response.Record.Objects {
		if object.Type == recType && object.Name == c.data.Name {
			ob = object
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
		c.logger.Println(ipurl, ": Http Error:", err.Error())
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

	resp, err = c.httpClient.PostForm(c.url, c.transformer.getRecEditValues(c.data, ob, respaddress))

	if err != nil {
		c.logger.Println(c.url, iptype, ": Http Error:", err.Error())
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
