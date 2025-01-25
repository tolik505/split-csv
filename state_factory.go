package split_csv

import (
	"bytes"
)

type stateInitializer interface {
	Init(
		s Splitter,
		fileName string,
		resultDirPath string,
	) *state
}

type stateFactory struct{}

func (f stateFactory) Init(
	s Splitter,
	fileName string,
	resultDirPath string,
) *state {
	var header []byte
	if s.WithHeader {
		header = make([]byte, 0)
	}

	return &state{
		s:             s,
		fileName:      fileName,
		resultDirPath: resultDirPath,
		isFirstLine:   true,
		chunk:         1,
		bulkBuffer:    bytes.NewBuffer(make([]byte, 0, s.bufferSize)),
		header:        header,
	}
}
