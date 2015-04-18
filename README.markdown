# CJToolkit CfUpdater

Dynamic IP Updater for Cloudflare DNS Records

## Installation

~~~
go get github.com/cjtoolkit/cfupdater
~~~

## Terminal Usage

~~~
CJToolkit CfUpdater Usage:

        Mandatory Options:

                -tkn            Cloudflare API Key
                -email          Cloudflare Email
                -z              Cloudflare Zone
                -name           Search Domain Name in Records

        Other Options:

                -ipv6           Enable IPv6 (AAAA) (Disabled by Default)
                -hour           Run Every x Hours (2 Hours by Default)
                -timeout        Specify API Timeout in Second (Default: 30)
                -debug          Enable Debug Mode (Disabled by Default)
~~~

## Example

~~~ sh
cfupdater \
	-tkn	"API Key Here" \
	-email	"email@exmaple.com" \
	-z		"example.com" \
	-name	"example.com"
~~~