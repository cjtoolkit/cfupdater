# CJToolkit CfUpdater

Dynamic IP Updater for CloudFlare DNS Records

## Installation

~~~
go get github.com/cjtoolkit/cfupdater
~~~

## Terminal Usage

~~~
CJToolkit CfUpdater Usage:

        Mandatory Options:

                -tkn            CloudFlare API Key
                -email          CloudFlare Email
                -z              CloudFlare Zone
                -name           Search Domain Name in Records

        Other Options:

                -hour           Run Every x Hours (2 Hours by Default)
                -timeout        Specify API Timeout in Second (Default: 30)
                -debug          Enable Debug Mode (Disabled by Default)
~~~

## Example

~~~ sh
#!/bin/zsh

cfupdater \
	-tkn	"API Key Here" \
	-email	"email@exmaple.com" \
	-z		"example.com" \
	-name	"example.com"
~~~
