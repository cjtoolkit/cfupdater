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

func getData() *data {
	return &data{
		timeout: TIMEOUT,
		hour:    HOUR,
	}
}

func getClientService(data *data, logger *log.Logger) *client {
	return &client{
		httpClient: &http.Client{Timeout: time.Duration(data.timeout) * time.Second},
		url:        API_URL,
		data:       *data,
		logger:     logger,
	}
}

func GetUpdater(data *data, logger *log.Logger) *Updater {
	return &Updater{
		client: getClientService(data, logger),
		logger: logger,
	}
}
