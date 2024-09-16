package tray

import (
	"encoding/json"
	"fmt"
	"komorebit/internal/contracts"
	"komorebit/internal/events"
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
	state        AppState
	menu         AppMenuItems
	eventManager *events.Manager
}

type AppState struct {
	activeWorkspaceIndex int
	activeWorkspaceName  string
	activeLayout         string
	isPaused             bool
	mu                   sync.Mutex
}

type AppMenuItems struct {
	activeWorkspaceIndicator *systray.MenuItem
	workspaceChangeBtns      map[string]*systray.MenuItem
	activeLayoutIndicator    *systray.MenuItem
	layoutChangeBtns         map[string]*systray.MenuItem
	pauseBtn                 *systray.MenuItem
	reloadBtn                *systray.MenuItem
	quitBtn                  *systray.MenuItem
}

func GetApp() *App {
	once.Do(func() {
		instance = &App{
			state: AppState{
				activeWorkspaceIndex: 0,
				activeWorkspaceName:  "",
				activeLayout:         "",
				isPaused:             false,
			},
			menu: AppMenuItems{
				layoutChangeBtns:    make(map[string]*systray.MenuItem),
				workspaceChangeBtns: make(map[string]*systray.MenuItem),
			},
		}
	})
	return instance
}

func (a *App) Init(eventManager *events.Manager) {
	a.eventManager = eventManager
	systray.SetIcon(icons.TildeIcon())
	systray.SetTitle("Komorebit")
	systray.SetTooltip("Komorebit")

	currentState, err := komorebic.Exec([]string{"state"})
	if err != nil {
		fmt.Println("Error getting komorebi state")
		systray.SetIcon(icons.SadIcon())
		return
	}
	var state contracts.State
	err = json.Unmarshal([]byte(currentState), &state)

	if err != nil {
		fmt.Println("Error unmarshalling komorebi state")
		systray.SetIcon(icons.SadIcon())
		return
	}

	activeMonitorIndex := int(state.Monitors.Focused)
	activeWorkspaceIndex := int(state.Monitors.Elements[activeMonitorIndex].Workspaces.Focused)
	activeWorkspaceName := state.Monitors.Elements[activeMonitorIndex].Workspaces.Elements[activeWorkspaceIndex].Name
	activeLayout := state.Monitors.Elements[activeMonitorIndex].Workspaces.Elements[activeWorkspaceIndex].Layout.Default

	a.menu.activeWorkspaceIndicator = systray.AddMenuItem("Active workspace: "+activeWorkspaceName, "Indicates the active workspace")
	workspaces := state.Monitors.Elements[activeMonitorIndex].Workspaces.Elements
	for _, workspace := range workspaces {
		btn := a.menu.activeWorkspaceIndicator.AddSubMenuItem(workspace.Name, workspace.Name)
		a.menu.workspaceChangeBtns[workspace.Name] = btn
		go a.handleWorkspaceChangeButton(btn, workspace.Name)
	}

	a.menu.activeLayoutIndicator = systray.AddMenuItem("Current layout: "+activeLayout, "Indicates the current layout")
	layoutOptions := []string{"bsp", "columns", "rows", "vertical-stack", "horizontal-stack", "ultrawide-vertical-stack", "grid", "right-main-vertical-stack"}
	for _, layout := range layoutOptions {
		btn := a.menu.activeLayoutIndicator.AddSubMenuItem(layout, layout)
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

func (a *App) HandleEvent(eventData contracts.EventData) {
	if eventData.Event.Type == "TogglePause" {
		a.togglePause()
		return
	}

	activeMonitorIndex := int(eventData.State.Monitors.Focused)
	activeWorkspaceIndex := int(eventData.State.Monitors.Elements[activeMonitorIndex].Workspaces.Focused)
	activeWorkspaceName := eventData.State.Monitors.Elements[activeMonitorIndex].Workspaces.Elements[activeWorkspaceIndex].Name
	activeLayout := eventData.State.Monitors.Elements[activeMonitorIndex].Workspaces.Elements[activeWorkspaceIndex].Layout.Default

	a.setActiveWorkspaceIndex(activeWorkspaceIndex)
	a.setActiveWorkspaceName(activeWorkspaceName)
	a.setActiveLayout(activeLayout)

}

func (a *App) handleLayoutChangeButton(button *systray.MenuItem, layout string) {
	for {
		<-button.ClickedCh
		fmt.Println("Requesting layout change")
		komorebic.Exec([]string{"change-layout", layout})
		fmt.Println("Finished layout change")
	}
}

func (a *App) handleWorkspaceChangeButton(button *systray.MenuItem, workspace string) {
	for {
		<-button.ClickedCh
		fmt.Println("Requesting workspace change")
		komorebic.Exec([]string{"focus-named-workspace", workspace})
		fmt.Println("Finished workspace change")
	}
}

func (a *App) handlePauseButton() {
	for {
		<-a.menu.pauseBtn.ClickedCh
		fmt.Println("Requesting pause")
		komorebic.Exec([]string{"toggle-pause"})
		fmt.Println("Finished pausing")
	}
}

func (a *App) handleReloadButton() {
	for {
		<-a.menu.reloadBtn.ClickedCh
		fmt.Println("Requesting reload")
		systray.SetIcon(icons.TildeIcon())
		a.stopEvents()
		komorebic.Exec([]string{"stop"})
		_, err := komorebic.Exec([]string{"start"})
		if err != nil {
			systray.SetIcon(icons.SadIcon())
			fmt.Println("Error reloading komorebi")
		} else {
			a.initEvents()
			fmt.Println("Finished reloading")
		}
	}
}

func (a *App) handleQuitButton() {
	<-a.menu.quitBtn.ClickedCh
	fmt.Println("Requesting quit")
	systray.Quit()
	fmt.Println("Finished quitting")
}

func (a *App) setActiveWorkspaceIndex(index int) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeWorkspaceIndex = index
	a.updateWorkspaceIcon()
}

func (a *App) setActiveWorkspaceName(name string) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeWorkspaceName = name
	if a.state.activeWorkspaceName != "" {
		a.menu.activeWorkspaceIndicator.SetTitle("Active workspace: " + a.state.activeWorkspaceName)
	}
}

func (a *App) setActiveLayout(layout string) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeLayout = layout
	if a.state.activeLayout != "" {
		a.menu.activeLayoutIndicator.SetTitle("Current layout: " + a.state.activeLayout)
	}
}

func (a *App) togglePause() {
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

func (a *App) stopEvents() {
	events.Unsubscribe()
	a.eventManager.Stop()
}

func (a *App) initEvents() {
	newEventManager := events.NewManager(a)
	newEventManager.Start()
	a.eventManager = newEventManager
	events.Subscribe()
}
