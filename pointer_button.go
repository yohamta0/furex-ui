package furex

type pointerButton int

const (
	pointerPrimary pointerButton = iota
	pointerSecondary
	pointerButtonCount
)

// Adapt button-specific handlers to the common hit testing and dispatch path.
func pointerButtonCallbacks(handler Handler, button pointerButton) (press func(int, int) bool, release func(int, int), ok bool) {
	switch button {
	case pointerPrimary:
		if h, ok := handler.(PointerPrimaryButtonHandler); ok {
			return h.HandleJustPressedPointerButtonPrimary, h.HandleJustReleasedPointerButtonPrimary, true
		}
	case pointerSecondary:
		if h, ok := handler.(PointerSecondaryButtonHandler); ok {
			return h.HandleJustPressedPointerButtonSecondary, h.HandleJustReleasedPointerButtonSecondary, true
		}
	}
	return nil, nil, false
}
