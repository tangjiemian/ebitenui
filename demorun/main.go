package main

import (
	"fmt"
	"github.com/blizzy78/ebitenui/demorun/gui"
	"log"
	"sort"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	goimage "image"
	_ "image/png"

	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
)

type Game struct {
	ui *ebitenui.UI
}

type pageContainer struct {
	widget    widget.PreferredSizeLocateableWidget
	titleText *widget.Text
	flipBook  *widget.FlipBook
}

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Ebiten UI Demo")
	ebiten.SetWindowResizable(true)
	ebiten.SetScreenClearedEveryFrame(false)

	ui, closeUI, err := createUI()
	if err != nil {
		log.Fatal(err)
	}

	defer closeUI()

	game := &Game{
		ui: ui,
	}

	err = ebiten.RunGame(game)
	if err != nil {
		log.Print(err)
	}
}

func createUI() (*ebitenui.UI, func(), error) {
	res, err := gui.NewUIResources()
	if err != nil {
		return nil, nil, err
	}

	drag := &gui.DragContents{
		Res: res,
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(res.Background))

	toolTips := gui.ToolTipContents{
		Tips: map[widget.HasWidget]string{},
		Res:  res,
	}

	toolTip := widget.NewToolTip(
		widget.ToolTipOpts.Container(rootContainer),
		widget.ToolTipOpts.ContentsCreater(&toolTips),
	)

	dnd := widget.NewDragAndDrop(
		widget.DragAndDropOpts.Container(rootContainer),
		widget.DragAndDropOpts.ContentsCreater(drag),
	)

	rootContainer.AddChild(headerContainer(res))

	var ui *ebitenui.UI
	rootContainer.AddChild(demoContainer(res, &toolTips, toolTip, dnd, drag, func() *ebitenui.UI {
		return ui
	}))

	urlContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Padding(widget.Insets{
			Left:  25,
			Right: 25,
		}),
	)))
	rootContainer.AddChild(urlContainer)

	urlContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("github.com/blizzy78/ebitenui", res.Text.SmallFace, res.Text.DisabledColor)))

	ui = &ebitenui.UI{
		Container: rootContainer,

		ToolTip: toolTip,

		DragAndDrop: dnd,
	}

	return ui, func() {
		res.Close()
	}, nil
}

func headerContainer(res *gui.UiResources) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(15))),
	)

	c.AddChild(header("Ebiten UI Demo", res,
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
	))

	c2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
		)),
	)
	c.AddChild(c2)

	c2.AddChild(widget.NewText(
		widget.TextOpts.Text("This program is a showcase of Ebiten UI widgets and layouts.", res.Text.Face, res.Text.IdleColor)))

	return c
}

func header(label string, res *gui.UiResources, opts ...widget.ContainerOpt) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(append(opts, []widget.ContainerOpt{
		widget.ContainerOpts.BackgroundImage(res.Header.Background),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.Header.Padding))),
	}...)...)

	c.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text(label, res.Header.Face, res.Header.Color),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	return c
}

func demoContainer(res *gui.UiResources, toolTips *gui.ToolTipContents, toolTip *widget.ToolTip, dnd *widget.DragAndDrop, drag *gui.DragContents,
	ui func() *ebitenui.UI) widget.PreferredSizeLocateableWidget {

	demoContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(20, 0),
		)))

	pages := []interface{}{
		buttonPage(res),
		checkboxPage(res),
		listPage(res),
		comboButtonPage(res),
		tabBookPage(res),
		gridLayoutPage(res),
		rowLayoutPage(res),
		sliderPage(res),
		toolTipPage(res, toolTips, toolTip),
		dragAndDropPage(res, dnd, drag),
		textInputPage(res),
		radioGroupPage(res),
		windowPage(res, ui),
		anchorLayoutPage(res),
	}

	collator := collate.New(language.English)
	sort.Slice(pages, func(a int, b int) bool {
		p1 := pages[a].(*page)
		p2 := pages[b].(*page)
		return collator.CompareString(p1.title, p2.title) < 0
	})

	pageContainer := newPageContainer(res)

	pageList := widget.NewList(
		widget.ListOpts.Entries(pages),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(*page).title
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.Image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(res.List.Track, res.List.Handle),
			widget.SliderOpts.HandleSize(res.List.HandleSize),
			widget.SliderOpts.TrackPadding(res.List.TrackPadding),
		),
		widget.ListOpts.EntryColor(res.List.Entry),
		widget.ListOpts.EntryFontFace(res.List.Face),
		widget.ListOpts.EntryTextPadding(res.List.EntryPadding),
		widget.ListOpts.HideHorizontalSlider(),

		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			pageContainer.setPage(args.Entry.(*page))
		}))
	demoContainer.AddChild(pageList)

	demoContainer.AddChild(pageContainer.widget)

	pageList.SetSelectedEntry(pages[0])

	return demoContainer
}

