package updater

const (
	URL_IPV4 = "https://icanhazip.com/"
	URL_IPV6 = "https://ipv6.icanhazip.com/"
)

type client struct {
	client httpClientInterface
}

func GetObjects() (ipv4, ipv6 *Object) {
	ipv4 = nil
	ipv6 = nil
	return
}

func getUrlAndType(ob *Object) (url, _type string) {
	url = URL_IPV4
	_type = "ipv4"
	if ob.Type == "AAAA" {
		url = URL_IPV6
		_type = "ipv6"
	}
}

func RunOn(ob *Object) {

}
