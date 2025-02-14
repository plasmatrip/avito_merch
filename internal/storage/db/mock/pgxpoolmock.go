// Code generated by MockGen. DO NOT EDIT.
// Source: pgxpool.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

// MockPgxPool is a mock of PgxPool interface.
type MockPgxPool struct {
	ctrl     *gomock.Controller
	recorder *MockPgxPoolMockRecorder
}

// MockPgxPoolMockRecorder is the mock recorder for MockPgxPool.
type MockPgxPoolMockRecorder struct {
	mock *MockPgxPool
}

// NewMockPgxPool creates a new mock instance.
func NewMockPgxPool(ctrl *gomock.Controller) *MockPgxPool {
	mock := &MockPgxPool{ctrl: ctrl}
	mock.recorder = &MockPgxPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPgxPool) EXPECT() *MockPgxPoolMockRecorder {
	return m.recorder
}

// Acquire mocks base method.
func (m *MockPgxPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Acquire", ctx)
	ret0, _ := ret[0].(*pgxpool.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Acquire indicates an expected call of Acquire.
func (mr *MockPgxPoolMockRecorder) Acquire(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acquire", reflect.TypeOf((*MockPgxPool)(nil).Acquire), ctx)
}

// AcquireAllIdle mocks base method.
func (m *MockPgxPool) AcquireAllIdle(ctx context.Context) []*pgxpool.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcquireAllIdle", ctx)
	ret0, _ := ret[0].([]*pgxpool.Conn)
	return ret0
}

// AcquireAllIdle indicates an expected call of AcquireAllIdle.
func (mr *MockPgxPoolMockRecorder) AcquireAllIdle(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcquireAllIdle", reflect.TypeOf((*MockPgxPool)(nil).AcquireAllIdle), ctx)
}

// AcquireFunc mocks base method.
func (m *MockPgxPool) AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcquireFunc", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcquireFunc indicates an expected call of AcquireFunc.
func (mr *MockPgxPoolMockRecorder) AcquireFunc(ctx, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcquireFunc", reflect.TypeOf((*MockPgxPool)(nil).AcquireFunc), ctx, f)
}

// Begin mocks base method.
func (m *MockPgxPool) Begin(ctx context.Context) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin", ctx)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockPgxPoolMockRecorder) Begin(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockPgxPool)(nil).Begin), ctx)
}

// BeginTx mocks base method.
func (m *MockPgxPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, txOptions)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockPgxPoolMockRecorder) BeginTx(ctx, txOptions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockPgxPool)(nil).BeginTx), ctx, txOptions)
}

// Close mocks base method.
func (m *MockPgxPool) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockPgxPoolMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPgxPool)(nil).Close))
}

// Config mocks base method.
func (m *MockPgxPool) Config() *pgxpool.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*pgxpool.Config)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockPgxPoolMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockPgxPool)(nil).Config))
}

// CopyFrom mocks base method.
func (m *MockPgxPool) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyFrom", ctx, tableName, columnNames, rowSrc)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CopyFrom indicates an expected call of CopyFrom.
func (mr *MockPgxPoolMockRecorder) CopyFrom(ctx, tableName, columnNames, rowSrc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyFrom", reflect.TypeOf((*MockPgxPool)(nil).CopyFrom), ctx, tableName, columnNames, rowSrc)
}

// Exec mocks base method.
func (m *MockPgxPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range arguments {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockPgxPoolMockRecorder) Exec(ctx, sql interface{}, arguments ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, arguments...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockPgxPool)(nil).Exec), varargs...)
}

// Ping mocks base method.
func (m *MockPgxPool) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPgxPoolMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPgxPool)(nil).Ping), ctx)
}

// Query mocks base method.
func (m *MockPgxPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockPgxPoolMockRecorder) Query(ctx, sql interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockPgxPool)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockPgxPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockPgxPoolMockRecorder) QueryRow(ctx, sql interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockPgxPool)(nil).QueryRow), varargs...)
}

// Reset mocks base method.
func (m *MockPgxPool) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset.
func (mr *MockPgxPoolMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockPgxPool)(nil).Reset))
}

// SendBatch mocks base method.
func (m *MockPgxPool) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendBatch", ctx, b)
	ret0, _ := ret[0].(pgx.BatchResults)
	return ret0
}

// SendBatch indicates an expected call of SendBatch.
func (mr *MockPgxPoolMockRecorder) SendBatch(ctx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendBatch", reflect.TypeOf((*MockPgxPool)(nil).SendBatch), ctx, b)
}

// Stat mocks base method.
func (m *MockPgxPool) Stat() *pgxpool.Stat {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(*pgxpool.Stat)
	return ret0
}

// Stat indicates an expected call of Stat.
func (mr *MockPgxPoolMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockPgxPool)(nil).Stat))
}
