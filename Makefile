VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
# GOBASE := $(shell pwd)
# GOPATH := $(GOBASE)/vendor:$(GOBASE)
# GOBIN := $(GOBASE)/bin
# GOFILES := $(wildcard *.go)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID := /tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## build: Build the binary.
build:
	echo "  >  Building binary..."
	go build $(LDFLAGS) -o bin/autosigner -mod=vendor

## clean: Clean build files. Runs `go clean` internally.
clean:
	echo "  >  Cleaning build cache"
	rm bin/autosigner 2> /dev/null
	go clean

## test: Run tests
test:
	echo "  >  Running tests"
	go test -v -mod=vendor ./...

.PHONY: build clean test help
all: help
help: Makefile
	@echo
	@echo " Choose a command run for make:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
