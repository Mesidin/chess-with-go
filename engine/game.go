package engine

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoPiece     = errors.New("no piece on that square")
	ErrWrongSide   = errors.New("not that side's turn")
	ErrIllegal     = errors.New("illegal move")
	ErrGameOver    = errors.New("game is over")
	ErrNothingUndo = errors.New("nothing to undo")
)

type Move struct {
	From, To    int
	Promo       Kind
	Capture     Piece
	Castle      bool
	EnPassant   bool
	DoublePawn  bool
	PrevEP      int
	PrevCastle  castle
	PrevHalf    int
}

func (m Move) String() string {
	s := SquareName(m.From) + SquareName(m.To)
	if m.Promo != Empty {
		s += strings.ToLower(kindLetter(m.Promo))
	}
	return s
}

func kindLetter(k Kind) string {
	switch k {
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	case Queen:
		return "Q"
	default:
		return ""
	}
}

type castle uint8

const (
	WhiteKingside castle = 1 << iota
	WhiteQueenside
	BlackKingside
	BlackQueenside
)

type Result uint8

const (
	InProgress Result = iota
	WhiteWin
	BlackWin
	Draw
)

func (r Result) String() string {
	switch r {
	case WhiteWin:
		return "White wins"
	case BlackWin:
		return "Black wins"
	case Draw:
		return "Draw"
	default:
		return "In progress"
	}
}

type Game struct {
	board   [64]Piece
	toMove  Color
	castle  castle
	ep      int // -1 none
	half    int
	full    int
	history []Move
	result  Result
	reason  string
}

