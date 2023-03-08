export GO111MODULE=on

.PHONY: build

FMAS_VERSION := latest
DOCKER_REPOSITORY := eroshiva/
GOLANGCI_LINTERS_VERSION := v1.51.1
BENCH_TIME := 10s

build: # @HELP build the Go binaries and run all validations (default)
	go mod tidy
	go mod vendor
	go build -mod=vendor -o build/_output/fractal-mas ./cmd/fractal-mas

install: # @HELP install newly build package in a local environment - now it can be used as a local command
install: build
	cd cmd/fractal-mas/ && go install && cd ../../

install-gobenchdata: # @HELP installs gobench data tool which parses go bench data into file
	go install go.bobheadxi.dev/gobenchdata@latest

example: # @HELP runs a unit test, which generates a random system model and plots a graph to showcase it
example: build
	./build/_output/fractal-mas --example


bench: # @HELP benchmark the codebase in classic way measure time of the function execution
bench: build
	./build/_output/fractal-mas --benchFMAS --hardcoded

# ToDo - build infrastructure with parsing around it..
bench-with-Docker: # @HELP benchmark the codebase wrapped in a Docker container
bench-with-Docker: image
	docker run -it ${DOCKER_REPOSITORY}fractal-mas-generator:${FMAS_VERSION}

# ToDo - build infrastructure with parsing around it..
gobench: # @HELP benchmark the codebase with gobench
gobench: build install-gobenchdata
#	go test -v -bench=. ./... -cpu=4 -count=100 -benchtime=${BENCH_TIME} -benchmem -timeout 0m | gobenchdata --json ./data/benchmarks.json
	go test -bench . -benchmem ./... -timeout 0m -benchtime=${BENCH_TIME} -count=10 | gobenchdata --json ./data/benchmarks.json
	gobenchdata web generate ./data
	cd data && gobenchdata web serve

generate_figures: # @HELP generates figures based on the benchmarked data. It needs an exact name of the file carrying data!
generate_figures: build
	./build/_output/fractal-mas --generateFigures benchmark_2023-03-08_07:32:08.json

linters-install: # @HELP install linters locally for verification
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINTERS_VERSION}

linters: # @HELP perform linting to verify codebase
linters: linters-install
	golangci-lint run --timeout 5m

test: # @HELP test the codebase
test: build linters
	go test -race -count=100 gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/...

# @HELP runs compiled binary and generates an example of Fractal MAIS
run: example

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./data/assets/ ./data/gobenchdata-web.yml ./data/index.html ./data/overrides.css
	go clean -cache -testcache

# ToDo - fix Dockerfile
image: # @HELP builds a Docker image
	docker build --platform linux/amd64 . -f build/fractal-mas/Dockerfile \
		-t ${DOCKER_REPOSITORY}fractal-mas-generator:${FMAS_VERSION}
