package deno

import (
	"bytes"
)

// dropCR is from bufio.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// ScanStderr is the extended version of ScanLines.
// It additionally treats the deno permission prompt as a line.
func ScanStderr(data []byte, atEOF bool) (advance int, token []byte, err error) {
	prompt := []byte("(y = yes, allow; n = no, deny) > ")

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}

	if i := bytes.Index(data, prompt); i >= 0 {
		return i + len(prompt), dropCR(data[0 : i+len(prompt)]), nil
	}

	if atEOF {
		return len(data), dropCR(data), nil
	}

	return 0, nil, nil
}
