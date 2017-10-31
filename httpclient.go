package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	u, _ := url.Parse(hc.Url)
	var resp *http.Response
	var err error
	if strings.HasPrefix(imgurl, "//") {
		resp, err = http.Get(u.Scheme + ":" + imgurl)
	} else {
		resp, err = http.Get(u.Scheme + "://" + u.Hostname() + imgurl)
	}
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
	title := html.UnescapeString(found[1])
	return title
}

func (hc *HttpClient) GetLocation() string {
	re := regexp.MustCompile("<span class=\"hplaAttr\">(.+?)</span>")
	found := re.FindStringSubmatch(hc.htmlSrc)
	if len(found) < 1 {
		return ""
	}
	location := html.UnescapeString(found[1])
	return location
}

func (hc *HttpClient) GetArticle() (title string, subtitle string, body string) {
	re := regexp.MustCompile("<div class=\"hplatt\">(.+?)</div>")
	m := re.FindStringSubmatch(hc.htmlSrc)
	if len(m) < 1 {
		return
	}
	title = html.UnescapeString(m[1])
	re = regexp.MustCompile("<div class=\"hplats\">(.+?)</div>")
	subtitle = html.UnescapeString(re.FindStringSubmatch(hc.htmlSrc)[1])
	re = regexp.MustCompile("<div id=\"hplaSnippet\">(.+?)</div>")
	body = html.UnescapeString(re.FindStringSubmatch(hc.htmlSrc)[1])
	return
}
