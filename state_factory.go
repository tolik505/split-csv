package split_csv

import (
	"bytes"
	"io"
)

type stateInitializer interface {
	Init(
		s Splitter,
		inputFilePath string,
		fileName string,
		ext string,
		resultDirPath string,
		inputFile io.ReadCloser,
	) *state
}

type stateFactory struct{}

func (f stateFactory) Init(
	s Splitter,
	inputFilePath string,
	fileName string,
	ext string,
	resultDirPath string,
	inputFile io.ReadCloser,
) *state {
	var header []byte
	if s.WithHeader {
		header = make([]byte, 0)
	}

	return &state{
		s:             s,
		inputFilePath: inputFilePath,
		fileName:      fileName,
		ext:           ext,
		resultDirPath: resultDirPath,
		inputFile:     inputFile,
		isFirstLine:   true,
		chunk:         1,
		bulkBuffer:    bytes.NewBuffer(make([]byte, 0, s.bufferSize)),
		header:        header,
	}
}
