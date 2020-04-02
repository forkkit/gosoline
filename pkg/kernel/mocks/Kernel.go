// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import kernel "github.com/applike/gosoline/pkg/kernel"
import mock "github.com/stretchr/testify/mock"

// Kernel is an autogenerated mock type for the Kernel type
type Kernel struct {
	mock.Mock
}

// Add provides a mock function with given fields: name, module, opts
func (_m *Kernel) Add(name string, module kernel.Module, opts ...kernel.ModuleOption) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name, module)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// AddFactory provides a mock function with given fields: factory
func (_m *Kernel) AddFactory(factory kernel.ModuleFactory) {
	_m.Called(factory)
}

// Run provides a mock function with given fields:
func (_m *Kernel) Run() {
	_m.Called()
}

// Stop provides a mock function with given fields: reason
func (_m *Kernel) Stop(reason string) {
	_m.Called(reason)
}
