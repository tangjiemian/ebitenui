package widget

import (
	"image/color"
	"testing"

	"github.com/blizzy78/ebitenui/event"
	"github.com/matryer/is"
)

func TestSelectComboButton_SetSelectedEntry(t *testing.T) {
	is := is.New(t)

	var eventArgs *SelectComboButtonEntrySelectedEventArgs
	numEvents := 0

	b := newSelectComboButton(t,
		SelectComboButtonOpts.EntryLabelFunc(func(e interface{}) string {
			return "label " + e.(string)
		}),

		SelectComboButtonOpts.EntrySelectedHandler(func(args *SelectComboButtonEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	entry := "foo"
	b.SetSelectedEntry(entry)
	event.ExecuteDeferredActions()

	is.Equal(b.SelectedEntry(), entry)
	is.Equal(eventArgs.Entry, entry)
	is.Equal(b.Label(), "label foo")

	b.SetSelectedEntry(entry)
	event.ExecuteDeferredActions()

	is.Equal(numEvents, 1)

	entry2 := "bar"
	b.SetSelectedEntry(entry2)
	event.ExecuteDeferredActions()

	is.Equal(eventArgs.PreviousEntry, entry)
}

func TestSelectComboButton_ContentVisible_Click(t *testing.T) {
	is := is.New(t)

	b := newSelectComboButton(t)

	leftMouseButtonClick(b, t)
	is.True(b.ContentVisible())

	leftMouseButtonClick(b, t)
	is.True(!b.ContentVisible())
}

func TestSelectComboButton_ContentVisible_Programmatic(t *testing.T) {
	is := is.New(t)

	b := newSelectComboButton(t)

	b.SetContentVisible(true)
	event.ExecuteDeferredActions()

	is.True(b.ContentVisible())

	b.SetContentVisible(false)
	event.ExecuteDeferredActions()

	is.True(!b.ContentVisible())
}

func newSelectComboButton(t *testing.T, opts ...SelectComboButtonOpt) *SelectComboButton {
	t.Helper()

	b := NewSelectComboButton(append(opts,
		SelectComboButtonOpts.ComboButtonOpts(
			ComboButtonOpts.ButtonOpts(
				ButtonOpts.Image(&ButtonImage{
					Idle: newNineSliceEmpty(t),
				}),
				ButtonOpts.TextAndImage("", loadFont(t), &ButtonImageImage{
					Idle:     newImageEmpty(t),
					Disabled: newImageEmpty(t),
				}, &ButtonTextColor{
					Idle:     color.Transparent,
					Disabled: color.Transparent,
				}),
			),
			ComboButtonOpts.Content(newButton(t))),
	)...)

	event.ExecuteDeferredActions()
	render(b, t)
	return b
}
