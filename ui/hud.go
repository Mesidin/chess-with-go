package ui

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"

	"chess-with-go/engine"
)

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
	col := color.RGBA{0x5C, 0x3D, 0x22, 0xFF}
	if b.contains(mx, my) {
		col = color.RGBA{0x7A, 0x52, 0x2D, 0xFF}
	}
	vector.FillRect(dst, b.x, b.y, b.w, b.h, col, true)
	drawLabel(dst, b.label, b.x+b.w/2, b.y+b.h/2-10, 18, color.RGBA{0xF6, 0xEB, 0xD8, 0xFF}, true)
}

func drawPiece(dst *ebiten.Image, cx, cy float32, p engine.Piece) {
	fill := color.RGBA{0xF4, 0xF0, 0xE6, 0xFF}
	ink := color.RGBA{0x14, 0x12, 0x10, 0xFF}
	if p.Color == engine.Black {
		fill = color.RGBA{0x22, 0x20, 0x1C, 0xFF}
		ink = color.RGBA{0xF4, 0xF0, 0xE6, 0xFF}
	}
	vector.FillCircle(dst, cx, cy, 26, fill, true)
	vector.StrokeCircle(dst, cx, cy, 26, 1.5, color.RGBA{0, 0, 0, 80}, true)
	glyph := "P"
	switch p.Kind {
	case engine.Knight:
		glyph = "N"
	case engine.Bishop:
		glyph = "B"
	case engine.Rook:
		glyph = "R"
	case engine.Queen:
		glyph = "Q"
	case engine.King:
		glyph = "K"
	}
	drawLabel(dst, glyph, cx, cy-11, 22, ink, true)
}
