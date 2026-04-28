package app

import "hdr_switcher/app/internal/config"

type App struct {
	Name    string
	ExeName string
	Config  *config.Config
}
