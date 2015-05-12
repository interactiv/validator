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
	validPerson := &Person{Name: "John Doe", IsMarried: true, age: 25}
	Validator := validator.New()
	Errors := Validator.Validate(validPerson)
	e.Expect(len(Errors)).ToBe(0)
	invalidPerson := &Person{Name: "John Doe", IsMarried: false, age: 12}
	Errors = Validator.Validate(invalidPerson)
	e.Expect(Errors[0].Error()).ToBe("This value should be true")
	e.Expect(Errors[1].Error()).ToBe("This value should be greater than 15")

}

func TestGroups(t *testing.T) {
	e := expect.New(t)
	person := &Person{Name: "Mike", age: 20}
	errors := validator.New().Validate(person, "group1")
	e.Expect(len(errors)).ToBe(0)
	errors = validator.New().Validate(person, "group2")
	e.Expect(len(errors)).ToEqual(1)
}

/********************************/
/*         FIXTURES             */
/********************************/

type Person struct {
	Name      string
	IsMarried bool
	age       int
}

func (p Person) Age() int {
	return p.age
}

func (p *Person) SetAge(age int) *Person {
	p.age = age
	return p
}

func (p *Person) LoadValidatorMetadata(metadata *validator.Metadata) {
	metadata.AddFieldConstraint("Name", constraint.NotBlank().SetGroups([]string{"group1"})).
		AddFieldConstraint("IsMarried", constraint.True().SetGroups([]string{"group2"})).
		AddGetterConstraint("Age", constraint.GreaterThan(15).SetGroups([]string{"group1", "group2"}))
}
