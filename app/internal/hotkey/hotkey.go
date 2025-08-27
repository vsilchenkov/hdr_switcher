package hotkey

import (
	"errors"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
)

// Регистрирует глобальный хоткей через WinAPI.
func RegisterHotKey(hWnd win.HWND, id int, fsModifiers uint, vk uint) error {
	if !registerHotKeyRaw(uintptr(hWnd), int32(id), uint32(fsModifiers), uint32(vk)) {
		return errors.New("RegisterHotKey failed")
	}
	return nil
}

func UnregisterHotKeyRaw(hwnd uintptr, id int32) bool {
	r1, _, _ := procUnregisterHotKey.Call(
		hwnd,
		uintptr(id),
	)
	return r1 != 0
}

func registerHotKeyRaw(hwnd uintptr, id int32, fsModifiers uint32, vk uint32) bool {
	r1, _, _ := procRegisterHotKey.Call(
		hwnd,
		uintptr(id),
		uintptr(fsModifiers),
		uintptr(vk),
	)
	return r1 != 0
}
