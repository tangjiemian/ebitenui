package gui

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	fontFaceRegular = "demorun/fonts/NotoSans-Regular.ttf"
	fontFaceBold    = "demorun/fonts/NotoSans-Bold.ttf"
)

type Fonts struct {
	face         font.Face
	titleFace    font.Face
	bigTitleFace font.Face
	toolTipFace  font.Face
}

func LoadFonts() (*Fonts, error) {
	fontFace, err := LoadFont(fontFaceRegular, 20)
	if err != nil {
		return nil, err
	}

	titleFontFace, err := LoadFont(fontFaceBold, 24)
	if err != nil {
		return nil, err
	}

	bigTitleFontFace, err := LoadFont(fontFaceBold, 28)
	if err != nil {
		return nil, err
	}

	toolTipFace, err := LoadFont(fontFaceRegular, 15)
	if err != nil {
		return nil, err
	}

	return &Fonts{
		face:         fontFace,
		titleFace:    titleFontFace,
		bigTitleFace: bigTitleFontFace,
		toolTipFace:  toolTipFace,
	}, nil
}

func (f *Fonts) Close() {
	if f.face != nil {
		_ = f.face.Close()
	}

	if f.titleFace != nil {
		_ = f.titleFace.Close()
	}

	if f.bigTitleFace != nil {
		_ = f.bigTitleFace.Close()
	}
}

func LoadFont(path string, size float64) (font.Face, error) {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
