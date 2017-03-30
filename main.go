package main

/*
#include <windows.h>
void change_wallpaper(char path[255])
{
    SystemParametersInfo(0x0014, 0, path, 1);
}
*/
import "C"
import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"fmt"
	"html"
	"log"
	"unsafe"
)

const (
	URL  = "https://cn.bing.com/"
	DURL = "https://cn.bing.com/cnhp/life?IID=%s&IG=%s" // page containing description
)

func fetchWebPage(url string) string {
	log.Println("Connecting...")
	resp, err := http.Get(url)
	if err != nil {
		panic("Address unreachable " + err.Error())
	}
	webpagesrcbyte, _ := ioutil.ReadAll(resp.Body)
	log.Println("Fetched page")
	return string(webpagesrcbyte)
}

func saveImage(src string) string {
	re := regexp.MustCompile("g_img=\\{url:\\s\"(.+\\.jpg)\"")
	imgurl := re.FindStringSubmatch(src)[1]
	log.Println(imgurl)
	resp, err := http.Get("https://cn.bing.com" + imgurl)
	if err != nil {
		log.Println("Unable to download image " + err.Error())
		return ""
	}
	fileName := imgurl[strings.LastIndex(imgurl, "/")+1:]
	imgdata, _ := ioutil.ReadAll(resp.Body)
	cdir, _ := os.Getwd()
	if _, err := os.Stat(cdir + "/wallpapers"); os.IsNotExist(err) {
		os.Mkdir(cdir+"/wallpapers", 0666)
	}
	path := cdir + "/wallpapers/" + time.Now().Format("2006_01_02") + "_" + fileName
	err = ioutil.WriteFile(path, imgdata, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return path
}

func setWindowsWallPaper(path string) {
	path = strings.Replace(path, "/", "\\", -1)
	fmt.Println(path)

	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	C.change_wallpaper(cs)
}

func getIG(src string) string {
	re := regexp.MustCompile("IG:\"(\\w+?\\d+?)\"")
	found := re.FindStringSubmatch(src)
	if len(found) < 1 {
		return ""
	}
	log.Println(found[1])
	return found[1]
}

func getIID(src string) string {
	re := regexp.MustCompile("_iid=\"(\\w{4}\\.\\d{4})\">")
	found := re.FindStringSubmatch(src)
	if len(found) < 1 {
		return ""
	}
	log.Println(found[1])
	return found[1]
}

func getTitle(src string) string {
	re := regexp.MustCompile("<div class=\"hplaTtl\">(.+?)</div>")
	found := re.FindStringSubmatch(src)
	if len(found) < 1 {
		return ""
	}
	return found[1]
}

func getLocation(src string) string {
	re := regexp.MustCompile("<span class=\"hplaAttr\">(.+?)</span>")
	found := re.FindStringSubmatch(src)
	if len(found) < 1 {
		return ""
	}
	return found[1]
}

func getArticle(src string) (title string, subtitle string, body string) {
	re := regexp.MustCompile("<div class=\"hplatt\">(.+?)</div>")
	m := re.FindStringSubmatch(src)
	if len(m) < 1 {
		return
	}
	title = m[1]
	re = regexp.MustCompile("<div class=\"hplats\">(.+?)</div>")
	subtitle = re.FindStringSubmatch(src)[1]
	re = regexp.MustCompile("<div id=\"hplaSnippet\">(.+?)</div>")
	body = re.FindStringSubmatch(src)[1]
	return
}

func main() {
	for {
		webpagesrc := fetchWebPage(URL)
		//ioutil.WriteFile("log.html.txt", []byte(webpagesrc), 0666)

		path := saveImage(webpagesrc)
		if path == "" {
			continue
		}
		iid := getIID(webpagesrc)
		ig := getIG(webpagesrc)
		if ig == "" {
			continue
		}

		webpagesrc = fetchWebPage(fmt.Sprintf(DURL, iid, ig))
		title := getTitle(webpagesrc)
		location := getLocation(webpagesrc)
		_, _, a := getArticle(webpagesrc)
		img := readImage(path)
		a = html.UnescapeString(a)
		splitText(title)
		splitText(location)
		splitText(" ")
		splitText(a)
		fmt.Println(sl)
		drawText(img, sl)
		cdir, _ := os.Getwd()
		setWindowsWallPaper(cdir + "/out.jpg")
		log.Println("Done")
		break
	}
}
