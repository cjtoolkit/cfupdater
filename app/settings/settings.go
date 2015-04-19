package settings

import (
	"flag"
)

var (
	// Mandatory
	Tkn   = flag.String("tkn", "", "CloudFlare API Key")
	Email = flag.String("email", "", "CloudFlare Email")
	Z     = flag.String("z", "", "CloudFlare Zone")
	Name  = flag.String("name", "", "Search Domain Name in Records")

	// Optional
	Debug   = flag.Bool("debug", false, "Enable Debug Mode (Disabled by Default)")
	Timeout = flag.Int64("timeout", 30, "Specify API Timeout in Second (Default: 30)")
	Hour    = flag.Int64("hour", 2, "Run Every x Hours (2 Hours by Default)")
)

func init() {
	flag.Parse()
}
