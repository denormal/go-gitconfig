package gitconfig

type Error interface {
	error
	Property() string
	Underlying() error
}

type err struct {
	property string
	e        error
}

func NewError(property string, e error) Error {
	return &err{property, e}
}

func (e err) Property() string  { return e.property }
func (e err) Underlying() error { return e.e }
func (e err) Error() string     { return e.e.Error() }

// ensure err satisfies Error
var _ Error = &err{}
