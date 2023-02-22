export GO111MODULE=on

.PHONY: build

GOLANGCI_LINTERS_VERSION := v1.51.1
BENCH_TIME := 10s

build: # @HELP build the Go binaries and run all validations (default)
build:
	go mod tidy
	go mod vendor
	go build -mod=vendor -o build/_output/fractal-mas ./cmd/fractal-mas

# ToDo - build infrastructure with parsing around it..
bench: # @HELP benchmark the codebase in classic way measure time of the function execution
bench:

# ToDo - build infrastructure with parsing around it..
gobench: # @HELP benchmark the codebase with gobench
gobench:
	go test -v -bench=. ./... -count=100 -run=^# -benchtime=${BENCH_TIME} -benchmem
	# there is a room to parse output of benchmarking and process graphically

linters-install: # @HELP install linters locally for verification
linters-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINTERS_VERSION}

linters: # @HELP perform linting to verify codebase
linters: linters-install
	golangci-lint run --timeout 5m

test: # @HELP test the codebase
test: build linters
	go test -race -count=100 gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/...

clean:: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor
	go clean -testcache gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/...
