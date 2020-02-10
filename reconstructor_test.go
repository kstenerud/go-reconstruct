package reconstruct

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/kstenerud/go-equivalence"
)

func generateString(charCount int) string {
	charRange := int('z' - 'a')
	var result strings.Builder
	for i := 0; i < charCount; i++ {
		ch := 'a' + (i+charCount)%charRange
		result.WriteByte(byte(ch))
	}
	return result.String()
}

func generateBytes(length int) []byte {
	return []byte(generateString(length))
}

type InnerStruct struct {
	Inner int
}

type TestStruct struct {
	Bo    bool
	By    byte
	I     int
	I8    int8
	I16   int16
	I32   int32
	I64   int64
	U     uint
	U8    uint8
	U16   uint16
	U32   uint32
	U64   uint64
	F32   float32
	F64   float64
	Ar    [4]byte
	St    string
	Ba    []byte
	Sl    []interface{}
	M     map[interface{}]interface{}
	Pi    *int
	IS    InnerStruct
	ISP   *InnerStruct
	Time  time.Time
	URL   url.URL
	PTime *time.Time
	PURL  *url.URL
}

func newTestStruct(baseValue int) *TestStruct {
	this := new(TestStruct)
	this.Init(baseValue)
	return this
}

func (this *TestStruct) Init(baseValue int) {
	this.Bo = baseValue&1 == 1
	this.By = byte(baseValue + 1)
	this.I = baseValue + 2
	this.I8 = int8(baseValue + 3)
	this.I16 = int16(baseValue + 4)
	this.I32 = int32(baseValue + 5)
	this.I64 = int64(baseValue + 6)
	this.U = uint(baseValue + 7)
	this.U8 = uint8(baseValue + 8)
	this.U16 = uint16(baseValue + 9)
	this.U32 = uint32(baseValue + 10)
	this.U64 = uint64(baseValue + 11)
	this.F32 = float32(baseValue) + 20.5
	this.F64 = float64(baseValue) + 21.5
	this.Ar[0] = byte(baseValue + 30)
	this.Ar[1] = byte(baseValue + 31)
	this.Ar[2] = byte(baseValue + 32)
	this.Ar[3] = byte(baseValue + 33)
	this.St = generateString(baseValue + 5)
	this.Ba = generateBytes(baseValue + 1)
	this.M = make(map[interface{}]interface{})
	for i := 0; i < baseValue+2; i++ {
		this.Sl = append(this.Sl, i)
		this.M[fmt.Sprintf("key%v", i)] = i
	}
	v := baseValue
	this.Pi = &v
	this.IS.Inner = baseValue + 15
	this.ISP = new(InnerStruct)
	this.ISP.Inner = baseValue + 16

	testTime := time.Date(2000+baseValue, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	this.PTime = &testTime
	this.PURL, _ = url.Parse(fmt.Sprintf("http://example.com/%v", baseValue))
}

func InitMap(m map[interface{}]interface{}, baseValue int) {
	m["Bo"] = baseValue&1 == 1
	m["By"] = byte(baseValue + 1)
	m["I"] = baseValue + 2
	m["I8"] = int8(baseValue + 3)
	m["I16"] = int16(baseValue + 4)
	m["I32"] = int32(baseValue + 5)
	m["I64"] = int64(baseValue + 6)
	m["U"] = uint(baseValue + 7)
	m["U8"] = uint8(baseValue + 8)
	m["U16"] = uint16(baseValue + 9)
	m["U32"] = uint32(baseValue + 10)
	m["U64"] = uint64(baseValue + 11)
	m["F32"] = float32(baseValue) + 20.5
	m["F64"] = float64(baseValue) + 21.5
	m["Ar"] = []byte{byte(baseValue + 30), byte(baseValue + 31), byte(baseValue + 32), byte(baseValue + 33)}
	m["St"] = generateString(baseValue + 5)
	m["Ba"] = generateBytes(baseValue + 1)
	var s []interface{}
	mm := make(map[interface{}]interface{})
	for i := 0; i < baseValue+2; i++ {
		s = append(s, i)
		mm[fmt.Sprintf("key%v", i)] = i
	}
	m["Sl"] = s
	m["M"] = mm
	v := baseValue
	m["Pi"] = &v
	is := make(map[interface{}]interface{})
	is["Inner"] = baseValue + 15
	isp := make(map[interface{}]interface{})
	isp["Inner"] = baseValue + 16
	m["IS"] = is
	m["ISP"] = &isp
	var Time time.Time
	m["Time"] = Time
	var URL url.URL
	m["URL"] = URL
	testTime := time.Date(2000+baseValue, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	m["PTime"] = &testTime
	m["PURL"], _ = url.Parse(fmt.Sprintf("http://example.com/%v", baseValue))
}

func newTestMap(baseValue int) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	InitMap(m, baseValue)
	return m
}

func expectReconstructed(t *testing.T, s interface{}, d interface{}, expected interface{}) {
	if err := Reconstruct(s, d); err != nil {
		t.Error(err)
		return
	}
	re := reflect.ValueOf(expected)
	for re.Kind() == reflect.Interface || re.Kind() == reflect.Ptr {
		re = re.Elem()
	}
	realExpected := re.Interface()
	rd := reflect.ValueOf(d)
	for rd.Kind() == reflect.Interface || rd.Kind() == reflect.Ptr {
		rd = rd.Elem()
	}
	realD := rd.Interface()
	if !equivalence.IsEquivalent(realD, realExpected) {
		t.Errorf("Expected reconstructed %v but got %v", expected, realD)
	}
}

func assertReconstructionFails(t *testing.T, s interface{}, d interface{}) {
	if err := Reconstruct(s, d); err == nil {
		t.Errorf("Expected reconstruction to fail")
	}
}

func TestReconstructBool(t *testing.T) {
	s := true
	d := false
	expected := true
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructString(t *testing.T) {
	s := "a"
	d := "b"
	expected := "a"
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructBytes(t *testing.T) {
	s := []byte{1, 2, 3}
	d := []byte{4, 5, 6, 7, 8}
	expected := []byte{1, 2, 3}
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructTime(t *testing.T) {
	s := time.Date(2000, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	d := time.Date(2001, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	expected := time.Date(2000, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructURL(t *testing.T) {
	var s *url.URL
	var d *url.URL
	var expected *url.URL
	var err error
	if s, err = url.Parse("https://example1.com"); err != nil {
		panic(err)
	}
	if d, err = url.Parse("https://example2.com"); err != nil {
		panic(err)
	}
	if expected, err = url.Parse("https://example1.com"); err != nil {
		panic(err)
	}
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructInt8FromX(t *testing.T) {
	d := int8(1)
	expected := int8(15)

	sInt8 := int8(15)
	d = 1
	expectReconstructed(t, sInt8, &d, expected)

	sInt16 := int16(15)
	d = 1
	expectReconstructed(t, sInt16, &d, expected)

	sInt32 := int32(15)
	d = 1
	expectReconstructed(t, sInt32, &d, expected)

	sInt64 := int64(15)
	d = 1
	expectReconstructed(t, sInt64, &d, expected)

	sInt := int(15)
	d = 1
	expectReconstructed(t, sInt, &d, expected)

	sUint8 := uint8(15)
	d = 1
	expectReconstructed(t, sUint8, &d, expected)

	sUint16 := uint16(15)
	d = 1
	expectReconstructed(t, sUint16, &d, expected)

	sUint32 := uint32(15)
	d = 1
	expectReconstructed(t, sUint32, &d, expected)

	sUint64 := uint64(15)
	d = 1
	expectReconstructed(t, sUint64, &d, expected)

	sUint := uint(15)
	d = 1
	expectReconstructed(t, sUint, &d, expected)

	sFloat32 := float32(15)
	d = 1
	expectReconstructed(t, sFloat32, &d, expected)

	sFloat64 := float64(15)
	d = 1
	expectReconstructed(t, sFloat64, &d, expected)
}

func TestReconstructUint8FromX(t *testing.T) {
	d := uint8(1)
	expected := uint8(15)

	sInt8 := int8(15)
	d = 1
	expectReconstructed(t, sInt8, &d, expected)

	sInt16 := int16(15)
	d = 1
	expectReconstructed(t, sInt16, &d, expected)

	sInt32 := int32(15)
	d = 1
	expectReconstructed(t, sInt32, &d, expected)

	sInt64 := int64(15)
	d = 1
	expectReconstructed(t, sInt64, &d, expected)

	sInt := int(15)
	d = 1
	expectReconstructed(t, sInt, &d, expected)

	sUint8 := uint8(15)
	d = 1
	expectReconstructed(t, sUint8, &d, expected)

	sUint16 := uint16(15)
	d = 1
	expectReconstructed(t, sUint16, &d, expected)

	sUint32 := uint32(15)
	d = 1
	expectReconstructed(t, sUint32, &d, expected)

	sUint64 := uint64(15)
	d = 1
	expectReconstructed(t, sUint64, &d, expected)

	sUint := uint(15)
	d = 1
	expectReconstructed(t, sUint, &d, expected)

	sFloat32 := float32(15)
	d = 1
	expectReconstructed(t, sFloat32, &d, expected)

	sFloat64 := float64(15)
	d = 1
	expectReconstructed(t, sFloat64, &d, expected)
}

func TestReconstructFloat32FromX(t *testing.T) {
	d := float32(1)
	expected := float32(15)

	sInt8 := int8(15)
	d = 1
	expectReconstructed(t, sInt8, &d, expected)

	sInt16 := int16(15)
	d = 1
	expectReconstructed(t, sInt16, &d, expected)

	sInt32 := int32(15)
	d = 1
	expectReconstructed(t, sInt32, &d, expected)

	sInt64 := int64(15)
	d = 1
	expectReconstructed(t, sInt64, &d, expected)

	sInt := int(15)
	d = 1
	expectReconstructed(t, sInt, &d, expected)

	sUint8 := uint8(15)
	d = 1
	expectReconstructed(t, sUint8, &d, expected)

	sUint16 := uint16(15)
	d = 1
	expectReconstructed(t, sUint16, &d, expected)

	sUint32 := uint32(15)
	d = 1
	expectReconstructed(t, sUint32, &d, expected)

	sUint64 := uint64(15)
	d = 1
	expectReconstructed(t, sUint64, &d, expected)

	sUint := uint(15)
	d = 1
	expectReconstructed(t, sUint, &d, expected)

	sFloat32 := float32(15)
	d = 1
	expectReconstructed(t, sFloat32, &d, expected)

	sFloat64 := float64(15)
	d = 1
	expectReconstructed(t, sFloat64, &d, expected)
}

func TestReconstructIntXFromInt(t *testing.T) {
	sInt := int64(15)

	d8 := int8(1)
	expected8 := int8(15)
	expectReconstructed(t, sInt, &d8, expected8)

	d16 := int16(1)
	expected16 := int16(15)
	expectReconstructed(t, sInt, &d16, expected16)

	d32 := int32(1)
	expected32 := int32(15)
	expectReconstructed(t, sInt, &d32, expected32)

	d64 := int64(1)
	expected64 := int64(15)
	expectReconstructed(t, sInt, &d64, expected64)

	d := int(1)
	expected := int(15)
	expectReconstructed(t, sInt, &d, expected)
}

func TestReconstructUintXFromInt(t *testing.T) {
	sInt := int64(15)

	d8 := uint8(1)
	expected8 := uint8(15)
	expectReconstructed(t, sInt, &d8, expected8)

	d16 := uint16(1)
	expected16 := uint16(15)
	expectReconstructed(t, sInt, &d16, expected16)

	d32 := uint32(1)
	expected32 := uint32(15)
	expectReconstructed(t, sInt, &d32, expected32)

	d64 := uint64(1)
	expected64 := uint64(15)
	expectReconstructed(t, sInt, &d64, expected64)

	d := uint(1)
	expected := uint(15)
	expectReconstructed(t, sInt, &d, expected)
}

func TestReconstructFloatXFromInt(t *testing.T) {
	sInt := int64(15)

	d32 := float32(1)
	expected32 := float32(15)
	expectReconstructed(t, sInt, &d32, expected32)

	d64 := float64(1)
	expected64 := float64(15)
	expectReconstructed(t, sInt, &d64, expected64)
}

func TestReconstructSlice(t *testing.T) {
	src := []interface{}{1, 2, 3}
	dst := []int8{0}
	expected := []int8{1, 2, 3}
	expectReconstructed(t, src, &dst, expected)
}

func TestReconstructMap(t *testing.T) {
	src := map[interface{}]interface{}{1: 1, 2: 2, 3: 3}
	dst := map[interface{}]interface{}{"a": "abc"}
	expected := map[interface{}]interface{}{1: 1, 2: 2, 3: 3}
	expectReconstructed(t, src, &dst, expected)
}

func TestReconstructArray(t *testing.T) {
	src := [4]int{1, 2, 3, 4}
	dst := [4]int{5, 6, 7, 8}
	expected := [4]int{1, 2, 3, 4}
	expectReconstructed(t, src, &dst, expected)
}

func TestReconstructStruct(t *testing.T) {
	src := newTestMap(2)
	dst := newTestStruct(1)
	expected := newTestStruct(2)
	expectReconstructed(t, src, dst, expected)
}

func TestReconstructStructUninitialized(t *testing.T) {
	src := newTestMap(2)
	dst := new(TestStruct)
	expected := newTestStruct(2)
	expectReconstructed(t, src, dst, expected)
}

type TimeStruct struct {
	PTime *time.Time
}

func TestReconstructPTime(t *testing.T) {
	t1 := time.Date(2000, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	s := map[interface{}]interface{}{"PTime": &t1}
	d := TimeStruct{}
	expected := TimeStruct{&t1}
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructIntPInt(t *testing.T) {
	s := 20
	var d *int
	expected := 20
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructUintPInt(t *testing.T) {
	s := uint8(20)
	var d *int
	expected := 20
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructFloatPInt(t *testing.T) {
	s := 20.0
	var d *int
	expected := 20
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructUIntPUInt(t *testing.T) {
	s := uint16(20)
	var d *uint
	expected := 20
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructIntPUInt(t *testing.T) {
	s := int32(20)
	var d *uint
	expected := 20
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructFloatPUInt(t *testing.T) {
	s := float32(20)
	var d *uint
	expected := 20
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructFloatPFloat(t *testing.T) {
	s := 1.2
	var d *float64
	expected := 1.2
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructIntPFloat(t *testing.T) {
	s := 1
	var d *float32
	expected := 1
	expectReconstructed(t, s, &d, expected)
}

func TestReconstructUIntPFloat(t *testing.T) {
	s := uint(1)
	var d *float64
	expected := 1
	expectReconstructed(t, s, &d, expected)
}
