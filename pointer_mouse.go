package furex

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MousePointerSource struct {
}

func NewMousePointerSource() PointerSource {
	return MousePointerSource{}
}

func (m MousePointerSource) Update(time.Duration) {
}

func (m MousePointerSource) ReadPosition() (x, y int) {
	return ebiten.CursorPosition()
}

func (m MousePointerSource) IsJustPressed() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (m MousePointerSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}
