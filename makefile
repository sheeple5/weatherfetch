BINARY_NAME := weatherfetch
BIN_DIR := ~/.local/bin

.PHONY: all build install clean

all: build

build:
	go build -o $(BINARY_NAME) .

install:
	go build -o $(BINARY_NAME) .
	mkdir -p $(BIN_DIR)
	mv $(BINARY_NAME) $(BIN_DIR)

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BIN_DIR)/$(BINARY_NAME)
