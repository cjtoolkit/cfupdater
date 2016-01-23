package updater

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUpdater(t *testing.T) {
	// start let

	shouldNotBeCalledFn := func(name string) {
		t.Errorf("'%s' should not be called!", name)
	}

	logger := newLoggerMock()
	client := newClientMock(shouldNotBeCalledFn)

	updater := &Updater{
		client: client,
		logger: logger,
	}

	reset := func() {
		logger.Buf.Reset()
		client.reset()
	}

	// end let

	Convey("Method 'GetObject' returns 'A' record (ipv4) and run updater against that record", t, func() {
		client.FnGetObject = func() (ipv4, ipv6 *Object) {
			return &Object{Type: "A"}, nil
		}

		client.FnRunOn = func(ob *Object) {
			So(ob.Type, ShouldEqual, "A")
		}

		updater.ExecuteUpdater()

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("Running CfUpdater! (IPv4 Only)"))
	})

	reset()

	Convey("Method 'GetObject' returns 'AAAA' record (ipv6) and run updater against that record", t, func() {
		client.FnGetObject = func() (ipv4, ipv6 *Object) {
			return nil, &Object{Type: "AAAA"}
		}

		client.FnRunOn = func(ob *Object) {
			So(ob.Type, ShouldEqual, "AAAA")
		}

		updater.ExecuteUpdater()

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("Running CfUpdater! (IPv6 Only)"))
	})

	reset()

	Convey("Method 'GetObject' returns both 'A' and 'AAAA' record and run updater against both of them", t, func() {
		client.FnGetObject = func() (ipv4, ipv6 *Object) {
			return &Object{Type: "A"}, &Object{Type: "AAAA"}
		}

		typeMap := map[string]bool{
			"A":    true,
			"AAAA": true,
		}

		client.FnRunOn = func(ob *Object) {
			if !typeMap[ob.Type] {
				t.Errorf("'%s' shouldn't exist as Object.Type", ob.Type)
			}
		}

		updater.ExecuteUpdater()

		So(logger.Buf.String(), ShouldEqual, fmt.Sprintln("Running CfUpdater! (IPv4 and IPv6)"))
	})

	reset()

	Convey("Method 'GetObject' returns neither 'A' or 'AAAA', does not run updater and calls Fatal on logger", t, func() {
		defer func() {
			So(fmt.Sprint(recover()), ShouldEqual, "Neither A or AAAA records were found!")
		}()

		client.FnGetObject = func() (ipv4, ipv6 *Object) {
			return
		}

		updater.ExecuteUpdater()
	})
}
