# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean

# Main package and executable name
PACKAGE = ./cmd/gomodoro
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