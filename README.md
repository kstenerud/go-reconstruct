Reconstruct
===========

Package reconstruct provides basic tools for deconstructing and reconstructing
data to/from go objects (structs, lists, maps, strings, scalars, etc). It's
expected to be used in tandem with other packages that provide serialization
or generation of data structures in a neutral format (interface slices and
interface maps). This package provides the "last mile" to deconstruct or
reconstruct objects based on data in this neutral format.

Neutral format data is data consisting of only the following types, and
combinations thereof:

- bool
- int (any size)
- uint (any size)
- float (any size)
- string
- []byte
- time.Time
- *url.URL
- []interface{}
- map[interface{}]interface{}

These provide the fundamental types required for deconstruction and
reconstruction of all types and structs in go.


Usage
-----

```golang
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

	// Prints:
	// Deconstructing &{Example 50 map[a:{0.5} b:{0.25} c:{0.75}]}
	// Resulting ad-hoc object: map[InnerStructsByName:map[a:map[Proportion:0.5] b:map[Proportion:0.25] c:map[Proportion:0.75]] Name:Example Number:50]
	// Reconstructed object: &{Example 50 map[a:{0.5} b:{0.25} c:{0.75}]}
}
```


License
-------

MIT License:

Copyright 2020 Karl Stenerud

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.