package gui

import (
	"image/color"
	"strconv"

	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
)

const (
	backgroundColor = "131a22"

	textIdleColor     = "dff4ff"
	textDisabledColor = "5a7a91"

	labelIdleColor     = textIdleColor
	labelDisabledColor = textDisabledColor

	buttonIdleColor     = textIdleColor
	buttonDisabledColor = labelDisabledColor

	listSelectedBackground         = "4b687a"
	listDisabledSelectedBackground = "2a3944"

	headerColor = textIdleColor

	textInputCaretColor         = "e7c34b"
	textInputDisabledCaretColor = "766326"

	toolTipColor = backgroundColor

	separatorColor = listDisabledSelectedBackground
)

type UiResources struct {
	Fonts *Fonts

	Background *image.NineSlice

	SeparatorColor color.Color

	Text        *TextResources
	Button      *buttonResources
	Label       *labelResources
	Checkbox    *checkboxResources
	ComboButton *comboButtonResources
	List        *listResources
	Slider      *sliderResources
	Panel       *panelResources
	TabBook     *tabBookResources
	Header      *headerResources
	TextInput   *textInputResources
	ToolTip     *toolTipResources
}

type TextResources struct {
	IdleColor     color.Color
	DisabledColor color.Color
	Face          font.Face
	TitleFace     font.Face
	BigTitleFace  font.Face
	SmallFace     font.Face
}

type buttonResources struct {
	Image   *widget.ButtonImage
	Text    *widget.ButtonTextColor
	Face    font.Face
	Padding widget.Insets
}

type checkboxResources struct {
	Image   *widget.ButtonImage
	Graphic *widget.CheckboxGraphicImage
	Spacing int
}

type labelResources struct {
	Text *widget.LabelColor
	Face font.Face
}

type comboButtonResources struct {
	Image   *widget.ButtonImage
	Text    *widget.ButtonTextColor
	Face    font.Face
	Graphic *widget.ButtonImageImage
	Padding widget.Insets
}

type listResources struct {
	Image        *widget.ScrollContainerImage
	Track        *widget.SliderTrackImage
	TrackPadding widget.Insets
	Handle       *widget.ButtonImage
	HandleSize   int
	Face         font.Face
	Entry        *widget.ListEntryColor
	EntryPadding widget.Insets
}

type sliderResources struct {
	TrackImage *widget.SliderTrackImage
	Handle     *widget.ButtonImage
	HandleSize int
}

type panelResources struct {
	Image   *image.NineSlice
	Padding widget.Insets
}

type tabBookResources struct {
	IdleButton     *widget.ButtonImage
	SelectedButton *widget.ButtonImage
	ButtonFace     font.Face
	ButtonText     *widget.ButtonTextColor
	ButtonPadding  widget.Insets
}

type headerResources struct {
	Background *image.NineSlice
	Padding    widget.Insets
	Face       font.Face
	Color      color.Color
}

type textInputResources struct {
	Image   *widget.TextInputImage
	Padding widget.Insets
	Face    font.Face
	Color   *widget.TextInputColor
}

type toolTipResources struct {
	background *image.NineSlice
	padding    widget.Insets
	face       font.Face
	color      color.Color
}

func NewUIResources() (*UiResources, error) {
	background := image.NewNineSliceColor(hexToColor(backgroundColor))

	fonts, err := LoadFonts()
	if err != nil {
		return nil, err
	}

	button, err := newButtonResources(fonts)
	if err != nil {
		return nil, err
	}

	checkbox, err := newCheckboxResources()
	if err != nil {
		return nil, err
	}

	comboButton, err := newComboButtonResources(fonts)
	if err != nil {
		return nil, err
	}

	list, err := newListResources(fonts)
	if err != nil {
		return nil, err
	}

	slider, err := newSliderResources()
	if err != nil {
		return nil, err
	}

	panel, err := newPanelResources()
	if err != nil {
		return nil, err
	}

	tabBook, err := newTabBookResources(fonts)
	if err != nil {
		return nil, err
	}

	header, err := newHeaderResources(fonts)
	if err != nil {
		return nil, err
	}

	textInput, err := newTextInputResources(fonts)
	if err != nil {
		return nil, err
	}

	toolTip, err := newToolTipResources(fonts)
	if err != nil {
		return nil, err
	}

	return &UiResources{
		Fonts: fonts,

		Background: background,

		SeparatorColor: hexToColor(separatorColor),

		Text: &TextResources{
			IdleColor:     hexToColor(textIdleColor),
			DisabledColor: hexToColor(textDisabledColor),
			Face:          fonts.face,
			TitleFace:     fonts.titleFace,
			BigTitleFace:  fonts.bigTitleFace,
			SmallFace:     fonts.toolTipFace,
		},

		Button:      button,
		Label:       newLabelResources(fonts),
		Checkbox:    checkbox,
		ComboButton: comboButton,
		List:        list,
		Slider:      slider,
		Panel:       panel,
		TabBook:     tabBook,
		Header:      header,
		TextInput:   textInput,
		ToolTip:     toolTip,
	}, nil
}

