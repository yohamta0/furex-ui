package furex

import (
	"image"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddChildUpdateRemove(t *testing.T) {
	view := &View{
		Width:      100,
		Height:     100,
		Direction:  Row,
		Justify:    JustifyStart,
		AlignItems: AlignItemStart,
	}

	mock := &mockHandler{}
	child := &View{
		Width:   10,
		Height:  10,
		Handler: mock,
	}
	require.Equal(t, view, view.AddChild(child))
	require.True(t, view.isDirty)

	view.Update()
	require.True(t, mock.IsUpdated)

	view.Draw(nil)
	require.Equal(t, image.Rect(0, 0, 10, 10), mock.Frame)

	require.True(t, view.RemoveChild(child))
	require.Equal(t, 0, len(view.children))
}

func TestUpdateWithSize(t *testing.T) {
	view := &View{
		Width:      100,
		Height:     100,
		Direction:  Row,
		Justify:    JustifyCenter,
		AlignItems: AlignItemCenter,
	}

	mock := &mockHandler{}
	child := &View{
		Width:   10,
		Height:  10,
		Handler: mock,
	}
	require.Equal(t, view, view.AddChild(child))

	view.UpdateWithSize(200, 200)
	require.True(t, mock.IsUpdated)

	view.Draw(nil)
	require.Equal(t, image.Rect(95, 95, 105, 105), mock.Frame)

}

func TestAddToParent(t *testing.T) {
	root := &View{
		Width:      100,
		Height:     100,
		Direction:  Row,
		Justify:    JustifyStart,
		AlignItems: AlignItemStart,
	}

	mock := &mockHandler{}

	child := (&View{
		Width:   10,
		Height:  10,
		Handler: mock,
	})

	require.Equal(t, child, child.AddTo(root))

	root.Update()
	require.True(t, mock.IsUpdated)

}

func TestAddChild(t *testing.T) {
	view := &View{
		Width:      100,
		Height:     100,
		Direction:  Row,
		Justify:    JustifyStart,
		AlignItems: AlignItemStart,
	}

	mocks := [2]mockHandler{}
	require.Equal(t, view, view.AddChild(
		&View{
			Width:   10,
			Height:  10,
			Handler: &mocks[0],
		},
		&View{
			Width:   10,
			Height:  10,
			Handler: &mocks[1],
		},
	))

	view.Update()
	require.True(t, mocks[0].IsUpdated)
	require.True(t, mocks[1].IsUpdated)

	view.Draw(nil)
	require.Equal(t, image.Rect(0, 0, 10, 10), mocks[0].Frame)
	require.Equal(t, image.Rect(10, 0, 20, 10), mocks[1].Frame)

	view.RemoveAll()
	require.Equal(t, 0, len(view.children))
}

func TestExpandBoxValues(t *testing.T) {
	v := &View{}

	tests := []struct {
		name   string
		values []int
		top    int
		right  int
		bottom int
		left   int
		ok     bool
	}{
		{
			name:   "one value applies to all sides",
			values: []int{10},
			top:    10, right: 10, bottom: 10, left: 10,
			ok: true,
		},
		{
			name:   "two values vertical and horizontal",
			values: []int{10, 20},
			top:    10, right: 20, bottom: 10, left: 20,
			ok: true,
		},
		{
			name:   "three values top horizontal bottom",
			values: []int{10, 20, 30},
			top:    10, right: 20, bottom: 30, left: 20,
			ok: true,
		},
		{
			name:   "four values clockwise",
			values: []int{10, 20, 30, 40},
			top:    10, right: 20, bottom: 30, left: 40,
			ok: true,
		},
		{
			name:   "no values returns ok false",
			values: nil,
			ok:     false,
		},
		{
			name:   "too many values returns ok false",
			values: []int{1, 2, 3, 4, 5},
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			top, right, bottom, left, ok := v.expandBoxValues(tt.values...)

			require.Equal(t, tt.top, top)
			require.Equal(t, tt.right, right)
			require.Equal(t, tt.bottom, bottom)
			require.Equal(t, tt.left, left)
			require.Equal(t, tt.ok, ok)
		})
	}
}

func TestView_SetMargin(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		top    int
		right  int
		bottom int
		left   int
	}{
		{
			name:   "single value",
			values: []int{10},
			top:    10,
			right:  10,
			bottom: 10,
			left:   10,
		},
		{
			name:   "two values",
			values: []int{10, 20},
			top:    10,
			right:  20,
			bottom: 10,
			left:   20,
		},
		{
			name:   "three values",
			values: []int{10, 20, 30},
			top:    10,
			right:  20,
			bottom: 30,
			left:   20,
		},
		{
			name:   "four values",
			values: []int{10, 20, 30, 40},
			top:    10,
			right:  20,
			bottom: 30,
			left:   40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &View{}

			v.SetMargin(tt.values...)

			require.Equal(t, tt.top, v.MarginTop)
			require.Equal(t, tt.right, v.MarginRight)
			require.Equal(t, tt.bottom, v.MarginBottom)
			require.Equal(t, tt.left, v.MarginLeft)
		})
	}
}

type CountingHandler struct {
	Times int
}

func (t *CountingHandler) Update(v *View) {
	t.Times++
}

func TestUpdateOnlyOnce(t *testing.T) {
	rootHandler := &CountingHandler{}
	nestedHandler := &CountingHandler{}

	// given
	rootView := &View{
		Handler: rootHandler,
	}
	nestedView := &View{
		Handler: nestedHandler,
	}
	rootView.addChild(nestedView)

	// when
	rootView.Update()

	// then
	require.True(t, rootHandler.Times == 1)
	require.True(t, nestedHandler.Times == 1)
}
