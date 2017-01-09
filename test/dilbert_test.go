package main

import "testing"

func Test_validateDate(t *testing.T) {
	type testArgs []string

	var validateDateTests = []struct {
		test      testArgs
		stripDate string
		err       error
	}{
		{testArgs{""}, "", errNumOfArgs},
		{testArgs{"", "1"}, "", errInvalidDate},
		{testArgs{"", "2016-99-30"}, "", errInvalidDate},
		{testArgs{"", "1983-02-29"}, "", errInvalidDate},
		{testArgs{"", "2016-01-30"}, "2016-01-30", nil},
	}

	for _, tt := range validateDateTests {
		stripDate, err := validateDate(tt.test)
		if stripDate != tt.stripDate && err != tt.err {
			t.Errorf("Expected %q, %v; got %q %v", tt.stripDate, tt.err, stripDate, err)
		}
	}
}

func Test_getStripPage(t *testing.T) {
	var getStripTests = []struct {
		test string
		err  string
	}{
		{"kcuc", "Error loading page http://dilbert.com/strip/kcuc"},
		{"1983-02-28", "Error loading page http://dilbert.com/strip/1983-02-28"},
	}

	for _, tt := range getStripTests {
		_, err := getStripPage(tt.test)
		if err.Error() != tt.err {
			t.Errorf("Expected %q; got %q", tt.err, err.Error())
		}
	}
}
