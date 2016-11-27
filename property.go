package gitconfig

import (
	"errors"
	"strconv"
	"strings"
)

var (
	InvalidBooleanError = errors.New("invalid boolean value")
	InvalidIntegerError = errors.New("invalid integer value")
)

// Property represents the name/value pair for a configuration property.
type Property interface {
	// Name returns the name of the property.
	Name() string

	// String returns the string representation of the property value.
	String() string

	// Bool returns the boolean value of the property. If the property value
	// is not a valid boolean, ok will be false.
	Bool() (value bool, ok bool)

	// List returns the list representation of the property. List splits the
	// string representation of the property value at colons (":").
	List() (list []string, ok bool)

	// Int returns the integer representation of the property. If the property
	// value is not a valid integer, ok will be false.
	Int() (int, bool)
}

// property is the implementation of the Property interface.
type property struct {
	name string
	v    string
}

// NewProperty returns a Property instance with the given name and value
func NewProperty(name, v string) (Property, Error) {
	return &property{name, v}, nil
} // NewProperty()

// Name returns the name of the property.
func (p property) Name() string { return p.name }

// String returns the string representation of the property value.
func (p property) String() string { return p.v }

// Bool returns the boolean value of the property. If the property value
// is not a valid boolean, the second return value will be false.
func (p property) Bool() (bool, bool) {
	// can we convert this property into a boolean?
	_bool, _err := NewBool(p.name, p.v)
	if _err != nil {
		return false, false
	} else {
		return _bool.Bool()
	}
} // Bool()

// List returns the list representation of the property. List splits the
// string representation of the property value at colons (":").
func (p property) List() ([]string, bool) {
	// can we convert this property into a boolean?
	//		- NewList() should never return an error
	_list, _err := NewList(p.name, p.v)
	if _err != nil {
		return nil, false
	} else {
		return _list.List()
	}
} // List()

// Int returns the integer representation of the property. If the property
// value is not a valid integer, the second return value will be false.
func (p property) Int() (int, bool) {
	// can we convert this property into an integer?
	_int, _err := NewInt(p.name, p.v)
	if _err != nil {
		return 0, false
	} else {
		return _int.Int()
	}
} // Int()

// boolean is the implementation of a boolean property
type boolean struct {
	Property
	b bool
}

// NewBool returns a property representing a boolean value. If the string v
// is not a valid boolean (e.g. "1", "on", "false", "no", etc) NewBool will
// return an Error.
func NewBool(name, v string) (Property, Error) {
	switch v {
	// true cases
	case "1":
		fallthrough
	case "on":
		fallthrough
	case "yes":
		fallthrough
	case "true":
		// NewProperty() should never return an error
		_property, _err := NewProperty(name, v)
		return &boolean{_property, true}, _err

	// false cases
	case "0":
		fallthrough
	case "off":
		fallthrough
	case "no":
		fallthrough
	case "false":
		// NewProperty() should never return an error
		_property, _err := NewProperty(name, v)
		return &boolean{_property, false}, _err
	}

	return nil, NewError(name, InvalidBooleanError)
} // NewBool()

// Bool returns the boolean value for the boolean property.
func (b boolean) Bool() (bool, bool) { return b.b, true }

// list is the implementation of a list property.
type list struct {
	Property
	l []string
}

// NewList returns a property representing a list value.
func NewList(name, v string) (Property, Error) {
	// split the string on ":"
	//		- NewProperty() should never return an error
	_list := strings.Split(v, ":")
	_property, _err := NewProperty(name, v)

	return &list{_property, _list}, _err
} // NewList()

// List returns the list representation of the list property.
func (l list) List() ([]string, bool) { return l.l, true }

// integer is the implementation of an integer property.
type integer struct {
	Property
	i int
}

// NewInt returns a property representing an integer value.
func NewInt(name, v string) (Property, Error) {
	// attempt to parse the integer property
	//		- NewProperty() should never return an error
	_int, _err := strconv.Atoi(v)
	if _err != nil {
		return nil, NewError(name, InvalidIntegerError)
	}
	_property, _error := NewProperty(name, v)

	return &integer{_property, _int}, _error
} // NewInt()

// Int returns the integer representation of an integer property.
func (i integer) Int() (int, bool) { return i.i, true }
