package bot

import (
	"math/rand"

	"chess-with-go/engine"
)

// Bot is a 1-ply capture-aware opponent. It is a beginner sparring partner.
type Bot struct {
	rng *rand.Rand
}

func New(rng *rand.Rand) *Bot {
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	return &Bot{rng: rng}
}

func (b *Bot) Pick(g *engine.Game) (engine.Move, bool) {
	moves := g.LegalMoves()
	if len(moves) == 0 {
		return engine.Move{}, false
	}
	best := -1e18
	var picks []engine.Move
	me := g.ToMove()
	for _, m := range moves {
		s := b.score(g, m, me)
		if s > best {
			best = s
			picks = []engine.Move{m}
		} else if s == best {
			picks = append(picks, m)
		}
	}
	return picks[b.rng.Intn(len(picks))], true
}

func (b *Bot) score(g *engine.Game, m engine.Move, me engine.Color) float64 {
	trial := g.Clone()
	if err := trial.PlayMove(m); err != nil {
		return -1e9
	}
	s := material(trial, me) + b.rng.Float64()*8
	if m.Capture.Kind != engine.Empty {
		s += 40 + float64(val(m.Capture.Kind))
	}
	if m.Promo != engine.Empty {
		s += float64(val(m.Promo))
	}
	if trial.Result() == winFor(me) {
		return 1e7
	}
	if trial.InCheck(me.Opponent()) {
		s += 35
	}
	// Don't hang the piece we just moved if a cheaper recapture exists.
	if hanging(trial, m.To, me) {
		s -= float64(val(trial.At(m.To).Kind)) * 0.85
	}
	toFile, toRank := float64(engine.File(m.To)), float64(engine.Rank(m.To))
	center := abs(toFile-3.5) + abs(toRank-3.5)
	s += 12 - center*2
	return s
}

func winFor(c engine.Color) engine.Result {
	if c == engine.White {
		return engine.WhiteWin
	}
	return engine.BlackWin
}

func material(g *engine.Game, me engine.Color) float64 {
	var n float64
	for sq := 0; sq < 64; sq++ {
		p := g.At(sq)
		if p.Kind == engine.Empty {
			continue
		}
		v := float64(val(p.Kind))
		if p.Color == me {
			n += v
		} else {
			n -= v
		}
	}
	return n
}

func hanging(g *engine.Game, sq int, me engine.Color) bool {
	p := g.At(sq)
	if p.Kind == engine.Empty || p.Color != me {
		return false
	}
	for _, m := range g.LegalMoves() {
		if m.To == sq && val(m.Capture.Kind) >= val(p.Kind)-50 {
			return true
		}
	}
	return false
}

func val(k engine.Kind) int {
	switch k {
	case engine.Pawn:
		return 100
	case engine.Knight, engine.Bishop:
		return 320
	case engine.Rook:
		return 500
	case engine.Queen:
		return 900
	case engine.King:
		return 20000
	default:
		return 0
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
