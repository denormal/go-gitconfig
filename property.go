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

type Property interface {
	Name() string
	String() string
	Bool() (bool, bool)
	List() ([]string, bool)
	Int() (int, bool)
}

type property struct {
	name string
	v    string
}

func NewProperty(name, v string) (Property, Error) {
	return &property{name, v}, nil
} // NewProperty()

func (p property) Name() string   { return p.name }
func (p property) String() string { return p.v }

func (p property) Bool() (bool, bool) {
	// can we convert this property into a boolean?
	_bool, _err := NewBool(p.name, p.v)
	if _err != nil {
		return false, false
	} else {
		return _bool.Bool()
	}
} // Bool()

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

func (p property) Int() (int, bool) {
	// can we convert this property into an integer?
	_int, _err := NewInt(p.name, p.v)
	if _err != nil {
		return 0, false
	} else {
		return _int.Int()
	}
} // Int()

type boolean struct {
	Property
	b bool
}

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

func (b boolean) Bool() (bool, bool) { return b.b, true }

type list struct {
	Property
	l []string
}

func NewList(name, v string) (Property, Error) {
	// split the string on ":"
	//		- NewProperty() should never return an error
	_list := strings.Split(v, ":")
	_property, _err := NewProperty(name, v)

	return &list{_property, _list}, _err
} // NewList()

func (l list) List() ([]string, bool) { return l.l, true }

type integer struct {
	Property
	i int
}

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

func (i integer) Int() (int, bool) { return i.i, true }
