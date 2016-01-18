package updater

import (
	"log"
	"time"
)

/*
Implements:
	UpdaterInterface
*/
type Updater struct {
	client clientInterface
	logger loggerInterface
}

func (u *Updater) ExecuteUpdater() {
	ipv4, ipv6 := u.client.GetObjects()

	switch {
	case ipv4 != nil && ipv6 == nil:
		u.logger.Println("Running CfUpdater! (IPv4 Only)")
		u.client.RunOn(ipv4)

	case ipv4 == nil && ipv6 != nil:
		u.logger.Println("Running CfUpdater! (IPv6 Only)")
		u.client.RunOn(ipv6)

	case ipv4 != nil && ipv6 != nil:
		u.logger.Println("Running CfUpdater! (IPv4 and IPv6)")
		go u.client.RunOn(ipv6)
		time.Sleep(100 * time.Millisecond)
		u.client.RunOn(ipv4)

	default:
		u.logger.Fatal("Neither A or AAAA records were found!")
	}
}
