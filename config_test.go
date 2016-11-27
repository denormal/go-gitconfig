package gitconfig_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/denormal/go-gitconfig"
)

var (
	_ALL    []gitconfig.Property
	_ORDER  []string
	_STRING = "" +
		"p.a=-123\n" +
		"p.b=1234\n" +
		"p.c=:b:\n" +
		"p.d=?\n" +
		"p.f=a:b:c\n" +
		"p.g=0\n" +
		"p.h=no\n" +
		"p.i=off\n" +
		"p.j=false\n" +
		"p.k=1\n" +
		"p.l=on\n" +
		"p.m=yes\n" +
		"p.n=true\n"
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

	// ensure All() returns the correct number of properties
	if len(_all) != len(_ALL) {
		t.Fatalf(
			"unexpected size from All(); expected %d, got %d",
			len(_ALL), len(_all),
		)
	}
	for _i, _got := range _all {
		// compare the properties by name
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
				"%q: unexpected nil property; expected %v",
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
				"%q: unexpected Get() property; expected nil, got %v",
				_get, _got.Name(),
			)
		}
	}
} // TestConfigGet()

func TestConfigFind(t *testing.T) {
	// ensure Find() behaves
	find(_PROPERTIES, "p.", t)
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
		t.Fatal("unexpected string; expected same property, got different")
	}
} // TestConfigString()

func init() {
	// build the list of all properties
	_ALL = make([]gitconfig.Property, 0)
	for _, _p := range _PROPERTIES {
		if _p.p != nil {
			_ALL = append(_ALL, _p.p)
		}
	}

	// generate the expected order
	//		- this is the sort order of the property names
	_ORDER = make([]string, 0, len(_ALL))
	for _, _p := range _ALL {
		_ORDER = append(_ORDER, _p.Name())
	}
	sort.Strings(_ORDER)

	// generate a random permutation of all test properties
	rand.Seed(time.Now().UnixNano())
	for _i := range _ALL {
		_j := rand.Intn(_i + 1)
		_ALL[_i], _ALL[_j] = _ALL[_j], _ALL[_i]
	}
} // init()

//
// helper function
//

func find(tests []*ptest, prefix string, t *testing.T) {
	// build a map of the results expected for these tests
	//		- extract the
	_map := make(map[string]gitconfig.Property)
	for _, _test := range tests {
		if _test.p != nil {
			_map[_test.n] = _test.p
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
	for _, _property := range _ALL {
		_find := _config.Find(_property.Name())
		if len(_find) != 1 {
			t.Fatalf(
				"%q: unexpected Find(); expected %d results, got %d",
				_property.Name(), 1, len(_find),
			)
		} else if _find[0].Name() != _property.Name() {
			t.Fatalf(
				"%q: unexpected Find(); expected %q, got %q",
				_property.Name(), _property.Name(), _find[0].Name(),
			)
		}
	}
} // find()
