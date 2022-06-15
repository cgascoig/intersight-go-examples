package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	poller, err := NewIntersightAlarmPoller(os.Getenv("IS_KEY_ID"), os.Getenv("IS_KEY_FILE"))
	if err != nil {
		log.Fatalf("creating new poller: %v", err)
	}

	log.Error(poller.Run())
}
