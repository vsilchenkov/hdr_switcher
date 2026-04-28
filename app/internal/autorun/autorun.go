package autorun

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

type Autorun struct {
	ExeName string
}

const pathAutorun = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`

func (a Autorun) ToggleToStartup() (bool, error) {

	execPath, err := os.Executable()
	if err != nil {
		return false, err
	}
	execPath = filepath.Clean(execPath)

	k, err := registry.OpenKey(registry.CURRENT_USER, pathAutorun, registry.SET_VALUE)
	if err != nil {
		return false, err
	}
	defer k.Close()

	isStartup, _ := a.IsStartup()

	if !isStartup {
		err = k.SetStringValue(a.ExeName, execPath)
	} else {
		err = k.DeleteValue(a.ExeName)
	}

	if err != nil {
		return false, err
	}

	return !isStartup, nil
}

func (a Autorun) IsStartup() (bool, error) {

	k, err := registry.OpenKey(registry.CURRENT_USER, pathAutorun, registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer k.Close()

	_, _, err = k.GetStringValue(a.ExeName)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a Autorun) Enable() error {

	enabled, _ := a.IsStartup()
	if enabled {
		return nil
	}

	_, err := a.ToggleToStartup()
	if err != nil {
		return fmt.Errorf("enable autorun: %w", err)
	}
	return nil

}

func (a Autorun) Disable() error {

	enabled, _ := a.IsStartup()
	if !enabled {
		return nil
	}

	_, err := a.ToggleToStartup()
	if err != nil {
		return fmt.Errorf("disable autorun: %w", err)
	}
	return nil

}
