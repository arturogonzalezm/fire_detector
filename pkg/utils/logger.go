package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func SetupLogger() {
	// Generate a file name with the current timestamp
	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("logs/fire_detection_%s.log", timestamp)

	// Create the log file
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v\n", err)
	}

	// Set the log output and format
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}
