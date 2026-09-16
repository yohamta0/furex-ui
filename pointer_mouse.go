package furex

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MousePrimaryPointerSource struct{}
type MouseSecondaryPointerSource struct{}

func NewMousePrimaryPointerSource() PointerSource {
	return MousePrimaryPointerSource{}
}
func NewMouseSecondaryPointerSource() PointerSource {
	return MouseSecondaryPointerSource{}
}

func (m MousePrimaryPointerSource) Update(time.Duration) {
}

func (m MousePrimaryPointerSource) ReadPosition() (int, int) {
	return ebiten.CursorPosition()
}

func (m MousePrimaryPointerSource) IsJustPressed() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (m MousePrimaryPointerSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

func (m MouseSecondaryPointerSource) Update(time.Duration) {
}

func (m MouseSecondaryPointerSource) ReadPosition() (int, int) {
	return ebiten.CursorPosition()
}

func (m MouseSecondaryPointerSource) IsJustPressed() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
}

func (m MouseSecondaryPointerSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight)
}
