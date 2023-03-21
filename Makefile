MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := bin/nake.$(UNAME)
.ONESHELL:

UNAME := $(shell uname -s | tr A-Z a-z)
GIT_REF := $(shell git describe --match="" --always --dirty=+)
GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF) 2> /dev/null)
PACKAGE := $(shell go list)

.PHONY: help
help:  ## Show this help
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ##/:/g' | column -t -s:

.git/hooks/pre-commit:  ## Install pre-commit checks
	pre-commit install

.PHONY: check
check: .git/hooks/pre-commit ## Run precommit checks
	pre-commit run --all

.PHONY: test
test:  ## Run go test
	go test -v ./...

bin/nake.darwin:  ## Build the application binary for current OS

bin/nake.%:  ## Build the application binary for target OS, for example bin/nake.linux
	GOOS=$* go build -o $@ -ldflags "-X $(PACKAGE)/version=$(GIT_TAG)+$(GIT_REF)" main.go

.PHONY: install
install: bin/nake.$(UNAME) ## Install the binary
	rm -f ~/.local/bin/nake
	cp $< ~/.local/bin/nake

.PHONY: fixtures
fixtures:
	go test ./... -generate
