package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"chess-with-go/engine"
	"chess-with-go/learn"
)

const learnTop float32 = 108

func (a *App) startLearn() {
	r, err := learn.NewRun()
	if err != nil {
		return
	}
	a.course = r
	a.eng = r.Game()
	a.sel = -1
	a.targets = nil
	a.scene = sceneLearn
}

func (a *App) syncLearn() {
	if a.course == nil {
		return
	}
	a.eng = a.course.Game()
	a.sel = a.course.Sel()
	if a.sel >= 0 && a.eng != nil {
		a.targets = a.eng.LegalFrom(a.sel)
	} else {
		a.targets = nil
	}
}

func (a *App) updateLearn(clicked bool, mx, my int) {
	if a.course == nil {
		a.scene = sceneMenu
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		a.scene = sceneMenu
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		_ = a.course.Next()
		a.syncLearn()
		if a.course.Graduated() {
			a.useBot = true
			a.userColor = engine.White
			a.scene = sceneSetup
		}
		return
	}
	if !clicked {
		return
	}
	for _, b := range a.learnNav {
		if !b.contains(mx, my) {
			continue
		}
		switch b.id {
		case "back":
			_ = a.course.Prev()
			a.syncLearn()
		case "next":
			_ = a.course.Next()
			a.syncLearn()
			if a.course.Graduated() {
				a.useBot = true
				a.userColor = engine.White
				a.scene = sceneSetup
			}
		case "menu":
			a.scene = sceneMenu
		}
		return
	}
	if sq, ok := hitSquareLearn(mx, my); ok {
		if a.course.Click(sq) {
			a.flash = 12
		}
		a.syncLearn()
	}
}

func hitSquareLearn(mx, my int) (int, bool) {
	return hitSquareAt(mx, my, 48, 118, 66)
}

func (a *App) drawLearn(screen *ebiten.Image, mx, my int) {
	if a.course == nil {
		return
	}
	les := a.course.Lesson()
	n := a.course.Index() + 1
	if a.course.Graduated() {
		n = a.course.Count()
	}
	drawLabel(screen, fmt.Sprintf("Learn  ·  %d / %d  ·  %s", n, a.course.Count(), les.Title), 36, 16, 16, colAccent, false)
	drawLabel(screen, les.Blurb, 36, 40, 13, colMuted, false)
	prompt := a.course.Prompt()
	col := colInk
	if a.course.Hint() != "" && !a.course.Complete() {
		col = colAccent
	}
	drawLabel(screen, prompt, 36, 64, 16, col, false)

	drawBoardAt(screen, a, 48, 118, 66, a.course.Marks())
	for _, b := range a.learnNav {
		drawBtn(screen, b, mx, my)
	}
}
