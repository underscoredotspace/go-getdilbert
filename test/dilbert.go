package main

import (
	"errors"
	"log"
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

	log.Println(stripDate)
}

/*
	* Check supplied date
	Load page relevant to date
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
	if _, err = time.Parse("2006-01-02", args[1]); err != nil {
		return "", errInvalidDate
	}

	return args[1], nil
}
