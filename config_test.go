package gitconfig_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/denormal/go-gitconfig"
)

var (
	_ALL    []gitconfig.Value
	_ORDER  []string
	_STRING = "" +
		"b.b=1\n" +
		"b.c=on\n" +
		"b.d=yes\n" +
		"b.e=true\n" +
		"b.f=0\n" +
		"b.g=off\n" +
		"b.h=no\n" +
		"b.i=false\n" +
		"i.b=1\n" +
		"i.d=0\n" +
		"i.f=-123\n" +
		"i.g=1234\n" +
		"l.a=?\n" +
		"l.b=1\n" +
		"l.c=on\n" +
		"l.d=yes\n" +
		"l.e=true\n" +
		"l.f=0\n" +
		"l.g=off\n" +
		"l.h=no\n" +
		"l.i=false\n" +
		"l.j=a:b:c\n" +
		"l.k=:b:\n" +
		"v.a=?\n" +
		"v.b=1\n" +
		"v.c=on\n" +
		"v.d=yes\n" +
		"v.e=true\n" +
		"v.f=0\n" +
		"v.g=off\n" +
		"v.h=no\n" +
		"v.i=false\n"
)

func TestNewConfig(t *testing.T) {
	// ensure NewConfig() behaves as expected
	_config := gitconfig.NewConfig(_ALL)
	if _config == nil {
		t.Fatalf("unexpected nil Config; expected instance")
	}
} // TestNewConfig()

func TestConfigAll(t *testing.T) {
	_config := gitconfig.NewConfig(_ALL)
	_all := _config.All()

	// ensure All() returns the correct number of values
	if len(_all) != len(_ALL) {
		t.Fatalf(
			"unexpected size from All(); expected %d, got %d",
			len(_ALL), len(_all),
		)
	}
	for _i, _got := range _all {
		// compare the values by name
		//		- we compare other attributes later
		_expected := _ORDER[_i]
		if _got.Name() != _expected {
			t.Fatalf(
				"unexpected All() order; expected %q, got %q at index %d",
				_expected, _got.Name(), _i,
			)
		}
	}
} // TestConfigAll()

func TestConfigGet(t *testing.T) {
	_config := gitconfig.NewConfig(_ALL)

	// ensure Get() retrieves the data we expect
	for _, _expected := range _ALL {
		_got, _ok := _config.Get(_expected.Name())
		if !_ok {
			t.Fatalf(
				"%q: unexpected Get() failure; expected %v, got %v",
				_expected.Name(), true, _ok,
			)
		} else if _got == nil {
			t.Fatalf(
				"%q: unexpected nil value; expected %v",
				_expected.Name(), _expected,
			)
		} else if _got.Name() != _expected.Name() {
			t.Fatalf(
				"%q: unexpected name; expected %q, got %q",
				_expected.Name(), _expected.Name(), _got.Name(),
			)
		} else if _got != _expected {
			t.Fatalf(
				"%q: unexpected Get() return; expected %v, got %v",
				_expected.Name(), _expected, _got,
			)
		}
	}

	// ensure Get() behaves as expected when given unknown names
	for _, _name := range _ORDER {
		_get := _name + time.Now().String()
		_got, _ok := _config.Get(_get)
		if _ok {
			t.Fatalf(
				"%q: unexpected Get() success; expected %v, got %v",
				_get, false, _ok,
			)
		} else if _got != nil {
			t.Fatalf(
				"%q: unexpected Get() value; expected nil, got %v",
				_get, _got.Name(),
			)
		}
	}
} // TestConfigGet()

func TestConfigFind(t *testing.T) {
	// ensure Find() behaves
	find(_VALUE, "v.", t)
	find(_BOOL, "b.", t)
	find(_INT, "i.", t)
	find(_LIST, "l.", t)
} // TestConfigFind()

func TestConfigString(t *testing.T) {
	_config := gitconfig.NewConfig(_ALL)

	// ensure String() does what is expected
	_string := _config.String()
	if _string == "" {
		t.Fatal("unexpected empty string")
	} else if _string != _STRING {
		t.Fatalf(
			"expected string; expected %q, got %q",
			_STRING, _string,
		)
	} else if _string != _config.String() {
		t.Fatal("unexpected string; expected same value, got different")
	}
} // TestConfigString()

func init() {
	// combine all value tests into a single list
	_all := make([]*vtest, 0)
	_all = append(_all, _BOOL...)
	_all = append(_all, _INT...)
	_all = append(_all, _LIST...)
	_all = append(_all, _VALUE...)

	// build the list of all values
	_ALL = make([]gitconfig.Value, 0)
	for _, _v := range _all {
		if _v.v != nil {
			_ALL = append(_ALL, _v.v)
		}
	}

	// generate the expected order
	//		- this is the sort order of the value names
	_ORDER = make([]string, 0, len(_ALL))
	for _, _v := range _ALL {
		_ORDER = append(_ORDER, _v.Name())
	}
	sort.Strings(_ORDER)

	// generate a random permutation of all test values
	rand.Seed(time.Now().UnixNano())
	for _i := range _ALL {
		_j := rand.Intn(_i + 1)
		_ALL[_i], _ALL[_j] = _ALL[_j], _ALL[_i]
	}
} // init()

//
// helper function
//

func find(tests []*vtest, prefix string, t *testing.T) {
	// build a map of the results expected for these tests
	//		- extract the
	_map := make(map[string]gitconfig.Value)
	for _, _test := range tests {
		if _test.v != nil {
			_map[_test.n] = _test.v
		}
	}

	// ensure Find() returns all the expected tests
	//		- first we add "*" to the prefix to create the wildcard pattern
	_config := gitconfig.NewConfig(_ALL)
	_pattern := prefix + "*"
	_find := _config.Find(_pattern)
	if len(_find) != len(_map) {
		t.Fatalf(
			"%q: unexpected Find(); expected %d results, got %d",
			_pattern, len(_map), len(_find),
		)
	}

	// now, try Find() with just the prefix
	//		- this should fail
	_find = _config.Find(prefix)
	if len(_find) != 0 {
		t.Fatalf(
			"%q: unexpected Find(); expected %d results, got %d",
			prefix, 0, len(_find),
		)
	}

	// finally, Find() should behave just as Get() when given an exact match
	for _, _value := range _ALL {
		_find := _config.Find(_value.Name())
		if len(_find) != 1 {
			t.Fatalf(
				"%q: unexpected Find(); expected %d results, got %d",
				_value.Name(), 1, len(_find),
			)
		} else if _find[0].Name() != _value.Name() {
			t.Fatalf(
				"%q: unexpected Find(); expected %q, got %q",
				_value.Name(), _value.Name(), _find[0].Name(),
			)
		}
	}
} // find()
