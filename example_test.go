package split_csv_test

import (
	"fmt"
	splitCsv "github.com/tolik505/split-csv"
)

func ExampleSplitCsv() {
	splitter := splitCsv.New()
	splitter.FileChunkSize = 800 //in bytes
	result, _ := splitter.Split("testdata/test.csv", "testdata/")
	fmt.Println(result)
	// Output: [testdata/test_1.csv testdata/test_2.csv testdata/test_3.csv]
}
