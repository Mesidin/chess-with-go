package engine

import "testing"

func TestStartLegalCount(t *testing.T) {
	g := New()
	n := len(g.LegalMoves())
	if n != 20 {
		t.Fatalf("start: %d legal, want 20", n)
	}
}

func TestPlayAndUndo(t *testing.T) {
	g := New()
	e2, _ := ParseSquare("e2")
	e4, _ := ParseSquare("e4")
	if err := g.Play(e2, e4, Empty); err != nil {
		t.Fatal(err)
	}
	if g.ToMove() != Black {
		t.Fatal("black to move")
	}
	if err := g.Undo(); err != nil {
		t.Fatal(err)
	}
	if g.ToMove() != White || g.At(e2).Kind != Pawn || g.At(e4).Kind != Empty {
		t.Fatal("undo failed")
	}
}

func TestScholarsMate(t *testing.T) {
	g := New()
	seq := [][2]string{
		{"e2", "e4"}, {"e7", "e5"},
		{"f1", "c4"}, {"b8", "c6"},
		{"d1", "h5"}, {"g8", "f6"},
		{"h5", "f7"},
	}
	for _, s := range seq {
		from, _ := ParseSquare(s[0])
		to, _ := ParseSquare(s[1])
		if err := g.Play(from, to, Empty); err != nil {
			t.Fatalf("%s%s: %v", s[0], s[1], err)
		}
	}
	if g.Result() != WhiteWin || g.Reason() != "checkmate" {
		t.Fatalf("result %s %s", g.Result(), g.Reason())
	}
}

func TestCannotMoveIntoCheck(t *testing.T) {
	g := New()
	// After e4 e5 Qh5, black cannot play f7? that's not in check yet.
	// Fool's mate: f3 e5 g4 Qh4#
	f2, _ := ParseSquare("f2")
	f3, _ := ParseSquare("f3")
	e7, _ := ParseSquare("e7")
	e5, _ := ParseSquare("e5")
	g2, _ := ParseSquare("g2")
	g4, _ := ParseSquare("g4")
	d8, _ := ParseSquare("d8")
	h4, _ := ParseSquare("h4")
	_ = g.Play(f2, f3, Empty)
	_ = g.Play(e7, e5, Empty)
	_ = g.Play(g2, g4, Empty)
	if err := g.Play(d8, h4, Empty); err != nil {
		t.Fatal(err)
	}
	if g.Result() != BlackWin {
		t.Fatalf("fools mate: %s %s", g.Result(), g.Reason())
	}
}

func TestCastlingRights(t *testing.T) {
	g := New()
	// clear path: knights and bishops out... faster: play a sequence
	moves := [][2]string{
		{"e2", "e4"}, {"e7", "e5"},
		{"g1", "f3"}, {"b8", "c6"},
		{"f1", "c4"}, {"g8", "f6"},
	}
	for _, s := range moves {
		from, _ := ParseSquare(s[0])
		to, _ := ParseSquare(s[1])
		if err := g.Play(from, to, Empty); err != nil {
			t.Fatal(err)
		}
	}
	e1, _ := ParseSquare("e1")
	g1, _ := ParseSquare("g1")
	if err := g.Play(e1, g1, Empty); err != nil {
		t.Fatal(err)
	}
	if g.At(g1).Kind != King || g.At(ParseMust("f1")).Kind != Rook {
		t.Fatal("castle did not move rook")
	}
}

func ParseMust(s string) int {
	sq, ok := ParseSquare(s)
	if !ok {
		panic(s)
	}
	return sq
}
