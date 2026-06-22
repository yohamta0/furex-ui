package furex

import "time"

var (
	CurrentPointerSource = NewMousePointerSource()
)

type PointerSource interface {
	Update(dt time.Duration)
	ReadPosition() (x, y int)
	IsJustPressed() bool
	IsJustReleased() bool
}

func UpdatePointer(dt time.Duration) {
	CurrentPointerSource.Update(dt)
}
