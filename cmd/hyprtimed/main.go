package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"hyprtime/internal/daemon/api"
	"hyprtime/internal/daemon/database"
	"hyprtime/internal/daemon/tracker"
	"hyprtime/internal/logger"
)

func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "hyprtime", "daemon.sock")
}

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	debug := flag.Bool("debug", false, "Enable debug logging (implies verbose)")
	quiet := flag.Bool("quiet", false, "Suppress all non-error logs")
	flag.Parse()

	logger.Init()

	if *debug {
		logger.SetLogLevel(logger.LogLevelDebug)
	} else if *verbose {
		logger.SetLogLevel(logger.LogLevelVerbose)
	} else if *quiet {
		logger.SetLogLevel(logger.LogLevelQuiet)
	}

	logger.Info("Starting Hyprtime Daemon")

	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	defer db.Close()

	socketPath := getSocketPath()
	apiServer, err := api.NewServer(db, socketPath)
	if err != nil {
		logger.Fatal("Failed to create API server: %v", err)
	}

	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Error("API server error: %v", err)
		}
	}()

	tr := tracker.New(db)
	if err := tr.Start(); err != nil {
		logger.Fatal("Failed to start tracker: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutting down...")

	tr.Stop()
	apiServer.Close()

	time.Sleep(500 * time.Millisecond)
	logger.Info("Daemon stopped")
}
