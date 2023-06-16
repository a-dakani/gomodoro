CURR_DIR:=$(shell pwd)
BIN_DIR=/usr/local/bin

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean

# Main package and executable name
PACKAGE = ./cmd/gomodoro-cli
EXECUTABLE = gomodoro
BUILD_PATH = ./build

# Run target
run: clean build
	$(BUILD_PATH)/$(EXECUTABLE)

# Build target
build: clean
	$(GOBUILD) -o $(BUILD_PATH)/$(EXECUTABLE) $(PACKAGE)

# Clean target
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_PATH)

install: build
	sudo cp $(BUILD_PATH)/$(EXECUTABLE) $(BIN_DIR)/$(EXECUTABLE)

uninstall:
	sudo rm -rf $(BIN_DIR)/$(EXECUTABLE)

update:
	git pull origin main
	make install