package main

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

func initMRPISPlayer() (*mpris.Player, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	names, err := mpris.List(conn)
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return nil, errors.New("no player found")
	}

	return mpris.New(conn, names[0]), nil
}

func main() {
	a := app.NewWithID("io.github.jacalz.fympris")
	a.SetIcon(theme.MediaMusicIcon())
	w := a.NewWindow("Fympris")

	player, err := initMRPISPlayer()
	if err != nil {
		fyne.LogError("Failed to get player", err)
		return
	}

	controller := newMediaController(player)
	w.SetContent(controller.createUI())

	w.Resize(fyne.NewSize(400, 400))
	w.ShowAndRun()
}
