package notify

import (
	"fmt"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// Простейшее уведомление через всплывающий balloon трея.
// systray сам по себе balloon не делает, поэтому сделаем через WinAPI Shell_NotifyIcon.
func Send(title, message string) {
	// Создаём временный balloon через NIF_INFO
	var nid win.NOTIFYICONDATA
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	// Для упрощения используем message-only окно systray (не строго необходимо)
	// В трее уже есть иконка systray, поэтому просто отправим инфо.
	copyUTF16ToFixed64(&nid.SzInfoTitle, title)
	copyUTF16ToFixed256(&nid.SzInfo, message)
	nid.UFlags = win.NIF_INFO
	// Тип: NIIF_INFO (информационное)
	nid.DwInfoFlags = win.NIIF_INFO
	// Вызовем Modify (или Add) — на практике Windows покажет balloon, если есть активный значок
	win.Shell_NotifyIcon(win.NIM_MODIFY, &nid)
	fmt.Printf("title:%s, message:%s\n", title, message)
}

func copyUTF16ToFixed64(dst *[64]uint16, s string) {
	u, _ := windows.UTF16FromString(s) // уже с завершающим 0

	// очистка буфера
	for i := range dst[:] {
		dst[i] = 0
	}

	n := len(u)
	if n > len(dst) {
		n = len(dst)
	}
	copy(dst[:], u[:n])
}


func copyUTF16ToFixed256(dst *[256]uint16, s string) {
	u, _ := windows.UTF16FromString(s)

	for i := range dst[:] {
		dst[i] = 0
	}

	n := len(u)
	if n > len(dst) {
		n = len(dst)
	}
	copy(dst[:], u[:n])
}
