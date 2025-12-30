.PHONY: help docker-build docker-test docker-run docker-shell test tidy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

docker-build: ## Build the Docker image with PostgreSQL
	docker build -t nullable-postgres-test -f tests/Dockerfile .

docker-run: docker-build ## Build and run with interactive output
	docker run --rm -it nullable-postgres-test

docker-shell: docker-build ## Open a shell in the container
	docker run --rm -it nullable-postgres-test /bin/sh

test: docker-build ## Run all Go tests (including database tests with Docker)
	docker run --rm -t nullable-postgres-test

tidy: ## Tidy Go modules
	go mod tidy
