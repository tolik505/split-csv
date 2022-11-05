<p align="center"><a href="https://godoc.org/github.com/tolik505/split-csv" target="_blank" rel="noopener noreferrer"><img width="250" src="https://repository-images.githubusercontent.com/212197147/d2207900-e626-11e9-827b-6faac4005ac1" alt="Vue logo"></a></p>

# split_csv [![GoDoc](https://godoc.org/github.com/tolik505/split-csv?status.svg)](https://godoc.org/github.com/tolik505/split-csv)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv?ref=badge_shield)
Fast and efficient Golang package for splitting large csv files on smaller chunks by size in bytes.


## Features:
- Super fast splitting. Splitting of 700MB+ file on 4 chunks takes less than 1 sec!
- Allocates minimum memory.
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
	splitter.FileChunkSize = 20000000 //in bytes (20MB)
	s.WithHeader = false //copying of header in chunks is disabled
	result, _ := splitter.Split("testdata/test.csv", "testdata/")
	fmt.Println(result)
	// Output: [testdata/test_1.csv testdata/test_2.csv testdata/test_3.csv]
}
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftolik505%2Fsplit-csv?ref=badge_large)