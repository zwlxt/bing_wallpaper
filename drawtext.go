package main

import (
	"bufio"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/font"

	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const (
	size             = 20
	dpi              = 72
	spacing          = 1.5
	wordsPerLine     = 16
	offsetx, offsety = 1500, 50
	textColor        = 0xe000
)

func readImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	image, err := jpeg.Decode(f)
	if err != nil {
		panic(err)
	}
	return image
}

func readFont(path string) *truetype.Font {
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

func getOffset(canvas image.Image) (x, y int) {
	x, y = canvas.Bounds().Min.X+offsetx, canvas.Bounds().Min.Y+offsety
	return
}

func getWidthAndHeight(lines int) (w, h int) {
	w = wordsPerLine * (size + 2)
	h = int((spacing + size + 8) * float64(lines+1))
	return
}

func drawText(canvas image.Image, text []string) {
	rgba := image.NewRGBA(canvas.Bounds())                         // image to draw, same size as input image
	draw.Draw(rgba, rgba.Bounds(), canvas, image.ZP, draw.Src)     // copy input image to it
	mask := image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 100}) // background of text
	x0, y0 := getOffset(canvas)                                    // x, y coordinate
	w, h := getWidthAndHeight(len(text))                           // calculate width of the text

	maskrect := image.Rect(x0, y0, x0+w, y0+h)           // background location and size
	draw.Draw(rgba, maskrect, mask, image.ZP, draw.Over) // draw background

	// font rendering
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(readFont("inziu-SC-regular.ttc"))
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	tc := image.NewUniform(color.Alpha16{textColor})
	c.SetSrc(tc)
	// c.SetHinting(font.HintingNone)
	c.SetHinting(font.HintingFull)

	textX, textY := x0, y0 // text x, y
	pt := freetype.Pt(textX+10, textY+10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, err := c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	// output to temp file
	outFile, err := os.Create("out.jpg")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = jpeg.Encode(b, rgba, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func isASCIIChar(c rune) bool {
	return c >= 0x20 && c <= 0x7e
}

func getASCIICharCount(r []rune) int {
	result := 0
	for _, c := range r {
		if c >= 0x20 && c <= 0x7e {
			result++
		}
	}
	return result
}

var sl []string

func splitText(s string) {
	r := []rune(s)
	if len(r) > wordsPerLine {
		ac := getASCIICharCount(r[:wordsPerLine])
		forward := 0
		i := 0
		for ; forward < ac; i++ {
			if !isASCIIChar(r[wordsPerLine+i]) {
				forward += 2
			} else {
				forward++
			}
		}
		sl = append(sl, string(r[:wordsPerLine+i]))
		splitText(string(r[wordsPerLine+i:]))
	} else {
		sl = append(sl, s)
	}
}

// func main() {
// 	img := readImage("D:/dev/Projects/Go/bing_wallpaper/wallpapers/2017_03_24_NoronhaTwoBrothers_ZH-CN10642407566_1920x1080.jpg")
// 	text := "Go 编程语言是一个使得程序员更加有效率的开源项目。Go 是有表达力、简洁、清晰和有效率的。它的并行机制使其很容易编写多核和网络应用，而新奇的类型系统允许构建有弹性的模块化程序。Go 编译到机器码非常快速，同时具有便利的垃圾回收和强大的运行时反射。它是快速的、静态类型编译语言，但是感觉上是动态类型的，解释型语言。"
// 	drawText(img, splitText(text))
// 	for _, s := range splitText(text) {
// 		fmt.Println(s)
// 	}
// }
