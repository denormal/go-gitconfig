package gitconfig_test

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/denormal/go-gitconfig"
)

func TestGitConfig(t *testing.T) {
	// ensure NewGitConfig() works as expected
	//		- the current directory should be a working copy
	//		- these tests should work with current drectory and ""
	_cwd, _err := os.Getwd()
	if _err != nil {
		t.Fatalf(
			"unable to determine current working directory: %s",
			_err.Error(),
		)
	}

	for _, _path := range []string{_cwd, ""} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}
	}

	// ensure NewGitConfig() works with a temporary directory that is
	// not a part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	_config, _err := gitconfig.NewGitConfig(_dir)
	if _err != nil {
		t.Fatalf(
			"%q: unexpected error from NewGitConfig: %s",
			_dir, _err.Error(),
		)
	} else if _config == nil {
		t.Fatalf("%q: unexpected nil from NewGitConfig", _dir)
	}

} // TestGitConfig()

func TestGitConfigLocal(t *testing.T) {
	// these tests should pass whether we are in the current directory
	// (which should be a working copy), or in a temporary directory that is
	// not part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	for _, _path := range []string{"", _dir} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}

		// ensure Local() behaves as expected
		_local := _config.Local()
		if _local == nil {
			t.Fatal("unexpected nil from GitConfig.Local()")
		}

		// if we're in a working copy, then we should attempt to retrieve the
		// local configuration manually and compare the two
		_is, _error := gitconfig.IsWorkingCopy(_path)
		if _error != nil {
			if _error != gitconfig.MissingWorkingCopyError {
				t.Fatalf(
					"unexpected error from InWorkingCopy(): %s",
					_error.Error(),
				)
			}
		}
		if _is {
			_l, _error := gitconfig.NewLocalConfig(_path)
			if _error != nil {
				t.Fatalf(
					"unexpected error from NewLocalConfig(): %s",
					_error.Error(),
				)
			} else if _l == nil {
				t.Fatalf("unexpected nil Config; expected instance")
			} else if len(_l.All()) != len(_local.All()) {
				t.Fatalf(
					"config mismatch; expected %d items, found %d",
					len(_local.All()), len(_l.All()),
				)
			}
		} else if len(_local.All()) != 0 {
			t.Fatalf(
				"%q: expected no local config; found %d items",
				_path, len(_local.All()),
			)
		}
	}
} // TestGitConfigLocal()

func TestGitConfigGlobal(t *testing.T) {
	// these tests should pass whether we are in the current directory
	// (which should be a working copy), or in a temporary directory that is
	// not part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	for _, _path := range []string{"", _dir} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}

		// ensure Global() behaves as expected
		_global := _config.Global()
		if _global == nil {
			t.Fatal("unexpected nil from GitConfig.Global()")
		}

		// attempt to retrieve the global configuration manually
		// and compare the two
		_g, _error := gitconfig.NewGlobalConfig()
		if _error != nil {
			t.Fatalf(
				"unexpected error from NewGlobalConfig(): %s",
				_error.Error(),
			)
		} else if _g == nil {
			t.Fatalf("unexpected nil Config; expected instance")
		} else if len(_g.All()) != len(_global.All()) {
			t.Fatalf(
				"config mismatch; expected %d items, found %d",
				len(_global.All()), len(_g.All()),
			)
		} else {
			// ensure the retrieved values are the same
			_values := _g.All()
			for _i, _got := range _global.All() {
				_expected := _values[_i]
				if _got.Name() != _expected.Name() {
					t.Fatalf(
						"value name mismatch; expected %q, got %q at index %d",
						_expected.Name(), _got.Name(), _i,
					)
				} else if _got.String() != _expected.String() {
					t.Fatalf(
						"value name mismatch; expected %q, got %q at index %d",
						_expected.String(), _got.String(), _i,
					)
				}
			}
		}
	}
} // TestGitConfigGlobal()

func TestGitConfigSystem(t *testing.T) {
	// these tests should pass whether we are in the current directory
	// (which should be a working copy), or in a temporary directory that is
	// not part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	for _, _path := range []string{"", _dir} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}

		// ensure System() behaves as expected
		_system := _config.System()
		if _system == nil {
			t.Fatal("unexpected nil from GitConfig.System()")
		}

		// attempt to retrieve the system configuration manually
		// and compare the two
		_s, _error := gitconfig.NewSystemConfig()
		if _error != nil {
			t.Fatalf(
				"unexpected error from NewSystemConfig(): %s",
				_error.Error(),
			)
		} else if _s == nil {
			t.Fatalf("unexpected nil Config; expected instance")
		} else if len(_s.All()) != len(_system.All()) {
			t.Fatalf(
				"config mismatch; expected %d items, found %d",
				len(_system.All()), len(_s.All()),
			)
		} else {
			// ensure the retrieved values are the same
			_values := _s.All()
			for _i, _got := range _system.All() {
				_expected := _values[_i]
				if _got.Name() != _expected.Name() {
					t.Fatalf(
						"value name mismatch; expected %q, got %q at index %d",
						_expected.Name(), _got.Name(), _i,
					)
				} else if _got.String() != _expected.String() {
					t.Fatalf(
						"value name mismatch; expected %q, got %q at index %d",
						_expected.String(), _got.String(), _i,
					)
				}
			}
		}
	}
} // TestGitConfigSystem()

