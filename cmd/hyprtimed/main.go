package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hyprtime/internal/database"
	"hyprtime/internal/logger"
	"hyprtime/internal/tracker"
)

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

	logger.Info("Starting Hyprland Screen Time Daemon")

	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	defer db.Close()

	tr := tracker.New(db)

	if err := tr.Start(); err != nil {
		logger.Fatal("Failed to start tracker: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutting down...")
	tr.Stop()
	time.Sleep(500 * time.Millisecond)
	logger.Info("Daemon stopped")
}
