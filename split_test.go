package split_csv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tolik505/split-csv/mocks"
)

var filesDefaultFlow = []string{
	"testdata/result_default/test_1.csv",
	"testdata/result_default/test_2.csv",
	"testdata/result_default/test_3.csv",
}

var filesDefaultReaderFlow = []string{
	"testdata/result_custom_reader/test_1.csv",
	"testdata/result_custom_reader/test_2.csv",
	"testdata/result_custom_reader/test_3.csv",
}

var filesDefaultFlowMultiline = []string{
	"testdata/result_default/test_multiline_cells_1.csv",
	"testdata/result_default/test_multiline_cells_2.csv",
	"testdata/result_default/test_multiline_cells_3.csv",
}

var filesWithoutHeader = []string{
	"testdata/result_without_header/test_1.csv",
	"testdata/result_without_header/test_2.csv",
	"testdata/result_without_header/test_3.csv",
}

var filesWithoutHeaderMultiline = []string{
	"testdata/result_without_header/test_multiline_cells_1.csv",
	"testdata/result_without_header/test_multiline_cells_2.csv",
	"testdata/result_without_header/test_multiline_cells_3.csv",
}

var filesSmallBuffer = []string{
	"testdata/result_small_buffer/test_1.csv",
	"testdata/result_small_buffer/test_2.csv",
	"testdata/result_small_buffer/test_3.csv",
}

var filesSmallBufferMultiline = []string{
	"testdata/result_small_buffer/test_multiline_cells_1.csv",
	"testdata/result_small_buffer/test_multiline_cells_2.csv",
	"testdata/result_small_buffer/test_multiline_cells_3.csv",
}

var filesForExampleTest = []string{
	"testdata/test_1.csv",
	"testdata/test_2.csv",
	"testdata/test_3.csv",
}

type testReader struct {
	dataCh chan []byte
	buf    []byte
}

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

type stateFactoryMock struct {
	BulkBufferMock *mocks.Buffer
}

func (f *stateFactoryMock) Init(
	s Splitter,
	fileName string,
	resultDirPath string,
) *state {
	chunkFileMock := &mocks.WriteCloser{}
	chunkFileMock.EXPECT().Write([]byte("brokenLine")).Return(0, errors.New("write error"))

	return &state{
		s:             s,
		fileName:      fileName,
		resultDirPath: resultDirPath,
		isFirstLine:   true,
		chunk:         1,
		bulkBuffer:    f.BulkBufferMock,
		brokenLine:    []byte("brokenLine"),
		chunkFile:     chunkFileMock,
		chunkFilePath: "/chunkFile",
	}
}