func TestGitConfigAll(t *testing.T) {
	// these tests should pass whether we are in the current directory
	// (which should be a working copy), or in a temporary directory that is
	// not part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	for _, _path := range []string{"", _dir} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}

		// manually retrieve the local, global & system configs
		_local, _global, _system := lgs(t, _path)

		// now, manually build the output of All()
		//		- this is an overlay of system, then global, then local
		//		  with each successive config overwriting the previous
		_map := make(map[string]gitconfig.Value)
		for _, _value := range _system.All() {
			_map[_value.Name()] = _value
		}
		for _, _value := range _global.All() {
			_map[_value.Name()] = _value
		}
		for _, _value := range _local.All() {
			_map[_value.Name()] = _value
		}

		// create the ordered list of values
		_expect := make([]gitconfig.Value, 0, len(_map))
		for _, _value := range _map {
			_expect = append(_expect, _value)
		}
		sort.Sort(gitconfig.Values(_expect))

		// ensure All() returns the expected values
		_all := _config.All()
		if len(_all) != len(_expect) {
			t.Fatalf(
				"%q: unexpected All() return; expected %d items, got %d",
				_path, len(_expect), len(_all),
			)
		} else {
			for _i, _expected := range _expect {
				_got := _all[_i]
				if _got.Name() != _expected.Name() {
					t.Fatalf(
						"unexpected name; expected %q, got %q at index %d",
						_expected.Name(), _got.Name(), _i,
					)
				} else if _got.String() != _expected.String() {
					t.Fatalf(
						"%q: unexpected All() return; expected %v, got %v",
						_expected.Name(), _expected, _got,
					)
				}

			}
		}
	}
} // TestGitConfigAll()

func TestGitConfigGet(t *testing.T) {
	// these tests should pass whether we are in the current directory
	// (which should be a working copy), or in a temporary directory that is
	// not part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	for _, _path := range []string{"", _dir} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}

		// ensure we can Get() all values returned by All()
		for _, _expected := range _config.All() {
			_got, _ok := _config.Get(_expected.Name())
			if !_ok {
				t.Fatalf(
					"%q: unexpected Get() failure; expected %v, got %v",
					_expected.Name(), true, _ok,
				)
			} else if _got == nil {
				t.Fatalf(
					"%q: unexpected nil value; expected %v",
					_expected.Name(), _expected,
				)
			} else if _got.Name() != _expected.Name() {
				t.Fatalf(
					"%q: unexpected name; expected %q, got %q",
					_expected.Name(), _expected.Name(), _got.Name(),
				)
			} else if _got != _expected {
				t.Fatalf(
					"%q: unexpected Get() return; expected %v, got %v",
					_expected.Name(), _expected, _got,
				)
			}
		}

		// ensure Get() behaves as expected when given unknown names
		for _, _value := range _config.All() {
			_get := _value.Name() + time.Now().String()
			_got, _ok := _config.Get(_get)
			if _ok {
				t.Fatalf(
					"%q: unexpected Get() success; expected %v, got %v",
					_get, false, _ok,
				)
			} else if _got != nil {
				t.Fatalf(
					"%q: unexpected Get() value; expected nil, got %v",
					_get, _got.Name(),
				)
			}
		}
	}
} // TestGitConfigGet()

func TestGitConfigFind(t *testing.T) {
	// these tests should pass whether we are in the current directory
	// (which should be a working copy), or in a temporary directory that is
	// not part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	for _, _path := range []string{"", _dir} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}

		// first, ensure Find() behaves like Get() when given an exact match
		_all := _config.All()
		for _, _value := range _all {
			_find := _config.Find(_value.Name())
			if len(_find) != 1 {
				t.Fatalf(
					"%q: unexpected Find(); expected %d results, got %d",
					_value.Name(), 1, len(_find),
				)
			} else if _find[0].Name() != _value.Name() {
				t.Fatalf(
					"%q: unexpected Find(); expected %q, got %q",
					_value.Name(), _value.Name(), _find[0].Name(),
				)
			}
		}

		// now, perform prefix searches
		_prefixes := prefixes(_all)
		for _prefix, _values := range _prefixes {
			_find := _config.Find(_prefix)
			if len(_find) != len(_values) {
				t.Fatalf(
					"%q: unexpected Find(); expected %d results, got %d",
					_prefix, len(_values), len(_find),
				)
			} else {
				for _i := 0; _i < len(_find); _i++ {
					_expected := _values[_i]
					_got := _find[_i]
					if _got.Name() != _expected.Name() {
						t.Fatalf(
							"%q: unexpected Find(); expected %q, got %q",
							_prefix, _expected.Name(), _got.Name(),
						)
					} else if _got.String() != _expected.String() {
						t.Fatalf(
							"%q: unexpected Find(); expected %q, got %q",
							_prefix, _expected.String(), _got.String(),
						)
					}
				}
			}
		}

		// finally, perform Find() with unknown prefixes
		for _p, _ := range _prefixes {
			_prefix := time.Now().String() + _p
			_find := _config.Find(_prefix)
			if len(_find) != 0 {
				t.Fatalf(
					"%q: unexpected Find(); expected %d results, got %d",
					_prefix, 0, len(_find),
				)
			}
		}
	}
} // TestGitConfigFind()

