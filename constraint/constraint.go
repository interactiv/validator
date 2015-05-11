// Copyrights 2015 mparaiso <mparaiso@online.fr>
// License MIT
// version 0.1

package constraint

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
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

// Email returns an email constraint
func Email() Constraint {
	c := new(email)
	return c
}

type email struct {
}

// Validate returns an error if the constraint is violated
func (c email) Validate(value interface{}) error {
	var ok bool
	var val string
	if val, ok = value.(string); ok != true {
		return errors.New(CannotValidateNonStringMessage)
	}
	if !EmailRegexp.MatchString(val) {
		return errors.New(EmailMessage)
	}
	return nil
}

// Length returns an length constraint
func Length(min int, max int) Constraint {
	c := &length{min, max}
	return c
}

type length struct {
	min int
	max int
}

// Validate returns an error if the constraint is violated
func (c length) Validate(value interface{}) error {
	var ok bool
	var val string
	if val, ok = value.(string); ok != true {
		return errors.New(CannotValidateNonStringMessage)
	}
	if c.min == c.max {
		if c.min != len(val) {
			return fmt.Errorf(ExactLengthMessage, c.min)
		}
	} else {
		if !(c.min <= len(val)) {
			return fmt.Errorf(MinMessage, c.min)
		}
		if !(len(val) <= c.max) {
			return fmt.Errorf(MaxMessage, c.max)
		}
	}
	return nil
}

// URL returns an url constraint
func URL() *URLConstraint {
	c := new(URLConstraint)
	c.protocols = []string{}
	return c
}

// URLConstraint represents an url constraint
type URLConstraint struct {
	protocols []string
}

// Protocols return the protocols supported by the constraint
func (c URLConstraint) Protocols() []string {
	return c.protocols
}

// SetProtocols sets the protocols supported by the constraint
func (c *URLConstraint) SetProtocols(protocols []string) *URLConstraint {
	c.protocols = protocols
	return c
}

// Validate returns an error if the constraint is violated
func (c URLConstraint) Validate(value interface{}) error {
	var ok bool
	var val string
	if val, ok = value.(string); ok != true {
		return errors.New(CannotValidateNonStringMessage)
	}
	if parsedURL, err := url.Parse(val); err != nil {
		return errors.New(URLMessage)
	} else if c.protocols != nil && len(c.protocols) > 0 {
		for _, protocol := range c.protocols {
			if parsedURL.Scheme == protocol {
				return nil
			}
			return errors.New(URLMessage)
		}
	}
	return nil
}

func Regexp(pattern *regexp.Regexp) *RegexpConstraint {
	return &RegexpConstraint{pattern, true}
}

type RegexpConstraint struct {
	pattern *regexp.Regexp
	match   bool
}

func (c RegexpConstraint) Match() bool {
	return c.match
}

func (c *RegexpConstraint) SetMatch(match bool) *RegexpConstraint {
	c.match = match
	return c
}

// Validate returns an error if the constraint is violated
func (c *RegexpConstraint) Validate(value interface{}) error {
	var ok bool
	var val string
	if val, ok = value.(string); ok != true {
		return errors.New(CannotValidateNonStringMessage)
	}
	if c.match && !c.pattern.MatchString(val) {
		return errors.New(RegexpMatchMessage)
	}
	if !c.match && c.pattern.MatchString(val) {
		return errors.New(RegexpMatchMessage)
	}
	return nil
}

func Range(min, max float64) Constraint {
	return &rangeConstraint{min, max}
}

type rangeConstraint struct {
	min float64
	max float64
}

// Validate returns an error if the constraint is violated
func (rc *rangeConstraint) Validate(value interface{}) error {
	valFloat64, err := ToFloat64(value)
	if err != nil {
		return errors.New(ErrorNotNumberMessage)
	}

	if !(rc.min <= valFloat64) {
		return fmt.Errorf(RangeMinMessage, fmt.Sprint(rc.min))
	}
	if !(valFloat64 <= rc.max) {
		return fmt.Errorf(RangeMaxMessage, fmt.Sprint(rc.max))

	}
	return nil
}

func EqualTo(val interface{}) Constraint {
	return &equalTo{val}
}

type equalTo struct {
	value interface{}
}

// Validate returns an error if the constraint is violated
func (c *equalTo) Validate(value interface{}) error {
	if c.value != value {
		return fmt.Errorf(EqualToMessage, c.value)
	}
	return nil
}

func NotEqualTo(val interface{}) Constraint {
	return &notEqualTo{val}
}

type notEqualTo struct {
	value interface{}
}

// Validate returns an error if the constraint is violated
func (c *notEqualTo) Validate(value interface{}) error {
	if c.value == value {
		return fmt.Errorf(NotEqualToMessage, c.value)
	}
	return nil
}

func LessThan(value float64) Constraint {
	return &lessThan{value}
}

type lessThan struct {
	value float64
}

// Validate returns an error if the constraint is violated
func (c lessThan) Validate(value interface{}) error {
	if val, err := ToFloat64(value); err != nil {
		errors.New(ErrorNotNumberMessage)
	} else if val >= c.value {
		return fmt.Errorf(LessThanMessage, c.value)
	}
	return nil
}

func LessThanOrEqual(value float64) Constraint {
	return &lessThanOrEqual{value}
}

