package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
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
		err  error
	}{
		{"kcuc", err404},
		{"1983-02-28", err404},
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
		test     string
		checksum string
		err      error
	}{
		{"http://assets.amuniversal.com/ecf9a570ae6b01341f1d005056a9545d", "2d9cc2fdd9dbc7b5f864ccf59c278aae", nil},
		{"http://google.com/ecf9a570ae6b01341f1d005056a9545", "d41d8cd98f00b204e9800998ecf8427e", err404},
	}
	for _, tt := range getStripImageTests {
		stripImage, err := getStripImage(tt.test)
		sum := md5.Sum(stripImage)
		base16 := hex.EncodeToString(sum[:])
		if err != tt.err || base16 != tt.checksum {
			t.Errorf("Expected %v, %v; got %v, %v", tt.checksum, tt.err, base16, err)
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
		stripDate  string
		err        error
	}{
		{[]byte("test content"), "2016-01-11", nil},
		{[]byte("test content"), "2016-01-11", os.ErrExist},
		{[]byte("test content"), "2016-01-32", errors.New("parsing time \"2016-01-32\": day out of range")},
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
