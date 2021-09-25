package gui

import (
	"image"

	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type SizedPanel struct {
	width     int
	height    int
	container *widget.Container
}

func NewSizedPanel(w int, h int, opts ...widget.ContainerOpt) *SizedPanel {
	return &SizedPanel{
		width:     w,
		height:    h,
		container: widget.NewContainer(opts...),
	}
}

func (p *SizedPanel) GetWidget() *widget.Widget {
	return p.container.GetWidget()
}

func (p *SizedPanel) PreferredSize() (int, int) {
	return p.width, p.height
}

func (p *SizedPanel) SetLocation(rect image.Rectangle) {
	p.container.SetLocation(rect)
}

func (p *SizedPanel) Render(screen *ebiten.Image, def widget.DeferredRenderFunc) {
	p.container.Render(screen, def)
}

func (p *SizedPanel) Container() *widget.Container {
	return p.container
}
