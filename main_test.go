package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckStatusCode(t *testing.T) {
	t.Run("Successfully validates success status code", func(t *testing.T) {
		statusCode := 200
		exp := "Success"
		s, err := checkStatusCode(statusCode)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates forbidden status code", func(t *testing.T) {
		statusCode := 403
		exp := "forbidden"
		s, err := checkStatusCode(statusCode)
		require.EqualError(t, err, exp)
		require.Empty(t, s)
	})

	t.Run("Successfully validates not found status code", func(t *testing.T) {
		statusCode := 404
		exp := "not found"
		s, err := checkStatusCode(statusCode)
		require.EqualError(t, err, exp)
		require.Empty(t, s)
	})

	t.Run("Successfully validates service unavailable status code", func(t *testing.T) {
		statusCode := 503
		exp := "service unavailable"
		s, err := checkStatusCode(statusCode)
		require.EqualError(t, err, exp)
		require.Empty(t, s)
	})
}

func TestParseCreateEvent(t *testing.T) {
	t.Run("Successfully validates CreateEvent for repository", func(t *testing.T) {
		payload := `{
						"ref": null,
						"ref_type": "repository",
						"master_branch": "main",
						"description": "A brand new repository",
						"pusher_type": "user"
					}`
		reponame := "devUser/my-repo"
		exp := fmt.Sprintf("Created new repository %s", reponame)
		s, err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates CreateEvent for branch", func(t *testing.T) {
		payload := `{
						"ref": "feature-branch",
						"ref_type": "branch",
						"master_branch": "main",
						"description": "A repository for an awesome project",
						"pusher_type": "user"
					}`
		reponame := "devUser/my-repo"
		exp := fmt.Sprintf("Created new branch %s", reponame)
		s, err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates CreateEvent for tag", func(t *testing.T) {
		payload := `{
						"ref": "v1.0.0",
						"ref_type": "tag",
						"master_branch": "main",
						"description": "Release version 1.0.0",
						"pusher_type": "user"
					}`
		reponame := "devUser/my-repo"
		exp := fmt.Sprintf("Created new tag %s", reponame)
		s, err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates error for CreateEvent", func(t *testing.T) {
		payload := `{
						"ref": "v1.0.0",
						"ref_type": "",
						"master_branch": "main",
						"description": "Release version 1.0.0",
						"pusher_type": "user"
					}`
		reponame := "sample_repo"
		s, err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.EqualError(t, err, "unable to parse, reference type is empty")
		require.Empty(t, s)
	})
}

func TestParseDeleteEvent(t *testing.T) {
	t.Run("Successfully validates DeleteEvent for branch", func(t *testing.T) {
		payload := `{
						"ref": "feature-branch",
						"ref_type": "branch",
						"pusher_type": "user"
					}`
		reponame := "devUser/my-repo"
		exp := fmt.Sprintf("Deleted branch %s\n", reponame)
		s, err := parseDeleteEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})
	t.Run("Successfully validates DeleteEvent for tag", func(t *testing.T) {
		payload := `{
						"ref": "v1.0.0",
						"ref_type": "tag",
						"pusher_type": "user"
					}`
		reponame := "devUser/my-repo"
		exp := fmt.Sprintf("Deleted tag %s\n", reponame)
		s, err := parseDeleteEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates error for DeleteEvent", func(t *testing.T) {
		payload := `{
						"ref": "v1.0.0",
						"ref_type": "",
						"pusher_type": "user"
					}`
		reponame := "devUser/my-repo"
		s, err := parseDeleteEvent(json.RawMessage(payload), reponame)
		require.EqualError(t, err, "unable to parse, reference type is empty")
		require.Empty(t, s)
	})
}

func TestParseIssuesEvent(t *testing.T) {
	t.Run("Successfully validates IssuesEvent for open issue", func(t *testing.T) {
		payload := `{
					"action": "opened",
					"issue": {
						"id": 55667788,
						"number": 42,
						"title": "Bug: Application crashes on startup",
						"state": "open",
						"body": "The application crashes when attempting to launch on Windows 11.",
						"user": {
						"login": "devUser",
						"id": 112233
						},
						"labels": [
						{
							"id": 121212,
							"name": "bug",
							"color": "d73a4a"
						}
						]
					}
				}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Application crashes on startup"
		exp := fmt.Sprintf("Issue %d. %s for %s is opened\n", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})

	t.Run("Successfully validates IssuesEvent for edited issue", func(t *testing.T) {
		payload := `{
						"action": "edited",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"body": "Updated description with additional crash logs.",
							"state": "open"
						},
						"changes": {
							"title": {
							"from": "Bug: Application crashes on startup"
							},
							"body": {
							"from": "The application crashes when attempting to launch on Windows 11."
							}
						}
					}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Crash on startup (Updated)"
		exp := fmt.Sprintf("Issue %d. %s for %s is edited\n", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})

	t.Run("Successfully validates IssuesEvent for closed issue", func(t *testing.T) {
		payload := `{
						"action": "closed",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"state": "closed"
						}
					}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Crash on startup (Updated)"
		exp := fmt.Sprintf("Issue %d. %s for %s is closed\n", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})

	t.Run("Successfully validates IssuesEvent for reopened issue", func(t *testing.T) {
		payload := `{
						"action": "reopened",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"state": "open"
						}
					}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Crash on startup (Updated)"
		exp := fmt.Sprintf("Issue %d. %s for %s is reopened\n", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})

	t.Run("Successfully validates IssuesEvent for assigned issue", func(t *testing.T) {
		payload := `{
						"action": "assigned",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"state": "open"
						},
						"assignee": {
							"login": "maintainerUser",
							"id": 445566
						}
					}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Crash on startup (Updated)"
		assignee := "maintainerUser"
		exp := fmt.Sprintf("Issue %d. %s for %s is assigned to %s\n", issuenum, issuetitle, reponame, assignee)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})

	t.Run("Successfully validates IssuesEvent for unassigned issue", func(t *testing.T) {
		payload := `{
						"action": "unassigned",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"state": "open"
						},
						"assignee": {
							"login": "maintainerUser",
							"id": 445566
						}
					}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Crash on startup (Updated)"
		assignee := "maintainerUser"
		exp := fmt.Sprintf("Issue %d. %s for %s is unassigned from %s\n", issuenum, issuetitle, reponame, assignee)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})

	t.Run("Successfully validates IssuesEvent for labeled issue", func(t *testing.T) {
		payload := `{
						"action": "labeled",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"state": "open"
						},
						"label": {
							"id": 234567,
							"name": "priority: high",
							"color": "ff0000"
						}
					}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Crash on startup (Updated)"
		label := "priority: high"
		exp := fmt.Sprintf("Issue %d. %s for %s is labeled as %s\n", issuenum, issuetitle, reponame, label)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})

	t.Run("Successfully validates IssuesEvent for unlabeled issue", func(t *testing.T) {
		payload := `{
						"action": "unlabeled",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"state": "open"
						},
						"label": {
							"id": 234567,
							"name": "priority: high",
							"color": "ff0000"
						}
					}`
		reponame := "devUser/awesome-project"
		issuenum := 42
		issuetitle := "Bug: Crash on startup (Updated)"
		label := "priority: high"
		exp := fmt.Sprintf("Issue %d. %s for %s is unlabeled from %s\n", issuenum, issuetitle, reponame, label)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, s, exp)
	})
}
