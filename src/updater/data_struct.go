package updater

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

type editRes struct {
	Result string `json:"result"`
}

type data struct {
	tkn     string
	email   string
	z       string
	name    string
	timeout int64
	hour    int64
}
