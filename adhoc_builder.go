package reconstruct

import (
	"net/url"
	"time"
)

type adhocBuilderState int

const (
	adhocBuilderStateTopLevel = iota
	adhocBuilderStateList
	adhocBuilderStateMapKey
	adhocBuilderStateMapValue
)

type adhocBuilderContext struct {
	State  adhocBuilderState
	List   []interface{}
	Map    map[interface{}]interface{}
	MapKey interface{}
}

// Builds an ad-hoc data type, consisting of data in neutral data format.
type AdhocBuilder struct {
	topLevelObject interface{}
	context        adhocBuilderContext
	contextStack   []adhocBuilderContext
}

func (this *AdhocBuilder) getCurrentState() adhocBuilderState {
	return this.context.State
}

func (this *AdhocBuilder) changeState(newState adhocBuilderState) {
	this.context.State = newState
}

func (this *AdhocBuilder) stackState(newState adhocBuilderState) {
	this.contextStack = append(this.contextStack, this.context)
	this.context.State = newState
}

func (this *AdhocBuilder) unstackState() {
	this.context = this.contextStack[len(this.contextStack)-1]
	this.contextStack = this.contextStack[:len(this.contextStack)-1]
}

func (this *AdhocBuilder) addValue(value interface{}) error {
	switch this.getCurrentState() {
	case adhocBuilderStateTopLevel:
		this.topLevelObject = value
	case adhocBuilderStateList:
		this.context.List = append(this.context.List, value)
	case adhocBuilderStateMapKey:
		this.context.MapKey = value
		this.changeState(adhocBuilderStateMapValue)
	case adhocBuilderStateMapValue:
		this.context.Map[this.context.MapKey] = value
		this.changeState(adhocBuilderStateMapKey)
	}
	return nil
}

// ---------
// Callbacks
// ---------

func (this *AdhocBuilder) OnNil() error {
	return this.addValue(nil)
}

func (this *AdhocBuilder) OnBool(value bool) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnInt(value int64) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnUint(value uint64) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnFloat(value float64) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnString(value string) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnBytes(value []byte) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnURI(value *url.URL) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnTime(value time.Time) error {
	return this.addValue(value)
}

func (this *AdhocBuilder) OnListBegin() error {
	this.stackState(adhocBuilderStateList)
	this.context.List = make([]interface{}, 0)
	return nil
}

func (this *AdhocBuilder) OnMapBegin() error {
	this.stackState(adhocBuilderStateMapKey)
	this.context.Map = make(map[interface{}]interface{})
	return nil
}

func (this *AdhocBuilder) OnContainerEnd() error {
	var v interface{} = this.context.Map
	if this.getCurrentState() == adhocBuilderStateList {
		v = this.context.List
	}
	this.unstackState()
	this.addValue(v)
	return nil
}

func (this *AdhocBuilder) OnMarker(name uint64) error {
	// TODO
	return nil
}

func (this *AdhocBuilder) OnReference(name uint64) error {
	// TODO
	return nil
}

// ----------
// Public API
// ----------

// Get the built object. Call this after using the AdhocBuilder as an
// ObjectIteratorCallbacks to fill it.
func (this *AdhocBuilder) GetObject() interface{} {
	return this.topLevelObject
}
