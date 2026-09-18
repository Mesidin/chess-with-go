#!/bin/sh
# Install Chess with Go binary and desktop entry for Omarchy / Linux main app menu.
set -e
cd "$(dirname "$0")/.."

mkdir -p ~/.local/bin
mkdir -p ~/.local/share/applications

if [ -f dist/chess-with-go-linux ]; then
    cp dist/chess-with-go-linux ~/.local/bin/chess-with-go
else
    go build -o ~/.local/bin/chess-with-go ./cmd/chess-with-go
fi

chmod +x ~/.local/bin/chess-with-go
cp assets/chess-with-go.desktop ~/.local/share/applications/

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database ~/.local/share/applications/
fi

echo "Installed Chess with Go to ~/.local/bin/chess-with-go"
echo "Added desktop shortcut to ~/.local/share/applications/chess-with-go.desktop"
echo "Chess with Go should now appear in your Omarchy / Linux main apps menu!"
