package engine

import (
	"fmt"
	"strconv"
	"strings"
)

func pieceFromLetter(r rune) (Piece, bool) {
	c := White
	if r >= 'a' && r <= 'z' {
		c = Black
		r -= 32
	}
	var k Kind
	switch r {
	case 'P':
		k = Pawn
	case 'N':
		k = Knight
	case 'B':
		k = Bishop
	case 'R':
		k = Rook
	case 'Q':
		k = Queen
	case 'K':
		k = King
	default:
		return Piece{}, false
	}
	return Piece{Color: c, Kind: k}, true
}

// ParseFEN builds a game from a FEN string (placement through side to move at minimum).
func ParseFEN(fen string) (*Game, error) { // placement / side / castle / ep / half / full
	parts := strings.Fields(fen)
	if len(parts) < 1 {
		return nil, fmt.Errorf("empty FEN")
	}
	g := &Game{ep: -1, full: 1}
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("FEN needs 8 ranks")
	}
	for i, row := range ranks {
		rank := 7 - i
		file := 0
		for _, r := range row {
			if r >= '1' && r <= '8' {
				file += int(r - '0')
				continue
			}
			p, ok := pieceFromLetter(r)
			if !ok || file > 7 {
				return nil, fmt.Errorf("bad FEN square on rank %d", rank)
			}
			g.board[Square(file, rank)] = p
			file++
		}
		if file != 8 {
			return nil, fmt.Errorf("rank %d has %d files", rank, file)
		}
	}
	g.toMove = White
	if len(parts) > 1 && parts[1] == "b" {
		g.toMove = Black
	}
	if len(parts) > 2 && parts[2] != "-" {
		for _, r := range parts[2] {
			switch r {
			case 'K':
				g.castle |= WhiteKingside
			case 'Q':
				g.castle |= WhiteQueenside
			case 'k':
				g.castle |= BlackKingside
			case 'q':
				g.castle |= BlackQueenside
			}
		}
	}
	if len(parts) > 3 && parts[3] != "-" {
		if sq, ok := ParseSquare(parts[3]); ok {
			g.ep = sq
		}
	}
	if len(parts) > 4 {
		if n, err := strconv.Atoi(parts[4]); err == nil {
			g.half = n
		}
	}
	if len(parts) > 5 {
		if n, err := strconv.Atoi(parts[5]); err == nil {
			g.full = n
		}
	}
	return g, nil
}
