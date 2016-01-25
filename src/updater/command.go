package updater

import (
	"github.com/cjtoolkit/cli"
	"io/ioutil"
	"log"
)

/*
Implements:
	CommandInterface in "github.com/cjtoolkit/cli"
*/
type command struct {
	data  *Data
	debug bool
}

func (c *command) CommandSetup(com *cli.Command) {
	c.data = getDataStorage()

	com.Init("cf:updater", "Cloudflare IP Updater")

	com.AddMandatoryArgument("tkn", "CloudFlare API Key", &c.data.Tkn)
	com.AddMandatoryArgument("email", "CloudFlare Email", &c.data.Email)
	com.AddMandatoryArgument("z", "CloudFlare Zone", &c.data.Z)
	com.AddMandatoryArgument("name", "Search Domain Name in Records", &c.data.Name)

	com.AddOptionalArgument("timeout", "Specify API Timeout in Second (Default: 30)", &c.data.Timeout)
	com.AddOptionalArgument("hour", "Run Every x Hours (2 Hours by Default)", &c.data.Hour)

	com.AddOption("debug", "Enable Debug Mode (Disabled by Default)", &c.debug)
}

func (c *command) CommandExec(l *log.Logger) {
	l.Println("Running CfUpdater Daemon")

	if !c.debug {
		l.SetOutput(ioutil.Discard)
	}

	GetUpdaterService(c.data, l).ExecuteUpdater()
}

func init() {
	cli.Register(&command{})
}
