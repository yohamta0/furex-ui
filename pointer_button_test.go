package furex

import (
	"image"
	"testing"

	"github.com/stretchr/testify/require"
)

type secondaryButtonHandler struct {
	pressed  int
	released int
}

func (h *secondaryButtonHandler) HandleJustPressedPointerButtonSecondary(int, int) bool {
	h.pressed++
	return true
}

func (h *secondaryButtonHandler) HandleJustReleasedPointerButtonSecondary(int, int) {
	h.released++
}

func TestPointerSecondaryButtonHandler(t *testing.T) {
	h := &secondaryButtonHandler{}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	require.True(t, ct.handlePointerButtonPressed(pointerSecondary, 20, 20))
	ct.handlePointerButtonReleased(pointerSecondary, 40, 40)

	require.Equal(t, 1, h.pressed)
	require.Equal(t, 1, h.released)
}

func TestPointerSecondaryButtonMiss(t *testing.T) {
	h := &secondaryButtonHandler{}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	require.False(t, ct.handlePointerButtonPressed(pointerSecondary, 40, 40))
	ct.handlePointerButtonReleased(pointerSecondary, 20, 20)

	require.Zero(t, h.pressed)
	require.Zero(t, h.released)
}
