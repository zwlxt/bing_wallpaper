package main

import (
	"os"
	"path/filepath"

	"fmt"
	"log"
)

const (
	URL  = "https://cn.bing.com/?FORM=HPENCN&setmkt=zh-cn&setlang=zh-cn"
	DURL = "https://cn.bing.com/cnhp/life?IID=%s&IG=%s" // page containing description
)

func main() {
	for {
		installDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		wallpaperDir := installDir + "/wallpapers/"
		fsStorage := &FileSystemStorage{Dir: wallpaperDir}

		hc1 := &HttpClient{Url: URL}
		hc1.FetchWebPage()
		fileName, imgdata := hc1.GetImage()
		fsStorage.Save(imgdata, fileName)
		setWindowsWallPaper(wallpaperDir + fileName)
		ig := hc1.GetIG()
		if ig == "" {
			continue
		}
		iid := hc1.GetIID()
		if iid == "" {
			continue
		}

		hc2 := &HttpClient{Url: fmt.Sprintf(DURL, iid, ig)}
		hc2.FetchWebPage()
		fmt.Println(hc2.GetTitle())
		fmt.Println(hc2.GetLocation())

		// webpagesrc = fetchWebPage(fmt.Sprintf(DURL, iid, ig))
		// title := getTitle(webpagesrc)
		// location := getLocation(webpagesrc)
		// _, _, a := getArticle(webpagesrc)
		// img := readImage(path)
		// a = html.UnescapeString(a)
		// splitText(title)
		// splitText(location)
		// splitText(" ")
		// splitText(a)
		// fmt.Println(sl)
		// drawText(img, sl)
		// cdir, _ := os.Getwd()
		// setWindowsWallPaper(cdir + "/out.jpg")
		log.Println("Done")
		break
	}
}
