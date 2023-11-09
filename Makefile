.PHONY: all help build clean clean-test unit-tests int-test test package clean-glide glide swagger-json swagger-doc clean-mocks certs clean-certs

all: build

help:
	@echo ""
	@echo "Available tasks:"
	@echo "    build                 Build the Hiota Agent Binary for Linux Distributions"
	@echo "    docker                Build the Hiota Agent Binary for Linux Distributions"
	@echo "    test                  Run Unit and Integration Tests"
	@echo ""

build: ## Build the Hiota Agent binary file for Linux Distributions
	rm -rf vendor
	go mod tidy
	go build main.go

docker:
	docker build . -t    app-server
