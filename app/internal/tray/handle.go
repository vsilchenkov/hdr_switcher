package tray

import (
	"fmt"
	"hdr_switcher/app/internal/hdr"
	"hdr_switcher/app/internal/logging"
	"hdr_switcher/app/internal/notify"
	"hdr_switcher/assets"
	"log/slog"

	"github.com/energye/systray"
)

func onClickShowStatus(items menuItems) {
	state, err := hdr.GetHDRState()
	if err != nil {
		notify.ShowBalloon("", fmt.Sprintf("Статус HDR: ошибка — %v", err))
	} else {
		notify.ShowBalloon("", fmt.Sprintf("Статус HDR: %s", state))
	}
	updateUI(items, state)
}

func onClicktoggle(items menuItems) {

	if err := hdr.ToggleHDR(); err != nil {
		notify.ShowBalloon("", fmt.Sprintf("Ошибка переключения HDR: %v", err))
	} else {
		updateUI(items, "")
	}
}

func updateUI(items menuItems, state string) {
	if state == "" {
		s, err := hdr.GetHDRState()
		if err != nil {
			slog.Error("Не удалось получить статус HDR для обновления UI", logging.Err(err))
			return
		}
		state = s
	}
	if state == hdr.StateOn {
		items.toggle.Check()
		systray.SetIcon(assets.IconHDROn)
	} else {
		items.toggle.Uncheck()
		systray.SetIcon(assets.IconHDROff)
	}
}
