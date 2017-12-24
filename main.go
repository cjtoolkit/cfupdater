package main

import (
	"log"
	"time"

	"github.com/cjtoolkit/cfupdater/src/cloudflare"
	"github.com/cjtoolkit/cfupdater/src/config"
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

	minute := time.Duration(config.GetConfig().Minute) * time.Minute

	for {
		for _, updater := range updaters {
			updater.RunUpdater()
		}
		time.Sleep(minute)
	}
}
