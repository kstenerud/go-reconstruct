package reconstruct

import (
	"net/url"
	"reflect"
	"time"
)

func (this *RootBuilder) GetBuiltObject() interface{} {
	if this.object.IsValid() && this.object.CanInterface() {
		return this.object.Interface()
	}
	return nil
}

// RootBuilder adapts ObjectIteratorCallbacks to ObjectBuilder, coordinates the
// build, and provides GetBuiltObject() for fetching the final result.
type RootBuilder struct {
	dstType        reflect.Type
	currentBuilder ObjectBuilder
	object         reflect.Value
}

// -----------
// RootBuilder
// -----------

func newRootBuilder(dstType reflect.Type) *RootBuilder {
	this := &RootBuilder{
		dstType: dstType,
		object:  reflect.New(dstType).Elem(),
	}

	builder := getTopLevelBuilderForType(dstType)
	this.currentBuilder = builder.CloneFromTemplate(this, this)

	return this
}

func (this *RootBuilder) setCurrentBuilder(builder ObjectBuilder) {
	this.currentBuilder = builder
}

// -------------
// ObjectBuilder
// -------------

func (this *RootBuilder) PostCacheInitBuilder() {
}

func (this *RootBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}
func (this *RootBuilder) Nil(dst reflect.Value) {
	this.currentBuilder.Nil(dst)
}
func (this *RootBuilder) Bool(value bool, dst reflect.Value) {
	this.currentBuilder.Bool(value, dst)
}
func (this *RootBuilder) Int(value int64, dst reflect.Value) {
	this.currentBuilder.Int(value, dst)
}
func (this *RootBuilder) Uint(value uint64, dst reflect.Value) {
	this.currentBuilder.Uint(value, dst)
}
func (this *RootBuilder) Float(value float64, dst reflect.Value) {
	this.currentBuilder.Float(value, dst)
}
func (this *RootBuilder) String(value string, dst reflect.Value) {
	this.currentBuilder.String(value, dst)
}
func (this *RootBuilder) Bytes(value []byte, dst reflect.Value) {
	this.currentBuilder.Bytes(value, dst)
}
func (this *RootBuilder) URI(value *url.URL, dst reflect.Value) {
	this.currentBuilder.URI(value, dst)
}
func (this *RootBuilder) Time(value time.Time, dst reflect.Value) {
	this.currentBuilder.Time(value, dst)
}
func (this *RootBuilder) List() {
	this.currentBuilder.List()
}
func (this *RootBuilder) Map() {
	this.currentBuilder.Map()
}
func (this *RootBuilder) End() {
	this.currentBuilder.End()
}
func (this *RootBuilder) Marker(id interface{}) {
	panic("TODO")
}
func (this *RootBuilder) Reference(id interface{}) {
	panic("TODO")
}
func (this *RootBuilder) PrepareForListContents() {
	panic("BUG")
}
func (this *RootBuilder) PrepareForMapContents() {
	panic("BUG")
}
func (this *RootBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.object = value
}

// -----------------------
// ObjectIteratorCallbacks
// -----------------------

func (this *RootBuilder) OnNil() error {
	this.Nil(this.object)
	return nil
}
func (this *RootBuilder) OnBool(value bool) error {
	this.Bool(value, this.object)
	return nil
}
func (this *RootBuilder) OnInt(value int64) error {
	this.Int(value, this.object)
	return nil
}
func (this *RootBuilder) OnUint(value uint64) error {
	this.Uint(value, this.object)
	return nil
}
func (this *RootBuilder) OnFloat(value float64) error {
	this.Float(value, this.object)
	return nil
}
func (this *RootBuilder) OnComplex(value complex128) error {
	panic("TODO")
	return nil
}
func (this *RootBuilder) OnString(value string) error {
	this.String(value, this.object)
	return nil
}
func (this *RootBuilder) OnBytes(value []byte) error {
	this.Bytes(value, this.object)
	return nil
}
func (this *RootBuilder) OnURI(value *url.URL) error {
	this.URI(value, this.object)
	return nil
}
func (this *RootBuilder) OnTime(value time.Time) error {
	this.Time(value, this.object)
	return nil
}
func (this *RootBuilder) OnListBegin() error {
	this.List()
	return nil
}
func (this *RootBuilder) OnMapBegin() error {
	this.Map()
	return nil
}
func (this *RootBuilder) OnContainerEnd() error {
	this.End()
	return nil
}
func (this *RootBuilder) OnMarker(id interface{}) error {
	this.Marker(id)
	return nil
}
func (this *RootBuilder) OnReference(id interface{}) error {
	this.Reference(id)
	return nil
}
