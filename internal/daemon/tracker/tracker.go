package tracker

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"hyprtime/internal/shared/database"
	"hyprtime/internal/daemon/ipc"
	"hyprtime/internal/logger"
)

type Tracker struct {
	db             *sql.DB
	ipc            *ipc.HyprlandIPC
	currentSession *Session
	eventChan      chan string
	stopChan       chan bool
	wg             sync.WaitGroup
	mu             sync.Mutex
	lastUpdate     time.Time
}

type Session struct {
	appID         int64
	class         string
	windowTitle   string
	windowAddress string
	startTime     time.Time
}

func New(db *sql.DB) *Tracker {
	return &Tracker{
		db:        db,
		eventChan: make(chan string, 100),
		stopChan:  make(chan bool),
	}
}

func (t *Tracker) Start() error {
	logger.Info("Initializing Hyprland IPC...")

	var err error
	t.ipc, err = ipc.New()
	if err != nil {
		return fmt.Errorf("failed to initialize Hyprland IPC: %w", err)
	}

	logger.Verbose("Subscribing to Hyprland events...")
	if err := t.ipc.SubscribeToEvents(t.eventChan); err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	t.wg.Add(1)
	go t.processEvents()

	t.wg.Add(1)
	go t.periodicUpdate()

	logger.Info("Tracking started successfully")

	go t.handleFocusChange()

	return nil
}

func (t *Tracker) Stop() {
	logger.Info("Stopping tracker...")

	t.mu.Lock()
	if t.currentSession != nil {
		now := time.Now()
		var duration int64
		if t.lastUpdate.IsZero() {
			duration = int64(now.Sub(t.currentSession.startTime).Seconds())
		} else {
			duration = int64(now.Sub(t.lastUpdate).Seconds())
		}

		if err := database.UpdateAppTime(t.db, t.currentSession.appID, duration); err != nil {
			logger.Error("Error updating app time: %v", err)
		} else {
			logger.Verbose("Final update for: %s (%.1fs)", t.currentSession.class, float64(duration))
		}
		t.currentSession = nil
	}
	t.mu.Unlock()

	close(t.stopChan)
	t.wg.Wait()

	logger.Info("Tracker stopped")
}

func (t *Tracker) processEvents() {
	defer t.wg.Done()

	for {
		select {
		case <-t.stopChan:
			return
		case event := <-t.eventChan:
			if event == "" {
				continue
			}
			eventType, data := ipc.ParseEvent(event)
			logger.Debug("Event received: %s >> %s", eventType, data)

			switch eventType {
			case "activewindow":
				go t.handleFocusChange()
			case "activewindowv2":
				go t.handleFocusChange()
			case "openwindow":
				go t.handleWindowOpen(data)
			case "closewindow":
				go t.handleWindowClose(data)
			}
		}
	}
}

func (t *Tracker) periodicUpdate() {
	defer t.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C:
			t.mu.Lock()
			if t.currentSession != nil {
				now := time.Now()
				var duration int64
				if t.lastUpdate.IsZero() {
					duration = int64(now.Sub(t.currentSession.startTime).Seconds())
				} else {
					duration = int64(now.Sub(t.lastUpdate).Seconds())
				}

				if duration > 0 {
					if err := database.UpdateAppTime(t.db, t.currentSession.appID, duration); err != nil {
						logger.Error("Error in periodic update: %v", err)
					} else {
						logger.Verbose("Periodic update: %s (%.1fs)", t.currentSession.class, float64(duration))
						t.lastUpdate = now
					}
				}
			}
			t.mu.Unlock()
		}
	}
}

func (t *Tracker) handleFocusChange() {
	window, err := t.ipc.GetActiveWindow()
	if err != nil {
		logger.Error("Error getting active window: %v", err)
		return
	}

	if window.Class == "" || window.Address == "" {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.currentSession != nil && t.currentSession.windowAddress == window.Address {
		return
	}

	if t.currentSession != nil {
		now := time.Now()
		var duration int64
		if t.lastUpdate.IsZero() {
			duration = int64(now.Sub(t.currentSession.startTime).Seconds())
		} else {
			duration = int64(now.Sub(t.lastUpdate).Seconds())
		}

		if err := database.UpdateAppTime(t.db, t.currentSession.appID, duration); err != nil {
			logger.Error("Error updating app time: %v", err)
		} else {
			logger.Verbose("Updated: %s (%.1fs)", t.currentSession.class, float64(duration))
		}
	}

	appID, err := database.GetOrCreateApp(t.db, window.Class)
	if err != nil {
		logger.Error("Error getting/creating app: %v", err)
		return
	}

	t.currentSession = &Session{
		appID:         appID,
		class:         window.Class,
		windowTitle:   window.Title,
		windowAddress: window.Address,
		startTime:     time.Now(),
	}
	t.lastUpdate = time.Now()

	logger.Info("Focused: %s (%s)", window.Class, window.Title)
}

func (t *Tracker) handleWindowOpen(data string) {
	parts := strings.SplitN(data, ",", 4)
	if len(parts) < 3 {
		return
	}

	class := parts[2]
	if class == "" {
		return
	}

	appID, err := database.GetOrCreateApp(t.db, class)
	if err != nil {
		logger.Error("Error getting/creating app: %v", err)
		return
	}

	if err := database.IncrementOpenCount(t.db, appID); err != nil {
		logger.Error("Error incrementing open count: %v", err)
		return
	}

	logger.Debug("Window opened: %s", class)
}

func (t *Tracker) handleWindowClose(data string) {
	address := strings.TrimSpace(data)

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.currentSession != nil && t.currentSession.windowAddress == address {
		now := time.Now()
		var duration int64
		if t.lastUpdate.IsZero() {
			duration = int64(now.Sub(t.currentSession.startTime).Seconds())
		} else {
			duration = int64(now.Sub(t.lastUpdate).Seconds())
		}

		if err := database.UpdateAppTime(t.db, t.currentSession.appID, duration); err != nil {
			logger.Error("Error updating app time: %v", err)
		} else {
			logger.Verbose("Window closed: %s (%.1fs)", t.currentSession.class, float64(duration))
		}
		t.currentSession = nil
		t.lastUpdate = time.Time{}
	}
}
