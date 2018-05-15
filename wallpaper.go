package main

import (
	"image"
	"io/ioutil"
	"log"
	"net/http"
)

type WallPaper struct {
	img image.Image
}

func FromURL(URL string) *WallPaper {
	resp, err := http.Get(URL)
	if err != nil {
		log.Println("Unable to download image " + err.Error())
		return nil
	}
	imgdata, _ := ioutil.ReadAll(resp.Body)
	return &WallPaper{decode(imgdata)}
}

func (wallpaper *WallPaper) FromLocalStorage(storage StorageManager, fileName string) *WallPaper {
	b := storage.Load(fileName)
	return &WallPaper{decode(b)}
}

func (wallpaper *WallPaper) SetImage(img image.Image) {
	wallpaper.img = img
}

func (wallpaper *WallPaper) SetTextDrawer(textDrawer TextDrawer) {
	wallpaper.img = textDrawer.Draw(wallpaper.img)
}

func (wallpaper WallPaper) SaveToFile(storage StorageManager, fileName string, quality int) {
	b := encode(wallpaper.img, quality)
	storage.Save(b, fileName)
}

func (wallpaper WallPaper) Image() image.Image {
	return wallpaper.img
}
