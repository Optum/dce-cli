// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	c_o_r_s "github.com/Optum/dce-cli/client/c_o_r_s"
	mock "github.com/stretchr/testify/mock"

	runtime "github.com/go-openapi/runtime"
)

// ClientService is an autogenerated mock type for the ClientService type
type ClientService struct {
	mock.Mock
}

// OptionsAccounts provides a mock function with given fields: params
func (_m *ClientService) OptionsAccounts(params *c_o_r_s.OptionsAccountsParams) (*c_o_r_s.OptionsAccountsOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsAccountsOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsAccountsParams) *c_o_r_s.OptionsAccountsOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsAccountsOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsAccountsParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsAccountsID provides a mock function with given fields: params
func (_m *ClientService) OptionsAccountsID(params *c_o_r_s.OptionsAccountsIDParams) (*c_o_r_s.OptionsAccountsIDOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsAccountsIDOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsAccountsIDParams) *c_o_r_s.OptionsAccountsIDOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsAccountsIDOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsAccountsIDParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsAuth provides a mock function with given fields: params
func (_m *ClientService) OptionsAuth(params *c_o_r_s.OptionsAuthParams) (*c_o_r_s.OptionsAuthOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsAuthOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsAuthParams) *c_o_r_s.OptionsAuthOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsAuthOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsAuthParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsAuthFile provides a mock function with given fields: params
func (_m *ClientService) OptionsAuthFile(params *c_o_r_s.OptionsAuthFileParams) (*c_o_r_s.OptionsAuthFileOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsAuthFileOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsAuthFileParams) *c_o_r_s.OptionsAuthFileOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsAuthFileOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsAuthFileParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsLeases provides a mock function with given fields: params
func (_m *ClientService) OptionsLeases(params *c_o_r_s.OptionsLeasesParams) (*c_o_r_s.OptionsLeasesOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsLeasesOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsLeasesParams) *c_o_r_s.OptionsLeasesOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsLeasesOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsLeasesParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsLeasesAuth provides a mock function with given fields: params
func (_m *ClientService) OptionsLeasesAuth(params *c_o_r_s.OptionsLeasesAuthParams) (*c_o_r_s.OptionsLeasesAuthOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsLeasesAuthOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsLeasesAuthParams) *c_o_r_s.OptionsLeasesAuthOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsLeasesAuthOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsLeasesAuthParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsLeasesID provides a mock function with given fields: params
func (_m *ClientService) OptionsLeasesID(params *c_o_r_s.OptionsLeasesIDParams) (*c_o_r_s.OptionsLeasesIDOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsLeasesIDOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsLeasesIDParams) *c_o_r_s.OptionsLeasesIDOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsLeasesIDOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsLeasesIDParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsLeasesIDAuth provides a mock function with given fields: params
func (_m *ClientService) OptionsLeasesIDAuth(params *c_o_r_s.OptionsLeasesIDAuthParams) (*c_o_r_s.OptionsLeasesIDAuthOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsLeasesIDAuthOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsLeasesIDAuthParams) *c_o_r_s.OptionsLeasesIDAuthOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsLeasesIDAuthOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsLeasesIDAuthParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptionsUsage provides a mock function with given fields: params
func (_m *ClientService) OptionsUsage(params *c_o_r_s.OptionsUsageParams) (*c_o_r_s.OptionsUsageOK, error) {
	ret := _m.Called(params)

	var r0 *c_o_r_s.OptionsUsageOK
	if rf, ok := ret.Get(0).(func(*c_o_r_s.OptionsUsageParams) *c_o_r_s.OptionsUsageOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*c_o_r_s.OptionsUsageOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*c_o_r_s.OptionsUsageParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetTransport provides a mock function with given fields: transport
func (_m *ClientService) SetTransport(transport runtime.ClientTransport) {
	_m.Called(transport)
}
