package tray

import (
	"fmt"
	"hdr_switcher/app/internal/hdr"
	"hdr_switcher/app/internal/notify"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"	

	"github.com/getlantern/systray"

	gHotkey "golang.design/x/hotkey"
)

const (
	// CTRL + F12
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	VK_G        = 0x47
	VK_F12      = 0x7B
	VK_F11      = 0x7A
	VK_F10      = 0x79
	VK_F9       = 0x78

	HOTKEY_ID = 1
)

var hk *gHotkey.Hotkey

func Run() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Запускаем горутину, которая будет ждать сигнала.
	go func() {
		sig := <-sigs
		log.Printf("Получен сигнал: %s. Запускаем процедуру выхода.", sig)
		// Инициируем корректное завершение работы systray.
		// Это приведет к вызову onExit.
		systray.Quit()
	}()

	log.Println("Запуск приложения в системном трее...")
	// Запуск блокирующего цикла systray.
	// onExit будет вызван, когда systray.Quit() сработает.
	systray.Run(onReady, onExit)
	log.Println("Приложение завершило работу.")
}

func onReady() {

	// Устанавливаем иконку сразу после запуска
	// systray.SetIcon(iconData) // <-- ВОТ ЭТА СТРОКА

	// runtime.LockOSThread()

	systray.SetTitle("HDR Toggle")
	systray.SetTooltip("Ctrl+F12: Toggle HDR")
	//systray.SetTooltip("Ctrl+Alt+F9: Toggle HDR")

	mToggle := systray.AddMenuItem("Toggle HDR (Ctrl+F12)", "Переключить HDR")
	// mToggle := systray.AddMenuItem("Toggle HDR (Ctrl+Alt+F9)", "Переключить HDR")

	mOn := systray.AddMenuItem("Force ON", "Включить HDR")
	mOff := systray.AddMenuItem("Force OFF", "Выключить HDR")
	systray.AddSeparator()
	mStatus := systray.AddMenuItem("Show status", "Показать состояние HDR")
	mOpenFolder := systray.AddMenuItem("Open app folder", "Открыть папку приложения")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Выход")

	// Регистрируем глобальный хоткей
	// mods := MOD_CONTROL | MOD_ALT
	// if !hotkey.RegisterHotKey(0, HOTKEY_ID, uint(mods), VK_F9) {
	// 	// err := windows.GetLastError()
	// 	// notify.Send("HDR Toggle", fmt.Sprintf("Не удалось зарегистрировать хоткей: %v", err))
	// 	log.Fatal("RegisterHotKey failed") // часто ERROR_HOTKEY_ALREADY_REGISTERED
	// } else {
	// 	go messageLoop()
	// }

	hk = gHotkey.New(
		//[]gHotkey.Modifier{gHotkey.ModCtrl, gHotkey.ModAlt},
		//gHotkey.KeyF9,
		[]gHotkey.Modifier{gHotkey.ModCtrl}, 	gHotkey.KeyF12,
	)

	err := hk.Register()
	if err != nil {
		log.Printf("Не удалось зарегистрировать хоткей: %v", err)
		notify.Send("HDR Toggle", fmt.Sprintf("Не удалось зарегистрировать хоткей: %v", err))
	} else {
		log.Println("Хоткей Ctrl+Alt+F9 успешно зарегистрирован.")
	}

	// Обработчики пунктов меню
	go func() {
		for {
			select {
			case <-hk.Keydown():
				if err := hdr.ToggleHDR(); err != nil {
					notify.Send("HDR Toggle", fmt.Sprintf("Ошибка переключения: %v", err))
				}
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
	cleanup()	
}

func cleanup() {
	// Снимаем регистрацию хоткея при выходе
	if hk != nil {
		err := hk.Unregister()
		if err != nil {
			log.Printf("Не удалось отменить регистрацию хоткея: %v", err)
		}
	}
}

// Цикл обработки сообщений, ловим WM_HOTKEY.
// func messageLoop() {
// 	var msg win.MSG
// 	for {
// 		ret := win.GetMessage(&msg, 0, 0, 0)
// 		if ret == 0 || ret == -1 {
// 			break
// 		}
// 		if msg.Message == win.WM_HOTKEY {
// 			if id := int(msg.WParam); id == HOTKEY_ID {
// 				if err := hdr.ToggleHDR(); err != nil {
// 					notify.Send("HDR Toggle", fmt.Sprintf("Ошибка переключения: %v", err))
// 				}
// 			}
// 		}
// 		win.TranslateMessage(&msg)
// 		win.DispatchMessage(&msg)
// 	}
// }

func openAppFolder() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("не удалось получить путь к exe: %w", err)
	}
	dir := filepath.Dir(exe)
	cmd := exec.Command("explorer.exe", dir)
	return cmd.Start()
}
