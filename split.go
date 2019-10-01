package split_csv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type FileSplit struct {
	FileChunkSize int
	WithHeader    bool
}

func (s FileSplit) Split(filePath string, resultPath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open file %s : %v", filePath, err)
		return nil, errors.New(msg)
	}
	defer file.Close()
	var result []string
	stat, err := file.Stat()
	if err != nil {
		msg := fmt.Sprintf("Couldn't get file stat %s : %v", filePath, err)
		return nil, errors.New(msg)
	}
	fileSize := stat.Size()
	if fileSize <= int64(s.FileChunkSize) {
		return []string{filePath}, nil
	}
	bufferSize := os.Getpagesize() * 128
	if s.FileChunkSize < bufferSize {
		bufferSize = s.FileChunkSize / 4
	}

	var m runtime.MemStats
	start := time.Now()
	firstLine := true
	chunk := 1
	var header []byte
	var brokenLine []byte
	bufBulk := make([]byte, bufferSize)
	bulkBuffer := bytes.NewBuffer(make([]byte, 0, bufferSize))
	var chunkFilePath string
	var chunkFile *os.File
	fileName, ext := getFileName(filePath)
	resPath := prepareResPath(resultPath)
	for {
		//Read bulk from file
		size, err := file.Read(bufBulk)
		if err == io.EOF {
			_, err = chunkFile.Write(brokenLine)
			if err != nil {
				msg := fmt.Sprintf("Couldn't write chunk file %s : %v", chunkFilePath, err)
				return nil, errors.New(msg)
			}
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Couldn't read file bulk %s : %v", filePath, err)
			return nil, errors.New(msg)
		}
		buffer := bytes.NewBuffer(bufBulk[:size])
		bulkBuffer.Reset()
		if len(brokenLine) > 0 {
			bulkBuffer.Write(brokenLine)
			brokenLine = []byte{}
		}
		for {
			//Read line from bulk
			bytesLine, err := buffer.ReadBytes('\n')
			if err == io.EOF {
				brokenLine = bytesLine
				break
			}
			if err != nil {
				msg := fmt.Sprintf("Couldn't read byte from file %s : %v", filePath, err)
				return nil, errors.New(msg)
			}
			if firstLine && s.WithHeader {
				firstLine = false
				header = bytesLine
				continue
			}
			bulkBuffer.Write(bytesLine)
		}
		//Save lines from bulk to a new file
		chunkFilePath = resPath + fileName + "_" + strconv.Itoa(chunk) + "." + ext
		stat, err := os.Stat(chunkFilePath)
		if os.IsNotExist(err) {
			chunkFile, err = os.Create(chunkFilePath)
			if err != nil {
				msg := fmt.Sprintf("Couldn't create file %s : %v", chunkFilePath, err)
				return nil, errors.New(msg)
			}
			defer chunkFile.Close()
			_, err = chunkFile.Write(header)
			if err != nil {
				msg := fmt.Sprintf("Couldn't write header of chunk file %s : %v", chunkFilePath, err)
				return nil, errors.New(msg)
			}
			result = append(result, chunkFilePath)
		}
		_, err = chunkFile.Write(bulkBuffer.Bytes())
		if err != nil {
			msg := fmt.Sprintf("Couldn't write chunk file %s : %v", chunkFilePath, err)
			return nil, errors.New(msg)
		}
		stat, _ = os.Stat(chunkFilePath)
		if stat.Size() > int64(s.FileChunkSize-bufferSize) {
			chunk++
		}
	}

	elapsed := time.Since(start)
	log.Printf("Splitting of %s took %s \n", file.Name(), elapsed)
	runtime.ReadMemStats(&m)
	log.Println("Alloc:", m.Alloc,
		"TotalAlloc:", m.TotalAlloc,
		"Sys:", m.Sys,
		"HeapAlloc:", m.HeapAlloc)

	return result, nil
}

func getFileName(path string) (string, string) {
	split := strings.Split(path, "/")
	name := split[len(split)-1]
	split = strings.Split(name, ".")
	ext := split[len(split)-1]
	nSplit := split[:len(split)-1]
	name = strings.Join(nSplit, "")

	return name, ext
}

func prepareResPath(path string) string {
	p := []byte(path)
	if p[len(p)-1] != '/' {
		p = append(p, '/')
	}

	return string(p)
}
