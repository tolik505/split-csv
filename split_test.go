package split_csv

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var filesDefaultFlow = []string{
	"testdata/result_default/test_1.csv",
	"testdata/result_default/test_2.csv",
	"testdata/result_default/test_3.csv",
}
var filesWithoutHeader = []string{
	"testdata/result_without_header/test_1.csv",
	"testdata/result_without_header/test_2.csv",
	"testdata/result_without_header/test_3.csv",
}
var filesSmallBuffer = []string{
	"testdata/result_small_buffer/test_1.csv",
	"testdata/result_small_buffer/test_2.csv",
	"testdata/result_small_buffer/test_3.csv",
}

func setUp(t *testing.T) {
	files := append(filesDefaultFlow, filesWithoutHeader...)
	files = append(files, filesSmallBuffer...)
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
	input := "testdata/test.csv"
	t.Run("Default flow", func(t *testing.T) {
		s := New()
		s.FileChunkSize = 800
		result, err := s.Split(input, "testdata/result_default")
		assertResult(t, result, filesDefaultFlow)
		assert.Nil(t, err)
	})
	t.Run("Without headers", func(t *testing.T) {
		s := New()
		s.FileChunkSize = 800
		s.WithHeader = false
		result, err := s.Split(input, "testdata/result_without_header")
		assertResult(t, result, filesWithoutHeader)
		assert.Nil(t, err)
	})
	t.Run("With small buffer", func(t *testing.T) {
		s := New()
		s.FileChunkSize = 800
		s.bufferSize = 100
		result, err := s.Split(input, "testdata/result_small_buffer/")
		assertResult(t, result, filesSmallBuffer)
		assert.Nil(t, err)
	})
	t.Run("Big file chunk", func(t *testing.T) {
		s := New()
		s.FileChunkSize = 1000000
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.Equal(t, err, ErrBigFileChunkSize)
	})
	t.Run("Small file chunk error", func(t *testing.T) {
		s := New()
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.Equal(t, err, ErrSmallFileChunkSize)
	})
	setUp(t)
}

func assertResult(t *testing.T, result []string, expected []string) {
	for i, item := range expected {
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

	assert.Equal(t, expected, result)
}
