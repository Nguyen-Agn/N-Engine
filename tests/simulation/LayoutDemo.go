package main

import (
	"image/color"
	"log"

	"autoworld/domain"
	"autoworld/domain/bridge"
	"autoworld/modules/napi"
	layout "autoworld/modules/nlayout"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenW = 800
	screenH = 600
)

type ColorBox struct {
	napi.IObject
	napi.Spr
	napi.Pos

	// col lÃ  mÃ u fill cá»§a Ã´
	col color.RGBA
	// w, h lÃ  kÃ­ch thÆ°á»›c Ã´ (pixel)
	w, h int
}

func NewColorBox(name string, x, y, w, h int, col color.RGBA) *ColorBox {
	b := &ColorBox{col: col, w: w, h: h}
	napi.Obj.NewObjectAndResgiter(b, name, "pos spr", "main")
	b.SetX(float32(x))
	b.SetY(float32(y))

	return b
}

func (b *ColorBox) Create() {
	img := ebiten.NewImage(b.w, b.h)
	img.Fill(b.col)

	spr := newSolidSprite(img, b.w, b.h)
	b.SetSprite("default", spr)
	b.SetCurrentSprite("default")
	b.SetScaleX(1)
	b.SetScaleY(1)

	b.SetOffsetX(float32(b.w) / 2)
	b.SetOffsetY(float32(b.h) / 2)
}

type BackgroundBox struct {
	napi.IObject
	napi.Back
}

func NewBackground() *BackgroundBox {
	bg := &BackgroundBox{}
	napi.Obj.NewObject(bg, "backg", "back sce-main")
	bg.SetColor(color.RGBA{18, 18, 28, 255})

	return bg
}

func (b *BackgroundBox) Create() {}

type HeaderBar struct {
	napi.IObject
	napi.Spr
	napi.Pos

	w, h int
	col  color.RGBA
}

func NewHeaderBar(name string, x, y, w, h int, col color.RGBA) *HeaderBar {
	hb := &HeaderBar{w: w, h: h, col: col}
	napi.Obj.NewObject(hb, name, "pos spr sce-main")
	hb.SetX(float32(x))
	hb.SetY(float32(y))
	return hb
}

func (hb *HeaderBar) Create() {
	img := ebiten.NewImage(hb.w, hb.h)
	img.Fill(hb.col)
	spr := newSolidSprite(img, hb.w, hb.h)
	hb.SetSprite("default", spr)
	hb.SetCurrentSprite("default")

	// Äáº·t Offset Ä‘á»ƒ render top-left
	hb.SetOffsetX(float32(hb.w) / 2)
	hb.SetOffsetY(float32(hb.h) / 2)
}

func setupLayoutDemo() {
	buildSection(
		"[1] Row Â· JustifyCenter Â· AlignCenter",
		0, 20, 400, 270,
		layout.DirRow|layout.AlignCenter|layout.JustifyCenter,
		10,
		[]boxSpec{
			{"r1-a", 80, 60, color.RGBA{99, 179, 237, 255}},
			{"r1-b", 100, 90, color.RGBA{154, 117, 234, 255}},
			{"r1-c", 60, 50, color.RGBA{72, 201, 176, 255}},
		},
	)

	buildSection(
		"[2] Row Â· JustifyEnd Â· AlignEnd",
		400, 20, 400, 270,
		layout.DirRow|layout.AlignEnd|layout.JustifyEnd,
		8,
		[]boxSpec{
			{"r2-a", 80, 70, color.RGBA{248, 113, 113, 255}},
			{"r2-b", 90, 50, color.RGBA{251, 191, 36, 255}},
			{"r2-c", 70, 80, color.RGBA{52, 211, 153, 255}},
		},
	)

	buildSection(
		"[3] Column Â· JustifyCenter Â· AlignCenter",
		0, 320, 400, 270,
		layout.DirColumn|layout.AlignCenter|layout.JustifyCenter,
		12,
		[]boxSpec{
			{"r3-a", 200, 50, color.RGBA{167, 243, 208, 255}},
			{"r3-b", 280, 40, color.RGBA{253, 186, 116, 255}},
			{"r3-c", 150, 55, color.RGBA{196, 181, 253, 255}},
		},
	)

	buildSection(
		"[4] Column Â· JustifyStart Â· AlignEnd",
		400, 320, 400, 270,
		layout.DirColumn|layout.AlignEnd|layout.JustifyStart,
		10,
		[]boxSpec{
			{"r4-a", 180, 50, color.RGBA{103, 232, 249, 255}},
			{"r4-b", 240, 45, color.RGBA{249, 168, 212, 255}},
			{"r4-c", 120, 60, color.RGBA{134, 239, 172, 255}},
		},
	)

	NewHeaderBar("divider-h", 0, 300, screenW, 4, color.RGBA{60, 60, 80, 255})
	NewHeaderBar("divider-v", 400, 0, 4, screenH, color.RGBA{60, 60, 80, 255})

	NewHeaderBar("header", 0, 0, screenW, 20, color.RGBA{30, 30, 50, 255})
}

