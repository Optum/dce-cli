// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// FileSystemer is an autogenerated mock type for the FileSystemer type
type FileSystemer struct {
	mock.Mock
}

// GetDefaultConfigFile provides a mock function with given fields:
func (_m *FileSystemer) GetDefaultConfigFile() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetHomeDir provides a mock function with given fields:
func (_m *FileSystemer) GetHomeDir() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsExistingFile provides a mock function with given fields: path
func (_m *FileSystemer) IsExistingFile(path string) bool {
	ret := _m.Called(path)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ReadFromFile provides a mock function with given fields: path
func (_m *FileSystemer) ReadFromFile(path string) string {
	ret := _m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// WriteToYAMLFile provides a mock function with given fields: path, _struct
func (_m *FileSystemer) WriteToYAMLFile(path string, _struct interface{}) {
	_m.Called(path, _struct)
}
