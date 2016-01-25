package updater

import (
	"net/url"
)

const (
	REC_LOAD_ALL = "rec_load_all"
	REC_EDIT     = "rec_edit"
)

type transformer struct{}

func (_ transformer) getRecLoadAllValues(data Data) url.Values {
	return url.Values{
		"a":     {REC_LOAD_ALL},
		"tkn":   {data.Tkn},
		"email": {data.Email},
		"z":     {data.Z},
	}
}

func (_ transformer) getRecEditValues(data Data, ob *Object, respaddress string) url.Values {
	return url.Values{
		"a":            {REC_EDIT},
		"z":            {data.Z},
		"type":         {ob.Type},
		"id":           {ob.Id},
		"name":         {ob.Name},
		"content":      {respaddress},
		"ttl":          {ob.Ttl},
		"service_mode": {ob.ServiceMode},
		"email":        {data.Email},
		"tkn":          {data.Tkn},
	}
}