func newPageContainer(res *gui.UiResources) *pageContainer {
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.Panel.Padding),
			widget.RowLayoutOpts.Spacing(15))),
	)

	titleText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("", res.Text.TitleFace, res.Text.IdleColor))
	c.AddChild(titleText)

	flipBook := widget.NewFlipBook(
		widget.FlipBookOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		}))),
	)
	c.AddChild(flipBook)

	return &pageContainer{
		widget:    c,
		titleText: titleText,
		flipBook:  flipBook,
	}
}

func (p *pageContainer) setPage(page *page) {
	p.titleText.Label = page.title
	p.flipBook.SetPage(page.content)
	p.flipBook.RequestRelayout()
}

func newCheckbox(label string, changedHandler widget.CheckboxChangedHandlerFunc, res *gui.UiResources) *widget.LabeledCheckbox {
	return widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(res.Checkbox.Spacing),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(res.Checkbox.Image)),
			widget.CheckboxOpts.Image(res.Checkbox.Graphic),
			widget.CheckboxOpts.ChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				if changedHandler != nil {
					changedHandler(args)
				}
			})),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text(label, res.Label.Face, res.Label.Text)))
}

func newPageContentContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
}

func newListComboButton(entries []interface{}, buttonLabel widget.SelectComboButtonEntryLabelFunc, entryLabel widget.ListEntryLabelFunc,
	entrySelectedHandler widget.ListComboButtonEntrySelectedHandlerFunc, res *gui.UiResources) *widget.ListComboButton {

	return widget.NewListComboButton(
		widget.ListComboButtonOpts.SelectComboButtonOpts(
			widget.SelectComboButtonOpts.ComboButtonOpts(
				widget.ComboButtonOpts.ButtonOpts(
					widget.ButtonOpts.Image(res.ComboButton.Image),
					widget.ButtonOpts.TextPadding(res.ComboButton.Padding),
				),
			),
		),
		widget.ListComboButtonOpts.Text(res.ComboButton.Face, res.ComboButton.Graphic, res.ComboButton.Text),
		widget.ListComboButtonOpts.ListOpts(
			widget.ListOpts.Entries(entries),
			widget.ListOpts.ScrollContainerOpts(
				widget.ScrollContainerOpts.Image(res.List.Image),
			),
			widget.ListOpts.SliderOpts(
				widget.SliderOpts.Images(res.List.Track, res.List.Handle),
				widget.SliderOpts.HandleSize(res.List.HandleSize),
				widget.SliderOpts.TrackPadding(res.List.TrackPadding)),
			widget.ListOpts.EntryFontFace(res.List.Face),
			widget.ListOpts.EntryColor(res.List.Entry),
			widget.ListOpts.EntryTextPadding(res.List.EntryPadding),
		),
		widget.ListComboButtonOpts.EntryLabelFunc(buttonLabel, entryLabel),
		widget.ListComboButtonOpts.EntrySelectedHandler(entrySelectedHandler))
}

func newList(entries []interface{}, res *gui.UiResources, widgetOpts ...widget.WidgetOpt) *widget.List {
	return widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.Image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(res.List.Track, res.List.Handle),
			widget.SliderOpts.HandleSize(res.List.HandleSize),
			widget.SliderOpts.TrackPadding(res.List.TrackPadding),
		),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.EntryFontFace(res.List.Face),
		widget.ListOpts.EntryColor(res.List.Entry),
		widget.ListOpts.EntryTextPadding(res.List.EntryPadding),
	)
}

func newSeparator(res *gui.UiResources, ld interface{}) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(res.SeparatorColor)),
	))

	return c
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

