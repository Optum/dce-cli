// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import time "time"

// Durationer is an autogenerated mock type for the Durationer type
type Durationer struct {
	mock.Mock
}

// ExpandEpochTime provides a mock function with given fields: str
func (_m *Durationer) ExpandEpochTime(str string) (int64, error) {
	ret := _m.Called(str)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(str)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(str)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseDuration provides a mock function with given fields: str
func (_m *Durationer) ParseDuration(str string) (time.Duration, error) {
	ret := _m.Called(str)

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func(string) time.Duration); ok {
		r0 = rf(str)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(str)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
