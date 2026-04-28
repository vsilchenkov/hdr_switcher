package notify

import (
	"hdr_switcher/app/internal/app"

	"github.com/gen2brain/beeep"
)

func Init(c *app.App) {
	beeep.AppName = c.Name
}

const icon = "info.png"

func ShowBalloon(title, text string) error {
	err := beeep.Notify(title, text, icon)
	return err
}
