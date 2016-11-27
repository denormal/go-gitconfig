package gitconfig

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type Config interface {
	All() []Value
	Get(string) (Value, bool)
	Find(string) []Value
	String() string
}

type config struct {
	c   map[string]Value
	all []Value
}

func NewConfig(v []Value) Config {
	c := &config{}
	c.c = make(map[string]Value)
	c.all = make([]Value, 0)
	for _, _v := range v {
		c.c[_v.Name()] = _v
		c.all = append(c.all, _v)
	}

	sort.Sort(Values(c.all))

	return c
} // NewConfig()

func (c config) All() []Value {
	return c.all
} // All()

func (c config) Get(name string) (Value, bool) {
	_value, _ok := c.c[name]
	return _value, _ok
} // Get()

func (c config) Find(pattern string) []Value {
	// does the pattern end in "*"?
	//		- if not, then this is just a Get() call in disguise
	if !strings.HasSuffix(pattern, "*") {
		_value, _ok := c.Get(pattern)
		if _ok {
			return []Value{_value}
		} else {
			return []Value{}
		}
	}

	// otherwise, remove the '*' from the pattern and look for config
	// values with names that share the resulting prefix
	_pattern := strings.TrimSuffix(pattern, "*")
	_values := []Value{}
	for _, _value := range c.all {
		if strings.HasPrefix(_value.Name(), _pattern) {
			_values = append(_values, _value)
		}
	}

	return _values
} // Find()

func (c config) String() string {
	_bytes := bytes.NewBuffer(nil)
	for _, _value := range c.all {
		_bytes.WriteString(
			fmt.Sprintln(_value.Name() + "=" + _value.String()),
		)
	}

	return _bytes.String()
} // String()
