package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func getDilbertImage(path string) (myerr error) {
	fmt.Println(path)
	return
}

func getDilbertImagePath(date string) (path string, myerr error) {
	var imgRegEx = regexp.MustCompile(`data-image="http:\/\/assets\.amuniversal.com\/([a-z0-9]+)"\s`)
	if res, err := http.Get("http://dilbert.com/strip/" + date); err == nil {
		if res.StatusCode == 200 {
			// If page loaded, open document body
			body, err := ioutil.ReadAll(res.Body)
			if err == nil {
				defer res.Body.Close()

				// Find the image path for the strip using RegEx
				if imgRegExMatches := imgRegEx.FindAllSubmatch(body, 1); len(imgRegExMatches) == 1 {
					path = "http://assets.amuniversal.com/" + string(imgRegExMatches[0][1])
				} else {
					myerr = errors.New("Failed to find image path for " + date)
				}
			} else {
				myerr = err
			}
		} else {
			myerr = errors.New("image for " + date + " " + res.Status)
		}
	} else {
		myerr = err
	}
	return
}

func main() {
	date := "2016-10-12"
	if path, err := getDilbertImagePath(date); err == nil {
		getDilbertImage(path)
	} else {
		log.Fatal(err)
	}
	// if err := getDilbertImage(path); err == nil {
	// 	fmt.Println("Dilbert image for " + date + " downloaded")
	// }

	/*
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
	*/
}
