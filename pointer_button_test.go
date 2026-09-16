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

func TestButtonInputOwnership(t *testing.T) {
	newState := func() (*containerEmbed, *child, *mockHandler, *image.Rectangle) {
		h := &mockHandler{}
		c := &child{
			item:           &View{Handler: h},
			bounds:         image.Rect(10, 10, 30, 30),
			handledTouchID: -1,
		}
		ct := &containerEmbed{children: []*child{c}}
		return ct, c, h, ct.childFrame(c)
	}

	t.Run("touch ignores pointer release", func(t *testing.T) {
		ct, c, h, frame := newState()

		require.True(t, c.HandleJustPressedTouchID(frame, 1, 20, 20))
		ct.handlePointerButtonReleased(pointerPrimary, 20, 20)
		require.False(t, h.IsReleased)

		c.HandleJustReleasedTouchID(frame, 1, 20, 20)
		require.True(t, h.IsReleased)
	})

	t.Run("pointer ignores touch release", func(t *testing.T) {
		ct, c, h, frame := newState()

		require.True(t, ct.handlePointerButtonPressed(pointerPrimary, 20, 20))
		c.HandleJustReleasedTouchID(frame, 1, 20, 20)
		require.False(t, h.IsReleased)

		ct.handlePointerButtonReleased(pointerPrimary, 20, 20)
		require.True(t, h.IsReleased)
	})
}
