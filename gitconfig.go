package gitconfig

import (
	"bytes"
	"fmt"
	"sort"
)

// GitConfig is the interface to git configuration, encompassing local, global
// and system configuration for a git working copy.
type GitConfig interface {
	Config

	// Local returns the local git configuration for the git working copy.
	Local() Config

	// System returns the system git configuration.
	System() Config

	// Global returns the global git configuration for the current user.
	Global() Config
}

// gc is the implementation of the GitConfig interface
type gc struct {
	local  Config
	system Config
	global Config
	all    []Property
}

// NewGitConfig returns a GitConfig instance representing the git working copy
// path. If path is not part of a git working copy, the local configuration
// of the GitConfig will be empty. An Error is returned if there is a problem
// accessing the path, or if there is a a problem parsing the configuration
// properties.
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

// Local returns the local git configuration for the git working copy.
func (g gc) Local() Config { return g.local }

// System returns the system git configuration.
func (g gc) System() Config { return g.system }

// Global returns the global git configuration for the current user.
func (g gc) Global() Config { return g.global }

// All returns the ordered list of all properties associated with this
// configuration. The GitConfig properties are comprised of the local,
// global and system properties, with local properties taking priority over
// global properties, and global properties taking priority over system
// properties.
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

// Get attempts to retrieve the property with the specified name from the
// current configuration, returning the property and ok set to true if it
// exists. Otherwise, ok will be false.
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

// Find returns the list of all configuration properties with names matching
// the given pattern. If the pattern ends with "*", the rest of the pattern
// will be treated as a prefix, with Find returning all properties whose name
// shares the prefix. If pattern does not end with "*", Find behaves as
// Get and looks for the exact property name.
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

// String returns a string representation of the configuration, returning
// the properties in name order.
func (g gc) String() string {
	_bytes := bytes.NewBuffer(nil)
	for _, _property := range g.All() {
		_bytes.WriteString(
			fmt.Sprintln(_property.Name() + "=" + _property.String()),
		)
	}

	return _bytes.String()
} // String()

// ensure gc implemented GitConfig
var _ GitConfig = &gc{}
