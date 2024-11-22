set quiet := true

MAIN_PACKAGE_PATH := "."
BINARY_NAME := "builder"

[private]
help:
    just --list --unsorted

# Run dev build
run: 
  go run .

# Run dev build with debug screens
debug:
  go run -tags debug .

# Fetch and organize dependencies
tidy:
  go mod tidy

test:
  go test -v ./...

# Build binary for current OS/Arch
build:
  fyne package 

# Build specifically for iOS simulator
build-ios-simulator:
  fyne package -os iossimulator

# Build for web via WASM
build-web:
  fyne package -os web

# Serve the WASM version of the app locally
serve:
  fyne serve
  
# Run all go:generate directives
generate:
  go generate
