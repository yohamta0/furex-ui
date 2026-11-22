// Scroll shows two independent scrollable lists built from HTML.
//
// Wheel and click callbacks receive only the pointer position, so the list
// and row handlers embed furex.ViewRef to reach their own view: the list
// reads its frame and rows to clamp scrolling, and a row finds the list it
// belongs to when clicked.
package main

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/furex/v2"
)

const (
	screenWidth    = 480
	screenHeight   = 360
	listWidth      = 200
	rowHeight      = 30
	scrollbarWidth = 6
)

var (
	backgroundColor = color.RGBA{25, 30, 40, 255}
	listColor       = color.RGBA{40, 46, 60, 255}
	rowColor        = color.RGBA{60, 70, 90, 255}
	hoverColor      = color.RGBA{80, 95, 125, 255}
	selectedColor   = color.RGBA{70, 130, 220, 255}
	thumbColor      = color.RGBA{150, 160, 180, 255}
)

var (
	fruits = []string{
		"Apple", "Apricot", "Banana", "Blueberry", "Cherry", "Coconut",
		"Grape", "Kiwi", "Lemon", "Lime", "Mango", "Melon", "Orange",
		"Peach", "Pear", "Plum", "Raspberry", "Strawberry",
	}
	colors = []string{
		"Red", "Orange", "Yellow", "Green", "Blue", "Indigo", "Violet",
		"Black", "White",
	}
)

type game struct {
	ui *furex.View
}

func newGame() *game {
	html := fmt.Sprintf(`<body>
  <div style="flex-direction: row; justify-content: space-around; align-items: center">
    %s
    %s
  </div>
</body>`, listHTML("fruits", 8, fruits), listHTML("colors", 4, colors))

	ui := furex.Parse(html, &furex.ParseOptions{
		Width:  screenWidth,
		Height: screenHeight,
		// A ViewRef handler belongs to one view, so each tag gets a new
		// handler from a factory.
		Components: furex.ComponentsMap{
			"list": func() furex.Handler { return &list{} },
			"row":  func() furex.Handler { return &row{} },
		},
	})
	return &game{ui: ui}
}

// listHTML returns a list tall enough to show visibleRows of items.
func listHTML(id string, visibleRows int, items []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<list id=%q style="flex-direction: column; width: %dpx; height: %dpx">`,
		id, listWidth, visibleRows*rowHeight)
	for _, item := range items {
		fmt.Fprintf(&b, `<row label=%q style="width: %dpx; height: %dpx"></row>`,
			item, listWidth-scrollbarWidth-2, rowHeight)
	}
	b.WriteString(`</list>`)
	return b.String()
}

func (g *game) Update() error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	g.ui.Draw(screen)
	ebitenutil.DebugPrintAt(screen, "Scroll a list with the wheel. Click a row to select it.", 16, 16)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Fruit: %s   Color: %s",
		selectedLabel(g.ui, "fruits"), selectedLabel(g.ui, "colors")), 16, screenHeight-28)
}

func (g *game) Layout(int, int) (int, int) { return screenWidth, screenHeight }

func selectedLabel(ui *furex.View, listID string) string {
	sel := ui.MustGetByID(listID).Handler.(*list).selected
	if sel == nil {
		return "-"
	}
	return sel.Attrs["label"]
}

// list shows the rows that fit in its frame and scrolls them by whole rows.
type list struct {
	furex.ViewRef

	// offset is the index of the first visible row. It stays fractional so
	// that small touchpad deltas add up.
	offset   float64
	selected *furex.View
}

var (
	_ furex.Drawer       = (*list)(nil)
	_ furex.Updater      = (*list)(nil)
	_ furex.WheelHandler = (*list)(nil)
)

// HandleWheel scrolls within the rows this list can show. It returns false at
// either end so that an enclosing handler may take the scroll.
func (l *list) HandleWheel(_, _ int, _, dy float64) bool {
	v := l.View()
	maxOffset := float64(max(0, rowCount(v)-visibleRows(v)))
	next := min(max(l.offset-dy, 0), maxOffset)
	if next == l.offset {
		return false
	}
	l.offset = next
	return true
}

// Update hides the rows outside the visible window.
func (l *list) Update(v *furex.View) {
	first := int(l.offset)
	last := first + visibleRows(v)
	i := 0
	for r := range v.Children() {
		d := furex.DisplayNone
		if first <= i && i < last {
			d = furex.DisplayFlex
		}
		if r.Display != d {
			r.SetDisplay(d)
		}
		i++
	}
}

func (l *list) Draw(screen *ebiten.Image, frame image.Rectangle, v *furex.View) {
	fillRect(screen, frame, listColor)

	total, visible := rowCount(v), visibleRows(v)
	if total <= visible {
		return
	}
	track := float32(frame.Dy())
	thumbY := float32(frame.Min.Y) + track*float32(int(l.offset))/float32(total)
	thumbH := track * float32(visible) / float32(total)
	vector.FillRect(screen, float32(frame.Max.X-scrollbarWidth), thumbY,
		scrollbarWidth, thumbH, thumbColor, false)
}

func rowCount(v *furex.View) int {
	n := 0
	for range v.Children() {
		n++
	}
	return n
}

func visibleRows(v *furex.View) int {
	return v.Frame().Dy() / rowHeight
}

// row is one selectable entry of a list.
type row struct {
	furex.ViewRef
}

var (
	_ furex.Drawer                      = (*row)(nil)
	_ furex.PointerPrimaryButtonHandler = (*row)(nil)
)

// HandleJustPressedPointerButtonPrimary selects this row in the list that
// contains it.
func (r *row) HandleJustPressedPointerButtonPrimary(int, int) bool {
	v := r.View()
	parentList(v).selected = v
	return true
}

func (r *row) HandleJustReleasedPointerButtonPrimary(int, int) {}

func (r *row) Draw(screen *ebiten.Image, frame image.Rectangle, v *furex.View) {
	c := rowColor
	switch {
	case parentList(v).selected == v:
		c = selectedColor
	case image.Pt(ebiten.CursorPosition()).In(frame):
		c = hoverColor
	}
	// Leave a 1px gap between rows.
	fillRect(screen, image.Rect(frame.Min.X, frame.Min.Y, frame.Max.X, frame.Max.Y-1), c)
	ebitenutil.DebugPrintAt(screen, v.Attrs["label"], frame.Min.X+8, frame.Min.Y+8)
}

func parentList(v *furex.View) *list {
	return v.Parent().Handler.(*list)
}

func fillRect(screen *ebiten.Image, r image.Rectangle, c color.Color) {
	vector.FillRect(screen, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()), c, false)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("furex: ViewRef scroll")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
