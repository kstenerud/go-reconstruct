package reconstruct

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func newURI(uriString string) *url.URL {
	uri, err := url.Parse(uriString)
	if err != nil {
		fmt.Printf("ERROR ERROR ERROR BUG: Bad URL (%v): %v", uriString, err)
		panic(err)
	}
	return uri
}

func reportPanic(function func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	function()
	return
}

func assertNoPanic(t *testing.T, function func()) {
	if err := reportPanic(function); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func assertPanics(t *testing.T, function func()) {
	if err := reportPanic(function); err == nil {
		t.Errorf("Expected an error")
	}
}

func generateString(charCount int, startIndex int) string {
	charRange := int('z' - 'a')
	var object strings.Builder
	for i := 0; i < charCount; i++ {
		ch := 'a' + (i+charCount+startIndex)%charRange
		object.WriteByte(byte(ch))
	}
	return object.String()
}

func generateBytes(length int, startIndex int) []byte {
	return []byte(generateString(length, startIndex))
}
