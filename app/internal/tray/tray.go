package tray

import (
	"fmt"
	"hdr_switcher/app/internal/hdr"
	"hdr_switcher/app/internal/hotkey"
	"hdr_switcher/app/internal/notify"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
)

const (
	// CTRL + F12
	MOD_CONTROL = 0x0002
	MOD_ALT     = 0x0001
	VK_G        = 0x47
	VK_F12      = 0x7B

	HOTKEY_ID = 1
)

func Run() {
	// Запуск системного трея и регистрации хоткея
	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetTitle("HDR Toggle")
	systray.SetTooltip("Ctrl+F12: Toggle HDR via HDRCmd")

	mToggle := systray.AddMenuItem("Toggle HDR (Ctrl+F12)", "Переключить HDR")
	mOn := systray.AddMenuItem("Force ON", "Включить HDR")
	mOff := systray.AddMenuItem("Force OFF", "Выключить HDR")
	systray.AddSeparator()
	mStatus := systray.AddMenuItem("Show status", "Показать состояние HDR")
	mOpenFolder := systray.AddMenuItem("Open app folder", "Открыть папку приложения")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Выход")

	// Регистрируем глобальный хоткей
	if err := hotkey.RegisterHotKey(win.HWND(0), HOTKEY_ID, MOD_CONTROL, VK_F12); err != nil {
		notify.Send("HDR Toggle", fmt.Sprintf("Не удалось зарегистрировать хоткей Ctrl+F12: %v", err))
		log.Printf("registerHotKey error: %v", err)
	} else {
		go messageLoop()
	}

	// Обработчики пунктов меню
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				if err := hdr.ToggleHDR(); err != nil {
					notify.Send("HDR Toggle", fmt.Sprintf("Ошибка переключения: %v", err))
				}
			case <-mOn.ClickedCh:
				if err := hdr.SetHDR(true); err != nil {
					notify.Send("HDR Toggle", fmt.Sprintf("Ошибка включения: %v", err))
				} else {
					notify.Send("HDR Toggle", "HDR включён")
				}
			case <-mOff.ClickedCh:
				if err := hdr.SetHDR(false); err != nil {
					notify.Send("HDR Toggle", fmt.Sprintf("Ошибка выключения: %v", err))
				} else {
					notify.Send("HDR Toggle", "HDR выключен")
				}
			case <-mStatus.ClickedCh:
				state, err := hdr.GetHDRState()
				if err != nil {
					notify.Send("HDR Toggle", fmt.Sprintf("Статус: ошибка — %v", err))
				} else {
					notify.Send("HDR Toggle", fmt.Sprintf("Статус: %s", state))
				}
			case <-mOpenFolder.ClickedCh:
				_ = openAppFolder()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Снимаем хоткей
	hotkey.UnregisterHotKeyRaw(uintptr(win.HWND(0)), int32(HOTKEY_ID))
}

// Цикл обработки сообщений, ловим WM_HOTKEY.
func messageLoop() {
	var msg win.MSG
	for {
		ret := win.GetMessage(&msg, 0, 0, 0)
		if ret == 0 || ret == -1 {
			break
		}
		if msg.Message == win.WM_HOTKEY {
			if id := int(msg.WParam); id == HOTKEY_ID {
				if err := hdr.ToggleHDR(); err != nil {
					notify.Send("HDR Toggle", fmt.Sprintf("Ошибка переключения: %v", err))
				}
			}
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}

func openAppFolder() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("не удалось получить путь к exe: %w", err)
	}
	dir := filepath.Dir(exe)
	cmd := exec.Command("explorer.exe", dir)
	return cmd.Start()
}
