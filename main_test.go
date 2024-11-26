package main

import (
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
		require.EqualError(t, err, "Forbidden\n")
	})

	t.Run("Successfully validates not found status code", func(t *testing.T) {
		statusCode := 404
		err := checkStatusCode(statusCode)
		require.EqualError(t, err, "Not Found\n")
	})

	t.Run("Successfully validates service unavailable status code", func(t *testing.T) {
		statusCode := 503
		err := checkStatusCode(statusCode)
		require.EqualError(t, err, "Service Unavailable\n")
	})
}
