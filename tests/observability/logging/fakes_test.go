package tests

import (
	"context"
	"errors"
	"log/slog"
)

var errEventLoggerTest = errors.New("event failed")

type recordedLog struct {
	level   slog.Level
	message string
	attrs   map[string]any
}

type recordingSlogHandler struct {
	records []recordedLog
}

func (h *recordingSlogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *recordingSlogHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]any)
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	})

	h.records = append(h.records, recordedLog{
		level:   record.Level,
		message: record.Message,
		attrs:   attrs,
	})
	return nil
}

func (h *recordingSlogHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *recordingSlogHandler) WithGroup(string) slog.Handler {
	return h
}
