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

createblockchain: build
	if [ -z "$(ADDRESS)" ]; then echo "Usage: make createblockchain ADDRESS='your address'"; exit 1; fi
	$(BUILD_DIR)/$(BINARY_NAME) createblockchain -address "$(ADDRESS)"

getbalance: build
	if [ -z "$(ADDRESS)" ]; then echo "Usage: make getbalance ADDRESS='your address'"; exit 1; fi
	$(BUILD_DIR)/$(BINARY_NAME) getbalance -address "$(ADDRESS)"

send: build
	if [ -z "$(FROM)" ] || [ -z "$(TO)" ] || [ -z "$(AMOUNT)" ]; then \
		echo "Usage: make send FROM='source address' TO='destination address' AMOUNT='amount to send'"; exit 1; \
	fi
	$(BUILD_DIR)/$(BINARY_NAME) send -from "$(FROM)" -to "$(TO)" -amount "$(AMOUNT)"

.PHONY: build run addblock printchain createblockchain getbalance send