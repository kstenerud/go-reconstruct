package reconstruct

import (
	"net/url"
	"reflect"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	timeType  = reflect.TypeOf(time.Time{})
	urlType   = reflect.TypeOf(url.URL{})
	pURLType  = reflect.TypeOf((*url.URL)(nil))
	bytesType = reflect.TypeOf([]uint8{})
)

func isFieldExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}
