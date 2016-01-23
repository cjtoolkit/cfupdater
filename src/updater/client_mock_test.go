package updater

/*
Implements:
	clientInterface
*/
type clientMock struct {
	Error       func(name string)
	FnGetObject func() (ipv4, ipv6 *Object)
	FnRunOn     func(ob *Object)
}

func newClientMock(Error func(name string)) *clientMock {
	c := &clientMock{
		Error: Error,
	}
	c.reset()
	return c
}

func (c *clientMock) reset() {
	c.FnGetObject = func() (ipv4, ipv6 *Object) {
		c.Error("GetObject")
		return
	}

	c.FnRunOn = func(ob *Object) {
		c.Error("RunOn")
	}
}

func (c *clientMock) GetObjects() (ipv4, ipv6 *Object) {
	ipv4, ipv6 = c.FnGetObject()
	return
}

func (c *clientMock) RunOn(ob *Object) {
	c.FnRunOn(ob)
}
