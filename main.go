package main

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"
	"log"
)

func main() {
	var config Config
	installDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configFile := installDir + "/config.yml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Println("Using default config")
		config = Default()
		config.Save(configFile)
	} else {
		log.Println("Using config: " + configFile)
		config.Load(configFile)
	}

	wallpaperDir := config.WallpaperDir + "/wallpapers/"
	fsStorage := &FileSystemStorage{Dir: wallpaperDir}
	var ig, iid string
	var wallpaper *WallPaper
	hasErr := true
	for i := 0; i < 10; i++ {
		bp := NewBingPage()
		imgURL := bp.ImageURL()
		wallpaper = FromURL(imgURL)
		if wallpaper == nil {
			continue
		}
		fileName := imgURL[strings.LastIndex(imgURL, "/")+1:]
		wallpaper.SaveToFile(fsStorage, fileName, 100)
		ig = bp.IG()
		if ig == "" {
			log.Println("Unable to get IG, Retry")
			continue
		}
		iid = bp.IID()
		if iid == "" {
			log.Println("Unable to get IID, Retry")
			continue
		}
		hasErr = false
		break
	}
	if hasErr {
		panic("Failed after 10 retries")
	}
	for i := 0; i < 10; i++ {
		bip := NewBingInfoPage(iid, ig)
		title := bip.Title()
		location := bip.Location()
		_, _, article := bip.Article()
		fmt.Println(title)
		fmt.Println(location)
		fmt.Println(article)
		if config.TextDrawerEnabled {
			wallpaper.SetTextDrawer(&WordWrappingTextDrawer{
				Config: config.TextDrawerConfig,
				Text:   title + "," + location + "\n" + article,
			})
		}
		wallpaper.SaveToFile(fsStorage, "wp_out.jpg", 100)
		absWallpaperPath, err := filepath.Abs(wallpaperDir + "/wp_out.jpg")
		if err != nil {
			panic(err)
		}
		setWindowsWallPaper(absWallpaperPath)
		log.Println("Done")
		break
	}
}
