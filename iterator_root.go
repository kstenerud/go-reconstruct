package reconstruct

import (
	"reflect"
)

func NewRootObjectIterator(useReferences bool, callbacks ObjectIteratorCallbacks) *RootObjectIterator {
	this := new(RootObjectIterator)
	this.Init(useReferences, callbacks)
	return this
}

func (this *RootObjectIterator) Init(useReferences bool, callbacks ObjectIteratorCallbacks) {
	this.useReferences = useReferences
	this.callbacks = callbacks
}

func (this *RootObjectIterator) Iterate(value interface{}) error {
	rv := reflect.ValueOf(value)
	iterator := getIteratorForType(rv.Type())
	iterator = iterator.CloneFromTemplate(this)
	return iterator.Iterate(rv)
}

// Iterates depth-first recursively through an object, notifying callbacks as it
// encounters data.
type RootObjectIterator struct {
	foundReferences map[reflect.Value]bool
	namedReferences map[reflect.Value]uint32
	nextMarkerName  uint32
	callbacks       ObjectIteratorCallbacks
	useReferences   bool
}

// ----------
// References
// ----------

func (this *RootObjectIterator) checkReferenceExists(value reflect.Value) (alreadyExists bool) {
	if _, ok := this.foundReferences[value]; ok {
		this.foundReferences[value] = true
		return true
	}

	this.foundReferences[value] = false
	return false
}

func (this *RootObjectIterator) findReferencesInValue(value reflect.Value) {
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

func (this *RootObjectIterator) findReferences(value reflect.Value) {
	if this.useReferences {
		this.foundReferences = make(map[reflect.Value]bool)
		this.namedReferences = make(map[reflect.Value]uint32)
		this.findReferencesInValue(value)
	}
}

func (this *RootObjectIterator) addReference(v reflect.Value) (didAddReferenceObject bool) {
	if this.useReferences {
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
	}
	return false
}
