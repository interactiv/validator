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

func (v *Validator) Validate(loader ValidatorMetadataLoader) (errors []string) {
	metadata := &Metadata{constraints: []constraint.Constraint{}}
	loader.LoadValidatorMetadata(metadata)
	for _, Constraint := range metadata.constraints {
		if !Constraint.Validate(loader) {
			errors = append(errors, Constraint.Message())
		}
	}
	return errors
}

type Metadata struct {
	constraints []constraint.Constraint
}

func (m *Metadata) AddFieldConstraint(field string, Constraint constraint.Constraint) *Metadata {
	m.constraints = append(m.constraints, constraint.NewFieldConstraint(field, Constraint))
	return m
}

type ValidationError struct {
}

func (ve *ValidationError) Error() string {
	return "error"
}
