package ioutil_test

import (
	"bytes"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/authgear/authgear-deno/pkg/ioutil"
)

func TestLimitedWriter(t *testing.T) {
	Convey("LimitedWriter", t, func() {
		cases := []struct {
			n        int64
			input    []byte
			exceeded bool
		}{
			{0, nil, false},
			{0, []byte("1"), true},
			{1, []byte("1"), false},
			{1, []byte("12"), true},
			{3, []byte("12"), false},
		}

		for _, c := range cases {
			var buf bytes.Buffer
			w := ioutil.LimitWriter(&buf, c.n)
			written, err := io.Copy(w, bytes.NewReader(c.input))
			So(err, ShouldBeNil)
			So(written, ShouldEqual, len(c.input))
			So(w.Exceeded, ShouldEqual, c.exceeded)
		}
	})
}
