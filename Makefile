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
	@docker build -t loan-repayment-calculator:v1.0.0 .