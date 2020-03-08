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
	this.findReferences(value)
	rv := reflect.ValueOf(value)
	iterator := getIteratorForType(rv.Type())
	iterator = iterator.CloneFromTemplate(this)
	return iterator.Iterate(rv)
}

// Iterates depth-first recursively through an object, notifying callbacks as it
// encounters data.
type RootObjectIterator struct {
	foundReferences map[uintptr]bool
	namedReferences map[uintptr]uint32
	nextMarkerName  uint32
	callbacks       ObjectIteratorCallbacks
	useReferences   bool
}

func (this *RootObjectIterator) findReferences(value interface{}) {
	if this.useReferences {
		this.foundReferences = FindDuplicatePointers(value)
		this.namedReferences = make(map[uintptr]uint32)
	}
}

func (this *RootObjectIterator) addReference(v reflect.Value) (didAddReferenceObject bool) {
	if this.useReferences {
		ptr := v.Pointer()
		if this.foundReferences[ptr] {
			var name uint32
			var exists bool
			if name, exists = this.namedReferences[ptr]; !exists {
				name = this.nextMarkerName
				this.nextMarkerName++
				this.namedReferences[ptr] = name
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