func setUp(t *testing.T) {
	files := append(filesDefaultFlow, filesDefaultFlowMultiline...)
	files = append(files, filesDefaultReaderFlow...)
	files = append(files, filesWithoutHeader...)
	files = append(files, filesWithoutHeaderMultiline...)
	files = append(files, filesSmallBuffer...)
	files = append(files, filesSmallBufferMultiline...)
	files = append(files, filesForExampleTest...)
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

func Test_Split_integration(t *testing.T) {
	setUp(t)
	input := "testdata/test.csv"
	inputMultiline := "testdata/test_multiline_cells.csv"
	t.Run("Default flow", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 800
		s.bufferSize = 1000
		result, err := s.Split(input, "testdata/result_default")
		assertResult(t, result, filesDefaultFlow)
		assert.Nil(t, err)
	})
	t.Run("Default flow (multiline cells)", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 800
		result, err := s.Split(inputMultiline, "testdata/result_default")
		assertResult(t, result, filesDefaultFlowMultiline)
		assert.Nil(t, err)
	})
	t.Run("Without headers", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 800
		s.WithHeader = false
		result, err := s.Split(input, "testdata/result_without_header")
		assertResult(t, result, filesWithoutHeader)
		assert.Nil(t, err)
	})
	t.Run("Without headers (multiline cells)", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 800
		s.WithHeader = false
		result, err := s.Split(inputMultiline, "testdata/result_without_header")
		assertResult(t, result, filesWithoutHeaderMultiline)
		assert.Nil(t, err)
	})
	t.Run("With small buffer", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 800
		s.bufferSize = 100
		result, err := s.Split(input, "testdata/result_small_buffer/")
		assertResult(t, result, filesSmallBuffer)
		assert.Nil(t, err)
	})
	t.Run("With small buffer (multiline cells)", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 800
		s.bufferSize = 90
		result, err := s.Split(inputMultiline, "testdata/result_small_buffer/")
		assertResult(t, result, filesSmallBufferMultiline)
		assert.Nil(t, err)
	})
	t.Run("Wrong separator", func(t *testing.T) {
		s := New()
		s.Separator = "Ω"
		s.FileChunkSize = 800
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.Equal(t, err, errors.New("only one-byte separators are supported"))
	})
	t.Run("Big file chunk error", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 1000000
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.Equal(t, err, errors.New("file chunk size is bigger than input file"))
	})
	t.Run("Small file chunk error", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.Equal(t, err, errors.New("file chunk size is too small"))
	})
	t.Run("saveBulkToFile error", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 2000
		s.bufferSize = 1000
		result, err := s.Split(input, "wrong")

		assert.Nil(t, result)
		assert.EqualError(
			t,
			err,
			"Couldn't create file wrong/test_1.csv: open wrong/test_1.csv: no such file or directory",
		)
	})
	t.Run("readLinesFromBulk error", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 100
		result, err := s.Split(input, "wrong")

		assert.Nil(t, result)
		assert.EqualError(
			t,
			err,
			"Couldn't create file wrong/test_1.csv: open wrong/test_1.csv: no such file or directory",
		)
	})
	t.Run("File Stat error", func(t *testing.T) {
		fileOpMock := mocks.NewFileOperator(t)
		fileOpMock.EXPECT().Stat(input).Return(nil, errors.New("error"))
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 100
		s.fileOp = fileOpMock
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.EqualError(t, err, "Couldn't get file stat testdata/test.csv : error")
	})
	t.Run("File Open error", func(t *testing.T) {
		fileOpMock := mocks.NewFileOperator(t)
		stat, _ := os.Stat(input)
		fileOpMock.EXPECT().Stat(input).Return(stat, nil)
		fileOpMock.EXPECT().Open(input).Return(nil, errors.New("error"))
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 100
		s.fileOp = fileOpMock
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.EqualError(t, err, "Couldn't open file testdata/test.csv : error")
	})
	t.Run("File Read error", func(t *testing.T) {
		fileOpMock := mocks.NewFileOperator(t)
		stat, _ := os.Stat(input)
		fileOpMock.EXPECT().Stat(input).Return(stat, nil)
		fileMock := mocks.NewReadCloser(t)
		fileOpMock.EXPECT().Open(input).Return(fileMock, nil)
		fileMock.EXPECT().Read(mock.Anything).Return(0, errors.New("error"))
		fileMock.EXPECT().Close().Return(nil)
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 100
		s.fileOp = fileOpMock
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.EqualError(t, err, "Couldn't read file bulk: error")
	})
	t.Run("write brokenLine to the bulk buffer error", func(t *testing.T) {
		fileOpMock := mocks.NewFileOperator(t)
		stat, _ := os.Stat(input)
		fileOpMock.EXPECT().Stat(input).Return(stat, nil)
		fileMock := mocks.NewReadCloser(t)
		fileOpMock.EXPECT().Open(input).Return(fileMock, nil)
		fileMock.EXPECT().Read(mock.Anything).Return(0, io.EOF)
		fileMock.EXPECT().Close().Return(nil)
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 100
		s.fileOp = fileOpMock
		bulkBufferMock := &mocks.Buffer{}
		bulkBufferMock.EXPECT().
			Write([]byte("brokenLine")).
			Return(0, errors.New("buffer write error"))
		s.stateFactory = &stateFactoryMock{
			BulkBufferMock: bulkBufferMock,
		}
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.EqualError(
			t,
			err,
			"Couldn't write brokenLine to the bulk buffer: buffer write error",
		)
	})
	t.Run("saveBulkToFile error after writing a broken line", func(t *testing.T) {
		fileOpMock := mocks.NewFileOperator(t)
		stat, _ := os.Stat(input)
		fileOpMock.EXPECT().Stat(input).Return(stat, nil)
		fileMock := mocks.NewReadCloser(t)
		fileOpMock.EXPECT().Open(input).Return(fileMock, nil)
		fileMock.EXPECT().Read(mock.Anything).Return(0, io.EOF)
		fileMock.EXPECT().Close().Return(nil)
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 100
		s.fileOp = fileOpMock
		bulkBufferMock := &mocks.Buffer{}
		bulkBufferMock.EXPECT().Write([]byte("brokenLine")).Return(0, nil)
		fileOpMock.EXPECT().Stat("test_1.csv").Return(nil, errors.New("not exist"))
		fileOpMock.EXPECT().IsNotExist(errors.New("not exist")).Return(true)
		fileOpMock.EXPECT().Create("test_1.csv").Return(nil, errors.New("test error"))
		s.stateFactory = &stateFactoryMock{
			BulkBufferMock: bulkBufferMock,
		}
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.EqualError(t, err, "Couldn't create file test_1.csv: test error")
	})

	setUp(t)
}

