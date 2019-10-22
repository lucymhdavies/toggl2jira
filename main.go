package main

import (
	"fmt"
	"os"
	"strings"
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

		// Only care about entries not yet logged in JIRA
		if !contains(entry.Tags, "Logged in JIRA") {

			descParts := strings.Fields(entry.Description)
			// Make the assumption that the first word in the description is
			// a JIRA ticket
			jiraTicket := descParts[0]

			log.Infof("%s (%s) - %s", entry.Start, jiraTicket, d)

		}
	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
