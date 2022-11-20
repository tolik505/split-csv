// Package split_csv implements splitting of csv files on chunks by size in bytes
package split_csv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// minFileChunkSize min file chunk size in bytes
const minFileChunkSize = 100

var (
	ErrWrongSeparator     = errors.New("only one-character separators are supported")
	ErrSmallFileChunkSize = errors.New("file chunk size is too small")
	ErrBigFileChunkSize   = errors.New("file chunk size is bigger than input file")
)

// Splitter struct which contains options for splitting
// FileChunkSize - a size of chunk in bytes, should be set by client
// WithHeader - whether split csv with header (true by default)
type Splitter struct {
	FileChunkSize int //in bytes
	WithHeader    bool
	Separator     string
	bufferSize    int //in bytes
	fileOp        fileOperator
	stateFactory  stateInitializer
}

// New initializes Splitter struct
func New() Splitter {
	return Splitter{
		WithHeader:   true,
		Separator:    ",",
		bufferSize:   os.Getpagesize() * 128,
		fileOp:       fileOp{},
		stateFactory: stateFactory{},
	}
}

// Split splits file in smaller chunks
func (s Splitter) Split(inputFilePath string, outputDirPath string) ([]string, error) {
	_, err := strconv.Unquote(`'` + s.Separator + `'`)
	if err != nil {
		return nil, ErrWrongSeparator
	}
	if s.FileChunkSize < minFileChunkSize {
		return nil, ErrSmallFileChunkSize
	}

	stat, err := s.fileOp.Stat(inputFilePath)
	if err != nil {
		msg := fmt.Sprintf("Couldn't get file stat %s : %v", inputFilePath, err)
		return nil, errors.New(msg)
	}
	fileSize := stat.Size()
	if fileSize <= int64(s.FileChunkSize) {
		return nil, ErrBigFileChunkSize
	}

	file, err := s.fileOp.Open(inputFilePath)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open file %s : %v", inputFilePath, err)
		return nil, errors.New(msg)
	}
	defer file.Close()

	bufBulk := make([]byte, s.bufferSize)
	fileName, ext := getFileName(inputFilePath)
	st := s.stateFactory.Init(
		s,
		inputFilePath,
		fileName,
		ext,
		prepareResultDirPath(outputDirPath),
		file,
	)
	isFirstBulk := true
	for {
		//Read bulk from file
		size, err := file.Read(bufBulk)
		if err == io.EOF {
			_, err = st.chunkFile.Write(st.brokenLine)
			if err != nil {
				msg := fmt.Sprintf("Couldn't write chunk file %s : %v", st.chunkFilePath, err)
				return nil, errors.New(msg)
			}
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Couldn't read file bulk %s : %v", inputFilePath, err)
			return nil, errors.New(msg)
		}
		st.fileBuffer = bytes.NewBuffer(bufBulk[:size])
		if len(st.brokenLine) > 0 {
			if _, err := st.bulkBuffer.Write(st.brokenLine); err != nil {
				msg := fmt.Sprintf("Couldn't write broken line to the bulk buffer: %v", err)
				return nil, errors.New(msg)
			}
			st.brokenLine = []byte{}
		}

		if isFirstBulk {
			st.columnsCount = countColumns(bufBulk, []byte(s.Separator)[0])
			isFirstBulk = false
		}

		err = s.readLinesFromBulk(st)
		if err != nil {
			return nil, err
		}

		err = s.saveBulkToFile(st)
		if err != nil {
			return nil, err
		}
	}
	st.chunkFile.Close()

	return st.result, nil
}

