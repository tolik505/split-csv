package split_csv

import "io"

type buffer interface {
	io.Writer
	ReadBytes(delim byte) (line []byte, err error)
	Bytes() []byte
	Len() int
	Reset()
}
