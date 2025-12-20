.PHONY: build install deploy clean test

all: build install deploy clean

# Build the binary
build:
	go build -o mdrender

# Install to GOPATH/bin
install:
	go install

# Move to /usr/local/bin
deploy:
	sudo mv mdrender /usr/local/bin

# Clean build artifacts
clean:
	rm -f mdrender

# Quick test with stdin
test: build
	@echo "# Test\nThis is a **test** with \`code\`" | ./mdrender

