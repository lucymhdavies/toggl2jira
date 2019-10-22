package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jason0x43/go-toggl"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func main() {
	togglUser := os.Getenv("TOGGL_USER")
	togglPass := os.Getenv("TOGGL_PASS")
	togglDuration := os.Getenv("TOGGL_DURATION")

	duration, err := time.ParseDuration(togglDuration) // 2 weeks
	if err != nil {
		log.Fatalf("Failed to parse duration: %s", err)
	}

	now := time.Now()
	then := now.Add(-duration)

	session, err := toggl.NewSession(togglUser, togglPass)
	if err != nil {
		log.Fatalf("Failed to start Toggl session: %s", err)
	}

	log.Infof("Getting time entries between %v and %v", then, now)
	entries, err := session.GetTimeEntries(then, now)
	if err != nil {
		log.Fatalf("Failed to load time entries: %s", err)
	}

	for _, entry := range entries {
		d, err := time.ParseDuration(fmt.Sprintf("%vs", entry.Duration))
		if err != nil {
			log.Fatalf("Failed to parse time entry duration: %s", err)
		}
		log.Infof("%s (%s) - %s", entry.Start, d, entry.Description)
	}

}
