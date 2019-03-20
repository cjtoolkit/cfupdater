# This project is now archived and is no longer maintained

I came to realise there are simplier ways to achieving this, just use a simple shell script
with curl and just use cron to schedule the script.  That way you get full control of where you
sources your WAN IP address from.

```sh
#!/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
IP=$(dig -4 +short myip.opendns.com @resolver1.opendns.com)
if [ "$IP" != "$(cat IpAddress/Server)" ]; then
        echo $IP > IpAddress/Server
        curl -X PUT \
        https://api.cloudflare.com/client/v4/zones/:zone_identifier/dns_records/:identifier \
        -H 'Content-Type: application/json' \
        -H 'X-Auth-Email: name@exmaple.com' \
        -H 'X-Auth-Key: Cloudflare API Key' \
        -H 'cache-control: no-cache' \
        -d "{\"type\":\"A\",\"name\":\"example.com\",\"content\":\"$IP\"}"
fi
```

[Api Doc for Updating DNS Record](https://api.cloudflare.com/#dns-records-for-a-zone-update-dns-record)

If you want to distribute your IP address file across multiple systems over LAN and WAN, you might want to checkout
https://syncthing.net/ it's awesome; and yes it's written in Go.  You could write your own version in Node, but I
strongly advices against that. 😉

Oh here a neat little trick of using the file with ssh.

```sh
$ ssh -o 'HostKeyAlias myhost' username@$(cat ~/IpAddress/Server)
```

I tested it, it's work wonderfully, make sure you add the alias to `~/.ssh/known_hosts`. 😊

Cheers,
Chris Jackson

# CJToolkit CfUpdater

Dynamic IP Updater for CloudFlare DNS Record

## Requirement

Google Go 1.5 or above.

## Installation

~~~
go get github.com/cjtoolkit/cfupdater
~~~

## How to Use

First save the configuration to `/home/:username/.config/cfupdater/config.json` (example below), than compile
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

It's is only designed to work with one `A` and one `AAAA` on one network, if you have other
domain names on the same network, use `CNAME` for other domain name and reference it against
the base name.
