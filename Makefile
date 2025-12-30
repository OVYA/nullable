.PHONY: help test tidy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Run all tests (including PostgreSQL integration tests)
	cd tests && go test -v ./...

tidy: ## Tidy Go modules
	go mod tidy
	cd tests && go mod tidy
