export GO111MODULE=on

.PHONY: build

FMAIS_VERSION := latest
DOCKER_REPOSITORY := eroshiva/
GOLANGCI_LINTERS_VERSION := v1.52.2
BENCH_TIME := 10s

build: ## @HELP build the Go binaries and run all validations (default)
	go mod tidy
	go mod vendor
	go build -mod=vendor -o build/_output/fractal-mais ./cmd/fractal-mais

install: build ## @HELP install newly build package in a local environment - now it can be used as a local command
	cd cmd/fractal-mais/ && go install && cd ../../

install-gobenchdata: ## @HELP installs gobench data tool which parses go bench data into file
	go install go.bobheadxi.dev/gobenchdata@latest

example: build ## @HELP runs a unit test, which generates a random system model and plots a graph to showcase it
	./build/_output/fractal-mais --example


bench: build ## @HELP benchmark the codebase in classic way (measure time of the function execution)
	./build/_output/fractal-mais --benchmark --hardcoded

bench-sm: build ## @HELP benchmark the codebase in a classic way (measure time of the function execution)
	./build/_output/fractal-mais --benchFMAIS --hardcoded

bench-rm: build ## @HELP benchmark the ME-ERT-CORE Reliability Model in a classic way (measure time of the function execution)
	./build/_output/fractal-mais --benchMeErtCORE --hardcoded

gobench: build install-gobenchdata ## @HELP benchmark the codebase with gobench
	go test -bench . -benchmem ./... -timeout 0m -benchtime=${BENCH_TIME} -count=10 | gobenchdata --json ./data/benchmarks.json
	gobenchdata web generate ./data
	cd data && gobenchdata web serve

generate-figures: build ## @HELP generates figures based on the benchmarked data. It needs an exact name of the file carrying data!
	./build/_output/fractal-mais --generateFigures benchmark_fmas_2023-03-26_01:59:08.json --generateFigures benchmark_meertcore_2023-04-04_21:14:18.csv

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

image: ## @HELP builds a Docker image
	docker build --platform linux/amd64 . -f build/fractal-mais/Dockerfile \
		-t ${DOCKER_REPOSITORY}fractal-mais-generator:${FMAIS_VERSION}

docker-bench: image ## @HELP benchmarks the whole codebase wrapped in a Docker container
	docker run --rm -v ~/go/src/gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/data:/usr/local/bin/data \
		-v ~/go/src/gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/figures:/usr/local/bin/figures \
		${DOCKER_REPOSITORY}fractal-mais-generator:${FMAIS_VERSION} --benchmark --hardcoded --docker
