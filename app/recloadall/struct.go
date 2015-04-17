package recloadall

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cjtoolkit/cfupdater/app/dfmt"
	"github.com/cjtoolkit/cfupdater/app/settings"
)

type recloadall struct {
	Msg      *string   `json:"msg"`
	Result   string    `json:"result"`
	Response *response `json:"response"`
}

type response struct {
	Record *record `json:"recs"`
}

type record struct {
	Objects []*Object `json:"objs"`
}

type Object struct {
	Id          string `json:"rec_id"`
	Content     string `json:"content"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Ttl         string `json:"ttl"`
	ServiceMode string `json:"service_mode"`
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
		dfmt.Println("Getting IP Adddress (", t, ") from", url_)
		resp, err := http.DefaultClient.Get(url_)
		if err != nil {
			dfmt.Println(url_, "Timed out")
			time.Sleep(time.Duration(*settings.Hour) * time.Hour)
			continue
		}
		io.Copy(buf, resp.Body)
		tempaddress := strings.TrimSpace(buf.String())
		dfmt.Println(tempaddress)
		if tempaddress == address {
			dfmt.Println("Current IP (", t, ") is up-to-date!")
			goto sleep
		}
		address = tempaddress
		dfmt.Println("Updating IP Address (", t, ")!")
		http.DefaultClient.PostForm("https://www.cloudflare.com/api_json.html", url.Values{
			"a":            {"rec_edit"},
			"z":            {*settings.Z},
			"type":         {ob.Type},
			"id":           {ob.Id},
			"name":         {ob.Name},
			"content":      {address},
			"ttl":          {ob.Ttl},
			"service_mode": {ob.ServiceMode},
			"email":        {*settings.Email},
			"tkn":          {*settings.Tkn},
		})
	sleep:
		buf.Reset()
		time.Sleep(time.Duration(*settings.Hour) * time.Hour)
	}
}
