package gitconfig

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Config is the interface to a block of git configuration, either local,
// global or system.
type Config interface {
	// All returns a  list of all properties associated with this configuration.
	All() []Property

	// Get attempts to retrieve the property with the specified name from the
	// current configuration, returning the property or nil if no property with
	// that name is found.
	Get(name string) Property

	// Find returns the list of all configuration properties with names matching
	// the given pattern. If the pattern ends with "*", the rest of the pattern
	// will be treated as a prefix, with Find returning all properties whose
	// name shares the prefix. If pattern does not end with "*", Find behaves as
	// Get and looks for the exact property name.
	Find(pattern string) []Property

	// String returns a string representation of the configuration, returning
	// the properties in name order.
	String() string
}

// config is the implementation of the git configuration block
type config struct {
	c   map[string]Property
	all []Property
}

// NewConfig returns the configuration instance for the list of configuration
// properties p. If p contains properties with the same name, the property
// listed last will be the property set in the returned Config instance.
func NewConfig(p []Property) Config {
	// build the name -> property lookup as well as the "all" list
	c := &config{}
	c.c = make(map[string]Property)
	for _, _p := range p {
		c.c[_p.Name()] = _p
	}

	// extract the uniquely named properties
	c.all = make([]Property, 0)
	for _, _p := range c.c {
		c.all = append(c.all, _p)
	}

	// sort the "all" list of properties
	sort.Sort(properties(c.all))

	return c
} // NewConfig()

// All returns the ordered list of all properties associated with this
// configuration.
func (c config) All() []Property {
	return c.all
} // All()

// Get attempts to retrieve the property with the specified name from the
// current configuration, returning the property or nil if no property with
// that name is found.
func (c config) Get(name string) Property {
	_property, _ok := c.c[name]
	if _ok {
		return _property
	} else {
		return nil
	}
} // Get()

// Find returns the list of all configuration properties with names matching
// the given pattern. If the pattern ends with "*", the rest of the pattern
// will be treated as a prefix, with Find returning all properties whose name
// shares the prefix. If pattern does not end with "*", Find behaves as
// Get and looks for the exact property name.
func (c config) Find(pattern string) []Property {
	// does the pattern end in "*"?
	//		- if not, then this is just a Get() call in disguise
	if !strings.HasSuffix(pattern, "*") {
		_property := c.Get(pattern)
		if _property != nil {
			return []Property{_property}
		} else {
			return []Property{}
		}
	}

	// otherwise, remove the '*' from the pattern and look for config
	// properties with names that share the resulting prefix
	_pattern := strings.TrimSuffix(pattern, "*")
	_properties := []Property{}
	for _, _property := range c.all {
		if strings.HasPrefix(_property.Name(), _pattern) {
			_properties = append(_properties, _property)
		}
	}

	return _properties
} // Find()

// String returns a string representation of the configuration, returning
// the properties in name order.
func (c config) String() string {
	_bytes := bytes.NewBuffer(nil)
	for _, _property := range c.all {
		_bytes.WriteString(
			fmt.Sprintln(_property.Name() + "=" + _property.String()),
		)
	}

	return _bytes.String()
} // String()

// ensure config conforms to the Config interface
var _ Config = &config{}
