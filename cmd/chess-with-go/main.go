package main

import (
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"chess-with-go/internal/meta"
	"chess-with-go/ui"
)

func main() {
	ebiten.SetWindowSize(720, 820)
	ebiten.SetWindowTitle(meta.Title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(ui.NewApp()); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
}
