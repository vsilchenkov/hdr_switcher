//go:build windows
// +build windows

package main

import (
	_ "embed"
	build "hdr_switcher/app/buld.go"
	"hdr_switcher/app/internal/app"
	"hdr_switcher/app/internal/autorun"
	"hdr_switcher/app/internal/config"
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/notify"
	"hdr_switcher/app/internal/tray"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	appName    = "HDR switcher"
	appExeName = "hdr_switcher"
)

func main() {

	logging.Init(logging.Config{
		OutputInFile: true})

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			slog.Error("Panic recovered",
				slog.Any("panic", r),
				slog.String("stack", string(stack)),
			)
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		notify.ShowBalloon("HDR Toggle", "Error loading configuration")
		slog.Error("Ошибка загрузки конфигурации",
			logging.Err(err))
		os.Exit(1)
	}

	app := &app.App{
		Name:    appName,
		ExeName: appExeName,
		Config:  cfg,
	}

	notify.Init(app)

	b := build.Build{
		Name: app.ExeName}

	single, handle, err := b.IsSingleInstance()
	if err != nil {
		slog.Error("Ошибка создания мьютекса", logging.Err(err))
		os.Exit(1)
	}

	if !single {
		slog.Info("Приложение уже запущено!")
		notify.ShowBalloon("HDR Toggle", "The app is already running")
		os.Exit(1)
	}

	defer windows.CloseHandle(handle)

	autorun := &autorun.Autorun{
		ExeName: app.ExeName,
	}

	tray.Run(app, autorun)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		tray.Quit()
	}()
}
