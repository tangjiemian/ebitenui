package gui

import (
	"github.com/blizzy78/ebitenui/widget"
)

type DragContents struct {
	Res *UiResources

	sources []*widget.Widget
	targets []*widget.Widget

	text *widget.Text
}

func (d *DragContents) Create(srcWidget widget.HasWidget, srcX int, srcY int) (widget.DragWidget, interface{}) {
	if !d.IsSource(srcWidget.GetWidget()) {
		return nil, nil
	}

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(d.Res.ToolTip.background),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(d.Res.ToolTip.padding),
		)),
	)

	d.text = widget.NewText(widget.TextOpts.Text("Drag Me!", d.Res.ToolTip.face, d.Res.ToolTip.color))
	c.AddChild(d.text)

	return c, nil
}

func (d *DragContents) Update(target widget.HasWidget, _ int, _ int, _ interface{}) {
	if target != nil && d.IsTarget(target.GetWidget()) {
		d.text.Label = "* DROP ME! *"
	} else {
		d.text.Label = "Drag Me!"
	}
}

func (d *DragContents) AddSource(s widget.HasWidget) {
	d.sources = append(d.sources, s.GetWidget())
}

func (d *DragContents) AddTarget(t widget.HasWidget) {
	d.targets = append(d.targets, t.GetWidget())
}

func (d *DragContents) IsSource(w *widget.Widget) bool {
	for _, s := range d.sources {
		if s == w {
			return true
		}
	}

	p := w.Parent()
	if p == nil {
		return false
	}

	return d.IsSource(p)
}

func (d *DragContents) IsTarget(w *widget.Widget) bool {
	for _, t := range d.targets {
		if t == w {
			return true
		}
	}

	p := w.Parent()
	if p == nil {
		return false
	}

	return d.IsTarget(p)
}
