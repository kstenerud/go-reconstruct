package reconstruct

import (
	"reflect"
)

// FindDuplicatePointers looks deep within an object for pointers that are used
// multiple times, marking slice, map, or pointer members that point to the same
// instance of an object.
//
// The returned map has boolean true entries for every reflect.Value that is a
// duplicate pointer. Values that are not duplicates will either not be present
// in the map, or will have a false value. Either way, fetching will return
// false if the key is not a duplicate pointer.
func FindDuplicatePointers(value interface{}) (foundPtrs map[reflect.Value]bool) {
	foundPtrs = make(map[reflect.Value]bool)
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

func findDuplicatePtrsInValue(value reflect.Value, foundPtrs map[reflect.Value]bool) {
	if !isSearchableKind(value.Kind()) {
		return
	}

	switch value.Kind() {
	case reflect.Interface:
		findDuplicatePtrsInValue(value.Elem(), foundPtrs)
	case reflect.Ptr:
		if !checkPtrAlreadyFound(value, foundPtrs) {
			findDuplicatePtrsInValue(value.Elem(), foundPtrs)
		}
	case reflect.Map:
		if !checkPtrAlreadyFound(value, foundPtrs) {
			if isSearchableKind(value.Type().Elem().Kind()) {
				iter := value.MapRange()
				for iter.Next() {
					findDuplicatePtrsInValue(iter.Value(), foundPtrs)
				}
			}
		}
	case reflect.Slice:
		if !checkPtrAlreadyFound(value, foundPtrs) {
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
			findDuplicatePtrsInValue(value.Field(i), foundPtrs)
		}
	}
}

func checkPtrAlreadyFound(value reflect.Value, foundPtrs map[reflect.Value]bool) (alreadyExists bool) {
	if _, ok := foundPtrs[value]; ok {
		foundPtrs[value] = true
		return true
	}

	foundPtrs[value] = false
	return false
}
