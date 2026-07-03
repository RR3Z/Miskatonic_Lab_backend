package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func setupLogger() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: "15:04:05",
	}))

	slog.SetDefault(logger)
}
