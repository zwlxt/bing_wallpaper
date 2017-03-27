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

const (
	size             = 25
	dpi              = 72
	spacing          = 1.5
	wordsPerLine     = 15
	offsetx, offsety = 1500, 100
)

func getOffset(canvas image.Image) (x, y int) {
	x, y = canvas.Bounds().Min.X+offsetx, canvas.Bounds().Min.Y+offsety
	return
}

func getWidthAndHeight(lines int) (w, h int) {
	w = wordsPerLine * (size + 2)
	h = int((spacing + size + 10) * float64(lines+1))
	return
}

func drawText(canvas image.Image, text []string) {
	rgba := image.NewRGBA(canvas.Bounds())
	draw.Draw(rgba, rgba.Bounds(), canvas, image.ZP, draw.Src)
	mask := image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 100})
	x0, y0 := getOffset(canvas)
	w, h := getWidthAndHeight(len(text))

	maskrect := image.Rect(x0, y0, x0+w, y0+h)
	draw.Draw(rgba, maskrect, mask, image.ZP, draw.Over)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(readFont("C:/Windows/Fonts/inziu-SC-regular.ttc"))
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	tc := image.NewUniform(color.Alpha16{0xC000})
	c.SetSrc(tc)
	// c.SetHinting(font.HintingNone)
	c.SetHinting(font.HintingFull)

	textAnchorX, textAnchorY := x0, y0
	pt := freetype.Pt(textAnchorX+10, textAnchorY+10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, err := c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

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

func getASCIICharCount(r []rune) int {
	result := 0
	for _, c := range r {
		if c >= 0x20 && c <= 0x7e {
			result++
		}
	}
	return result
}

func splitText(text string) []string {
	var result []string
	textRune := []rune(text)
	for i := 0; i < len(textRune); i += wordsPerLine {
		if len(textRune)-i < wordsPerLine {
			result = append(result, string(textRune[i:]))
		} else {
			ac := getASCIICharCount(textRune[i : i+wordsPerLine])
			result = append(result, string(textRune[i:i+wordsPerLine+ac]))
			i += ac
		}
	}
	return result
}

// func main() {
// 	img := readImage("D:/dev/Projects/Go/bing_wallpaper/wallpapers/2017_03_24_NoronhaTwoBrothers_ZH-CN10642407566_1920x1080.jpg")
// 	text := "Go 编程语言是一个使得程序员更加有效率的开源项目。Go 是有表达力、简洁、清晰和有效率的。它的并行机制使其很容易编写多核和网络应用，而新奇的类型系统允许构建有弹性的模块化程序。Go 编译到机器码非常快速，同时具有便利的垃圾回收和强大的运行时反射。它是快速的、静态类型编译语言，但是感觉上是动态类型的，解释型语言。"
// 	drawText(img, splitText(text))
// 	for _, s := range splitText(text) {
// 		fmt.Println(s)
// 	}
// }
