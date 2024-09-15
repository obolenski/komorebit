package main

import (
	"fmt"
	"komorebit/internal/komorebi"
	"komorebit/internal/tray"

	"github.com/getlantern/systray"
)

func main() {
	onExit := func() {
		komorebi.Unsubscribe()
		fmt.Println("Bye")
	}

	onReady := func() {
		fmt.Println("Hi")
		app := tray.GetApp()
		app.Init()
		go komorebi.HandleEvents()
		komorebi.Subscribe()
	}

	systray.Run(onReady, onExit)
}
