run:
	direnv allow
	docker compose up --build

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
	--url "postgresql://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@localhost:5432/$${POSTGRES_DB}?search_path=public&sslmode=disable" \
	--to "file://./database/atlas.hcl" \


ecr_region = $$(terraform -chdir=infrastructure/terraform/ output -raw region)
ecr_url = $$(terraform -chdir=infrastructure/terraform/ output -raw ecr_url)
push_image: 
	aws ecr get-login-password --region $(ecr_region) | docker login --username AWS --password-stdin $(ecr_url)
	docker build --tag polls-app:latest .
	docker tag polls-app:latest $(ecr_url)):polls-app
	docker push $(ecr_url)):polls-app