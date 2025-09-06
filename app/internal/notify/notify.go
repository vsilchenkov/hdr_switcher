package notify

import (
	"github.com/gen2brain/beeep"
)

func init() {
	beeep.AppName = "HDR switcher"
}

func ShowBalloon(title, text string) error {
	err := beeep.Notify(title, text, "info.png")
	return err
}
