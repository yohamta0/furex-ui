package furex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type viewRefHandler struct {
	ViewRef

	updated    bool
	updateView *View
}

var (
	_ ViewBinder = (*viewRefHandler)(nil)
	_ Updater    = (*viewRefHandler)(nil)
)

func (h *viewRefHandler) Update(v *View) {
	h.updated = true
	h.updateView = h.View()
}

// eventOnlyViewRefHandler has no Update, so only the binding done by furex
// gives it access to its view.
type eventOnlyViewRefHandler struct {
	ViewRef

	enterView *View
}

func (h *eventOnlyViewRefHandler) HandlePointerEnter(x, y int) bool {
	h.enterView = h.View()
	return true
}

func (h *eventOnlyViewRefHandler) HandlePointerLeave() {}

func TestViewRefIsNilBeforeFirstUpdate(t *testing.T) {
	h := &viewRefHandler{}
	root := &View{Width: 100, Height: 100, Handler: h}

	require.Nil(t, h.View())

	root.Update()

	require.Equal(t, root, h.View())
}

func TestViewRefBindsBeforeHandlerUpdate(t *testing.T) {
	h := &viewRefHandler{}
	root := &View{Width: 100, Height: 100, Handler: h}

	root.Update()

	require.True(t, h.updated)
	require.Equal(t, root, h.updateView)
}

func TestViewRefBindsBeforeEvents(t *testing.T) {
	setPrimaryPointerSource(t, &plainPointerSource{x: 5, y: 5})
	h := &eventOnlyViewRefHandler{}
	child := &View{Width: 10, Height: 10, Handler: h}
	root := &View{Width: 100, Height: 100}
	root.AddChild(child)

	root.Update()

	require.Equal(t, child, h.enterView)
}

func TestViewRefStaysBoundAcrossTicks(t *testing.T) {
	h := &viewRefHandler{}
	root := &View{Width: 100, Height: 100, Handler: h}

	require.NotPanics(t, func() {
		root.Update()
		root.Update()
	})
	require.Equal(t, root, h.View())
}

func TestViewRefRebindsAfterHandlerReplacement(t *testing.T) {
	h1 := &viewRefHandler{}
	h2 := &viewRefHandler{}
	root := &View{Width: 100, Height: 100, Handler: h1}
	root.Update()

	root.Handler = h2
	root.Update()

	require.Equal(t, root, h2.View())
}

func TestViewRefPanicsWhenBoundToDifferentView(t *testing.T) {
	h := &viewRefHandler{}
	v1 := &View{Width: 10, Height: 10, Handler: h}
	v1.Update()

	v2 := &View{Width: 10, Height: 10, Handler: h}
	require.PanicsWithValue(t,
		"furex: ViewRef is already bound to another view; create one handler instance per view (e.g. with a ComponentsMap factory)",
		func() { v2.Update() })
}

func TestViewRefBindNilClearsBinding(t *testing.T) {
	h := &viewRefHandler{}
	v1 := &View{Width: 10, Height: 10, Handler: h}
	v1.Update()

	h.BindView(nil)
	require.Nil(t, h.View())

	v2 := &View{Width: 10, Height: 10, Handler: h}
	require.NotPanics(t, func() { v2.Update() })
	require.Equal(t, v2, h.View())
}

func TestViewRefBindsParsedViews(t *testing.T) {
	h := &viewRefHandler{}
	root := Parse(`<body><box id="box"></box></body>`, &ParseOptions{
		Width:      100,
		Height:     100,
		Components: ComponentsMap{"box": h},
	})

	root.Update()

	box, ok := root.GetByID("box")
	require.True(t, ok)
	require.Equal(t, box, h.View())
}

func TestViewRefSharedInstancePanics(t *testing.T) {
	h := &viewRefHandler{}
	root := Parse(`<body><box><box></box></box></body>`, &ParseOptions{
		Components: ComponentsMap{"box": h},
	})

	require.Panics(t, func() { root.Update() })
}
