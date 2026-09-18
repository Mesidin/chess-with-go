package ui

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"chess-with-go/bot"
	"chess-with-go/engine"
	"chess-with-go/internal/meta"
	"chess-with-go/internal/theme"
	"chess-with-go/learn"
)

const (
	screenW = 720
	screenH = 820
)

type scene int

const (
	sceneMenu scene = iota
	sceneSetup
	sceneSkins
	scenePlay
	sceneLearn
)

type App struct {
	scene    scene
	eng      *engine.Game
	sel      int // -1 none
	targets  []engine.Move
	flash    int
	startBtn button
	skinsBtn button
	quitBtn  button
	playBtns []button
	skinBack  button
	skin      Skin
	wantQuit  bool
	learnBtn  button
	setupPlay button
	setupBack button
	setupBtns []button
	learnNav  []button
	course    *learn.Run
	bot       *bot.Bot
	useBot    bool
	userColor engine.Color
	botWait   int
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
	a := &App{
		sel:       -1,
		skin:      skinByID("classic"),
		bot:       bot.New(rand.New(rand.NewSource(time.Now().UnixNano()))),
		userColor: engine.White,
	}
	cx := float32(screenW)/2 - 140
	a.startBtn = button{cx, 260, 280, 52, "Start game", "start"}
	a.learnBtn = button{cx, 324, 280, 52, "Learn to play", "learn"}
	a.skinsBtn = button{cx, 388, 280, 52, "Board skins", "skins"}
	a.quitBtn = button{cx, 452, 280, 52, "Quit", "quit"}
	a.skinBack = button{cx, float32(screenH - 90), 280, 44, "Back", "back"}
	a.setupPlay = button{cx, 520, 280, 48, "Play", "play"}
	a.setupBack = button{cx, 580, 280, 44, "Back", "back"}
	a.setupBtns = []button{
		{cx, 280, 280, 44, "Humans", "hh"},
		{cx, 336, 280, 44, "White vs bot", "hw"},
		{cx, 392, 280, 44, "Black vs bot", "hb"},
	}
	a.playBtns = []button{
		{24, float32(screenH - 70), 110, 40, "Undo (U)", "undo"},
		{144, float32(screenH - 70), 120, 40, "New game", "new"},
		{274, float32(screenH - 70), 120, 40, "Skins", "skins"},
		{580, float32(screenH - 70), 110, 40, "Menu", "menu"},
	}
	a.learnNav = []button{
		{24, float32(screenH - 70), 110, 40, "Back", "back"},
		{148, float32(screenH - 70), 140, 40, "Next (N)", "next"},
		{580, float32(screenH - 70), 110, 40, "Menu", "menu"},
	}
	return a
}

func (a *App) Layout(int, int) (int, int) { return screenW, screenH }

