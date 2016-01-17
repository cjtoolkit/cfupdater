package updater

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
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
	parameters Parameters
	logger     loggerInterface
}

func (c *client) getObjects() (ipv4, ipv6 *Object) {
	resp, err := c.httpClient.PostForm(c.url, url.Values{
		"a":     {"rec_load_all"},
		"tkn":   {c.parameters.tkn},
		"email": {c.parameters.email},
		"z":     {c.parameters.z},
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
		if object.Type == "A" && object.Name == c.parameters.name {
			ipv4 = object
			break
		}
	}

	for _, object := range rec.Response.Record.Objects {
		if object.Type == "AAAA" && object.Name == c.parameters.name {
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
}

func (c *client) runOn(
	ob *Object,
	cfurl string,
	parameters Parameters,
	logger loggerInterface,
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
	tempaddress := strings.TrimSpace(string(b))
	b = nil

	if net.ParseIP(tempaddress) == nil {
		c.logger.Println(tempaddress, "is not a valid IP Address")
		return
	}

	if tempaddress == *address {
		return
	}

	resp, err = c.httpClient.PostForm(cfurl, url.Values{
		"a":            {"rec_edit"},
		"z":            {*parameters.z},
		"type":         {ob.Type},
		"id":           {ob.Id},
		"name":         {ob.Name},
		"content":      {ob.Content},
		"ttl":          {ob.Ttl},
		"service_mode": {ob.ServiceMode},
		"email":        {parameters.email},
		"tkn":          {parameters.tkn},
	})

	if err != nil {
		c.logger.Println(cfurl, ": Timed out (", iptype, ")")
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

	*address = tempaddress
}

func (c *client) RunOn(ob *Object) {
	c.Lock()
	cfurl := c.url
	parameters := c.parameters
	logger := c.logger
	c.Unlock()

	ipurl, iptype := c.getUrlAndType(ob)
	address := strings.TrimSpace(ob.Content)

	for {
		c.runOn(ob, cfurl, parameters, logger, ipurl, iptype, &address)
		time.Sleep(time.Duration(parameters.hour) * time.Hour)
	}
}
