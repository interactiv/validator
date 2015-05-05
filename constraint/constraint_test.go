// Copyrights 2015 mparaiso <mparaiso@online.fr>
// License MIT

package constraint_test

import (
	"reflect"
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
}
