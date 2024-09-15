package komorebi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"komorebit/internal/komorebic"
	"komorebit/internal/tray"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

func HandleEvents() {
	pipeName := `\\.\pipe\komorebit`

	config := &winio.PipeConfig{
		SecurityDescriptor: "",
		MessageMode:        true,
		InputBufferSize:    4096,
		OutputBufferSize:   4096,
	}

	listener, err := winio.ListenPipe(pipeName, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating named pipe listener: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Named pipe server listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}

		fmt.Println("Client connected")
		go handleClient(conn)
	}
}

func Subscribe() {
	komorebic.Exec([]string{"subscribe-pipe", "komorebit"})
	fmt.Println("Subscribed to komorebi pipe")
}

func Unsubscribe() {
	komorebic.Exec([]string{"unsubscribe-pipe", "komorebit"})
	fmt.Println("Unsubscribed from komorebi pipe")
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		notification, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		var EventData EventData
		err = json.Unmarshal(notification, &EventData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling JSON: %v\n", err)
			continue
		}

		app := tray.GetApp()

		if EventData.Event.Type == "TogglePause" {
			app.TogglePause()
			continue
		}

		activeMonitorIndex := int(EventData.State.Monitors.Focused)
		activeWorkspaceIndex := int(EventData.State.Monitors.Elements[activeMonitorIndex].Workspaces.Focused)
		activeLayout := EventData.State.Monitors.Elements[activeMonitorIndex].Workspaces.Elements[activeWorkspaceIndex].Layout.Default

		app.SetActiveWorkspaceIndex(activeWorkspaceIndex)
		app.SetActiveLayout(activeLayout)

	}
}
