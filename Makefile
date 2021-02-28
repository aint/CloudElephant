COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags |cut -d- -f1)

LDFLAGS = -ldflags "-X github.com/aint/CloudElephant/cmd.gitTag=${TAG} -X github.com/aint/CloudElephant/cmd.gitCommit=${COMMIT} -X github.com/aint/CloudElephant/cmd.gitBranch=${BRANCH}"

.PHONY: help deps build test clean

help: ## Display this help screen
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Download the dependencies
	go mod download

build: ## Build the executable
	deps
	go build ${LDFLAGS} -v -o ce .

clean: ## Remove the executable
	rm ce

test: ## Run tests
	go test -v ./...

lint: deps ## Lint the source code
	golangci-lint run --timeout 5m -E golint

