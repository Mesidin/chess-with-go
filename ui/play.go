package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"chess-with-go/engine"
)

func hitSquareAt(mx, my, ox, oy, cell int) (int, bool) {
	fx := mx - ox
	fy := my - oy
	if fx < 0 || fy < 0 || fx >= 8*cell || fy >= 8*cell {
		return 0, false
	}
	file := fx / cell
	rankFromTop := fy / cell
	return engine.Square(file, 7-rankFromTop), true
}

func drawBoardAt(screen *ebiten.Image, a *App, ox, oy, cell float32, marks []int) {
	if a.eng == nil {
		return
	}
	light, dark, hi, tgt := a.skin.Light, a.skin.Dark, a.skin.Hi, a.skin.Tgt
	marked := map[int]bool{}
	for _, sq := range marks {
		marked[sq] = true
	}
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			x := ox + float32(f)*cell
			y := oy + float32(7-r)*cell
			col := light
			if (f+r)%2 == 1 {
				col = dark
			}
			vector.FillRect(screen, x, y, cell, cell, col, true)
			sq := engine.Square(f, r)
			if a.sel == sq {
				vector.FillRect(screen, x, y, cell, cell, hi, true)
			}
			if marked[sq] {
				vector.StrokeRect(screen, x+3, y+3, cell-6, cell-6, 2, colAccent, true)
			}
			for _, m := range a.targets {
				if m.To == sq {
					vector.FillCircle(screen, x+cell/2, y+cell/2, 8, tgt, true)
				}
			}
			p := a.eng.At(sq)
			if p.Kind != engine.Empty {
				drawPiece(screen, x+cell/2, y+cell/2, cell, p, a.skin)
			}
		}
	}
	files := "abcdefgh"
	for i := 0; i < 8; i++ {
		drawLabel(screen, string(files[i]), ox+float32(i)*cell+cell/2, oy+8*cell+6, 13, colMuted, true)
		drawLabel(screen, fmtInt(i+1), ox-18, oy+float32(7-i)*cell+cell/2-8, 13, colMuted, true)
	}
}

func fmtInt(n int) string {
	return string(rune('0' + n))
}

func (a *App) updateSetup(clicked bool, mx, my int) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		a.scene = sceneMenu
		return
	}
	if !clicked {
		return
	}
	if a.setupBack.contains(mx, my) {
		a.scene = sceneMenu
		return
	}
	if a.setupPlay.contains(mx, my) {
		a.startPlay()
		return
	}
	for _, b := range a.setupBtns {
		if !b.contains(mx, my) {
			continue
		}
		switch b.id {
		case "hh":
			a.useBot = false
		case "hw":
			a.useBot = true
			a.userColor = engine.White
		case "hb":
			a.useBot = true
			a.userColor = engine.Black
		}
	}
}

func (a *App) drawSetup(screen *ebiten.Image, mx, my int) {
	drawLabel(screen, "New game", screenW/2, 80, 28, colInk, true)
	drawLabel(screen, "Who plays", screenW/2, 230, 14, colMuted, true)
	for _, b := range a.setupBtns {
		drawBtn(screen, b, mx, my)
		sel := (b.id == "hh" && !a.useBot) ||
			(b.id == "hw" && a.useBot && a.userColor == engine.White) ||
			(b.id == "hb" && a.useBot && a.userColor == engine.Black)
		if sel {
			vectorStrokeRect(screen, b.x-3, b.y-3, b.w+6, b.h+6, colAccent)
		}
	}
	drawBtn(screen, a.setupPlay, mx, my)
	drawBtn(screen, a.setupBack, mx, my)
}

func (a *App) startPlay() {
	a.eng = engine.New()
	a.sel = -1
	a.targets = nil
	a.botWait = 0
	a.scene = scenePlay
	if a.botToPlay() {
		a.botWait = 24
	}
}

func (a *App) botToPlay() bool {
	if !a.useBot || a.eng == nil || a.eng.Over() {
		return false
	}
	return a.eng.ToMove() != a.userColor
}

func (a *App) botMove() {
	m, ok := a.bot.Pick(a.eng)
	if !ok {
		return
	}
	_ = a.eng.PlayMove(m)
}

func (a *App) undoPlay() {
	if a.eng == nil {
		return
	}
	_ = a.eng.Undo()
	if a.useBot {
		_ = a.eng.Undo()
	}
	a.sel = -1
	a.targets = nil
}
