#!/bin/bash

# Build CSS first
npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify


# Enable CGO for SQLite support
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o my-journal