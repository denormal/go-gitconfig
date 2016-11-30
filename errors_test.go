package gitconfig_test

import (
	"testing"

	"github.com/denormal/go-gitconfig"
	"github.com/denormal/go-gittools"
)

func TestErrors(t *testing.T) {
	if gitconfig.MissingGitError != gittools.MissingGitError {
		t.Fatal(
			"expected gittools.MissingGitError to be the same as " +
				"gitconfig.MissingGitError; difference found",
		)
	}
	if gitconfig.MissingWorkingCopyError != gittools.MissingWorkingCopyError {
		t.Fatal(
			"expected gittools.MissingWorkingCopyError to be the same as " +
				"gitconfig.MissingWorkingCopyError; difference found",
		)
	}
} // TestErrors()
