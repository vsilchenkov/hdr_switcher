package notify

import (
	"github.com/gen2brain/beeep"
)

func Init(appName string) {
	beeep.AppName = appName
}

const icon = "info.png"

func ShowBalloon(title, text string) error {
	err := beeep.Notify(title, text, icon)
	return err
}
