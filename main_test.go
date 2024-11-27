package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckStatusCode(t *testing.T) {
	t.Run("Successfully validates success status code", func(t *testing.T) {
		statusCode := 200
		err := checkStatusCode(statusCode)
		require.Nil(t, err)
	})

	t.Run("Successfully validates forbidden status code", func(t *testing.T) {
		statusCode := 403
		err := checkStatusCode(statusCode)
		require.EqualError(t, err, "forbidden")
	})

	t.Run("Successfully validates not found status code", func(t *testing.T) {
		statusCode := 404
		err := checkStatusCode(statusCode)
		require.EqualError(t, err, "not found")
	})

	t.Run("Successfully validates service unavailable status code", func(t *testing.T) {
		statusCode := 503
		err := checkStatusCode(statusCode)
		require.EqualError(t, err, "service unavailable")
	})
}

func TestParseCreateEvent(t *testing.T) {
	t.Run("Successfully validates CreateEvent for repository", func(t *testing.T) {
		payload := `{"ref_type": "repository"}`
		reponame := "sample_repo"
		err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
	})

	t.Run("Successfully validates CreateEvent for branch", func(t *testing.T) {
		payload := `{"ref_type": "branch"}`
		reponame := "sample_repo"
		err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
	})

	t.Run("Successfully validates CreateEvent for tag", func(t *testing.T) {
		payload := `{"ref_type": "tag"}`
		reponame := "sample_repo"
		err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.Nil(t, err)
	})

	t.Run("Successfully validates error for CreateEvent", func(t *testing.T) {
		payload := `{"ref_type": ""}`
		reponame := "sample_repo"
		err := parseCreateEvent(json.RawMessage(payload), reponame)
		require.EqualError(t, err, "unable to parse, reference type is empty")
	})
}
