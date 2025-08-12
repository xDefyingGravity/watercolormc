package internal

import (
	"github.com/gen2brain/beeep"
)

func Notify(title string, message string) error {
	return beeep.Notify(title, message, "")
}
