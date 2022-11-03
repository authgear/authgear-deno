package deno

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func changeExtension(p string, newExt string) string {
	ext := path.Ext(p)
	prefix := strings.TrimSuffix(p, ext)
	return prefix + newExt
}

func shouldEqualContent(actual interface{}, expected ...interface{}) string {
	f1, err := os.Open(actual.(string))
	if err != nil {
		return err.Error()
	}
	defer f1.Close()

	f2, err := os.Open(expected[0].(string))
	if err != nil {
		return err.Error()
	}
	defer f2.Close()

	content1, err := io.ReadAll(f1)
	if err != nil {
		return err.Error()
	}

	content2, err := io.ReadAll(f2)
	if err != nil {
		return err.Error()
	}

	return ShouldResemble(content1, content2)
}

func TestRunner(t *testing.T) {
	Convey("Runner", t, func() {
		ctx := context.Background()
		runner := &Runner{
			RunnerScript: "./runner.ts",
			Permissioner: DisallowIPPolicy(
				DisallowGlobalUnicast,
				DisallowInterfaceLocalMulticast,
				DisallowLinkLocalUnicast,
				DisallowLinkLocalMulticast,
				DisallowLoopback,
				DisallowMulticast,
				DisallowPrivate,
				DisallowUnspecified,
			),
		}

		targetScripts, err := filepath.Glob("./testdata/good/*.ts")
		So(err, ShouldBeNil)
		for _, p := range targetScripts {
			Convey(p, func() {
				opts := RunOptions{
					TargetScript: p,
					Input:        changeExtension(p, ".in"),
					Output:       changeExtension(p, ".out"),
				}
				err := runner.RunFile(ctx, opts)
				So(err, ShouldBeNil)

				expected := opts.Output + ".expected"
				So(opts.Output, shouldEqualContent, expected)
			})
		}

		targetScripts, err = filepath.Glob("./testdata/bad/*.ts")
		So(err, ShouldBeNil)
		for _, p := range targetScripts {
			Convey(p, func() {
				opts := RunOptions{
					TargetScript: p,
					Input:        changeExtension(p, ".in"),
					Output:       changeExtension(p, ".out"),
				}
				err := runner.RunFile(ctx, opts)
				var exitError *exec.ExitError
				So(err, ShouldHaveSameTypeAs, exitError)
				exitError = err.(*exec.ExitError)
				So(exitError.ExitCode(), ShouldEqual, 1)
			})
		}
	})
}
