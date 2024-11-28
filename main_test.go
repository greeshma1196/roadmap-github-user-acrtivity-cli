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
		exp := "unable to parse, reference type is empty"
		s, err := parseDeleteEvent(json.RawMessage(payload), reponame)
		require.EqualError(t, err, exp)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is opened", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is edited", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is closed", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is reopened", issuenum, issuetitle, reponame)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is assigned to %s", issuenum, issuetitle, reponame, assignee)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is unassigned from %s", issuenum, issuetitle, reponame, assignee)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is labeled as %s", issuenum, issuetitle, reponame, label)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
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
		exp := fmt.Sprintf("Issue %d. %s for %s is unlabeled from %s", issuenum, issuetitle, reponame, label)
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates error for IssuesEvent", func(t *testing.T) {
		payload := `{
						"action": "",
						"issue": {
							"id": 55667788,
							"number": 42,
							"title": "Bug: Crash on startup (Updated)",
							"state": "open"
						}
					}`
		reponame := "devUser/my-repo"
		exp := "unable to parse"
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.EqualError(t, err, exp)
		require.Empty(t, s)
	})
}

func TestParsePullRequestEvent(t *testing.T) {
	t.Run("Successfully validates PullRequestEvent for open pull request", func(t *testing.T) {
		payload := `{
						"action": "opened",
						"number": 42,
						"pull_request": {
							"id": 789456,
							"url": "https://api.github.com/repos/devUser/awesome-project/pulls/42",
							"title": "Add feature X",
							"state": "open",
							"body": "This PR adds feature X with detailed description.",
							"merged": false
						}
					}`
		reponame := "devUser/awesome-project"
		prnum := 42
		prtitle := "Add feature X"
		prurl := "https://api.github.com/repos/devUser/awesome-project/pulls/42"
		exp := fmt.Sprintf("Pull request %d. %s for %s is opened at %s", prnum, prtitle, reponame, prurl)
		s, err := parsePullRequestEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates PullRequestEvent for closed pull request", func(t *testing.T) {
		payload := `{
						"action": "closed",
						"number": 43,
						"pull_request": {
							"id": 789457,
							"url": "https://api.github.com/repos/collabUser/project-repo/pulls/43",
							"title": "Fix bug Y",
							"state": "closed",
							"body": "Bug Y has been fixed.",
							"merged": true
						}
					}`
		reponame := "devUser/awesome-project"
		prnum := 43
		prtitle := "Fix bug Y"
		prurl := "https://api.github.com/repos/collabUser/project-repo/pulls/43"
		exp := fmt.Sprintf("Pull request %d. %s for %s is closed at %s", prnum, prtitle, reponame, prurl)
		s, err := parsePullRequestEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates PullRequestEvent for reopened pull request", func(t *testing.T) {
		payload := `{
						"action": "reopened",
						"number": 44,
						"pull_request": {
							"id": 789458,
							"url": "https://api.github.com/repos/newUser/new-repo/pulls/44",
							"title": "Improve docs",
							"state": "open",
							"body": "Reopening PR for further review.",
							"merged": false
						}
					}`
		reponame := "devUser/awesome-project"
		prnum := 44
		prtitle := "Improve docs"
		prurl := "https://api.github.com/repos/newUser/new-repo/pulls/44"
		exp := fmt.Sprintf("Pull request %d. %s for %s is reopened at %s", prnum, prtitle, reponame, prurl)
		s, err := parsePullRequestEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates PullRequestEvent for assigned pull request", func(t *testing.T) {
		payload := `{
						"action": "assigned",
						"number": 45,
						"pull_request": {
							"id": 789459,
							"url": "https://api.github.com/repos/reviewerUser/review-repo/pulls/45",
							"title": "Add CI/CD pipeline",
							"state": "open",
							"body": "Adding continuous integration setup.",
							"merged": false
						},
						"assignee": {
							"login": "devUser",
							"id": 123456
						}
					}`
		reponame := "devUser/awesome-project"
		prnum := 45
		prtitle := "Add CI/CD pipeline"
		prurl := "https://api.github.com/repos/reviewerUser/review-repo/pulls/45"
		prassignee := "devUser"
		exp := fmt.Sprintf("Pull request %d. %s for %s is assigned to %s, %s", prnum, prtitle, reponame, prassignee, prurl)
		s, err := parsePullRequestEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates PullRequestEvent for synchronized pull request", func(t *testing.T) {
		payload := `{
						"action": "synchronize",
						"number": 46,
						"pull_request": {
							"id": 789460,
							"url": "https://api.github.com/repos/leadMaintainer/core-repo/pulls/46",
							"title": "Update README",
							"state": "open",
							"body": "Updated README with more details.",
							"merged": false
						}
					}`
		reponame := "devUser/awesome-project"
		prnum := 46
		prtitle := "Update README"
		prurl := "https://api.github.com/repos/leadMaintainer/core-repo/pulls/46"
		exp := fmt.Sprintf("Pull request %d. %s for %s is synchronized, %s", prnum, prtitle, reponame, prurl)
		s, err := parsePullRequestEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates PullRequestEvent for synchronized pull request", func(t *testing.T) {
		payload := `{
						"action": "synchronize",
						"number": 46,
						"pull_request": {
							"id": 789460,
							"url": "https://api.github.com/repos/leadMaintainer/core-repo/pulls/46",
							"title": "Update README",
							"state": "open",
							"body": "Updated README with more details.",
							"merged": false
						}
					}`
		reponame := "devUser/awesome-project"
		prnum := 46
		prtitle := "Update README"
		prurl := "https://api.github.com/repos/leadMaintainer/core-repo/pulls/46"
		exp := fmt.Sprintf("Pull request %d. %s for %s is synchronized, %s", prnum, prtitle, reponame, prurl)
		s, err := parsePullRequestEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})

	t.Run("Successfully validates error for PullRequestEvent", func(t *testing.T) {
		payload := `{
						"action": "",
						"number": 46,
						"pull_request": {
							"id": 789460,
							"url": "https://api.github.com/repos/leadMaintainer/core-repo/pulls/46",
							"title": "Update README",
							"state": "open",
							"body": "Updated README with more details.",
							"merged": false
						}
					}`
		reponame := "devUser/my-repo"
		exp := "unable to parse"
		s, err := parseIssuesEvent(json.RawMessage(payload), reponame)
		require.EqualError(t, err, exp)
		require.Empty(t, s)
	})
}

