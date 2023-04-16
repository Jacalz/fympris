package main

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

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
	a := app.New()
	w := a.NewWindow("MPRIS")

	player, err := initMRPISPlayer()
	if err != nil {
		fyne.LogError("Failed to get player", err)
		return
	}

	previous := &widget.Button{Icon: theme.MediaSkipPreviousIcon(), OnTapped: func() {
		err := player.Previous()
		if err != nil {
			fyne.LogError("Could not skip to previous", err)
		}
	}}

	status, err := player.GetPlaybackStatus()
	if err != nil {
		fyne.LogError("Failed to get playback status", err)
		return
	}

	icon := theme.MediaPauseIcon()
	if status == mpris.PlaybackPaused {
		icon = theme.MediaPlayIcon()
	}

	playOrPause := &widget.Button{Icon: icon}
	playOrPause.OnTapped = func() {
		err := player.PlayPause()
		if err != nil {
			fyne.LogError("Could not change playback mode", err)
			return
		}

		status, err := player.GetPlaybackStatus()
		if err != nil {
			fyne.LogError("Failed to get playback status", err)
			return
		}

		if status == mpris.PlaybackPaused {
			playOrPause.Icon = theme.MediaPlayIcon()
		} else {
			playOrPause.Icon = theme.MediaPauseIcon()
		}

		playOrPause.Refresh()
	}

	next := &widget.Button{Icon: theme.MediaSkipNextIcon(), OnTapped: func() {
		err := player.Next()
		if err != nil {
			fyne.LogError("Could not skip to next", err)
		}
	}}

	w.SetContent(container.NewHBox(previous, playOrPause, next))
	w.ShowAndRun()
}
