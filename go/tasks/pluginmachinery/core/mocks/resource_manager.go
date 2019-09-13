// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import core "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
import mock "github.com/stretchr/testify/mock"

// ResourceManager is an autogenerated mock type for the ResourceManager type
type ResourceManager struct {
	mock.Mock
}

// AllocateResource provides a mock function with given fields: ctx, namespace, allocationToken
func (_m *ResourceManager) AllocateResource(ctx context.Context, namespace string, allocationToken string) (core.AllocationStatus, error) {
	ret := _m.Called(ctx, namespace, allocationToken)

	var r0 core.AllocationStatus
	if rf, ok := ret.Get(0).(func(context.Context, string, string) core.AllocationStatus); ok {
		r0 = rf(ctx, namespace, allocationToken)
	} else {
		r0 = ret.Get(0).(core.AllocationStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, namespace, allocationToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReleaseResource provides a mock function with given fields: ctx, namespace, allocationToken
func (_m *ResourceManager) ReleaseResource(ctx context.Context, namespace string, allocationToken string) error {
	ret := _m.Called(ctx, namespace, allocationToken)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, namespace, allocationToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}