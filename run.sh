#!/bin/bash

# Load environment variables from .env file
export $(cat .env | xargs)

# Run the server
go run cmd/main.go 