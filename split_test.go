package split_csv

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var files = []string{
	"./testdata/result/test_1.csv",
	"./testdata/result/test_2.csv",
	"./testdata/result/test_3.csv",
}

func setUp(t *testing.T) {
	for _, file := range files {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			continue
		}
		err = os.Remove(file)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_fs_split(t *testing.T) {
	setUp(t)
	input := "./testdata/test.csv"
	fs := FileSplit{FileChunkSize: 1000}
	result, _ := fs.Split(input, "./testdata/result")
	for i, item := range files {
		if i == 3 {
			break
		}
		fileActual, err := os.Open(item)
		if err != nil {
			t.Error(err)
		}
		defer fileActual.Close()
		fileExp, err := os.Open(item + ".expected")
		if err != nil {
			t.Error(err)
		}
		defer fileExp.Close()
		statActual, err := fileActual.Stat()
		if err != nil {
			t.Error(err)
		}
		statExp, err := fileExp.Stat()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, statExp.Size(), statActual.Size())
	}

	assert.Equal(t, files, result)
}
