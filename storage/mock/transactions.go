// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	storage "github.com/onflow/flow-go/storage"
)

// Transactions is an autogenerated mock type for the Transactions type
type Transactions struct {
	mock.Mock
}

// BatchStore provides a mock function with given fields: tx, batch
func (_m *Transactions) BatchStore(tx *flow.TransactionBody, batch storage.ReaderBatchWriter) error {
	ret := _m.Called(tx, batch)

	if len(ret) == 0 {
		panic("no return value specified for BatchStore")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.TransactionBody, storage.ReaderBatchWriter) error); ok {
		r0 = rf(tx, batch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ByID provides a mock function with given fields: txID
func (_m *Transactions) ByID(txID flow.Identifier) (*flow.TransactionBody, error) {
	ret := _m.Called(txID)

	if len(ret) == 0 {
		panic("no return value specified for ByID")
	}

	var r0 *flow.TransactionBody
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (*flow.TransactionBody, error)); ok {
		return rf(txID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.TransactionBody); ok {
		r0 = rf(txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.TransactionBody)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: tx
func (_m *Transactions) Store(tx *flow.TransactionBody) error {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for Store")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.TransactionBody) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransactions creates a new instance of Transactions. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactions(t interface {
	mock.TestingT
	Cleanup(func())
}) *Transactions {
	mock := &Transactions{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
