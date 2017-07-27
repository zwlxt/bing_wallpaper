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
		dbStorage := NewSqliteStorage(installDir + "/bingwallpaper.db")

		hc1 := &HttpClient{Url: URL}
		hc1.FetchWebPage()
		fileName, imgdata := hc1.GetImage()
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
		title := hc2.GetTitle()
		location := hc2.GetLocation()
		_, _, description := hc2.GetArticle()
		fmt.Println(title)
		fmt.Println(location)
		dbStorage.AdditionDescription(title, location, description)
		dbStorage.Save(imgdata, fileName)

		log.Println("Done")
		break
	}
}