// readLinesFromBulk reads bulk line by line
func (s Splitter) readLinesFromBulk(st *state) error {
	for {
		bytesLine, err := st.fileBuffer.ReadBytes('\n')
		if err == io.EOF {
			st.brokenLine = bytesLine
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Couldn't read bytes from buffer of file %s : %v", st.inputFilePath, err)
			return errors.New(msg)
		}
		if st.firstLine && st.s.WithHeader {
			st.firstLine = false
			st.header = bytesLine
			continue
		}
		if _, err := st.bulkBuffer.Write(bytesLine); err != nil {
			msg := fmt.Sprintf("Couldn't write to the bulk buffer: %v", err)
			return errors.New(msg)
		}
		if st.s.FileChunkSize < st.s.bufferSize && st.bulkBuffer.Len() >= (st.s.FileChunkSize-len(st.header)) {
			separator := []byte(s.Separator)[0]
			lineColumnsCount := countColumns(bytesLine, separator)
			if lineColumnsCount != st.columnsCount || isIncompleteLine(bytesLine, separator) {
				continue
			}
			err = s.saveBulkToFile(st)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isIncompleteLine(line []byte, separator byte) bool {
	openingQuote := false
	previousSeparator := false
	previousQuote := false
	quotesCount := 0

	for i := 0; i < len(line); i++ {
		if openingQuote && line[i] != '"' && previousQuote && quotesCount%2 == 0 {
			openingQuote = false
		}
		switch line[i] {
		case '"':
			quotesCount++
			if previousSeparator || i == 0 {
				openingQuote = true
			} else if !openingQuote {
				quotesCount--
			}
			previousQuote = true
		case separator:
			if !openingQuote {
				previousSeparator = true
			}
		case '\n':
			if openingQuote {
				return true
			}
		case ' ':
			break
		default:
			previousSeparator = false
		}
		if line[i] != '"' {
			previousQuote = false
		}
	}

	return false
}

func countColumns(content []byte, separator byte) int {
	result := 1
	openingQuote := false
	previousSeparator := false
	previousQuote := false
	quotesCount := 0
loop:
	for i := 0; i < len(content); i++ {
		if openingQuote && content[i] != '"' && previousQuote && quotesCount%2 == 0 {
			openingQuote = false
		}
		switch content[i] {
		case '"':
			quotesCount++
			if previousSeparator || i == 0 {
				openingQuote = true
			}
			previousQuote = true
		case separator:
			if !openingQuote {
				previousSeparator = true
				result++
			}
		case '\n':
			previousQuote = false
			if !openingQuote {
				break loop
			}
		case ' ':
			break
		default:
			previousSeparator = false
		}
		if content[i] != '"' {
			previousQuote = false
		}
	}

	return result
}

// saveBulkToFile saves lines from bulk to a new file
func (s Splitter) saveBulkToFile(st *state) error {
	st.chunkFilePath = st.resultDirPath + st.fileName + "_" + strconv.Itoa(st.chunk) + "." + st.ext
	stat, err := s.fileOp.Stat(st.chunkFilePath)
	if s.fileOp.IsNotExist(err) {
		chunkFile, err := s.fileOp.Create(st.chunkFilePath)
		if err != nil {
			msg := fmt.Sprintf("Couldn't create file %s : %v", st.chunkFilePath, err)
			return errors.New(msg)
		}
		st.setChunkFile(chunkFile)
		_, err = st.chunkFile.Write(st.header)
		if err != nil {
			msg := fmt.Sprintf("Couldn't write header of chunk file %s : %v", st.chunkFilePath, err)
			return errors.New(msg)
		}
		st.result = append(st.result, st.chunkFilePath)
	}
	_, err = st.chunkFile.Write(st.bulkBuffer.Bytes())
	if err != nil {
		msg := fmt.Sprintf("Couldn't write chunk file %s : %v", st.chunkFilePath, err)
		return errors.New(msg)
	}
	stat, _ = s.fileOp.Stat(st.chunkFilePath)
	if stat.Size() > int64(st.s.FileChunkSize-st.s.bufferSize) {
		st.chunk++
	}
	st.bulkBuffer.Reset()

	return nil
}

// getFileName extracts name and extension from path
func getFileName(path string) (string, string) {
	filenameArr := strings.Split(filepath.Base(path), ".")
	if len(filenameArr) == 2 {
		return filenameArr[0], filenameArr[1]
	}

	return filenameArr[0], ""
}

// prepareResultDirPath adds '/' to the end of path if needed
func prepareResultDirPath(path string) string {
	if path == "" {
		return ""
	}
	p := []byte(path)
	if p[len(p)-1] != os.PathSeparator {
		p = append(p, os.PathSeparator)
	}

	return string(p)
}
