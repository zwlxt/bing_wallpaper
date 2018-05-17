package main

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type TextDrawer interface {
	Draw(img image.Image) image.Image
}

type TextDrawerConfig struct {
	FontFile          string
	FontSize          float64
	OffsetX, OffsetY  int
	LineSpacing       int
	TextWidth         int
	TextColor         color.RGBA
	BackgroundColor   color.RGBA
	BackgroundPadding int
}

type WordWrappingTextDrawer struct {
	Config TextDrawerConfig
	Text   string
}

func (d WordWrappingTextDrawer) Draw(img image.Image) image.Image {
	fontFile := d.Config.FontFile
	fontSize := d.Config.FontSize
	width := d.Config.TextWidth
	lineSpacing := d.Config.LineSpacing
	bgColor := d.Config.BackgroundColor
	bgPadding := d.Config.BackgroundPadding
	textColor := d.Config.TextColor
	x := d.Config.OffsetX
	y := d.Config.OffsetY
	canvas := image.NewRGBA(img.Bounds())
	draw.Draw(canvas, img.Bounds(), img, image.ZP, draw.Src)
	ff := fontFace(fontFile, fontSize)
	lines := wordWrap(d.Text, width, ff)
	height := paragraphHeight(lines, ff, lineSpacing)
	drawBackground(canvas, bgColor, image.Rect(x-bgPadding, y-bgPadding,
		x+bgPadding+width, y+bgPadding+height))
	drawTextWordWrap(canvas, lines, ff, textColor, lineSpacing, x, y)
	return canvas
}

func fontFace(filename string, size float64) font.Face {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Read file error " + filename + ":" + err.Error())
	}
	f, err := freetype.ParseFont(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(f, &truetype.Options{Size: size})
}

func drawBackground(canvas draw.Image, c color.Color, rect image.Rectangle) {
	bg := image.NewUniform(c)
	draw.Draw(canvas, rect, bg, image.ZP, draw.Over)
}

func paragraphHeight(text []string, ff font.Face, lineSpacing int) int {
	return ff.Metrics().Ascent.Floor() +
		(ff.Metrics().Height.Floor()+lineSpacing)*len(text) - lineSpacing
}

func wordWrap(text string, width int, ff font.Face) []string {
	lineWidth := fixed.I(0)
	rs := []rune(text)
	line := ""
	lines := make([]string, 0)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		advance, ok := ff.GlyphAdvance(r)
		if !ok {
			// skipping unknown character
			continue
		}

		if lineWidth+advance < fixed.I(width) {
			line += string(r)
			lineWidth += advance
			if r == '\n' { // handle line breakers
				// remove line breakers to prevent showing its glyph in some fonts
				line = line[:len(line)-1]
				lines = append(lines, line)
				line = ""
				lineWidth = fixed.I(0)
			}
			if i == len(rs)-1 { // last loop
				lines = append(lines, line)
			}
		} else {
			lines = append(lines, line)
			line = ""
			lineWidth = fixed.I(0)
			i--
		}
	}
	return lines
}

func drawTextWordWrap(canvas draw.Image, lines []string,
	ff font.Face, tc color.Color, lineSpacing, x, y int) {

	point := fixed.Point26_6{
		// X offset
		X: fixed.I(x),
		// Y offset of glyph
		// This value is accepted by font.Drawer as the Y value of baseline,
		// so Ascent value must be added
		Y: ff.Metrics().Ascent + fixed.I(y),
	}
	drawer := &font.Drawer{
		Src: image.NewUniform(tc),
		Dst: canvas,
		// Note that this is the baseline location
		Dot:  point,
		Face: ff,
	}

	for _, line := range lines {
		drawer.DrawString(line)
		point.Y += ff.Metrics().Height
		point.Y += fixed.I(lineSpacing)
		drawer.Dot = point
	}
}
