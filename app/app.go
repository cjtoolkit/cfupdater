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

		for _, name := range [...]string{"hour", "timeout", "debug"} {
			f := flag.Lookup(name)
			if f.Name == "timeout" {
				fmt.Print("\t\t-", f.Name, "\t", f.Usage)
			} else {
				fmt.Print("\t\t-", f.Name, "\t\t", f.Usage)
			}
			fmt.Println()
		}

		os.Exit(0)
	}

	ipv4, ipv6 := cf.Get()

	switch {

	case ipv4 != nil && ipv6 == nil:
		fmt.Println("Running CfUpdater! (IPv4 Only)")
		ipv4.Run()

	case ipv4 == nil && ipv6 != nil:
		fmt.Println("Running CfUpdater! (IPv6 Only)")
		ipv6.Run()

	case ipv4 != nil && ipv6 != nil:
		fmt.Println("Running CfUpdater! (IPv4 and IPv6)")
		go ipv6.Run()
		ipv4.Run()

	default:
		fmt.Println("Neither A or AAAA records were found!")
		os.Exit(1)

	}
}
