package split_csv

import (
	"bytes"
	"os"
)

type state struct {
	s             Splitter
	inputFilePath string
	fileName      string
	ext           string
	resultDirPath string
	inputFile     *os.File
	chunkFile     *os.File
	chunkFilePath string
	header        []byte
	firstLine     bool
	brokenLine    []byte
	chunk         int
	bulkBuffer    *bytes.Buffer
	fileBuffer    *bytes.Buffer
	result        []string
}

func (s *state) setChunkFile(file *os.File) {
	if s.chunkFile != nil {
		s.chunkFile.Close()
	}
	s.chunkFile = file
}
