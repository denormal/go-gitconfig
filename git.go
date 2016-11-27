package gitconfig

import (
	"strings"
)

var _CONFIG = []string{"config", "--list"}

// NewLocalConfig returns the Config instance for the local git configuration,
// for the repository represented by path. If there is a problem extracting
// this configuration, and Error is returned. If path is "", the current
// working directory of the process will be used.
func NewLocalConfig(path string) (Config, error) {
	return get(path, "--local")
} // getLocal()

// NewSystemConfig returns the Config instance for the system git configuration.
// If there is a problem extracting this configuration, and Error is returned.
func NewSystemConfig() (Config, error) {
	return get("", "--system")
} // NewSystemConfig()

// NewGlobalConfig returns the Config instance for the global git configuration.
// If there is a problem extracting this configuration, and Error is returned.
func NewGlobalConfig() (Config, error) {
	return get("", "--global")
} // NewGlobalConfig()

//
// private functions
//

// get returns the list of configuration properties for the "git config"
// command executed in the given path with the supplied flag. An Error is
// returned if there is a problem executing git, or parsing a property.
func get(path, flag string) (Config, error) {
	// add the flag to the argument list
	_args := _CONFIG
	_args = append(_args, flag)

	// attempt to execute the "git config" command
	_output, _err := execute(path, _args)
	if _err != nil {
		return nil, _err
	}

	// parse the configuration output into properties
	_lines := strings.Split(string(_output), "\n")
	_properties := make([]Property, 0, len(_lines))
	for _, _line := range _lines {
		// remove possible "\r" line ending
		_line = strings.TrimSpace(_line)
		_parts := strings.SplitN(_line, "=", 2)
		switch len(_parts) {
		case 1: // we have a name only
			if _parts[0] != "" {
				_property := NewProperty(_parts[0], "")
				_properties = append(_properties, _property)
			}
		case 2: // we have a name and property
			if _parts[0] != "" {
				_property := NewProperty(_parts[0], _parts[1])
				_properties = append(_properties, _property)
			}
		}
	}

	return NewConfig(_properties), nil
} // get()
