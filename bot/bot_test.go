package bot

import (
	"math/rand"
	"testing"

	"chess-with-go/engine"
)

func TestPickIsLegal(t *testing.T) {
	g := engine.New()
	b := New(rand.New(rand.NewSource(2)))
	m, ok := b.Pick(g)
	if !ok {
		t.Fatal("no move")
	}
	if err := g.PlayMove(m); err != nil {
		t.Fatal(err)
	}
}

func TestPrefersCapture(t *testing.T) {
	g, err := engine.ParseFEN("4k3/8/8/3p4/4P3/8/8/4K3 w - - 0 1")
	if err != nil {
		t.Fatal(err)
	}
	b := New(rand.New(rand.NewSource(1)))
	m, ok := b.Pick(g)
	if !ok {
		t.Fatal("no move")
	}
	if engine.SquareName(m.From)+engine.SquareName(m.To) != "e4d5" {
		t.Fatalf("expected e4xd5, got %s", m)
	}
}
