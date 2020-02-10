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

// Builds an ad-hoc data type, consisting of data in neutral data format.
type AdhocBuilder struct {
	topLevelObject interface{}
	containerStack []interface{}
	mapKeyStack    []interface{}
	currentList    []interface{}
	currentMap     map[interface{}]interface{}
	currentMapKey  interface{}
	state          []adhocBuilderState
}

func (this *AdhocBuilder) getCurrentState() adhocBuilderState {
	if len(this.state) == 0 {
		return adhocBuilderStateTopLevel
	}
	return this.state[len(this.state)-1]
}

func (this *AdhocBuilder) changeState(newState adhocBuilderState) {
	this.state[len(this.state)-1] = newState
}

func (this *AdhocBuilder) stackState(newState adhocBuilderState) {
	switch this.getCurrentState() {
	case adhocBuilderStateList:
		this.containerStack = append(this.containerStack, this.currentList)
	case adhocBuilderStateMapValue:
		this.containerStack = append(this.containerStack, this.currentMap)
	}
	this.mapKeyStack = append(this.mapKeyStack, this.currentMapKey)
	this.state = append(this.state, newState)
}

func (this *AdhocBuilder) unstackState() {
	this.state = this.state[:len(this.state)-1]
	this.currentMapKey = this.mapKeyStack[len(this.mapKeyStack)-1]
	this.mapKeyStack = this.mapKeyStack[:len(this.mapKeyStack)-1]

	switch this.getCurrentState() {
	case adhocBuilderStateList:
		this.currentList = this.containerStack[len(this.containerStack)-1].([]interface{})
		this.containerStack = this.containerStack[:len(this.containerStack)-1]
	case adhocBuilderStateMapValue:
		this.currentMap = this.containerStack[len(this.containerStack)-1].(map[interface{}]interface{})
		this.containerStack = this.containerStack[:len(this.containerStack)-1]
	}
}

func (this *AdhocBuilder) stackContainer(container interface{}) {
	this.containerStack = append(this.containerStack, container)
}

func (this *AdhocBuilder) unstackContainer() {
	this.state = this.state[:len(this.state)-1]
}

func (this *AdhocBuilder) addValue(value interface{}) error {
	switch this.getCurrentState() {
	case adhocBuilderStateTopLevel:
		this.topLevelObject = value
	case adhocBuilderStateList:
		this.currentList = append(this.currentList, value)
	case adhocBuilderStateMapKey:
		this.currentMapKey = value
		this.changeState(adhocBuilderStateMapValue)
	case adhocBuilderStateMapValue:
		this.currentMap[this.currentMapKey] = value
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
	this.currentList = this.currentList[:0]
	return nil
}

func (this *AdhocBuilder) OnListEnd() error {
	v := this.currentList
	this.unstackState()
	this.addValue(v)
	return nil
}

func (this *AdhocBuilder) OnMapBegin() error {
	this.stackState(adhocBuilderStateMapKey)
	this.currentMap = make(map[interface{}]interface{})
	return nil
}

func (this *AdhocBuilder) OnMapEnd() error {
	v := this.currentMap
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
