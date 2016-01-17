package updater

type clientInterface interface {
	GetObjects() (ipv4, ipv6 *Object)
	RunOn(ob *Object)
}
