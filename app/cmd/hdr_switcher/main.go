//go:build windows
// +build windows

package main

import (
	_ "embed"
	build "hdr_switcher/app/buld.go"
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/notify"
	"hdr_switcher/app/internal/tray"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/windows"
)

const AppName = "HDR switcher"

func main() {

	logging.Init(logging.Config{
		OutputInFile: false})
		
	notify.Init(AppName)

	single, handle, err := build.IsSingleInstance(AppName)
	if err != nil {
		slog.Error("Ошибка создания мьютекса", logging.Err(err))
		os.Exit(1)
	}

	if !single {
		slog.Info("Приложение уже запущено!")
		notify.ShowBalloon("","The app is already running")
		os.Exit(1)
	}

    defer windows.CloseHandle(handle)

	tray.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		tray.Quit()
	}()

}
