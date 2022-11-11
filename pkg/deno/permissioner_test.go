package deno_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/authgear/authgear-deno/pkg/deno"

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
			{`{"name":"net","host":"1.1.1.1"}`, false, "global unicast: 1.1.1.1"},
			{`{"name":"net","host":"8.8.8.8"}`, false, "global unicast: 8.8.8.8"},
			{`{"name":"net","host":"0.0.0.0"}`, false, "unspecified: 0.0.0.0"},
		}
		ctx := context.Background()

		p := deno.DisallowIPPolicy(
			deno.DisallowGlobalUnicast,
			deno.DisallowInterfaceLocalMulticast,
			deno.DisallowLinkLocalUnicast,
			deno.DisallowLinkLocalMulticast,
			deno.DisallowLoopback,
			deno.DisallowMulticast,
			deno.DisallowPrivate,
			deno.DisallowUnspecified,
		)

		for _, c := range cases {
			Convey(c.descriptor, func() {
				var pd deno.PermissionDescriptor
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
