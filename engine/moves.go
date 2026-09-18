package engine

func (g *Game) LegalMoves() []Move {
	var out []Move
	side := g.toMove
	for _, m := range g.pseudo() {
		g.apply(m)
		ok := !g.InCheck(side)
		g.unapply(m)
		if ok {
			out = append(out, m)
		}
	}
	return out
}

func (g *Game) LegalFrom(from int) []Move {
	var out []Move
	for _, m := range g.LegalMoves() {
		if m.From == from {
			out = append(out, m)
		}
	}
	return out
}

func (g *Game) snapshot(m Move) Move {
	m.PrevCastle = g.castle
	m.PrevEP = g.ep
	m.PrevHalf = g.half
	return m
}

func (g *Game) pseudo() []Move {
	var out []Move
	for sq := 0; sq < 64; sq++ {
		p := g.board[sq]
		if p.Kind == Empty || p.Color != g.toMove {
			continue
		}
		switch p.Kind {
		case Pawn:
			out = append(out, g.pawnMoves(sq, p.Color)...)
		case Knight:
			out = append(out, g.leapMoves(sq, p.Color, knightD)...)
		case King:
			out = append(out, g.leapMoves(sq, p.Color, kingD)...)
			out = append(out, g.castles(sq, p.Color)...)
		case Bishop:
			out = append(out, g.slideMoves(sq, p.Color, bishopD)...)
		case Rook:
			out = append(out, g.slideMoves(sq, p.Color, rookD)...)
		case Queen:
			out = append(out, g.slideMoves(sq, p.Color, bishopD)...)
			out = append(out, g.slideMoves(sq, p.Color, rookD)...)
		}
	}
	return out
}

var knightD = []int{17, 15, 10, 6, -6, -10, -15, -17}
var kingD = []int{9, 8, 7, 1, -1, -7, -8, -9}
var bishopD = []int{9, 7, -7, -9}
var rookD = []int{8, 1, -1, -8}

func onBoard(from, to int) bool {
	if to < 0 || to > 63 {
		return false
	}
	df := File(to) - File(from)
	if df > 2 || df < -2 {
		return false
	}
	return true
}

func (g *Game) leapMoves(from int, c Color, deltas []int) []Move {
	var out []Move
	for _, d := range deltas {
		to := from + d
		if !onBoard(from, to) {
			continue
		}
		q := g.board[to]
		if q.Kind != Empty && q.Color == c {
			continue
		}
		m := g.snapshot(Move{From: from, To: to, Capture: q})
		out = append(out, m)
	}
	return out
}

func (g *Game) slideMoves(from int, c Color, deltas []int) []Move {
	var out []Move
	for _, d := range deltas {
		to := from
		for {
			fromSq := to
			to = to + d
			if !onBoard(fromSq, to) {
				break
			}
			q := g.board[to]
			if q.Kind != Empty {
				if q.Color != c {
					out = append(out, g.snapshot(Move{From: from, To: to, Capture: q}))
				}
				break
			}
			out = append(out, g.snapshot(Move{From: from, To: to}))
		}
	}
	return out
}

func (g *Game) pawnMoves(from int, c Color) []Move {
	var out []Move
	dir, start, promo := 8, 1, 7
	if c == Black {
		dir, start, promo = -8, 6, 0
	}
	one := from + dir
	if one >= 0 && one <= 63 && g.board[one].Kind == Empty {
		out = append(out, g.promos(from, one, Piece{}, c, Rank(one) == promo)...)
		two := from + 2*dir
		if Rank(from) == start && two >= 0 && two <= 63 && g.board[two].Kind == Empty {
			out = append(out, g.snapshot(Move{From: from, To: two, DoublePawn: true}))
		}
	}
	for _, df := range []int{-1, 1} {
		if File(from)+df < 0 || File(from)+df > 7 {
			continue
		}
		to := from + dir + df
		if to < 0 || to > 63 {
			continue
		}
		q := g.board[to]
		if q.Kind != Empty && q.Color != c {
			out = append(out, g.promos(from, to, q, c, Rank(to) == promo)...)
		}
		if g.ep == to {
			out = append(out, g.snapshot(Move{From: from, To: to, EnPassant: true, Capture: Piece{Color: c.Opponent(), Kind: Pawn}}))
		}
	}
	return out
}

func (g *Game) promos(from, to int, cap Piece, c Color, promo bool) []Move {
	if !promo {
		return []Move{g.snapshot(Move{From: from, To: to, Capture: cap})}
	}
	var out []Move
	for _, k := range []Kind{Queen, Rook, Bishop, Knight} {
		out = append(out, g.snapshot(Move{From: from, To: to, Capture: cap, Promo: k}))
	}
	return out
}

func (g *Game) castles(from int, c Color) []Move {
	var out []Move
	rank := 0
	ks, qs := WhiteKingside, WhiteQueenside
	if c == Black {
		rank = 7
		ks, qs = BlackKingside, BlackQueenside
	}
	if from != Square(4, rank) {
		return nil
	}
	if g.InCheck(c) {
		return nil
	}
	empty := func(files ...int) bool {
		for _, f := range files {
			if g.board[Square(f, rank)].Kind != Empty {
				return false
			}
		}
		return true
	}
	safe := func(files ...int) bool {
		for _, f := range files {
			if g.attacked(Square(f, rank), c.Opponent()) {
				return false
			}
		}
		return true
	}
	if g.castle&ks != 0 && empty(5, 6) && safe(5, 6) {
		out = append(out, g.snapshot(Move{From: from, To: Square(6, rank), Castle: true}))
	}
	if g.castle&qs != 0 && empty(1, 2, 3) && safe(2, 3) {
		out = append(out, g.snapshot(Move{From: from, To: Square(2, rank), Castle: true}))
	}
	return out
}

func (g *Game) attacked(sq int, by Color) bool {
	// Pawns
	dir := -8
	if by == White {
		dir = 8
	}
	for _, df := range []int{-1, 1} {
		from := sq - dir - df
		if File(sq)-df < 0 || File(sq)-df > 7 {
			continue
		}
		if from >= 0 && from <= 63 && g.board[from].Is(by, Pawn) {
			return true
		}
	}
	for _, d := range knightD {
		from := sq - d
		if onBoard(sq, from) && g.board[from].Is(by, Knight) {
			return true
		}
	}
	for _, d := range kingD {
		from := sq - d
		if onBoard(sq, from) && g.board[from].Is(by, King) {
			return true
		}
	}
	if g.slideAttack(sq, by, bishopD, Bishop, Queen) {
		return true
	}
	if g.slideAttack(sq, by, rookD, Rook, Queen) {
		return true
	}
	return false
}

func (g *Game) slideAttack(sq int, by Color, deltas []int, kinds ...Kind) bool {
	ok := func(k Kind) bool {
		for _, t := range kinds {
			if k == t {
				return true
			}
		}
		return false
	}
	for _, d := range deltas {
		from := sq
		for {
			prev := from
			from -= d
			if !onBoard(prev, from) {
				break
			}
			p := g.board[from]
			if p.Kind == Empty {
				continue
			}
			if p.Color == by && ok(p.Kind) {
				return true
			}
			break
		}
	}
	return false
}
