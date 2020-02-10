package reconstruct

import (
	"testing"
	"time"

	"github.com/kstenerud/go-equivalence"
)

func expectAdhocObject(t *testing.T, function func(*AdhocBuilder), expected interface{}) {
	builder := new(AdhocBuilder)
	function(builder)
	actual := builder.GetObject()
	if !equivalence.IsEquivalent(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestAdhocNil(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnNil()
	}, nil)
}

func TestAdhocBool(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnBool(true)
	}, true)
}

func TestAdhocInt(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnInt(-1)
	}, -1)
}

func TestAdhocUint(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnUint(1)
	}, 1)
}

func TestAdhocFloat(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnFloat(5.5)
	}, 5.5)
}

func TestAdhocString(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnString("test")
	}, "test")
}

func TestAdhocBytes(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnBytes([]byte{1, 2, 3})
	}, []byte{1, 2, 3})
}

func TestAdhocURI(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnURI(newURI("http://example.com"))
	}, newURI("http://example.com"))
}

func TestAdhocTime(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnTime(time.Date(2000, time.Month(1), 1, 1, 1, 1, 0, time.UTC))
	}, time.Date(2000, time.Month(1), 1, 1, 1, 1, 0, time.UTC))
}

func TestAdhocListEmpty(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnListBegin()
		builder.OnListEnd()
	}, []interface{}{})
}

func TestAdhocListListEmpty(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnListBegin()
		builder.OnListBegin()
		builder.OnListEnd()
		builder.OnListEnd()
	}, []interface{}{[]interface{}{}})
}

func TestAdhocList(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnListBegin()
		builder.OnBool(false)
		builder.OnInt(1)
		builder.OnString("blah")
		builder.OnListEnd()
	}, []interface{}{false, 1, "blah"})
}

func TestAdhocMapEmpty(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnMapBegin()
		builder.OnMapEnd()
	}, map[interface{}]interface{}{})
}

func TestAdhocMapMapEmpty(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnMapBegin()
		builder.OnBool(true)
		builder.OnMapBegin()
		builder.OnMapEnd()
		builder.OnMapEnd()
	}, map[interface{}]interface{}{true: map[interface{}]interface{}{}})
}

func TestAdhocMap(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnMapBegin()
		builder.OnString("key1")
		builder.OnInt(-50)
		builder.OnFloat(1.5)
		builder.OnString("value")
		builder.OnMapEnd()
	}, map[interface{}]interface{}{"key1": -50, 1.5: "value"})
}

func TestAdhocListMap(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnListBegin()
		builder.OnBool(false)
		builder.OnInt(1)
		builder.OnString("blah")

		builder.OnMapBegin()
		builder.OnString("key1")
		builder.OnInt(-50)
		builder.OnFloat(1.5)
		builder.OnString("value")
		builder.OnMapEnd()

		builder.OnBool(true)
		builder.OnListEnd()
	}, []interface{}{
		false,
		1,
		"blah",
		map[interface{}]interface{}{"key1": -50, 1.5: "value"},
		true,
	})
}

func TestAdhocMapList(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnMapBegin()
		builder.OnString("key1")
		builder.OnInt(-50)
		builder.OnFloat(1.5)
		builder.OnString("value")
		builder.OnString("a list")

		builder.OnListBegin()
		builder.OnFloat(1.1)
		builder.OnBool(false)
		builder.OnListEnd()

		builder.OnInt(100)
		builder.OnInt(200)
		builder.OnMapEnd()
	}, map[interface{}]interface{}{
		"key1": -50,
		1.5:    "value",
		"a list": []interface{}{
			1.1,
			false,
		},
		100: 200,
	})
}

func TestAdhocMapMap(t *testing.T) {
	expectAdhocObject(t, func(builder *AdhocBuilder) {
		builder.OnMapBegin()
		builder.OnString("key1")
		builder.OnInt(-50)

		builder.OnString("the map")
		builder.OnMapBegin()
		builder.OnFloat(1.5)
		builder.OnString("value")
		builder.OnMapEnd()

		builder.OnBool(true)
		builder.OnBool(false)

		builder.OnMapEnd()
	}, map[interface{}]interface{}{
		"key1": -50,
		"the map": map[interface{}]interface{}{
			1.5: "value",
		},
		true: false,
	})
}
