build:
	go build -o github-activity
test:
	go test ./... -coverprofile cover.out -v
lint: 
	golangci-lint run -v