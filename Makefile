# Makefile

BINARY_NAME = goblockchain
BUILD_DIR = ./bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME)

run: build
	$(BUILD_DIR)/$(BINARY_NAME)

addblock: build
	if [ -z "$(DATA)" ]; then echo "Usage: make addblock DATA='your data'"; exit 1; fi
	$(BUILD_DIR)/$(BINARY_NAME) addblock -data "$(DATA)"

printchain: build
	$(BUILD_DIR)/$(BINARY_NAME) printchain

.PHONY: build run addblock printchain
