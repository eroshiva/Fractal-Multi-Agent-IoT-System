export GO111MODULE=on

FMAIS_VERSION := main
DOCKER_REPOSITORY := eroshiva/
GOLANGCI_LINTERS_VERSION := v1.56.2
BENCH_TIME := 10s

.PHONY: help build
help: # Credits to https://gist.github.com/prwhite/8168133 for this handy oneliner
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build the Go binaries and run all validations (default)
	go mod tidy
	go mod vendor
	go build -mod=vendor -o build/_output/fractal-mais ./cmd/fractal-mais

install: build ## Install newly build package in a local environment - now it can be used as a local command
	cd cmd/fractal-mais/ && go install && cd ../../

install-gobenchdata: ## Installs gobench data tool which parses go bench data into file
	go install go.bobheadxi.dev/gobenchdata@latest

example: build ## Runs a unit test, which generates a random system model and plots a graph to showcase it
	./build/_output/fractal-mais --example

bench: build ## Benchmark the codebase in classic way (measure time of the function execution)
	./build/_output/fractal-mais --benchmark --hardcoded

bench-sm: build ## Benchmark the codebase in a classic way (measure time of the function execution)
	./build/_output/fractal-mais --benchFMAIS --hardcoded

bench-rm: build ## Benchmark the ME-ERT-CORE Reliability Model (both, per definition and optimized) in a classic way (measure time of the function execution)
	./build/_output/fractal-mais --benchMeErtCORE --hardcoded
	./build/_output/fractal-mais --benchMeErtCOREoptimized --maxNumInstances 101 --appNumber 101

bench-rm-optimized: build ## Benchmark the ME-ERT-CORE (optimized) Reliability Model in a classic way (measure time of the function execution)
	./build/_output/fractal-mais --benchMeErtCOREoptimized --maxNumInstances 101 --appNumber 101

gobench: build install-gobenchdata ## Benchmark the codebase with gobench
	go test -bench . -benchmem ./... -timeout 0m -benchtime=${BENCH_TIME} -count=10 | gobenchdata --json ./data/benchmarks.json
	gobenchdata web generate ./data
	cd data && gobenchdata web serve

figures: generate-figures generate-joint-figure generate-joint-reliability-figure ## Generates all figures from the time complexity estimation and measurement

generate-figures: build ## Generates figures based on the benchmarked data. It needs an exact name of the file carrying data!
	./build/_output/fractal-mais --generateFigures benchmark_fmais_2023-03-26_01-59-08.json --generateFigures benchmark_meertcore_2023-04-04_21-14-18.csv

generate-joint-figure: build ## Generates joint figure based on the benchmarked data. It needs an exact name of the file carrying data!
	./build/_output/fractal-mais --generateJointFigure benchmark_meertcore_optimized_2023-08-05_07-50-20.json --generateJointFigure benchmark_meertcore_per_definition_2023-08-05_07-50-20.json

generate-joint-reliability-figure: build ## Generates joint figure based on the benchmarked data. It needs an exact name of the file carrying data!
	./build/_output/fractal-mais --meertcore --generateJointFigure me-ert-core_fmais_depth_2.json --generateJointFigure me-ert-core_fmais_depth_4.json


linters-install: ## Install linters locally for verification
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINTERS_VERSION}

linters: linters-install ## Perform linting to verify codebase
	golangci-lint run --timeout 5m

test: build linters ## Test the codebase
	go test -coverprofile=./data/test-cover.out -race gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/...

run: example ## Runs compiled binary and generates an example of Fractal MAIS

clean: ## Remove all the build artifacts
	rm -rf ./build/_output ./vendor ./data/assets/ ./data/gobenchdata-web.yml ./data/index.html ./data/overrides.css ./data/benchmarks.json ./data/test-cover.out
	rm -rf ./figures/measurement_*.eps ./figures/measurement_*.png
	go clean -cache -testcache

image: ## Builds a Docker image
	docker build --platform linux/amd64 . -f build/fractal-mais/Dockerfile \
		-t ${DOCKER_REPOSITORY}fractal-mais-generator:${FMAIS_VERSION}

docker-bench: image ## Benchmarks the whole codebase wrapped in a Docker container
	docker run --rm -v ~/go/src/gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/data:/usr/local/bin/data \
		-v ~/go/src/gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/figures:/usr/local/bin/figures \
		${DOCKER_REPOSITORY}fractal-mais-generator:${FMAIS_VERSION} --benchmark --hardcoded --docker

measurement: build ## Runs measurement for FMAIS of depth 2, 3 and 4 and large-scale FMAIS measurement
	./build/_output/fractal-mais --runMeasurement

docker-measurement: image ## Runs measurement in a Docker container
	docker run --rm -v ~/go/src/gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/data:/usr/local/bin/data \
		-v ~/go/src/gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/figures:/usr/local/bin/figures \
		${DOCKER_REPOSITORY}fractal-mais-generator:${FMAIS_VERSION} --runMeasurement

kind: image ## Builds and image and uploads to kind cluster
	kind load docker-image ${DOCKER_REPOSITORY}fractal-mais-generator:${FMAIS_VERSION}

helm-lint: ## Validates correctness of the helm charts
	helm lint ./charts

helm-render: ## Renders helm charts to check if all values are propagated correctly
	helm template ./charts

helm-check: ## Dry runs a helm charts to check if charts are deployable
	helm install -n fractal-mais --dry-run fmais ./charts

helm-deploy: ## Deploys Fractal MAIS in a running Kubernetes cluster and performs measurement
	kubectl create namespace fractal-mais
	helm install -n fractal-mais fmais ./charts

helm-uninstall: ## Uninstalls deployed charts
	helm uninstall -n fractal-mais fmais
	kubectl delete namespace fractal-mais
