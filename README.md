# Chess with Go

A chess game written in Go. Same idea as [Go with Go](../go-with-go): a rules engine plus a native window.

This is a **start**, not a full twin of Go with Go. Human vs human on one machine, legal moves (including castling, en passant, promotion to queen, check, mate, stalemate). PGN, a bot, a learn mode, and puzzles are later work.

© 2026 Bradley Erickson

## Requirements

Go 1.26 or later. Ebitengine v2.10 is pure Go on desktop (no C compiler).

## Run

From this folder:

Windows PowerShell:

```
go test ./...
go run .\cmd\chess-with-go
```

macOS and Omarchy:

```
go test ./...
go run ./cmd/chess-with-go
```

On Omarchy the window uses X11 (XWayland). If it does not open: `sudo pacman -S go libx11 libglvnd mesa`.

Click a piece, then a highlighted square. **U** undo. Promotion is to queen for now.

## Later (not in this start)

PGN save/load, clocks, a bot, learn-to-play, tactics puzzles, terminal UI — the same shape as Go with Go, in later sessions so weekly usage stays bounded.
