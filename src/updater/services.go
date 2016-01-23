package updater

import (
	"log"
	"net/http"
	"time"
)

const (
	API_URL = "https://www.cloudflare.com/api_json.html"
	TIMEOUT = int64(30)
	HOUR    = int64(2)
)

func getDataStorage() *Data {
	return &Data{
		Timeout: TIMEOUT,
		Hour:    HOUR,
	}
}

func getClientService(data *Data, logger *log.Logger) *client {
	return &client{
		httpClient: &http.Client{Timeout: time.Duration(data.Timeout) * time.Second},
		url:        API_URL,
		data:       *data,
		logger:     logger,
	}
}

func GetUpdaterService(data *Data, logger *log.Logger) *Updater {
	return &Updater{
		client: getClientService(data, logger),
		logger: logger,
	}
}
