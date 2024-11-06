set quiet := true

MAIN_PACKAGE_PATH := "."
BINARY_NAME := "app-builder"

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

# Build binary for current OS/Arch
build:
  go build -o=./build/{{ BINARY_NAME }} {{ MAIN_PACKAGE_PATH }}

# Run all go:generate directives
generate:
  go generate
