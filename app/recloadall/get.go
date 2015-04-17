package recloadall

import (
	"encoding/json"
	"net/http"
	"net/url"

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

	err = json.NewDecoder(resp.Body).Decode(&rec)
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
