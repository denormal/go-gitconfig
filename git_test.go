package gitconfig_test

import (
	"testing"

	"github.com/denormal/go-gitconfig"
	"github.com/denormal/go-gittools"
)

func TestNewLocalConfig(t *testing.T) {
	// do we have git installed?
	if !gittools.HasGit() {
		t.Skip("git not installed")
	}

	// are we in a git repository?
	_is, _ := gittools.InWorkingCopy()
	if _is {
		_config, _err := gitconfig.NewLocalConfig("")
		if _err != nil {
			t.Fatalf(
				"unexpected error from NewLocalConfig(): %s",
				_err.Error(),
			)
		} else if _config == nil {
			t.Fatal(
				"expected non-empty return from NewLocalConfig(); nil found",
			)
		}
	} else {
		t.Skip("working directory not within Git repository")
	}
} // TestNewLocalConfig()

func TestNewGlobalConfig(t *testing.T) {
	// do we have git installed?
	if !gittools.HasGit() {
		t.Skip("git not installed")
	}

	// attempt to load the global configuration
	_config, _err := gitconfig.NewGlobalConfig()
	if _err != nil {
		t.Fatalf("unexpected error from NewGlobalConfig(): %s", _err.Error())
	} else if _config == nil {
		t.Fatal("expected non-empty return from NewGlobalConfig(); nil found")
	}
} // TestNewGlobalConfig()

func TestNewSystemConfig(t *testing.T) {
	// do we have git installed?
	if !gittools.HasGit() {
		t.Skip("git not installed")
	}

	// attempt to load the system configuration
	_config, _err := gitconfig.NewSystemConfig()
	if _err != nil {
		t.Fatalf("unexpected error from NewSystemConfig(): %s", _err.Error())
	} else if _config == nil {
		t.Fatal("expected non-empty return from NewSystemConfig(); nil found")
	}
} // TestNewSystemConfig()
