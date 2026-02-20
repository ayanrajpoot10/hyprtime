package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"hyprtime/internal/logger"
)

// HyprlandIPC handles communication with Hyprland via Unix socket
type HyprlandIPC struct {
	socketPath string
	eventPath  string
}

// ActiveWindow represents the currently active window
type ActiveWindow struct {
	Address      string `json:"address"`
	Class        string `json:"class"`
	Title        string `json:"title"`
	InitialClass string `json:"initialClass"`
	InitialTitle string `json:"initialTitle"`
}

// New creates a new Hyprland IPC client
func New() (*HyprlandIPC, error) {
	instance := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if instance == "" {
		return nil, fmt.Errorf("HYPRLAND_INSTANCE_SIGNATURE not set. Are you running Hyprland?")
	}

	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/run/user/" + os.Getenv("UID")
	}

	socketPath := fmt.Sprintf("%s/hypr/%s/.socket.sock", runtimeDir, instance)
	eventPath := fmt.Sprintf("%s/hypr/%s/.socket2.sock", runtimeDir, instance)

	return &HyprlandIPC{
		socketPath: socketPath,
		eventPath:  eventPath,
	}, nil
}

// GetActiveWindow retrieves the currently active window
func (h *HyprlandIPC) GetActiveWindow() (*ActiveWindow, error) {
	conn, err := net.Dial("unix", h.socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Hyprland socket: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("j/activewindow"))
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	buf := make([]byte, 8192)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var window ActiveWindow
	if err := json.Unmarshal(buf[:n], &window); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &window, nil
}

// SubscribeToEvents subscribes to Hyprland events
func (h *HyprlandIPC) SubscribeToEvents(eventChan chan<- string) error {
	conn, err := net.Dial("unix", h.eventPath)
	if err != nil {
		return fmt.Errorf("failed to connect to event socket: %w", err)
	}

	go func() {
		defer conn.Close()
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				logger.Error("Error reading events: %v", err)
				return
			}
			line = strings.TrimSpace(line)
			if line != "" {
				eventChan <- line
			}
		}
	}()

	return nil
}

// ParseEvent parses an event string from Hyprland
func ParseEvent(event string) (string, string) {
	parts := strings.SplitN(event, ">>", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
