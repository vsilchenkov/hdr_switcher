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

func onClickShowStatus() {
	state, err := hdr.GetHDRState()
	if err != nil {
		notify.ShowBalloon("", fmt.Sprintf("Статус HDR: ошибка — %v", err))
	} else {
		notify.ShowBalloon("", fmt.Sprintf("Статус HDR: %s", state))
	}
}

func onClicktoggle(item *systray.MenuItem) {

	if err := hdr.ToggleHDR(); err != nil {
		notify.ShowBalloon("", fmt.Sprintf("Ошибка переключения HDR: %v", err))
	} else {
		updateUI(item)
	}
}

func updateUI(item *systray.MenuItem) {
	state, err := hdr.GetHDRState()
	if err != nil {
		slog.Error("Не удалось получить статус HDR для обновления UI", logging.Err(err))
		return
	}
	if state == hdr.StateOn {
		item.Check()
		systray.SetIcon(assets.IconHDROn)
	} else {
		item.Uncheck()
		systray.SetIcon(assets.IconHDROff)
	}
}
