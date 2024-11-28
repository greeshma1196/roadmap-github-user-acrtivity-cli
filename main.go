package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Event []struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
}

type CreateEvent struct {
	RefType string `json:"ref_type"`
}

type DeleteEvent struct {
	RefType string `json:"ref_type"`
}

type IssuesEvent struct {
	Action string `json:"action"`
	Issue  struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	} `json:"issue"`
	Assignee struct {
		Login string `json:"login"`
	} `json:"assignee"`
	Label struct {
		Name string `json:"name"`
	} `json:"label"`
}

type PullRequestEvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Url   string `json:"url"`
		Title string `json:"title"`
	} `json:"pull_request"`
	Assignee struct {
		Login string `json:"login"`
	} `json:"assignee"`
}

func checkStatusCode(statusCode int) (string, error) {
	var s string
	if statusCode == 200 || statusCode == 304 {
		s = "Success"
	} else if statusCode == 403 {
		return "", fmt.Errorf("forbidden")
	} else if statusCode == 404 {
		return "", fmt.Errorf("not found")
	} else if statusCode == 503 {
		return "", fmt.Errorf("service unavailable")
	}
	return s, nil
}

func parseCreateEvent(payload json.RawMessage, reponame string) (string, error) {
	var s string
	var cresp CreateEvent
	if err := json.Unmarshal(payload, &cresp); err != nil {
		panic(err)
	}

	if cresp.RefType == "repository" {
		s = fmt.Sprintf("Created new repository %s", reponame)
	} else if cresp.RefType == "branch" {
		s = fmt.Sprintf("Created new branch %s", reponame)
	} else if cresp.RefType == "tag" {
		s = fmt.Sprintf("Created new tag %s", reponame)
	} else {
		return "", fmt.Errorf("unable to parse, reference type is empty")
	}

	return s, nil
}

func parseDeleteEvent(payload json.RawMessage, reponame string) (string, error) {
	var s string
	var cresp DeleteEvent
	if err := json.Unmarshal(payload, &cresp); err != nil {
		panic(err)
	}

	if cresp.RefType == "branch" {
		s = fmt.Sprintf("Deleted branch %s\n", reponame)
	} else if cresp.RefType == "tag" {
		s = fmt.Sprintf("Deleted tag %s\n", reponame)
	} else {
		return "", fmt.Errorf("unable to parse, reference type is empty")
	}

	return s, nil
}

func parseIssuesEvent(payload json.RawMessage, reponame string) (string, error) {
	var s string
	var cresp IssuesEvent
	if err := json.Unmarshal(payload, &cresp); err != nil {
		panic(err)
	}

	if cresp.Action == "opened" {
		s = fmt.Sprintf("Issue %d. %s for %s is opened", cresp.Issue.Number, cresp.Issue.Title, reponame)
	} else if cresp.Action == "edited" {
		s = fmt.Sprintf("Issue %d. %s for %s is edited", cresp.Issue.Number, cresp.Issue.Title, reponame)
	} else if cresp.Action == "closed" {
		s = fmt.Sprintf("Issue %d. %s for %s is closed", cresp.Issue.Number, cresp.Issue.Title, reponame)
	} else if cresp.Action == "reopened" {
		s = fmt.Sprintf("Issue %d. %s for %s is reopened", cresp.Issue.Number, cresp.Issue.Title, reponame)
	} else if cresp.Action == "assigned" {
		s = fmt.Sprintf("Issue %d. %s for %s is assigned to %s", cresp.Issue.Number, cresp.Issue.Title, reponame, cresp.Assignee.Login)
	} else if cresp.Action == "unassigned" {
		s = fmt.Sprintf("Issue %d. %s for %s is unassigned from %s", cresp.Issue.Number, cresp.Issue.Title, reponame, cresp.Assignee.Login)
	} else if cresp.Action == "labeled" {
		s = fmt.Sprintf("Issue %d. %s for %s is labeled as %s", cresp.Issue.Number, cresp.Issue.Title, reponame, cresp.Label.Name)
	} else if cresp.Action == "unlabeled" {
		s = fmt.Sprintf("Issue %d. %s for %s is unlabeled from %s", cresp.Issue.Number, cresp.Issue.Title, reponame, cresp.Label.Name)
	} else {
		return "", fmt.Errorf("unable to parse")
	}

	return s, nil
}

func parsePullRequestEvent(payload json.RawMessage, reponame string) (string, error) {
	var s string
	var cresp PullRequestEvent
	if err := json.Unmarshal(payload, &cresp); err != nil {
		panic(err)
	}

	if cresp.Action == "opened" {
		s = fmt.Sprintf("Pull request %d. %s for %s is opened at %s", cresp.Number, cresp.PullRequest.Title, reponame, cresp.PullRequest.Url)
	} else if cresp.Action == "closed" {
		s = fmt.Sprintf("Pull request %d. %s for %s is closed at %s", cresp.Number, cresp.PullRequest.Title, reponame, cresp.PullRequest.Url)
	} else if cresp.Action == "reopened" {
		s = fmt.Sprintf("Pull request %d. %s for %s is reopened at %s", cresp.Number, cresp.PullRequest.Title, reponame, cresp.PullRequest.Url)
	} else if cresp.Action == "assigned" {
		s = fmt.Sprintf("Pull request %d. %s for %s is assigned to %s, %s", cresp.Number, cresp.PullRequest.Title, reponame, cresp.Assignee.Login, cresp.PullRequest.Url)
	} else if cresp.Action == "synchronize" {
		s = fmt.Sprintf("Pull request %d. %s for %s is synchronized, %s", cresp.Number, cresp.PullRequest.Title, reponame, cresp.PullRequest.Url)
	} else {
		return "", fmt.Errorf("unable to parse")
	}

	return s, nil
}

func main() {
	username := os.Args[1]

	url := "https://api.github.com/users/" + username + "/events"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	_, err = checkStatusCode(resp.StatusCode)

	if err != nil {
		panic(err)
	}

	var cresp Event
	if err := json.NewDecoder(resp.Body).Decode(&cresp); err != nil {
		panic(err)
	}

	for _, event := range cresp {
		if event.Type == "CreateEvent" {
			s, err := parseCreateEvent(event.Payload, event.Repo.Name)
			if err != nil {
				panic(err)
			}
			fmt.Println(s)
		} else if event.Type == "DeleteEvent" {
			s, err := parseDeleteEvent(event.Payload, event.Repo.Name)
			if err != nil {
				panic(err)
			}
			fmt.Println(s)
		} else if event.Type == "IssuesEvent" {
			s, err := parseIssuesEvent(event.Payload, event.Repo.Name)
			if err != nil {
				panic(err)
			}
			fmt.Println(s)
		} else if event.Type == "PullRequestEvent" {
			s, err := parsePullRequestEvent(event.Payload, event.Repo.Name)
			if err != nil {
				panic(err)
			}
			fmt.Println(s)
		} else if event.Type == "PushEvent" {
			fmt.Printf("PushEvent: %s\n", event.Repo.Name)
		} else if event.Type == "ReleaseEvent" {
			fmt.Printf("ReleaseEvent: %s\n", event.Repo.Name)
		}
	}

}
