package main

import (
	"encoding/json"
	"os"
	"time"
)

type Log struct {
	DateTime time.Time
	Type     IntegrityChangeType
	Data     any
}

func SaveLogs(log Log, logPath string) error {
	logBytes, err := json.Marshal(log)
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer logFile.Close()

	if _, err = logFile.WriteString(string(logBytes) + "\n"); err != nil {
		return err
	}

	return nil
}
