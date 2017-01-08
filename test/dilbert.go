package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var errTooManyArgs = errors.New("One argument required - the date in format yyyy-mm-dd")
var errInvalidDate = errors.New("Invalid date provided")

func main() {
	stripDate, err := validateDate(os.Args)
	if err != nil {
		log.Fatalln(err.Error())
	}

	stripPageData, err := getStripPage(stripDate)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(stripPageData)
}

/*
	* Check supplied date
	* Load page relevant to date
	Find strip image with regex
	Load image
	Save image
*/

// validateDate takes os.Args as a slice and checks date in the valid format has been provided
//   no other args required, so returns error for too many args
func validateDate(args []string) (stripDate string, err error) {
	// Check to see we got one arg, no more no less
	if len(args[1:]) != 1 {
		return "", errTooManyArgs
	}

	// Attempt to parse os.Args[1] with our date format to ensure it will work
	if _, parseErr := time.Parse("2006-01-02", args[1]); parseErr != nil {
		log.Println(parseErr.Error())
		return "", errInvalidDate
	}

	return args[1], nil
}

func getStripPage(stripDate string) (stripPage []byte, err error) {
	// This prevents http.get from following redirects. In our case there is no 404 merely redirect to homepage.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	stripPageAddr := "http://dilbert.com/strip/" + stripDate
	res, err := client.Get(stripPageAddr)
	if err != nil {
		return
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Println(res.Header.Get("location") + " - " + res.Status)
		return nil, errors.New("Error loading page " + stripPageAddr)
	}

	stripPage, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}
