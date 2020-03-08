package reconstruct

import (
	"reflect"
)

// FindDuplicatePointers walks an object and its contents looking for pointers
// that are used multiple times, marking slices, maps, or pointers that point to
// the same instance of an object. Both exported and unexported members are
// examined and followed.
//
// The returned foundPtrs will map to true for every duplicate pointer found.
// Non-duplicates will either not be present in the map, or will map to false.
// Either way, foundPtrs[myPointer] will return false if myPointer is not a
// duplicate pointer.
func FindDuplicatePointers(value interface{}) (foundPtrs map[uintptr]bool) {
	foundPtrs = make(map[uintptr]bool)
	findDuplicatePtrsInValue(reflect.ValueOf(value), foundPtrs)
	return
}

func isSearchableKind(kind reflect.Kind) bool {
	const searchableKinds uint = (uint(1) << reflect.Interface) |
		(uint(1) << reflect.Ptr) |
		(uint(1) << reflect.Slice) |
		(uint(1) << reflect.Map) |
		(uint(1) << reflect.Array) |
		(uint(1) << reflect.Struct)

	return searchableKinds&(uint(1)<<kind) != 0
}

func findDuplicatePtrsInValue(value reflect.Value, foundPtrs map[uintptr]bool) {
	switch value.Kind() {
	case reflect.Interface:
		if !value.IsNil() {
			value = value.Elem()
			if isSearchableKind(value.Kind()) {
				findDuplicatePtrsInValue(value, foundPtrs)
			}
		}
	case reflect.Ptr:
		if !value.IsNil() && !checkPtrAlreadyFound(value, foundPtrs) {
			if isSearchableKind(value.Type().Elem().Kind()) {
				findDuplicatePtrsInValue(value.Elem(), foundPtrs)
			}
		}
	case reflect.Map:
		if !value.IsNil() && !checkPtrAlreadyFound(value, foundPtrs) {
			if isSearchableKind(value.Type().Elem().Kind()) {
				iter := value.MapRange()
				for iter.Next() {
					findDuplicatePtrsInValue(iter.Value(), foundPtrs)
				}
			}
		}
	case reflect.Slice:
		if !value.IsNil() && !checkPtrAlreadyFound(value, foundPtrs) {
			if isSearchableKind(value.Type().Elem().Kind()) {
				count := value.Len()
				for i := 0; i < count; i++ {
					findDuplicatePtrsInValue(value.Index(i), foundPtrs)
				}
			}
		}
	case reflect.Array:
		if isSearchableKind(value.Type().Elem().Kind()) {
			count := value.Len()
			for i := 0; i < count; i++ {
				findDuplicatePtrsInValue(value.Index(i), foundPtrs)
			}
		}
	case reflect.Struct:
		// TODO: How to eliminate non-searchable fields?
		count := value.NumField()
		for i := 0; i < count; i++ {
			field := value.Field(i)
			if isSearchableKind(field.Kind()) {
				findDuplicatePtrsInValue(field, foundPtrs)
			}
		}
	}
}

func checkPtrAlreadyFound(value reflect.Value, foundPtrs map[uintptr]bool) (alreadyExists bool) {
	ptr := value.Pointer()
	if _, ok := foundPtrs[ptr]; ok {
		foundPtrs[ptr] = true
		return true
	}

	foundPtrs[ptr] = false
	return false
}
