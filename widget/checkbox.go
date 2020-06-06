package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"
	"github.com/hajimehoshi/ebiten"
)

type Checkbox struct {
	ChangedEvent *event.Event

	buttonOpts []ButtonOpt
	image      *CheckboxGraphicImage
	triState   bool

	init   *MultiOnce
	button *Button
	state  CheckboxState
}

type CheckboxOpt func(c *Checkbox)

type CheckboxGraphicImage struct {
	Unchecked *ButtonImageImage
	Checked   *ButtonImageImage
	Greyed    *ButtonImageImage
}

type CheckboxState int

type CheckboxChangedEventArgs struct {
	Checkbox *Checkbox
	State    CheckboxState
}

type CheckboxChangedHandlerFunc func(args *CheckboxChangedEventArgs)

const (
	CheckboxUnchecked = CheckboxState(iota)
	CheckboxChecked
	CheckboxGreyed
)

const CheckboxOpts = checkboxOpts(true)

type checkboxOpts bool

func NewCheckbox(opts ...CheckboxOpt) *Checkbox {
	c := &Checkbox{
		ChangedEvent: &event.Event{},

		init: &MultiOnce{},
	}

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	return c
}

func (o checkboxOpts) WithButtonOpts(opts ...ButtonOpt) CheckboxOpt {
	return func(c *Checkbox) {
		c.buttonOpts = append(c.buttonOpts, opts...)
	}
}

func (o checkboxOpts) WithImage(i *CheckboxGraphicImage) CheckboxOpt {
	return func(c *Checkbox) {
		c.image = i
	}
}

func (o checkboxOpts) WithTriState() CheckboxOpt {
	return func(c *Checkbox) {
		c.triState = true
	}
}

func (o checkboxOpts) WithChangedHandler(f CheckboxChangedHandlerFunc) CheckboxOpt {
	return func(c *Checkbox) {
		c.ChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*CheckboxChangedEventArgs))
		})
	}
}

func (c *Checkbox) GetWidget() *Widget {
	c.init.Do()
	return c.button.GetWidget()
}

func (c *Checkbox) PreferredSize() (int, int) {
	c.init.Do()
	return c.button.PreferredSize()
}

func (c *Checkbox) SetLocation(rect image.Rectangle) {
	c.init.Do()
	c.button.SetLocation(rect)
}

func (c *Checkbox) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	c.init.Do()
	c.button.SetupInputLayer(def)
}

func (c *Checkbox) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	c.init.Do()

	c.button.GraphicImage = c.state.graphicImage(c.image)

	c.button.Render(screen, def)
}

func (c *Checkbox) createWidget() {
	c.button = NewButton(
		append(c.buttonOpts, []ButtonOpt{
			ButtonOpts.WithGraphic(c.image.Unchecked.Idle),

			ButtonOpts.WithClickedHandler(func(args *ButtonClickedEventArgs) {
				c.SetState(c.state.Advance(c.triState))
			}),
		}...)...)
	c.buttonOpts = nil
}

func (c *Checkbox) State() CheckboxState {
	return c.state
}

func (c *Checkbox) SetState(s CheckboxState) {
	if s == CheckboxGreyed && !c.triState {
		panic("non-tri state checkbox cannot be in greyed state")
	}

	if s != c.state {
		c.state = s

		c.ChangedEvent.Fire(&CheckboxChangedEventArgs{
			Checkbox: c,
			State:    s,
		})
	}
}

func (s CheckboxState) Advance(triState bool) CheckboxState {
	if s == CheckboxUnchecked {
		return CheckboxChecked
	}

	if s == CheckboxChecked {
		if triState {
			return CheckboxGreyed
		}

		return CheckboxUnchecked
	}

	return CheckboxUnchecked
}

func (s CheckboxState) graphicImage(i *CheckboxGraphicImage) *ButtonImageImage {
	if s == CheckboxChecked {
		return i.Checked
	}

	if s == CheckboxGreyed {
		return i.Greyed
	}

	return i.Unchecked
}
