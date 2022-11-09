package ioutil

import (
	"io"
)

// LimitedWriter stops writting to the underlying W when
// N is non-positive.
// Unlike [io.LimitedReader], it does not return error when the limit exceeds.
type LimitedWriter[T io.Writer] struct {
	W        T
	N        int64
	Exceeded bool
}

func LimitWriter[T io.Writer](w T, n int64) *LimitedWriter[T] {
	return &LimitedWriter[T]{W: w, N: n}
}

func (w *LimitedWriter[T]) Write(p []byte) (n int, err error) {
	n = len(p)
	if w.N > 0 {
		if int64(len(p)) > w.N {
			p = p[:w.N]
			w.Exceeded = true
		}
		w.N -= int64(len(p))
		_, err = w.W.Write(p)
	} else {
		w.Exceeded = true
	}
	return
}
