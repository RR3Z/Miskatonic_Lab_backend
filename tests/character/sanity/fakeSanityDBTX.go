package tests

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FakeSanityDBTX struct {
	QueryRowCalls    int
	LastQueryRowArgs []any
	QueryRowData     []any
	QueryRowResults  []FakeSanityQueryRowResult
}

type FakeSanityQueryRowResult struct {
	Err  error
	Data []any
}

func (f *FakeSanityDBTX) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (f *FakeSanityDBTX) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (f *FakeSanityDBTX) QueryRow(_ context.Context, _ string, args ...interface{}) pgx.Row {
	f.QueryRowCalls++
	f.LastQueryRowArgs = args

	if len(f.QueryRowResults) >= f.QueryRowCalls {
		result := f.QueryRowResults[f.QueryRowCalls-1]
		return fakeSanityRow{data: result.Data, err: result.Err}
	}

	return fakeSanityRow{data: f.QueryRowData}
}

type fakeSanityRow struct {
	data []any
	err  error
}

func (r fakeSanityRow) Scan(dest ...any) error {
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
