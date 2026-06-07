package tests

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FakeLuckDBTX struct {
	QueryRowCalls    int
	LastQueryRowArgs []any
	QueryRowData     []any
	QueryRowResults  []FakeLuckQueryRowResult
}

type FakeLuckQueryRowResult struct {
	Err  error
	Data []any
}

func (f *FakeLuckDBTX) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (f *FakeLuckDBTX) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (f *FakeLuckDBTX) QueryRow(_ context.Context, _ string, args ...interface{}) pgx.Row {
	f.QueryRowCalls++
	f.LastQueryRowArgs = args

	if len(f.QueryRowResults) >= f.QueryRowCalls {
		result := f.QueryRowResults[f.QueryRowCalls-1]
		return fakeLuckRow{data: result.Data, err: result.Err}
	}

	return fakeLuckRow{data: f.QueryRowData}
}

type fakeLuckRow struct {
	data []any
	err  error
}

func (r fakeLuckRow) Scan(dest ...any) error {
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
