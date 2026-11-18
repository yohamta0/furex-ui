package main

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/furex/v2"
)

type ScreenSize struct {
	Width  int
	Height int
}

type Game struct {
	initOnce   sync.Once
	screenSize ScreenSize
	gameUI     *furex.View
}

func (g *Game) Update() error {
	g.initOnce.Do(func() {
		g.setupUI()
	})

	width, height := ebiten.WindowSize()
	g.gameUI.UpdateWithSize(width, height)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3d, 0x55, 0x0c, 0xff})
	g.gameUI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenSize.Width = outsideWidth
	g.screenSize.Height = outsideHeight
	return g.screenSize.Width, g.screenSize.Height
}

func NewGame() (*Game, error) {
	game := &Game{}
	return game, nil
}

func (g *Game) setupUI() {
	g.gameUI = &furex.View{
		Width:      g.screenSize.Width,
		Height:     g.screenSize.Height,
		Direction:  furex.Column,
		Justify:    furex.JustifyCenter,
		AlignItems: furex.AlignItemStretch,
	}

	g.gameUI.AddChild(&furex.View{
		Height: 50,
		Handler: &Box{
			Color: color.RGBA{0xdd, 0xdd, 0xdd, 0xff},
		},
	})

	g.gameUI.AddChild(&furex.View{
		Grow: 1,
		Handler: &Box{
			Color: color.RGBA{0x81, 0xb6, 0x22, 0xff},
		},
	})

	g.gameUI.AddChild(&furex.View{
		Height: 50,
		Handler: &Box{
			Color: color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
	})
}

type Box struct {
	Color color.Color
}

func (b *Box) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	vector.FillRect(
		screen,
		float32(frame.Min.X),
		float32(frame.Min.Y),
		float32(frame.Size().X),
		float32(frame.Size().Y),
		b.Color,
		false,
	)
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game, err := NewGame()
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
