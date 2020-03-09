Reconstruct
===========

Package reconstruct provides basic tools for deconstructing and reconstructing
data to/from go objects (structs, lists, maps, strings, scalars, etc). It's
expected to be used in tandem with other packages that provide serialization
or generation of data structures. This package provides the "last mile" to
deconstruct or reconstruct objects based on the following data events:

 * Nil
 * Bool
 * Int
 * Uint
 * Float
 * String
 * Bytes
 * URI
 * Time
 * List
 * Map
 * End
 * Marker
 * Reference


Usage
-----

```golang
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

	builder := reconstruct.NewBuilderFor(value)
	useReferences := false
	if err := reconstruct.IterateObject(value, useReferences, builder); err != nil {
		// TODO: Handle this
	}

	rebuilt := builder.GetBuiltObject()

	fmt.Printf("Reconstructed object: %v\n", describe.Describe(rebuilt, 4))
}
```

#### Output:

```
TODO
```


License
-------

MIT License:

Copyright 2020 Karl Stenerud

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
