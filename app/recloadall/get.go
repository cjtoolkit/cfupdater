package recloadall

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cjtoolkit/cfupdater/app/dfmt"
	"github.com/cjtoolkit/cfupdater/app/settings"
)

func Get() (ipv4, ipv6 *Object) {
	resp, err := http.DefaultClient.PostForm("https://www.cloudflare.com/api_json.html", url.Values{
		"a":     {"rec_load_all"},
		"tkn":   {*settings.Tkn},
		"email": {*settings.Email},
		"z":     {*settings.Z},
	})

	if err != nil {
		panic(err)
	}

	rec := &recloadall{}

	err = json.NewDecoder(resp.Body).Decode(rec)
	resp.Body.Close()

	if err != nil {
		panic(err)
	}

	dfmt.Println("Result: ", rec.Result)
	dfmt.Println()

	if rec.Result != "success" {
		str := ""
		if rec.Msg != nil {
			str = *rec.Msg
		}
		panic("CFUpdater: rec_load_all failed. " + str)
	}

	dfmt.Println("Record: ", rec.Response.Record.Objects)
	dfmt.Println()

	name := *settings.Name

	// Search IPv4

	for _, object := range rec.Response.Record.Objects {
		dfmt.Println(object)
		if object.Type == "A" && object.Name == name {
			ipv4 = object
			break
		}
	}

	if !*settings.IPv6 {
		return
	}

	// Search IPv6

	for _, object := range rec.Response.Record.Objects {
		dfmt.Println(object)
		if object.Type == "AAAA" && object.Name == name {
			ipv6 = object
			break
		}
	}

	return
}

func (ob *Object) Run() {
	url_ := "https://icanhazip.com/"
	t := "ipv4"
	if ob.Type == "AAAA" {
		url_ = "https://ipv6.icanhazip.com/"
		t = "ipv6"
	}

	buf := &bytes.Buffer{}
	address := strings.TrimSpace(ob.Content)

	dfmt.Println(address)

	for {
		tempaddress := ""
		var _editRes *editRes
		dfmt.Println("Getting IP Adddress (", t, ") from", url_)
		resp, err := http.DefaultClient.Get(url_)
		if err != nil {
			dfmt.Println(url_, ": Timed out")
			goto sleep
		}
		io.Copy(buf, resp.Body)
		resp.Body.Close()
		tempaddress = strings.TrimSpace(buf.String())
		buf.Reset()
		dfmt.Println("IP Adddres pull from", url_, ":", tempaddress)
		if tempaddress == address {
			dfmt.Println("Current IP (", t, ") is up-to-date!")
			goto sleep
		}

		dfmt.Println("Updating IP Address (", t, ")!")
		resp, err = http.DefaultClient.PostForm("https://www.cloudflare.com/api_json.html", url.Values{
			"a":            {"rec_edit"},
			"z":            {*settings.Z},
			"type":         {ob.Type},
			"id":           {ob.Id},
			"name":         {ob.Name},
			"content":      {tempaddress},
			"ttl":          {ob.Ttl},
			"service_mode": {ob.ServiceMode},
			"email":        {*settings.Email},
			"tkn":          {*settings.Tkn},
		})

		if err != nil {
			dfmt.Println("https://www.cloudflare.com/api_json.html : Timed out (", t, ")")
			goto sleep
		}

		err = json.NewDecoder(resp.Body).Decode(_editRes)
		if err != nil {
			dfmt.Println("CfUpdater json decoder failed (", t, ")")
			goto sleep
		}

		_editRes = &editRes{}

		if _editRes.Result != "success" {
			dfmt.Println("CfUpdater failed to update IP address (", t, ")")
			goto sleep
		}

		address = tempaddress
	sleep:
		time.Sleep(time.Duration(*settings.Hour) * time.Hour)
	}
}
