package gitconfig_test

import (
	"testing"

	"github.com/denormal/go-gitconfig"
)

type ptest struct {
	n   string
	p   gitconfig.Property
	b   *bool
	i   *int
	s   string
	l   []string
	e   error
	err gitconfig.Error
}

func (p *ptest) Test(t *testing.T) {
	// do we expect an error?
	if p.e != nil {
		if p.err == nil {
			t.Errorf(
				"%q: expected New error %q; none found",
				p.n, p.e.Error(),
			)
		} else if p.err.Underlying() != p.e {
			t.Errorf(
				"%q: New error mismatch; expected %q, got %q",
				p.n, p.e.Error(), p.err.Error(),
			)
		}
	} else {

		// ensure the name is as expected
		if p.p.Name() != p.n {
			t.Errorf(
				"%q: unexpected name; expected %q, got %q",
				p.n, p.n, p.p.Name(),
			)
		}

		// ensure the string form is as expected
		if p.p.String() != p.s {
			t.Errorf(
				"%q: unexpected string; expected %q, got %q",
				p.n, p.s, p.p.String(),
			)
		}

		// do we expect a boolean value?
		if p.b != nil {
			_bool, _ok := p.p.Bool()
			if !_ok {
				t.Errorf(
					"%q: boolean failure; expected %v, got %v",
					p.n, true, _ok,
				)
			} else if *p.b != _bool {
				t.Errorf(
					"%q: boolean mismatch; expected %v, got %v",
					*p.b, _bool,
				)
			}
		} else {
			_bool, _ok := p.p.Bool()
			if _ok {
				t.Errorf(
					"%q: unexpected boolean success; expected %v, got %v",
					p.n, false, _ok,
				)
			} else if _bool != false {
				t.Errorf(
					"%q: unexpected boolean return; expected %v, got %v",
					p.n, false, _bool,
				)
			}
		}

		// do we expect an integer value?
		if p.i != nil {
			_int, _ok := p.p.Int()
			if !_ok {
				t.Errorf(
					"%q: integer failure; expected %v, got %v",
					p.n, true, _ok,
				)
			} else if *p.i != _int {
				t.Errorf(
					"%q: integer mismatch; expected %v, got %v",
					*p.i, _int,
				)
			}
		} else {
			_int, _ok := p.p.Int()
			if _ok {
				t.Errorf(
					"%q: unexpected integer success; expected %v, got %v",
					p.n, false, _ok,
				)
			} else if _int != 0 {
				t.Errorf(
					"%q: unexpected integer return; expected %v, got %v",
					p.n, 0, _int,
				)
			}
		}

		// do we expect a list value?
		if p.i != nil {
			_list := p.p.List()
			if len(p.l) != len(_list) {
				t.Errorf(
					"%q: list length mismatch; expected %v, got %v",
					p.n, len(p.l), len(_list),
				)
			} else {
				for _i, _v := range p.l {
					if _list[_i] != _v {
						t.Errorf(
							"%q: list item mismatch; "+
								"expected %v, got %v for item %d",
							p.n, _v, _list[_i], _i,
						)
					}
				}
			}
		}
	}
} // Test()

