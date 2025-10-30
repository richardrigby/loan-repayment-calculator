gen:
	docker build -f Dockerfile.buf -t buf-with-plugins . && docker run --rm -v "$(PWD)":/workspace buf-with-plugins generate

run: ## Run the application natively
	@go run cmd/loan-repayment-calculator/main.go

run-docker: ## Build and run the application in a Docker container
	@docker compose up --build -d loan-repayment-calculator

stop: ## Stop running containers
	@docker compose stop

test: ## Run tests
	@go test ./...

clean: ## Stop and remove containers, networks, volumes, and images created by `up`.
	@docker compose down --volumes --remove-orphans

build: ## Build the Docker image
	@COMMIT_SHA=$$(git rev-parse --short HEAD); echo $$COMMIT_SHA; perl -i.bak -pe "s|(image:\s*\S+):\S+|\\1:$$COMMIT_SHA|g" deployments/k8s/deployment.yaml; rm -f deployments/k8s/deployment.yaml.bak; docker build -t loan-repayment-calculator:$$COMMIT_SHA .

run-k8s: ## Deploy the application to a local Kubernetes cluster
	@kubectl apply -f deployments/k8s/

clean-k8s: ## Remove the application from the local Kubernetes cluster
	@kubectl delete -f deployments/k8s/

port-forward: ## Forward ports for local access to the Kubernetes service
	@kubectl port-forward svc/loan-repayment-calculator 50051:50051 50052:50052

.PHONY: gen run run-docker stop test clean build run-k8s clean-k8s