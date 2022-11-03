package deno

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAllowRemoteIP(t *testing.T) {
	Convey("AllowRemoteIP", t, func() {
		cases := []struct {
			descriptor string
			expected   bool
			err        string
		}{
			{`{"name":"run"}`, false, ""},
			{`{"name":"read"}`, false, ""},
			{`{"name":"write"}`, false, ""},
			{`{"name":"net"}`, false, ""},
			{`{"name":"env"}`, false, ""},
			{`{"name":"sys"}`, false, ""},
			{`{"name":"ffi"}`, false, ""},
			{`{"name":"hrtime"}`, false, ""},

			{`{"name":"net","host":"127.0.0.1"}`, false, ""},
			{`{"name":"net","host":"127.0.0.1:8080"}`, false, ""},
			{`{"name":"net","host":"[::1]"}`, false, ""},
			{`{"name":"net","host":"[::1]:8080"}`, false, ""},
			{`{"name":"net","host":"1.1.1.1"}`, false, "ip is global unicast"},
			{`{"name":"net","host":"8.8.8.8"}`, false, "ip is global unicast"},
			{`{"name":"net","host":"0.0.0.0"}`, false, "ip is unspecified"},
		}
		ctx := context.Background()

		p := DisallowIPPolicy(
			DisallowGlobalUnicast,
			DisallowInterfaceLocalMulticast,
			DisallowLinkLocalUnicast,
			DisallowLinkLocalMulticast,
			DisallowLoopback,
			DisallowMulticast,
			DisallowPrivate,
			DisallowUnspecified,
		)

		for _, c := range cases {
			Convey(c.descriptor, func() {
				var pd PermissionDescriptor
				err := json.Unmarshal([]byte(c.descriptor), &pd)
				So(err, ShouldBeNil)
				actual, err := p.RequestPermission(ctx, pd)
				So(actual, ShouldEqual, c.expected)
				if c.expected == false && c.err != "" {
					So(err, ShouldBeError, c.err)
				}
			})
		}
	})
}
