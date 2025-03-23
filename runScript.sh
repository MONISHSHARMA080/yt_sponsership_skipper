#!/bin/bash

# Define the paths
ROOT_DIR=$(pwd)
WEBSITE_DIR="${ROOT_DIR}/website"

# Check if website directory exists
if [ ! -d "$WEBSITE_DIR" ]; then
    echo "Error: Website directory not found at $WEBSITE_DIR"
    exit 1
fi

# First, open a terminal with the Go server
gnome-terminal --tab --title="Go Server" --working-directory="$ROOT_DIR" -- bash -c "echo 'Starting Go server...'; go run .; exec bash" \
    --tab --title="Frontend" --working-directory="$WEBSITE_DIR" -- bash -c "echo 'Starting frontend...'; nrd; exec bash" \
    --tab --title="Ngrok" --working-directory="$ROOT_DIR" -- bash -c "echo 'Starting ngrok tunnel...'; ngrok http 8080; exec bash"

echo "Development environment started!"