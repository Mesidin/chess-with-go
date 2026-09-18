package ui

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"

	"chess-with-go/internal/theme"
)

var (
	colHUD     color.RGBA
	colBtn     color.RGBA
	colBtnHot  color.RGBA
	colBtnText color.RGBA
	colInk     color.RGBA
	colMuted   color.RGBA
	colAccent color.RGBA
)

func init() {
	applyChrome(theme.Current())
}

func applyChrome(p theme.Palette) {
	colHUD = p.Background
	colBtn = p.LighterBackground
	colBtnHot = p.Accent
	colBtnText = p.Foreground
	colInk = p.Foreground
	colMuted = p.Muted
	colAccent = p.Accent
}

var faceSrc *text.GoTextFaceSource

func init() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}
	faceSrc = src
}

func face(size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: faceSrc, Size: size}
}

func drawLabel(dst *ebiten.Image, s string, x, y, size float32, col color.Color, center bool) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(col)
	if center {
		op.PrimaryAlign = text.AlignCenter
	}
	text.Draw(dst, s, face(float64(size)), op)
}

func drawBtn(dst *ebiten.Image, b button, mx, my int) {
	col := colBtn
	if b.contains(mx, my) {
		col = colBtnHot
	}
	vector.FillRect(dst, b.x, b.y, b.w, b.h, col, true)
	drawLabel(dst, b.label, b.x+b.w/2, b.y+b.h/2-10, 18, colBtnText, true)
}


