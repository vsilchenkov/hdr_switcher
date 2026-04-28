package tray

import (
	"fmt"
	"hdr_switcher/app/internal/app"
	"hdr_switcher/app/internal/autorun"
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/notify"
	"hdr_switcher/assets"
	"log/slog"
	"sync"
	"time"

	"github.com/energye/systray"

	gHotkey "golang.design/x/hotkey"
)

const (
	titleTray = "HDR Toggle"
)

type menuItems struct {
	toggle      *systray.MenuItem
	status      *systray.MenuItem
	autorunItem *systray.MenuItem
	openFolder  *systray.MenuItem
	quit        *systray.MenuItem

	app     *app.App
	autorun *autorun.Autorun
	hk      *hotKey
}

type hotKey struct {
	*gHotkey.Hotkey
	sync.Mutex
}

func Run(app *app.App, autorun *autorun.Autorun) {

	hk := &hotKey{}

	slog.Info("Запуск приложения в системном трее...")
	systray.Run(func() { onReady(app, hk, autorun) }, func() { onExit(hk) })
}

func Quit() {
	systray.Quit()
}

func onReady(app *app.App, hk *hotKey, autorun *autorun.Autorun) {

	systray.SetIcon(assets.IconHDROff)

	items := menuItems{
		app:     app,
		autorun: autorun,
		hk:      hk,
	}

	hotKeySwitchName := app.Config.Hotkey.String()

	systray.SetTitle(titleTray)
	systray.SetTooltip(fmt.Sprintf("%s: Toggle HDR", hotKeySwitchName))

	items.toggle = systray.AddMenuItem(fmt.Sprintf("Toggle HDR (%s)", hotKeySwitchName), "Переключить HDR")
	items.toggle.Click(func() { items.onClicktoggle() })
	items.status = systray.AddMenuItem("Show status", "Показать состояние HDR")
	items.status.Click(func() { items.onClickShowStatus() })

	systray.AddSeparator()
	items.openFolder = systray.AddMenuItem("Open app folder", "Открыть папку приложения")
	items.openFolder.Click(items.openAppFolder)

	items.autorunItem = systray.AddMenuItem("Autorun", "Автозагрузка")
	items.autorunItem.Click(items.changeAutorun)
	items.updateAutorunUI()

	systray.AddSeparator()
	items.quit = systray.AddMenuItem("Quit", "Выход")
	items.quit.Click(func() { systray.Quit() })

	systray.SetOnClick(func(menu systray.IMenu) {
		items.onClicktoggle()
	})

	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})

	items.registerHotKey()

	// Обновление UI
	items.updateUI("")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in UI update loop",
					slog.Any("panic", r))
			}
		}()
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			items.updateUI("")
		}
	}()

}

func (m *menuItems) registerHotKey() {
	m.hk.Lock()
	defer m.hk.Unlock()

	if m.hk.Hotkey != nil {
		_ = m.hk.Hotkey.Unregister()
		m.hk.Hotkey = nil
	}

	mods, key, err := m.app.Config.Hotkey.Parse()
	if err != nil {
		slog.Error("Некорректная конфигурация хоткея", logging.Err(err))
		notify.ShowBalloon("HDR Toggle", "Error in hotkey configuration")
		return
	}

	m.hk.Hotkey = gHotkey.New(mods, key)

	err = m.hk.Hotkey.Register()
	if err != nil {
		i := "Не удалось зарегистрировать хоткей " + m.app.Config.Hotkey.String()
		slog.Error(i, logging.Err(err))
		notify.ShowBalloon("HDR Toggle", i)
		m.hk.Hotkey = nil
		return
	}

	slog.Debug("Хоткей успешно зарегистрирован",
		slog.String("Hotkey", m.app.Config.Hotkey.String()))

	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in hotkey listener",
					slog.Any("panic", r))
				notify.ShowBalloon("HDR Toggle", "An error occurred in hotkey listener")
			}
		}()
		for range m.hk.Keydown() {
			m.onClicktoggle()
		}
	}()
}

func onExit(hk *hotKey) {
	cleanup(hk)
}

func cleanup(hk *hotKey) {
	hk.Lock()
	defer hk.Unlock()

	// Снимаем регистрацию хоткея при выходе
	if hk.Hotkey != nil {
		err := hk.Hotkey.Unregister()
		if err != nil {
			slog.Error("Не удалось отменить регистрацию хоткея",
				logging.Err(err))
		}
		hk.Hotkey = nil
	}
}
