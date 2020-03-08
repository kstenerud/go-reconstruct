// +build disabled

package reconstruct

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func assertIterateBuild(t *testing.T, expected interface{}) {
	builder := NewBuilderFor(expected)
	useReferences := true
	if err := IterateObject(expected, useReferences, builder); err != nil {
		t.Error(err)
		return
	}

	actual := builder.GetBuiltObject()
	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func TestIterateBuildBasic(t *testing.T) {
	assertIterateBuild(t, true)
	assertIterateBuild(t, int(1))
	assertIterateBuild(t, int8(1))
	assertIterateBuild(t, int16(1))
	assertIterateBuild(t, int32(1))
	assertIterateBuild(t, int64(1))
	assertIterateBuild(t, uint(1))
	assertIterateBuild(t, uint8(1))
	assertIterateBuild(t, uint16(1))
	assertIterateBuild(t, uint32(1))
	assertIterateBuild(t, uint64(1))
	assertIterateBuild(t, float32(1))
	assertIterateBuild(t, float64(1))
	assertIterateBuild(t, "test")
	assertIterateBuild(t, []int{1, 2, 3})
}

type IterateBuildInnerStruct struct {
	Inner int
}

type IterateBuildTester struct {
	unexported int
	Bo         bool
	By         byte
	I          int
	I8         int8
	I16        int16
	I32        int32
	I64        int64
	U          uint
	U8         uint8
	U16        uint16
	U32        uint32
	U64        uint64
	F32        float32
	F64        float64
	Ar         [4]byte
	St         string
	Ba         []byte
	Sl         []interface{}
	M          map[interface{}]interface{}
	Pi         *int
	IS         IterateBuildInnerStruct
	ISP        *IterateBuildInnerStruct
	Time       time.Time
	URL        url.URL
	PTime      *time.Time
	PURL       *url.URL
}

func newIterateBuildTestStruct(baseValue int) *IterateBuildTester {
	this := new(IterateBuildTester)
	this.Init(baseValue)
	return this
}

func (this *IterateBuildTester) Init(baseValue int) {
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
	this.St = generateString(baseValue+5, baseValue)
	this.Ba = generateBytes(baseValue+1, baseValue)
	this.M = make(map[interface{}]interface{})
	for i := 0; i < baseValue+2; i++ {
		this.Sl = append(this.Sl, i)
		this.M[fmt.Sprintf("key%v", i)] = i
	}
	v := baseValue
	this.Pi = &v
	this.IS.Inner = baseValue + 15
	this.ISP = new(IterateBuildInnerStruct)
	this.ISP.Inner = baseValue + 16

	testTime := time.Date(2000+baseValue, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	this.PTime = &testTime
	this.PURL, _ = url.Parse(fmt.Sprintf("http://example.com/%v", baseValue))
}

func TestIterateBuild(t *testing.T) {
	assertIterateBuild(t, *newIterateBuildTestStruct(1))
}

func TestIterateBuildEmpty(t *testing.T) {
	assertIterateBuild(t, *new(IterateBuildTester))
}
