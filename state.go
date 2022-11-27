package split_csv

import (
	"io"
)

type state struct {
	s             Splitter
	inputFilePath string
	fileName      string
	ext           string
	resultDirPath string
	inputFile     io.ReadCloser
	chunkFile     io.WriteCloser
	chunkFilePath string
	header        []byte
	isFirstLine   bool
	brokenLine    []byte
	chunk         int
	bulkBuffer    buffer // to buffer a bulk to be stored as a chunk file
	fileBuffer    buffer // to buffer a chunk of the input file
	columnsCount  int
	result        []string
}

func (s *state) setChunkFile(file io.WriteCloser) {
	if s.chunkFile != nil {
		s.chunkFile.Close()
	}
	s.chunkFile = file
}
