# CJToolkit CfUpdater

Dynamic IP Updater for CloudFlare DNS Record

## Requirement

Google Go 1.5 or above.

## Installation

~~~
go get github.com/cjtoolkit/cfupdater
~~~

## How to Use

First save the configuration to `/home/:username/.cfupdater/config.json` (example below), than compile
and run cfupdater.

### Example config

~~~json
{
  "email": "hello@example.com",
  "api_key": "API Key Here",
  "zone": "example.com",
  "name": "example.com"
}
~~~

## Note

You can obtain an API Key from https://www.cloudflare.com/a/account/my-account.

If you want to run this on windows, make sure that `HOME` is set in environmental variables.
