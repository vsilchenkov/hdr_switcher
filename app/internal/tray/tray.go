package tray

import (
	"fmt"
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/notify"
	"hdr_switcher/assets"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/getlantern/systray"

	gHotkey "golang.design/x/hotkey"
	//https://github.com/energye/systray
)

const (
	titleTray         = "HDR Toggle"
	hotKeySwitch_Name = "Ctrl+F12"
)

type menuItems struct {
	toggle     *systray.MenuItem
	status     *systray.MenuItem
	openFolder *systray.MenuItem
	quit       *systray.MenuItem
}

var hk *gHotkey.Hotkey

func Run() {

	slog.Info("Запуск приложения в системном трее...")
	systray.Run(onReady, onExit)
}

func Quit() {
	systray.Quit()
}

func onReady() {

	systray.SetIcon(assets.IconHDROff)

	items := menuItems{}

	systray.SetTitle(titleTray)
	systray.SetTooltip("Ctrl+F12: Toggle HDR")

	items.toggle = systray.AddMenuItem("Toggle HDR (Ctrl+F12)", "Переключить HDR")
	items.status = systray.AddMenuItem("Show status", "Показать состояние HDR")
	systray.AddSeparator()

	items.openFolder = systray.AddMenuItem("Open app folder", "Открыть папку приложения")
	systray.AddSeparator()
	items.quit = systray.AddMenuItem("Quit", "Выход")

	registerHotKey()

	// Обновление UI
	updateUI(items.toggle)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateUI(items.toggle)
		}
	}()

	// Обработчики пунктов меню
	go events(&items)

}

func registerHotKey() {

	hk = gHotkey.New(
		[]gHotkey.Modifier{gHotkey.ModCtrl}, gHotkey.KeyF12,
	)

	err := hk.Register()
	if err != nil {
		i := "Не удалось зарегистрировать хоткей"
		slog.Error(i, logging.Err(err))
		notify.ShowBalloon("HDR Toggle", i)
	} else {
		slog.Debug("Хоткей успешно зарегистрирован", slog.String("Hotkey", hotKeySwitch_Name))
	}
}

func onExit() {
	cleanup()
}

func cleanup() {
	// Снимаем регистрацию хоткея при выходе
	if hk != nil {
		err := hk.Unregister()
		if err != nil {
			slog.Error("Не удалось отменить регистрацию хоткея", logging.Err(err))
		}
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
