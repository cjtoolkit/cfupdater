package app

import (
	"flag"
	"fmt"
	"github.com/cjtoolkit/cfupdater/app/cf"
	"github.com/cjtoolkit/cfupdater/app/settings"
	"os"
)

func Exec() {
	if *settings.Tkn == "" || *settings.Email == "" ||
		*settings.Z == "" || *settings.Name == "" {

		fmt.Println("CJToolkit CfUpdater Usage:")
		fmt.Println()

		fmt.Println("\tMandatory Options:")
		fmt.Println()

		for _, name := range [...]string{"tkn", "email", "z", "name"} {
			f := flag.Lookup(name)
			fmt.Print("\t\t-", f.Name, "\t\t", f.Usage)
			fmt.Println()
		}

		fmt.Println()

		fmt.Println("\tOther Options:")
		fmt.Println()

		for _, name := range [...]string{"ipv6", "hour", "debug"} {
			f := flag.Lookup(name)
			fmt.Print("\t\t-", f.Name, "\t\t", f.Usage)
			fmt.Println()
		}

		os.Exit(0)
	}

	fmt.Println("Running Cloudflare Updater!")
	ipv4, ipv6 := cf.Get()

	if ipv6 != nil {
		go ipv6.Run()
	}

	ipv4.Run()
}
