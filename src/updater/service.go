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

func GetParameters() *Parameters {
	return &Parameters{
		timeout: TIMEOUT,
		hour:    HOUR,
	}
}

func getClientService(parameters *Parameters, logger *log.Logger) *client {
	return &client{
		httpClient: &http.Client{Timeout: time.Duration(parameters.timeout) * time.Second},
		url:        API_URL,
		parameters: *parameters,
		logger:     logger,
	}
}

func GetUpdater(parameters *Parameters, logger *log.Logger) *Updater {
	return &Updater{
		parameters: *parameters,
		client:     getClientService(parameters, logger),
		logger:     logger,
	}
}
