package hotkey

import (
	"fmt"
	"hdr_switcher/app/internal/notify"
	"syscall"

	"golang.org/x/sys/windows"
)

var (

	// Library
	libuser32 = windows.NewLazySystemDLL("user32.dll")

	registerHotKey   = libuser32.NewProc("RegisterHotKey")
	unregisterHotKey = libuser32.NewProc("UnregisterHotKey")
)

func RegisterHotKey(hwnd windows.HWND, id int, fsModifiers, vk uint) bool {

	// runtime.LockOSThread()
	// defer runtime.UnlockOSThread()

	// runtime.LockOSThread()

	ret, _, e1 := syscall.SyscallN(registerHotKey.Addr(),
		4,
		uintptr(hwnd),
		uintptr(id),
		uintptr(fsModifiers),
		uintptr(vk),
		0,
		0)

	if ret == 0 {
		// e1 — это already GetLastError; он может быть syscall.Errno(1409) и т.п.
		notify.Send("HDR Toggle", fmt.Sprintf("Не удалось зарегистрировать хоткей: %v", e1))
		// Не падаем сразу — попробуем альтернативу
	}
	return ret != 0
}

func UnregisterHotKey(hwnd windows.HWND, id int) bool {
	ret, _, _ := syscall.SyscallN(unregisterHotKey.Addr(),
		2,
		uintptr(hwnd),
		uintptr(id),
		0)

	return ret != 0
}
