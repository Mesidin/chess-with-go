package engine

import "fmt"

type Color uint8

const (
	White Color = iota
	Black
)

func (c Color) String() string {
	if c == White {
		return "White"
	}
	return "Black"
}

func (c Color) Opponent() Color {
	if c == White {
		return Black
	}
	return White
}

type Kind uint8

const (
	Empty Kind = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

func (k Kind) String() string {
	switch k {
	case Pawn:
		return "pawn"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Rook:
		return "rook"
	case Queen:
		return "queen"
	case King:
		return "king"
	default:
		return "empty"
	}
}

type Piece struct {
	Color Color
	Kind  Kind
}

func (p Piece) Is(c Color, k Kind) bool { return p.Kind == k && p.Color == c }

func (p Piece) Letter() string {
	if p.Kind == Empty {
		return ""
	}
	var s string
	switch p.Kind {
	case Pawn:
		s = "P"
	case Knight:
		s = "N"
	case Bishop:
		s = "B"
	case Rook:
		s = "R"
	case Queen:
		s = "Q"
	case King:
		s = "K"
	}
	if p.Color == Black {
		return fmt.Sprintf("%s", toLower(s))
	}
	return s
}

func toLower(s string) string {
	if s == "" {
		return s
	}
	b := s[0]
	if b >= 'A' && b <= 'Z' {
		return string(b + 32)
	}
	return s
}

// Square is 0..63, a1 = 0, h1 = 7, a8 = 56.
func File(sq int) int { return sq % 8 }
func Rank(sq int) int { return sq / 8 }

func Square(file, rank int) int { return rank*8 + file }

func SquareName(sq int) string {
	if sq < 0 || sq > 63 {
		return "--"
	}
	return string(rune('a'+File(sq))) + string(rune('1'+Rank(sq)))
}

func ParseSquare(s string) (int, bool) {
	if len(s) != 2 {
		return 0, false
	}
	f := int(s[0] - 'a')
	r := int(s[1] - '1')
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return 0, false
	}
	return Square(f, r), true
}
