package reconstruct

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func assertRoundTrip(t *testing.T, v interface{}) {
	// Build an ad-hoc object from v
	builder := new(AdhocBuilder)
	iterator := NewObjectIterator(builder)
	if err := iterator.Iterate(v); err != nil {
		t.Error(err)
		return
	}
	adhocObject := builder.GetObject()

	// Create a new pointer-to-object of type v and fill it from the ad-hoc object
	dstType := reflect.TypeOf(v)
	dst := reflect.New(dstType).Interface()
	if err := Reconstruct(adhocObject, dst); err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(v, dst) {
		t.Errorf("Not equal: %v VS %v", describe.Describe(v), describe.Describe(dst))
	}
}

func newURI(v string) *url.URL {
	uri, err := url.Parse(v)
	if err != nil {
		panic(err)
	}
	return uri
}

func TestRoundtripBasicTypes(t *testing.T) {
	assertRoundTrip(t, true)
	assertRoundTrip(t, int(-1))
	assertRoundTrip(t, int8(-1))
	assertRoundTrip(t, int16(-1))
	assertRoundTrip(t, int32(-1))
	assertRoundTrip(t, int64(1))
	assertRoundTrip(t, uint(1))
	assertRoundTrip(t, uint8(1))
	assertRoundTrip(t, uint16(1))
	assertRoundTrip(t, uint32(1))
	assertRoundTrip(t, uint64(1))
	assertRoundTrip(t, float32(1.5))
	assertRoundTrip(t, float64(-8.1))
}

type SmallStruct struct {
	Value int
}

type PointerStruct struct {
	PInt      *int
	PURL      *url.URL
	PWhatever *interface{}
	PStruct   *SmallStruct
}

func TestRoundtripStructWithPointers(t *testing.T) {
	i := 1
	u := newURI("http://x.com")
	var w interface{} = "test"
	s := SmallStruct{1}
	v := PointerStruct{&i, u, &w, &s}
	assertRoundTrip(t, v)
}

func TestRoundtripNil(t *testing.T) {
	assertRoundTrip(t, []interface{}{nil})
	assertRoundTrip(t, map[interface{}]interface{}{1: nil})
	assertRoundTrip(t, PointerStruct{})
}

func TestRoundtripTime(t *testing.T) {
	assertRoundTrip(t, time.Now())
}

func TestRoundtripURI(t *testing.T) {
	assertRoundTrip(t, newURI("http://example.com/blah?something=5"))
}

func TestRoundtripListsArraysSlices(t *testing.T) {
	assertRoundTrip(t, "testing")
	assertRoundTrip(t, []byte{1, 2, 3, 4})
	assertRoundTrip(t, []int{5, 6, 7, 8})
	assertRoundTrip(t, []string{"abc", "def"})
	assertRoundTrip(t, []interface{}{"abc", 5000})
}

func TestRoundtripListList(t *testing.T) {
	assertRoundTrip(t, []interface{}{
		1, []interface{}{
			2,
		},
	})
}

func TestRoundtripMaps(t *testing.T) {
	assertRoundTrip(t, map[string]int{"test": 1})
}

func TestRoundtripStructs(t *testing.T) {
	assertRoundTrip(t, time.Now())
	assertRoundTrip(t, newURI("http://nowhere.com"))
}

func TestRoundtripDeepMapList(t *testing.T) {
	assertRoundTrip(t, map[string]interface{}{
		"aaa":  1,
		"blah": map[interface{}]interface{}{1.5: "x"},
	})
}

func TestRoundtripDeepMapList2(t *testing.T) {
	assertRoundTrip(t, map[string]interface{}{
		"test": 1,
		"inner": []interface{}{
			1, 2, "blah", map[interface{}]interface{}{
				1.5: "x",
			},
		},
	})
}

func TestRoundtripTestStruct(t *testing.T) {
	assertRoundTrip(t, newTestStruct(1))
}
