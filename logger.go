package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/logging"
)

var loggerName string = "HerodotusDev"

// For now, print logs both to stdout and to monitoring
func _logInternal(fullMessage string, sev logging.Severity) {
	projectID := os.Getenv("PROJECT_ID")

	ctx := context.Background()
	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer loggingClient.Close()

	if _loggerName, ok := os.LookupEnv("LOGGER_NAME"); ok {
		loggerName = _loggerName
	}

	logger := loggingClient.Logger(loggerName).StandardLogger(sev)

	log.Println(fullMessage)
	logger.Println(fullMessage)
}

func LogError(message string) {
	fullMessage := fmt.Sprintf("error: \n%v\n", message)
	_logInternal(fullMessage, logging.Error)
}

func LogInfo(message string) {
	fullMessage := fmt.Sprintf("info: \n%v\n", message)
	_logInternal(fullMessage, logging.Info)
}

func LogDebug(message string) {
	fullMessage := fmt.Sprintf("debug: \n%v\n", message)
	_logInternal(fullMessage, logging.Debug)
}

func LogWarning(message string) {
	fullMessage := fmt.Sprintf("warning: \n%v\n", message)
	_logInternal(fullMessage, logging.Warning)
}
