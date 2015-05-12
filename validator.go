// Copyrights 2015 mparaiso <mparaiso@online.fr>
// License MIT
// version 0.1

package validator

import (
	"github.com/interactiv/validator/constraint"
)

type ValidatorMetadataLoader interface {
	LoadValidatorMetadata(metadata *Metadata)
}
type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(loader ValidatorMetadataLoader, groups ...string) (errors []error) {
	metadata := &Metadata{constraints: []constraint.Constraint{}}
	loader.LoadValidatorMetadata(metadata)
	for _, Constraint := range metadata.constraints {
		if err := Constraint.Validate(loader); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

type Metadata struct {
	constraints []constraint.Constraint
}

// AddFieldConstraint adds a  constraint.fieldConstraint to validator.Metadata
func (m *Metadata) AddFieldConstraint(field string, Constraint constraint.Constraint) *Metadata {
	m.constraints = append(m.constraints, constraint.NewFieldConstraint(field, Constraint))
	return m
}

// AddGetterConstraint adds a  constraint.getterConstraint to validator.Metadata
func (m *Metadata) AddGetterConstraint(getter string, Constraint constraint.Constraint) *Metadata {
	m.constraints = append(m.constraints, constraint.NewGetterContraint(getter, Constraint))
	return m
}
