package main

import (
	"os"
	"time"

	"github.com/jason0x43/go-toggl"
	log "github.com/sirupsen/logrus"
)

func togglSession() *toggl.Session {
	togglUser := os.Getenv("TOGGL_USER")
	togglPass := os.Getenv("TOGGL_PASS")

	session, err := toggl.NewSession(togglUser, togglPass)
	if err != nil {
		log.Fatalf("Failed to start Toggl session: %s", err)
	}

	return &session
}

func togglEntries(session *toggl.Session) []toggl.TimeEntry {

	togglDuration := os.Getenv("TOGGL_DURATION")
	duration, err := time.ParseDuration(togglDuration) // 2 weeks
	if err != nil {
		log.Fatalf("Failed to parse duration: %s", err)
	}

	now := time.Now()
	then := now.Add(-duration)

	log.Infof("Getting time entries between %v and %v", then, now)
	entries, err := session.GetTimeEntries(then, now)
	if err != nil {
		log.Fatalf("Failed to load time entries: %s", err)
	}

	return entries
}
