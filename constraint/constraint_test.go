// Copyrights 2015 mparaiso <mparaiso@online.fr>
// License MIT

package constraint_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/interactiv/expect"
	"github.com/interactiv/validator/constraint"
)

type list []interface{}

func TestConstraints(t *testing.T) {
	e := expect.New(t)
	for _, fixture := range fixtures {
		error := fixture[0].(constraint.Constraint).Validate(fixture[1])
		t.Log(error)
		e.Expect(error == nil).ToBe(fixture[2].(bool))
	}
}

type Example struct{}

var fixtures = []list{
	list{constraint.Blank(), "", true},
	list{constraint.Blank(), "example", false},
	list{constraint.NotBlank(), "example", true},
	list{constraint.NotBlank(), "", false},
	list{constraint.NotNil(), nil, false},
	list{constraint.NotNil(), new(Example), true},
	list{constraint.Nil(), nil, true},
	list{constraint.Nil(), new(Example), false},
	list{constraint.Type(reflect.TypeOf(Example{})), Example{}, true},
	list{constraint.Type(reflect.TypeOf(5)), 10, true},
	list{constraint.Type(reflect.TypeOf("example")), 10, false},
	list{constraint.Email(), "johndoe@example.com", true},
	list{constraint.Email(), "johndoe@example", false},
	list{constraint.Length(5, 10), "example", true},
	list{constraint.Length(7, 7), "example", true},
	list{constraint.Length(7, 7), "examples", false},
	list{constraint.Length(4, 6), "example", false},
	list{constraint.URL().SetProtocols([]string{"https"}), "https://example.com", true},
	list{constraint.URL().SetProtocols([]string{"https"}), "http://example.com", false},
	list{constraint.URL().SetProtocols([]string{"http"}), "example.com", false},
	list{constraint.Regexp(regexp.MustCompile("[a-z A-Z]+\\s[a-z A-Z]+")), "John Doe", true},
	list{constraint.Regexp(regexp.MustCompile("[a-z A-Z]+\\s[a-z A-Z]+")).SetMatch(false), "John Doe", false},
	list{constraint.Regexp(regexp.MustCompile("[a-z A-Z]+\\s[a-z A-Z]+")), "Jane", false},
	list{constraint.Range(10, 15), 10, true},
	list{constraint.Range(10, 15), 10.5, true},
	list{constraint.Range(10, 15), float32(10.6), true},
	list{constraint.Range(10, 15), 9, false},
	list{constraint.Range(10, 15), 16, false},
	list{constraint.EqualTo("foo"), "foo", true},
	list{constraint.EqualTo(10), 10, true},
	list{constraint.EqualTo(10), "10", false},
	list{constraint.NotEqualTo(10), 11, true},
	list{constraint.NotEqualTo("example"), "examples", true},
	list{constraint.LessThan(10), 5, true},
	list{constraint.LessThanOrEqual(10), 10, true},
	list{constraint.GreaterThan(5), 10, true},
	list{constraint.GreaterThanOrEqual(5), 5, true},
}
