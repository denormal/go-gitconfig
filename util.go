package gitconfig

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

var (
	_EXE       = "git"
	_DIR       = []string{"rev-parse", "--git-dir"}
	_WORKING   = []string{"rev-parse", "--show-toplevel"}
	_ISWORKING = []string{"rev-parse", "--is-inside-work-tree"}

	MissingGitError         = errors.New("git executable not found")
	MissingWorkingCopyError = errors.New("git working copy not found")
)

// HasGit returns true if the host system has git installed, and if the
// git executable is located within the current user's PATH.
func HasGit() (bool, error) {
	_, _err := exec.LookPath(_EXE)
	if _err != nil {
		return false, _err
	}

	return true, nil
} // HasGit()

// InWorkingCopy returns true if the current directory is within a git
// working copy.
func InWorkingCopy() (bool, error) {
	_cwd, _err := os.Getwd()
	if _err != nil {
		return false, _err
	}

	// is the current directory in a git working copy?
	return IsWorkingCopy(_cwd)
} // InWorkingCopy()

// IsWorkingCopy returns true if path is within a git working copy.
func IsWorkingCopy(path string) (bool, error) {
	_output, _err := execute(path, _ISWORKING)
	if _err == nil {
		_lines := strings.Split(string(_output), "\n")
		for _, _line := range _lines {
			_line = strings.TrimSpace(_line)
			if _line != "true" {
				// we have a git working copy
				return true, nil
			}
		}

		// we don't have a working copy
		return false, MissingWorkingCopyError
	}

	// do we have git installed?
	//		- if we do, then we interpret the error as a missing working copy
	//		- this is a little dangerous, as other problems could be the cause
	//		  of the error, such as changes to the "git" API
	//		- however, this does give a better user experience at present
	//		- interrogating child process exist codes is difficult across
	//		  platforms, so for now we take this simplistic approach
	//		- it also saves having a dependency on "git" exit codes
	_has, _ := HasGit()
	if _has {
		return false, MissingWorkingCopyError
	}

	// either we encountered an error, or we don't have a working copy
	return false, _err
} // IsWorkingCopy()

// WorkingCopy returns the root of the working copy to which path belongs, if
// path is within a git working copy.
func WorkingCopy(path string) (string, error) {
	_output, _err := execute(path, _WORKING)
	if _err == nil {
		_lines := strings.Split(string(_output), "\n")
		for _, _line := range _lines {
			_line = strings.TrimSpace(_line)
			if _line != "" {
				// we have a git working copy
				return _line, nil
			}
		}

		// we don't have a working copy
		return "", MissingWorkingCopyError
	}

	// do we have git installed?
	//		- if we do, then we interpret the error as a missing working copy
	_has, _ := HasGit()
	if _has {
		return "", MissingWorkingCopyError
	}

	// either we encountered an error, or we don't have a working copy
	return "", _err
} // IsWorkingCopy()

// GitDir returns the git directory for the working copy in which path is
// located, or an error if the path cannot be resolved, or is not located
// within a working copy.
func GitDir(path string) (string, error) {
	// attempt to resolve the .git directory within the given path hierarchy
	_output, _err := execute(path, _DIR)
	if _err == nil {
		_lines := strings.Split(string(_output), "\n")
		for _, _line := range _lines {
			_line = strings.TrimSpace(_line)
			if _line != "" {
				return _line, nil
			}
		}
		return "", MissingWorkingCopyError
	}

	// we could not determine if we are in a git working copy
	return "", _err
} // GitDir()