func Test_SplitReader_integration(t *testing.T) {
	setUp(t)
	dataCh := make(chan []byte)
	reader := &testReader{dataCh: dataCh}
	data := []string{
		"Test header 1; Test header 2; Test header 3; Test header 4; Test header 5\n",
		"1; test value 1st; test value 1st; test value 1st; test value 1st\n",
		"2; test value 2nd; test value 2nd; test value 2nd; test value 2nd\n",
		"3; test value 3rd; test value 3rd; test value 3rd; test value 3rd\n",
		"4; test value 4th; test value 4th; test value 4th; test value 4th\n",
		"5; test value 5th; test value 5th; test value 5th; test value 5th\n",
		"6; test value 6th; test value 6th; test value 6th; test value 6th\n",
		"11; test value 1st; test value 1st; test value 1st; test value 1st\n",
		"12; test value 2nd; test value 2nd; test value 2nd; test value 2nd\n",
		"13; test value 3rd; test value 3rd; test value 3rd; test value 3rd\n",
		"14; test value 4th; test value 4th; test value 4th; test value 4th\n",
		"15; test value 5th; test value 5th; test value 5th; test value 5th\n",
		"16; test value 6th; test value 6th; test value 6th; test value 6th\n",
	}
	go func() {
		defer close(dataCh)
		for _, v := range data {
			dataCh <- []byte(v)
		}
	}()
	t.Run("Custom reader", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 400
		s.bufferSize = 100
		result, err := s.SplitReader(reader, "testdata/result_custom_reader", "test")
		assertResult(t, result, filesDefaultReaderFlow)
		assert.Nil(t, err)
	})

	t.Run("File reader", func(t *testing.T) {
		s := New()
		s.Separator = ";"
		s.FileChunkSize = 800
		s.bufferSize = 1000
		file, _ := s.fileOp.Open("testdata/test.csv")
		result, err := s.SplitReader(file, "testdata/result_default", "test")
		assertResult(t, result, filesDefaultFlow)
		assert.Nil(t, err)
	})
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

func Test_prepareResultDirPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "It adds backslash in the end of the path",
			args: args{path: "testdata/result"},
			want: "testdata/result/",
		},
		{
			name: "It doesn't add the second backslash in the end of the path",
			args: args{path: "testdata/result/"},
			want: "testdata/result/",
		},
		{
			name: "It returns an empty string when the path is empty",
			args: args{path: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				prepareResultDirPath(tt.args.path),
				"prepareResultDirPath(%v)",
				tt.args.path,
			)
		})
	}
}

func Test_getFileName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "It returns the file name and extension for the path with folders",
			args: args{
				path: "path/to/file.txt",
			},
			want: "file",
		},
		{
			name: "It returns the file name and extension for the path without folders",
			args: args{
				path: "file.txt",
			},
			want: "file",
		},
		{
			name: "It returns the file name and empty extension for the path without extension",
			args: args{
				path: "file",
			},
			want: "file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFileName(tt.args.path)
			assert.Equalf(t, tt.want, got, "getFileName(%v)", tt.args.path)
		})
	}
}

