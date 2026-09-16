package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/furex/v2"
)

type game struct{ ui *furex.View }

func (g *game) Update() error {
	if g.ui == nil {
		g.ui = (&furex.View{Width: 480, Height: 320, Direction: furex.Row, Justify: furex.JustifyCenter, AlignItems: furex.AlignItemCenter}).AddChild(
			&furex.View{Width: 180, Height: 120, Handler: &buttonBox{color.RGBA{70, 130, 220, 255}, "LEFT / RIGHT", 0, 0}},
			&furex.View{Width: 180, Height: 120, Handler: &buttonBox{color.RGBA{220, 120, 70, 255}, "SECOND BOX", 0, 0}},
		)
	}
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{25, 30, 40, 255})
	g.ui.Draw(screen)
}
func (g *game) Layout(w, h int) (int, int) { return w, h }

type buttonBox struct {
	color.Color
	label       string
	left, right int
}

var _ furex.Drawer = (*buttonBox)(nil)
var _ furex.PointerPrimaryButtonHandler = (*buttonBox)(nil)
var _ furex.PointerSecondaryButtonHandler = (*buttonBox)(nil)

func (b *buttonBox) Draw(screen *ebiten.Image, frame image.Rectangle, _ *furex.View) {
	ebitenutil.DrawRect(screen, float64(frame.Min.X), float64(frame.Min.Y), float64(frame.Dx()), float64(frame.Dy()), b.Color)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s\nLeft: %d  Right: %d", b.label, b.left, b.right), frame.Min.X+12, frame.Min.Y+12)
}
func (b *buttonBox) HandleJustPressedPointerButtonPrimary(int, int) bool   { b.left++; return true }
func (b *buttonBox) HandleJustReleasedPointerButtonPrimary(int, int)       {}
func (b *buttonBox) HandleJustPressedPointerButtonSecondary(int, int) bool { b.right++; return true }
func (b *buttonBox) HandleJustReleasedPointerButtonSecondary(int, int)     {}

func main() {
	ebiten.SetWindowSize(480, 320)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
