//go:build windows
// +build windows

package main

import (
	_ "embed"
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/tray"
	"os"
	"os/signal"
	"syscall"
)

////go:embed icon.ico
//var iconData []byte 16x16

func main() {

	logging.Init(logging.Config{
		OutputInFile: false})

	tray.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		tray.Quit()
	}()

}