func New() *Game {
	g := &Game{toMove: White, ep: -1, full: 1, castle: WhiteKingside | WhiteQueenside | BlackKingside | BlackQueenside}
	place := func(sq int, c Color, k Kind) { g.board[sq] = Piece{Color: c, Kind: k} }
	for f := 0; f < 8; f++ {
		place(Square(f, 1), White, Pawn)
		place(Square(f, 6), Black, Pawn)
	}
	back := []Kind{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for f, k := range back {
		place(Square(f, 0), White, k)
		place(Square(f, 7), Black, k)
	}
	return g
}

func (g *Game) At(sq int) Piece {
	if sq < 0 || sq > 63 {
		return Piece{}
	}
	return g.board[sq]
}

func (g *Game) AtFileRank(file, rank int) Piece { return g.At(Square(file, rank)) }
func (g *Game) ToMove() Color                   { return g.toMove }
func (g *Game) Result() Result                  { return g.result }
func (g *Game) Reason() string                  { return g.reason }
func (g *Game) Over() bool                      { return g.result != InProgress }
func (g *Game) FullMove() int                   { return g.full }
func (g *Game) Moves() []Move                   { return append([]Move(nil), g.history...) }

func (g *Game) Play(from, to int, promo Kind) error {
	if g.Over() {
		return ErrGameOver
	}
	legals := g.LegalMoves()
	var chosen *Move
	for i := range legals {
		m := &legals[i]
		if m.From == from && m.To == to {
			if promo == Empty && m.Promo != Empty {
				promo = Queen
			}
			if m.Promo == promo || (m.Promo == Empty && promo == Empty) {
				chosen = m
				if m.Promo == promo {
					break
				}
			}
		}
	}
	if chosen == nil {
		return ErrIllegal
	}
	g.apply(*chosen)
	g.history = append(g.history, *chosen)
	g.updateResult()
	return nil
}

func (g *Game) Undo() error {
	if len(g.history) == 0 {
		return ErrNothingUndo
	}
	m := g.history[len(g.history)-1]
	g.history = g.history[:len(g.history)-1]
	g.unapply(m)
	g.result = InProgress
	g.reason = ""
	return nil
}

func (g *Game) apply(m Move) {
	p := g.board[m.From]
	g.board[m.From] = Piece{}
	if m.EnPassant {
		dir := -8
		if p.Color == Black {
			dir = 8
		}
		g.board[m.To+dir] = Piece{}
	}
	if m.Castle {
		if m.To == Square(6, Rank(m.From)) { // kingside
			g.board[Square(5, Rank(m.From))] = g.board[Square(7, Rank(m.From))]
			g.board[Square(7, Rank(m.From))] = Piece{}
		} else {
			g.board[Square(3, Rank(m.From))] = g.board[Square(0, Rank(m.From))]
			g.board[Square(0, Rank(m.From))] = Piece{}
		}
	}
	if m.Promo != Empty {
		p.Kind = m.Promo
	}
	g.board[m.To] = p
	g.clearCastle(m, p)
	g.ep = -1
	if m.DoublePawn {
		g.ep = (m.From + m.To) / 2
	}
	if p.Kind == Pawn || m.Capture.Kind != Empty || m.EnPassant {
		g.half = 0
	} else {
		g.half++
	}
	if p.Color == Black {
		g.full++
	}
	g.toMove = p.Color.Opponent()
}

func (g *Game) unapply(m Move) {
	p := g.board[m.To]
	if m.Promo != Empty {
		p.Kind = Pawn
	}
	g.board[m.To] = Piece{}
	g.board[m.From] = p
	if m.EnPassant {
		dir := -8
		if p.Color == Black {
			dir = 8
		}
		cap := Piece{Color: p.Color.Opponent(), Kind: Pawn}
		g.board[m.To+dir] = cap
	} else if m.Capture.Kind != Empty {
		g.board[m.To] = m.Capture
	}
	if m.Castle {
		if m.To == Square(6, Rank(m.From)) {
			g.board[Square(7, Rank(m.From))] = g.board[Square(5, Rank(m.From))]
			g.board[Square(5, Rank(m.From))] = Piece{}
		} else {
			g.board[Square(0, Rank(m.From))] = g.board[Square(3, Rank(m.From))]
			g.board[Square(3, Rank(m.From))] = Piece{}
		}
	}
	g.castle = m.PrevCastle
	g.ep = m.PrevEP
	g.half = m.PrevHalf
	g.toMove = p.Color
	if p.Color == Black {
		g.full--
	}
}

func (g *Game) clearCastle(m Move, p Piece) {
	off := func(c castle) { g.castle &^= c }
	if p.Kind == King {
		if p.Color == White {
			off(WhiteKingside | WhiteQueenside)
		} else {
			off(BlackKingside | BlackQueenside)
		}
	}
	if p.Kind == Rook {
		switch m.From {
		case 0:
			off(WhiteQueenside)
		case 7:
			off(WhiteKingside)
		case 56:
			off(BlackQueenside)
		case 63:
			off(BlackKingside)
		}
	}
	switch m.To {
	case 0:
		off(WhiteQueenside)
	case 7:
		off(WhiteKingside)
	case 56:
		off(BlackQueenside)
	case 63:
		off(BlackKingside)
	}
}

func (g *Game) updateResult() {
	if len(g.LegalMoves()) > 0 {
		if g.half >= 100 {
			g.result = Draw
			g.reason = "50-move rule"
		}
		return
	}
	if g.InCheck(g.toMove) {
		if g.toMove == White {
			g.result = BlackWin
		} else {
			g.result = WhiteWin
		}
		g.reason = "checkmate"
		return
	}
	g.result = Draw
	g.reason = "stalemate"
}

func (g *Game) InCheck(c Color) bool {
	k := g.kingSq(c)
	if k < 0 {
		return false
	}
	return g.attacked(k, c.Opponent())
}

func (g *Game) kingSq(c Color) int {
	for sq := 0; sq < 64; sq++ {
		if g.board[sq].Is(c, King) {
			return sq
		}
	}
	return -1
}

func (g *Game) FEN() string {
	var b strings.Builder
	for r := 7; r >= 0; r-- {
		empty := 0
		for f := 0; f < 8; f++ {
			p := g.board[Square(f, r)]
			if p.Kind == Empty {
				empty++
				continue
			}
			if empty > 0 {
				fmt.Fprintf(&b, "%d", empty)
				empty = 0
			}
			b.WriteString(p.Letter())
		}
		if empty > 0 {
			fmt.Fprintf(&b, "%d", empty)
		}
		if r > 0 {
			b.WriteByte('/')
		}
	}
	if g.toMove == White {
		b.WriteString(" w ")
	} else {
		b.WriteString(" b ")
	}
	cr := ""
	if g.castle&WhiteKingside != 0 {
		cr += "K"
	}
	if g.castle&WhiteQueenside != 0 {
		cr += "Q"
	}
	if g.castle&BlackKingside != 0 {
		cr += "k"
	}
	if g.castle&BlackQueenside != 0 {
		cr += "q"
	}
	if cr == "" {
		cr = "-"
	}
	b.WriteString(cr)
	b.WriteByte(' ')
	if g.ep < 0 {
		b.WriteByte('-')
	} else {
		b.WriteString(SquareName(g.ep))
	}
	fmt.Fprintf(&b, " %d %d", g.half, g.full)
	return b.String()
}
