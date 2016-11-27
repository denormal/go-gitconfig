package gitconfig

// Error is the interface for configuration errors
type Error interface {
	error

	// Property() returns the name of the property this error is related to,
	// if any. If this error is not related to a property, Property() returns
	// the empty string.
	Property() string

	// Underlying returns the underlying error.
	Underlying() error
}

// err is the implementation of Error
type err struct {
	property string
	e        error
}

// NewError returns a configuration error related representing error e relating
// to the named property.
func NewError(property string, e error) Error {
	return &err{property, e}
}

// Property() returns the name of the property this error is related to,
// if any. If this error is not related to a property, Property() returns
// the empty string.
func (e err) Property() string { return e.property }

// Underlying returns the underlying error.
func (e err) Underlying() error { return e.e }

// Error returns the error string for this Error.
func (e err) Error() string { return e.e.Error() }

// ensure err satisfies Error
var _ Error = &err{}
