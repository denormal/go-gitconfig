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

type Value interface {
	Name() string
	String() string
	Bool() (bool, bool)
	List() ([]string, bool)
	Int() (int, bool)
}

type value struct {
	name string
	v    string
}

func NewValue(name, v string) (Value, Error) {
	return &value{name, v}, nil
} // NewValue()

func (v value) Name() string   { return v.name }
func (v value) String() string { return v.v }

func (v value) Bool() (bool, bool) {
	// can we convert this value into a boolean?
	_bool, _err := NewBool(v.name, v.v)
	if _err != nil {
		return false, false
	} else {
		return _bool.Bool()
	}
} // Bool()

func (v value) List() ([]string, bool) {
	// can we convert this value into a boolean?
	//		- NewList() should never return an error
	_list, _err := NewList(v.name, v.v)
	if _err != nil {
		return nil, false
	} else {
		return _list.List()
	}
} // List()

func (v value) Int() (int, bool) {
	// can we convert this value into an integer?
	_int, _err := NewInt(v.name, v.v)
	if _err != nil {
		return 0, false
	} else {
		return _int.Int()
	}
} // Int()

type boolean struct {
	Value
	b bool
}

func NewBool(name, v string) (Value, Error) {
	switch v {
	// true cases
	case "1":
		fallthrough
	case "on":
		fallthrough
	case "yes":
		fallthrough
	case "true":
		// NewValue() should never return an error
		_value, _err := NewValue(name, v)
		return &boolean{_value, true}, _err

	// false cases
	case "0":
		fallthrough
	case "off":
		fallthrough
	case "no":
		fallthrough
	case "false":
		// NewValue() should never return an error
		_value, _err := NewValue(name, v)
		return &boolean{_value, false}, _err
	}

	return nil, NewError(name, InvalidBooleanError)
} // NewBool()

func (b boolean) Bool() (bool, bool) { return b.b, true }

type list struct {
	Value
	l []string
}

func NewList(name, v string) (Value, Error) {
	// split the string on ":"
	//		- NewValue() should never return an error
	_list := strings.Split(v, ":")
	_value, _err := NewValue(name, v)

	return &list{_value, _list}, _err
} // NewList()

func (l list) List() ([]string, bool) { return l.l, true }

type integer struct {
	Value
	i int
}

func NewInt(name, v string) (Value, Error) {
	// attempt to parse the integer value
	//		- NewValue() should never return an error
	_int, _err := strconv.Atoi(v)
	if _err != nil {
		return nil, NewError(name, InvalidIntegerError)
	}
	_value, _error := NewValue(name, v)

	return &integer{_value, _int}, _error
} // NewInt()

func (i integer) Int() (int, bool) { return i.i, true }
