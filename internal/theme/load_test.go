package theme

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHarborTOML(t *testing.T) {
	dir := t.TempDir()
	data := `accent = "#e75a50"
foreground = "#efebdc"
background = "#1B1B1B"
`
	if err := os.WriteFile(filepath.Join(dir, "colors.toml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	p, ok := loadFromDir(dir)
	if !ok {
		t.Fatal("parse failed")
	}
	if p.Hex(p.Accent) != "#e75a50" {
		t.Fatalf("accent %s", p.Hex(p.Accent))
	}
}

func TestDefaultIsHarborDark(t *testing.T) {
	t.Setenv("OMARCHY_THEME_DIR", filepath.Join(t.TempDir(), "missing"))
	p := resolve(HarborDark())
	if p.Name != "harbordark" {
		t.Fatalf("name %s", p.Name)
	}
}
