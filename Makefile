GOOPTS := GO111MODULE=on GOARCH=amd64 CGO_ENABLED=1
BINARY_NAME = packer-plugin-register-ami
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor bin generate-mocks

all: help

build: linux windows darwin ## Build your project and put the output binary in bin/

linux: clean modules bin generate
	$(GOOPTS) GOOS=linux go build -o bin/$(BINARY_NAME)_linux_amd64

windows: clean modules bin generate
	$(GOOPTS) GOOS=windows go build -o bin/$(BINARY_NAME)_windows_amd64.exe

darwin: clean modules bin generate
	$(GOOPTS) GOOS=darwin go build -o bin/$(BINARY_NAME)_darwin_amd64

test: ## Run the tests of the project
	$(GOOPTS) go test -v -race ./...


vendor: ## Copy of all packages needed to support builds and tests in the vendor directory
	go mod vendor

lint: lint-releaser lint-go ## Run all available linters

lint-releaser: ## Use goreleaser/goreleaser to check .goreleaser.yml
	docker run --rm -v $(shell pwd):/app -w /app goreleaser/goreleaser check -f /app/.goreleaser.yml

lint-go: ## Use golintci-lint on your project
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s

modules: ## Install required GO modules
	go mod download

clean: ## Remove build related file
	rm -fr ./bin

tools: ## Install generate dependency
	go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@latest
	go install github.com/vektra/mockery/v2@latest

generate: tools ## Generate hcl2spec.go and mocks
	go generate ./...
	mockery --name SSMAPI --dir ./vendor/github.com/aws/aws-sdk-go/service/ssm/ssmiface	


bin:
	mkdir -p bin
	rm -f bin/*

help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)