<p align="center"><a href="https://godoc.org/github.com/tolik505/split-csv" target="_blank" rel="noopener noreferrer"><img width="250" src="https://repository-images.githubusercontent.com/212197147/d2207900-e626-11e9-827b-6faac4005ac1" alt="Vue logo"></a></p>

# Split csv

[![GoDoc](https://godoc.org/github.com/tolik505/split-csv?status.svg)](https://godoc.org/github.com/tolik505/split-csv)
[![Go Report Card](https://goreportcard.com/badge/github.com/tolik505/split-csv?style=flat-square)](https://goreportcard.com/report/github.com/tolik505/split-csv)
[![codecov](https://codecov.io/gh/tolik505/split-csv/branch/master/graph/badge.svg?token=YRJJN6J5XN)](https://codecov.io/gh/tolik505/split-csv)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv?ref=badge_shield)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/tolik505/split-csv/blob/master/LICENSE.MD)

Fast and efficient Golang package for splitting large csv files on smaller chunks by size in bytes.

## Features:

- Super-fast splitting. Splitting of 700MB+ file takes less than 1 sec!
- Allocates minimum memory regardless file size.
- Also accepts io.Reader as input.
- Supports multiline cells and headers (csv should follow the basic rules https://en.wikipedia.org/wiki/Comma-separated_values).
- Configurable destination folder.
- Disabling/enabling of copying a header in chunk files.

## Installation

Install:

```shell
go get -u github.com/tolik505/split-csv
```

Import:

```go
import splitCsv "github.com/tolik505/split-csv"
```

## Quickstart

```go
func ExampleSplitCsv() {
	splitter := splitCsv.New()
	splitter.Separator = ";"     // "," is by default
	splitter.FileChunkSize = 100000000 //in bytes (100MB)
	result, _ := splitter.Split("testdata/test.csv", "testdata/")
	fmt.Println(result)
	// Output: [testdata/test_1.csv testdata/test_2.csv testdata/test_3.csv]
}
```

If copying of a header in chunks is not needed then:

```go
func ExampleSplitCsv() {
	splitter := splitCsv.New()
	splitter.Separator = ";"     // "," is by default
	splitter.FileChunkSize = 20000000 //in bytes (20MB)
	s.WithHeader = false //copying of header in chunks is disabled
	result, _ := splitter.Split("testdata/test.csv", "testdata/")
	fmt.Println(result)
	// Output: [testdata/test_1.csv testdata/test_2.csv testdata/test_3.csv]
}
```

Or if you want to pass io.Reader instead of a file path:

```go
// First implement io.Reader interface with an appropriate logic for your use-case
type testReader struct {
	dataCh chan []byte
	buf    []byte
}
// Read listens to the data channel and populates p accordingly.
// When p is full remaining data goes to the buffer to be used in the next read cycle
func (r *testReader) Read(p []byte) (n int, err error) {
	pLen := len(p)
	i := 0
	for _, char := range r.buf {
		p[i] = char
		i++
	}
	r.buf = nil
	for bytes := range r.dataCh {
		for j, char := range bytes {
			p[i] = byte(char)
			i++
			if i >= pLen {
				r.buf = bytes[j+1:]

				return pLen, nil
			}
		}
	}

	return i, io.EOF
}
// In the example data is being sent to the channel which is consumed by the custom reader.
// In such way we can stream data to the splitter.
func ExampleSplitCsv() {
	dataCh := make(chan []byte)
	reader := &testReader{dataCh: dataCh}
	data := []string{
		"Test header 1; Test header 2; Test header 3; Test header 4; Test header 5\n",
		"1; test value 1st; test value 1st; test value 1st; test value 1st\n",
		"2; test value 2nd; test value 2nd; test value 2nd; test value 2nd\n",
		"3; test value 3rd; test value 3rd; test value 3rd; test value 3rd\n",
	}
	go func() {
		defer close(dataCh)
		for _, v := range data {
			dataCh <- []byte(v)
		}
	}()
	splitter := splitCsv.New()
	splitter.Separator = ";"     // "," is by default
	splitter.FileChunkSize = 100000000 //in bytes (100MB)
	result, _ := splitter.SplitReader(reader, "output/dir", "output_file_prefix")
	fmt.Println(result)
	// Output: [output/dir/test_1.csv output/dir/test_2.csv output/dir/test_3.csv]
}
```

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv?ref=badge_large)

