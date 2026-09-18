package learn

import "chess-with-go/engine"

type Step struct {
	Prompt   string
	From, To int
	AllowAny bool
	Look     bool // no move; student reads, then Next
	Hint     string
	Marks    []int
}

type Lesson struct {
	Title string
	Blurb string
	FEN   string // empty = start position
	Steps []Step
	Done  string
}

func sq(name string) int {
	s, ok := engine.ParseSquare(name)
	if !ok {
		panic(name)
	}
	return s
}

func Lessons() []Lesson {
	return []Lesson{
		{
			Title: "The board",
			Blurb: "White moves first. Pieces sit on squares. Click a piece, then a highlighted square.",
			Steps: []Step{{
				Prompt:   "Move any white pawn one or two squares forward.",
				AllowAny: true,
				Hint:     "Click a white pawn on the 2nd rank, then the empty square in front of it.",
			}},
			Done: "Each side moves one piece per turn (except castling, which moves two).",
		},
		{
			Title: "The pawn",
			Blurb: "Pawns walk forward, never back. From the starting rank they may step two squares.",
			Steps: []Step{{
				Prompt: "Play e2–e4, the most common first move.",
				From:   sq("e2"), To: sq("e4"),
				Marks: []int{sq("e2"), sq("e4")},
				Hint:  "Click the pawn on e2, then e4.",
			}},
			Done: "Pawns capture one square diagonally forward, not straight ahead.",
		},
		{
			Title: "The knight",
			Blurb: "Knights jump in an L: two squares one way and one to the side. They can leap over pieces.",
			Steps: []Step{{
				Prompt: "Develop the king's knight: g1–f3.",
				From:   sq("g1"), To: sq("f3"),
				Marks: []int{sq("g1"), sq("f3")},
				Hint:  "Click the knight on g1, then f3.",
			}},
			Done: "A knight on f3 watches the center and prepares kingside castling.",
		},
		{
			Title: "The bishop",
			Blurb: "Bishops slide any number of squares diagonally, staying on one color.",
			FEN:   "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			Steps: []Step{{
				Prompt: "Put the bishop on the open diagonal: f1–c4.",
				From:   sq("f1"), To: sq("c4"),
				Marks: []int{sq("f1"), sq("c4")},
				Hint:  "Click the bishop on f1, then c4.",
			}},
			Done: "That bishop now aims at Black's weak f7 pawn.",
		},
		{
			Title: "Capture",
			Blurb: "You capture by moving onto an enemy piece. That piece leaves the board.",
			FEN:   "4k3/8/8/3p4/4P3/8/8/4K3 w - - 0 1",
			Steps: []Step{{
				Prompt: "Take the black pawn: e4xd5.",
				From:   sq("e4"), To: sq("d5"),
				Marks: []int{sq("e4"), sq("d5")},
				Hint:  "Pawns capture diagonally. Click e4, then d5.",
			}},
			Done: "The captured pawn is gone. Only one piece occupies a square.",
		},
		{
			Title: "Check",
			Blurb: "A move that attacks the enemy king is check. The opponent must get out of check.",
			FEN:   "4k3/8/8/8/8/8/8/Q3K3 w - - 0 1",
			Steps: []Step{{
				Prompt: "Check the black king with the queen: a1–a8.",
				From:   sq("a1"), To: sq("a8"),
				Marks: []int{sq("a1"), sq("a8"), sq("e8")},
				Hint:  "Click the queen on a1, then a8.",
			}},
			Done: "Black must move the king, block the check, or capture the queen.",
		},
		{
			Title: "Castling",
			Blurb: "If the king and a rook have not moved, and the squares between are empty and safe, they can castle.",
			FEN:   "rnbqk2r/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1",
			Steps: []Step{{
				Prompt: "Castle kingside: click the king on e1, then g1.",
				From:   sq("e1"), To: sq("g1"),
				Marks: []int{sq("e1"), sq("g1"), sq("h1")},
				Hint:  "The king jumps two squares toward the rook; the rook jumps to the other side.",
			}},
			Done: "The king is safer, and the rook is developed. Queenside castling works the same way toward a1.",
		},
		{
			Title: "Checkmate",
			Blurb: "If the king is in check and has no legal move, the game is over.",
			FEN:   "r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4",
			Steps: []Step{{
				Prompt: "Scholar's mate: the queen on f7 is protected by the bishop. Black has no escape. Click Next.",
				Look:   true,
				Marks:  []int{sq("f7"), sq("e8"), sq("c4")},
			}},
			Done: "The queen on f7 is protected by the bishop on c4. The king cannot capture it. You know enough to play a game.",
		},
	}
}
