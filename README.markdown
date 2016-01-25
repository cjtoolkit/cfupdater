# CJToolkit CfUpdater

Dynamic IP Updater for CloudFlare DNS Records

## Installation

~~~
go get github.com/cjtoolkit/cfupdater
~~~

## Terminal Usage

~~~
CJToolkit CfUpdater Usage:

        cf:updater [options] tkn email z name [timeout hour]

        Mandatory Argument:
             email          CloudFlare Email (type: *string)
             name           Search Domain Name in Records (type: *string)
             tkn            CloudFlare API Key (type: *string)
             z              CloudFlare Zone (type: *string)

        Optional Argument:
             hour           Run Every x Hours (2 Hours by Default) (type: *int64)
             timeout        Specify API Timeout in Second (Default: 30) (type: *int64)

        Options:
             --debug        Enable Debug Mode (Disabled by Default)
~~~

## Example

~~~ sh
#!/bin/zsh

cfupdater cf:updater \
	"API Key Here" \
	"email@exmaple.com" \
	"example.com" \
	"example.com"
~~~
