export GO111MODULE=on

.PHONY: build

GOLANGCI_LINTERS_VERSION := v1.51.1

build: # @HELP build the Go binaries and run all validations (default)
build:
	go mod tidy
	go build -o build/_output/fractal-mas ./cmd/fractal-mas

linters-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINTERS_VERSION}

linters: linters-install
	golangci-lint run --timeout 5m

test: build linters
	go test -race gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/...

clean:: # @HELP remove all the build artifacts
	rm -rf ./build/_output
	go clean -testcache gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/...
