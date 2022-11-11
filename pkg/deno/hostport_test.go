package deno_test

import (
	"net"
	"testing"

	"github.com/authgear/authgear-deno/pkg/deno"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHostPort(t *testing.T) {
	Convey("HostPort", t, func() {
		hostport, err := deno.ParseHostPort("")
		So(err, ShouldBeNil)
		So(hostport, ShouldBeNil)

		hostport, err = deno.ParseHostPort("localhost")
		So(err, ShouldBeNil)
		So(hostport, ShouldResemble, &deno.HostPort{
			Host: "localhost",
		})

		hostport, err = deno.ParseHostPort("127.0.0.1")
		So(err, ShouldBeNil)
		So(hostport, ShouldResemble, &deno.HostPort{
			Host: "127.0.0.1",
			IPv4: net.IPv4(127, 0, 0, 1).To4(),
		})

		hostport, err = deno.ParseHostPort("[::1]")
		So(err, ShouldBeNil)
		So(hostport, ShouldResemble, &deno.HostPort{
			Host: "::1",
			IPv6: net.ParseIP("::1"),
		})

		hostport, err = deno.ParseHostPort("localhost:80")
		So(err, ShouldBeNil)
		So(hostport, ShouldResemble, &deno.HostPort{
			Host: "localhost",
			Port: "80",
		})

		hostport, err = deno.ParseHostPort("127.0.0.1:80")
		So(err, ShouldBeNil)
		So(hostport, ShouldResemble, &deno.HostPort{
			Host: "127.0.0.1",
			IPv4: net.IPv4(127, 0, 0, 1).To4(),
			Port: "80",
		})

		hostport, err = deno.ParseHostPort("[::1]:80")
		So(err, ShouldBeNil)
		So(hostport, ShouldResemble, &deno.HostPort{
			Host: "::1",
			IPv6: net.ParseIP("::1"),
			Port: "80",
		})
	})
}
