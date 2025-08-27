package logging

import (
	"log"
	"os"
	"path/filepath"
)

func Setup() {

	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	logPath := filepath.Join(dir, "app.log")

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(f)
	} else {
		log.Printf("cannot open log file: %v", err)
	}
	
}