func TestParsePushEvent(t *testing.T) {
	t.Run("Successfully validates PushEvent for a single commit", func(t *testing.T) {
		payload := `{
						"push_id": 2020202020,
						"size": 1,
						"distinct_size": 1,
						"ref": "refs/heads/feature-branch",
						"head": "abcd1234efgh5678ijkl9012mnop3456qrst7890",
						"before": "9876543210fedcba0987654321fedcba09876543",
						"commits": [
							{
							"sha": "abcd1234efgh5678ijkl9012mnop3456qrst7890",
							"author": {
								"email": "collabUser@example.com",
								"name": "collabUser"
							},
							"message": "Initial commit on feature-branch",
							"distinct": true,
							"url": "https://api.github.com/repos/collabUser/another-project/commits/abcd1234efgh5678ijkl9012mnop3456qrst7890"
							}
						]
					}`
		reponame := "devUser/awesome-project"
		size := 1
		exp := fmt.Sprintf("Pushed %d commit to %s", size, reponame)
		s, err := parsePushEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})
	t.Run("Successfully validates PushEvent for multiple commits", func(t *testing.T) {
		payload := `{
					"push_id": 1010101010,
					"size": 2,
					"distinct_size": 2,
					"ref": "refs/heads/main",
					"head": "a1b2c3d4e5f67890abcdef1234567890abcdef12",
					"before": "1234567890abcdef1234567890abcdef12345678",
					"commits": [
						{
						"sha": "a1b2c3d4e5f67890abcdef1234567890abcdef12",
						"author": {
							"email": "devUser@example.com",
							"name": "devUser"
						},
						"message": "Add feature X implementation",
						"distinct": true,
						"url": "https://api.github.com/repos/devUser/awesome-project/commits/a1b2c3d4e5f67890abcdef1234567890abcdef12"
						},
						{
						"sha": "1234abcd5678ef90abcdef1234567890abcdef13",
						"author": {
							"email": "collaborator@example.com",
							"name": "collaboratorUser"
						},
						"message": "Fix bug Y in feature X",
						"distinct": true,
						"url": "https://api.github.com/repos/devUser/awesome-project/commits/1234abcd5678ef90abcdef1234567890abcdef13"
						}
					]
				}`
		reponame := "devUser/awesome-project"
		size := 2
		exp := fmt.Sprintf("Pushed %d commits to %s", size, reponame)
		s, err := parsePushEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
		require.Equal(t, exp, s)
	})
	t.Run("Successfully validates error for PushEvent", func(t *testing.T) {
		payload := `{
					"push_id": 1010101010,
					"size": 0,
					"distinct_size": 2,
					"ref": "refs/heads/main",
					"head": "a1b2c3d4e5f67890abcdef1234567890abcdef12",
					"before": "1234567890abcdef1234567890abcdef12345678",
					"commits": [
						{
						"sha": "a1b2c3d4e5f67890abcdef1234567890abcdef12",
						"author": {
							"email": "devUser@example.com",
							"name": "devUser"
						},
						"message": "Add feature X implementation",
						"distinct": true,
						"url": "https://api.github.com/repos/devUser/awesome-project/commits/a1b2c3d4e5f67890abcdef1234567890abcdef12"
						},
						{
						"sha": "1234abcd5678ef90abcdef1234567890abcdef13",
						"author": {
							"email": "collaborator@example.com",
							"name": "collaboratorUser"
						},
						"message": "Fix bug Y in feature X",
						"distinct": true,
						"url": "https://api.github.com/repos/devUser/awesome-project/commits/1234abcd5678ef90abcdef1234567890abcdef13"
						}
					]
				}`
		reponame := "devUser/awesome-project"
		exp := "unable to parse"
		s, err := parsePushEvent(json.RawMessage(payload), reponame)
		require.EqualError(t, err, exp)
		require.Empty(t, s)
	})
}
