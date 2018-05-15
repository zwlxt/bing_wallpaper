package main

import (
	"bytes"
	"image"
	"image/jpeg"
)

func decode(buf []byte) image.Image {
	r := bytes.NewReader(buf)
	img, err := jpeg.Decode(r)
	if err != nil {
		panic(err)
	}
	return img
}

func encode(img image.Image, quality int) []byte {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
