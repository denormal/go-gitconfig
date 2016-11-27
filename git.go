package gitconfig

import (
	"strings"
)

var _CONFIG = []string{"config", "--list"}

func NewLocalConfig(path string) (Config, Error) {
	// attempt to retrieve the local configuration
	_properties, _err := get(path, "--local")
	if _err != nil {
		return nil, _err
	}

	return NewConfig(_properties), nil
} // getLocal()

func NewSystemConfig() (Config, Error) {
	_properties, _err := get("", "--system")
	if _err != nil {
		return nil, _err
	}

	return NewConfig(_properties), nil
} // NewSystemConfig()

func NewGlobalConfig() (Config, Error) {
	_properties, _err := get("", "--global")
	if _err != nil {
		return nil, _err
	}

	return NewConfig(_properties), nil
} // NewGlobalConfig()

//
// private functions
//

func get(path, flag string) ([]Property, Error) {
	// add the flag to the argument list
	_args := _CONFIG
	_args = append(_args, flag)

	// attempt to execute the "git config" command
	_output, _err := execute(path, _args)
	if _err != nil {
		return nil, NewError("", _err)
	}

	// parse the configuration output into properties
	_lines := strings.Split(string(_output), "\n")
	_rtn := make([]Property, 0, len(_lines))
	for _, _line := range _lines {
		// remove possible "\r" line ending
		_line = strings.TrimSpace(_line)
		_parts := strings.SplitN(_line, "=", 2)
		switch len(_parts) {
		case 1: // we have a name only
			if _parts[0] != "" {
				_property, _err := NewProperty(_parts[0], "")
				if _err != nil {
					return nil, _err
				}
				_rtn = append(_rtn, _property)
			}
		case 2: // we have a name and property
			if _parts[0] != "" {
				_property, _err := NewProperty(_parts[0], _parts[1])
				if _err != nil {
					return nil, _err
				}
				_rtn = append(_rtn, _property)
			}
		}
	}

	return _rtn, nil
} // get()
