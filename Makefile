# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Name of the executable
BINARY_NAME=gopl
INSTALL_DIR=/usr/local/bin

all: clean build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

install:
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)"
	cp $(BINARY_NAME) $(INSTALL_DIR)
	@echo "Adding $(INSTALL_DIR) to your PATH in .bashrc"
	echo 'export PATH=$$PATH:$(INSTALL_DIR)' >> ~/.zshrc
	source ~/.zshrc
