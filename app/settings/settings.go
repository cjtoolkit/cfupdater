package settings

import (
	"flag"
)

var (
	Tkn   = flag.String("tkn", "", "Cloudflare API Key")
	Email = flag.String("email", "", "Cloudflare Email")
	Z     = flag.String("z", "", "Cloudflare Zone")
	Name  = flag.String("name", "", "Search Domain Name in Records")
	Hour  = flag.Int64("hour", 2, "Run Every x Hours (2 Hours by Default)")

	IPv6  = flag.Bool("ipv6", false, "Enable IPv6 (AAAA)")
	Debug = flag.Bool("debug", false, "Enable Debug Mode")
)

func init() {
	flag.Parse()
}
