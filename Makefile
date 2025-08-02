APP_NAME := harvest
BINARY_NAME := $(APP_NAME)
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

BUILD_DIR := ./bin

INSTALL_DIR := /usr/local/bin

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o "bin/$(BINARY_NAME)" .

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo mkdir -p $(INSTALL_DIR)
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully!"
	@echo "You can now run: $(BINARY_NAME)"

uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled successfully!"

check-install:
	@echo "Checking if $(BINARY_NAME) is installed..."
	@which $(BINARY_NAME) || echo "$(BINARY_NAME) not found in PATH"
	@$(BINARY_NAME) --version 2>/dev/null || echo "$(BINARY_NAME) not working or no --version flag"

install-dev: build
	@echo "Creating development symlink for $(BINARY_NAME)..."
	@mkdir -p $(HOME)/.local/bin
	@ln -sf $(PWD)/$(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "Development symlink created!"
	@echo "Changes to your binary will be reflected immediately"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  install       - Install to system (requires sudo)"
	@echo "  install-go    - Install using 'go install'"
	@echo "  install-dev   - Create development symlink"
	@echo "  uninstall     - Remove from system"
	@echo "  check-install - Check if app is installed"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help"

.PHONY: build install install-go install-dev uninstall check-install clean help
