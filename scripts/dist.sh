#!/bin/sh
# Build portable binaries for Mac (Intel + Apple Silicon) and Linux.
set -e
cd "$(dirname "$0")/.."
mkdir -p dist
export CGO_ENABLED=0

build() {
	os=$1
	arch=$2
	out=$3
	echo "→ dist/$out"
	GOOS=$os GOARCH=$arch go build -trimpath -ldflags="-s -w" -o "dist/$out" ./cmd/chess-with-go
}

build darwin amd64 chess-with-go-mac-intel
build darwin arm64 chess-with-go-mac-apple
build linux amd64 chess-with-go-linux

echo
echo "Copy one file onto the target machine and make it executable:"
echo "  Mac (Apple Silicon):  dist/chess-with-go-mac-apple"
echo "  Mac (Intel):          dist/chess-with-go-mac-intel"
echo "  Omarchy / Linux:      dist/chess-with-go-linux"
echo
ls -lh dist
