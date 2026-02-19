package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hyprtimed/internal/database"
	"hyprtimed/internal/logger"
	"hyprtimed/internal/tracker"
)

func main() {
	// Parse command line flags
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	debug := flag.Bool("debug", false, "Enable debug logging (implies verbose)")
	quiet := flag.Bool("quiet", false, "Suppress all non-error logs")
	flag.Parse()

	// Initialize logger
	logger.Init()

	// Set log level based on flags
	if *debug {
		logger.SetLogLevel(logger.LogLevelDebug)
	} else if *verbose {
		logger.SetLogLevel(logger.LogLevelVerbose)
	} else if *quiet {
		logger.SetLogLevel(logger.LogLevelQuiet)
	}

	logger.Info("Starting Hyprland Screen Time Daemon (log level: %s)", logger.GetLogLevelString())

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize tracker
	tr := tracker.New(db)

	// Start tracking
	if err := tr.Start(); err != nil {
		logger.Fatal("Failed to start tracker: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutting down gracefully...")
	tr.Stop()
	time.Sleep(500 * time.Millisecond)
	logger.Info("Daemon stopped")
}
