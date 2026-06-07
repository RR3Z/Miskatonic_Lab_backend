package tests

import (
	"context"
	"errors"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FakeMagicDBTX struct {
	ExecCalls     int
	QueryCalls    int
	QueryRowCalls int

	LastExecSQL      string
	LastExecArgs     []any
	LastQuerySQL     string
	LastQueryArgs    []any
	LastQueryRowSQL  string
	LastQueryRowArgs []any

	ExecErr       error
	QueryErr      error
	QueryRowErr   error
	QueryRowsData [][]any
	QueryRowData  []any

	QueryRowResults []FakeMagicQueryRowResult
}

type FakeMagicQueryRowResult struct {
	Err  error
	Data []any
}

func (f *FakeMagicDBTX) Exec(_ context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	f.ExecCalls++
	f.LastExecSQL = sql
	f.LastExecArgs = args

	return pgconn.CommandTag{}, f.ExecErr
}

func (f *FakeMagicDBTX) Query(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	f.QueryCalls++
	f.LastQuerySQL = sql
	f.LastQueryArgs = args

	if f.QueryErr != nil {
		return nil, f.QueryErr
	}

	return &fakeMagicRows{rows: f.QueryRowsData}, nil
}

func (f *FakeMagicDBTX) QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row {
	f.QueryRowCalls++
	f.LastQueryRowSQL = sql
	f.LastQueryRowArgs = args

	if len(f.QueryRowResults) >= f.QueryRowCalls {
		result := f.QueryRowResults[f.QueryRowCalls-1]
		return fakeMagicRow{data: result.Data, err: result.Err}
	}

	return fakeMagicRow{data: f.QueryRowData, err: f.QueryRowErr}
}

type fakeMagicRow struct {
	data []any
	err  error
}

func (r fakeMagicRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	scanMagicValues(r.data, dest)
	return nil
}

type fakeMagicRows struct {
	rows  [][]any
	index int
}

func (r *fakeMagicRows) Close() {}

func (r *fakeMagicRows) Err() error {
	return nil
}

func (r *fakeMagicRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r *fakeMagicRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *fakeMagicRows) Next() bool {
	if r.index >= len(r.rows) {
		return false
	}

	r.index++
	return true
}

func (r *fakeMagicRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.rows) {
		return errors.New("Scan called without current magic row")
	}

	scanMagicValues(r.rows[r.index-1], dest)
	return nil
}

func (r *fakeMagicRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.rows) {
		return nil, errors.New("Values called without current magic row")
	}

	return r.rows[r.index-1], nil
}

func (r *fakeMagicRows) RawValues() [][]byte {
	return nil
}

func (r *fakeMagicRows) Conn() *pgx.Conn {
	return nil
}

func scanMagicValues(data []any, dest []any) {
	for i, value := range data {
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
}
