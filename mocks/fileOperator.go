// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	io "io"
	fs "io/fs"

	mock "github.com/stretchr/testify/mock"
)

// FileOperator is an autogenerated mock type for the FileOperator type
type FileOperator struct {
	mock.Mock
}

type FileOperator_Expecter struct {
	mock *mock.Mock
}

func (_m *FileOperator) EXPECT() *FileOperator_Expecter {
	return &FileOperator_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: name
func (_m *FileOperator) Create(name string) (io.WriteCloser, error) {
	ret := _m.Called(name)

	var r0 io.WriteCloser
	if rf, ok := ret.Get(0).(func(string) io.WriteCloser); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.WriteCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FileOperator_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type FileOperator_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - name string
func (_e *FileOperator_Expecter) Create(name interface{}) *FileOperator_Create_Call {
	return &FileOperator_Create_Call{Call: _e.mock.On("Create", name)}
}

func (_c *FileOperator_Create_Call) Run(run func(name string)) *FileOperator_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *FileOperator_Create_Call) Return(_a0 io.WriteCloser, _a1 error) *FileOperator_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// IsNotExist provides a mock function with given fields: err
func (_m *FileOperator) IsNotExist(err error) bool {
	ret := _m.Called(err)

	var r0 bool
	if rf, ok := ret.Get(0).(func(error) bool); ok {
		r0 = rf(err)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FileOperator_IsNotExist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsNotExist'
type FileOperator_IsNotExist_Call struct {
	*mock.Call
}

// IsNotExist is a helper method to define mock.On call
//   - err error
func (_e *FileOperator_Expecter) IsNotExist(err interface{}) *FileOperator_IsNotExist_Call {
	return &FileOperator_IsNotExist_Call{Call: _e.mock.On("IsNotExist", err)}
}

func (_c *FileOperator_IsNotExist_Call) Run(run func(err error)) *FileOperator_IsNotExist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error))
	})
	return _c
}

func (_c *FileOperator_IsNotExist_Call) Return(_a0 bool) *FileOperator_IsNotExist_Call {
	_c.Call.Return(_a0)
	return _c
}

// Open provides a mock function with given fields: name
func (_m *FileOperator) Open(name string) (io.ReadCloser, error) {
	ret := _m.Called(name)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FileOperator_Open_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Open'
type FileOperator_Open_Call struct {
	*mock.Call
}

// Open is a helper method to define mock.On call
//   - name string
func (_e *FileOperator_Expecter) Open(name interface{}) *FileOperator_Open_Call {
	return &FileOperator_Open_Call{Call: _e.mock.On("Open", name)}
}

func (_c *FileOperator_Open_Call) Run(run func(name string)) *FileOperator_Open_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *FileOperator_Open_Call) Return(_a0 io.ReadCloser, _a1 error) *FileOperator_Open_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Stat provides a mock function with given fields: name
func (_m *FileOperator) Stat(name string) (fs.FileInfo, error) {
	ret := _m.Called(name)

	var r0 fs.FileInfo
	if rf, ok := ret.Get(0).(func(string) fs.FileInfo); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fs.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FileOperator_Stat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stat'
type FileOperator_Stat_Call struct {
	*mock.Call
}

// Stat is a helper method to define mock.On call
//   - name string
func (_e *FileOperator_Expecter) Stat(name interface{}) *FileOperator_Stat_Call {
	return &FileOperator_Stat_Call{Call: _e.mock.On("Stat", name)}
}

func (_c *FileOperator_Stat_Call) Run(run func(name string)) *FileOperator_Stat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *FileOperator_Stat_Call) Return(_a0 fs.FileInfo, _a1 error) *FileOperator_Stat_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewFileOperator interface {
	mock.TestingT
	Cleanup(func())
}

// NewFileOperator creates a new instance of FileOperator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFileOperator(t mockConstructorTestingTNewFileOperator) *FileOperator {
	mock := &FileOperator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