//-------------------------

type page struct {
	title   string
	content widget.PreferredSizeLocateableWidget
}

func buttonPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	bs := []*widget.Button{}
	for i := 0; i < 3; i++ {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.Button.Image),
			widget.ButtonOpts.Text(fmt.Sprintf("Button %d", i+1), res.Button.Face, res.Button.Text),
			widget.ButtonOpts.TextPadding(res.Button.Padding),
		)
		c.AddChild(b)
		bs = append(bs, b)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		for _, b := range bs {
			b.GetWidget().Disabled = args.State == widget.CheckboxChecked
		}
	}, res))

	return &page{
		title:   "Button",
		content: c,
	}
}

func checkboxPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	cb1 := newCheckbox("Two-State Checkbox", nil, res)
	c.AddChild(cb1)

	cb2 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(res.Checkbox.Spacing),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(res.Checkbox.Image)),
			widget.CheckboxOpts.Image(res.Checkbox.Graphic),
			widget.CheckboxOpts.TriState()),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Tri-State Checkbox", res.Label.Face, res.Label.Text)))
	c.AddChild(cb2)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		cb2.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Checkbox",
		content: c,
	}
}

func listPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	listsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(10, 0))))
	c.AddChild(listsContainer)

	entries1 := []interface{}{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"}
	list1 := newList(entries1, res, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list1)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
	listsContainer.AddChild(buttonsContainer)

	bs := []*widget.Button{}
	for i := 0; i < 3; i++ {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.Button.Image),
			widget.ButtonOpts.TextPadding(res.Button.Padding),
			widget.ButtonOpts.Text(fmt.Sprintf("Action %d", i+1), res.Button.Face, res.Button.Text))
		buttonsContainer.AddChild(b)
		bs = append(bs, b)
	}

	entries2 := []interface{}{"Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen", "Twenty"}
	list2 := newList(entries2, res, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 220,
	}))
	listsContainer.AddChild(list2)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		list1.GetWidget().Disabled = args.State == widget.CheckboxChecked
		list2.GetWidget().Disabled = args.State == widget.CheckboxChecked
		for _, b := range bs {
			b.GetWidget().Disabled = args.State == widget.CheckboxChecked
		}
	}, res))

	return &page{
		title:   "List",
		content: c,
	}
}

func comboButtonPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	entries := []interface{}{}
	for i := 1; i <= 20; i++ {
		entries = append(entries, i)
	}

	cb := newListComboButton(
		entries,
		func(e interface{}) string {
			return fmt.Sprintf("Entry %d", e.(int))
		},
		func(e interface{}) string {
			return fmt.Sprintf("Entry %d", e.(int))
		},
		func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			c.RequestRelayout()
		},
		res)
	c.AddChild(cb)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		cb.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Combo Button",
		content: c,
	}
}

func tabBookPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	tabs := []*widget.TabBookTab{}

	for i := 0; i < 4; i++ {
		tc := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10))),
			widget.ContainerOpts.AutoDisableChildren())

		for j := 0; j < 3; j++ {
			b := widget.NewButton(
				widget.ButtonOpts.Image(res.Button.Image),
				widget.ButtonOpts.TextPadding(res.Button.Padding),
				widget.ButtonOpts.Text(fmt.Sprintf("Button %d on Tab %d", j+1, i+1), res.Button.Face, res.Button.Text))
			tc.AddChild(b)
		}

		tab := widget.NewTabBookTab(fmt.Sprintf("Tab %d", i+1), tc)
		if i == 2 {
			tab.Disabled = true
		}

		tabs = append(tabs, tab)
	}

	t := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tabs...),
		widget.TabBookOpts.TabButtonImage(res.TabBook.IdleButton, res.TabBook.SelectedButton),
		widget.TabBookOpts.TabButtonText(res.TabBook.ButtonFace, res.TabBook.ButtonText),
		widget.TabBookOpts.TabButtonOpts(widget.StateButtonOpts.ButtonOpts(widget.ButtonOpts.TextPadding(res.TabBook.ButtonPadding))),
		widget.TabBookOpts.TabButtonSpacing(10),
		widget.TabBookOpts.Spacing(15))
	c.AddChild(t)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Tab Book",
		content: c,
	}
}

func gridLayoutPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	bc := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(4),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(10, 10))))
	c.AddChild(bc)

	i := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 4; col++ {
			b := widget.NewButton(
				widget.ButtonOpts.Image(res.Button.Image),
				widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+i)), i+1), res.Button.Face, res.Button.Text))
			bc.AddChild(b)

			i++
		}
	}

	return &page{
		title:   "Grid Layout",
		content: c,
	}
}

func rowLayoutPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Horizontal", res.Text.Face, res.Text.IdleColor)))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(5))))
	c.AddChild(bc)

	for col := 0; col < 5; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.Image(res.Button.Image),
			widget.ButtonOpts.TextPadding(res.Button.Padding),
			widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+col)), col+1), res.Button.Face, res.Button.Text))
		bc.AddChild(b)
	}

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Vertical", res.Text.Face, res.Text.IdleColor)))

	bc = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5))))
	c.AddChild(bc)

	labels := []string{"Tiny", "Medium", "Very Large"}
	for _, l := range labels {
		b := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.Button.Image),
			widget.ButtonOpts.TextPadding(res.Button.Padding),
			widget.ButtonOpts.Text(l, res.Button.Face, res.Button.Text))
		bc.AddChild(b)
	}

	return &page{
		title:   "Row Layout",
		content: c,
	}
}

func sliderPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	pageSizes := []int{3, 10}
	sliders := []*widget.Slider{}

	for _, ps := range pageSizes {
		ps := ps

		sc := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(10))),
			widget.ContainerOpts.AutoDisableChildren(),
		)
		c.AddChild(sc)

		var text *widget.Label

		s := widget.NewSlider(
			widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),
			widget.SliderOpts.MinMax(1, 20),
			widget.SliderOpts.Images(res.Slider.TrackImage, res.Slider.Handle),
			widget.SliderOpts.HandleSize(res.Slider.HandleSize),
			widget.SliderOpts.PageSizeFunc(func() int {
				return ps
			}),
			widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
				text.Label = fmt.Sprintf("%d", args.Current)
			}),
		)
		sc.AddChild(s)
		sliders = append(sliders, s)

		text = widget.NewLabel(
			widget.LabelOpts.TextOpts(widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}))),
			widget.LabelOpts.Text(fmt.Sprintf("%d", s.Current), res.Label.Face, res.Label.Text),
		)
		sc.AddChild(text)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		for _, s := range sliders {
			s.GetWidget().Parent().Disabled = args.State == widget.CheckboxChecked
		}
	}, res))

	return &page{
		title:   "Slider",
		content: c,
	}
}

func toolTipPage(res *gui.UiResources, toolTips *gui.ToolTipContents, toolTip *widget.ToolTip) *page {
	c := newPageContentContainer()

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Hover over these buttons to see their tool tips.", res.Text.Face, res.Text.IdleColor)))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15))))
	c.AddChild(bc)

	for col := 0; col < 4; col++ {
		b := widget.NewButton(
			widget.ButtonOpts.Image(res.Button.Image),
			widget.ButtonOpts.TextPadding(res.Button.Padding),
			widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+col)), col+1), res.Button.Face, res.Button.Text))

		if col == 2 {
			b.GetWidget().Disabled = true
		}

		toolTips.Set(b, fmt.Sprintf("Tool tip for button %d", col+1))
		toolTips.WidgetsWithTime = append(toolTips.WidgetsWithTime, b)

		bc.AddChild(b)
	}

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	showTimeCheckbox := newCheckbox("Show additional infos in tool tips", func(args *widget.CheckboxChangedEventArgs) {
		toolTips.ShowTime = args.State == widget.CheckboxChecked
	}, res)
	toolTips.Set(showTimeCheckbox, "If enabled, tool tips will show system time for demonstration.")
	c.AddChild(showTimeCheckbox)

	stickyDelayedCheckbox := newCheckbox("Tool tips are sticky and delayed", func(args *widget.CheckboxChangedEventArgs) {
		toolTip.Sticky = args.State == widget.CheckboxChecked
		if args.State == widget.CheckboxChecked {
			toolTip.Delay = 800 * time.Millisecond
		} else {
			toolTip.Delay = 0
		}
	}, res)
	toolTips.Set(stickyDelayedCheckbox, "If enabled, tool tips do not show immediately and will not move with the cursor.")
	c.AddChild(stickyDelayedCheckbox)

	return &page{
		title:   "Tool Tip",
		content: c,
	}
}

