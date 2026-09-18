package ui

import "image/color"

// Skin restyles the board and piece tints. Chrome (menus, HUD) stays Omarchy.
type Skin struct {
	ID        string
	Name      string
	Blurb     string
	Light     color.RGBA
	Dark      color.RGBA
	Hi        color.RGBA
	Tgt       color.RGBA
	WhiteTint color.RGBA // {255,255,255,255} = un-tinted white army
	BlackTint color.RGBA // multiplied into Cburnett black sprites
	BlackLift float32    // >1 lifts dark pixels so a tint reads as a color
}

func (s Skin) whiteIdentity() bool {
	return s.WhiteTint.R > 240 && s.WhiteTint.G > 240 && s.WhiteTint.B > 240
}

func skins() []Skin {
	hi := func(c color.RGBA) color.RGBA { c.A = 0x99; return c }
	return []Skin{
		{
			ID: "classic", Name: "Classic",
			Blurb:     "Walnut and cream, natural Staunton.",
			Light:     color.RGBA{0xF0, 0xD9, 0xB5, 0xFF},
			Dark:      color.RGBA{0xB5, 0x88, 0x63, 0xFF},
			Hi:        hi(color.RGBA{0xF6, 0xE2, 0x7A, 0xFF}),
			Tgt:       color.RGBA{0x2E, 0x8B, 0x57, 0xCC},
			WhiteTint: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
			BlackTint: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
			BlackLift: 1,
		},
		{
			ID: "forest", Name: "Forest",
			Blurb:     "Tournament green, natural pieces.",
			Light:     color.RGBA{0xEE, 0xEE, 0xD2, 0xFF},
			Dark:      color.RGBA{0x76, 0x96, 0x56, 0xFF},
			Hi:        hi(color.RGBA{0xF6, 0xF1, 0x6D, 0xFF}),
			Tgt:       color.RGBA{0x1B, 0x5E, 0x20, 0xCC},
			WhiteTint: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
			BlackTint: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
			BlackLift: 1,
		},
		{
			ID: "harbor", Name: "Harbor",
			Blurb:     "Coral dark army on dock-tan squares.",
			Light:     color.RGBA{0xE1, 0xCE, 0x98, 0xFF},
			Dark:      color.RGBA{0x6A, 0x5C, 0x48, 0xFF},
			Hi:        hi(color.RGBA{0xE7, 0x5A, 0x50, 0xFF}),
			Tgt:       color.RGBA{0xE7, 0x5A, 0x50, 0xCC},
			WhiteTint: color.RGBA{0xFF, 0xF6, 0xE8, 0xFF},
			BlackTint: color.RGBA{0xE7, 0x5A, 0x50, 0xFF},
			BlackLift: 7,
		},
		{
			ID: "gold", Name: "Gold rush",
			Blurb:     "Black-gold pieces on claim-stake squares.",
			Light:     color.RGBA{0xED, 0xC5, 0x31, 0xFF},
			Dark:      color.RGBA{0x5C, 0x3E, 0x0A, 0xFF},
			Hi:        hi(color.RGBA{0xFF, 0xEE, 0x69, 0xFF}),
			Tgt:       color.RGBA{0xC9, 0xA2, 0x27, 0xCC},
			WhiteTint: color.RGBA{0xFF, 0xF4, 0xC2, 0xFF},
			BlackTint: color.RGBA{0x8A, 0x6A, 0x12, 0xFF},
			BlackLift: 6,
		},
		{
			ID: "ocean", Name: "Ocean",
			Blurb:     "Navy pieces on sea-glass.",
			Light:     color.RGBA{0xD6, 0xEA, 0xF8, 0xFF},
			Dark:      color.RGBA{0x2E, 0x75, 0xA3, 0xFF},
			Hi:        hi(color.RGBA{0x7F, 0xDB, 0xFF, 0xFF}),
			Tgt:       color.RGBA{0x00, 0x7A, 0xA5, 0xCC},
			WhiteTint: color.RGBA{0xF2, 0xFB, 0xFF, 0xFF},
			BlackTint: color.RGBA{0x0B, 0x3D, 0x5C, 0xFF},
			BlackLift: 6.5,
		},
		{
			ID: "neon", Name: "Neon",
			Blurb:     "Ice vs magenta on a night board.",
			Light:     color.RGBA{0x3A, 0x3A, 0x4A, 0xFF},
			Dark:      color.RGBA{0x1A, 0x1A, 0x28, 0xFF},
			Hi:        hi(color.RGBA{0xE0, 0x40, 0xFB, 0xFF}),
			Tgt:       color.RGBA{0x00, 0xE5, 0xFF, 0xCC},
			WhiteTint: color.RGBA{0xB3, 0xE5, 0xFC, 0xFF},
			BlackTint: color.RGBA{0xE0, 0x40, 0xFB, 0xFF},
			BlackLift: 8,
		},
	}
}

func skinByID(id string) Skin {
	for _, s := range skins() {
		if s.ID == id {
			return s
		}
	}
	return skins()[0]
}
