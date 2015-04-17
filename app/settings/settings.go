package settings

import (
	"flag"
)

var (
	Tkn   = flag.String("tkn", "", "Cloudflare API Key")
	Email = flag.String("email", "", "Cloudflare Email")
	Z     = flag.String("z", "", "Cloudflare Zone")
	Name  = flag.String("name", "", "Search Domain Name")
	Hour  = flag.Int64("hour", 2, "Run Every x Hours")
	IPv6  = flag.Bool("ipv6", false, "IPv6")
	Debug = flag.Bool("debug", false, "Debug Mode")
)

func init() {
	flag.Parse()
}
