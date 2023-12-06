package mock

import "github.com/dwethmar/lingo/pkg/database"

var _ database.RowHandler = &RowHandler{}
var _ database.RowsHandler = &RowsHandler{}

// RowHandler is a mock implementation of the database.RowHandler interface.
type RowHandler struct {
	NextFunc func() bool
	ErrFunc  func() error
	ScanFunc func(dest ...interface{}) error
}

func (m *RowHandler) Next() bool {
	if m.NextFunc == nil {
		panic("call to unimplemented method Next on RowHandler")
	}
	return m.NextFunc()
}

func (m *RowHandler) Err() error {
	if m.ErrFunc == nil {
		panic("call to unimplemented method Err on RowHandler")
	}
	return m.ErrFunc()
}

func (m *RowHandler) Scan(dest ...interface{}) error {
	if m.ScanFunc == nil {
		panic("call to unimplemented method Scan on RowHandler")
	}
	return m.ScanFunc(dest...)
}

// RowsHandler is a mock implementation of the database.RowsHandler interface.
type RowsHandler struct {
	CloseFunc func() error
	ErrFunc   func() error
	NextFunc  func() bool
	ScanFunc  func(dest ...interface{}) error
}

func (m *RowsHandler) Close() error {
	if m.CloseFunc == nil {
		panic("call to unimplemented method Close on RowsHandler")
	}
	return m.CloseFunc()
}

func (m *RowsHandler) Err() error {
	if m.ErrFunc == nil {
		panic("call to unimplemented method Err on RowsHandler")
	}
	return m.ErrFunc()
}

func (m *RowsHandler) Next() bool {
	if m.NextFunc == nil {
		panic("call to unimplemented method Next on RowsHandler")
	}
	return m.NextFunc()
}

func (m *RowsHandler) Scan(dest ...interface{}) error {
	if m.ScanFunc == nil {
		panic("call to unimplemented method Scan on RowsHandler")
	}
	return m.ScanFunc(dest...)
}
