package main

import (
	"os"

	"github.com/go-jira/jira"
	log "github.com/sirupsen/logrus"
)

func jiraSession() *jira.Jira {

	jiraUser := os.Getenv("JIRA_USER")
	jiraPass := os.Getenv("JIRA_PASS")
	jiraURL := os.Getenv("JIRA_URL")

	jiraClient := jira.NewJira(jiraURL)

	auth := &jira.AuthOptions{
		Username: jiraUser,
		Password: jiraPass,
	}

	_, err := jiraClient.NewSession(auth)
	if err != nil {
		log.Fatalf("Failed to start JIRA session: %s", err)
	}

	return jiraClient
}
