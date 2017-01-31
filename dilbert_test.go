package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_validateDate(t *testing.T) {
	type testArgs []string

	var validateDateTests = []struct {
		test testArgs
		err  string
	}{
		{testArgs{"", "1"}, `parsing time "1" as "2006-01-02": cannot parse "1" as "2006"`},
		{testArgs{"", "2016-99-30"}, `parsing time "2016-99-30": month out of range`},
		{testArgs{"", "1983-02-29"}, `parsing time "1983-02-29": day out of range`},
		{testArgs{"", "2016-01-30"}, ""},
	}

	for _, tt := range validateDateTests {
		_, err := validateDate(tt.test)
		if err != nil {
			if err.Error() != tt.err {
				t.Errorf("Expected %v; got %v", tt.err, err.Error())
			}
		} else {
			if tt.err != "" {
				t.Errorf("Expected %v; got %v", tt.err, err)
			}
		}
	}
}

func Test_getStripPage(t *testing.T) {
	var getStripTests = []struct {
		test time.Time
		err  error
	}{
		{time.Date(1983, 02, 28, 0, 0, 0, 0, time.UTC), err404},
		{time.Date(2009, 04, 29, 0, 0, 0, 0, time.UTC), nil},
	}

	for _, tt := range getStripTests {
		_, err := getStripPage(tt.test)
		if err != tt.err {
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
		{`some other text`, "", errFailedToFindImageAddr},
	}

	for _, tt := range getStripImageAddrTests {
		stripImageAddr, err := getStripImageAddr([]byte(tt.test))
		if stripImageAddr != tt.stripImageAddr && err != tt.err {
			t.Errorf("Expected %q, %v; got %q %v", tt.stripImageAddr, tt.err, stripImageAddr, err)
		}
	}
}

func Test_getStripImage(t *testing.T) {
	var getStripImageTests = []struct {
		test string
		err  error
	}{
		{"http://assets.amuniversal.com/ecf9a570ae6b01341f1d005056a9545d", nil},
		{"http://google.com/ecf9a570ae6b01341f1d005056a9545", err404},
	}
	for _, tt := range getStripImageTests {
		_, err := getStripImage(tt.test)
		if err != tt.err {
			t.Errorf("Expected %v; got %v", tt.err, err)
		}
	}
}

func Test_saveStripImage(t *testing.T) {
	testDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testDir)

	var saveStripImageTests = []struct {
		stripImage []byte
		stripDate  time.Time
		err        error
	}{
		{[]byte("test content"), time.Date(2016, 01, 11, 0, 0, 0, 0, time.UTC), nil},
		{[]byte("test content"), time.Date(2016, 01, 11, 0, 0, 0, 0, time.UTC), os.ErrExist},
	}
	for _, tt := range saveStripImageTests {
		err := saveStripImage(testDir, tt.stripImage, tt.stripDate)
		if err == nil && tt.err != err {
			t.Errorf("Expected '%v'; got '%v'", tt.err, err)
		}
		if err != nil && tt.err.Error() != err.Error() {
			t.Errorf("Expected '%v'; got '%v'", tt.err.Error(), err.Error())
		}
	}
}
