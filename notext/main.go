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
	"unsafe"
)

func main() {
	resp, err := http.Get("https://cn.bing.com/")
	if err != nil {
		panic("Address unreachable " + err.Error())
	}
	webpagesrcbyte, _ := ioutil.ReadAll(resp.Body)
	webpagesrc := string(webpagesrcbyte)
	re := regexp.MustCompile("g_img=\\{url:\\s\"(.+\\.jpg)\"")
	imgurl := re.FindStringSubmatch(webpagesrc)[1]
	resp, err = http.Get("https://cn.bing.com" + imgurl)
	if err != nil {
		panic("Unable to download image " + err.Error())
	}
	fileName := imgurl[strings.LastIndex(imgurl, "/")+1:]
	imgdata, _ := ioutil.ReadAll(resp.Body)
	cdir, _ := os.Getwd()
	path := cdir + "/wallpapers/" + time.Now().Format("2006_01_02") + "_" + fileName
	err = ioutil.WriteFile(path, imgdata, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	path = strings.Replace(path, "/", "\\", -1)
	fmt.Println(path)

	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	C.change_wallpaper(cs)
}
