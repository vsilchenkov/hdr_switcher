package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	OutputInFile bool
}

func Init(c Config) {

	var output io.Writer

	if !c.OutputInFile {
		output = os.Stdout
	} else {

		exe, err := os.Executable()
		if err != nil {
			fmt.Printf("os.Executable failed: %v", err)
		}
		dir := filepath.Dir(exe)
		logPath := filepath.Join(dir, "app.log")

		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("cannot open log file: %v", err)
			output = os.Stdout
		} else {
			output = f
		}
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// Преобразуем время в строку с нужным форматом
				t := a.Value.Time()
				formatted := t.Format("2006-01-02 15:04:05")
				return slog.String(a.Key, formatted)
			}
			return a
		},
	})
	logger := slog.New(handler)

	slog.SetDefault(logger)

}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
