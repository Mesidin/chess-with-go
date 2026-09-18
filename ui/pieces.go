package ui

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"chess-with-go/engine"
)

// Wikipedia / Cburnett Staunton pieces (Colin M.L. Burnett, CC BY-SA 3.0).
//
//go:embed pieces/*.png
var pieceFS embed.FS

var pieceImg [2][7]*ebiten.Image

func init() {
	load := func(name string) *ebiten.Image {
		b, err := pieceFS.ReadFile("pieces/" + name)
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			panic(err)
		}
		return ebiten.NewImageFromImage(img)
	}
	pieceImg[engine.White][engine.King] = load("wK.png")
	pieceImg[engine.White][engine.Queen] = load("wQ.png")
	pieceImg[engine.White][engine.Rook] = load("wR.png")
	pieceImg[engine.White][engine.Bishop] = load("wB.png")
	pieceImg[engine.White][engine.Knight] = load("wN.png")
	pieceImg[engine.White][engine.Pawn] = load("wP.png")
	pieceImg[engine.Black][engine.King] = load("bK.png")
	pieceImg[engine.Black][engine.Queen] = load("bQ.png")
	pieceImg[engine.Black][engine.Rook] = load("bR.png")
	pieceImg[engine.Black][engine.Bishop] = load("bB.png")
	pieceImg[engine.Black][engine.Knight] = load("bN.png")
	pieceImg[engine.Black][engine.Pawn] = load("bP.png")
}

func drawPiece(dst *ebiten.Image, cx, cy, cell float32, p engine.Piece, skin Skin) {
	if p.Kind == engine.Empty {
		return
	}
	img := pieceImg[p.Color][p.Kind]
	if img == nil {
		return
	}
	size := cell * 0.90
	b := img.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(float64(size)/float64(b.Dx()), float64(size)/float64(b.Dy()))
	op.GeoM.Translate(float64(cx-size/2), float64(cy-size/2))
	if p.Color == engine.White && !skin.whiteIdentity() {
		t := skin.WhiteTint
		op.ColorScale.Scale(float32(t.R)/255, float32(t.G)/255, float32(t.B)/255, 1)
	}
	if p.Color == engine.Black {
		t := skin.BlackTint
		lift := skin.BlackLift
		if lift < 1 {
			lift = 1
		}
		op.ColorScale.Scale(float32(t.R)/255*lift, float32(t.G)/255*lift, float32(t.B)/255*lift, 1)
	}
	dst.DrawImage(img, op)
}
