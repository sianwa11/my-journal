#!/bin/bash

# Enable CGO for SQLite support
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o my-journal