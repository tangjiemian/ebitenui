package widget

import "image"

type rowLayout struct {
	direction Direction
	padding   Insets
	spacing   int
	dirty     bool
}

type RowLayoutOpt func(f *rowLayout)

type RowLayoutData struct {
	Position  RowLayoutPosition
	Stretch   bool
	MaxWidth  int
	MaxHeight int
}

type RowLayoutPosition int

const (
	RowLayoutPositionStart = RowLayoutPosition(iota)
	RowLayoutPositionCenter
	RowLayoutPositionEnd
)

const RowLayoutOpts = rowLayoutOpts(true)

type rowLayoutOpts bool

func NewRowLayout(opts ...RowLayoutOpt) Layouter {
	r := &rowLayout{}

	for _, o := range opts {
		o(r)
	}

	return r
}

func (o rowLayoutOpts) Direction(d Direction) RowLayoutOpt {
	return func(r *rowLayout) {
		r.direction = d
	}
}

func (o rowLayoutOpts) Padding(i Insets) RowLayoutOpt {
	return func(r *rowLayout) {
		r.padding = i
	}
}

func (o rowLayoutOpts) Spacing(s int) RowLayoutOpt {
	return func(f *rowLayout) {
		f.spacing = s
	}
}

func (f *rowLayout) PreferredSize(widgets []HasWidget) (int, int) {
	r := image.Rectangle{}
	f.layout(widgets, image.Rectangle{}, false, func(w HasWidget, wr image.Rectangle) {
		r = r.Union(wr)
	})
	return r.Dx() + f.padding.Dx(), r.Dy() + f.padding.Dy()
}

func (r *rowLayout) Layout(widgets []HasWidget, rect image.Rectangle) {
	if !r.dirty {
		return
	}

	defer func() {
		r.dirty = false
	}()

	r.layout(widgets, rect, true, func(w HasWidget, wr image.Rectangle) {
		w.(Locateable).SetLocation(wr)
	})
}

func (r *rowLayout) layout(widgets []HasWidget, rect image.Rectangle, usePosition bool, locationFunc func(w HasWidget, wr image.Rectangle)) {
	if len(widgets) == 0 {
		return
	}

	rect = r.padding.Apply(rect)
	x, y := 0, 0

	for _, widget := range widgets {
		wx, wy := x, y
		ww, wh := r.preferredSizeWidget(widget)

		ld := widget.GetWidget().LayoutData
		if rld, ok := ld.(*RowLayoutData); ok {
			wx, wy, ww, wh = r.applyLayoutData(rld, wx, wy, ww, wh, usePosition, rect, x, y)
		}

		if _, ok := widget.(Locateable); ok {
			wr := image.Rect(0, 0, ww, wh)
			wr = wr.Add(rect.Min)
			wr = wr.Add(image.Point{wx, wy})
			locationFunc(widget, wr)
		}

		if r.direction == DirectionHorizontal {
			x += ww + r.spacing
		} else {
			y += wh + r.spacing
		}
	}
}

func (r *rowLayout) preferredSizeWidget(w HasWidget) (int, int) {
	var ww int
	var wh int
	if p, ok := w.(PreferredSizer); ok {
		ww, wh = p.PreferredSize()
	} else {
		ww, wh = 50, 50
	}
	return ww, wh
}

func (r *rowLayout) applyLayoutData(ld *RowLayoutData, wx int, wy int, ww int, wh int, usePosition bool, rect image.Rectangle, x int, y int) (int, int, int, int) {
	if usePosition {
		ww, wh = r.applyStretch(ld, ww, wh, rect)
	}

	ww, wh = r.applyMaxSize(ld, ww, wh)

	if usePosition {
		wx, wy = r.applyPosition(ld, wx, wy, ww, wh, rect, x, y)
	}

	return wx, wy, ww, wh
}

func (r *rowLayout) applyStretch(ld *RowLayoutData, ww int, wh int, rect image.Rectangle) (int, int) {
	if !ld.Stretch {
		return ww, wh
	}

	if r.direction == DirectionHorizontal {
		wh = rect.Dy()
	} else {
		ww = rect.Dx()
	}

	return ww, wh
}

func (r *rowLayout) applyMaxSize(ld *RowLayoutData, ww int, wh int) (int, int) {
	if ld.MaxWidth > 0 && ww > ld.MaxWidth {
		ww = ld.MaxWidth
	}

	if ld.MaxHeight > 0 && wh > ld.MaxHeight {
		wh = ld.MaxHeight
	}

	return ww, wh
}

func (r *rowLayout) applyPosition(ld *RowLayoutData, wx int, wy int, ww int, wh int, rect image.Rectangle, x int, y int) (int, int) {
	switch ld.Position {
	case RowLayoutPositionCenter:
		if r.direction == DirectionHorizontal {
			wy = y + (rect.Dy()-wh)/2
		} else {
			wx = x + (rect.Dx()-ww)/2
		}

	case RowLayoutPositionEnd:
		if r.direction == DirectionHorizontal {
			wy = y + rect.Dy() - wh
		} else {
			wx = x + rect.Dx() - ww
		}
	}

	return wx, wy
}

func (r *rowLayout) MarkDirty() {
	r.dirty = true
}
