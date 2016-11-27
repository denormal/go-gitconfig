package gitconfig

import (
	"bytes"
	"fmt"
	"sort"
)

type GitConfig interface {
	Config
	Local() Config
	System() Config
	Global() Config
}

type gc struct {
	local  Config
	system Config
	global Config
	all    []Property
}

func NewGitConfig(path string) (GitConfig, Error) {
	var (
		_local Config
		_err   Error
	)

	// are we in a git repository?
	_is, _error := IsWorkingCopy(path)
	if _err != nil {
		if _err != MissingWorkingCopyError {
			return nil, NewError("", _error)
		}
	}

	if !_is {
		_local = NewConfig(nil)
	} else {
		_local, _err = NewLocalConfig(path)
		if _err != nil {
			return nil, _err
		}
	}
	_system, _err := NewSystemConfig()
	if _err != nil {
		return nil, _err
	}
	_global, _err := NewGlobalConfig()
	if _err != nil {
		return nil, _err
	}

	return &gc{_local, _system, _global, nil}, nil
} // NewGitConfig()

func (g gc) Local() Config  { return g.local }
func (g gc) System() Config { return g.system }
func (g gc) Global() Config { return g.global }

func (g gc) All() []Property {
	if g.all == nil {
		_all := make([]Property, 0)
		_all = append(_all, g.Local().All()...)
		_all = append(_all, g.Global().All()...)
		_all = append(_all, g.System().All()...)

		// generate a unique list of properties
		_rtn := make([]Property, 0)
		_map := make(map[string]bool)
		for _, _property := range _all {
			if !_map[_property.Name()] {
				_rtn = append(_rtn, _property)
				_map[_property.Name()] = true
			}
		}

		sort.Sort(Properties(_rtn))
		g.all = _rtn
	}

	return g.all
} // All()

func (g gc) Get(name string) (Property, bool) {
	_property, _ok := g.Local().Get(name)
	if !_ok {
		_property, _ok = g.Global().Get(name)
		if !_ok {
			_property, _ok = g.System().Get(name)
		}
	}

	return _property, _ok
} // Get()

func (g gc) Find(pattern string) []Property {
	_config := make([]Property, 0)
	_config = append(_config, g.Local().Find(pattern)...)
	_config = append(_config, g.Global().Find(pattern)...)
	_config = append(_config, g.System().Find(pattern)...)

	// generate the unique list of properties
	_rtn := make([]Property, 0)
	_map := make(map[string]bool)
	for _, _property := range _config {
		if !_map[_property.Name()] {
			_rtn = append(_rtn, _property)
			_map[_property.Name()] = true
		}
	}

	// sort the return properties to make Find() predictable
	sort.Sort(Properties(_rtn))

	return _rtn
} // Find()

func (g gc) String() string {
	_bytes := bytes.NewBuffer(nil)
	for _, _property := range g.All() {
		_bytes.WriteString(
			fmt.Sprintln(_property.Name() + "=" + _property.String()),
		)
	}

	return _bytes.String()
} // String()
