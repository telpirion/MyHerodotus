package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/logging"
)

var (
	loggingClient *logging.Client
	logger        *log.Logger
)

func startLogging(projectID, loggerName string) error {
	if loggingClient == nil {
		ctx := context.Background()
		loggingClient, err := logging.NewClient(ctx, projectID)
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer loggingClient.Close()
	}
	logger = loggingClient.Logger(loggerName).StandardLogger(logging.Info)
	return nil
}

// For now, print logs both to stdout and to monitoring
func _logInternal(fullMessage string) {
	log.Println(fullMessage)
	logger.Println(fullMessage)
}

func LogError(message string) {
	fullMessage := fmt.Sprintf("error: \n%v\n", message)
	_logInternal(fullMessage)
}

func LogInfo(message string) {
	fullMessage := fmt.Sprintf("info: \n%v\n", message)
	_logInternal(fullMessage)
}

func LogDebug(message string) {
	fullMessage := fmt.Sprintf("debug: \n%v\n", message)
	_logInternal(fullMessage)
}

func LogWarning(message string) {
	fullMessage := fmt.Sprintf("warning: \n%v\n", message)
	_logInternal(fullMessage)
}
