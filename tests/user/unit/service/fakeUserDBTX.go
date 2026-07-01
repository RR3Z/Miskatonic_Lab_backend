package tests

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FakeUserDBTX struct {
	ExecCalls    int
	LastExecArgs []any
	ExecErr      error

	QueryRowCalls    int
	LastQueryRowArgs []any
	QueryRowData     []any
	QueryRowResults  []FakeUserQueryRowResult
}

type FakeUserQueryRowResult struct {
	Err  error
	Data []any
}

func (f *FakeUserDBTX) Exec(_ context.Context, _ string, args ...interface{}) (pgconn.CommandTag, error) {
	f.ExecCalls++
	f.LastExecArgs = args
	if f.ExecErr != nil {
		return pgconn.CommandTag{}, f.ExecErr
	}
	return pgconn.CommandTag{}, nil
}

func (f *FakeUserDBTX) Query(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (f *FakeUserDBTX) QueryRow(_ context.Context, _ string, args ...interface{}) pgx.Row {
	f.QueryRowCalls++
	f.LastQueryRowArgs = args

	if len(f.QueryRowResults) >= f.QueryRowCalls {
		result := f.QueryRowResults[f.QueryRowCalls-1]
		return fakeUserRow{data: result.Data, err: result.Err}
	}

	return fakeUserRow{data: f.QueryRowData}
}

type fakeUserRow struct {
	data []any
	err  error
}

func (r fakeUserRow) Scan(dest ...any) error {
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