func TestSplitter_saveBulkToFile(t *testing.T) {
	type args struct {
		st     *state
		fileOp func(t *testing.T) fileOperator
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "It fails to create a file chunk",
			args: args{
				st: &state{},
				fileOp: func(t *testing.T) fileOperator {
					foMock := mocks.NewFileOperator(t)
					foMock.EXPECT().Stat("_0.csv").Return(nil, errors.New("isNotExist"))
					foMock.EXPECT().IsNotExist(errors.New("isNotExist")).Return(true)
					foMock.EXPECT().Create("_0.csv").Return(nil, errors.New("error"))

					return foMock
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(t, err, "Couldn't create file _0.csv: error")
			},
		},
		{
			name: "It fails to write the header into file chunk",
			args: args{
				st: &state{
					header: []byte{123},
				},
				fileOp: func(t *testing.T) fileOperator {
					foMock := mocks.NewFileOperator(t)
					foMock.EXPECT().Stat("_0.csv").Return(nil, errors.New("isNotExist"))
					foMock.EXPECT().IsNotExist(errors.New("isNotExist")).Return(true)
					chunkFile := mocks.NewWriteCloser(t)
					foMock.EXPECT().Create("_0.csv").Return(chunkFile, nil)
					chunkFile.EXPECT().Write([]byte{123}).Return(0, errors.New("error"))

					return foMock
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(
					t,
					err,
					"Couldn't write header of chunk file _0.csv : error",
				)
			},
		},
		{
			name: "It fails to write the buffer into file chunk",
			args: args{
				st: &state{
					bulkBuffer: bytes.NewBuffer([]byte{234}),
				},
				fileOp: func(t *testing.T) fileOperator {
					foMock := mocks.NewFileOperator(t)
					foMock.EXPECT().Stat("_0.csv").Return(nil, errors.New("isNotExist"))
					foMock.EXPECT().IsNotExist(errors.New("isNotExist")).Return(true)
					chunkFile := mocks.NewWriteCloser(t)
					foMock.EXPECT().Create("_0.csv").Return(chunkFile, nil)
					chunkFile.EXPECT().Write([]byte(nil)).Return(0, nil)
					chunkFile.EXPECT().Write([]byte{234}).Return(0, errors.New("error"))

					return foMock
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(t, err, "Couldn't write chunk file _0.csv : error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Splitter{
				fileOp: tt.args.fileOp(t),
			}
			tt.wantErr(
				t,
				s.saveBulkToFile(tt.args.st),
				fmt.Sprintf("saveBulkToFile(%v)", tt.args.st),
			)
		})
	}
}

func TestSplitter_readLinesFromBulk(t *testing.T) {
	type args struct {
		fileOp     func(t *testing.T) fileOperator
		fileBuffer func(t *testing.T) buffer
		bulkBuffer func(t *testing.T) buffer
	}
	defaultFileOp := func(t *testing.T) fileOperator { return fileOp{} }
	defaultBulkBuffer := func(t *testing.T) buffer { return bytes.NewBuffer([]byte{0}) }
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "It fails to read from buffer",
			args: args{
				fileBuffer: func(t *testing.T) buffer {
					fbMock := mocks.NewBuffer(t)
					fbMock.EXPECT().ReadBytes(uint8('\n')).Return(nil, errors.New("error"))

					return fbMock
				},
				fileOp:     defaultFileOp,
				bulkBuffer: defaultBulkBuffer,
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(
					t,
					err,
					"Couldn't read bytes from buffer: error",
				)
			},
		},
		{
			name: "It fails to write to the bulk buffer",
			args: args{
				fileBuffer: func(t *testing.T) buffer {
					fbMock := mocks.NewBuffer(t)
					fbMock.EXPECT().ReadBytes(uint8('\n')).Return([]byte("base line"), nil)

					return fbMock
				},
				fileOp: defaultFileOp,
				bulkBuffer: func(t *testing.T) buffer {
					fbMock := mocks.NewBuffer(t)
					fbMock.EXPECT().Write([]byte("base line")).Return(0, errors.New("error"))

					return fbMock
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(t, err, "Couldn't write to the bulk buffer: error")
			},
		},
		{
			name: "It fails on saveBulkToFile error",
			args: args{
				fileBuffer: func(t *testing.T) buffer {
					fbMock := mocks.NewBuffer(t)
					fbMock.EXPECT().ReadBytes(uint8('\n')).Return([]byte("base; line"), nil)

					return fbMock
				},
				bulkBuffer: func(t *testing.T) buffer {
					fbMock := mocks.NewBuffer(t)
					fbMock.EXPECT().Write([]byte("base; line")).Return(0, nil)
					fbMock.EXPECT().Len().Return(10)

					return fbMock
				},
				fileOp: func(t *testing.T) fileOperator {
					foMock := mocks.NewFileOperator(t)
					foMock.EXPECT().Stat("_0.csv").Return(nil, errors.New("isNotExist"))
					foMock.EXPECT().IsNotExist(errors.New("isNotExist")).Return(true)
					foMock.EXPECT().Create("_0.csv").Return(nil, errors.New("error"))

					return foMock
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(t, err, "Couldn't create file _0.csv: error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Splitter{
				fileOp:        tt.args.fileOp(t),
				FileChunkSize: 1,
				bufferSize:    10,
				Separator:     ";",
			}
			st := &state{
				s:            s,
				fileBuffer:   tt.args.fileBuffer(t),
				bulkBuffer:   tt.args.bulkBuffer(t),
				columnsCount: 2,
			}
			_, err := s.readLinesFromBulk(st)
			tt.wantErr(t, err, fmt.Sprintf("readLinesFromBulk(%v)", st))
		})
	}
}

func Test_countCompletedColumns(t *testing.T) {
	type args struct {
		bulkBytes []byte
		separator byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Simple rows",
			args: args{
				bulkBytes: []byte(
					`Test header 1; Test header 2; Test header 3; Test header 4; Test header 5
1; test value; test value; test value; test value`,
				),
				separator: ';',
			},
			want: 5,
		},
		{
			name: "Complex rows",
			args: args{
				bulkBytes: []byte(`""; test ""value""; """"; """test;abc
multiline;multiline;
value"; "test
value"
16; test value; test value; test value; test value`),
				separator: ';',
			},
			want: 5,
		},
		{
			name: "Complete line",
			args: args{
				bulkBytes: []byte(`1; test value; test value; test value; "test; value"
`),
				separator: ';',
			},
			want: 5,
		},
		{
			name: "Incomplete line",
			args: args{
				bulkBytes: []byte(`1; test value; "test; value"; test value; "test value
`),
				separator: ';',
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				countCompletedColumns(tt.args.bulkBytes, tt.args.separator),
				"countCompletedColumns(%v, %v)",
				tt.args.bulkBytes,
				tt.args.separator,
			)
		})
	}
}

func Test_isCompletingLine(t *testing.T) {
	type args struct {
		line      []byte
		separator byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Incomplete line 1",
			args: args{
				line: []byte(`value"; "test
`),
				separator: ';',
			},
			want: false,
		},
		{
			name: "Incomplete line 2",
			args: args{
				line: []byte(`""; test ""value""; """"; """test;abc
`),
				separator: ';',
			},
			want: false,
		},
		{
			name: "In not completing line",
			args: args{
				line:      []byte(`14; test value; test value; test value`),
				separator: ';',
			},
			want: false,
		},
		{
			name: "Completing line 1",
			args: args{
				line: []byte(`va;lue"
`),
				separator: ';',
			},
			want: true,
		},
		{
			name: "Completing line 2",
			args: args{
				line:      []byte(`lines"; "Test; header 4"; Test header 5`),
				separator: ';',
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				isCompletingLine(tt.args.line, tt.args.separator),
				"isCompletingLine(%v, %v)",
				tt.args.line,
				tt.args.separator,
			)
		})
	}
}
