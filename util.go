package xlog

import (
	"go/build"
	"strings"
)

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func trimSrcpath(s string) string {
	var r string
	r = strings.TrimPrefix(s, build.Default.GOROOT+"/src/")
	if r != s {
		return r
	}
	r = strings.TrimPrefix(s, build.Default.GOPATH+"/src/")
	if r != s {
		return r
	}
	return s
}

func trimDirs(s string) string {
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == '/' {
			return s[i+1:]
		}
	}
	return s
}