func newSolidSprite(img *ebiten.Image, w, h int) domain.ISpriteLW {
	spr := bridge.NewSpriteLW(w, h)
	spr.AddImage(img)
	return spr
}

type boxSpec struct {
	name string
	w, h int
	col  color.RGBA
}

func buildSection(sectionName string, rx, ry, rw, rh int, config, gap int, specs []boxSpec) {
	// Viá»n section (1px dÃ y, dÃ¹ng HeaderBar Ä‘Æ¡n giáº£n)
	NewHeaderBar(sectionName+"-border-t", rx, ry, rw, 1, color.RGBA{70, 70, 100, 255})
	NewHeaderBar(sectionName+"-border-b", rx, ry+rh-1, rw, 1, color.RGBA{70, 70, 100, 255})
	NewHeaderBar(sectionName+"-border-l", rx, ry, 1, rh, color.RGBA{70, 70, 100, 255})
	NewHeaderBar(sectionName+"-border-r", rx+rw-1, ry, 1, rh, color.RGBA{70, 70, 100, 255})

	root := layout.NewDivWithNameAndWidthAndHeight("root-"+sectionName, rx, ry, rw, rh, "px")
	root.SetLayoutConfig(config)
	root.SetGap(gap)

	type entry struct {
		spec boxSpec
		a    *layout.A
	}
	entries := make([]entry, 0, len(specs))

	for _, spec := range specs {
		// A khÃ´ng cÃ³ object â†’ chá»‰ dÃ¹ng Ä‘á»ƒ tÃ­nh tá»a Ä‘á»™
		anchor := layout.NewA(spec.name, nil, spec.w, spec.h)
		root.AddChildren(anchor)
		entries = append(entries, entry{spec: spec, a: anchor})
	}

	// TÃ­nh toÃ¡n layout â†’ anchor.BoxX(), BoxY() sáº½ chá»©a tá»a Ä‘á»™ chÃ­nh xÃ¡c
	root.ComputeLayout(rx, ry, rw, rh)

	// Táº¡o ColorBox táº¡i Ä‘Ãºng tá»a Ä‘á»™ layout Ä‘Ã£ tÃ­nh
	for _, e := range entries {
		NewColorBox(e.spec.name, e.a.BoxX(), e.a.BoxY(), e.spec.w, e.spec.h, e.spec.col)
	}
}

func main() {
	// 1. Khá»Ÿi táº¡o engine
	napi.Game.Init(napi.GameConfig{
		Title:  "AutoWorld â€” nLayout Visual Demo",
		Width:  screenW,
		Height: screenH,
	})

	// 2. Táº¡o scene (chá»‰ Physical Map, khÃ´ng cáº§n GUI vÃ¬ layout tÃ­nh tá»a Ä‘á»™ screen-space)
	_, err := napi.Scene.NewSceneAndGo("main", "gui-800-600")
	if err != nil {
		log.Fatalf("KhÃ´ng thá»ƒ khá»Ÿi táº¡o Scene: %v", err)
	}

	// 3. Ná»n + cÃ¡c Ã´ layout
	NewBackground()
	setupLayoutDemo()

	// 4. Cháº¡y vÃ²ng láº·p game
	napi.Game.GameStart()
}
