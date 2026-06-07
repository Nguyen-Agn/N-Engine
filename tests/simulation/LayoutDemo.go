package main

import (
	"image/color"
	"log"

	"autoworld/domain"
	"autoworld/domain/bridge"
	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
	layout "autoworld/modules/nlayout"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenW = 800
	screenH = 600
)

type ColorBox struct {
	ncom.Object
	ncom.Spr
	ncom.Pos
	ncom.Box

	col color.RGBA
}

func NewColorBox(name string, x, y, w, h int, col color.RGBA) *ColorBox {
	b := &ColorBox{col: col}
	napi.Obj.NewObject(b, name, "pos spr sce-main box")
	b.SetX(float32(x))
	b.SetY(float32(y))

	b.SetBoxW(float32(w))
	b.SetBoxH(float32(h))
	return b
}

func (b *ColorBox) Create() {
	b.SetSprite("default", napi.Assert.GetSprite("nine-slice"))
	b.SetCurrentSprite("default")

	b.SetScaleX(b.BoxW() / float32(b.GetCurrentSprite().Width()))
	b.SetScaleY(b.BoxH() / float32(b.GetCurrentSprite().Height()))

	b.Set9Slice(true, "4")

	b.SetOffsetX(float32(b.BoxW()) / 2)
	b.SetOffsetY(float32(b.BoxH()) / 2)
}

type BackgroundBox struct {
	napi.IObject
	ncom.Back
}

type HeaderBar struct {
	napi.IObject
	ncom.Spr
	ncom.Pos
	ncom.Box
	col color.RGBA
}

func NewHeaderBar(name string, x, y, w, h int, col color.RGBA) *HeaderBar {
	hb := &HeaderBar{col: col}
	napi.Obj.NewObject(hb, name, "pos spr sce-main box")
	hb.SetX(float32(x))
	hb.SetY(float32(y))

	hb.SetBoxW(float32(w))
	hb.SetBoxH(float32(h))

	return hb
}

func (hb *HeaderBar) Create() {
	img := ebiten.NewImage(int(hb.BoxW()), int(hb.BoxH()))
	img.Fill(hb.col)
	spr := newSolidSprite(img, int(hb.BoxW()), int(hb.BoxH()))
	hb.SetSprite("default", spr)
	hb.SetCurrentSprite("default")

	hb.SetOffsetX(float32(hb.BoxW()) / 2)
	hb.SetOffsetY(float32(hb.BoxH()) / 2)
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

		anchor := layout.NewA(spec.name, nil, spec.w, spec.h)
		root.AddChildren(anchor)
		entries = append(entries, entry{spec: spec, a: anchor})
	}

	root.ComputeLayout(rx, ry, rw, rh)

	for _, e := range entries {
		NewColorBox(e.spec.name, e.a.BoxX(), e.a.BoxY(), e.spec.w, e.spec.h, e.spec.col)
	}
}

func main() {
	napi.Game.Init(napi.GameConfig{
		Title:  "AutoWorld nLayout Visual Demo",
		Width:  screenW,
		Height: screenH,
	})

	napi.Game.LoadFromFile("./SharedObject/Config.toml")

	_, err := napi.Scene.NewSceneAndGo("main", "gui-800-600")
	if err != nil {
		log.Fatalf("No Scene: %v", err)
	}

	setupLayoutDemo()

	// 4. Cháº¡y vÃ²ng láº·p game
	napi.Game.GameStart()
}
