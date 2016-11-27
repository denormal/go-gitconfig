package gitconfig_test

import (
	"testing"

	"github.com/denormal/go-gitconfig"
)

type ptest struct {
	n string
	p gitconfig.Property
	b *bool
	i *int
	s string
	l []string
}

func (p *ptest) Test(t *testing.T) {
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
		_bool, _err := p.p.Bool()
		if _err != nil {
			t.Errorf(
				"%q: boolean failure; expected %v, got %q",
				p.n, nil, _err.Error(),
			)
		} else if *p.b != _bool {
			t.Errorf(
				"%q: boolean mismatch; expected %v, got %v",
				*p.b, _bool,
			)
		}
	} else {
		_bool, _err := p.p.Bool()
		if _err == nil {
			t.Errorf(
				"%q: unexpected boolean success; expected error, got %q",
				p.n, _err,
			)
		} else if _err != gitconfig.InvalidBooleanError {
			t.Errorf(
				"%q: unexpected boolean error; expected %q, got %q",
				p.n, gitconfig.InvalidBooleanError.Error(), _err.Error(),
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
		_int, _err := p.p.Int()
		if _err != nil {
			t.Errorf(
				"%q: integer failure; expected %v, got %q",
				p.n, true, _err.Error(),
			)
		} else if *p.i != _int {
			t.Errorf(
				"%q: integer mismatch; expected %v, got %v",
				*p.i, _int,
			)
		}
	} else {
		_int, _err := p.p.Int()
		if _err == nil {
			t.Errorf(
				"%q: unexpected integer success; expected error, got %v",
				p.n, _err,
			)
		} else if _err != gitconfig.InvalidIntegerError {
			t.Errorf(
				"%q: unexpected integer error; expected %q, got %q",
				p.n, gitconfig.InvalidIntegerError.Error(), _err.Error(),
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
		// _P( name, property, bool, int, list )
		_P("p.a", "-123", nil, &__123, nil),
		_P("p.b", "1234", nil, &_1234, nil),
		_P("p.c", ":b:", nil, nil, []string{"", "b", ""}),
		_P("p.d", "?", nil, nil, nil),
		_P("p.f", "a:b:c", nil, nil, []string{"a", "b", "c"}),
		_P("p.g", "0", &_false, &_0, nil),
		_P("p.h", "no", &_false, nil, nil),
		_P("p.i", "off", &_false, nil, nil),
		_P("p.j", "false", &_false, nil, nil),
		_P("p.k", "1", &_true, &_1, nil),
		_P("p.l", "on", &_true, nil, nil),
		_P("p.m", "yes", &_true, nil, nil),
		_P("p.n", "true", &_true, nil, nil),
	}
)

func TestProperty(t *testing.T) {
	// iterate over the property tests
	for _, _test := range _PROPERTIES {
		_test.Test(t)
	}
} // TestProperty()

//
// helper functions
//

// _P returns the NewProperty test case
func _P(n, v string, b *bool, i *int, l []string) *ptest {
	_property := gitconfig.NewProperty(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &ptest{n, _property, b, i, v, _list}
} // _P()
