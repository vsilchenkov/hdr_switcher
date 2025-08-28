//go:build windows
// +build windows

package main

import (
	_ "embed" 
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/tray"
)

////go:embed icon.ico
//var iconData []byte 16x16

func main() {

	logging.Setup()

	tray.Run()

}
