package logging

import (
	"time"
)

type buffer []byte

func newBuffer() *buffer {
	b := make(buffer, 0)
	return &b
}

func itoa(buf *buffer, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// format yyyy/mm/dd hh:mm:ss level file:line:
func (buf *buffer) formatHeader(t time.Time, l Level, file string, line int) {
	// date
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	// time
	hour, minute, second := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, minute, 2)
	*buf = append(*buf, ':')
	itoa(buf, second, 2)
	*buf = append(*buf, ' ')
	// level
	*buf = append(*buf, l.String()...)
	*buf = append(*buf, ' ')
	// file:line
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			file = file[i+1:]
			break
		}
	}
	*buf = append(*buf, file...)
	*buf = append(*buf, ':')
	itoa(buf, line, -1)
	*buf = append(*buf, ": "...)
}

func (buf *buffer) WriteString(s string) {
	*buf = append(*buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		*buf = append(*buf, '\n')
	}
}

func (buf *buffer) Bytes() []byte {
	return *buf
}

func (buf *buffer) Reset() {
	*buf = (*buf)[:0]
}
