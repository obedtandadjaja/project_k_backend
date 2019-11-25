test:
	go clean; \
  go build; \
  go test -v -race ./...
