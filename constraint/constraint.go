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
	SetGroups(group []string) Constraint
	Groups() []string
}

type BaseConstraint struct {
	groups []string
}

func (c *BaseConstraint) Groups() []string {
	if c.groups == nil {
		c.groups = []string{}
	}
	return c.groups
}
func (c *BaseConstraint) SetGroups(groups []string) Constraint {
	c.groups = groups
	return c
}
func (c *BaseConstraint) Validate(value interface{}) error {
	return nil
}

// NewFieldConstraint returns a constraint for a field of an struct
func NewFieldConstraint(fieldName string, constraint Constraint) Constraint {
	return &fieldConstraint{fieldName: fieldName, constraint: constraint}
}

type fieldConstraint struct {
	BaseConstraint
	fieldName  string
	constraint Constraint
}

// Validate validates a field constraint
func (fc *fieldConstraint) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if reflect.Struct != v.Kind() {
		log.Panicf("%v is not a struct", fmt.Sprint(value))
	}
	err := fc.constraint.Validate(v.FieldByName(fc.fieldName).Interface())

	return err
}

func NewGetterContraint(getterName string, constraint Constraint) Constraint {
	return &getterConstraint{getterName: getterName, constraint: constraint}
}

type getterConstraint struct {
	BaseConstraint
	getterName string
	constraint Constraint
}

// Validate returns an error if the constraint is violated
func (gc getterConstraint) Validate(Struct interface{}) error {
	value := reflect.ValueOf(Struct)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("%v is not a struct", Struct)
	}
	result := value.MethodByName(gc.getterName).Call([]reflect.Value{})[0].Interface()
	log.Print(result)
	return gc.constraint.Validate(result)

}

// NotBlank returns a notBlank constraint
func NotBlank() Constraint {
	c := new(notBlank)
	return c
}

type notBlank struct {
	BaseConstraint
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
	BaseConstraint
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
	BaseConstraint
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
	BaseConstraint
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
	BaseConstraint
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
	BaseConstraint
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
	BaseConstraint
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
	BaseConstraint
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
	c := &length{min: min, max: max}
	return c
}

type length struct {
	BaseConstraint
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
	BaseConstraint
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

	return &RegexpConstraint{pattern: pattern, match: true}
}

type RegexpConstraint struct {
	BaseConstraint
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
	return &rangeConstraint{min: min, max: max}
}

type rangeConstraint struct {
	BaseConstraint
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
	return &equalTo{value: val}
}

type equalTo struct {
	BaseConstraint
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
	return &notEqualTo{value: val}
}

type notEqualTo struct {
	BaseConstraint
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
	return &lessThan{value: value}
}

type lessThan struct {
	BaseConstraint
	value float64
}

// Validate returns an error if the constraint is violated
func (c lessThan) Validate(value interface{}) error {
	if val, err := ToFloat64(value); err != nil {
		errors.New(ErrorNotNumberMessage)
	} else if val >= c.value {
		return fmt.Errorf(LessThanMessage, fmt.Sprint(c.value))
	}
	return nil
}

func LessThanOrEqual(value float64) Constraint {
	return &lessThanOrEqual{value: value}
}

type lessThanOrEqual struct {
	BaseConstraint
	value float64
}

// Validate returns an error if the constraint is violated
func (c *lessThanOrEqual) Validate(value interface{}) error {
	if val, err := ToFloat64(value); err != nil {
		errors.New(ErrorNotNumberMessage)
	} else if val > c.value {
		return fmt.Errorf(LessThanOrEqualMessage, fmt.Sprint(c.value))
	}
	return nil
}

func GreaterThan(value float64) Constraint {
	return &greaterThan{value: value}
}

type greaterThan struct {
	BaseConstraint
	value float64
}

// Validate returns an error if the constraint is violated
func (c *greaterThan) Validate(value interface{}) error {
	if val, err := ToFloat64(value); err != nil {
		return errors.New(ErrorNotNumberMessage)
	} else if val <= c.value {
		return fmt.Errorf(GreaterThanMessage, fmt.Sprint(c.value))
	}
	return nil
}

