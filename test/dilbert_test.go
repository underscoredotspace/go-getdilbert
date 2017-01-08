package main

import "testing"

type testArgs []string

func Test_validateDate(t *testing.T) {
	_, err := validateDate(testArgs{""})
	if err == nil || err != errTooManyArgs {
		t.Error("Expected errTooManyArgs")
	}

	err = nil
	_, err = validateDate(testArgs{"", "1"})
	if err == nil || err != errInvalidDate {
		t.Error("Expected errInvalidDate")
	}

	var stripDate string
	err = nil
	stripDate, err = validateDate(testArgs{"", "2016-01-30"})
	if stripDate != "2016-01-30" {
		t.Error("Expected 2016-01-30")
	}

	err = nil
	_, err = validateDate(testArgs{"", "2016-99-30"})
	if err == nil || err != errInvalidDate {
		t.Error("Expected errInvalidDate")
	}

	err = nil
	_, err = validateDate(testArgs{"", "1983-02-29"})
	if err == nil || err != errInvalidDate {
		t.Error("Expected errInvalidDate")
	}
}

func Test_getStripPage(t *testing.T) {
	_, err := getStripPage("garbage")
	if err == nil || err.Error() != "Error loading page http://dilbert.com/strip/garbage" {
		t.Error("Expected \"Error loading page http://dilbert.com/strip/garbage\"")
	}

	err = nil
	_, err = getStripPage("1983-02-28")
	if err == nil || err.Error() != "Error loading page http://dilbert.com/strip/1983-02-28" {
		t.Error("Expected \"Error loading page http://dilbert.com/strip/1983-02-28\"")
	}
}
