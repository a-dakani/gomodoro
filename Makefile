CURR_DIR:=$(shell pwd)
BIN_DIR=/usr/local/bin

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean

# Main package and executable name
PACKAGE_CLI = ./cmd/gomodoro-cli
PACKAGE_API = ./cmd/gomodoro-api
EXECUTABLE_CLI = gomodoro-cli
EXECUTABLE_API = gomodoro-api
BUILD_PATH = ./build

# Run CLI target
run-cli: clean build-cli
	$(BUILD_PATH)/$(EXECUTABLE_CLI)

# Build CLI target
build-cli: clean
	$(GOBUILD) -o $(BUILD_PATH)/$(EXECUTABLE) $(PACKAGE_CLI)

# Run CLI target with reflex to hot reload
watch-api:
	ulimit -n 1000
	reflex -s -r '\.go$$' make run-api

# Run API target
run-api: clean build-api
	$(BUILD_PATH)/$(EXECUTABLE_API)

# Build API target
build-api: clean
	$(GOBUILD) -o $(BUILD_PATH)/$(EXECUTABLE) $(PACKAGE_API)

# Clean target
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_PATH)


# Install CLI target
install: build-cli
	sudo cp $(BUILD_PATH)/$(EXECUTABLE_CLI) $(BIN_DIR)/$(EXECUTABLE_CLI)

# Uninstall CLI target
uninstall-cli:
	sudo rm -rf $(BIN_DIR)/$(EXECUTABLE_CLI)

# Update CLI target
update-cli:
	git pull origin main
	make uninstall-cli
	make install




