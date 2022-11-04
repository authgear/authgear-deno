package deno

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
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

		Convey("RunFile", func() {
			targetScripts, err := filepath.Glob("./testdata/good/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					opts := RunFileOptions{
						TargetScript: p,
						Input:        changeExtension(p, ".in"),
						Output:       changeExtension(p, ".out"),
					}
					_, err := runner.RunFile(ctx, opts)
					So(err, ShouldBeNil)

					expected := opts.Output + ".expected"
					So(opts.Output, shouldEqualContent, expected)
				})
			}

			targetScripts, err = filepath.Glob("./testdata/bad/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					opts := RunFileOptions{
						TargetScript: p,
						Input:        changeExtension(p, ".in"),
						Output:       changeExtension(p, ".out"),
					}
					_, err := runner.RunFile(ctx, opts)
					var runError *RunFileError
					var exitError *exec.ExitError
					So(errors.As(err, &runError), ShouldBeTrue)
					So(errors.As(err, &exitError), ShouldBeTrue)
					So(exitError.ExitCode(), ShouldEqual, 1)
				})
			}
		})

		Convey("RunGoValue", func() {
			targetScripts, err := filepath.Glob("./testdata/good/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					targetScriptBytes, err := ioutil.ReadFile(p)
					So(err, ShouldBeNil)
					targetScript := string(targetScriptBytes)

					inputBytes, err := ioutil.ReadFile(changeExtension(p, ".in"))
					So(err, ShouldBeNil)

					var input interface{}
					err = json.Unmarshal(inputBytes, &input)
					So(err, ShouldBeNil)

					opts := RunGoValueOptions{
						TargetScript: targetScript,
						Input:        input,
					}

					runGoValueResult, err := runner.RunGoValue(ctx, opts)
					So(err, ShouldBeNil)

					expectedBytes, err := ioutil.ReadFile(changeExtension(p, ".out.expected"))
					So(err, ShouldBeNil)

					actualBytes, err := json.Marshal(runGoValueResult.Output)
					So(err, ShouldBeNil)

					So(string(actualBytes), ShouldEqualJSON, string(expectedBytes))
				})
			}
		})
	})
}
