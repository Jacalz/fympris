package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

type mediaController struct {
	artistAndTitle *widget.RichText
	icon           *canvas.Image

	previous *widget.Button
	playback *widget.Button
	next     *widget.Button

	player *mpris.Player
}

func (m *mediaController) playPrevious() {
	err := m.player.Previous()
	if err != nil {
		fyne.LogError("Could not skip to previous", err)
	}
}

func (m *mediaController) playNext() {
	err := m.player.Next()
	if err != nil {
		fyne.LogError("Could not skip to next", err)
	}
}

func (m *mediaController) changePlaybackStatus() {
	err := m.player.PlayPause()
	if err != nil {
		fyne.LogError("Could not change playback mode", err)
	}
}

func (m *mediaController) createUI() *fyne.Container {

	buttons := container.NewHBox(m.previous, m.playback, m.next)
	centeredButtons := container.NewCenter(buttons)

	width := buttons.MinSize().Width
	m.icon.SetMinSize(fyne.NewSize(width, width))

	centeredArtistAndTitle := container.NewCenter(m.artistAndTitle)
	content := container.NewBorder(nil, centeredButtons, nil, nil, m.icon)
	return container.NewBorder(centeredArtistAndTitle, nil, nil, nil, content)
}

func (m *mediaController) metadataChanged(status *dbus.Signal) {
	data := status.Body[1].(map[string]dbus.Variant)

	if val, ok := data["PlaybackStatus"]; ok {
		status := val.Value().(string)
		if status == string(mpris.PlaybackPlaying) {
			m.playback.Icon = theme.MediaPauseIcon()
		} else {
			m.playback.Icon = theme.MediaPlayIcon()
		}

		m.playback.Refresh()
	}

	if val, ok := data["Metadata"]; ok {
		metadata := val.Value().(map[string]dbus.Variant)

		artist := strings.Join(metadata["xesam:artist"].Value().([]string), ",")
		title := metadata["xesam:title"].Value().(string)
		m.artistAndTitle.ParseMarkdown(fmt.Sprintf("**%s**: %s", artist, title))

		if val, ok = metadata["mpris:artUrl"]; ok {
			iconURI := val.Value().(string)
			m.icon.File = iconURI[len("file://"):]
			m.icon.Refresh()
		}
	}
}

func newMediaController(player *mpris.Player) *mediaController {
	artist, title, iconURI := getMetadata(player)

	icon := canvas.NewImageFromURI(iconURI)
	icon.FillMode = canvas.ImageFillContain

	controller := &mediaController{
		artistAndTitle: widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**: %s", artist, title)),
		icon:           icon,
		player:         player,
	}

	controller.previous = &widget.Button{Icon: theme.MediaSkipPreviousIcon(), OnTapped: controller.playPrevious}
	controller.next = &widget.Button{Icon: theme.MediaSkipNextIcon(), OnTapped: controller.playNext}
	controller.playback = &widget.Button{Icon: currentPlaybackIcon(player), OnTapped: controller.changePlaybackStatus}

	setUpMetadataChangeListener(controller)

	return controller
}

func getMetadata(player *mpris.Player) (artist string, title string, iconURI fyne.URI) {
	metadata, err := player.GetMetadata()
	if err != nil {
		fyne.LogError("Could not find metadata", err)
		return "", "", nil
	}

	iconPath := metadata["mpris:artUrl"].Value().(string)
	iconURI, err = storage.ParseURI(iconPath)
	if err != nil {
		fyne.LogError("Failed to parse artwork url", err)
		return "", "", nil
	}

	artist = strings.Join(metadata["xesam:artist"].Value().([]string), ",")
	title = metadata["xesam:title"].Value().(string)
	return artist, title, iconURI
}

func currentPlaybackIcon(player *mpris.Player) fyne.Resource {
	status, err := player.GetPlaybackStatus()
	if err != nil {
		fyne.LogError("Could not find playback status", err)
		return theme.MediaPlayIcon()
	}

	if status == mpris.PlaybackPlaying {
		return theme.MediaPauseIcon()
	}

	return theme.MediaPlayIcon()
}

func setUpMetadataChangeListener(controller *mediaController) {
	statusChanged := make(chan *dbus.Signal)
	err := controller.player.OnSignal(statusChanged)
	if err != nil {
		fyne.LogError("Could not set up change listener", err)
		return
	}

	go func(changes <-chan *dbus.Signal, controller *mediaController) {
		for {
			status := <-changes
			controller.metadataChanged(status)
		}
	}(statusChanged, controller)
}
