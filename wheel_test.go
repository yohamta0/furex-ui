package furex

import (
	"image"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type wheelHandler struct {
	scrolls int
	x, y    int
	dx, dy  float64
	handled bool
}

func (h *wheelHandler) HandleWheel(x, y int, dx, dy float64) bool {
	h.scrolls++
	h.x, h.y = x, y
	h.dx, h.dy = dx, dy
	return h.handled
}

type plainPointerSource struct {
	x, y int
}

func (s *plainPointerSource) Update(time.Duration) {}

func (s *plainPointerSource) ReadPosition() (int, int) {
	return s.x, s.y
}

func (s *plainPointerSource) IsJustPressed() bool {
	return false
}

func (s *plainPointerSource) IsJustReleased() bool {
	return false
}

type wheelPointerSource struct {
	plainPointerSource
	dx, dy float64
}

func (s *wheelPointerSource) ReadWheel() (float64, float64) {
	return s.dx, s.dy
}

func setPrimaryPointerSource(t *testing.T, source PointerSource) {
	t.Helper()
	prev := CurrentPrimaryPointerSource
	CurrentPrimaryPointerSource = source
	t.Cleanup(func() {
		CurrentPrimaryPointerSource = prev
	})
}

func TestWheelHandler(t *testing.T) {
	h := &wheelHandler{handled: true}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	require.True(t, ct.handleWheel(20, 20, 0, 1))
	require.Equal(t, 1, h.scrolls)
	require.Equal(t, 20, h.x)
	require.Equal(t, 20, h.y)
	require.Zero(t, h.dx)
	require.Equal(t, 1.0, h.dy)
}

func TestWheelHandlerMiss(t *testing.T) {
	h := &wheelHandler{handled: true}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	require.False(t, ct.handleWheel(40, 40, 0, 1))
	require.Zero(t, h.scrolls)
}

func TestWheelHandlerNotConsumed(t *testing.T) {
	h := &wheelHandler{}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	require.False(t, ct.handleWheel(20, 20, 0, 1))
	require.Equal(t, 1, h.scrolls)
}

func TestWheelNestedChildConsumesScroll(t *testing.T) {
	parentHandler := &wheelHandler{handled: true}
	childHandler := &wheelHandler{handled: true}
	ct := nestedContainer(parentHandler, childHandler)

	require.True(t, ct.handleWheel(100, 100, 0, -1))
	require.Equal(t, 1, childHandler.scrolls)
	require.Zero(t, parentHandler.scrolls)
}

func TestWheelNestedChildFallsThroughToParent(t *testing.T) {
	parentHandler := &wheelHandler{handled: true}
	childHandler := &wheelHandler{}
	ct := nestedContainer(parentHandler, childHandler)

	require.True(t, ct.handleWheel(100, 100, 0, -1))
	require.Equal(t, 1, childHandler.scrolls)
	require.Equal(t, 1, parentHandler.scrolls)
}

func TestWheelDisplayNoneSkipped(t *testing.T) {
	h := &wheelHandler{handled: true}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h, Display: DisplayNone},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	require.False(t, ct.handleWheel(20, 20, 0, 1))
	require.Zero(t, h.scrolls)
}

func TestWheelEvents(t *testing.T) {
	setPrimaryPointerSource(t, &wheelPointerSource{
		plainPointerSource: plainPointerSource{x: 20, y: 20},
		dy:                 1,
	})
	h := &wheelHandler{handled: true}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	ct.handleWheelEvents()

	require.Equal(t, 1, h.scrolls)
	require.Equal(t, 20, h.x)
	require.Equal(t, 20, h.y)
	require.Equal(t, 1.0, h.dy)
}

func TestWheelEventsZeroOffset(t *testing.T) {
	setPrimaryPointerSource(t, &wheelPointerSource{
		plainPointerSource: plainPointerSource{x: 20, y: 20},
	})
	h := &wheelHandler{handled: true}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	ct.handleWheelEvents()
	require.Zero(t, h.scrolls)
}

func TestWheelEventsWithoutWheelSource(t *testing.T) {
	setPrimaryPointerSource(t, &plainPointerSource{x: 20, y: 20})
	h := &wheelHandler{handled: true}
	ct := &containerEmbed{children: []*child{{
		item:   &View{Handler: h},
		bounds: image.Rect(10, 10, 30, 30),
	}}}

	ct.handleWheelEvents()
	require.Zero(t, h.scrolls)
}
