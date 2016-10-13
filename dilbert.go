package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func getDilbertImage(date string, path string) (filename string, myerr error) {
	if res, err := http.Get(path); err == nil {
		if image, err := ioutil.ReadAll(res.Body); err == nil {
			defer res.Body.Close()
			filename = "./" + date + ".gif"
			if err := ioutil.WriteFile(filename, image, 0644); err != nil {
				myerr = err
			}
		} else {
			myerr = err
		}
	} else {
		myerr = err
	}
	return
}

func getDilbertImagePath(date string) (path string, myerr error) {
	var imgRegEx = regexp.MustCompile(`data-image="http:\/\/assets\.amuniversal.com\/([a-z0-9]+)"\s`)
	if res, err := http.Get("http://dilbert.com/strip/" + date); err == nil {
		if res.StatusCode == 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err == nil {
				defer res.Body.Close()

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
	//date := "2016-10-12"
	date := time.Now().Format("2006-01-02")
	if path, err := getDilbertImagePath(date); err == nil {
		fmt.Println(path, "found")
		if filename, err := getDilbertImage(date, path); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(filename, "saved")
		}
	} else {
		log.Fatal(err)
	}
}
