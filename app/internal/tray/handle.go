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

// func events(items *menuItems) {

// 	for {
// 		select {
// 		case <-hk.Keydown():
// 			onClicktoggle(items.toggle)
// 		case <-items.toggle.ClickedCh:
// 			onClicktoggle(items.toggle)
// 		case <-items.status.ClickedCh:

// 		case <-items.openFolder.ClickedCh:
// 			_ = openAppFolder()
// 		case <-items.quit:
// 			systray.Quit()
// 			return
// 		}
// 	}
// }
func onClickShowStatus() {
	state, err := hdr.GetHDRState()
	if err != nil {
		notify.ShowBalloon(titleTray, fmt.Sprintf("Статус: ошибка — %v", err))
	} else {
		notify.ShowBalloon(titleTray, fmt.Sprintf("Статус: %s", state))
	}
}
func onClicktoggle(item *systray.MenuItem) {

	if err := hdr.ToggleHDR(); err != nil {
		notify.ShowBalloon(titleTray, fmt.Sprintf("Ошибка переключения: %v", err))
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
