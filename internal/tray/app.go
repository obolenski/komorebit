package tray

import (
	"fmt"
	"komorebit/internal/contracts"
	"komorebit/internal/events"
	"komorebit/internal/icons"
	"komorebit/internal/komorebic"
	"strconv"
	"sync"
	"time"

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
	activeMonitorIndex   int
	activeWorkspaceIndex int
	activeLayout         string
	isPaused             bool
	isSad                bool
	isInit               bool
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
				activeMonitorIndex:   0,
				activeWorkspaceIndex: 0,
				activeLayout:         "",
				isPaused:             false,
				isSad:                false,
				isInit:               true,
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

	a.menu.activeWorkspaceIndicator = systray.AddMenuItem("Active workspace: ???", "Indicates the active workspace")
	a.menu.activeWorkspaceIndicator.SetIcon(icons.WorkspacesIcon())
	workspaceButtonCount := 10
	for i := 0; i < workspaceButtonCount; i++ {
		btnTitle := "Workspace " + strconv.Itoa(i+1)
		btn := a.menu.activeWorkspaceIndicator.AddSubMenuItem(btnTitle, btnTitle)
		a.menu.workspaceChangeBtns[btnTitle] = btn
		go a.handleWorkspaceChangeButton(btn, i)
	}

	a.menu.activeLayoutIndicator = systray.AddMenuItem("Current layout: ???", "Indicates the current layout")
	a.menu.activeLayoutIndicator.SetIcon(icons.LayoutsIcon())
	layoutOptions := []string{"bsp", "columns", "rows", "vertical-stack", "horizontal-stack", "ultrawide-vertical-stack", "grid", "right-main-vertical-stack"}
	for _, layout := range layoutOptions {
		btn := a.menu.activeLayoutIndicator.AddSubMenuItem(layout, layout)
		a.menu.layoutChangeBtns[layout] = btn
		go a.handleLayoutChangeButton(btn, layout)
	}

	systray.AddSeparator()
	a.menu.pauseBtn = systray.AddMenuItem("Pause/unpause komorebi", "Pause/unpause komorebi (komorebic toggle-pause)")
	a.menu.pauseBtn.SetIcon(icons.PauseIcon())
	a.menu.reloadBtn = systray.AddMenuItem("Reload komorebi", "Reload komorebi (komorebic stop; komorebic start)")
	a.menu.reloadBtn.SetIcon(icons.RefreshIcon())
	systray.AddSeparator()
	a.menu.quitBtn = systray.AddMenuItem("Quit komorebit", "Quit the app (does not affect komorebi)")
	a.menu.quitBtn.SetIcon(icons.QuitIcon())

	go a.handlePauseButton()
	go a.handleReloadButton()
	go a.handleQuitButton()

	go a.monitorKomorebi()
}

func (a *App) HandleEvent(eventData contracts.EventData) {
	if eventData.Event.Type == "TogglePause" {
		a.togglePause()
		return
	} else if eventData.Event.Type == "KomorebiStopped" {
		a.beSad()
		return
	}

	activeMonitorIndex := int(eventData.State.Monitors.Focused)
	activeWorkspaceIndex := int(eventData.State.Monitors.Elements[activeMonitorIndex].Workspaces.Focused)
	activeLayout := eventData.State.Monitors.Elements[activeMonitorIndex].Workspaces.Elements[activeWorkspaceIndex].Layout.Default

	a.setActiveWorkspaceIndex(activeWorkspaceIndex)
	a.setActiveLayout(activeLayout)
	a.setActiveMonitorIndex(activeMonitorIndex)
}

func (a *App) handleLayoutChangeButton(button *systray.MenuItem, layout string) {
	for {
		<-button.ClickedCh
		fmt.Println("Requesting layout change")
		komorebic.Exec([]string{"change-layout", layout})
		fmt.Println("Finished layout change")
	}
}

func (a *App) handleWorkspaceChangeButton(button *systray.MenuItem, workspace int) {
	for {
		<-button.ClickedCh
		requestedWorkspaceIndex := strconv.Itoa(workspace)
		fmt.Println("Requesting workspace change")
		komorebic.Exec([]string{"focus-workspace", requestedWorkspaceIndex})
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

func (a *App) setActiveMonitorIndex(index int) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeMonitorIndex = index
}

func (a *App) setActiveWorkspaceIndex(index int) {
	a.state.mu.Lock()
	defer a.state.mu.Unlock()
	a.state.activeWorkspaceIndex = index
	a.updateWorkspaceIcon()
	a.menu.activeWorkspaceIndicator.SetTitle("Active workspace: " + strconv.Itoa(index+1))
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

func (a *App) monitorKomorebi() {
	for {
		if a.state.isSad || a.state.isInit {
			a.hope()
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func (a *App) beSad() {
	a.stopEvents()
	a.state.mu.Lock()
	a.state.isSad = true
	a.state.mu.Unlock()
	systray.SetIcon(icons.SadIcon())
	go a.hope()
}

func (a *App) hope() {
	fmt.Println("Hoping...")
	_, err := komorebic.Exec([]string{"state"})
	if err == nil {
		a.beHappy()
	}
	time.Sleep(1 * time.Second)
}

func (a *App) beHappy() {
	a.state.mu.Lock()
	a.state.isSad = false
	a.state.isInit = false
	a.state.mu.Unlock()
	a.initEvents()
	a.updateIcon()
}

func (a *App) stopEvents() {
	events.Unsubscribe()
	if a.eventManager != nil {
		a.eventManager.Stop()
	}
}

func (a *App) initEvents() {
	newEventManager := events.NewManager(a)
	newEventManager.Start()
	a.eventManager = newEventManager
	events.Subscribe()
}
