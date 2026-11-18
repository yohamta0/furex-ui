package furex

import "time"

var (
	CurrentPrimaryPointerSource                 = NewMousePrimaryPointerSource()
	CurrentSecondaryPointerSource PointerSource = NewMouseSecondaryPointerSource()
)

type PointerSource interface {
	Update(dt time.Duration)
	ReadPosition() (x, y int)
	IsJustPressed() bool
	IsJustReleased() bool
}

// WheelSource represents an optional pointer source that supports wheel input.
type WheelSource interface {
	// ReadWheel returns the wheel offsets in the current tick.
	ReadWheel() (dx, dy float64)
}

// UpdatePointer updates the primary and secondary input sources.
func UpdatePointer(dt time.Duration) {
	CurrentPrimaryPointerSource.Update(dt)
	if CurrentSecondaryPointerSource != nil {
		CurrentSecondaryPointerSource.Update(dt)
	}
}
