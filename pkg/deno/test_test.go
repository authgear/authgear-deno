package deno

import (
	"io"
	"os"
	"path"
	"strings"

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

	return ShouldResemble(string(content1), string(content2))
}
