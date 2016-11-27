package gitconfig_test

import (
	"testing"

	"github.com/denormal/go-gitconfig"
)

type vtest struct {
	n   string
	v   gitconfig.Value
	b   *bool
	i   *int
	s   string
	l   []string
	e   error
	err gitconfig.Error
}

func (vt *vtest) Test(t *testing.T) {
	// do we expect an error?
	if vt.e != nil {
		if vt.err == nil {
			t.Errorf(
				"%q: expected New error %q; none found",
				vt.n, vt.e.Error(),
			)
		} else if vt.err.Underlying() != vt.e {
			t.Errorf(
				"%q: New error mismatch; expected %q, got %q",
				vt.n, vt.e.Error(), vt.err.Error(),
			)
		}
	} else {

		// ensure the name is as expected
		if vt.v.Name() != vt.n {
			t.Errorf(
				"%q: unexpected name; expected %q, got %q",
				vt.n, vt.n, vt.v.Name(),
			)
		}

		// ensure the string form is as expected
		if vt.v.String() != vt.s {
			t.Errorf(
				"%q: unexpected string; expected %q, got %q",
				vt.n, vt.s, vt.v.String(),
			)
		}

		// do we expect a boolean value?
		if vt.b != nil {
			_bool, _ok := vt.v.Bool()
			if !_ok {
				t.Errorf(
					"%q: boolean failure; expected %v, got %v",
					vt.n, true, _ok,
				)
			} else if *vt.b != _bool {
				t.Errorf(
					"%q: boolean mismatch; expected %v, got %v",
					*vt.b, _bool,
				)
			}
		} else {
			_bool, _ok := vt.v.Bool()
			if _ok {
				t.Errorf(
					"%q: unexpected boolean success; expected %v, got %v",
					vt.n, false, _ok,
				)
			} else if _bool != false {
				t.Errorf(
					"%q: unexpected boolean return; expected %v, got %v",
					vt.n, false, _bool,
				)
			}
		}

		// do we expect an integer value?
		if vt.i != nil {
			_int, _ok := vt.v.Int()
			if !_ok {
				t.Errorf(
					"%q: integer failure; expected %v, got %v",
					vt.n, true, _ok,
				)
			} else if *vt.i != _int {
				t.Errorf(
					"%q: integer mismatch; expected %v, got %v",
					*vt.i, _int,
				)
			}
		} else {
			_int, _ok := vt.v.Int()
			if _ok {
				t.Errorf(
					"%q: unexpected integer success; expected %v, got %v",
					vt.n, false, _ok,
				)
			} else if _int != 0 {
				t.Errorf(
					"%q: unexpected integer return; expected %v, got %v",
					vt.n, 0, _int,
				)
			}
		}

		// do we expect a list value?
		if vt.i != nil {
			_list, _ok := vt.v.List()
			if !_ok {
				t.Errorf(
					"%q: list failure; expected %v, got %v",
					vt.n, true, _ok,
				)
			} else if len(vt.l) != len(_list) {
				t.Errorf(
					"%q: list length mismatch; expected %v, got %v",
					vt.n, len(vt.l), len(_list),
				)
			} else {
				for _i, _v := range vt.l {
					if _list[_i] != _v {
						t.Errorf(
							"%q: list item mismatch; "+
								"expected %v, got %v for item %d",
							vt.n, _v, _list[_i], _i,
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

	// define the value tests
	//		- each test must have a unique name
	//		- tests are expected to start with X. to assist config_test.go
	_VALUE = []*vtest{
		// _V( name, value, bool, int, list, error )
		_V("v.a", "?", nil, nil, nil, nil),
		_V("v.b", "1", &_true, &_1, nil, nil),
		_V("v.c", "on", &_true, nil, nil, nil),
		_V("v.d", "yes", &_true, nil, nil, nil),
		_V("v.e", "true", &_true, nil, nil, nil),
		_V("v.f", "0", &_false, &_0, nil, nil),
		_V("v.g", "off", &_false, nil, nil, nil),
		_V("v.h", "no", &_false, nil, nil, nil),
		_V("v.i", "false", &_false, nil, nil, nil),
	}

	_BOOL = []*vtest{
		// _B( name, value, bool, int, list, error )
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

	_INT = []*vtest{
		// _I( name, value, bool, int, list, error )
		_I("i.a", "?", nil, nil, nil, gitconfig.InvalidIntegerError),
		_I("i.b", "1", &_true, &_1, nil, nil),
		_I("i.c", "on", &_true, nil, nil, gitconfig.InvalidIntegerError),
		_I("i.d", "0", &_false, &_0, nil, nil),
		_I("i.e", "off", &_false, nil, nil, gitconfig.InvalidIntegerError),
		_I("i.f", "-123", nil, &__123, nil, nil),
		_I("i.g", "1234", nil, &_1234, nil, nil),
	}

	_LIST = []*vtest{
		// _L( name, value, bool, int, list, error )
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

func TestValue(t *testing.T) {
	// iterate over the value tests
	for _, _test := range _VALUE {
		_test.Test(t)
	}
} // TestValue()

func TestBool(t *testing.T) {
	// iterate over the boolean tests
	for _, _test := range _BOOL {
		_test.Test(t)
	}
} // TestBool()

func TestInt(t *testing.T) {
	// iterate over the integer tests
	for _, _test := range _INT {
		_test.Test(t)
	}
} // TestInt()

func TestList(t *testing.T) {
	// iterate over the list tests
	for _, _test := range _LIST {
		_test.Test(t)
	}
} // TestList()

//
// helper functions
//

// _V returns the NewValue test case
func _V(n, v string, b *bool, i *int, l []string, e error) *vtest {
	_value, _err := gitconfig.NewValue(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &vtest{n, _value, b, i, v, _list, e, _err}
} // _V()

// _B returns the NewBool test case
func _B(n, v string, b *bool, i *int, l []string, e error) *vtest {
	_value, _err := gitconfig.NewBool(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &vtest{n, _value, b, i, v, _list, e, _err}
} // _B()

// _I returns the NewInt test case
func _I(n, v string, b *bool, i *int, l []string, e error) *vtest {
	_value, _err := gitconfig.NewInt(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &vtest{n, _value, b, i, v, _list, e, _err}
} // _I()

// _L returns the NewList test case
func _L(n, v string, b *bool, i *int, l []string, e error) *vtest {
	_value, _err := gitconfig.NewList(n, v)
	_list := l
	if _list == nil {
		_list = []string{v}
	}

	return &vtest{n, _value, b, i, v, _list, e, _err}
} // _L()
