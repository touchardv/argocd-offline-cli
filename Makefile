BUILD_DIR := $(shell pwd)/build
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
SOURCES := $(shell find . -name '*.go')

BINARY := argocd-offline-cli-$(GOOS)-$(GOARCH)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.DEFAULT_GOAL := build
.PHONY: build
build: $(BUILD_DIR)/$(BINARY)

$(BUILD_DIR)/$(BINARY): $(BUILD_DIR) $(SOURCES)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY) ./cmd

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	go clean ./...

.PHONY: run
run: $(BUILD_DIR)/$(BINARY)
	$(BUILD_DIR)/$(BINARY) appset

.PHONY: test
test:
	go test -v -cover -timeout 10s ./...
