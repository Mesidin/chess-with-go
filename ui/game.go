package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"chess-with-go/engine"
	"chess-with-go/internal/meta"
)

const (
	screenW = 720
	screenH = 820
)

type scene int

const (
	sceneMenu scene = iota
	scenePlay
)

type App struct {
	scene    scene
	eng      *engine.Game
	sel      int // -1 none
	targets  []engine.Move
	flash    int
	startBtn button
	quitBtn  button
	playBtns []button
	wantQuit bool
}

type button struct {
	x, y, w, h float32
	label, id  string
}

func (b button) contains(x, y int) bool {
	fx, fy := float32(x), float32(y)
	return fx >= b.x && fy >= b.y && fx < b.x+b.w && fy < b.y+b.h
}

func NewApp() *App {
	a := &App{sel: -1}
	cx := float32(screenW)/2 - 140
	a.startBtn = button{cx, 320, 280, 52, "Start game", "start"}
	a.quitBtn = button{cx, 392, 280, 52, "Quit", "quit"}
	a.playBtns = []button{
		{40, float32(screenH - 70), 120, 40, "Undo (U)", "undo"},
		{180, float32(screenH - 70), 140, 40, "New game", "new"},
		{560, float32(screenH - 70), 120, 40, "Menu", "menu"},
	}
	return a
}

func (a *App) Layout(int, int) (int, int) { return screenW, screenH }

func (a *App) Update() error {
	if a.wantQuit {
		return ebiten.Termination
	}
	if a.flash > 0 {
		a.flash--
	}
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	mx, my := ebiten.CursorPosition()
	switch a.scene {
	case sceneMenu:
		if clicked && a.startBtn.contains(mx, my) {
			a.eng = engine.New()
			a.sel = -1
			a.scene = scenePlay
		}
		if clicked && a.quitBtn.contains(mx, my) {
			a.wantQuit = true
		}
	case scenePlay:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			a.scene = sceneMenu
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyU) {
			_ = a.eng.Undo()
			a.sel = -1
		}
		if clicked {
			for _, b := range a.playBtns {
				if !b.contains(mx, my) {
					continue
				}
				switch b.id {
				case "undo":
					_ = a.eng.Undo()
					a.sel = -1
				case "new":
					a.eng = engine.New()
					a.sel = -1
				case "menu":
					a.scene = sceneMenu
				}
				return nil
			}
			if sq, ok := hitSquare(mx, my); ok && !a.eng.Over() {
				a.clickSq(sq)
			}
		}
	}
	return nil
}

func (a *App) clickSq(sq int) {
	if a.sel >= 0 {
		for _, m := range a.targets {
			if m.To == sq {
				promo := m.Promo
				if promo == engine.Empty {
					promo = engine.Empty
				}
				if err := a.eng.Play(a.sel, sq, promo); err != nil {
					a.flash = 12
				}
				a.sel = -1
				a.targets = nil
				return
			}
		}
	}
	p := a.eng.At(sq)
	if p.Kind != engine.Empty && p.Color == a.eng.ToMove() {
		a.sel = sq
		a.targets = a.eng.LegalFrom(sq)
		return
	}
	a.sel = -1
	a.targets = nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x2B, 0x1D, 0x12, 0xFF})
	mx, my := ebiten.CursorPosition()
	switch a.scene {
	case sceneMenu:
		drawLabel(screen, meta.Title, screenW/2, 160, 32, color.RGBA{0xF6, 0xEB, 0xD8, 0xFF}, true)
		drawLabel(screen, "A chess game written in Go", screenW/2, 210, 16, color.RGBA{0xC4, 0xB0, 0x90, 0xFF}, true)
		drawBtn(screen, a.startBtn, mx, my)
		drawBtn(screen, a.quitBtn, mx, my)
	case scenePlay:
		drawBoard(screen, a)
		status := a.eng.ToMove().String() + " to move"
		if a.eng.Over() {
			status = a.eng.Result().String()
			if a.eng.Reason() != "" {
				status += " · " + a.eng.Reason()
			}
		} else if a.eng.InCheck(a.eng.ToMove()) {
			status += " · check"
		}
		drawLabel(screen, status, 40, float32(screenH-110), 18, color.RGBA{0xF6, 0xEB, 0xD8, 0xFF}, false)
		for _, b := range a.playBtns {
			drawBtn(screen, b, mx, my)
		}
	}
	drawLabel(screen, meta.Copyright, screenW/2, float32(screenH-22), 13, color.RGBA{0xC4, 0xB0, 0x90, 0xFF}, true)
}

func hitSquare(mx, my int) (int, bool) {
	const origin = 60
	const cell = 75
	fx := mx - origin
	fy := my - origin
	if fx < 0 || fy < 0 || fx >= 8*cell || fy >= 8*cell {
		return 0, false
	}
	file := fx / cell
	rankFromTop := fy / cell
	rank := 7 - rankFromTop
	return engine.Square(file, rank), true
}

func drawBoard(screen *ebiten.Image, a *App) {
	const origin, cell float32 = 60, 75
	light := color.RGBA{0xF0, 0xD9, 0xB5, 0xFF}
	dark := color.RGBA{0xB5, 0x88, 0x63, 0xFF}
	hi := color.RGBA{0xF6, 0xF1, 0x6D, 0xA0}
	tgt := color.RGBA{0x2E, 0x8B, 0x57, 0xAA}
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			x := origin + float32(f)*cell
			y := origin + float32(7-r)*cell
			col := light
			if (f+r)%2 == 1 {
				col = dark
			}
			vector.FillRect(screen, x, y, cell, cell, col, true)
			sq := engine.Square(f, r)
			if a.sel == sq {
				vector.FillRect(screen, x, y, cell, cell, hi, true)
			}
			for _, m := range a.targets {
				if m.To == sq {
					vector.FillCircle(screen, x+cell/2, y+cell/2, 10, tgt, true)
				}
			}
			p := a.eng.At(sq)
			if p.Kind != engine.Empty {
				drawPiece(screen, x+cell/2, y+cell/2, p)
			}
		}
	}
	files := "abcdefgh"
	for i := 0; i < 8; i++ {
		drawLabel(screen, string(files[i]), origin+float32(i)*cell+cell/2, origin+8*cell+6, 13, color.RGBA{0xC4, 0xB0, 0x90, 0xFF}, true)
		drawLabel(screen, fmt.Sprintf("%d", i+1), origin-18, origin+float32(7-i)*cell+cell/2-8, 13, color.RGBA{0xC4, 0xB0, 0x90, 0xFF}, true)
	}
}
