// +build go1.12

package reconstruct

import (
	"reflect"
)

func mapRange(v reflect.Value) *reflect.MapIter {
	return v.MapRange()
}
