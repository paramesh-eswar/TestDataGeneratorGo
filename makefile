# Variables
APP_NAME = TestDataGenerator

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	APP_NAME = TestDataGenerator
endif
ifeq ($(OS),Windows_NT)
	APP_NAME = TestDataGenerator.exe
endif
GO_FILES = $(wildcard *.go)
GO_CMD = go
GO_BUILD = $(GO_CMD) build
GO_TEST = $(GO_CMD) test
GO_RUN = $(GO_CMD) run

# Default target: build
all: build

# Build target: compile the Go source files into the final executable
build: $(APP_NAME)

# Compile the Go source files
$(APP_NAME): $(GO_FILES)
	$(GO_BUILD) -o $(APP_NAME) $(GO_FILES)

# Run the application
run: build
	$(GO_RUN) ./... $(APP_NAME)

# Test target: run the tests for the project
test:
	$(GO_TEST) -v ./... -bench=. -cover

# Declare phony targets
.PHONY: all build test
