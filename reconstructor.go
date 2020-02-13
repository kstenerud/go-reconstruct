// Package reconstruct provides basic tools for deconstructing and reconstructing
// data to/from go objects (structs, lists, maps, strings, scalars, etc). It's
// expected to be used in tandem with other packages that provide serialization
// or generation of data structures in a neutral format (interface slices and
// interface maps). This package provides the "last mile" to deconstruct or
// reconstruct objects based on data in this neutral format.
//
// Neutral format data is data consisting of only the following types, and
// combinations thereof:
// - bool
// - int (any size)
// - uint (any size)
// - float (any size)
// - string
// - []byte
// - time.Time
// - *url.URL
// - []interface{}
// - map[interface{}]interface{}
//
// These provide the fundamental types required for deconstruction and
// reconstruction of all types and structs in go.
package reconstruct

import (
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type baseType uint

const (
	baseTypeNone = iota
	baseTypeInt
	baseTypeUint
	baseTypeFloat
	baseTypeComplex
	baseTypeIndirect
)

type objectReconstructor func(reflect.Value, reflect.Value)

type structElementReconstructor struct {
	elementIndex       int
	mapKey             reflect.Value
	reconstructElement objectReconstructor
}

var (
	baseTypes           [50]baseType
	knownReconstructors map[reflect.Type]objectReconstructor
)

func init() {
	baseTypes[baseType(reflect.Int)] = baseTypeInt
	baseTypes[baseType(reflect.Int8)] = baseTypeInt
	baseTypes[baseType(reflect.Int16)] = baseTypeInt
	baseTypes[baseType(reflect.Int32)] = baseTypeInt
	baseTypes[baseType(reflect.Int64)] = baseTypeInt
	baseTypes[baseType(reflect.Uint)] = baseTypeUint
	baseTypes[baseType(reflect.Uint8)] = baseTypeUint
	baseTypes[baseType(reflect.Uint16)] = baseTypeUint
	baseTypes[baseType(reflect.Uint32)] = baseTypeUint
	baseTypes[baseType(reflect.Uint64)] = baseTypeUint
	baseTypes[baseType(reflect.Float32)] = baseTypeFloat
	baseTypes[baseType(reflect.Float64)] = baseTypeFloat
	baseTypes[baseType(reflect.Complex64)] = baseTypeComplex
	baseTypes[baseType(reflect.Complex128)] = baseTypeComplex
	baseTypes[baseType(reflect.Interface)] = baseTypeIndirect
	baseTypes[baseType(reflect.Ptr)] = baseTypeIndirect

	var interfaceForType interface{}

	knownReconstructors = make(map[reflect.Type]objectReconstructor)
	knownReconstructors[reflect.TypeOf(int(0))] = reconstructInt
	knownReconstructors[reflect.TypeOf(int8(0))] = reconstructInt
	knownReconstructors[reflect.TypeOf(int16(0))] = reconstructInt
	knownReconstructors[reflect.TypeOf(int32(0))] = reconstructInt
	knownReconstructors[reflect.TypeOf(int64(0))] = reconstructInt
	knownReconstructors[reflect.TypeOf(uint(0))] = reconstructUint
	knownReconstructors[reflect.TypeOf(uint8(0))] = reconstructUint
	knownReconstructors[reflect.TypeOf(uint16(0))] = reconstructUint
	knownReconstructors[reflect.TypeOf(uint32(0))] = reconstructUint
	knownReconstructors[reflect.TypeOf(uint64(0))] = reconstructUint
	knownReconstructors[reflect.TypeOf(float32(0))] = reconstructFloat
	knownReconstructors[reflect.TypeOf(float64(0))] = reconstructFloat
	knownReconstructors[reflect.TypeOf(true)] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf(complex64(complex(1, 1)))] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf(complex128(complex(1, 1)))] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf(time.Time{})] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf("")] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf([]byte{})] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf(url.URL{})] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf(interfaceForType)] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf([]interface{}{})] = reconstructIdenticalKind
	knownReconstructors[reflect.TypeOf(map[interface{}]interface{}{})] = reconstructIdenticalKind
}

func getTypeIndirections(t reflect.Type) (indirections []reflect.Type) {
	pointer := t
	halfSpeedPointer := t
	moveHalfSpeed := false
	indirections = append(indirections, pointer)
	for pointer.Kind() == reflect.Ptr {
		pointer = pointer.Elem()
		if moveHalfSpeed {
			halfSpeedPointer = halfSpeedPointer.Elem()
		}
		if pointer == halfSpeedPointer {
			panic(fmt.Errorf("%v: Cannot handle type with looping indirections", t))
		}
		indirections = append(indirections, pointer)
		moveHalfSpeed = !moveHalfSpeed
	}
	return
}

func getConcreteObject(v reflect.Value) reflect.Value {
	for baseTypes[v.Kind()] == baseTypeIndirect {
		elem := v.Elem()
		if !elem.IsValid() {
			return v
		}
		v = elem
	}
	return v
}

