SHELL := /bin/bash
DATE ?= $(shell date '+%Y-%m-%d')
BASE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION ?= $(shell git describe --tags --always --match=* 2> /dev/null)
VERSION_HASH = $(shell git rev-parse --short HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

go = GOGC=off go
MODULE = $(shell env GO111MODULE=on go list -m)

LDFLAGS += -X "$(MODULE)/common.Version=$(VERSION)" -X "$(MODULE)/common.CommitSHA=$(VERSION_HASH)" -X "$(MODULE)/common.BuildDate=$(DATE)"

.PHONY: build
build: ## Build
	$Q $(go) build -ldflags '$(LDFLAGS)' -o cw
