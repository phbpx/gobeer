SHELL := /bin/bash

# ==============================================================================
# help

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# ==============================================================================
# Setup

GOBIN := $(shell go version)

check.go:
	@go version >/dev/null 2>&1 || (echo "ERROR: go is not installed" && exit 1)

## Install go tools
setup.dev: check.go
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# ==============================================================================
# Test

## Execute tests
test:
	go test ./... -count=1 -coverprofile=coverage.out -v

## Execute tests with coverage visualization
cover: test
	go tool cover -html=coverage.out

## Execute static check
lint:
	staticcheck -checks=all ./...

# ==============================================================================
# Dev

## Run local database
run.db:
	docker-compose up -d
