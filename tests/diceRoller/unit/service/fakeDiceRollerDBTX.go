package tests

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FakeDiceRollerDBTX struct {
	QueryRowCalls    int
	LastQueryRowArgs []any
	QueryRowData     []any
	QueryRowResults  []FakeDiceRollerQueryRowResult

	QueryCalls    int
	LastQueryArgs []any
	QueryRows     [][]any
	QueryErr      error
	RowsScanErr   error

	ExecCalls    int
	LastExecArgs []any
	ExecErr      error
}

type FakeDiceRollerQueryRowResult struct {
	Err  error
	Data []any
}

func (f *FakeDiceRollerDBTX) Exec(_ context.Context, _ string, args ...interface{}) (pgconn.CommandTag, error) {
	f.ExecCalls++
	f.LastExecArgs = args
	if f.ExecErr != nil {
		return pgconn.CommandTag{}, f.ExecErr
	}
	return pgconn.CommandTag{}, nil
}

func (f *FakeDiceRollerDBTX) Query(_ context.Context, _ string, args ...interface{}) (pgx.Rows, error) {
	f.QueryCalls++
	f.LastQueryArgs = args
	if f.QueryErr != nil {
		return nil, f.QueryErr
	}
	return &fakeDiceRows{data: f.QueryRows, scanErr: f.RowsScanErr}, nil
}

func (f *FakeDiceRollerDBTX) QueryRow(_ context.Context, _ string, args ...interface{}) pgx.Row {
	f.QueryRowCalls++
	f.LastQueryRowArgs = args

	if len(f.QueryRowResults) >= f.QueryRowCalls {
		result := f.QueryRowResults[f.QueryRowCalls-1]
		return fakeDiceRow{data: result.Data, err: result.Err}
	}

	return fakeDiceRow{data: f.QueryRowData}
}

type fakeDiceRow struct {
	data []any
	err  error
}

func (r fakeDiceRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	for i, value := range r.data {
		if i >= len(dest) {
			break
		}

		target := reflect.ValueOf(dest[i])
		if target.Kind() != reflect.Pointer || target.IsNil() {
			continue
		}

		targetValue := target.Elem()
		sourceValue := reflect.ValueOf(value)
		if sourceValue.IsValid() && sourceValue.Type().AssignableTo(targetValue.Type()) {
			targetValue.Set(sourceValue)
		}
	}

	return nil
}

type fakeDiceRows struct {
	data    [][]any
	current int
	scanErr error
}

func (r *fakeDiceRows) Close()                                       {}
func (r *fakeDiceRows) Err() error                                   { return nil }
func (r *fakeDiceRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *fakeDiceRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *fakeDiceRows) Values() ([]any, error)                       { return nil, nil }
func (r *fakeDiceRows) RawValues() [][]byte                          { return nil }
func (r *fakeDiceRows) Conn() *pgx.Conn                              { return nil }

func (r *fakeDiceRows) Next() bool {
	if r.current >= len(r.data) {
		return false
	}
	r.current++
	return true
}

func (r *fakeDiceRows) Scan(dest ...any) error {
	if r.scanErr != nil {
		return r.scanErr
	}

	row := r.data[r.current-1]
	for i, value := range row {
		if i >= len(dest) {
			break
		}

		target := reflect.ValueOf(dest[i])
		if target.Kind() != reflect.Pointer || target.IsNil() {
			continue
		}

		targetValue := target.Elem()
		sourceValue := reflect.ValueOf(value)
		if sourceValue.IsValid() && sourceValue.Type().AssignableTo(targetValue.Type()) {
			targetValue.Set(sourceValue)
		}
	}

	return nil
}
