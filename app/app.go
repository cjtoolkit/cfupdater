package app

import (
	"fmt"
	"github.com/cjtoolkit/cfupdater/app/recloadall"
)

func Exec() {
	fmt.Println("Running Cloudflare Updater!")
	ipv4, ipv6 := recloadall.Get()

	if ipv6 != nil {
		go ipv6.Run()
	}

	ipv4.Run()
}
