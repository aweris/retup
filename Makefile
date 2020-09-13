# Ensure everything works even if GOPATH is not set, which is often the case.
GOPATH       ?= $(shell go env GOPATH)
GOBIN        ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO           ?= $(shell which go)

VERSION      := $(strip $(shell [ -d .git ] && git describe --abbrev=0))
LONG_VERSION := $(strip $(shell [ -d .git ] && git describe --always --tags --dirty))
BUILD_DATE   := $(shell date -u +"%Y-%m-%dT%H:%M:%S%Z")
VCS_REF      := $(strip $(shell [ -d .git ] && git rev-parse HEAD))

GO_PACKAGES   = $(shell go list ./... | grep -v -E '/vendor/|/test')
GO_FILES     := $(shell find . -type f -name '*.go' -not -path './vendor/*')

GOBUILD      := $(GO) build -mod=vendor
GOINSTALL    := $(GO) install -mod=vendor
GOMOD        := $(GO) mod
GOCLEAN      := $(GO) clean
GOTEST       := $(GO) test
GOFMT        := gofmt
GOLANGCILINT := $(GOBIN)/golangci-lint
LDFLAGS      := '-s -w -X main.version=$(VERSION) -X main.commit=$(VCS_REF) -X main.date=$(BUILD_DATE)'
TAGS         := netgo

# Project directories
ROOT_DIR     := $(CURDIR)
BUILD_DIR    := $(ROOT_DIR)/build

# Helper variables
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

default: help

.PHONY: retup
retup: ## Runs retup target
retup: vendor $(BUILD_DIR) main.go $(wildcard *.go) $(wildcard */*.go); $(info $(M) running retup )
	$(Q) CGO_ENABLED=0 $(GOBUILD) -a -ldflags $(LDFLAGS) -tags $(TAGS) -o $(BUILD_DIR)/$@ .

.PHONY: vendor
vendor: ## Updates vendored copy of dependencies
vendor: ; $(info $(M) running go mod vendor)
	$(Q) $(GOMOD) tidy
	$(Q) $(GOMOD) vendor

.PHONY: clean
clean: ## Cleanup everything
clean: ; $(info $(M) cleaning )
	$(Q) $(GOCLEAN)
	$(Q) $(shell rm -rf $(BUILD_DIR))

.PHONY: lint
lint: ## Runs golangci-lint analysis
lint: $(GOLANGCI_LINT) ; $(info $(M) running lint )
	$(Q) $(GOLANGCILINT) run -v --enable-all --skip-dirs tmp -c .golangci.yml

.PHONY: fix
fix: ## Runs golangci-lint fix
fix: $(GOLANGCILINT) fmt ; $(info $(M) running fix )
	$(Q) $(GOLANGCILINT) run --fix --enable-all --skip-dirs tmp -c .golangci.yml

.PHONY: fmt
fmt: ## Runs gofmt
fmt: ; $(info $(M) running format )
	$(Q) $(GOFMT) -w -s $(GO_FILES)

.PHONY: test
test: ## Runs go test
test: ; $(info $(M) runnig tests)
	$(Q) $(GOTEST) -race -cover -v ./...

.PHONY: install-tools
install-tools: ## Install tools
install-tools: vendor ; $(info $(M) installing tools)
	$(Q) $(GOINSTALL) $(shell cat tools.go | grep _ | awk -F'"' '{print $$2}')

.PHONY: help
help: ## shows this help message
	$(Q) awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m\t %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

$(BUILD_DIR): ; $(info $(M) creating build directory)
	$(Q) $(shell mkdir -p $@)