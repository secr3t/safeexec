#!/bin/bash
set -e

PROJECT_ROOT=$(pwd)
ASSETS_DIR="$PROJECT_ROOT/internal/assets"
mkdir -p "$ASSETS_DIR"

PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64" "windows/amd64")

echo "Building watchdog for multiple platforms to internal/assets..."

for PLATFORM in "${PLATFORMS[@]}"; do
    OS=${PLATFORM%/*}
    ARCH=${PLATFORM#*/}
    
    OUTPUT_NAME="watchdog-$OS-$ARCH"
    if [ "$OS" == "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
    fi
    
    echo "  - Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o "$ASSETS_DIR/$OUTPUT_NAME" ./cmd/watchdog
done

echo "Done. Binaries are in $ASSETS_DIR"
