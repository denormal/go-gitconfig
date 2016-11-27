package gitconfig_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/denormal/go-gitconfig"
)

func TestHasGit(t *testing.T) {
	// do we have git installed?
	_has, _ := gitconfig.HasGit()
	if !_has {
		t.Skip("git not installed")
	}

	// amend the current PATH to be emtpy
	_path := os.Getenv("PATH")
	defer os.Setenv("PATH", _path)
	_err := os.Setenv("PATH", "")
	if _err != nil {
		t.Fatalf("unable to reset PATH: %s", _err.Error())
	}

	// HasGit should now fail
	_has, _err = gitconfig.HasGit()
	if _err == nil {
		t.Error("expected error; none found")
	} else {
		_error, _ok := _err.(*exec.Error)
		if !_ok || _error.Err != exec.ErrNotFound {
			t.Errorf("expected ErrNotFound; got %q", _err.Error())
		}
	}

	// ensure false is returned
	if _has {
		t.Error("expected false for HasGit() ; true found")
	}
} // TestHasGit()

func TestIsWorkingCopy(t *testing.T) {
	// we should be in a working copy
	_cwd, _err := os.Getwd()
	if _err != nil {
		t.Fatalf(
			"unable to determine current working directory: %s",
			_err.Error(),
		)
	}
	// using the current working directory or "" should give the same results
	for _, _path := range []string{_cwd, ""} {
		_is, _err := gitconfig.IsWorkingCopy(_path)
		if _err != nil {
			if _err != gitconfig.MissingWorkingCopyError {
				t.Errorf("%q: unexpected error: %s", _path, _err.Error())
			}
		} else if !_is {
			t.Errorf("%q: expected to be in working copy; none found", _path)
		}
	}

	// create a temporary directory
	//		- this should not be a working copy
	_dir, _err := ioutil.TempDir("", "")
	if _err != nil {
		t.Fatalf("unable to create temporary directory: %s", _err.Error())
	}
	defer os.RemoveAll(_dir)

	// ensure this is not a working directory
	_is, _err := gitconfig.IsWorkingCopy(_dir)
	if _err == nil {
		t.Errorf(
			"%q: expected error: %s",
			_dir, gitconfig.MissingWorkingCopyError.Error(),
		)
	} else if _is {
		t.Errorf(
			"%q: expected not to be in working copy; working copy found",
			_dir,
		)
	} else if _err != gitconfig.MissingWorkingCopyError {
		t.Errorf("%q: unexpected error: %s", _dir, _err.Error())
	}
} // TestIsWorkingCopy()

func TestInWorkingCopy(t *testing.T) {
	// we should be in a working copy
	//		- if we're not, ensure it's correctly reported
	_in, _err := gitconfig.InWorkingCopy()
	if _err != nil {
		if _err != gitconfig.MissingWorkingCopyError {
			t.Errorf("unexpected error: %s", _err.Error())
		}
	} else if !_in {
		t.Errorf("expected to be in working copy; none found")
	}
} // TestInWorkingCopy()

func TestGitDir(t *testing.T) {
	// we should be in a working copy
	_cwd, _err := os.Getwd()
	if _err != nil {
		t.Fatalf(
			"unable to determine current working directory: %s",
			_err.Error(),
		)
	}

	// using the current working directory or "" should give the same results
	for _, _path := range []string{_cwd, ""} {
		// only perform this test if this is not a working copy
		_is, _err := gitconfig.IsWorkingCopy(_path)
		if _err != nil {
			if _err != gitconfig.MissingWorkingCopyError {
				t.Fatalf("unexpected error: %s", _err.Error())
			}
		}

		if !_is {
			t.Skipf("%q: not a working copy", _path)
		} else {
			_dir, _err := gitconfig.GitDir(_path)
			if _err != nil {
				t.Errorf("%q: unexpected error: %s", _path, _err.Error())
			} else if _dir == "" {
				t.Errorf(
					"%q: expected to be in working copy; none found",
					_path,
				)
			}
		}
	}
} // TestGitDir()

func TestWorkingCopy(t *testing.T) {
	// we should be in a working copy
	_cwd, _err := os.Getwd()
	if _err != nil {
		t.Fatalf(
			"unable to determine current working directory: %s",
			_err.Error(),
		)
	}

	// using the current working directory or "" should give the same results
	for _, _path := range []string{_cwd, ""} {
		// only perform this test if this is not a working copy
		_is, _err := gitconfig.IsWorkingCopy(_path)
		if _err != nil {
			if _err != gitconfig.MissingWorkingCopyError {
				t.Errorf("unexpected error: %s", _err.Error())
			}
		}

		if !_is {
			t.Skipf("%q: not a working copy", _path)
		} else {
			_wc, _err := gitconfig.WorkingCopy(_path)
			if _err != nil {
				t.Errorf("%q: unexpected error: %s", _path, _err.Error())
			} else if _wc == "" {
				t.Errorf(
					"%q: expected to be in working copy; none found",
					_path,
				)
			}
		}
	}
} // TestWorkingCopy()
