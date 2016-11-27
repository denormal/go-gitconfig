package gitconfig

import (
	"os"
	"os/exec"
	"path/filepath"
)

// execute attempts to execute the git command from the path directory (or
// the parent of path if path represents a file), with args arguments. If
// git is executed successfully, the slice of bytes output to STDOUT by git
// will be returned, otherwise an error will be returned.
func execute(path string, args []string) ([]byte, error) {
	// is git installed?
	_has, _err := HasGit()
	if _err != nil {
		return nil, _err
	} else if !_has {
		return nil, MissingGitError
	}

	// if we have a path, attempt to change into it before executing
	// the git command
	if path != "" {
		var _dir string

		// do we have a file or a directory?
		_info, _err := os.Stat(path)
		if _err != nil {
			return nil, _err
		} else if _info.IsDir() {
			_dir = path
		} else {
			_dir, _ = filepath.Split(path)
		}

		_cwd, _err := os.Getwd()
		if _err != nil {
			return nil, _err
		}
		defer os.Chdir(_cwd)

		// attempt to change into the given path
		_err = os.Chdir(_dir)
		if _err != nil {
			return nil, _err
		}
	}

	// execute the git command
	return exec.Command(_EXE, args...).Output()
} // execute()
