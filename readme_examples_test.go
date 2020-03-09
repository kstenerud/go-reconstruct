package reconstruct

import (
	"fmt"
	"testing"

	"github.com/kstenerud/go-describe"
)

type ExampleInnerStruct struct {
	Proportion float32
}

type ExampleStruct struct {
	Name               string
	Number             int
	InnerStructsByName map[string]ExampleInnerStruct
}

func Demonstrate() {
	value := new(ExampleStruct)
	value.Name = "Example"
	value.Number = 50
	value.InnerStructsByName = make(map[string]ExampleInnerStruct)
	value.InnerStructsByName["a"] = ExampleInnerStruct{0.5}
	value.InnerStructsByName["b"] = ExampleInnerStruct{0.25}
	value.InnerStructsByName["c"] = ExampleInnerStruct{0.75}

	fmt.Printf("Connecting iterator and builder to deconstruct and reconstruct %v\n", describe.Describe(value, 4))

	builder := NewBuilderFor(value)
	useReferences := false
	if err := IterateObject(value, useReferences, builder); err != nil {
		// TODO: Handle this
	}

	rebuilt := builder.GetBuiltObject()

	fmt.Printf("Reconstructed object: %v\n", describe.Describe(rebuilt, 4))
}

func TestReadmeExamples(t *testing.T) {
	Demonstrate()
}
