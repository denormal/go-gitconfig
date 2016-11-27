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

	all []Value
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

func (g gc) All() []Value {
	if g.all == nil {
		_all := make([]Value, 0)
		_all = append(_all, g.Local().All()...)
		_all = append(_all, g.Global().All()...)
		_all = append(_all, g.System().All()...)

		// generate a unique list of values
		_rtn := make([]Value, 0)
		_map := make(map[string]bool)
		for _, _value := range _all {
			if !_map[_value.Name()] {
				_rtn = append(_rtn, _value)
				_map[_value.Name()] = true
			}
		}

		sort.Sort(Values(_rtn))
		g.all = _rtn
	}

	return g.all
} // All()

func (g gc) Get(name string) (Value, bool) {
	_value, _ok := g.Local().Get(name)
	if !_ok {
		_value, _ok = g.Global().Get(name)
		if !_ok {
			_value, _ok = g.System().Get(name)
		}
	}

	return _value, _ok
} // Get()

func (g gc) Find(pattern string) []Value {
	_config := make([]Value, 0)
	_config = append(_config, g.Local().Find(pattern)...)
	_config = append(_config, g.Global().Find(pattern)...)
	_config = append(_config, g.System().Find(pattern)...)

	// generate the unique list of values
	_rtn := make([]Value, 0)
	_map := make(map[string]bool)
	for _, _value := range _config {
		if !_map[_value.Name()] {
			_rtn = append(_rtn, _value)
			_map[_value.Name()] = true
		}
	}

	// sort the return values to make Find() predictable
	sort.Sort(Values(_rtn))

	return _rtn
} // Find()

func (g gc) String() string {
	_bytes := bytes.NewBuffer(nil)
	for _, _value := range g.All() {
		_bytes.WriteString(
			fmt.Sprintln(_value.Name() + "=" + _value.String()),
		)
	}

	return _bytes.String()
} // String()
