package main

import (
	"fmt"
	"komorebit/internal/events"
	"komorebit/internal/tray"

	"github.com/getlantern/systray"
)

func main() {
	app := tray.GetApp()
	eventManager := events.NewManager(app)

	onExit := func() {
		events.Unsubscribe()
		fmt.Println("Bye")
	}

	onReady := func() {
		fmt.Println("Hi")
		app.Init(eventManager)
	}

	systray.Run(onReady, onExit)
}
