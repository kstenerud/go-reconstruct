// +build disabled

package reconstruct

import (
	"fmt"
	"reflect"
	"testing"
)

func deconstructAndReconstruct(value interface{}) error {
	fmt.Printf("Deconstructing %v\n", value)
	// Build an ad-hoc object from value
	builder := new(AdhocBuilder)
	iterator := NewObjectIterator(builder)
	if err := iterator.Iterate(value); err != nil {
		return err
	}
	adhocObject := builder.GetObject()
	fmt.Printf("Resulting ad-hoc object: %v\n", adhocObject)

	// Create a new pointer-to-object of type v and fill it from the ad-hoc object.
	// You could also use a concrete object.
	pointer := reflect.New(reflect.TypeOf(value))
	dst := pointer.Interface()
	if err := Reconstruct(adhocObject, dst); err != nil {
		return err
	}
	fmt.Printf("Reconstructed object: %v\n", pointer.Elem())
	return nil
}

type ExampleInnerStruct struct {
	Proportion float32
}

type ExampleStruct struct {
	Name               string
	Number             int
	InnerStructsByName map[string]ExampleInnerStruct
}

func TestReadmeExamples(t *testing.T) {
	value := new(ExampleStruct)
	value.Name = "Example"
	value.Number = 50
	value.InnerStructsByName = make(map[string]ExampleInnerStruct)
	value.InnerStructsByName["a"] = ExampleInnerStruct{0.5}
	value.InnerStructsByName["b"] = ExampleInnerStruct{0.25}
	value.InnerStructsByName["c"] = ExampleInnerStruct{0.75}
	if err := deconstructAndReconstruct(value); err != nil {
		t.Error(err)
	}
}
