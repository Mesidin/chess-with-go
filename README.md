# Chess with Go

A chess game written in Go. Same idea as [Go with Go](../go-with-go): a rules engine plus a native window.

Human vs human or a simple bot, legal moves (castling, en passant, promotion to queen, check, mate, stalemate), and a learn-to-play course. PGN, clocks, and tactics puzzles are later work.

© 2026 Bradley Erickson

## Requirements

Go 1.26 or later. Ebitengine v2.10 is pure Go on desktop (no C compiler).

Menus and HUD follow the **active Omarchy theme** (`colors.toml` under `~/.local/state/omarchy/current/theme`). Off Omarchy, that chrome is [Harbor Dark](https://github.com/HANCORE-linux/omarchy-harbordark-theme).

The **board and pieces** are a separate skin (Classic, Forest, Harbor, Gold rush, Ocean, Neon) under **Board skins**. Omarchy does not recolor the squares.

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

**Start game** picks Humans, White vs bot, or Black vs bot. Click a piece, then a highlighted square. **U** undo (vs the bot, your move and the reply). Promotion is to queen for now.

**Learn to play** walks through pawns, knights, bishops, capture, check, castling, and mate. Highlighted squares are the ones to use.

Pieces are Colin M.L. Burnett’s Staunton set (the Wikipedia chess diagrams), [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/).

## Later (not in this start)

PGN save/load, clocks, a bot, learn-to-play, tactics puzzles, terminal UI — the same shape as Go with Go, in later sessions so weekly usage stays bounded.
