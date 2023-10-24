// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/julianinsua/the_simp_bank/internal/database (interfaces: Store)

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	database "github.com/julianinsua/the_simp_bank/internal/database"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddToAccountBalance mocks base method.
func (m *MockStore) AddToAccountBalance(arg0 context.Context, arg1 database.AddToAccountBalanceParams) (database.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToAccountBalance", arg0, arg1)
	ret0, _ := ret[0].(database.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddToAccountBalance indicates an expected call of AddToAccountBalance.
func (mr *MockStoreMockRecorder) AddToAccountBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToAccountBalance", reflect.TypeOf((*MockStore)(nil).AddToAccountBalance), arg0, arg1)
}

// CreateAccount mocks base method.
func (m *MockStore) CreateAccount(arg0 context.Context, arg1 database.CreateAccountParams) (database.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1)
	ret0, _ := ret[0].(database.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockStoreMockRecorder) CreateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockStore)(nil).CreateAccount), arg0, arg1)
}

// CreateEntry mocks base method.
func (m *MockStore) CreateEntry(arg0 context.Context, arg1 database.CreateEntryParams) (database.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntry", arg0, arg1)
	ret0, _ := ret[0].(database.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEntry indicates an expected call of CreateEntry.
func (mr *MockStoreMockRecorder) CreateEntry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntry", reflect.TypeOf((*MockStore)(nil).CreateEntry), arg0, arg1)
}

// CreateTransfer mocks base method.
func (m *MockStore) CreateTransfer(arg0 context.Context, arg1 database.CreateTransferParams) (database.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransfer", arg0, arg1)
	ret0, _ := ret[0].(database.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransfer indicates an expected call of CreateTransfer.
func (mr *MockStoreMockRecorder) CreateTransfer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransfer", reflect.TypeOf((*MockStore)(nil).CreateTransfer), arg0, arg1)
}

// DeleteAccount mocks base method.
func (m *MockStore) DeleteAccount(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockStoreMockRecorder) DeleteAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockStore)(nil).DeleteAccount), arg0, arg1)
}

// GetAccount mocks base method.
func (m *MockStore) GetAccount(arg0 context.Context, arg1 uuid.UUID) (database.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(database.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockStoreMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockStore)(nil).GetAccount), arg0, arg1)
}

// GetAccountEntries mocks base method.
func (m *MockStore) GetAccountEntries(arg0 context.Context, arg1 database.GetAccountEntriesParams) ([]database.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountEntries", arg0, arg1)
	ret0, _ := ret[0].([]database.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountEntries indicates an expected call of GetAccountEntries.
func (mr *MockStoreMockRecorder) GetAccountEntries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountEntries", reflect.TypeOf((*MockStore)(nil).GetAccountEntries), arg0, arg1)
}

// GetAccountForUpdate mocks base method.
func (m *MockStore) GetAccountForUpdate(arg0 context.Context, arg1 uuid.UUID) (database.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountForUpdate", arg0, arg1)
	ret0, _ := ret[0].(database.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountForUpdate indicates an expected call of GetAccountForUpdate.
func (mr *MockStoreMockRecorder) GetAccountForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountForUpdate", reflect.TypeOf((*MockStore)(nil).GetAccountForUpdate), arg0, arg1)
}

// GetAccountsList mocks base method.
func (m *MockStore) GetAccountsList(arg0 context.Context, arg1 database.GetAccountsListParams) ([]database.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountsList", arg0, arg1)
	ret0, _ := ret[0].([]database.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountsList indicates an expected call of GetAccountsList.
func (mr *MockStoreMockRecorder) GetAccountsList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountsList", reflect.TypeOf((*MockStore)(nil).GetAccountsList), arg0, arg1)
}

// GetEntry mocks base method.
func (m *MockStore) GetEntry(arg0 context.Context, arg1 uuid.UUID) (database.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntry", arg0, arg1)
	ret0, _ := ret[0].(database.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockStoreMockRecorder) GetEntry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockStore)(nil).GetEntry), arg0, arg1)
}

// GetTransfer mocks base method.
func (m *MockStore) GetTransfer(arg0 context.Context, arg1 uuid.UUID) (database.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransfer", arg0, arg1)
	ret0, _ := ret[0].(database.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransfer indicates an expected call of GetTransfer.
func (mr *MockStoreMockRecorder) GetTransfer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransfer", reflect.TypeOf((*MockStore)(nil).GetTransfer), arg0, arg1)
}

// TransferTx mocks base method.
func (m *MockStore) TransferTx(arg0 context.Context, arg1 database.TransferTxParams) (database.TransferTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferTx", arg0, arg1)
	ret0, _ := ret[0].(database.TransferTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransferTx indicates an expected call of TransferTx.
func (mr *MockStoreMockRecorder) TransferTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferTx", reflect.TypeOf((*MockStore)(nil).TransferTx), arg0, arg1)
}

// UpdateAccountBalance mocks base method.
func (m *MockStore) UpdateAccountBalance(arg0 context.Context, arg1 database.UpdateAccountBalanceParams) (database.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccountBalance", arg0, arg1)
	ret0, _ := ret[0].(database.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccountBalance indicates an expected call of UpdateAccountBalance.
func (mr *MockStoreMockRecorder) UpdateAccountBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccountBalance", reflect.TypeOf((*MockStore)(nil).UpdateAccountBalance), arg0, arg1)
}
