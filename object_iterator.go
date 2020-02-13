package reconstruct

import (
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// These callbacks will be called as the iterator iterates through an object
// and discovers data.
type ObjectIteratorCallbacks interface {
	OnNil() error
	OnBool(value bool) error
	OnInt(value int64) error
	OnUint(value uint64) error
	OnFloat(value float64) error
	OnString(value string) error
	OnBytes(value []byte) error
	OnURI(value *url.URL) error
	OnTime(value time.Time) error
	OnListBegin() error
	OnListEnd() error
	OnMapBegin() error
	OnMapEnd() error
	// TODO: Marker and reference
	OnMarker(name uint64) error
	OnReference(name uint64) error
}

// Iterates recursively through an object, notifying a callbacks object as it
// encounters data.
type ObjectIterator struct {
	foundReferences map[reflect.Value]bool
	namedReferences map[reflect.Value]uint32
	nextMarkerName  uint32
	callbacks       ObjectIteratorCallbacks
}

// -----------------
// Reference Finding
// -----------------

func (this *ObjectIterator) checkReferenceExists(value reflect.Value) (alreadyExists bool) {
	if _, ok := this.foundReferences[value]; ok {
		this.foundReferences[value] = true
		return true
	}

	this.foundReferences[value] = false
	return false
}

func (this *ObjectIterator) findReferencesInValue(value reflect.Value) {
	if !value.IsValid() {
		return
	}

	switch value.Kind() {
	case reflect.Interface:
		this.findReferencesInValue(value.Elem())
	case reflect.Ptr:
		if !this.checkReferenceExists(value) {
			this.findReferencesInValue(value.Elem())
		}
	case reflect.Map:
		if !this.checkReferenceExists(value) {
			iter := value.MapRange()
			for iter.Next() {
				this.findReferencesInValue(iter.Key())
				this.findReferencesInValue(iter.Value())
			}
		}
	case reflect.Slice:
		if !this.checkReferenceExists(value) {
			count := value.Len()
			for i := 0; i < count; i++ {
				this.findReferencesInValue(value.Index(i))
			}
		}
	case reflect.Array:
		count := value.Len()
		for i := 0; i < count; i++ {
			this.findReferencesInValue(value.Index(i))
		}
	case reflect.Struct:
		count := value.NumField()
		for i := 0; i < count; i++ {
			this.findReferencesInValue(value.Field(i))
		}
	}
}

func (this *ObjectIterator) findReferences(value reflect.Value) {
	this.foundReferences = make(map[reflect.Value]bool)
	this.namedReferences = make(map[reflect.Value]uint32)
	this.findReferencesInValue(value)
}

// ---------
// Iteration
// ---------

func (this *ObjectIterator) iterateSliceOrArray(v reflect.Value) (err error) {
	t := v.Type()

	if t.Elem().Kind() == reflect.Uint8 {
		if v.Kind() == reflect.Array {
			if !v.CanAddr() {
				tempSlice := make([]byte, v.Len())
				tempLen := v.Len()
				for i := 0; i < tempLen; i++ {
					tempSlice[i] = v.Index(i).Interface().(uint8)
				}
				return this.callbacks.OnBytes(tempSlice)
			}
			v = v.Slice(0, v.Len())
		}
		return this.callbacks.OnBytes(v.Bytes())
	}

	if err = this.callbacks.OnListBegin(); err != nil {
		return
	}
	length := v.Len()
	for i := 0; i < length; i++ {
		if err = this.iterateValue(v.Index(i)); err != nil {
			return
		}
	}
	if err = this.callbacks.OnListEnd(); err != nil {
		return
	}
	return
}

func (this *ObjectIterator) iterateMap(value reflect.Value) (err error) {
	if err = this.callbacks.OnMapBegin(); err != nil {
		return
	}

	iter := value.MapRange()
	for iter.Next() {
		if err = this.iterateValue(iter.Key()); err != nil {
			return
		}
		if err = this.iterateValue(iter.Value()); err != nil {
			return
		}
	}

	if err = this.callbacks.OnMapEnd(); err != nil {
		return
	}

	return
}

func (this *ObjectIterator) iterateStruct(v reflect.Value) (err error) {
	t := v.Type()
	if t.Name() == "Time" && t.PkgPath() == "time" {
		return this.callbacks.OnTime(v.Interface().(time.Time))
	}
	if t.Name() == "URL" && t.PkgPath() == "net/url" {
		realValue := v.Interface().(url.URL)
		return this.callbacks.OnURI(&realValue)
	}

	if err = this.callbacks.OnMapBegin(); err != nil {
		return
	}

	fieldCount := t.NumField()
	for i := 0; i < fieldCount; i++ {
		this.callbacks.OnString(t.Field(i).Name)
		this.iterateValue(v.Field(i))
	}

	if err = this.callbacks.OnMapEnd(); err != nil {
		return
	}

	return nil
}

func (this *ObjectIterator) addReference(v reflect.Value) (didAddReferenceObject bool) {
	if this.foundReferences[v] {
		var name uint32
		var exists bool
		if name, exists = this.namedReferences[v]; !exists {
			name = this.nextMarkerName
			this.nextMarkerName++
			this.namedReferences[v] = name
			this.callbacks.OnMarker(uint64(name))
			return false
		} else {
			this.callbacks.OnReference(uint64(name))
			return true
		}
	}
	return false
}

func (this *ObjectIterator) iterateValue(v reflect.Value) error {
	if !v.IsValid() {
		return this.callbacks.OnNil()
	}

	switch v.Kind() {
	case reflect.Bool:
		return this.callbacks.OnBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return this.callbacks.OnInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return this.callbacks.OnUint(v.Uint())
	case reflect.Float32, reflect.Float64:
		return this.callbacks.OnFloat(v.Float())

	case reflect.Complex64, reflect.Complex128:
		cv := v.Complex()
		if err := this.callbacks.OnListBegin(); err != nil {
			return err
		}
		if err := this.callbacks.OnFloat(real(cv)); err != nil {
			return err
		}
		if err := this.callbacks.OnFloat(imag(cv)); err != nil {
			return err
		}
		return this.callbacks.OnListEnd()
	case reflect.String:
		return this.callbacks.OnString(v.String())
	case reflect.Array:
		return this.iterateSliceOrArray(v)
	case reflect.Slice:
		if v.IsNil() {
			return this.callbacks.OnNil()
		}
		if this.addReference(v) {
			return nil
		}
		return this.iterateSliceOrArray(v)
	case reflect.Map:
		if v.IsNil() {
			return this.callbacks.OnNil()
		}
		if this.addReference(v) {
			return nil
		}
		return this.iterateMap(v)
	case reflect.Struct:
		return this.iterateStruct(v)
	case reflect.Ptr:
		if v.IsNil() {
			return this.callbacks.OnNil()
		}
		if this.addReference(v) {
			return nil
		}
		return this.iterateValue(v.Elem())
	case reflect.Interface:
		return this.iterateValue(v.Elem())
	case reflect.Chan,
		reflect.Func,
		reflect.Uintptr,
		reflect.UnsafePointer:
		// Do nothing
		return nil

	default:
		return fmt.Errorf("%v: Unhandled type", v.Kind())
	}
}

func (this *ObjectIterator) iterate(value reflect.Value) error {
	this.findReferences(value)
	return this.iterateValue(value)
}

// ----------
// Public API
// ----------

func NewObjectIterator(callbacks ObjectIteratorCallbacks) *ObjectIterator {
	this := new(ObjectIterator)
	this.Init(callbacks)
	return this
}

func (this *ObjectIterator) Init(callbacks ObjectIteratorCallbacks) {
	this.callbacks = callbacks
}

func (this *ObjectIterator) Iterate(value interface{}) error {
	return this.iterate(reflect.ValueOf(value))
}

// Iterate over an object (recursively), calling the callbacks as it encounters
// data.
func IterateObject(value interface{}, callbacks ObjectIteratorCallbacks) error {
	iter := NewObjectIterator(callbacks)
	return iter.Iterate(value)
}
