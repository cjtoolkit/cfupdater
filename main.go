package main

import (
	_ "github.com/cjtoolkit/cfupdater/src"
	"github.com/cjtoolkit/cli"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "Cf:Updater", log.LstdFlags)
	cli.Start(logger, "CfUpdater")
}
