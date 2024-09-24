package events

import (
	"bufio"
	"encoding/json"
	"fmt"
	"komorebit/internal/contracts"
	"komorebit/internal/komorebic"
	"net"
	"os"
	"sync"

	"github.com/Microsoft/go-winio"
)

type EventHandler interface {
	HandleEvent(contracts.EventData)
}

type Manager struct {
	handler  EventHandler
	stopChan chan struct{}
	stopped  bool
	mu       sync.Mutex
}

func NewManager(handler EventHandler) *Manager {
	return &Manager{
		handler:  handler,
		stopChan: make(chan struct{}),
		stopped:  false,
	}
}

func (m *Manager) Start() {
	go m.handleEvents()
}

func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.stopped {
		close(m.stopChan)
		m.stopped = true
	}
}

func (m *Manager) Restart() {
	m.Stop()
	m.mu.Lock()
	m.stopChan = make(chan struct{})
	m.stopped = false
	m.mu.Unlock()
	m.Start()
}

func (m *Manager) handleEvents() {
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

	go func() {
		<-m.stopChan
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-m.stopChan:
				return
			default:
				fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
				continue
			}
		}
		fmt.Println("Client connected")
		go m.handleClient(conn)
	}
}

func (m *Manager) handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		notification, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}
		var eventData contracts.EventData
		err = json.Unmarshal(notification, &eventData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling JSON: %v\n", err)
			continue
		}
		m.handler.HandleEvent(eventData)
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
