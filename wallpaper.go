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
	"strings"

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

func NewWallPaper(config Config) *WallPaper {
	return &WallPaper{
		config: config,
	}
}

func (wallpaper *WallPaper) AddText(text string) {
	canvas := wallpaper.img
	rgba := image.NewRGBA(canvas.Bounds())                         // image to draw, same size as input image
	draw.Draw(rgba, rgba.Bounds(), canvas, image.ZP, draw.Src)     // copy input image to it
	mask := image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 100}) // background of text
	x0, y0 := wallpaper.getOffset(canvas)                          // x, y coordinate

	// font rendering
	c := freetype.NewContext()
	c.SetDPI(wallpaper.config.FontDPI)
	c.SetFont(wallpaper.loadFont(wallpaper.config.FontFile))
	c.SetFontSize(wallpaper.config.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	tc := image.NewUniform(color.Alpha16{wallpaper.config.TextColor})
	c.SetSrc(tc)
	c.SetHinting(font.HintingFull)

	textX, textY := x0, y0 // text x, y
	pt := freetype.Pt(textX+10, textY+10+int(c.PointToFixed(wallpaper.config.FontSize)>>6))

	lines := wallpaper.splitText(text)
	for _, s := range lines {
		_, err := c.DrawString(s, pt)
		if err != nil {
			panic(err)
		}
		pt.Y += c.PointToFixed(wallpaper.config.FontSize * wallpaper.config.FontSpacing)
	}

	w, h := wallpaper.getWidthAndHeight(len(lines)) // calculate width of the text

	maskrect := image.Rect(x0, y0, x0+w, y0+h)           // background location and size
	draw.Draw(rgba, maskrect, mask, image.ZP, draw.Over) // draw background
	wallpaper.img = rgba
}

var sl = make([]string, 0)

func (wallpaper WallPaper) splitText(s string) []string {
	wordsPerLine := wallpaper.config.WordsPerLine

	// detect line breakers
	lineBreakerIndex := strings.Index(s, "\n")
	if lineBreakerIndex != -1 {
		sl = append(sl, s[:lineBreakerIndex])
		s = s[lineBreakerIndex+1:]
	}

	r := []rune(s)
	if len(r) > wordsPerLine {
		ac := wallpaper.getASCIICharCount(r[:wordsPerLine])
		forward := 0
		i := 0

		for ; forward < ac; i++ {
			if wordsPerLine+i >= len(r) {
				// the entire paragraph is shorter than width
				break
			}

			if !wallpaper.isASCIIChar(r[wordsPerLine+i]) {
				forward += 2
			} else {
				forward++
			}
		}
		sl = append(sl, string(r[:wordsPerLine+i]))
		wallpaper.splitText(string(r[wordsPerLine+i:]))
	} else {
		sl = append(sl, s)
	}
	return sl
}

func (WallPaper) isASCIIChar(c rune) bool {
	return c >= 0x20 && c <= 0x7e || c == '\n'
}

func (wallpaper WallPaper) getASCIICharCount(r []rune) int {
	result := 0
	for _, c := range r {
		if wallpaper.isASCIIChar(c) {
			result++
		}
	}
	return result
}

func (wallpaper *WallPaper) Decode(buf []byte) {
	r := bytes.NewReader(buf)
	image, err := jpeg.Decode(r)
	if err != nil {
		panic(err)
	}
	wallpaper.img = image
}

func (wallpaper *WallPaper) Encode() []byte {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, wallpaper.img, &jpeg.Options{Quality: 100})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
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

func (wallpaper WallPaper) getOffset(canvas image.Image) (x, y int) {
	x = canvas.Bounds().Min.X + wallpaper.config.OffsetX
	y = canvas.Bounds().Min.Y + wallpaper.config.OffsetY
	return
}

func (wallpaper WallPaper) getWidthAndHeight(lines int) (w, h int) {
	w = wallpaper.config.WordsPerLine * (int(wallpaper.config.FontSize) + 2)
	h = int((wallpaper.config.FontSpacing + wallpaper.config.FontSize + 8) * float64(lines+1))
	return
}
