package furex

import (
	"image"
	"testing"

	"github.com/stretchr/testify/require"
)

func nestedContainer(parentHandler, childHandler Handler) *containerEmbed {
	childView := &View{Handler: childHandler}
	parentView := &View{Handler: parentHandler}
	parentView.children = []*child{{
		item:   childView,
		bounds: image.Rect(50, 50, 150, 150),
	}}
	return &containerEmbed{children: []*child{{
		item:   parentView,
		bounds: image.Rect(0, 0, 200, 200),
	}}}
}

type enterLeaveHandler struct {
	enters int
	leaves int
}

func (h *enterLeaveHandler) HandlePointerEnter(x, y int) bool {
	h.enters++
	return true
}

func (h *enterLeaveHandler) HandlePointerLeave() {
	h.leaves++
}

type swipeHandler struct {
	swiped bool
	dir    SwipeDirection
}

func (h *swipeHandler) HandleSwipe(dir SwipeDirection) {
	h.swiped = true
	h.dir = dir
}

func TestTouchNestedChildConsumesPress(t *testing.T) {
	parentHandler := &mockHandler{}
	childHandler := &mockHandler{}
	ct := nestedContainer(parentHandler, childHandler)

	require.True(t, ct.HandleJustPressedTouchID(1, 100, 100))
	require.True(t, childHandler.IsPressed)
	require.False(t, parentHandler.IsPressed)
}

func TestTouchSwipeStartingOnNestedChild(t *testing.T) {
	parentHandler := &swipeHandler{}
	childHandler := &mockHandler{}
	ct := nestedContainer(parentHandler, childHandler)

	require.True(t, ct.HandleJustPressedTouchID(1, 100, 100))
	require.True(t, childHandler.IsPressed)

	ct.HandleJustReleasedTouchID(1, 160, 100)
	require.True(t, parentHandler.swiped)
	require.Equal(t, SwipeDirectionRight, parentHandler.dir)
}

func TestPointerEnterNestedChildKeepsParentEntered(t *testing.T) {
	parentHandler := &enterLeaveHandler{}
	childHandler := &enterLeaveHandler{}
	ct := nestedContainer(parentHandler, childHandler)

	require.True(t, ct.handlePointerEnterLeave(100, 100))
	require.Equal(t, 1, parentHandler.enters)
	require.Equal(t, 1, childHandler.enters)

	require.False(t, ct.handlePointerEnterLeave(300, 300))
	require.Equal(t, 1, parentHandler.leaves)
	require.Equal(t, 1, childHandler.leaves)
}

func TestIsInside(t *testing.T) {
	for _, tt := range []struct {
		r    image.Rectangle
		x, y int
		want bool
	}{
		{
			r: image.Rect(0, 0, 10, 10),
			x: 0, y: 0,
			want: true,
		},
		{
			r: image.Rect(10, 10, 10, 10),
			x: 10, y: 10,
			want: true,
		},
		{
			r: image.Rect(10, 10, 20, 20),
			x: 20, y: 20,
			want: true,
		},
	} {
		require.Equal(t, tt.want, isInside(&tt.r, tt.x, tt.y))
	}
}
