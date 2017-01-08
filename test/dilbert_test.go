package main

import "testing"

type testArgs []string

func Test_validateDate(t *testing.T) {
	_, err := validateDate(testArgs{""})
	if err != errTooManyArgs {
		t.Error("Expected errTooManyArgs")
	}

	err = nil
	_, err = validateDate(testArgs{"", "1"})
	if err != errInvalidDate {
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
	if err != errInvalidDate {
		t.Error("Expected errInvalidDate")
	}
}
