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

                -ipv6           Enable IPv6 (AAAA)
                -hour           Run Every x Hours (2 Hours by Default)
                -debug          Enable Debug Mode
~~~

## Example

~~~ sh
cfupdater \
	-tkn	"API Key Here" \
	-email	"email@exmaple.com" \
	-z		"example.com" \
	-name	"example.com"
~~~