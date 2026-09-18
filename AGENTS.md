# Chess with Go — agent notes

A chess game written in Go. Legal moves (castling, en passant, promotion to queen, check, mate, stalemate), Staunton piece sprites, a simple bot, and a learn-to-play course.

PGN, clocks, and tactics puzzles are later work unless the user asks.

## Layout

- `engine/` — rules. No UI imports.
- `bot/` — 1-ply capture-aware sparring partner.
- `learn/` — interactive beginner course.
- `ui/` — Ebitengine window and Cburnett piece PNGs (`ui/pieces/`).
- `internal/theme/` — Omarchy palette loader.
- `internal/meta/` — copyright and title.
- `cmd/chess-with-go/` — entry.

## Run

```
go test ./...
go run ./cmd/chess-with-go
```

On Windows PowerShell use `.\cmd\chess-with-go`. Do not kill a running window to rebuild unless the user asks.

## Omarchy theming

Menus, HUD, and buttons follow the **active Omarchy theme**. The chessboard and piece colors are a separate **board skin** (Classic, Forest, Harbor, Gold rush, Ocean, Neon) chosen in **Board skins**. Do not drive square or piece colors from Omarchy.

1. Read `colors.toml` (else `alacritty.toml`) from, in order:
   - `$OMARCHY_THEME_DIR`
   - `~/.local/state/omarchy/current/theme`
   - `~/.config/omarchy/current/theme`
   - `~/.config/omarchy/themes/<theme.name>`
2. If none of those exist, use the baked-in **Harbor Dark** palette (`internal/theme.HarborDark`, [HANCORE-linux/omarchy-harbordark-theme](https://github.com/HANCORE-linux/omarchy-harbordark-theme)).
3. Reload about once a second so `omarchy theme set` retints a running game.

Piece **sprites** stay Cburnett Staunton; skins tint them. There is no halo behind black pieces.

## Product rules

- Home: Start game, Learn to play, Board skins, Quit. Start opens Humans / White vs bot / Black vs bot.
- Click a piece, then a legal square. **U** undo (vs bot, two plies). Promotion is queen until someone asks otherwise.
- Copyright: © 2026 Bradley Erickson (`internal/meta`).
- Piece art: Colin M.L. Burnett / Wikipedia, CC BY-SA 3.0 (`ui/pieces/NOTICE.txt`).

## Tests

`go test ./...` must stay green. Engine tests cover start legal-move count, fool’s mate, scholar’s mate, and castling.
