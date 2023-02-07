package deno_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/authgear/authgear-deno/pkg/deno"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChecker(t *testing.T) {
	Convey("Checker", t, func() {
		ctx := context.Background()
		checker := &deno.Checker{}

		Convey("CheckFile", func() {
			targetScripts, err := filepath.Glob("./testdata/checker/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					expectedStderr, err := os.ReadFile(changeExtension(p, ".stderr"))
					So(err, ShouldBeNil)

					opts := deno.CheckFileOptions{
						TargetScript:         p,
						IsUnstableAPIAllowed: false,
					}
					err = checker.CheckFile(ctx, opts)

					if len(expectedStderr) <= 0 {
						So(err, ShouldBeNil)
					} else {
						var checkError *deno.CheckFileError
						So(errors.As(err, &checkError), ShouldBeTrue)
						So(checkError.Stderr, ShouldEqual, string(expectedStderr))
					}
				})
			}
		})

		Convey("CheckSnippet", func() {
			targetScripts, err := filepath.Glob("./testdata/checker/*.ts")
			So(err, ShouldBeNil)
			for _, p := range targetScripts {
				Convey(p, func() {
					targetScriptBytes, err := os.ReadFile(p)
					So(err, ShouldBeNil)

					expectedStderr, err := os.ReadFile(changeExtension(p, ".stderr"))
					So(err, ShouldBeNil)

					opts := deno.CheckSnippetOptions{
						TargetScript:         string(targetScriptBytes),
						IsUnstableAPIAllowed: false,
					}
					err = checker.CheckSnippet(ctx, opts)

					if len(expectedStderr) <= 0 {
						So(err, ShouldBeNil)
					} else {
						var checkError *deno.CheckFileError
						So(errors.As(err, &checkError), ShouldBeTrue)
						So(checkError.Stderr, ShouldEqual, string(expectedStderr))
					}
				})
			}
		})
	})
}
