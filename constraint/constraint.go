// Copyrights 2015 mparaiso <mparaiso@online.fr>
// License MIT
// version 0.1

package constraint

type Constraint interface {
	Message() string
	Groups() []*Group
	Validate(interface{}) bool
}

type BaseConstraint struct {
	message string
	groups  []*Group
}

func (c BaseConstraint) Message() string {
	return "shouldn't be blank"
}

func (c BaseConstraint) Groups() []*Group {
	return c.groups
}

func (c BaseConstraint) Validate(value interface{}) bool {
	return true
}

func NewFieldConstraint(fieldName string, constraint Constraint) Constraint {
	return &FieldConstraint{fieldName: fieldName, constraint: constraint}
}

type FieldConstraint struct {
	BaseConstraint
	fieldName  string
	constraint Constraint
}

type Group struct {
	Name string
}

type notBlank struct {
	BaseConstraint
}

func NotBlank() Constraint {
	c := new(notBlank)
	c.groups = []*Group{}
	return c
}
