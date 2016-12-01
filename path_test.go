package gitconfig_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/denormal/go-gitconfig"
	"github.com/denormal/go-gittools"
)

func TestPath(t *testing.T) {
	// skip this test if git is not installed
	if !gittools.HasGit() {
		t.Skip("git not installed")
	}

	// ensure New() works as expected
	//		- the current directory should be a working copy
	//		- these tests should work with current drectory and ""
	_cwd, _err := os.Getwd()
	if _err != nil {
		t.Fatalf(
			"unable to determine current working directory: %s",
			_err.Error(),
		)
	}

	// ensure Path() returns the expected value
	for _, _path := range []string{_cwd, ""} {
		_config, _err := gitconfig.NewWithPath(_path)
		if _err != nil {
			t.Fatalf(
				"%q: unexpected error from NewWithPath: %s",
				_path, _err.Error(),
			)
		} else if _config == nil {
			t.Fatalf("%q: unexpected nil from NewWithPath", _path)
		} else if _config.Path() != _cwd {
			t.Fatalf(
				"path mismatch: expected %q, got %q",
				_cwd, _config.Path(),
			)
		}
	}

	// ensure New() works with a temporary directory that is
	// not a part of a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf(
			"unable to create temporary directory: %s",
			_err.Error(),
		)
	}
	defer os.RemoveAll(_dir)

	_config, _err := gitconfig.NewWithPath(_dir)
	if _err != nil {
		t.Fatalf(
			"%q: unexpected error from NewWithPath: %s",
			_dir, _err.Error(),
		)
	} else if _config == nil {
		t.Fatalf("%q: unexpected nil from NewWithPath", _dir)
	} else if _config.Path() != _dir {
		t.Fatalf(
			"path mismatch: expected %q, got %q",
			_dir, _config.Path(),
		)
	}
} // TestPath()