func TestGitConfigString(t *testing.T) {
	// these tests should pass whether we are in the current directory
	// (which should be a working copy), or in a temporary directory that is
	// not part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	for _, _path := range []string{"", _dir} {
		_config, _err := gitconfig.NewGitConfig(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewGitConfig: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewGitConfig", _path)
		}

		// ensure String() returns a non-empty string
		_string := _config.String()
		if _string == "" {
			t.Fatalf("%q: unexpected empty string", _path)
		}

		// manually retrieve the local, global & system configs
		_local, _global, _system := lgs(t, _path)

		// build a map of lines from each config String()
		_l := lookup(_local.String())
		_g := lookup(_global.String())
		_s := lookup(_system.String())

		// ensure every line form String() is present in one of the
		// String() outputs from the individual configs
		for _, _line := range lines(_string) {
			if _l[_line] || _g[_line] || _s[_line] {
				continue
			}

			// this line is missing from the individual configs
			t.Fatalf("%q: String() missing line: %q", _path, _line)
		}
	}
} // TestGitConfigString()

//
// helper methods
//

func lgs(t *testing.T, path string) (gitconfig.Config, gitconfig.Config, gitconfig.Config) {
	// retrieve the local, global & system configs
	//		- if we're not in a working copy, then we don't attempt to
	//		  retrieve the local configuration
	_local := gitconfig.NewConfig(nil)
	_is, _error := gitconfig.IsWorkingCopy(path)
	if _error != nil {
		if _error != gitconfig.MissingWorkingCopyError {
			t.Fatalf(
				"unexpected error from InWorkingCopy(): %s",
				_error.Error(),
			)
		}
	} else if _is {
		_local, _error = gitconfig.NewLocalConfig(path)
		if _error != nil {
			t.Fatalf(
				"unexpected error from NewLocalConfig(): %s",
				_error.Error(),
			)
		} else if _local == nil {
			t.Fatalf("unexpected nil Config; expected instance")
		}
	}

	_global, _error := gitconfig.NewGlobalConfig()
	if _error != nil {
		t.Fatalf(
			"unexpected error from NewGlobalConfig(): %s",
			_error.Error(),
		)
	} else if _global == nil {
		t.Fatalf("unexpected nil Config; expected instance")
	}

	_system, _error := gitconfig.NewSystemConfig()
	if _error != nil {
		t.Fatalf(
			"unexpected error from NewSystemConfig(): %s",
			_error.Error(),
		)
	} else if _system == nil {
		t.Fatalf("unexpected nil Config; expected instance")
	}

	// return the local, global and system configurations
	return _local, _global, _system
} // lgs()

func lines(content string) []string {
	// split the content into lines
	//		- first on "\n"
	//		- then remove dangling "\r" (if any)
	_lines := strings.Split(content, "\n")
	for _i := 0; _i < len(_lines); _i++ {
		_lines[_i] = strings.TrimSpace(_lines[_i])
	}

	return _lines
} // lines()

func lookup(content string) map[string]bool {
	_lines := lines(content)
	_map := make(map[string]bool)
	for _, _line := range _lines {
		_map[_line] = true
	}

	return _map
} // lookup()

func prefixes(values []gitconfig.Value) map[string][]gitconfig.Value {
	_map := make(map[string][]gitconfig.Value)
	for _, _value := range values {
		_prefix := prefix(_value.Name())
		_, _ok := _map[_prefix]
		if !_ok {
			_map[_prefix] = make([]gitconfig.Value, 0)
		}
		_map[_prefix] = append(_map[_prefix], _value)
	}

	// sort the expected return values
	for _, _list := range _map {
		sort.Sort(gitconfig.Values(_list))
	}

	return _map
} // prefixes()

func prefix(name string) string {
	_parts := strings.Split(name, ".")
	return _parts[0] + ".*"
} // prefix()