var (
	// test values
	_true  = true
	_false = false
	_1     = 1
	_0     = 0
	__123  = -123
	_1234  = 1234

	// define the property tests
	//		- each test must have a unique name
	//		- tests are expected to start with X. to assist config_test.go
	_PROPERTIES = []*ptest{
		// _P( name, property, bool, int, list, error )
		_P("p.a", "?", nil, nil, nil, nil),
		_P("p.b", "1", &_true, &_1, nil, nil),
		_P("p.c", "on", &_true, nil, nil, nil),
		_P("p.d", "yes", &_true, nil, nil, nil),
		_P("p.e", "true", &_true, nil, nil, nil),
		_P("p.f", "0", &_false, &_0, nil, nil),
		_P("p.g", "off", &_false, nil, nil, nil),
		_P("p.h", "no", &_false, nil, nil, nil),
		_P("p.i", "false", &_false, nil, nil, nil),
	}

	_BOOLS = []*ptest{
		// _B( name, property, bool, int, list, error )
		_B("b.a", "?", nil, nil, nil, gitconfig.InvalidBooleanError),
		_B("b.b", "1", &_true, &_1, nil, nil),
		_B("b.c", "on", &_true, nil, nil, nil),
		_B("b.d", "yes", &_true, nil, nil, nil),
		_B("b.e", "true", &_true, nil, nil, nil),
		_B("b.f", "0", &_false, &_0, nil, nil),
		_B("b.g", "off", &_false, nil, nil, nil),
		_B("b.h", "no", &_false, nil, nil, nil),
		_B("b.i", "false", &_false, nil, nil, nil),
	}

	_INTS = []*ptest{
		// _I( name, property, bool, int, list, error )
		_I("i.a", "?", nil, nil, nil, gitconfig.InvalidIntegerError),
		_I("i.b", "1", &_true, &_1, nil, nil),
		_I("i.c", "on", &_true, nil, nil, gitconfig.InvalidIntegerError),
		_I("i.d", "0", &_false, &_0, nil, nil),
		_I("i.e", "off", &_false, nil, nil, gitconfig.InvalidIntegerError),
		_I("i.f", "-123", nil, &__123, nil, nil),
		_I("i.g", "1234", nil, &_1234, nil, nil),
	}

	_LISTS = []*ptest{
		// _L( name, property, bool, int, list, error )
		_L("l.a", "?", nil, nil, nil, nil),
		_L("l.b", "1", &_true, &_1, nil, nil),
		_L("l.c", "on", &_true, nil, nil, nil),
		_L("l.d", "yes", &_true, nil, nil, nil),
		_L("l.e", "true", &_true, nil, nil, nil),
		_L("l.f", "0", &_false, &_0, nil, nil),
		_L("l.g", "off", &_false, nil, nil, nil),
		_L("l.h", "no", &_false, nil, nil, nil),
		_L("l.i", "false", &_false, nil, nil, nil),
		_L("l.j", "a:b:c", nil, nil, []string{"a", "b", "c"}, nil),
		_L("l.k", ":b:", nil, nil, []string{"", "b", ""}, nil),
	}
)

func TestProperty(t *testing.T) {
	// iterate over the property tests
	for _, _test := range _PROPERTIES {
		_test.Test(t)
	}
} // TestProperty()

func TestBool(t *testing.T) {
	// iterate over the boolean tests
	for _, _test := range _BOOLS {
		_test.Test(t)
	}
} // TestBool()

func TestInt(t *testing.T) {
	// iterate over the integer tests
	for _, _test := range _INTS {
		_test.Test(t)
	}
} // TestInt()

func TestList(t *testing.T) {
	// iterate over the list tests
	for _, _test := range _LISTS {
		_test.Test(t)
	}
} // TestList()

//
// helper functions
//

// _P returns the NewProperty test case
func _P(n, v string, b *bool, i *int, l []string, e error) *ptest {
	_property := gitconfig.NewProperty(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &ptest{n, _property, b, i, v, _list, e, nil}
} // _P()

// _B returns the NewBool test case
func _B(n, v string, b *bool, i *int, l []string, e error) *ptest {
	_property, _err := gitconfig.NewBool(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &ptest{n, _property, b, i, v, _list, e, _err}
} // _B()

// _I returns the NewInt test case
func _I(n, v string, b *bool, i *int, l []string, e error) *ptest {
	_property, _err := gitconfig.NewInt(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &ptest{n, _property, b, i, v, _list, e, _err}
} // _I()

// _L returns the NewList test case
func _L(n, v string, b *bool, i *int, l []string, e error) *ptest {
	_property := gitconfig.NewList(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &ptest{n, _property, b, i, v, _list, e, nil}
} // _L()
