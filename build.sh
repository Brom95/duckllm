#!/bin/bash

# Define the Go project directory
PROJECT_DIR=$(pwd)

# Define the output directory for the builds
OUTPUT_DIR="$PROJECT_DIR/builds"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Define the platforms you want to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/386"
)

# Loop through each platform and build the project
for PLATFORM in "${PLATFORMS[@]}"; do
    OS=${PLATFORM%/*}
    ARCH=${PLATFORM#*/}

    # Set the output binary name
    OUTPUT_BINARY="$OUTPUT_DIR/duckllm-$OS-$ARCH"

    # Append .exe for Windows platforms
    if [[ "$OS" == "windows" ]]; then
        OUTPUT_BINARY+=".exe"
    fi

    # Build the Go project
    echo "Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT_BINARY" "$PROJECT_DIR"

    if [ $? -ne 0 ]; then
        echo "Failed to build for $OS/$ARCH"
        exit 1
    fi
done

echo "Builds completed successfully!"
