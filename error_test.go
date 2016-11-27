package gitconfig_test

import (
	"errors"
	"testing"

	"github.com/denormal/go-gitconfig"
)

const (
	_ERROR    = "this is my test error"
	_PROPERTY = "this is my property"
)

func TestError(t *testing.T) {
	_err := errors.New(_ERROR)
	_error := gitconfig.NewError(_PROPERTY, _err)

	// ensure the property can be retrieved
	if _error.Property() != _PROPERTY {
		t.Errorf(
			"error property mismatch; expected %q, got %q",
			_PROPERTY, _error.Property(),
		)
	}

	// ensure the underlying error can be retrieved
	if _error.Underlying() != _err {
		t.Errorf(
			"underlying error mismatch; expected %v, got %v",
			_err, _error.Underlying(),
		)
	}

	// ensure the error message can be retrieved
	if _error.Error() != _ERROR {
		t.Errorf(
			"error mismatch; expected %q, got %q",
			_ERROR, _error.Error(),
		)
	}
} // TestError()
