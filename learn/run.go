package learn

import "chess-with-go/engine"

type Run struct {
	lessons   []Lesson
	index     int
	step      int
	game      *engine.Game
	sel       int
	hint      string
	complete  bool
	graduated bool
}

func NewRun() (*Run, error) {
	r := &Run{lessons: Lessons(), sel: -1}
	if err := r.load(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Run) load() error {
	if r.index >= len(r.lessons) {
		r.graduated = true
		r.complete = true
		return nil
	}
	les := r.lessons[r.index]
	var g *engine.Game
	var err error
	if les.FEN != "" {
		g, err = engine.ParseFEN(les.FEN)
	} else {
		g = engine.New()
	}
	if err != nil {
		return err
	}
	r.game = g
	r.step = 0
	r.sel = -1
	r.hint = ""
	r.complete = false
	return nil
}

func (r *Run) Game() *engine.Game { return r.game }
func (r *Run) Hint() string       { return r.hint }
func (r *Run) Complete() bool     { return r.complete }
func (r *Run) Graduated() bool    { return r.graduated }
func (r *Run) Index() int         { return r.index }
func (r *Run) Count() int         { return len(r.lessons) }
func (r *Run) Sel() int           { return r.sel }

func (r *Run) Lesson() Lesson {
	if r.index >= len(r.lessons) {
		return Lesson{Title: "Ready to play", Done: "Start a game from the menu. Try White vs bot."}
	}
	return r.lessons[r.index]
}

func (r *Run) Step() Step {
	les := r.Lesson()
	if r.complete || r.step >= len(les.Steps) {
		return Step{}
	}
	return les.Steps[r.step]
}

func (r *Run) Prompt() string {
	if r.graduated || r.complete {
		return r.Lesson().Done
	}
	if r.hint != "" {
		return r.hint
	}
	return r.Step().Prompt
}

func (r *Run) Marks() []int {
	if r.complete || r.graduated {
		return nil
	}
	st := r.Step()
	out := append([]int(nil), st.Marks...)
	if st.From >= 0 && !st.AllowAny && !st.Look {
		out = append(out, st.From, st.To)
	}
	return out
}

func (r *Run) advance() {
	r.hint = ""
	r.sel = -1
	r.step++
	if r.step >= len(r.Lesson().Steps) {
		r.complete = true
	}
}

func (r *Run) Click(sq int) (flash bool) {
	if r.game == nil || r.complete || r.graduated {
		return false
	}
	st := r.Step()
	if st.Look {
		return false
	}
	if r.sel < 0 {
		p := r.game.At(sq)
		if p.Kind == engine.Empty || p.Color != r.game.ToMove() {
			r.hint = st.Hint
			return true
		}
		if !st.AllowAny && sq != st.From {
			r.hint = st.Hint
			return true
		}
		r.sel = sq
		r.hint = ""
		return false
	}
	from := r.sel
	r.sel = -1
	if st.AllowAny {
		if err := r.game.Play(from, sq, engine.Empty); err != nil {
			r.hint = st.Hint
			return true
		}
		r.advance()
		return false
	}
	if from != st.From || sq != st.To {
		r.hint = st.Hint
		return true
	}
	if err := r.game.Play(from, sq, engine.Empty); err != nil {
		r.hint = err.Error()
		return true
	}
	r.advance()
	return false
}

func (r *Run) Next() error {
	if r.index < len(r.lessons)-1 {
		r.index++
		return r.load()
	}
	r.index = len(r.lessons)
	r.graduated = true
	r.complete = true
	return nil
}

func (r *Run) Prev() error {
	if r.index > 0 {
		r.index--
	}
	r.graduated = false
	return r.load()
}