func getObjectOrNewInstance(object reflect.Value) reflect.Value {
	if object.Kind() == reflect.Ptr {
		ptr := reflect.New(object.Type().Elem())
		object.Set(ptr)
		return object.Elem()
	}
	return object
}

func fillInt(value int64, dst reflect.Value) int64 {
	dst.SetInt(value)
	return dst.Int()
}

func fillUint(value uint64, dst reflect.Value) uint64 {
	dst.SetUint(value)
	return dst.Uint()
}

func fillFloat(value float64, dst reflect.Value) float64 {
	dst.SetFloat(value)
	return dst.Float()
}

// For bool, string, []byte, uri, time
func reconstructIdenticalKind(s reflect.Value, d reflect.Value) {
	dConcrete := getConcreteObject(d)
	if dConcrete.Kind() != s.Kind() {
		if dConcrete.Kind() != reflect.Ptr {
			panic(fmt.Errorf("%v doesn't have expected kind %v (kind was %v)", s, dConcrete.Kind(), s.Kind()))
		}
		if s.CanAddr() {
			// Copy pointer instead
			s = s.Addr()
		} else {
			// d is a nil pointer, so make a new object to reassign
			v := reflect.New(dConcrete.Type().Elem())
			v.Elem().Set(s)
			s = v
		}
	}
	dConcrete.Set(s)
	return
}

func reconstructFloat(s reflect.Value, d reflect.Value) {
	dConcrete := getObjectOrNewInstance(getConcreteObject(d))
	sflag := baseTypes[s.Kind()]
	switch sflag {
	case baseTypeFloat:
		v := s.Float()
		if fillFloat(v, dConcrete) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	case baseTypeInt:
		v := s.Int()
		if int64(fillFloat(float64(v), dConcrete)) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	case baseTypeUint:
		v := s.Uint()
		if uint64(fillFloat(float64(v), dConcrete)) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	}
	panic(fmt.Errorf("Type %v is not compatible with type %v", s.Type(), dConcrete.Type()))
}

func reconstructInt(s reflect.Value, d reflect.Value) {
	dConcrete := getObjectOrNewInstance(getConcreteObject(d))
	sflag := baseTypes[s.Kind()]
	switch sflag {
	case baseTypeFloat:
		v := s.Float()
		if float64(fillInt(int64(v), dConcrete)) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	case baseTypeInt:
		v := s.Int()
		if fillInt(v, dConcrete) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	case baseTypeUint:
		// Note: Not checking for negative bit
		v := s.Uint()
		if uint64(fillInt(int64(v), dConcrete)) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	}
	panic(fmt.Errorf("Type %v is not compatible with type %v", s.Type(), dConcrete.Type()))
}

func reconstructUint(s reflect.Value, d reflect.Value) {
	dConcrete := getObjectOrNewInstance(getConcreteObject(d))
	sflag := baseTypes[s.Kind()]
	switch sflag {
	case baseTypeFloat:
		v := s.Float()
		if float64(fillUint(uint64(v), dConcrete)) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	case baseTypeInt:
		// Note: Allowing negative values
		v := s.Int()
		if int64(fillUint(uint64(v), dConcrete)) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	case baseTypeUint:
		v := s.Uint()
		if fillUint(v, dConcrete) != v {
			panic(fmt.Errorf("%v cannot fit into type %v", v, dConcrete.Kind()))
		}
		return
	}
	panic(fmt.Errorf("Type %v is not compatible with type %v", s.Type(), dConcrete.Type()))
}

func setToNil(v reflect.Value) {
	v.Set(reflect.Zero(v.Type()))
}

func newStructReconstructor(t reflect.Type) objectReconstructor {
	var elementReconstructors []structElementReconstructor
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// TODO: handle Anonymous (embedded field)
		if field.PkgPath == "" { // PkgPath is empty for exported fields
			var elementReconstructor structElementReconstructor
			elementReconstructor.elementIndex = i
			elementReconstructor.mapKey = reflect.ValueOf(field.Name)
			elementReconstructor.reconstructElement = getReconstructorForType(field.Type)
			elementReconstructors = append(elementReconstructors, elementReconstructor)
		}
	}

	return func(s reflect.Value, d reflect.Value) {
		d = getObjectOrNewInstance(getConcreteObject(d))
		for _, elementReconstructor := range elementReconstructors {
			sValue := s.MapIndex(elementReconstructor.mapKey)
			if sValue.IsNil() {
				setToNil(d)
			} else {
				sValue = getConcreteObject(sValue)
				dValue := getConcreteObject(d.Field(elementReconstructor.elementIndex))
				elementReconstructor.reconstructElement(sValue, dValue)
			}
		}
	}
}

