export GO111MODULE=on

.PHONY: build

FMAS_VERSION := latest
DOCKER_REPOSITORY := eroshiva/
GOLANGCI_LINTERS_VERSION := v1.51.1
BENCH_TIME := 10s

build: ## @HELP build the Go binaries and run all validations (default)
	go mod tidy
	go mod vendor
	go build -mod=vendor -o build/_output/fractal-mas ./cmd/fractal-mas

install: build ## @HELP install newly build package in a local environment - now it can be used as a local command
	cd cmd/fractal-mas/ && go install && cd ../../

install-gobenchdata: ## @HELP installs gobench data tool which parses go bench data into file
	go install go.bobheadxi.dev/gobenchdata@latest

example: build ## @HELP runs a unit test, which generates a random system model and plots a graph to showcase it
	./build/_output/fractal-mas --example


bench: build ## @HELP benchmark the codebase in classic way (measure time of the function execution)
	./build/_output/fractal-mas --benchmark --hardcoded

bench_sm: build ## @HELP benchmark the codebase in a classic way (measure time of the function execution)
	./build/_output/fractal-mas --benchFMAS --hardcoded

bench_rm: build ## @HELP benchmark the ME-ERT-CORE Reliability Model in a classic way (measure time of the function execution)
	./build/_output/fractal-mas --benchMeErtCORE --hardcoded

# ToDo - build infrastructure with parsing around it..
bench-with-Docker: ## @HELP benchmark the codebase wrapped in a Docker container
bench-with-Docker: image
	docker run -it ${DOCKER_REPOSITORY}fractal-mas-generator:${FMAS_VERSION}

gobench: build install-gobenchdata ## @HELP benchmark the codebase with gobench
	go test -bench . -benchmem ./... -timeout 0m -benchtime=${BENCH_TIME} -count=10 | gobenchdata --json ./data/benchmarks.json
	gobenchdata web generate ./data
	cd data && gobenchdata web serve

generate_figures: build ## @HELP generates figures based on the benchmarked data. It needs an exact name of the file carrying data!
	./build/_output/fractal-mas --generateFigures benchmark_2023-03-12_10:55:56.json

linters-install: ## @HELP install linters locally for verification
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINTERS_VERSION}

linters: linters-install ## @HELP perform linting to verify codebase
	golangci-lint run --timeout 5m

test: build linters ## @HELP test the codebase
	go test -race -count=100 gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/...

run: example ## @HELP runs compiled binary and generates an example of Fractal MAIS

clean: ## @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./data/assets/ ./data/gobenchdata-web.yml ./data/index.html ./data/overrides.css ./data/benchmarks.json
	go clean -cache -testcache

image: build ## @HELP builds a Docker image
	docker build --platform linux/amd64 . -f build/fractal-mas/Dockerfile \
		-t ${DOCKER_REPOSITORY}fractal-mas-generator:${FMAS_VERSION}