type lessThanOrEqual struct {
	value float64
}

// Validate returns an error if the constraint is violated
func (c *lessThanOrEqual) Validate(value interface{}) error {
	if val, err := ToFloat64(value); err != nil {
		errors.New(ErrorNotNumberMessage)
	} else if val > c.value {
		return fmt.Errorf(LessThanOrEqualMessage, c.value)
	}
	return nil
}

func GreaterThan(value float64) Constraint {
	return &greaterThan{value}
}

type greaterThan struct {
	value float64
}

// Validate returns an error if the constraint is violated
func (c *greaterThan) Validate(value interface{}) error {
	if val, err := ToFloat64(value); err != nil {
		errors.New(ErrorNotNumberMessage)
	} else if val <= c.value {
		return fmt.Errorf(GreaterThanMessage, c.value)
	}
	return nil
}

// GreaterThanOrEqual returns a great than equal constraint
func GreaterThanOrEqual(value float64) Constraint {
	return &greaterThanOrEqual{value, GreaterThanOrEqualMessage}
}

type greaterThanOrEqual struct {
	value   float64
	message string
}

// Message returns the error message
func (c greaterThanOrEqual) Message() string {
	return c.message
}

// Message sets the error message
func (c *greaterThanOrEqual) SetMessage(message string) *greaterThanOrEqual {
	c.message = message
	return c
}

// Validate returns an error if the constraint is violated
func (c greaterThanOrEqual) Validate(value interface{}) error {
	if val, err := ToFloat64(value); err != nil {
		errors.New(ErrorNotNumberMessage)
	} else if val < c.value {
		return fmt.Errorf(c.message, c.value)
	}
	return nil
}

func Choice(choices []interface{}) *choice {
	return &choice{choices: choices}
}

type choice struct {
	choices         []interface{}
	multiple        bool
	min             int
	max             int
	message         string
	multipleMessage string
	minMessage      string
	maxMessage      string
	strict          bool
}

// GetChoices returns a []interface{}
func (choice choice) GetChoices() []interface{} {
	return choice.choices
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetChoices(choices []interface{}) *choice {
	choice.choices = choices
	return choice
}

// GetMultiple returns a bool
func (choice choice) GetMultiple() bool {
	return choice.multiple
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetMultiple(multiple bool) *choice {
	choice.multiple = multiple
	return choice
}

// GetMin returns a int
func (choice choice) GetMin() int {
	return choice.min
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetMin(min int) *choice {
	choice.min = min
	return choice
}

// GetMax returns a int
func (choice choice) GetMax() int {
	return choice.max
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetMax(max int) *choice {
	choice.max = max
	return choice
}

// GetMessage returns a string
func (choice choice) GetMessage() string {
	return choice.message
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetMessage(message string) *choice {
	choice.message = message
	return choice
}

// GetMultipleMessage returns a string
func (choice choice) GetMultipleMessage() string {
	return choice.multipleMessage
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetMultipleMessage(multipleMessage string) *choice {
	choice.multipleMessage = multipleMessage
	return choice
}

// GetMinMessage returns a string
func (choice choice) GetMinMessage() string {
	return choice.minMessage
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetMinMessage(minMessage string) *choice {
	choice.minMessage = minMessage
	return choice
}

// GetMaxMessage returns a string
func (choice choice) GetMaxMessage() string {
	return choice.maxMessage
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetMaxMessage(maxMessage string) *choice {
	choice.maxMessage = maxMessage
	return choice
}

// GetStrict returns a bool
func (choice choice) GetStrict() bool {
	return choice.strict
}

// Setchoice sets *choice.choice and returns *choice
func (choice *choice) SetStrict(strict bool) *choice {
	choice.strict = strict
	return choice
}

func (c choice) Validate(values interface{}) error {
	switch t := values.(type) {
	case []interface{}:
		for _, choice := range c.choices {
			for _, c := range t {
				if choice == c {
					return nil
				}
			}
		}
	default:
		return errors.New(ErrorNotArrayMessage)
	}
	return fmt.Errorf("this value should be %s", c.GetChoices())
}

/***********/
/* HELPERS */
/***********/

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("Cant convert %v to float64")
	}
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
	EmailMessage                   = " This value is not a valid email address"
	MinMessage                     = "This value is too short. It should have %d characters or more."
	MaxMessage                     = "This value is too long. It should have %d characters or less"
	ExactLengthMessage             = "This value should have exactly %d characters"
	URLMessage                     = "This value is not a valid URL."
	RegexpMatchMessage             = "This value is not valid"
	RangeMinMessage                = "This value should be %s or more"
	RangeMaxMessage                = "This value should be %s or less"
	ErrorNotNumberMessage          = "This value should be a valid number"
	ErrorNotArrayMessage           = "This value should be a valid array or slice"
	EqualToMessage                 = "This value should be equal to %v"
	NotEqualToMessage              = "This value should not be equal to %v"
	LessThanMessage                = "This value should be less than %d"
	LessThanOrEqualMessage         = "This value should be less than or equal to %d"
	GreaterThanMessage             = "This value should be greater than %d"
	GreaterThanOrEqualMessage      = "This value should be greater than or equal to %d"
)

var (
	// EmailRegexp represents an email pattern
	EmailRegexp = regexp.MustCompile(".+\\@.+\\..+")
)
