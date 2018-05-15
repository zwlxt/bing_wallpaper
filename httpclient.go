package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func getPage(url string) string {
	log.Println("Connecting...")
	resp, err := http.Get(url)
	if err != nil {
		panic("Address unreachable " + err.Error())
	}
	webpagesrcbyte, _ := ioutil.ReadAll(resp.Body)
	log.Println("Fetched page")
	return string(webpagesrcbyte)
}
