package main

import (
	"errors"
	"testing"
)

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
			t.Errorf("Expected %q; got %v", tt.err, err)
		}
	}
}

func Test_getStripImageAddr(t *testing.T) {
	var getStripImageAddrTests = []struct {
		test           string
		stripImageAddr string
		err            error
	}{
		{`data-url="http://dilbert.com/strip/2017-01-10" data-image="http://assets.amuniversal.com/ecf9a570ae6b01341f1d005056a9545d" data-date="January 10, 2017" `, "http://assets.amuniversal.com/ecf9a570ae6b01341f1d005056a9545d", nil},
		{`some other text`, "", errors.New("404 - Page Not Found")},
	}

	for _, tt := range getStripImageAddrTests {
		stripImageAddr, err := getStripImageAddr([]byte(tt.test))
		if stripImageAddr != tt.stripImageAddr && err != tt.err {
			t.Errorf("Expected %q, %v; got %q %v", tt.stripImageAddr, tt.err, stripImageAddr, err)
		}
	}
}
