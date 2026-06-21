package furex

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GamePadPointerSource struct {
	id          ebiten.GamepadID
	pressButton ebiten.GamepadButton
	speed       float64

	x, y int

	limitBounds bool
	maxX, maxY  int
}

func NewGamePadPointerSource(
	id ebiten.GamepadID,
	pressButton ebiten.GamepadButton,
	speed float64,
) *GamePadPointerSource {
	return &GamePadPointerSource{
		id:          id,
		pressButton: pressButton,
		speed:       speed,
	}
}

func (g *GamePadPointerSource) WithBounds(width, height int) *GamePadPointerSource {
	g.maxX, g.maxY = width, height
	g.limitBounds = true
	return g
}

func (g *GamePadPointerSource) Update(dt time.Duration) {
	horizontal := ebiten.StandardGamepadAxisValue(
		g.id,
		ebiten.StandardGamepadAxisLeftStickHorizontal,
	)
	vertical := ebiten.StandardGamepadAxisValue(
		g.id,
		ebiten.StandardGamepadAxisLeftStickVertical,
	)

	const deadZone = 0.15
	if math.Abs(horizontal) < deadZone {
		horizontal = 0
	}
	if math.Abs(vertical) < deadZone {
		vertical = 0
	}

	g.x += int(horizontal * g.speed * dt.Seconds())
	g.y += int(vertical * g.speed * dt.Seconds())

	g.clampCoordinates()
}

func (g *GamePadPointerSource) clampCoordinates() {
	if !g.limitBounds {
		return
	}

	if g.x < 0 {
		g.x = 0
	}
	if g.x > g.maxX {
		g.x = g.maxX
	}

	if g.y < 0 {
		g.y = 0
	}
	if g.y > g.maxY {
		g.y = g.maxY
	}
}

func (g *GamePadPointerSource) ReadPosition() (x, y int) {
	return g.x, g.y
}

func (g *GamePadPointerSource) IsJustPressed() bool {
	return inpututil.IsGamepadButtonJustPressed(g.id, g.pressButton)
}

func (g *GamePadPointerSource) IsJustReleased() bool {
	return inpututil.IsGamepadButtonJustReleased(g.id, g.pressButton)
}
