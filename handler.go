package furex

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Handler represents a component that can be added to a container.
type Handler interface{}

// Drawer represents a component that can be added to a container.
type Drawer interface {
	// Draw function draws the content of the component inside the frame.
	Draw(screen *ebiten.Image, frame image.Rectangle, v *View)
}

// Updater represents a component that updates by one tick.
type Updater interface {
	// Update updates the state of the component by one tick.
	Update(v *View)
}

// ViewBinder is implemented by handlers that need their own View in event
// callbacks. Embed ViewRef to implement it.
// The library calls BindView when the view is added to a parent, which
// includes every view created by Parse, and every tick before the handler's
// Update and before dispatching any events.
type ViewBinder interface {
	BindView(v *View)
}

// ViewRef stores the View its handler is attached to.
// Embed it in a handler to implement ViewBinder, and set that handler as a
// pointer: a value does not implement ViewBinder, so View would stay nil.
//
// View returns nil until the handler is bound: when its view is added to a
// parent or on the next Update, whichever comes first. BindView(nil) clears
// the binding, which is only needed when a handler is explicitly moved to
// another view.
// A handler that embeds ViewRef must not be copied after it has been bound.
type ViewRef struct {
	view *View
}

// BindView binds the handler to v. It panics if the handler is already bound
// to a different view. Passing nil clears the binding.
func (r *ViewRef) BindView(v *View) {
	if v == nil {
		r.view = nil
		return
	}
	if r.view != nil && r.view != v {
		panic(fmt.Sprintf("furex: ViewRef is already bound to another view (tag %q, id %q); create one handler instance per view (e.g. with a ComponentsMap factory)", v.TagName, v.ID))
	}
	r.view = v
}

// View returns the View the handler is attached to.
func (r *ViewRef) View() *View {
	return r.view
}

// DrawHandler represents a component that can be added to a container.
// Deprectead: use Drawer instead
type DrawHandler interface {
	// HandleDraw function draws the content of the component inside the frame.
	// The frame parameter represents the location (x,y) and size (width,height) relative to the window (0,0).
	HandleDraw(screen *ebiten.Image, frame image.Rectangle)
}

// UpdateHandler represents a component that updates by one tick.
// Deprectead: use Updater instead
type UpdateHandler interface {
	// Updater updates the state of the component by one tick.
	HandleUpdate()
}

// ButtonHandler represents a button component.
type ButtonHandler interface {
	// HandlePress handle the event when user just started pressing the button
	// The parameter (x, y) is the location relative to the window (0,0).
	// touchID is the unique ID of the touch.
	// If the button is pressed by a mouse, touchID is -1.
	HandlePress(x, y int, t ebiten.TouchID)

	// HandleRelease handle the event when user just released the button.
	// The parameter (x, y) is the location relative to the window (0,0).
	// The parameter isCancel is true when the touch/left click is released outside of the button.
	HandleRelease(x, y int, isCancel bool)
}

// NotButton represents a component that is not a button.
// TODO: update HandlePress to return bool in the next major version.
type NotButton interface {
	// IsButton returns true if the handler is a button.
	IsButton() bool
}

// TouchHandler represents a component that handle touches.
type TouchHandler interface {
	// HandleJustPressedTouchID handles the touchID just pressed and returns true if it handles the TouchID
	HandleJustPressedTouchID(touch ebiten.TouchID, x, y int) bool
	// HandleJustReleasedTouchID handles the touchID just released
	// Should be called only when it handled the TouchID when pressed
	HandleJustReleasedTouchID(touch ebiten.TouchID, x, y int)
}

// PointerHandler represents a component that handle pointer move.
type PointerHandler interface {
	// HandlePointer handles the pointer move and returns true if it handled the pointer move.
	// The parameter (x, y) is the location relative to the window (0,0).
	HandlePointer(x, y int) bool
}

// PointerPrimaryButtonHandler represents a component that handle mouse button left click.
type PointerPrimaryButtonHandler interface {
	// HandleJustPressedPointerButtonPrimary handle primary pointer button just pressed.
	// The parameter (x, y) is the location relative to the window (0,0).
	// It returns true if it handles the pointer move.
	HandleJustPressedPointerButtonPrimary(x, y int) bool
	// HandleJustReleasedPointerButtonPrimary handles the pointer button just released.
	// The parameter (x, y) is the location relative to the window (0,0).
	HandleJustReleasedPointerButtonPrimary(x, y int)
}

// PointerSecondaryButtonHandler handles secondary pointer (mouse right button) clicks.
type PointerSecondaryButtonHandler interface {
	// HandleJustPressedPointerButtonSecondary returns true to consume the press.
	// Coordinates are relative to the window (0,0).
	HandleJustPressedPointerButtonSecondary(x, y int) bool
	// HandleJustReleasedPointerButtonSecondary is called only after consuming a press,
	// even when released outside the view. Coordinates are relative to the window.
	HandleJustReleasedPointerButtonSecondary(x, y int)
}

// WheelHandler represents a component that handles mouse wheel or touchpad scroll.
type WheelHandler interface {
	// HandleWheel handles the wheel scroll and returns true if it handles the scroll.
	// The parameter (x, y) is the location relative to the window (0,0).
	// The parameters (dx, dy) are the wheel offsets in the current tick.
	HandleWheel(x, y int, dx, dy float64) bool
}

// PointerEnterLeaveHandler represents a component that handle pointer enter/leave.
type PointerEnterLeaveHandler interface {
	// HandlePointerEnter handles the pointer enter.
	HandlePointerEnter(x, y int) bool
	// HandlePointerLeave handles the pointer leave.
	HandlePointerLeave()
}

// SwipeHandler represents different swipe directions.
type SwipeDirection int

const (
	SwipeDirectionLeft SwipeDirection = iota
	SwipeDirectionRight
	SwipeDirectionUp
	SwipeDirectionDown
)

// SwipeHandler represents a component that handle swipe.
type SwipeHandler interface {
	// HandleSwipe handles swipes.
	HandleSwipe(dir SwipeDirection)
}

type handler struct {
	opts HandlerOpts
}

// HandlerOpts represents the options for a handler.
type HandlerOpts struct {
	Update        func(v *View)
	Draw          func(screen *ebiten.Image, frame image.Rectangle, v *View)
	HandlePress   func(x, y int, t ebiten.TouchID)
	HandleRelease func(x, y int, isCancel bool)
}

// NewHandler creates a new handler.
func NewHandler(opts HandlerOpts) Handler {
	return &handler{opts: opts}
}

func (h *handler) Update(v *View) {
	if h.opts.Update != nil {
		h.opts.Update(v)
	}
}

func (h *handler) Draw(screen *ebiten.Image, frame image.Rectangle, v *View) {
	if h.opts.Draw != nil {
		h.opts.Draw(screen, frame, v)
	}
}

func (h *handler) HandlePress(x, y int, t ebiten.TouchID) {
	if h.opts.HandlePress != nil {
		h.opts.HandlePress(x, y, t)
	}
}

func (h *handler) HandleRelease(x, y int, isCancel bool) {
	if h.opts.HandleRelease != nil {
		h.opts.HandleRelease(x, y, isCancel)
	}
}
