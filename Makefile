.DEFAULT_GOAL=help

build-binaries: ## build binaries for supported platforms
	goreleaser build --snapshot --clean

test: ## test tracer with example.com
	@echo "-----------------------------------------------------\n"
	@go run cmd/gotrace/*.go myip
	@echo "\n-----------------------------------------------------"
	@go run cmd/gotrace/*.go route example.com

help: ## show this help
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ": "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
