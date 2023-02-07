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

type RunnerTestConfig struct {
	IsUnstableAPIAllowed bool `json:"is_unstable_api_allowed"`
}

var defaultRunnerConfig = RunnerTestConfig{
	IsUnstableAPIAllowed: false,
}

func readRunnerTestConfig(path string) (*RunnerTestConfig, error) {
	configBytes, err := os.ReadFile(changeExtension(path, ".config.json"))
	var testConfig RunnerTestConfig
	if errors.Is(err, os.ErrNotExist) {
		return &defaultRunnerConfig, nil
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configBytes, &testConfig)
	if err != nil {
		return nil, err
	}
	return &testConfig, nil
}

func TestRunner(t *testing.T) {
	Convey("Runner", t, func() {
		ctx := context.Background()
		runner := &deno.Runner{
			RunnerScript: "./runner.ts",
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

		Convey("RunFile", func() {
			targetScripts, err := filepath.Glob("./testdata/runner/good/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					testConfig, err := readRunnerTestConfig(p)
					So(err, ShouldBeNil)

					opts := deno.RunFileOptions{
						TargetScript:         p,
						Input:                changeExtension(p, ".in"),
						Output:               changeExtension(p, ".out"),
						IsUnstableAPIAllowed: testConfig.IsUnstableAPIAllowed,
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
					testConfig, err := readRunnerTestConfig(p)
					So(err, ShouldBeNil)
					opts := deno.RunFileOptions{
						TargetScript:         p,
						Input:                changeExtension(p, ".in"),
						Output:               changeExtension(p, ".out"),
						IsUnstableAPIAllowed: testConfig.IsUnstableAPIAllowed,
					}
					_, err = runner.RunFile(ctx, opts)
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

					testConfig, err := readRunnerTestConfig(p)
					So(err, ShouldBeNil)

					opts := deno.RunGoValueOptions{
						TargetScript:         targetScript,
						Input:                input,
						IsUnstableAPIAllowed: testConfig.IsUnstableAPIAllowed,
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
