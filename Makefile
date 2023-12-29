# Makefile

# Binary name
BINARY=b7s-attributes

# Build the project
build:
	go build -o $(BINARY) cmd/b7s-attributes/main.go

# Clean up
clean:
	rm -f $(BINARY)

.PHONY: build clean
