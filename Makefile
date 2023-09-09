test:
	go test ./...

test_coverage:
	go test ./... -cover

dep:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	curl -sSf https://atlasgo.sh | sh

vet:
	go vet ./...

lint:
	golangci-lint run --enable-all

db_migrate:
	atlas schema apply \
	--url "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?search_path=public&sslmode=disable" \
	--to "file://./database/atlas.hcl" \
