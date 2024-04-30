package deno_test

import (
	"encoding/json"
	"testing"

	"github.com/authgear/authgear-deno/pkg/deno"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLineToPermissionDescriptor(t *testing.T) {
	Convey("lineToPermissionDescriptor", t, func() {
		cases := []struct {
			line     string
			expected string
		}{
			{`Deno requests run access.`, `{"name":"run"}`},
			{`Deno requests read access.`, `{"name":"read"}`},
			{`Deno requests write access.`, `{"name":"write"}`},
			{`Deno requests network access.`, `{"name":"net"}`},
			{`Deno requests env access.`, `{"name":"env"}`},
			{`Deno requests sys access.`, `{"name":"sys"}`},
			{`Deno requests ffi access.`, `{"name":"ffi"}`},
			{`Deno requests access to high precision time.`, `{"name":"hrtime"}`},

			{`Deno requests run access to "sh".`, `{"name":"run","command":"sh"}`},
			{`Deno requests read access to "/".`, `{"name":"read","path":"/"}`},
			{`Deno requests write access to "/".`, `{"name":"write","path":"/"}`},
			{`Deno requests network access to "localhost".`, `{"name":"net","host":"localhost"}`},
			{`Deno requests env access to "PATH".`, `{"name":"env","variable":"PATH"}`},
			{`Deno requests ffi access to "/".`, `{"name":"ffi","path":"/"}`},

			{`Deno requests sys access to "hostname".`, `{"name":"sys","kind":"hostname"}`},
			{`Deno requests sys access to "loadavg".`, `{"name":"sys","kind":"loadavg"}`},
			{`Deno requests sys access to "systemMemoryInfo".`, `{"name":"sys","kind":"systemMemoryInfo"}`},
			{`Deno requests sys access to "networkInterfaces".`, `{"name":"sys","kind":"networkInterfaces"}`},
			{`Deno requests sys access to "osRelease".`, `{"name":"sys","kind":"osRelease"}`},
			{`Deno requests sys access to "uid".`, `{"name":"sys","kind":"uid"}`},
			{`Deno requests sys access to "gid".`, `{"name":"sys","kind":"gid"}`},
		}

		for _, c := range cases {
			Convey(c.line, func() {
				d, ok := deno.LineToPermissionDescriptor(c.line)
				So(ok, ShouldBeTrue)
				b, _ := json.Marshal(d)
				So(string(b), ShouldEqualJSON, c.expected)
			})
		}
	})
}