func (a *App) Update() error {
	applyChrome(theme.Tick())
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
			a.scene = sceneSetup
		}
		if clicked && a.learnBtn.contains(mx, my) {
			a.startLearn()
		}
		if clicked && a.skinsBtn.contains(mx, my) {
			a.scene = sceneSkins
		}
		if clicked && a.quitBtn.contains(mx, my) {
			a.wantQuit = true
		}
	case sceneSetup:
		a.updateSetup(clicked, mx, my)
	case sceneLearn:
		a.updateLearn(clicked, mx, my)
	case sceneSkins:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || (clicked && a.skinBack.contains(mx, my)) {
			if a.eng != nil {
				a.scene = scenePlay
			} else {
				a.scene = sceneMenu
			}
			return nil
		}
		if clicked {
			for i, s := range skins() {
				if skinHit(i).contains(mx, my) {
					a.skin = s
					break
				}
			}
		}
	case scenePlay:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			a.scene = sceneMenu
			return nil
		}
		if a.botToPlay() {
			if a.botWait > 0 {
				a.botWait--
				return nil
			}
			a.botMove()
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyU) {
			a.undoPlay()
		}
		if clicked {
			for _, b := range a.playBtns {
				if !b.contains(mx, my) {
					continue
				}
				switch b.id {
				case "undo":
					a.undoPlay()
				case "new":
					a.startPlay()
				case "menu":
					a.scene = sceneMenu
				case "skins":
					a.scene = sceneSkins
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
				} else if a.botToPlay() {
					a.botWait = 18
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
	applyChrome(theme.Current())
	screen.Fill(colHUD)
	mx, my := ebiten.CursorPosition()
	switch a.scene {
	case sceneMenu:
		drawLabel(screen, meta.Title, screenW/2, 160, 32, colInk, true)
		drawLabel(screen, "A chess game written in Go", screenW/2, 210, 16, colMuted, true)
		drawBtn(screen, a.startBtn, mx, my)
		drawBtn(screen, a.learnBtn, mx, my)
		drawBtn(screen, a.skinsBtn, mx, my)
		drawBtn(screen, a.quitBtn, mx, my)
	case sceneSetup:
		a.drawSetup(screen, mx, my)
	case sceneLearn:
		a.drawLearn(screen, mx, my)
	case sceneSkins:
		a.drawSkins(screen, mx, my)
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
		drawLabel(screen, status, 40, float32(screenH-110), 18, colInk, false)
		for _, b := range a.playBtns {
			drawBtn(screen, b, mx, my)
		}
	}
	drawLabel(screen, meta.Copyright, screenW/2, float32(screenH-22), 13, colMuted, true)
}

func hitSquare(mx, my int) (int, bool) {
	return hitSquareAt(mx, my, 60, 60, 75)
}

func drawBoard(screen *ebiten.Image, a *App) {
	drawBoardAt(screen, a, 60, 60, 75, nil)
}

func skinHit(i int) button {
	const (
		x, w, h, gap float32 = 80, 560, 72, 10
		top                  = 120
	)
	return button{x, top + float32(i)*(h+gap), w, h, "", skins()[i].ID}
}

func (a *App) drawSkins(screen *ebiten.Image, mx, my int) {
	drawLabel(screen, "Board skins", screenW/2, 48, 28, colInk, true)
	drawLabel(screen, "Squares and piece colors. Menus still follow Omarchy.", screenW/2, 84, 14, colMuted, true)
	for i, s := range skins() {
		b := skinHit(i)
		drawBtn(screen, button{b.x, b.y, b.w, b.h, "", s.ID}, mx, my)
		if a.skin.ID == s.ID {
			vectorStrokeRect(screen, b.x-2, b.y-2, b.w+4, b.h+4, colAccent)
		}
		const cell float32 = 22
		ox, oy := b.x+16, b.y+14
		for r := 0; r < 2; r++ {
			for f := 0; f < 2; f++ {
				c := s.Light
				if (f+r)%2 == 1 {
					c = s.Dark
				}
				vector.FillRect(screen, ox+float32(f)*cell, oy+float32(r)*cell, cell, cell, c, true)
			}
		}
		drawPiece(screen, ox+cell*0.5, oy+cell*1.5, cell, engine.Piece{Color: engine.White, Kind: engine.King}, s)
		drawPiece(screen, ox+cell*1.5, oy+cell*0.5, cell, engine.Piece{Color: engine.Black, Kind: engine.King}, s)
		drawLabel(screen, s.Name, b.x+90, b.y+12, 18, colInk, false)
		drawLabel(screen, s.Blurb, b.x+90, b.y+38, 13, colMuted, false)
	}
	drawBtn(screen, a.skinBack, mx, my)
}

func vectorStrokeRect(dst *ebiten.Image, x, y, w, h float32, col color.Color) {
	vector.StrokeLine(dst, x, y, x+w, y, 2, col, true)
	vector.StrokeLine(dst, x+w, y, x+w, y+h, 2, col, true)
	vector.StrokeLine(dst, x+w, y+h, x, y+h, 2, col, true)
	vector.StrokeLine(dst, x, y+h, x, y, 2, col, true)
}
