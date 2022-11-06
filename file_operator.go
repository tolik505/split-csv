package split_csv

import (
	"io"
	"os"
)

type fileOperator interface {
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
}

type fileOp struct{}

func (f fileOp) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func (f fileOp) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (f fileOp) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f fileOp) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
