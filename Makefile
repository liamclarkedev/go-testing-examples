# MallocNanoZone is a Intel mac issue, see https://github.com/golang/go/issues/49138
test:
	@MallocNanoZone=0 go test -race ./...

lint:
	@golangci-lint run ./...

cover:
	@go test -coverprofile cover.out ./...

report:
	@go tool cover -html=cover.out
