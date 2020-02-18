package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiradata"
	"github.com/jason0x43/go-toggl"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func main() {

	toggl.DisableLog()

	jiraSession := jiraSession()
	togglSession := togglSession()

	for _, entry := range togglEntries(togglSession) {

		// Only care about entries not yet logged in JIRA
		if entry.HasTag("Logged in JIRA") || entry.HasTag("toggl2jira ignore") || entry.HasTag("Not Logging in JIRA") {
			continue
		}

		// Don't care about entries currently in progress
		if entry.IsRunning() {
			log.Infof("Skipping running entry: %s", entry.Description)
			continue
		}

		descParts := strings.Fields(entry.Description)
		// Make the assumption that the first word in the description is
		// a JIRA ticket
		jiraTicket := descParts[0]

		if entry.Duration < 60 {
			entry.Duration = 60
		}

		timeSpent, err := time.ParseDuration(fmt.Sprintf("%vs", entry.Duration))
		if err != nil {
			log.Fatalf("Failed to parse time entry duration: %s", err)
		}

		issue, err := jiraSession.GetIssue(jiraTicket, &jira.IssueOptions{})
		if err != nil {
			log.Warnf("Could not find issue: %s", jiraTicket)
			continue
		}

		started := entry.Start.Format("2006-01-02T15:04:05.000-0700")

		log.Infof("%v - %v - %v - %v", jiraTicket, issue.Fields["summary"], timeSpent, started)

		worklog := &jiradata.Worklog{
			Comment:          "Autologged with toggl2jira",
			Started:          started,
			TimeSpentSeconds: int(entry.Duration),
		}
		_, err = jiraSession.AddIssueWorklog(jiraTicket, worklog)
		if err != nil {
			log.Errorf("Could not log time in JIRA: %s", err)
			continue
		}
		log.Infof("JIRA updated")

		entry.AddTag("toggl2jira")
		entry.AddTag("Logged in JIRA")
		_, err = togglSession.UpdateTimeEntry(entry)
		if err != nil {
			log.Errorf("Could not update time entry: %s", err)
			continue
		}
		log.Infof("Toggl updated")

	}

}
