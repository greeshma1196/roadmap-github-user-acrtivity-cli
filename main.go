package main

import (
	"fmt"
	"net/http"
	"os"
)

func checkStatusCode(statusCode int) error {

	if statusCode == 200 || statusCode == 304 {
		fmt.Printf("Success\n")
	} else if statusCode == 403 {
		return fmt.Errorf("Forbidden\n")
	} else if statusCode == 404 {
		return fmt.Errorf("Not Found\n")
	} else if statusCode == 503 {
		return fmt.Errorf("Service Unavailable\n")
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

}
