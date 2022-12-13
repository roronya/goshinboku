package jira

import (
	"fmt"
	"os"
	"testing"
)

func TestGetIssue(t *testing.T) {
	user := os.Getenv("JIRA_USER")
	pass := os.Getenv("JIRA_PASSWORD")
	server := os.Getenv("JIRA_SERVER")
	client, _ := NewClient(user, pass, server)
	issue, _ := GetIssue(client, "AIROES-11051")
	fmt.Printf("%#v\n", issue.Fields.Status)
}
