package main

import (
	"fmt"
	"html"
	"log"
	"regexp"
)

const (
	BaseURL = "https://cn.bing.com"
	InfoURL = "https://cn.bing.com/cnhp/life?IID=%s&IG=%s"
)

type BingPage struct {
	src string
}

func NewBingPage() *BingPage {
	return &BingPage{getPage(BaseURL)}
}

func (b BingPage) IID() string {
	re := regexp.MustCompile("_iid=\"(\\w{4}\\.\\d{4})\">")
	found := re.FindStringSubmatch(b.src)
	if len(found) < 1 {
		return ""
	}
	log.Println("IID:" + found[1])
	return found[1]
}

func (b BingPage) IG() string {
	re := regexp.MustCompile("IG:\"(\\w+?\\d+?)\"")
	found := re.FindStringSubmatch(b.src)
	if len(found) < 1 {
		return ""
	}
	log.Println("IG:" + found[1])
	return found[1]
}

func (b BingPage) ImageURL() string {
	re := regexp.MustCompile("g_img=\\{url:\\s\"(.+\\.jpg)\"")
	m := re.FindStringSubmatch(b.src)
	if len(m) < 1 {
		return ""
	}
	imgurl := BaseURL + m[1]
	log.Println(imgurl)
	return imgurl
}

type BingInfoPage struct {
	src string
}

func NewBingInfoPage(iid, ig string) *BingInfoPage {
	return &BingInfoPage{getPage(fmt.Sprintf(InfoURL, iid, ig))}
}

func (b BingInfoPage) Title() string {
	re := regexp.MustCompile("<div class=\"hplaTtl\">(.+?)</div>")
	found := re.FindStringSubmatch(b.src)
	if len(found) < 1 {
		return ""
	}
	title := html.UnescapeString(found[1])
	return title
}

func (b BingInfoPage) Location() string {
	re := regexp.MustCompile("<span class=\"hplaAttr\">(.+?)</span>")
	found := re.FindStringSubmatch(b.src)
	if len(found) < 1 {
		return ""
	}
	location := html.UnescapeString(found[1])
	return location
}

func (b BingInfoPage) Article() (title string, subtitle string, body string) {
	re := regexp.MustCompile("<div class=\"hplatt\">(.+?)</div>")
	m := re.FindStringSubmatch(b.src)
	if len(m) < 1 {
		return
	}
	title = html.UnescapeString(m[1])
	re = regexp.MustCompile("<div class=\"hplats\">(.+?)</div>")
	subtitle = html.UnescapeString(re.FindStringSubmatch(b.src)[1])
	re = regexp.MustCompile("<div id=\"hplaSnippet\">(.+?)</div>")
	body = html.UnescapeString(re.FindStringSubmatch(b.src)[1])
	return
}
