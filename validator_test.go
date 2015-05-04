// Copyrights 2015 mparaiso <mparaiso@online.fr>
// License MIT
// version 0.1
package validator_test

import (
	"testing"

	"github.com/interactiv/expect"
	"github.com/interactiv/validator"
	"github.com/interactiv/validator/constraint"
)

func TestFactory(t *testing.T) {
	e := expect.New(t)
	validPerson := &Person{Name: "John Doe"}
	Validator := validator.New()
	Errors := Validator.Validate(validPerson)
	e.Expect(len(Errors)).ToBe(0)
	invalidPerson := &Person{Name: "John Doe"}
	Errors = Validator.Validate(invalidPerson)
	e.Expect(len(Errors)).ToBeGreaterThan(0)
}

/********************************/
/*         FIXTURES             */
/********************************/

type Person struct {
	Name string
}

func (p *Person) LoadValidatorMetadata(metadata *validator.Metadata) {
	metadata.AddFieldConstraint("Name", constraint.NotBlank())
}
