package gitconfig_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/denormal/go-gitconfig"
)

func TestGit(t *testing.T) {
	_git, _err := gitconfig.Git()
	if _err != nil {
		t.Fatalf("git executable not found: %s", _err.Error())
	} else if _git == "" {
		t.Fatalf("git executable found, but path empty: %q", _git)
	}

	// ensure the git path exists and is executable
	_info, _err := os.Stat(_git)
	if _err != nil {
		t.Fatalf("git executable could not be inspected: %q", _err.Error())
	} else if _info.IsDir() {
		t.Fatalf("git path is directory not file: %s", _git)
	} else if _info.Mode()|0100 == 0 {
		t.Fatalf("git path is not executable by caller: %s", _git)
	}

	// amend the current PATH to be empty
	_path := os.Getenv("PATH")
	defer os.Setenv("PATH", _path)
	_err = os.Setenv("PATH", "")
	if _err != nil {
		t.Fatalf("unable to reset PATH: %s", _err.Error())
	}

	// run the test again to ensure Git() fails
	_git, _err = gitconfig.Git()
	if _err == nil {
		t.Error("expected error looking for git; none found")
	} else {
		_error, _ok := _err.(*exec.Error)
		if !_ok || _error.Err != exec.ErrNotFound {
			t.Errorf("expected ErrNotFound; got %q", _err.Error())
		}
	}

	// ensure the path to git is empty
	if _git != "" {
		t.Error("expected empty path from Git(); found %q", _git)
	}
} // TestGit()

func TestHasGit(t *testing.T) {
	// do we have git installed?
	if !gitconfig.HasGit() {
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
	if gitconfig.HasGit() {
		t.Error("unexpected success: HasGit() with PATH %q", os.Getenv("PATH"))
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
