package autorun

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const pathAutorun = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`

func ToggleToStartup(appName string) (bool, error) {

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

	isStartup, _ := IsStartup(appName)

	if !isStartup {
		err = k.SetStringValue(appName, execPath)
	} else {
		err = k.DeleteValue(appName)
	}

	if err != nil {
		return false, err
	}

	return !isStartup, nil
}

func IsStartup(appName string) (bool, error) {

	k, err := registry.OpenKey(registry.CURRENT_USER, pathAutorun, registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer k.Close()

	_, _, err = k.GetStringValue(appName)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
