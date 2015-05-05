// Copyrights 2015 mparaiso <mparaiso@online.fr>
// License MIT
// version 0.1

package constraint

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// Constraint represents a constraint that can be validated
type Constraint interface {
	Validate(interface{}) error
}

// NewFieldConstraint returns a constraint for a field of an struct
func NewFieldConstraint(fieldName string, constraint Constraint) Constraint {
	return &FieldConstraint{fieldName: fieldName, constraint: constraint}
}

// FieldConstraint Represents a field constraint
type FieldConstraint struct {
	fieldName  string
	constraint Constraint
}

// FieldError is a field error implementing the Error interface
type FieldError struct {
	error
	fieldName  string
	typeString string
}

// Error returns an error message
func (fe FieldError) Error() string {
	return fe.error.Error()
}

// FieldName returns the name the field in the struct validated
func (fe FieldError) FieldName() string {
	return fe.fieldName
}

// Type returns the type of the struct
func (fe FieldError) Type() string {
	return fe.typeString
}

// Validate validates a field constraint
func (fc *FieldConstraint) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if reflect.Struct != v.Kind() {
		log.Panicf("v% is not a struct", value)
	}
	err := fc.constraint.Validate(v.FieldByName(fc.fieldName).Interface())
	if err != nil {
		return FieldError{error: err, fieldName: fc.fieldName, typeString: v.Type().String()}
	}
	return err
}

// NotBlank returns a notBlank constraint
func NotBlank() Constraint {
	c := new(notBlank)
	return c
}

type notBlank struct {
}

func (nb *notBlank) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if len(v.String()) <= 0 {
			return errors.New(NotBlankMessage)
		}
	default:
		return errors.New(CannotValidateNonStringMessage)
	}
	return nil
}

// Blank returns a blank constraint
func Blank() Constraint {
	c := new(blank)
	return c
}

type blank struct {
}

func (nb *blank) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if len(v.String()) > 0 {
			return errors.New(BlankMessage)
		}
	default:
		return errors.New(CannotValidateNonStringMessage)
	}
	return nil
}

// NotNil returns a notNil constraint
func NotNil() Constraint {
	c := new(notNil)
	return c
}

type notNil struct {
}

func (c *notNil) Validate(value interface{}) error {
	if value == nil {
		return errors.New(NotNillMessage)
	}
	return nil
}

// Nil returns a nil constraint
func Nil() Constraint {
	c := new(nill)
	return c
}

type nill struct {
}

func (c *nill) Validate(value interface{}) error {
	if value != nil {
		return errors.New(NillMessage)
	}
	return nil
}

// True returns a true constraint
func True() Constraint {
	c := new(isTrue)
	return c
}

type isTrue struct {
}

func (c *isTrue) Validate(value interface{}) error {
	if value != true {
		return errors.New(TrueMessage)
	}
	return nil
}

// False returns a false constraint
func False() Constraint {
	c := new(isFalse)
	return c
}

type isFalse struct {
}

func (c *isFalse) Validate(value interface{}) error {
	if value != false {
		return errors.New(FalseMessage)
	}
	return nil
}

// Type returns a type constraint
func Type(theType reflect.Type) Constraint {
	c := new(isType)
	c.theType = theType
	return c
}

type isType struct {
	theType reflect.Type
}

// Validate returns an error if the constraint is violated
func (c *isType) Validate(value interface{}) error {
	if !c.theType.AssignableTo(reflect.TypeOf(value)) {
		return fmt.Errorf(TypeMessage, c.theType.String())
	}
	return nil
}

// validation error messages
const (
	NotBlankMessage                = "This value should not be blank"
	NotNillMessage                 = "This value should not be nil"
	NillMessage                    = "This value should be nil"
	BlankMessage                   = "This value should be blank"
	CannotValidateNonStringMessage = "Cannot validate this value (not a string)"
	TrueMessage                    = "This value should be true"
	FalseMessage                   = "This value should be false"
	TypeMessage                    = "This value should be of type %s"
)
