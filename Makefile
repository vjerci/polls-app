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
	--to "file://./database/atlas.hcl"

# for infra_deploy you need installed aws_cli, terraform, atlas, helm and yq
# running this command will spin up eks cluster with fully functional polls application (be carefull because infrastructure is not free)
infra_deploy:
	echo "be carefull because infrastructure is not free, make sure to run \"make infra_destroy\" after you are done testing the app to destroy provisioned infrastructure"
	echo "setting up infrastructure"
	make infra_init
	echo "waiting 100s for ecrs to come online"
	sleep 100
	make push_api
	make push_front
	make kubctl_login
	make helm_deploy
	echo "waiting 200s for postgre to come online"
	sleep 200 
	make db_migrate_kubernetes
	make get_endpoint

infra_destroy:
	terraform -chdir=infrastructure/terraform apply -destroy

infra_init:
	terraform -chdir=infrastructure/terraform apply

commit_id=$$(git rev-parse  HEAD)
# these variables should be kept in ci secrets (they are dynamic for ease of deployment)
deployment_region = $$(terraform -chdir=infrastructure/terraform/ output -raw region)
cluster_name = $$(terraform -chdir=infrastructure/terraform/ output -raw cluster_name)
api_ecr_url = $$(terraform -chdir=infrastructure/terraform/ output -raw api_ecr_url)
push_api: 
	aws ecr get-login-password --region $(deployment_region) | docker login --username AWS --password-stdin $(api_ecr_url)
	docker buildx build \
		--platform linux/amd64 \
		--tag polls-api:$(commit_id) \
		.
	docker tag polls-api:$(commit_id) $(api_ecr_url):$(commit_id)
	docker push $(api_ecr_url):$(commit_id)

front_ecr_url = $$(terraform -chdir=infrastructure/terraform/ output -raw front_ecr_url)
push_front:
	aws ecr get-login-password --region $(deployment_region) | docker login --username AWS --password-stdin $(front_ecr_url)
	docker buildx build \
		--platform linux/amd64 \
		--tag polls-front:$(commit_id) \
		--build-arg NEXT_PUBLIC_API_URL="/api" \
		--build-arg NEXT_PUBLIC_GOOGLE_CLIENT_ID=$(GOOGLE_CLIENT_ID) \
		./front
	docker tag polls-front:$(commit_id) $(front_ecr_url):$(commit_id)
	docker push $(front_ecr_url):$(commit_id)

kubctl_login:
	 aws eks --region $(deployment_region) update-kubeconfig --name $(cluster_name)

jwt_key = $$(terraform -chdir=infrastructure/terraform/ output -raw jwt_key)
db_password = $$(terraform -chdir=infrastructure/terraform/ output -raw db_password)
helm_deploy:
	helm install --atomic polls \
		--set pod.api.repositoryUrl="${api_ecr_url}" \
		--set pod.api.dockerTag="${commit_id}" \
		--set pod.front.repositoryUrl="${front_ecr_url}" \
		--set pod.front.dockerTag="${commit_id}" \
		--set pod.api.JWTKey="${jwt_key}" \
		--set pod.db.password="${db_password}" \
		 ./infrastructure/helm

# these variables should be kept in ci secrets (they are dynamic for ease of deployment)
postgres_user=$$(yq .pod.db.user < ./infrastructure/helm/values.yaml)
postgres_db=$$(yq .pod.db.dbName < ./infrastructure/helm/values.yaml)
db_migrate_kubernetes:
	unset DEBUG && \
	echo "migrating db on cluster${cluster_name}" && \
	kubectl port-forward service/postgres 5432:5432 & \
	sleep 5 && \
	atlas schema apply \
		--url "postgresql://${postgres_user}:${db_password}@localhost:5432/${postgres_db}?search_path=public&sslmode=disable" \
		--to "file://./database/atlas.hcl" && \
	kill $$!

get_endpoint:
	echo "Your endpoint is:"
	kubectl get --no-headers=true ingress | awk -F " " '{print $$4}'
