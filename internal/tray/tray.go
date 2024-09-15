package tray

import (
	"fmt"
	"komorebit/internal/icons"
	"komorebit/internal/komorebic"
	"sync"

	"github.com/getlantern/systray"
)

var (
	instance *App
	once     sync.Once
)

type App struct {
	state AppState
	menu  AppMenuItems
}

type AppState struct {
	activeMonitorIndex   int
	activeWorkspaceIndex int
	activeLayout         string
	isPaused             bool
	mu                   sync.Mutex
}

type AppMenuItems struct {
	currentLayoutIndicator *systray.MenuItem
	layoutChangeBtns       map[string]*systray.MenuItem
	pauseBtn               *systray.MenuItem
	reloadBtn              *systray.MenuItem
	quitBtn                *systray.MenuItem
}

func GetApp() *App {
	once.Do(func() {
		instance = &App{
			state: AppState{
				activeMonitorIndex:   0,
				activeWorkspaceIndex: 0,
				activeLayout:         "",
				isPaused:             false,
			},
		}
	})
	return instance
}

func (a *App) Init() {
	systray.SetIcon(icons.TildeIcon())
	systray.SetTitle("Komorebit")
	systray.SetTooltip("Komorebit")

	a.menu.layoutChangeBtns = make(map[string]*systray.MenuItem)

	a.menu.currentLayoutIndicator = systray.AddMenuItem("Current layout: ???", "Indicates the current layout")
	var layoutOptions []string = []string{"bsp", "columns", "rows", "vertical-stack", "horizontal-stack", "ultrawide-vertical-stack", "grid", "right-main-vertical-stack"}
	for _, layout := range layoutOptions {
		btn := a.menu.currentLayoutIndicator.AddSubMenuItem(layout, layout)
		a.menu.layoutChangeBtns[layout] = btn
		go a.handleLayoutChangeButton(btn, layout)
	}

	systray.AddSeparator()

	a.menu.pauseBtn = systray.AddMenuItem("Pause/unpause komorebi", "Pause/unpause komorebi (komorebic toggle-pause)")
	a.menu.reloadBtn = systray.AddMenuItem("Reload komorebi", "Reload komorebi (komorebic stop; komorebic start)")

	systray.AddSeparator()

	a.menu.quitBtn = systray.AddMenuItem("Quit komorebit", "Quit the app (does not affect komorebi)")

	go a.handlePauseButton()
	go a.handleReloadButton()
	go a.handleQuitButton()
}

func (a *App) handleLayoutChangeButton(button *systray.MenuItem, layout string) {
	for {
		<-button.ClickedCh
		fmt.Println("Requesting layout change")
		komorebic.Exec([]string{"change-layout", layout})
		fmt.Println("Finished layout change")
	}
}

func (a *App) handlePauseButton() {
	for {
		<-a.menu.pauseBtn.ClickedCh
		fmt.Println("Requesting pause")
		komorebic.Exec([]string{"toggle-pause"})
		a.TogglePause()
		fmt.Println("Finished pausing")
	}
}

func (a *App) handleReloadButton() {
	for {
		<-a.menu.reloadBtn.ClickedCh
		fmt.Println("Requesting reload")
		komorebic.Exec([]string{"stop"})
		komorebic.Exec([]string{"start"})
		fmt.Println("Finished reloading")
	}
}

func (a *App) handleQuitButton() {
	<-a.menu.quitBtn.ClickedCh
	fmt.Println("Requesting quit")
	systray.Quit()
	fmt.Println("Finished quitting")
}

func (a *App) SetActiveMonitorIndex(index int) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeMonitorIndex = index
}

func (a *App) SetActiveWorkspaceIndex(index int) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeWorkspaceIndex = index
	a.updateWorkspaceIcon()
}

func (a *App) SetActiveLayout(layout string) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeLayout = layout
	if a.state.activeLayout != "" {
		a.menu.currentLayoutIndicator.SetTitle("Current layout: " + a.state.activeLayout)
	}
}

func (a *App) TogglePause() {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.isPaused = !a.state.isPaused
	a.updateIcon()
}

func (a *App) updateWorkspaceIcon() {
	systray.SetIcon(icons.WorkspaceIcon(a.state.activeWorkspaceIndex))
}

func (a *App) updateIcon() {
	if a.state.isPaused {
		systray.SetIcon(icons.PauseIcon())
	} else {
		a.updateWorkspaceIcon()
	}
}
