package gitconfig

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
	Config

	local  Config
	system Config
	global Config
}

// New returns a GitConfig instance representing the git working copy in
// the current working directory. If the current working directory is not
// part of a git working copy, the local configuration of the GitConfig will
// be nil. An error is returned if there is a problem accessing the current
// working directory, or if there is a problem parsing the configuration
// properties.
func New() (GitConfig, error) {
	return NewWithPath("")
} // New()

// NewWithPath returns a GitConfig instance representing the git working copy
// path. If path is not part of a git working copy, the local configuration
// of the GitConfig will be nil. An error is returned if there is a problem
// accessing the path, or if there is a problem parsing the configuration
// properties.
//
// If path is "", the current working directory of the process will be used.
func NewWithPath(path string) (GitConfig, error) {
	var _local Config

	// are we in a git repository?
	_is, _err := IsWorkingCopy(path)
	if _err != nil {
		if _err != MissingWorkingCopyError {
			return nil, _err
		}
	}

	// if we're in a git repository, attempt to load the local configuration
	if _is {
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

	// generate the combined properties
	//		- the properties are prioritised such that local properties
	//		  override global properties, that override system properties
	_all := []Property{}
	_all = append(_all, _system.All()...)
	_all = append(_all, _global.All()...)
	if _local != nil {
		_all = append(_all, _local.All()...)
	}
	_config := NewConfig(_all)

	return &gc{_config, _local, _system, _global}, nil
} // New()

// Local returns the local git configuration for the git working copy.
// If this GitConfig does not represent a working copy, Local will return nil.
func (g gc) Local() Config { return g.local }

// System returns the system git configuration.
func (g gc) System() Config { return g.system }

// Global returns the global git configuration for the current user.
func (g gc) Global() Config { return g.global }

// ensure gc implemented GitConfig
var _ GitConfig = &gc{}
