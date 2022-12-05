BINARY_NAME=test-inspector

# linting
LINTER_CONFIG := .golangci.yaml

# Linker
PACKAGE := "test-inspector"
VERSION := "$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH := "$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP := $(date '+%Y-%m-%dT%H:%M:%S')
# SUPABASE_KEY := ""

LDFLAGS := -ldflags "-X '${PACKAGE}/cmd.Version=${VERSION}' -X '${PACKAGE}/cmd.CommitHash=${COMMIT_HASH}' -X '${PACKAGE}/cmd.BuildTime=${BUILD_TIMESTAMP}' -X '${PACKAGE}/cmd.SupabaseKey=${SUPABASE_KEY}'"
LDFLAGS_DEV := -ldflags "-X '${PACKAGE}/cmd.SupabaseKey=${SUPABASE_KEY}'"

.PHONY: build
build: ## build the project
	GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o ./bin/darwin/amd64/${BINARY_NAME} main.go
	GOARCH=arm64 GOOS=darwin go build ${LDFLAGS} -o ./bin/darwin/arm64/${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o ./bin/linux/amd64/${BINARY_NAME} main.go
	GOARCH=arm64 GOOS=linux go build ${LDFLAGS} -o ./bin/linux/arm64/${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o ./bin/windows/amd64/${BINARY_NAME}.exe main.go
	GOARCH=arm64 GOOS=windows go build ${LDFLAGS} -o ./bin/windows/arm64/${BINARY_NAME}.exe main.go

.PHONY: build_dev
build_dev: ## build the project
	GOARCH=amd64 GOOS=darwin go build ${LDFLAGS_DEV} -o ./bin/darwin/amd64/${BINARY_NAME} main.go
	GOARCH=arm64 GOOS=darwin go build ${LDFLAGS_DEV} -o ./bin/darwin/arm64/${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=linux go build ${LDFLAGS_DEV} -o ./bin/linux/amd64/${BINARY_NAME} main.go
	GOARCH=arm64 GOOS=linux go build ${LDFLAGS_DEV} -o ./bin/linux/arm64/${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=windows go build ${LDFLAGS_DEV} -o ./bin/windows/amd64/${BINARY_NAME}.exe main.go
	GOARCH=arm64 GOOS=windows go build ${LDFLAGS_DEV} -o ./bin/windows/arm64/${BINARY_NAME}.exe main.go

.PHONY: run
run: ## run the darwin arm64 binary
	@./bin/darwin/arm64/${BINARY_NAME}

.PHONY: build_and_run
build_and_run: build run

.PHONY: clean
clean: ## clean built binaries
	@go clean
	@rm ./bin/darwin/arm64/${BINARY_NAME}
	@rm ./bin/linux/arm64/${BINARY_NAME}
	@rm ./bin/windows/arm64/${BINARY_NAME}.exe
	@rm ./bin/darwin/amd64/${BINARY_NAME}
	@rm ./bin/linux/amd64/${BINARY_NAME}
	@rm ./bin/windows/amd64/${BINARY_NAME}.exe

.PHONY: test
test: ## run tests
	@go test ./... -race

.PHONY: test_coverage
test_coverage: ## get coverage report
	@go test ./... -coverprofile=coverage.out

.PHONY: dep
dep: ## download dependencies
	@go mod download

.PHONY: vet
vet: ## run go vet on all packages
	@go vet

.PHONY: lint
lint: ## run all the lint tools
	@golangci-lint run

.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
