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

# Build binary for current OS/Arch
build:
  fyne package --name {{ BINARY_NAME }}

# Run all go:generate directives
generate:
  go generate
