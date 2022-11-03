package deno

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRunner(t *testing.T) {
	Convey("Runner", t, func() {
		ctx := context.Background()
		runner := NewRunner(DisallowIPPolicy(
			DisallowGlobalUnicast,
			DisallowInterfaceLocalMulticast,
			DisallowLinkLocalUnicast,
			DisallowLinkLocalMulticast,
			DisallowLoopback,
			DisallowMulticast,
			DisallowPrivate,
			DisallowUnspecified,
		))

		paths, err := filepath.Glob("./testdata/good/*")
		So(err, ShouldBeNil)
		for _, p := range paths {
			Convey(p, func() {
				err := runner.RunFile(ctx, p)
				So(err, ShouldBeNil)
			})
		}

		paths, err = filepath.Glob("./testdata/bad/*")
		So(err, ShouldBeNil)
		for _, p := range paths {
			Convey(p, func() {
				err := runner.RunFile(ctx, p)
				var exitError *exec.ExitError
				So(err, ShouldHaveSameTypeAs, exitError)
				exitError = err.(*exec.ExitError)
				So(exitError.ExitCode(), ShouldEqual, 1)
			})
		}
	})
}
