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

	"github.com/energye/systray"

	gHotkey "golang.design/x/hotkey"
)

const (
	titleTray         = "HDR Toggle"
	hotKeySwitch_Name = "Ctrl+Alt+F12"
)

type menuItems struct {
	toggle     *systray.MenuItem
	status     *systray.MenuItem
	openFolder *systray.MenuItem
	quit       *systray.MenuItem
}

var hk *gHotkey.Hotkey

func Run() {

	// notify.ShowBalloon("","Launching the application in the system tray")
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
	systray.SetTooltip(fmt.Sprintf("%s: Toggle HDR", hotKeySwitch_Name))

	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})

	systray.CreateMenu()
	items.toggle = systray.AddMenuItem(fmt.Sprintf("Toggle HDR (%s)", hotKeySwitch_Name), "Переключить HDR")
	items.toggle.Click(func() { onClicktoggle(items) })
	items.status = systray.AddMenuItem("Show status", "Показать состояние HDR")
	items.status.Click(func ()  { onClickShowStatus(items)	})

	systray.AddSeparator()
	items.openFolder = systray.AddMenuItem("Open app folder", "Открыть папку приложения")
	items.openFolder.Click(openAppFolder)

	systray.AddSeparator()
	items.quit = systray.AddMenuItem("Quit", "Выход")
	items.quit.Click(func() { systray.Quit() })

	systray.SetOnClick(func(menu systray.IMenu) {
		onClicktoggle(items)
	})

	registerHotKey(items)

	// Обновление UI
	updateUI(items, "")
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateUI(items, "")
		}
	}()

}

func registerHotKey(items menuItems) {

	hk = gHotkey.New(
		[]gHotkey.Modifier{gHotkey.ModCtrl, gHotkey.ModAlt}, gHotkey.KeyF12,
	)

	err := hk.Register()
	if err != nil {
		i := "Не удалось зарегистрировать хоткей"
		slog.Error(i, logging.Err(err))
		notify.ShowBalloon("HDR Toggle", i)
	} else {
		slog.Debug("Хоткей успешно зарегистрирован", slog.String("Hotkey", hotKeySwitch_Name))
	}

	go func() {
		for range hk.Keydown() {
			onClicktoggle(items)
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
			slog.Error("Не удалось отменить регистрацию хоткея", logging.Err(err))
		}
	}
}

func openAppFolder() {
	exe, err := os.Executable()
	if err != nil {
		slog.Error("open app folder", slog.Any("error", err))
		return
	}
	dir := filepath.Dir(exe)
	cmd := exec.Command("explorer.exe", dir)
	cmd.Start()
}
