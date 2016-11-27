package gitconfig

import (
	"errors"
	"strconv"
	"strings"
)

var (
	InvalidBooleanError = errors.New("invalid boolean value")
	InvalidIntegerError = errors.New("invalid integer value")

	// true and false values for boolean properties
	_True  = true
	_False = false
)

// Property represents the name/value pair for a configuration property.
type Property interface {
	// Name returns the name of the property.
	Name() string

	// String returns the string representation of the property value.
	String() string

	// Bool returns the boolean value of the property. If the property value
	// is not a valid boolean, an error will be returned.
	Bool() (bool, error)

	// List returns the list representation of the property. List splits the
	// string representation of the property value at colons ":".
	List() []string

	// Int returns the integer representation of the property. If the property
	// value is not a valid integer, an error will be returned.
	Int() (int, error)
}

// property is the implementation of the Property interface.
type property struct {
	name string
	v    string
	b    *bool
	i    *int
	l    []string
}

// NewProperty returns a Property instance with the given name and value v.
func NewProperty(name, v string) Property {
	return &property{name: name, v: v}
} // NewProperty()

// Name returns the name of the property.
func (p property) Name() string { return p.name }

// String returns the string representation of the property value.
func (p property) String() string { return p.v }

// Bool returns the boolean value of the property. If the property value
// is not a valid boolean, Bool returns the InvalidBooleanError.
func (p property) Bool() (bool, error) {
	if p.b == nil {
		p.b = boolean(p.v)
	}

	if p.b == nil {
		return false, InvalidBooleanError
	} else {
		return *p.b, nil
	}
} // Bool()

// List returns the list representation of the property. List splits the
// string representation of the property value at colons ":".
func (p property) List() []string {
	if p.l == nil {
		p.l = list(p.v)
	}

	return p.l
} // List()

// Int returns the integer representation of the property. If the property
// value is not a valid integer, Int returns the InvalidIntegerError.
func (p property) Int() (int, error) {
	if p.i == nil {
		p.i = integer(p.v)
	}

	if p.i == nil {
		return 0, InvalidIntegerError
	} else {
		return *p.i, nil
	}
} // Int()

//
// helper methods
//

// boolean converts the given string v into a boolean, if the string represnts
// a valid boolean value (such as "1", "true", "off", "no", etc). boolean
// returns nil otherwise.
func boolean(v string) *bool {
	switch v {
	// true cases
	case "1":
		fallthrough
	case "on":
		fallthrough
	case "yes":
		fallthrough
	case "true":
		return &_True

	// false cases
	case "0":
		fallthrough
	case "off":
		fallthrough
	case "no":
		fallthrough
	case "false":
		return &_False
	}

	//	return nil, InvalidBooleanError
	return nil
} // boolean()

// list returns the list representation of the value string s. Values are split
// on colons ":".
func list(v string) []string {
	// split the string on ":"
	return strings.Split(v, ":")
} // list()

// integer converts the given string v into an integer, if the string
// represents a valid integer value. integer returns nil otherwise.
func integer(v string) *int {
	// attempt to parse the integer property
	_int, _err := strconv.Atoi(v)
	if _err == nil {
		return &_int
	} else {
		return nil
	}
} // integer()
