package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var errInvalidDate = errors.New("Usage - " + os.Args[0] + " [yyyy-mm-dd]")
var errFailedToFindImageAddr = errors.New("Failed to find image path")
var err404 = errors.New("404 - Page Not Found")

func main() {
	stripDate, err := validateDate(os.Args)
	if err != nil {
		log.Fatalln(err.Error())
	}

	stripPage, err := getStripPage(stripDate)
	if err != nil {
		log.Fatalln(err.Error())
	}

	stripImageAddr, err := getStripImageAddr(stripPage)
	if err != nil {
		log.Fatalln(err.Error())
	}

	stripImage, err := getStripImage(stripImageAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = saveStripImage("./images/", stripImage, stripDate)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Saved")
}

// validateDate takes os.Args as a slice and checks date in the valid format has been provided
//   no other args required, so returns error for too many args
func validateDate(args []string) (stripDate string, err error) {
	// If we didn't get any arguments, assume stripDate is today
	if len(args[1:]) == 0 {
		return time.Now().Format("2006-01-02"), nil
	}

	// Attempt to parse os.Args[1] with our date format to ensure it will work
	if _, parseErr := time.Parse("2006-01-02", args[1]); parseErr != nil {
		return "", errInvalidDate
	}

	return args[1], nil
}

func getStripPage(stripDate string) (stripPage []byte, err error) {
	// This prevents http.get from following redirects, as with dilbert.com/strip/* there is no 404 merely redirect to homepage.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	stripPageAddr := "http://dilbert.com/strip/" + stripDate
	res, err := client.Get(stripPageAddr)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return nil, err404
	}

	stripPage, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

func getStripImageAddr(stripPage []byte) (stripImageAddr string, err error) {
	var imgRegEx = regexp.MustCompile(`data-image="http:\/\/assets\.amuniversal.com\/([a-z0-9]+)"`)
	if imgRegExMatches := imgRegEx.FindAllSubmatch(stripPage, 1); len(imgRegExMatches) == 1 {
		stripImageAddr = "http://assets.amuniversal.com/" + string(imgRegExMatches[0][1])
	} else {
		err = errFailedToFindImageAddr
	}
	return
}

func getStripImage(stripImageAddr string) (stripImage []byte, err error) {
	res, err := http.Get(stripImageAddr)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return nil, err404
	}

	stripImage, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

func saveStripImage(stripsPath string, stripImage []byte, stripDate string) error {
	// Parse date in to time.Time variable to allow us to pull elements for folder/filename
	stripDateParsed, err := time.Parse("2006-01-02", stripDate)
	if err != nil {
		return err
	}

	// Build the path, adding any missing slashes
	savePath := filepath.Join(stripsPath, stripDateParsed.Format("2006/01/02")+".gif")

	// Check to see if file already exists
	_, err = os.Stat(savePath)
	if err == nil {
		// No point going further if it does
		return os.ErrExist
	}

	// Extract directory from savePath
	saveDir, _ := filepath.Split(savePath)

	// Make saveDir including required parents
	err = os.MkdirAll(saveDir, 0744)
	if err != nil {
		return err
	}

	// Save the strip itself
	return ioutil.WriteFile(savePath, stripImage, 0644)
}