// GreaterThanOrEqual returns a great than equal constraint
func GreaterThanOrEqual(value float64) Constraint {
	return &greaterThanOrEqual{value: value, message: GreaterThanOrEqualMessage}
}

type greaterThanOrEqual struct {
	BaseConstraint
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
	return &choice{choices: choices, multiple: true, max: -1}
}

type choice struct {
	BaseConstraint
	choices         []interface{}
	multiple        bool
	min             int
	max             int
	message         string
	multipleMessage string
	minMessage      string
	maxMessage      string
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

// Validate returns an error if the constraint is violated
func (c choice) Validate(values interface{}) error {
	switch IsArrayorSlice(values) {
	case true:
		if array, err := ToInterfaceArray(values); err != nil {
			return err
		} else {
			return c.validateArray(array)
		}

	default:
		for _, choice := range c.choices {
			if values == choice {
				return nil
			}
		}
	}
	return errors.New(ChoiceMessage)
}

func (c choice) validateArray(values []interface{}) error {
	if len(values) < c.min {
		return errors.New(ChoiceMinMessage)
	}
	if c.max > 0 && len(values) > c.max {
		return errors.New(ChoiceMaxMessage)
	}
	for _, value := range values {
		index := -1
		for i, choice := range c.choices {
			if choice == value {
				index = i
				break
			}
		}
		if index < 0 {
			return errors.New(ChoiceMultipleMessage)
		}
	}
	return nil
}

func Count(min int, max int) *count {
	return &count{min: min, max: max}
}

type count struct {
	BaseConstraint
	min int
	max int
}

// GetMin returns a int
func (count count) GetMin() int {
	return count.min
}

// SetMin sets *count.count and returns *count
func (count *count) SetMin(min int) *count {
	count.min = min
	return count
}

// GetMax returns a float64
func (count count) GetMax() int {
	return count.max
}

// SetMax sets *count.count and returns *count
func (count *count) SetMax(max int) *count {
	count.max = max
	return count
}

func (count count) Validate(value interface{}) error {
	if f, err := ToInterfaceArray(value); err != nil {
		return err
	} else if count.min == count.max && len(f) != count.min {
		return fmt.Errorf(CountExactMessage, fmt.Sprint(count.min))
	} else if len(f) < count.min {
		return fmt.Errorf(CountMinMessage, fmt.Sprint(count.min))
	} else if count.max < len(f) {
		return fmt.Errorf(CountMaxMessage, fmt.Sprint(count.max))
	}
	return nil
}

/***********/
/* HELPERS */
/***********/

// ToInterfaceArray takes an array or slice and returns an interface slice or
// an error if the value isn't an array or a slice
func ToInterfaceArray(value interface{}) ([]interface{}, error) {
	if v, ok := value.([]interface{}); ok {
		return v, nil
	}
	if IsArrayorSlice(value) {
		v := reflect.ValueOf(value)
		r := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			r[i] = v.Index(i).Interface()
		}
		return r, nil

	}
	return nil, fmt.Errorf("%+v is not an array or a slice", value)

}

// IsArrayorSlice returns true if value is
// an array or a slice
func IsArrayorSlice(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

// ToFloat64 converts a number to a float64 or returns an error
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
		return 0, fmt.Errorf("Cant convert %s to float64", fmt.Sprint(value))
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
	EmailMessage                   = "This value is not a valid email address"
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
	LessThanMessage                = "This value should be less than %s"
	LessThanOrEqualMessage         = "This value should be less than or equal to %s"
	GreaterThanMessage             = "This value should be greater than %s"
	GreaterThanOrEqualMessage      = "This value should be greater than or equal to %s"
	ChoiceMessage                  = "The value you selected is not a valid choice"
	ChoiceMinMessage               = "You must select at least %s choices"
	ChoiceMaxMessage               = "You must select at most %s choices"
	ChoiceMultipleMessage          = "One or more of the given values is invalid"
	CountMinMessage                = "This collection should contain %s elements or more"
	CountMaxMessage                = "This collection should contain %s elements or less"
	CountExactMessage              = "This collection should contain exactly %s elements"
)

var (
	// EmailRegexp represents an email pattern
	EmailRegexp = regexp.MustCompile(".+\\@.+\\..+")
)