func dragAndDropPage(res *gui.UiResources, dnd *widget.DragAndDrop, drag *gui.DragContents) *page {
	c := newPageContentContainer()

	dndContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(30),
		)),
	)
	c.AddChild(dndContainer)

	sourcePanel := gui.NewSizedPanel(200, 200,
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.Panel.Padding))),
	)
	drag.AddSource(sourcePanel)
	dndContainer.AddChild(sourcePanel)

	sourcePanel.Container().AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text("Drag\nFrom\nHere", res.Text.Face, res.Text.DisabledColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	))

	targetPanel := gui.NewSizedPanel(200, 200,
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(res.Panel.Padding))),
	)
	drag.AddTarget(targetPanel)
	dndContainer.AddChild(targetPanel)

	targetText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text("Drop\nHere", res.Text.Face, res.Text.DisabledColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)

	targetPanel.Container().AddChild(targetText)

	dnd.DroppedEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.DragAndDropDroppedEventArgs)
		if !drag.IsTarget(a.Target.GetWidget()) {
			return
		}

		targetText.Label = "Thanks!"
		targetText.Color = res.Text.IdleColor

		time.AfterFunc(2500*time.Millisecond, func() {
			targetText.Label = "Drop\nHere"
			targetText.Color = res.Text.DisabledColor
		})
	})

	return &page{
		title:   "Drag & Drop",
		content: c,
	}
}

func textInputPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(res.TextInput.Image),
		widget.TextInputOpts.Color(res.TextInput.Color),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(res.TextInput.Face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.TextInput.Face, 2),
		),
	}

	t := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter text here"))...,
	)
	c.AddChild(t)

	tSecure := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter secure text here"),
		widget.TextInputOpts.Secure(true))...,
	)
	c.AddChild(tSecure)

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	c.AddChild(newCheckbox("Disabled", func(args *widget.CheckboxChangedEventArgs) {
		t.GetWidget().Disabled = args.State == widget.CheckboxChecked
		tSecure.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res))

	return &page{
		title:   "Text Input",
		content: c,
	}
}

func radioGroupPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	var cbs []*widget.Checkbox
	for i := 0; i < 5; i++ {
		cb := newCheckbox(fmt.Sprintf("Checkbox %d", i+1), nil, res)
		c.AddChild(cb)
		cbs = append(cbs, cb.Checkbox())
	}

	widget.NewRadioGroup(widget.RadioGroupOpts.Checkboxes(cbs...))

	return &page{
		title:   "Radio Group",
		content: c,
	}
}

func windowPage(res *gui.UiResources, ui func() *ebitenui.UI) *page {
	c := newPageContentContainer()

	b := widget.NewButton(
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.TextPadding(res.Button.Padding),
		widget.ButtonOpts.Text("Open Window", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openWindow(res, ui)
		}),
	)
	c.AddChild(b)

	return &page{
		title:   "Window",
		content: c,
	}
}

func openWindow(res *gui.UiResources, ui func() *ebitenui.UI) {
	var rw ebitenui.RemoveWindowFunc

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.Panel.Padding),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Modal Window", res.Text.BigTitleFace, res.Text.IdleColor),
	))

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("This window blocks all input to widgets below it.", res.Text.Face, res.Text.IdleColor),
	))

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	c.AddChild(bc)

	o2b := widget.NewButton(
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.TextPadding(res.Button.Padding),
		widget.ButtonOpts.Text("Open Another", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openWindow2(res, ui)
		}),
	)
	bc.AddChild(o2b)

	cb := widget.NewButton(
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.TextPadding(res.Button.Padding),
		widget.ButtonOpts.Text("Close", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	bc.AddChild(cb)

	w := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
	)

	ww, wh := ebiten.WindowSize()
	r := goimage.Rect(0, 0, ww*3/4, wh/3)
	r = r.Add(goimage.Point{ww / 4 / 2, wh * 2 / 3 / 2})
	w.SetLocation(r)

	rw = ui().AddWindow(w)
}