func newButtonResources(fonts *Fonts) (*buttonResources, error) {
	idle, err := LoadImageNineSlice("demorun/graphics/button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := LoadImageNineSlice("demorun/graphics/button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}

	pressed, err := LoadImageNineSlice("demorun/graphics/button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := LoadImageNineSlice("demorun/graphics/button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	i := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}

	return &buttonResources{
		Image: i,

		Text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		Face: fonts.face,

		Padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newCheckboxResources() (*checkboxResources, error) {
	idle, err := LoadImageNineSlice("demorun/graphics/checkbox-idle.png", 20, 0)
	if err != nil {
		return nil, err
	}

	hover, err := LoadImageNineSlice("demorun/graphics/checkbox-hover.png", 20, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := LoadImageNineSlice("demorun/graphics/checkbox-disabled.png", 20, 0)
	if err != nil {
		return nil, err
	}

	checked, err := LoadGraphicImages("demorun/graphics/checkbox-checked-idle.png", "demorun/graphics/checkbox-checked-disabled.png")
	if err != nil {
		return nil, err
	}

	unchecked, err := LoadGraphicImages("demorun/graphics/checkbox-unchecked-idle.png", "demorun/graphics/checkbox-unchecked-disabled.png")
	if err != nil {
		return nil, err
	}

	greyed, err := LoadGraphicImages("demorun/graphics/checkbox-greyed-idle.png", "demorun/graphics/checkbox-greyed-disabled.png")
	if err != nil {
		return nil, err
	}

	return &checkboxResources{
		Image: &widget.ButtonImage{
			Idle:     idle,
			Hover:    hover,
			Pressed:  hover,
			Disabled: disabled,
		},

		Graphic: &widget.CheckboxGraphicImage{
			Checked:   checked,
			Unchecked: unchecked,
			Greyed:    greyed,
		},

		Spacing: 10,
	}, nil
}

func newLabelResources(fonts *Fonts) *labelResources {
	return &labelResources{
		Text: &widget.LabelColor{
			Idle:     hexToColor(labelIdleColor),
			Disabled: hexToColor(labelDisabledColor),
		},
		Face: fonts.face,
	}
}

func newComboButtonResources(fonts *Fonts) (*comboButtonResources, error) {
	idle, err := LoadImageNineSlice("demorun/graphics/combo-button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := LoadImageNineSlice("demorun/graphics/combo-button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}

	pressed, err := LoadImageNineSlice("demorun/graphics/combo-button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := LoadImageNineSlice("demorun/graphics/combo-button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	i := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}

	arrowDown, err := LoadGraphicImages("demorun/graphics/arrow-down-idle.png", "demorun/graphics/arrow-down-disabled.png")
	if err != nil {
		return nil, err
	}

	return &comboButtonResources{
		Image: i,

		Text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		Face:    fonts.face,
		Graphic: arrowDown,

		Padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newListResources(fonts *Fonts) (*listResources, error) {
	idle, _, err := ebitenutil.NewImageFromFile("demorun/graphics/list-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, _, err := ebitenutil.NewImageFromFile("demorun/graphics/list-disabled.png")
	if err != nil {
		return nil, err
	}

	mask, _, err := ebitenutil.NewImageFromFile("demorun/graphics/list-mask.png")
	if err != nil {
		return nil, err
	}

	trackIdle, _, err := ebitenutil.NewImageFromFile("demorun/graphics/list-track-idle.png")
	if err != nil {
		return nil, err
	}

	trackDisabled, _, err := ebitenutil.NewImageFromFile("demorun/graphics/list-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, _, err := ebitenutil.NewImageFromFile("demorun/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, _, err := ebitenutil.NewImageFromFile("demorun/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	return &listResources{
		Image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(idle, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(disabled, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(mask, [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		},

		Track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(trackDisabled, [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},

		TrackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},

		Handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleIdle, 0, 5),
		},

		HandleSize: 5,
		Face:       fonts.face,

		Entry: &widget.ListEntryColor{
			Unselected:         hexToColor(textIdleColor),
			DisabledUnselected: hexToColor(textDisabledColor),

			Selected:         hexToColor(textIdleColor),
			DisabledSelected: hexToColor(textDisabledColor),

			SelectedBackground:         hexToColor(listSelectedBackground),
			DisabledSelectedBackground: hexToColor(listDisabledSelectedBackground),
		},

		EntryPadding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    2,
			Bottom: 2,
		},
	}, nil
}

func newSliderResources() (*sliderResources, error) {
	idle, _, err := ebitenutil.NewImageFromFile("demorun/graphics/slider-track-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, _, err := ebitenutil.NewImageFromFile("demorun/graphics/slider-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, _, err := ebitenutil.NewImageFromFile("demorun/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, _, err := ebitenutil.NewImageFromFile("demorun/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	handleDisabled, _, err := ebitenutil.NewImageFromFile("demorun/graphics/slider-handle-disabled.png")
	if err != nil {
		return nil, err
	}

	return &sliderResources{
		TrackImage: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(idle, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Hover:    image.NewNineSlice(idle, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Disabled: image.NewNineSlice(disabled, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
		},

		Handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleDisabled, 0, 5),
		},

		HandleSize: 6,
	}, nil
}

func newPanelResources() (*panelResources, error) {
	i, err := LoadImageNineSlice("demorun/graphics/panel-idle.png", 10, 10)
	if err != nil {
		return nil, err
	}

	return &panelResources{
		Image: i,
		Padding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		},
	}, nil
}

func newTabBookResources(fonts *Fonts) (*tabBookResources, error) {
	selectedIdle, err := LoadImageNineSlice("demorun/graphics/button-selected-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	selectedHover, err := LoadImageNineSlice("demorun/graphics/button-selected-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}

	selectedPressed, err := LoadImageNineSlice("demorun/graphics/button-selected-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	selectedDisabled, err := LoadImageNineSlice("demorun/graphics/button-selected-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	selected := &widget.ButtonImage{
		Idle:     selectedIdle,
		Hover:    selectedHover,
		Pressed:  selectedPressed,
		Disabled: selectedDisabled,
	}

	idle, err := LoadImageNineSlice("demorun/graphics/button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := LoadImageNineSlice("demorun/graphics/button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}

	pressed, err := LoadImageNineSlice("demorun/graphics/button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := LoadImageNineSlice("demorun/graphics/button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	unselected := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}

	return &tabBookResources{
		SelectedButton: selected,
		IdleButton:     unselected,
		ButtonFace:     fonts.face,

		ButtonText: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		ButtonPadding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newHeaderResources(fonts *Fonts) (*headerResources, error) {
	bg, err := LoadImageNineSlice("demorun/graphics/header.png", 446, 9)
	if err != nil {
		return nil, err
	}

	return &headerResources{
		Background: bg,

		Padding: widget.Insets{
			Left:   25,
			Right:  25,
			Top:    4,
			Bottom: 4,
		},

		Face:  fonts.bigTitleFace,
		Color: hexToColor(headerColor),
	}, nil
}

func newTextInputResources(fonts *Fonts) (*textInputResources, error) {
	idle, _, err := ebitenutil.NewImageFromFile("demorun/graphics/text-input-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, _, err := ebitenutil.NewImageFromFile("demorun/graphics/text-input-disabled.png")
	if err != nil {
		return nil, err
	}

	return &textInputResources{
		Image: &widget.TextInputImage{
			Idle:     image.NewNineSlice(idle, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
			Disabled: image.NewNineSlice(disabled, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
		},

		Padding: widget.Insets{
			Left:   8,
			Right:  8,
			Top:    4,
			Bottom: 4,
		},

		Face: fonts.face,

		Color: &widget.TextInputColor{
			Idle:          hexToColor(textIdleColor),
			Disabled:      hexToColor(textDisabledColor),
			Caret:         hexToColor(textInputCaretColor),
			DisabledCaret: hexToColor(textInputDisabledCaretColor),
		},
	}, nil
}

func newToolTipResources(fonts *Fonts) (*toolTipResources, error) {
	bg, _, err := ebitenutil.NewImageFromFile("demorun/graphics/tool-tip.png")
	if err != nil {
		return nil, err
	}

	return &toolTipResources{
		background: image.NewNineSlice(bg, [3]int{19, 6, 13}, [3]int{19, 5, 13}),

		padding: widget.Insets{
			Left:   15,
			Right:  15,
			Top:    10,
			Bottom: 10,
		},

		face:  fonts.toolTipFace,
		color: hexToColor(toolTipColor),
	}, nil
}

func (u *UiResources) Close() {
	u.Fonts.Close()
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.RGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
