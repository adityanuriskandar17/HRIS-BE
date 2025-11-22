#!/bin/bash

# Script to regenerate swagger docs and build the application
# Used by Air for auto-reload

# Add go bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Regenerate swagger docs
echo "Regenerating Swagger docs..."
swag init -g cmd/api/main.go -o docs

# Build the application
echo "Building application..."
go build -o ./tmp/main ./cmd/api

