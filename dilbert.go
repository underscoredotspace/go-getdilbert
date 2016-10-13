package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	var imgRegEx = regexp.MustCompile(`data-image="http:\/\/assets\.amuniversal.com\/([a-z0-9]+)"\s`)

	if res1, err := http.Get("https://op11.co.uk/dilbert.html"); err == nil && res1.StatusCode == 200 {
		// If page loaded, open document body
		body, err := ioutil.ReadAll(res1.Body)
		if err == nil {
			defer res1.Body.Close()

			// Find the image path for the strip using RegEx
			imgRegExMatches := imgRegEx.FindAllSubmatch(body, 1)
			if len(imgRegExMatches) == 1 {
				if len(imgRegExMatches[0]) == 2 {
					// pass
				} else {
					log.Fatal("Failed to find image match in body (1)")
				}
			} else {
				log.Fatal("Failed to find image match in body (2)")
			}
			imgURL := "http://assets.amuniversal.com/" + string(imgRegExMatches[0][1])
			if res2, err := http.Get(imgURL); err == nil && res2.StatusCode == 200 && res2.Header.Get("Content-Type") != "text/html; charset=utf-8" {
				image, err := ioutil.ReadAll(res2.Body)
				if err == nil {
					defer res2.Body.Close()
					fmt.Println("Loaded image")
					// save it to disk
					err := ioutil.WriteFile("./today.gif", image, 0644)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Fatal("Failed to load image")
				}
			} else {
				log.Fatal("Failed to get image from RegEx")
			}
		}
	} else {
		// Epic fail, we've not got page from dilbert.com
		log.Fatalf("%s\n", res1.Status)
	}
}
