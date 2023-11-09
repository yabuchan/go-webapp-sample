.PHONY: all help build clean clean-test unit-tests int-test test package clean-glide glide swagger-json swagger-doc clean-mocks certs clean-certs

all: build

help:
	@echo ""
	@echo "Available tasks:"
	@echo "    build                 Build the Hiota Agent Binary for Linux Distributions"
	@echo "    docker                Build the Hiota Agent Binary for Linux Distributions"
	@echo "    compose-set           Set log level. specify one of [debug, info, error, fatal, dpanic, panic].  e.g. make set debug"
	@echo "    compose-up            Start docker compose"
	@echo "    compose-down          Stop docker compose"
	@echo "    kube-set          	 Set log level. specify one of [debug, info, error, fatal, dpanic, panic].  e.g. make set debug"
	@echo ""

build: ## Build the Hiota Agent binary file for Linux Distributions
	rm -rf vendor
	go mod tidy
	go build main.go

docker:
	docker build . -t    app-server

compose-up:
	docker compose up

compose-down:
	docker compose down

compose-set:
	 @docker exec  go-webapp-sample_app-server_1 /bin/bash -c  'echo $(level) >  /logLevel'

kube-set:
	 pod=$(kubectl get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | grep app-server`);
	 kubectl exec $(pod) /bin/bash -c  'echo $(level) >  /logLevel';