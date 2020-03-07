package reconstruct

import (
	"net/url"
	"reflect"
	"time"
)

// -------
// []uint8
// -------

type uint8SliceIterator struct {
	root *RootObjectIterator
}

func newUInt8SliceIterator() ObjectIterator {
	return &uint8SliceIterator{}
}

func (this *uint8SliceIterator) PostCacheInitIterator() {
}

func (this *uint8SliceIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &uint8SliceIterator{root: root}
}

func (this *uint8SliceIterator) Iterate(v reflect.Value) error {
	return this.root.callbacks.OnBytes(v.Bytes())
}

// ----
// Time
// ----

type timeIterator struct {
	root *RootObjectIterator
}

func newTimeIterator() ObjectIterator {
	return &timeIterator{}
}

func (this *timeIterator) PostCacheInitIterator() {
}

func (this *timeIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &timeIterator{root: root}
}

func (this *timeIterator) Iterate(v reflect.Value) error {
	return this.root.callbacks.OnTime(v.Interface().(time.Time))
}

// ----
// *URL
// ----

type pURLIterator struct {
	root *RootObjectIterator
}

func newPURLIterator() ObjectIterator {
	return &pURLIterator{}
}

func (this *pURLIterator) PostCacheInitIterator() {
}

func (this *pURLIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &pURLIterator{root: root}
}

func (this *pURLIterator) Iterate(v reflect.Value) error {
	if v.IsNil() {
		return this.root.callbacks.OnNil()
	}
	return this.root.callbacks.OnURI(v.Interface().(*url.URL))
}

// ---
// URL
// ---

type urlIterator struct {
	root *RootObjectIterator
}

func newURLIterator() ObjectIterator {
	return &urlIterator{}
}

func (this *urlIterator) PostCacheInitIterator() {
}

func (this *urlIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &urlIterator{root: root}
}

func (this *urlIterator) Iterate(v reflect.Value) error {
	vCopy := v.Interface().(url.URL)
	return this.root.callbacks.OnURI(&vCopy)
}

// ----
// Bool
// ----

type boolIterator struct {
	root *RootObjectIterator
}

func newBoolIterator() ObjectIterator {
	return &boolIterator{}
}

func (this *boolIterator) PostCacheInitIterator() {
}

func (this *boolIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &boolIterator{root: root}
}

func (this *boolIterator) Iterate(v reflect.Value) error {
	return this.root.callbacks.OnBool(v.Bool())
}

// ---
// Int
// ---

type intIterator struct {
	root *RootObjectIterator
}

func newIntIterator() ObjectIterator {
	return &intIterator{}
}

func (this *intIterator) PostCacheInitIterator() {
}

func (this *intIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &intIterator{root: root}
}

func (this *intIterator) Iterate(v reflect.Value) error {
	return this.root.callbacks.OnInt(v.Int())
}

// ----
// Uint
// ----

type uintIterator struct {
	root *RootObjectIterator
}

func newUintIterator() ObjectIterator {
	return &uintIterator{}
}

func (this *uintIterator) PostCacheInitIterator() {
}

func (this *uintIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &uintIterator{root: root}
}

func (this *uintIterator) Iterate(v reflect.Value) error {
	return this.root.callbacks.OnUint(v.Uint())
}

// -----
// Float
// -----

type floatIterator struct {
	root *RootObjectIterator
}

func newFloatIterator() ObjectIterator {
	return &floatIterator{}
}

func (this *floatIterator) PostCacheInitIterator() {
}

func (this *floatIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &floatIterator{root: root}
}

func (this *floatIterator) Iterate(v reflect.Value) error {
	return this.root.callbacks.OnFloat(v.Float())
}

// ------
// String
// ------

type stringIterator struct {
	root *RootObjectIterator
}

func newStringIterator() ObjectIterator {
	return &stringIterator{}
}

func (this *stringIterator) PostCacheInitIterator() {
}

func (this *stringIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &stringIterator{root: root}
}

func (this *stringIterator) Iterate(v reflect.Value) error {
	return this.root.callbacks.OnString(v.String())
}
