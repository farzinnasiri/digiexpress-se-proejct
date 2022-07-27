# constants
APP_NAME := dlocator
APP_VERSION := 0.0.1
OWNER := digiexpress
ROOT := github.com/$(OWNER)/$(APP_NAME)

# api config
PROTO_DIR := api/proto
STUBS_DIR := pkg/api
PROTOC ?= protoc
PROTOC_OPTIONS ?= -I.
BINARY_NAME=dlocator
BINARY_PATH=bin

# linter config
LINTER_VERSION = v1.47.0

SOURCES = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*" -not -path "*.pb.go"))
PACKAGES := $(shell go list ./... | grep -v /vendor)
PROTOS = $(patsubst ./%,%,$(shell find $(PROTO_DIR) -name "*.proto"))
PBS = $(patsubst ./%,%,$(shell find $(STUBS_DIR) -name "*.pb.go"))

help: ## to display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

run: ## to start the application
	@go run ./cmd/$(APP_NAME)
.PHONY: run

.pre-check-protoc: ## to install the protoc compiler for golang
	@if [ -z "$$(which protoc-gen-go)" ]; then go get -v github.com/golang/protobuf/protoc-gen-go; fi
.PHONY: .pre-check-protoc

protoc: | .pre-check-protoc ## to generate all necessary protobuf files
	@#$(foreach proto,$(PROTOS),$(PROTOC) $(PROTOC_OPTIONS) --go_out=plugins=grpc:$(STUBS_DIR)/$(word 1,$(subst /, ,$(subst $(PROTO_DIR), ,$(proto)))) $(proto)${\n})
	@$(foreach proto,$(PROTOS),$(PROTOC) $(PROTOC_OPTIONS) --go_out=./internal/pkg/api/ --go-grpc_out=./internal/pkg/api/ api/proto/v1/*.proto)
.PHONY: protoc

bench:
	@go test -bench=. -count 5 -run=^# ./benchmarks

wire: | .pre-check-wire ## to generate wire file for dependency injection
	@wire ./internal/app/$(APP_NAME)
.PHONY: wire

.pre-check-wire:
	@if [ -z "$$(which wire)" ]; then \
  		go get -v github.com/google/wire/cmd/wire; \
  	fi
.PHONY: .pre-check-wire

devtools: .pre-check-wire .pre-check-test-tools .pre-check-protoc .pre-check-formatting-tools ## to install all the development tools
	$(info wire, gocov, protoc, gci, gofumpt are all installed)
.PHONY: devtools

dependencies: ## to install the dependencies
	@go mod tidy -compat=1.17
	@go mod download
	@go mod vendor
.PHONY: dependencies

clean: ## to remove generated files
	@-rm -rf $(APP_NAME)
	@-find . -type d -name mocks -exec rm -rf \{} +
.PHONY: clean

.pre-check-formatting-tools:
	@if [ -z "$$(which gofumpt)" ]; then go install mvdan.cc/gofumpt@latest; fi
	@#if [ -z "$$(which gci)" ]; then go get github.com/daixiang0/gci; fi
.PHONY: .pre-check-formatting-tools

fmt: | .pre-check-formatting-tools ## to run go formatter on all source codes across the project
	@gofumpt -l -w ${SOURCES}
	@#gci write ${SOURCES}
.PHONY: fmt

.bin/golangci-lint:
	@if [ -z "$$(which golangci-lint)" ]; then\
 		curl --verbose -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin/ $(LINTER_VERSION);\
 	else\
 		$(info golangci-lint was not found. going to install version: $(LINTER_VERSION))\
 		mkdir -p .bin;\
 		ln -s "$$(which golangci-lint)" $@;\
 	fi

lint: .bin/golangci-lint ## to lint the source files
	@.bin/golangci-lint run --config=.golangci-lint.yml ./internal/... ./cmd/...
.PHONY: lint

.pre-check-test-tools:
	@if [ -z "$$(which gocov)" ]; then go get -v github.com/axw/gocov/gocov; fi
.PHONY: .pre-check-test-tools

test: | .pre-check-test-tools ## to run all tests
	@go test ./...
.PHONY: test

coverage: coverage.cover coverage.html ## to run tests and generate test coverage data
	@gocov convert $< | gocov report
.PHONY: coverage

coverage.html: coverage.cover ## to output coverage data to an html
	@go tool cover -html=$< -o $@
.PHONY: coverage.html

coverage.cover: $(SRCS) $(PBS) Makefile
	@-rm -rfv .coverage
	@mkdir -p .coverage
	@$(foreach pkg,$(PACKAGES),go test -timeout 30s -short -covermode=count -coverprofile=.coverage/$(subst /,-,$(pkg)).cover $(pkg)${\n})
	@echo "mode: count" > $@
	@grep -h -v "^mode:" .coverage/*.cover >> $@
.PHONY: coverage.cover

.remove_empty_dirs: ## to remove empty directories across all project
	@-find . -type d -print | xargs rmdir 2>/dev/null | true
.PHONY: .remove_empty_dirs


build:
	@echo "Cleanning..."
	@echo $(BINARY_PATH)
	@rm -rf $(BINARY_PATH)
	@mkdir $(BINARY_PATH)
	@echo "Building DLocator..."
	@go build -o $(BINARY_PATH) -v ./...
	@echo "Building is finished"



# helpers

# a variable containing a new line e.g. ${\n} would emit a new line useful in foreach functions
define \n


endef
