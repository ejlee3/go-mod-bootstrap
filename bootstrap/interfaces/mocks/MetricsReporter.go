// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	metrics "github.com/rcrowley/go-metrics"
	mock "github.com/stretchr/testify/mock"
)

// MetricsReporter is an autogenerated mock type for the MetricsReporter type
type MetricsReporter struct {
	mock.Mock
}

// Report provides a mock function with given fields: registry, metricTags
func (_m *MetricsReporter) Report(registry metrics.Registry, metricTags map[string]map[string]string) error {
	ret := _m.Called(registry, metricTags)

	var r0 error
	if rf, ok := ret.Get(0).(func(metrics.Registry, map[string]map[string]string) error); ok {
		r0 = rf(registry, metricTags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}