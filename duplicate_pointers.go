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

func findDuplicatePtrsInValue(value reflect.Value, foundPtrs map[reflect.Value]bool) {
	if !value.IsValid() {
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
			iter := value.MapRange()
			for iter.Next() {
				findDuplicatePtrsInValue(iter.Key(), foundPtrs)
				findDuplicatePtrsInValue(iter.Value(), foundPtrs)
			}
		}
	case reflect.Slice:
		if !checkPtrAlreadyFound(value, foundPtrs) {
			count := value.Len()
			for i := 0; i < count; i++ {
				findDuplicatePtrsInValue(value.Index(i), foundPtrs)
			}
		}
	case reflect.Array:
		count := value.Len()
		for i := 0; i < count; i++ {
			findDuplicatePtrsInValue(value.Index(i), foundPtrs)
		}
	case reflect.Struct:
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
