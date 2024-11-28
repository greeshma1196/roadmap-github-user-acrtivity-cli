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

func checkStatusCode(statusCode int) error {

	if statusCode == 200 || statusCode == 304 {
		fmt.Printf("Success\n")
	} else if statusCode == 403 {
		return fmt.Errorf("forbidden")
	} else if statusCode == 404 {
		return fmt.Errorf("not found")
	} else if statusCode == 503 {
		return fmt.Errorf("service unavailable")
	}
	return nil
}

func parseCreateEvent(payload json.RawMessage, reponame string) error {
	var cresp CreateEvent
	if err := json.Unmarshal(payload, &cresp); err != nil {
		panic(err)
	}

	if cresp.RefType == "repository" {
		fmt.Printf("Created new repository %s\n", reponame)
	} else if cresp.RefType == "branch" {
		fmt.Printf("Created new branch %s\n", reponame)
	} else if cresp.RefType == "tag" {
		fmt.Printf("Created new tag %s\n", reponame)
	} else {
		return fmt.Errorf("unable to parse, reference type is empty")
	}

	return nil
}

func parseDeleteEvent(payload json.RawMessage, reponame string) error {
	var cresp DeleteEvent
	if err := json.Unmarshal(payload, &cresp); err != nil {
		panic(err)
	}

	if cresp.RefType == "branch" {
		fmt.Printf("Deleted branch %s\n", reponame)
	} else if cresp.RefType == "tag" {
		fmt.Printf("Deleted tag %s\n", reponame)
	} else {
		return fmt.Errorf("unable to parse, reference type is empty")
	}

	return nil
}

func main() {
	username := os.Args[1]

	url := "https://api.github.com/users/" + username + "/events"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode)

	if err != nil {
		panic(err)
	}

	var cresp Event
	if err := json.NewDecoder(resp.Body).Decode(&cresp); err != nil {
		panic(err)
	}

	for _, event := range cresp {
		if event.Type == "CreateEvent" {
			err := parseCreateEvent(event.Payload, event.Repo.Name)
			if err != nil {
				panic(err)
			}
		} else if event.Type == "DeleteEvent" {
			err := parseDeleteEvent(event.Payload, event.Repo.Name)
			if err != nil {
				panic(err)
			}
		} else if event.Type == "IssuesEvent" {
			fmt.Printf("IssuesEvent: %s\n", event.Repo.Name)
		} else if event.Type == "PullRequestEvent" {
			fmt.Printf("PullRequestEvent: %s\n", event.Repo.Name)
		} else if event.Type == "PushEvent" {
			fmt.Printf("PushEvent: %s\n", event.Repo.Name)
		} else if event.Type == "ReleaseEvent" {
			fmt.Printf("ReleaseEvent: %s\n", event.Repo.Name)
		}
	}

}
