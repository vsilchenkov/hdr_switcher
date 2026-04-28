package tray

import (
	"fmt"
	"hdr_switcher/app/internal/hdr"
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/notify"
	"hdr_switcher/assets"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/energye/systray"
)

func (m *menuItems) onClickShowStatus() {
	state, err := hdr.GetHDRState()
	if err != nil {
		notify.ShowBalloon("", fmt.Sprintf("Статус HDR: ошибка — %v", err))
	} else {
		notify.ShowBalloon("", fmt.Sprintf("Статус HDR: %s", state))
	}
	m.updateUI(state)
}

func (m *menuItems) onClicktoggle() {

	if err := hdr.ToggleHDR(); err != nil {
		notify.ShowBalloon("", fmt.Sprintf("Ошибка переключения HDR: %v", err))
	} else {
		m.updateUI("")
	}
}

func (m *menuItems) updateUI(state string) {
	if state == "" {
		s, err := hdr.GetHDRState()
		if err != nil {
			slog.Error("Не удалось получить статус HDR для обновления UI", logging.Err(err))
			return
		}
		state = s
	}
	if state == hdr.StateOn {
		m.toggle.Check()
		systray.SetIcon(assets.IconHDROn)
	} else {
		m.toggle.Uncheck()
		systray.SetIcon(assets.IconHDROff)
	}
}

func (m *menuItems) changeAutorun() {
	item := m.autorunItem
	var err error
	var msg string

	if item.Checked() {
		err = m.autorun.Disable()
		msg = "Autorun disabled"
	} else {
		err = m.autorun.Enable()
		msg = "Autorun enabled"
	}

	if err != nil {
		slog.Error("Ошибка изменения автозапуска", logging.Err(err))
		notify.ShowBalloon("", fmt.Sprintf("Error changed autorun: %v", err))
		return
	}

	notify.ShowBalloon("", msg)
	m.updateAutorunUI()
}

func (m *menuItems) updateAutorunUI() {

	enabled, err := m.autorun.IsStartup()
	if err != nil {
		slog.Error("Не удалось получить статус автозапуска для обновления UI",
			logging.Err(err))
		return
	}
	if enabled {
		m.autorunItem.Check()
	} else {
		m.autorunItem.Uncheck()
	}
}

func (m *menuItems) openAppFolder() {
	exe, err := os.Executable()
	if err != nil {
		slog.Error("open app folder", slog.Any("error", err))
		return
	}
	dir := filepath.Dir(exe)
	cmd := exec.Command("explorer.exe", dir)
	cmd.Start()
}
