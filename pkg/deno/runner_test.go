package deno_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/authgear/authgear-deno/pkg/deno"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRunner(t *testing.T) {
	Convey("Runner", t, func() {
		ctx := context.Background()
		runner := &deno.Runner{
			Permissioner: deno.DisallowIPPolicy(
				deno.DisallowGlobalUnicast,
				deno.DisallowInterfaceLocalMulticast,
				deno.DisallowLinkLocalUnicast,
				deno.DisallowLinkLocalMulticast,
				deno.DisallowLoopback,
				deno.DisallowMulticast,
				deno.DisallowPrivate,
				deno.DisallowUnspecified,
			),
		}

		Convey("change default export and give helpful error message", func() {
			opts := deno.RunFileOptions{
				TargetScript: "./testdata/runner/no-default-export/hook.ts",
				Input:        "./testdata/runner/no-default-export/hook.in",
			}
			_, err := runner.RunFile(ctx, opts)

			var runError *deno.RunFileError
			var exitError *exec.ExitError
			So(errors.As(err, &runError), ShouldBeTrue)
			So(errors.As(err, &exitError), ShouldBeTrue)
			So(exitError.ExitCode(), ShouldEqual, 1)
			So(runError.Stderr.W.String(), ShouldEqual, "The hook must export a default function. Check that you have `export default async function(...) { ... }` in your script.\r\n")
		})

		Convey("RunFile", func() {
			targetScripts, err := filepath.Glob("./testdata/runner/good/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					opts := deno.RunFileOptions{
						TargetScript: p,
						Input:        changeExtension(p, ".in"),
						Output:       changeExtension(p, ".out"),
					}
					result, err := runner.RunFile(ctx, opts)
					So(err, ShouldBeNil)

					expected := opts.Output + ".expected"
					So(opts.Output, shouldEqualContent, expected)

					actualStdout := result.Stdout.W.Bytes()
					expectedStdout, err := os.ReadFile(changeExtension(p, ".stdout"))
					So(err, ShouldBeNil)
					So(string(actualStdout), ShouldEqual, string(expectedStdout))
				})
			}

			targetScripts, err = filepath.Glob("./testdata/runner/bad/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					opts := deno.RunFileOptions{
						TargetScript: p,
						Input:        changeExtension(p, ".in"),
						Output:       changeExtension(p, ".out"),
					}
					_, err := runner.RunFile(ctx, opts)
					var runError *deno.RunFileError
					var exitError *exec.ExitError
					// TODO: I wanted to match the stderr as well. But 2 problems block me.
					// 1. The stderr contains some ASCII escape sequence that Deno uses to clear the screen.
					// 2. The stack trace of Deno prints the absolute path of the script, which makes it hard to perform matching on different machines.
					So(errors.As(err, &runError), ShouldBeTrue)
					So(errors.As(err, &exitError), ShouldBeTrue)
					So(exitError.ExitCode(), ShouldEqual, 1)
				})
			}
		})

		Convey("RunGoValue", func() {
			targetScripts, err := filepath.Glob("./testdata/runner/good/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					targetScriptBytes, err := os.ReadFile(p)
					So(err, ShouldBeNil)
					targetScript := string(targetScriptBytes)

					inputBytes, err := os.ReadFile(changeExtension(p, ".in"))
					So(err, ShouldBeNil)

					var input interface{}
					err = json.Unmarshal(inputBytes, &input)
					So(err, ShouldBeNil)

					opts := deno.RunGoValueOptions{
						TargetScript: targetScript,
						Input:        input,
					}

					runGoValueResult, err := runner.RunGoValue(ctx, opts)
					So(err, ShouldBeNil)

					expectedBytes, err := os.ReadFile(changeExtension(p, ".out.expected"))
					So(err, ShouldBeNil)

					actualBytes, err := json.Marshal(runGoValueResult.Output)
					So(err, ShouldBeNil)

					So(string(actualBytes), ShouldEqualJSON, string(expectedBytes))
				})
			}
		})
	})
}
