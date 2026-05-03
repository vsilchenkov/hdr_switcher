package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v2"
	gHotkey "golang.design/x/hotkey"
)

type HotkeyConfig struct {
	Modifiers []string `yaml:"Modifiers"`
	Key       string   `yaml:"Key"`
}

type Config struct {
	Hotkey HotkeyConfig `yaml:"Hotkey"`
	Debug  bool
}

var modifierMap = map[string]gHotkey.Modifier{
	"Ctrl":  gHotkey.ModCtrl,
	"Alt":   gHotkey.ModAlt,
	"Shift": gHotkey.ModShift,
	"Win":   gHotkey.ModWin,
}

var keyMap = map[string]gHotkey.Key{
	"F1":  gHotkey.KeyF1,
	"F2":  gHotkey.KeyF2,
	"F3":  gHotkey.KeyF3,
	"F4":  gHotkey.KeyF4,
	"F5":  gHotkey.KeyF5,
	"F6":  gHotkey.KeyF6,
	"F7":  gHotkey.KeyF7,
	"F8":  gHotkey.KeyF8,
	"F9":  gHotkey.KeyF9,
	"F11": gHotkey.KeyF11,
	"F12": gHotkey.KeyF12,
}

func LoadConfig() (*Config, error) {

	cfg := &Config{}
	parseFlags(cfg)

	path := "config.yml"

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func parseFlags(cfg *Config) {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Use debug")
	flag.Parse()
	cfg.Debug = debug
}

func (h HotkeyConfig) Parse() ([]gHotkey.Modifier, gHotkey.Key, error) {
	var mods []gHotkey.Modifier
	for _, m := range h.Modifiers {
		mod, ok := modifierMap[m]
		if !ok {
			return nil, 0, fmt.Errorf("неизвестный модификатор: %q", m)
		}
		mods = append(mods, mod)
	}

	key, ok := keyMap[h.Key]
	if !ok {
		return nil, 0, fmt.Errorf("неизвестная клавиша: %q", h.Key)
	}

	return mods, key, nil
}

func (h HotkeyConfig) String() string {
	var result strings.Builder
	for i, m := range h.Modifiers {
		if i > 0 {
			result.WriteString("+")
		}
		result.WriteString(m)
	}
	return result.String() + "+" + h.Key
}
