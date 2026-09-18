package learn

import (
	"testing"

	"chess-with-go/engine"
)

func TestLessonsLoad(t *testing.T) {
	for i, les := range Lessons() {
		var g *engine.Game
		var err error
		if les.FEN != "" {
			g, err = engine.ParseFEN(les.FEN)
		} else {
			g = engine.New()
		}
		if err != nil {
			t.Fatalf("%d %s: %v", i, les.Title, err)
		}
		for _, st := range les.Steps {
			if st.Look || st.AllowAny {
				continue
			}
			if err := g.Clone().Play(st.From, st.To, engine.Empty); err != nil {
				t.Fatalf("%s: %s%s illegal: %v", les.Title, engine.SquareName(st.From), engine.SquareName(st.To), err)
			}
		}
	}
}

func TestRunPawnLesson(t *testing.T) {
	r, err := NewRun()
	if err != nil {
		t.Fatal(err)
	}
	r.index = 1
	if err := r.load(); err != nil {
		t.Fatal(err)
	}
	r.Click(sq("e2"))
	r.Click(sq("e4"))
	if !r.Complete() {
		t.Fatalf("hint %s", r.Hint())
	}
}
