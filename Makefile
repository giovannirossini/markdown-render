.PHONY: build install clean test

# Build the binary
build:
	go build -o mdrender

# Install to GOPATH/bin
install:
	go install

# Clean build artifacts
clean:
	rm -f mdrender

# Quick test with stdin
test: build
	@echo "# Test\nThis is a **test** with \`code\`" | ./mdrender

