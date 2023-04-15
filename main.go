package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

func initMRPISPlayer() *mpris.Player {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	names, err := mpris.List(conn)
	if err != nil {
		panic(err)
	}
	if len(names) == 0 {
		log.Fatal("No player found")
	}

	name := names[0]
	return mpris.New(conn, name)
}

func main() {
	a := app.New()
	w := a.NewWindow("MPRIS")

	w.ShowAndRun()
}
