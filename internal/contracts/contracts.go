package contracts

type EventData struct {
	Event Event `json:"event"`
	State State `json:"state"`
}

type Event struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type State struct {
	Monitors                          Monitors    `json:"monitors"`
	IsPaused                          bool        `json:"is_paused"`
	ResizeDelta                       int         `json:"resize_delta"`
	NewWindowBehaviour                string      `json:"new_window_behaviour"`
	CrossMonitorMoveBehaviour         string      `json:"cross_monitor_move_behaviour"`
	UnmanagedWindowOperationBehaviour string      `json:"unmanaged_window_operation_behaviour"`
	WorkAreaOffset                    interface{} `json:"work_area_offset"`
	FocusFollowsMouse                 interface{} `json:"focus_follows_mouse"`
	MouseFollowsFocus                 bool        `json:"mouse_follows_focus"`
	HasPendingRaiseOp                 bool        `json:"has_pending_raise_op"`
}

type Monitors struct {
	Elements []Monitor `json:"elements"`
	Focused  int       `json:"focused"`
}

type Monitor struct {
	ID                        int         `json:"id"`
	Name                      string      `json:"name"`
	Device                    string      `json:"device"`
	DeviceID                  string      `json:"device_id"`
	Size                      Rect        `json:"size"`
	WorkAreaSize              Rect        `json:"work_area_size"`
	WorkAreaOffset            interface{} `json:"work_area_offset"`
	WindowBasedWorkAreaOffset interface{} `json:"window_based_work_area_offset"`
	WindowBasedWorkAreaLimit  int         `json:"window_based_work_area_offset_limit"`
	Workspaces                Workspaces  `json:"workspaces"`
}

type Workspaces struct {
	Elements []Workspace `json:"elements"`
	Focused  int         `json:"focused"`
}

type Workspace struct {
	Name              string        `json:"name"`
	Containers        Containers    `json:"containers"`
	MonocleContainer  interface{}   `json:"monocle_container"`
	MaximizedWindow   interface{}   `json:"maximized_window"`
	FloatingWindows   []interface{} `json:"floating_windows"`
	Layout            Layout        `json:"layout"`
	LatestLayout      []Rect        `json:"latest_layout"`
	Tile              bool          `json:"tile"`
	ApplyWindowOffset bool          `json:"apply_window_based_work_area_offset"`
}

type Containers struct {
	Elements []Container `json:"elements"`
	Focused  int         `json:"focused"`
}

type Container struct {
	ID      string  `json:"id"`
	Windows Windows `json:"windows"`
}

type Windows struct {
	Elements []Window `json:"elements"`
	Focused  int      `json:"focused"`
}

type Window struct {
	Hwnd  int    `json:"hwnd"`
	Title string `json:"title"`
	Exe   string `json:"exe"`
	Class string `json:"class"`
	Rect  Rect   `json:"rect"`
}

type Rect struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

type Layout struct {
	Default string `json:"Default"`
}
