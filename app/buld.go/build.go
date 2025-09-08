package build

import (
	"syscall"

	"golang.org/x/sys/windows"
)

func IsSingleInstance(appName string) (bool, windows.Handle, error) {

	mutexName := "Global\\630a2b47-6b7b-4928-9b12-88bc16312a93-" + appName

	// Создаём мьютекс с bInitialOwner = true
	mutex16, _ := syscall.UTF16PtrFromString(mutexName)
	h, err := windows.CreateMutex(nil, true, mutex16)
	if err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			// Мьютекс уже существует — второй экземпляр
			return false, 0, nil
		}
		return false, 0, err
	}

	// Успешно создан новый мьютекс — первый экземпляр
	return true, h, nil
}