func openWindow2(res *gui.UiResources, ui func() *ebitenui.UI) {
	var rw ebitenui.RemoveWindowFunc

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.Panel.Padding),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Second Window", res.Text.BigTitleFace, res.Text.IdleColor),
	))

	cb := widget.NewButton(
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.TextPadding(res.Button.Padding),
		widget.ButtonOpts.Text("Close", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	c.AddChild(cb)

	w := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
	)

	ww, wh := ebiten.WindowSize()
	r := goimage.Rect(0, 0, ww/2, wh/2)
	r = r.Add(goimage.Point{ww * 4 / 10, wh / 2 / 2})
	w.SetLocation(r)

	rw = ui().AddWindow(w)
}

func anchorLayoutPage(res *gui.UiResources) *page {
	c := newPageContentContainer()

	p := gui.NewSizedPanel(300, 220,
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(res.Panel.Padding),
		)),
	)
	c.AddChild(p)

	sp := gui.NewSizedPanel(50, 50,
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{})),
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
	)
	p.Container().AddChild(sp.Container())

	c.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	posC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(50),
		)),
	)
	c.AddChild(posC)

	hPosC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	posC.AddChild(hPosC)

	hPosC.AddChild(widget.NewLabel(widget.LabelOpts.Text("Horizontal", res.Label.Face, res.Label.Text)))

	labels := []string{"Start", "Center", "End"}
	hCBs := []*widget.Checkbox{}
	for _, l := range labels {
		cb := newCheckbox(l, nil, res)
		hPosC.AddChild(cb)
		hCBs = append(hCBs, cb.Checkbox())
	}

	widget.NewRadioGroup(
		widget.RadioGroupOpts.Checkboxes(hCBs...),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
			ald.HorizontalPosition = widget.AnchorLayoutPosition(indexCheckbox(hCBs, args.Active))
			sp.Container().GetWidget().LayoutData = ald
			p.Container().RequestRelayout()
		}),
	)

	vPosC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.AutoDisableChildren(),
	)
	posC.AddChild(vPosC)

	vPosC.AddChild(widget.NewLabel(widget.LabelOpts.Text("Vertical", res.Label.Face, res.Label.Text)))

	vCBs := []*widget.Checkbox{}
	for _, l := range labels {
		cb := newCheckbox(l, nil, res)
		vPosC.AddChild(cb)
		vCBs = append(vCBs, cb.Checkbox())
	}

	widget.NewRadioGroup(
		widget.RadioGroupOpts.Checkboxes(vCBs...),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
			ald.VerticalPosition = widget.AnchorLayoutPosition(indexCheckbox(vCBs, args.Active))
			sp.Container().GetWidget().LayoutData = ald
			p.Container().RequestRelayout()
		}),
	)

	stretchC := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	posC.AddChild(stretchC)

	stretchC.AddChild(widget.NewText(widget.TextOpts.Text("Stretch", res.Text.Face, res.Text.IdleColor)))

	stretchHorizontalCheckbox := newCheckbox("Horizontal", func(args *widget.CheckboxChangedEventArgs) {
		ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
		ald.StretchHorizontal = args.State == widget.CheckboxChecked
		sp.Container().GetWidget().LayoutData = ald
		p.Container().RequestRelayout()

		hPosC.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res)
	stretchC.AddChild(stretchHorizontalCheckbox)

	stretchVerticalCheckbox := newCheckbox("Vertical", func(args *widget.CheckboxChangedEventArgs) {
		ald := sp.Container().GetWidget().LayoutData.(widget.AnchorLayoutData)
		ald.StretchVertical = args.State == widget.CheckboxChecked
		sp.Container().GetWidget().LayoutData = ald
		p.Container().RequestRelayout()

		vPosC.GetWidget().Disabled = args.State == widget.CheckboxChecked
	}, res)
	stretchC.AddChild(stretchVerticalCheckbox)

	return &page{
		title:   "Anchor Layout",
		content: c,
	}
}

func indexCheckbox(cs []*widget.Checkbox, c *widget.Checkbox) int {
	for i, cb := range cs {
		if cb == c {
			return i
		}
	}
	return -1
}
