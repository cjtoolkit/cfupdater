package main

import (
	"github.com/cjtoolkit/cfupdater/src/cloudflare"
	"github.com/cjtoolkit/cfupdater/src/config"
	"log"
	"time"
)

func checkError(err error) {
	if nil != err {
		log.Fatalln(err.Error())
	}
}

func main() {
	checkError(config.InitConfig())

	log.Println("Running CfUpdater")

	records, err := cloudflare.NewDnsRecordsGetters().GetRecords()
	checkError(err)

	updaters := []cloudflare.RecordUpdater{}
	for _, record := range records {
		updaters = append(updaters, cloudflare.NewRecordUpdater(record))
	}

	hour := time.Duration(config.GetConfig().Hour) * time.Hour

	for {
		for _, updater := range updaters {
			updater.RunUpdater()
		}
		time.Sleep(hour)
	}
}
