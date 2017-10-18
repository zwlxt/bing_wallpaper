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
		ig := hc1.GetIG()
		if ig == "" {
			log.Println("Unable to get IG, Retry")
			continue
		}
		iid := hc1.GetIID()
		if iid == "" {
			log.Println("Unable to get IID, Retry")
			continue
		}

		hc2 := &HttpClient{Url: fmt.Sprintf(DURL, iid, ig)}
		hc2.FetchWebPage()
		title := hc2.GetTitle()
		location := hc2.GetLocation()
		_, _, article := hc2.GetArticle()
		fmt.Println(title)
		fmt.Println(location)

		wp := NewWallPaper(DefaultConfig())
		wp.Decode(imgdata)
		wp.AddText(title + location + "\n" + article)
		buf := wp.Encode()
		fsStorage.Save(buf, "wp_out.jpg")
		setWindowsWallPaper(installDir + "/wallpapers/" + "wp_out.jpg")
		log.Println("Done")
		break
	}
}
