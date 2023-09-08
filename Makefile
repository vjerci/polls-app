test:
	go test ./...

test_coverage:
	go test ./... -cover

dep:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

vet:
	go vet ./...

lint:
	golangci-lint run --enable-all