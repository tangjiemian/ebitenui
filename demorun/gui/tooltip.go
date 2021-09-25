package gui

import (
	"time"

	"github.com/blizzy78/ebitenui/widget"
)

type ToolTipContents struct {
	Tips            map[widget.HasWidget]string
	WidgetsWithTime []widget.HasWidget
	ShowTime        bool

	Res *UiResources

	text     *widget.TextToolTip
	timeText *widget.TextToolTip
}

func (t *ToolTipContents) Create(w widget.HasWidget) widget.ToolTipWidget {
	if _, ok := t.Tips[w]; !ok {
		return nil
	}

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(t.Res.ToolTip.background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(t.Res.ToolTip.padding),
			widget.RowLayoutOpts.Spacing(2),
		)))

	t.text = widget.NewTextToolTip(
		widget.TextToolTipOpts.TextOpts(
			widget.TextOpts.Text("", t.Res.ToolTip.face, t.Res.ToolTip.color),
		),
	)
	c.AddChild(t.text)

	if t.ShowTime && t.CanShowTime(w) {
		t.timeText = widget.NewTextToolTip(
			widget.TextToolTipOpts.TextOpts(
				widget.TextOpts.Text("", t.Res.ToolTip.face, t.Res.ToolTip.color),
			),
		)
		c.AddChild(t.timeText)
	}

	return c
}

func (t *ToolTipContents) Set(w widget.HasWidget, s string) {
	t.Tips[w] = s
}

func (t *ToolTipContents) Update(w widget.HasWidget) {
	t.text.Label = t.Tips[w]

	if !t.ShowTime || !t.CanShowTime(w) {
		return
	}

	t.timeText.Label = time.Now().Local().Format("2006-01-02 15:04:05")
}

func (t *ToolTipContents) CanShowTime(w widget.HasWidget) bool {
	for _, tw := range t.WidgetsWithTime {
		if tw == w {
			return true
		}
	}
	return false
}
