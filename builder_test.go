// +build disabled

package reconstruct

import (
	"net/url"
	"testing"
	"time"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func n() func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnNil() }
}
func b(value bool) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnBool(value) }
}
func i(value int64) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnInt(value) }
}
func u(value uint64) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnUint(value) }
}
func f(value float64) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnFloat(value) }
}
func s(value string) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnString(value) }
}
func bin(value []byte) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnBytes(value) }
}
func uri(value string) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnURI(newURI(value)) }
}
func tm(value time.Time) func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnTime(value) }
}
func l() func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnListBegin() }
}
func m() func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnMapBegin() }
}
func e() func(builder *RootBuilder) {
	return func(builder *RootBuilder) { builder.OnContainerEnd() }
}

// TODO: Marker, reference

func runBuildCmds(builder *RootBuilder, commands ...func(*RootBuilder)) {
	for _, cmd := range commands {
		cmd(builder)
	}
}

func runBuild(expected interface{}, commands ...func(*RootBuilder)) interface{} {
	builder := NewBuilderFor(expected)
	runBuildCmds(builder, commands...)
	return builder.GetBuiltObject()
}

func assertBuild(t *testing.T, expected interface{}, commands ...func(*RootBuilder)) {
	actual := runBuild(expected, commands...)
	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func assertBuildPanics(t *testing.T, template interface{}, commands ...func(*RootBuilder)) {
	assertPanics(t, func() {
		runBuild(template, commands...)
	})
}

func TestBuilderBasic(t *testing.T) {
	assertBuild(t, true, b(true))
	assertBuild(t, int(1), i(1))
	assertBuild(t, int8(1), i(1))
	assertBuild(t, int16(1), i(1))
	assertBuild(t, int32(1), i(1))
	assertBuild(t, int64(1), i(1))
	assertBuild(t, uint(1), u(1))
	assertBuild(t, uint8(1), u(1))
	assertBuild(t, uint16(1), u(1))
	assertBuild(t, uint32(1), u(1))
	assertBuild(t, uint64(1), u(1))
	assertBuild(t, float32(1.1), f(1.1))
	assertBuild(t, float64(1.1), f(1.1))
	assertBuild(t, "testing", s("testing"))
	assertBuild(t, []byte{1, 2, 3}, bin([]byte{1, 2, 3}))
	assertBuild(t, interface{}(1234), i(1234))
}

func TestBuilderTime(t *testing.T) {
	testTime := time.Now()
	assertBuild(t, testTime, tm(testTime))
}

func TestBuilderURI(t *testing.T) {
	testURI := "https://x.com"
	assertBuild(t, newURI(testURI), uri(testURI))
}

func TestBuilderBasicTypeFail(t *testing.T) {
	assertBuildPanics(t, true, n())
	assertBuildPanics(t, true, i(1))
	assertBuildPanics(t, true, u(1))
	assertBuildPanics(t, true, f(1))
	assertBuildPanics(t, true, s("1"))
	assertBuildPanics(t, true, bin([]byte{1}))
	assertBuildPanics(t, true, uri("x://x"))
	assertBuildPanics(t, true, tm(time.Now()))
	assertBuildPanics(t, true, l())
	assertBuildPanics(t, true, m())
	assertBuildPanics(t, true, e())

	assertBuildPanics(t, int(1), n())
	assertBuildPanics(t, int(1), b(true))
	assertBuildPanics(t, int(1), s("1"))
	assertBuildPanics(t, int(1), bin([]byte{1}))
	assertBuildPanics(t, int(1), uri("x://x"))
	assertBuildPanics(t, int(1), tm(time.Now()))
	assertBuildPanics(t, int(1), l())
	assertBuildPanics(t, int(1), m())
	assertBuildPanics(t, int(1), e())

	assertBuildPanics(t, uint(1), n())
	assertBuildPanics(t, uint(1), b(true))
	assertBuildPanics(t, uint(1), s("1"))
	assertBuildPanics(t, uint(1), bin([]byte{1}))
	assertBuildPanics(t, uint(1), uri("x://x"))
	assertBuildPanics(t, uint(1), tm(time.Now()))
	assertBuildPanics(t, uint(1), l())
	assertBuildPanics(t, uint(1), m())
	assertBuildPanics(t, uint(1), e())

	assertBuildPanics(t, float64(1), n())
	assertBuildPanics(t, float64(1), b(true))
	assertBuildPanics(t, float64(1), s("1"))
	assertBuildPanics(t, float64(1), bin([]byte{1}))
	assertBuildPanics(t, float64(1), uri("x://x"))
	assertBuildPanics(t, float64(1), tm(time.Now()))
	assertBuildPanics(t, float64(1), l())
	assertBuildPanics(t, float64(1), m())
	assertBuildPanics(t, float64(1), e())

	// TODO: What should allow nil?
	// assertBuildPanics(t, "", n())
	assertBuildPanics(t, "", b(true))
	assertBuildPanics(t, "", i(1))
	assertBuildPanics(t, "", u(1))
	assertBuildPanics(t, "", f(1))
	assertBuildPanics(t, "", bin([]byte{1}))
	assertBuildPanics(t, "", uri("x://x"))
	assertBuildPanics(t, "", tm(time.Now()))
	assertBuildPanics(t, "", l())
	assertBuildPanics(t, "", m())
	assertBuildPanics(t, "", e())

	// TODO: What should allow nil?
	// assertBuildPanics(t, []byte{}, n())
	assertBuildPanics(t, []byte{}, b(true))
	assertBuildPanics(t, []byte{}, i(1))
	assertBuildPanics(t, []byte{}, u(1))
	assertBuildPanics(t, []byte{}, f(1))
	assertBuildPanics(t, []byte{}, s("1"))
	assertBuildPanics(t, []byte{}, uri("x://x"))
	assertBuildPanics(t, []byte{}, tm(time.Now()))
	assertBuildPanics(t, []byte{}, l())
	assertBuildPanics(t, []byte{}, m())
	assertBuildPanics(t, []byte{}, e())

	// TODO: What should allow nil?
	// assertBuildPanics(t, newURI("x://x"), n())
	assertBuildPanics(t, newURI("x://x"), b(true))
	assertBuildPanics(t, newURI("x://x"), i(1))
	assertBuildPanics(t, newURI("x://x"), u(1))
	assertBuildPanics(t, newURI("x://x"), f(1))
	assertBuildPanics(t, newURI("x://x"), s("1"))
	assertBuildPanics(t, newURI("x://x"), bin([]byte{1}))
	assertBuildPanics(t, newURI("x://x"), tm(time.Now()))
	assertBuildPanics(t, newURI("x://x"), l())
	assertBuildPanics(t, newURI("x://x"), m())
	assertBuildPanics(t, newURI("x://x"), e())

	assertBuildPanics(t, time.Now(), n())
	assertBuildPanics(t, time.Now(), b(true))
	assertBuildPanics(t, time.Now(), i(1))
	assertBuildPanics(t, time.Now(), u(1))
	assertBuildPanics(t, time.Now(), f(1))
	assertBuildPanics(t, time.Now(), s("1"))
	assertBuildPanics(t, time.Now(), bin([]byte{1}))
	assertBuildPanics(t, time.Now(), uri("x://x"))
	assertBuildPanics(t, time.Now(), l())
	assertBuildPanics(t, time.Now(), m())
	assertBuildPanics(t, time.Now(), e())

	assertBuildPanics(t, []int{}, n())
	assertBuildPanics(t, []int{}, b(true))
	assertBuildPanics(t, []int{}, i(1))
	assertBuildPanics(t, []int{}, u(1))
	assertBuildPanics(t, []int{}, f(1))
	assertBuildPanics(t, []int{}, s("1"))
	assertBuildPanics(t, []int{}, bin([]byte{1}))
	assertBuildPanics(t, []int{}, uri("x://x"))
	assertBuildPanics(t, []int{}, tm(time.Now()))
	assertBuildPanics(t, []int{}, m())
	assertBuildPanics(t, []int{}, e())

	// TODO: Check if this is correct behavior to not panic
	// assertBuildPanics(t, map[int]int{}, n())
	assertBuildPanics(t, map[int]int{}, i(1))
	assertBuildPanics(t, map[int]int{}, u(1))
	assertBuildPanics(t, map[int]int{}, f(1))
	assertBuildPanics(t, map[int]int{}, s("1"))
	assertBuildPanics(t, map[int]int{}, bin([]byte{1}))
	assertBuildPanics(t, map[int]int{}, uri("x://x"))
	assertBuildPanics(t, map[int]int{}, tm(time.Now()))
	assertBuildPanics(t, map[int]int{}, l())
	assertBuildPanics(t, map[int]int{}, e())
}

func TestBuilderNumericConversion(t *testing.T) {
	assertBuild(t, int8(50), u(50))
	assertBuild(t, int16(50), f(50))
	assertBuild(t, uint32(50), i(50))
	assertBuild(t, uint64(50), f(50))
	assertBuild(t, float32(50), i(50))
	assertBuild(t, float64(50), u(50))
}

func TestBuilderNumericConversionFail(t *testing.T) {
	assertBuildPanics(t, int8(0), i(300))
	assertBuildPanics(t, int(0), f(3.5))
	assertBuildPanics(t, uint(0), i(-1))
	assertBuildPanics(t, uint(0), f(3.5))
	assertBuildPanics(t, float32(0), i(0x7fffffffffffffff))
	assertBuildPanics(t, float64(0), u(0xffffffffffffffff))
}

func TestBuilderSlice(t *testing.T) {
	assertBuild(t, []bool{false, true}, l(), b(false), b(true), e())
	assertBuild(t, []int8{1, 2, 3}, l(), i(1), u(2), f(3), e())
	assertBuild(t, []interface{}{false, 1, "test"}, l(), b(false), i(1), s("test"), e())
}

func TestBuilderArray(t *testing.T) {
	assertBuild(t, [2]bool{false, true}, l(), b(false), b(true), e())
	assertBuild(t, [3]int8{1, 2, 3}, l(), i(1), u(2), f(3), e())
	assertBuild(t, [3]interface{}{false, 1, "test"}, l(), b(false), i(1), s("test"), e())
}

func TestBuilderMap(t *testing.T) {
	assertBuild(t, map[string]bool{"true": true, "false": false}, m(), s("true"), b(true), s("false"), b(false), e())
	assertBuild(t, map[interface{}]interface{}{"false": false, 1: "one"}, m(), s("false"), b(false), i(1), s("one"), e())
}

func TestBuilderSliceSlice(t *testing.T) {
	assertBuild(t, [][]bool{{false, true}}, l(), l(), b(false), b(true), e(), e())
}

func TestBuilderMapMap(t *testing.T) {
	assertBuild(t, map[string]map[int]bool{"first": {1: true}}, m(), s("first"), m(), i(1), b(true), e(), e())
}

func TestBuilderSliceMap(t *testing.T) {
	assertBuild(t, []map[int]bool{{1: true}}, l(), m(), i(1), b(true), e(), e())
}

func TestBuilderMapSlice(t *testing.T) {
	assertBuild(t, map[string][]int{"first": {1}}, m(), s("first"), l(), i(1), e(), e())
}

type BuilderTestStruct struct {
	internal string
	ABool    bool
	AString  string
	AnInt    int
	AMap     map[int]int8
	ASlice   []string
}

func TestBuilderStruct(t *testing.T) {
	assertBuild(t,
		BuilderTestStruct{
			AString: "test",
			AnInt:   1,
			ABool:   true,
			AMap:    map[int]int8{1: 50},
			ASlice:  []string{"the slice"},
		},
		m(),
		s("AString"), s("test"),
		s("AMap"), m(), i(1), i(50), e(),
		s("AnInt"), i(1),
		s("ASlice"), l(), s("the slice"), e(),
		s("ABool"), b(true),
		e())
}

func TestBuilderStructIgnored(t *testing.T) {
	assertBuild(t, BuilderTestStruct{
		AString: "test",
		AnInt:   1,
		ABool:   true,
	}, m(), s("AString"), s("test"), s("Something"), i(5), s("AnInt"), i(1), s("ABool"), b(true), e())
}

func TestBuilderListStruct(t *testing.T) {
	assertBuild(t,
		[]BuilderTestStruct{
			BuilderTestStruct{
				AString: "test",
				AnInt:   1,
				ABool:   true,
				AMap:    map[int]int8{1: 50},
				ASlice:  []string{"the slice"},
			},
		},
		l(),
		m(),
		s("AString"), s("test"),
		s("AMap"), m(), i(1), i(50), e(),
		s("AnInt"), i(1),
		s("ASlice"), l(), s("the slice"), e(),
		s("ABool"), b(true),
		e(),
		e())
}

func TestBuilderMapStruct(t *testing.T) {
	assertBuild(t,
		map[string]BuilderTestStruct{
			"struct": BuilderTestStruct{
				AString: "test",
				AnInt:   1,
				ABool:   true,
				AMap:    map[int]int8{1: 50},
				ASlice:  []string{"the slice"},
			},
		},
		m(),
		s("struct"),
		m(),
		s("AString"), s("test"),
		s("AMap"), m(), i(1), i(50), e(),
		s("AnInt"), i(1),
		s("ASlice"), l(), s("the slice"), e(),
		s("ABool"), b(true),
		e(),
		e())
}

func TestBuilderMultipleComplexBuilds(t *testing.T) {
	v := BuilderTestStruct{
		AString: "test",
		AnInt:   1,
		ABool:   true,
		AMap:    map[int]int8{1: 50},
		ASlice:  []string{"the slice"},
	}

	for idx := 0; idx < 10; idx++ {
		assertBuild(t,
			v,
			m(),
			s("AString"), s("test"),
			s("AMap"), m(), i(1), i(50), e(),
			s("AnInt"), i(1),
			s("ASlice"), l(), s("the slice"), e(),
			s("ABool"), b(true),
			e())
	}
}

type BuilderPtrTestStruct struct {
	internal    string
	ABool       *bool
	AnInt       *int
	AnInt8      *int8
	AnInt16     *int16
	AnInt32     *int32
	AnInt64     *int64
	AUint       *uint
	AUint8      *uint8
	AUint16     *uint16
	AUint32     *uint32
	AUint64     *uint64
	AFloat32    *float32
	AFloat64    *float64
	AString     *string
	AnInterface *interface{}
}

func TestBuilderPtr(t *testing.T) {
	aBool := true
	anInt := int(1)
	anInt8 := int8(2)
	anInt16 := int16(3)
	anInt32 := int32(4)
	anInt64 := int64(5)
	aUint := uint(6)
	aUint8 := uint8(7)
	aUint16 := uint16(8)
	aUint32 := uint32(9)
	aUint64 := uint64(10)
	aFloat32 := float32(11.5)
	aFloat64 := float64(12.5)
	aString := "test"
	var anInterface interface{}
	anInterface = 100
	v := BuilderPtrTestStruct{
		ABool:       &aBool,
		AnInt:       &anInt,
		AnInt8:      &anInt8,
		AnInt16:     &anInt16,
		AnInt32:     &anInt32,
		AnInt64:     &anInt64,
		AUint:       &aUint,
		AUint8:      &aUint8,
		AUint16:     &aUint16,
		AUint32:     &aUint32,
		AUint64:     &aUint64,
		AFloat32:    &aFloat32,
		AFloat64:    &aFloat64,
		AString:     &aString,
		AnInterface: &anInterface,
	}
	assertBuild(t,
		v,
		m(),
		s("ABool"), b(true),
		s("AnInt"), i(1),
		s("AnInt8"), i(2),
		s("AnInt16"), i(3),
		s("AnInt32"), i(4),
		s("AnInt64"), i(5),
		s("AUint"), u(6),
		s("AUint8"), u(7),
		s("AUint16"), u(8),
		s("AUint32"), u(9),
		s("AUint64"), u(10),
		s("AFloat32"), f(11.5),
		s("AFloat64"), f(12.5),
		s("AString"), s("test"),
		s("AnInterface"), i(100),
		e())
}

type BuilderSliceTestStruct struct {
	internal    []string
	ABool       []bool
	AnInt       []int
	AnInt8      []int8
	AnInt16     []int16
	AnInt32     []int32
	AnInt64     []int64
	AUint       []uint
	AUint8      []uint8
	AUint16     []uint16
	AUint32     []uint32
	AUint64     []uint64
	AFloat32    []float32
	AFloat64    []float64
	AString     []string
	AnInterface []interface{}
}

func TestBuilderSliceOfStructs(t *testing.T) {
	v := []BuilderSliceTestStruct{
		BuilderSliceTestStruct{
			AnInt: []int{1},
		},
		BuilderSliceTestStruct{
			AnInt: []int{1},
		},
	}

	assertBuild(t,
		v,
		l(),
		m(),
		s("AnInt"), l(), i(1), e(),
		e(),
		m(),
		s("AnInt"), l(), i(1), e(),
		e(),
		e())
}

type SimpleTestStruct struct {
	IValue int
}

func TestListOfStruct(t *testing.T) {
	v := []*SimpleTestStruct{
		&SimpleTestStruct{
			IValue: 5,
		},
	}

	assertBuild(t,
		v,
		l(),
		m(),
		s("IValue"),
		i(5),
		e(),
		e())
}

type NilContainers struct {
	Bytes []byte
	Slice []interface{}
	Map   map[interface{}]interface{}
}

func TestNilContainers(t *testing.T) {
	v := NilContainers{}

	assertBuild(t, v,
		m(),
		s("Bytes"),
		n(),
		s("Slice"),
		n(),
		s("Map"),
		n(),
		e())
}

type PURLContainer struct {
	URL *url.URL
}

func TestPURLContainer(t *testing.T) {
	v := PURLContainer{newURI("http://x.com")}

	assertBuild(t, v,
		m(),
		s("URL"),
		uri("http://x.com"),
		e())
}

func TestNilPURLContainer(t *testing.T) {
	v := PURLContainer{}

	assertBuild(t, v,
		m(),
		s("URL"),
		n(),
		e())
}

func TestByteArrayBytes(t *testing.T) {
	assertBuild(t, [2]byte{1, 2},
		bin([]byte{1, 2}))
}

// TODO
// type SelfReferential struct {
// 	Self *SelfReferential
// }

// func TestSelfReferential(t *testing.T) {
// 	assertBuild(t, SelfReferential{},
// 		m(),
// 		s("Self"),
// 		n(),
// 		e())
// }
