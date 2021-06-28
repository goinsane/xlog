package xlog

import (
	"go/build"
	"os"
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

func trimSrcPath(s string) string {
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
		if s[i] == '/' || s[i] == os.PathSeparator {
			return s[i+1:]
		}
	}
	return s
}

type fmtState struct {
	Buffer []byte
	Wid    int
	WidOK  bool
	Prec   int
	PrecOK bool
	Flags  []rune
}

func (f *fmtState) Write(b []byte) (n int, err error) {
	f.Buffer = b
	return len(b), nil
}

func (f *fmtState) Width() (wid int, ok bool) {
	return f.Wid, f.WidOK
}

func (f *fmtState) Precision() (prec int, ok bool) {
	return f.Prec, f.PrecOK
}

func (f *fmtState) Flag(c int) bool {
	for _, flag := range f.Flags {
		if int(flag) == c {
			return true
		}
	}
	return false
}
