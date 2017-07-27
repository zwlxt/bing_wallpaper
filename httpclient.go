package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type HttpClient struct {
	Url, htmlSrc string
}

func (hc *HttpClient) FetchWebPage() {
	log.Println("Connecting...")
	resp, err := http.Get(hc.Url)
	if err != nil {
		panic("Address unreachable " + err.Error())
	}
	webpagesrcbyte, _ := ioutil.ReadAll(resp.Body)
	log.Println("Fetched page")
	hc.htmlSrc = string(webpagesrcbyte)
}

func (hc *HttpClient) GetImage() (string, []byte) {
	re := regexp.MustCompile("g_img=\\{url:\\s\"(.+\\.jpg)\"")
	imgurl := re.FindStringSubmatch(hc.htmlSrc)[1]
	log.Println(imgurl)
	resp, err := http.Get(hc.Url + imgurl)
	if err != nil {
		log.Println("Unable to download image " + err.Error())
	}
	fileName := imgurl[strings.LastIndex(imgurl, "/")+1:]
	imgdata, _ := ioutil.ReadAll(resp.Body)
	return fileName, imgdata
}

func (hc *HttpClient) GetIG() string {
	re := regexp.MustCompile("IG:\"(\\w+?\\d+?)\"")
	found := re.FindStringSubmatch(hc.htmlSrc)
	if len(found) < 1 {
		return ""
	}
	log.Println("IG:" + found[1])
	return found[1]
}

func (hc *HttpClient) GetIID() string {
	re := regexp.MustCompile("_iid=\"(\\w{4}\\.\\d{4})\">")
	found := re.FindStringSubmatch(hc.htmlSrc)
	if len(found) < 1 {
		return ""
	}
	log.Println("IID:" + found[1])
	return found[1]
}

func (hc *HttpClient) GetTitle() string {
	re := regexp.MustCompile("<div class=\"hplaTtl\">(.+?)</div>")
	found := re.FindStringSubmatch(hc.htmlSrc)
	if len(found) < 1 {
		return ""
	}
	return found[1]
}

func (hc *HttpClient) GetLocation() string {
	re := regexp.MustCompile("<span class=\"hplaAttr\">(.+?)</span>")
	found := re.FindStringSubmatch(hc.htmlSrc)
	if len(found) < 1 {
		return ""
	}
	return found[1]
}

func (hc *HttpClient) GetArticle() (title string, subtitle string, body string) {
	re := regexp.MustCompile("<div class=\"hplatt\">(.+?)</div>")
	m := re.FindStringSubmatch(hc.htmlSrc)
	if len(m) < 1 {
		return
	}
	title = m[1]
	re = regexp.MustCompile("<div class=\"hplats\">(.+?)</div>")
	subtitle = re.FindStringSubmatch(hc.htmlSrc)[1]
	re = regexp.MustCompile("<div id=\"hplaSnippet\">(.+?)</div>")
	body = re.FindStringSubmatch(hc.htmlSrc)[1]
	return
}
