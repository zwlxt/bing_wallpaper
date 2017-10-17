package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Config struct {
	FontFile         string
	FontSize         float64
	FontDPI          float64
	FontSpacing      float64
	WordsPerLine     int
	OffsetX, OffsetY int
	TextColor        uint16
}

func DefaultConfig() Config {
	return Config{
		FontFile:     "inziu-SC-regular.ttc",
		FontSize:     20,
		FontDPI:      72,
		FontSpacing:  1.5,
		WordsPerLine: 16,
		OffsetX:      1500,
		OffsetY:      50,
		TextColor:    0xe000,
	}
}

type WallPaper struct {
	img    image.Image
	config Config
}

func (this *WallPaper) AddText(text string) {
	canvas := this.img
	rgba := image.NewRGBA(canvas.Bounds())                         // image to draw, same size as input image
	draw.Draw(rgba, rgba.Bounds(), canvas, image.ZP, draw.Src)     // copy input image to it
	mask := image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 100}) // background of text
	x0, y0 := this.getOffset(canvas)                                    // x, y coordinate
	w, h := this.getWidthAndHeight(len(text))                           // calculate width of the text

	maskrect := image.Rect(x0, y0, x0+w, y0+h)           // background location and size
	draw.Draw(rgba, maskrect, mask, image.ZP, draw.Over) // draw background

	// font rendering
	c := freetype.NewContext()
	c.SetDPI(this.config.FontDPI)
	c.SetFont(this.loadFont(this.config.FontFile))
	c.SetFontSize(this.config.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	tc := image.NewUniform(color.Alpha16{this.config.TextColor})
	c.SetSrc(tc)
	// c.SetHinting(font.HintingNone)
	c.SetHinting(font.HintingFull)

	textX, textY := x0, y0 // text x, y
	pt := freetype.Pt(textX+10, textY+10+int(c.PointToFixed(this.config.FontSize)>>6))

	lines := this.splitText(text)
	for _, s := range lines {
		_, err := c.DrawString(s, pt)
		if err != nil {
			panic(err)
		}
		pt.Y += c.PointToFixed(this.config.FontSize * this.config.FontSpacing)
	}
	this.img = rgba
}

var sl = make([]string, 0)

func (this WallPaper) splitText(s string) []string {
	r := []rune(s)
	wordsPerLine := this.config.WordsPerLine
	if len(r) > wordsPerLine {
		ac := this.getASCIICharCount(r[:wordsPerLine])
		forward := 0
		i := 0
		for ; forward < ac; i++ {
			if wordsPerLine+i >= len(r) {
				break
			}
			if !this.isASCIIChar(r[wordsPerLine+i]) {
				forward += 2
			} else {
				forward++
			}
		}
		sl = append(sl, string(r[:wordsPerLine+i]))
		this.splitText(string(r[wordsPerLine+i:]))
	} else {
		sl = append(sl, s)
	}
	return sl
}

func (WallPaper) isASCIIChar(c rune) bool {
	return c >= 0x20 && c <= 0x7e
}

func (WallPaper) getASCIICharCount(r []rune) int {
	result := 0
	for _, c := range r {
		if c >= 0x20 && c <= 0x7e {
			result++
		}
	}
	return result
}

func (this *WallPaper) Decode(buf []byte) {
	r := bytes.NewReader(buf)
	image, err := jpeg.Decode(r)
	if err != nil {
		panic(err)
	}
	this.img = image
}

func (WallPaper) loadFont(path string) *truetype.Font {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = "C:/Windows/Fonts/" + path
	}
	log.Println("Font:" + path)
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	return font
}

func (this WallPaper) getOffset(canvas image.Image) (x, y int) {
	x = canvas.Bounds().Min.X + this.config.OffsetX
	y = canvas.Bounds().Min.Y + this.config.OffsetY
	return
}

func (this WallPaper) getWidthAndHeight(lines int) (w, h int) {
	w = this.config.WordsPerLine * (int(this.config.FontSize) + 2)
	h = int((this.config.FontSpacing + this.config.FontSize + 8) * float64(lines+1))
	return
}
