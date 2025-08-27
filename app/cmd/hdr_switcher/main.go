//go:build windows
// +build windows

package main

import (
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/tray"
)

func main() {

	logging.Setup()

	tray.Run()

}