func newMapReconstructor(t reflect.Type) objectReconstructor {
	kt := t.Key()
	vt := t.Elem()

	// Note: kt and vt will not both be of type interface{} because
	//       map[interface{}]interface{} was registered in init().

	if kt.Kind() == reflect.Interface {
		vReconstructor := getReconstructorForType(vt)

		return func(s reflect.Value, d reflect.Value) {
			d = getObjectOrNewInstance(getConcreteObject(d))
			newMap := reflect.MakeMapWithSize(t, s.Len())
			for iter := s.MapRange(); iter.Next(); {
				v := reflect.New(vt).Elem()
				vReconstructor(iter.Value(), v)
				newMap.SetMapIndex(iter.Key(), v)
			}
			d.Set(newMap)
		}
	}

	if vt.Kind() == reflect.Interface {
		kReconstructor := getReconstructorForType(kt)

		return func(s reflect.Value, d reflect.Value) {
			d = getObjectOrNewInstance(getConcreteObject(d))
			newMap := reflect.MakeMapWithSize(t, s.Len())
			for iter := s.MapRange(); iter.Next(); {
				k := reflect.New(kt).Elem()
				kReconstructor(iter.Key().Elem(), k)
				newMap.SetMapIndex(k, iter.Value())
			}
			d.Set(newMap)
		}
	}

	kReconstructor := getReconstructorForType(kt)
	vReconstructor := getReconstructorForType(vt)

	return func(s reflect.Value, d reflect.Value) {
		d = getObjectOrNewInstance(getConcreteObject(d))
		newMap := reflect.MakeMapWithSize(t, s.Len())
		for iter := s.MapRange(); iter.Next(); {
			k := reflect.New(kt).Elem()
			v := reflect.New(vt).Elem()
			kReconstructor(iter.Key().Elem(), k)
			vReconstructor(iter.Value().Elem(), v)
			newMap.SetMapIndex(k, v)
		}
		d.Set(newMap)
	}
}

func newSliceReconstructor(t reflect.Type) objectReconstructor {
	elemType := t.Elem()

	// Note: elemType will not be of type interface{} because []interface{} was
	//       registered in init().

	elemReconstructor := getReconstructorForType(elemType)

	return func(s reflect.Value, d reflect.Value) {
		d = getObjectOrNewInstance(getConcreteObject(d))
		newSlice := reflect.MakeSlice(t, 0, s.Len())
		for i := 0; i < s.Len(); i++ {
			v := reflect.New(elemType).Elem()
			elemReconstructor(s.Index(i).Elem(), v)
			newSlice = reflect.Append(newSlice, v)
		}
		d.Set(newSlice)
	}
}

func newArrayReconstructor(t reflect.Type) objectReconstructor {
	elemType := t.Elem()
	if elemType.Kind() == reflect.Interface {
		return reconstructIdenticalKind
	}

	elemReconstructor := getReconstructorForType(elemType)

	return func(s reflect.Value, d reflect.Value) {
		d = getObjectOrNewInstance(getConcreteObject(d))
		if s.Len() != t.Len() {
			panic(fmt.Errorf("Src length (%v) and dst array length (%v) don't match", s.Len(), t.Len()))
		}
		for i := 0; i < s.Len(); i++ {
			elemReconstructor(s.Index(i), d.Index(i))
		}
	}
}

func generateReconstructorForType(t reflect.Type) (reconstructor objectReconstructor) {
	indirections := getTypeIndirections(t)
	directType := indirections[len(indirections)-1]

	ok := false
	if reconstructor, ok = knownReconstructors[directType]; ok {
		knownReconstructors[t] = reconstructor
		return
	}

	switch directType.Kind() {
	case reflect.Map:
		reconstructor = newMapReconstructor(directType)
	case reflect.Slice:
		reconstructor = newSliceReconstructor(directType)
	case reflect.Array:
		reconstructor = newArrayReconstructor(directType)
	case reflect.Struct:
		reconstructor = newStructReconstructor(directType)
	case reflect.Interface:
		reconstructor = reconstructIdenticalKind
	default:
		panic(fmt.Errorf("Cannot generate type reconstructor for type %v (kind %v, direct kind %v)", t, t.Kind(), directType.Kind()))
	}

	knownReconstructors[directType] = reconstructor
	knownReconstructors[t] = reconstructor
	return
}

func getReconstructorForType(t reflect.Type) objectReconstructor {
	if reconstructor, ok := knownReconstructors[t]; ok {
		return reconstructor
	}
	return generateReconstructorForType(t)
}

// ----------
// Public API
// ----------

// If true, function Reconstruct() will panic instead of returning an error.
// This is only useful for debugging problems in the library.
var PanicOnError = false

// Reconstruct an object from data in a neutral format.
// src is expected to be neutral format data.
// dst is expected to be addressible.
// It's assumed that the src data is compatible with dst. This function will
// fail with an error if this is not the case. The errors will come directly
// from the reflect package, and may be a little cryptic in some cases.
// Use PanicOnError to help debug the cause.
func Reconstruct(src interface{}, dst interface{}) (err error) {
	defer func() {
		if !PanicOnError {
			if e := recover(); e != nil {
				var ok bool
				err, ok = e.(error)
				if !ok {
					err = fmt.Errorf("%v", e)
				}
			}
		}
	}()

	s := getConcreteObject(reflect.ValueOf(src))
	d := getConcreteObject(reflect.ValueOf(dst))
	if !d.CanSet() {
		panic(fmt.Errorf("dst is not settable"))
	}
	reconstructor := getReconstructorForType(d.Type())
	reconstructor(s, d)
	return
}
