package exception

import "fmt"

var _ error = (*InvalidInputError)(nil)

var _ ValidateMessenger = (*InvalidInputError)(nil)

type ValidateMessenger interface {
	ShouldNotEmpty() error
	ShouldEqualTo(value interface{}) error
	ShouldGreaterThan(value interface{}) error
	ShouldLessThan(value interface{}) error
	ShouldBetweenRange(min, max interface{}) error
	ShouldInRange(min, max interface{}) error
}

type InvalidInputError struct {
	KeyName string
	Message string
}

func NewInvalidInputError(fieldName string) *InvalidInputError {
	return &InvalidInputError{KeyName: fieldName}
}

func (e InvalidInputError) Error() string {
	return fmt.Sprintf("input `%s` is invalid: %s", e.KeyName, e.Message)
}
func (e InvalidInputError) ShouldNotEmpty() error {
	e.Message = "should not empty"
	return InvalidInputError{
		KeyName: e.KeyName,
		Message: e.Message,
	}
}

func (e InvalidInputError) ShouldEqualTo(value interface{}) error {
	e.Message = fmt.Sprintf("should equal to %v", value)
	return InvalidInputError{
		KeyName: e.KeyName,
		Message: e.Message,
	}
}

func (e InvalidInputError) ShouldGreaterThan(value interface{}) error {
	e.Message = fmt.Sprintf("should greater than %v", value)
	return InvalidInputError{
		KeyName: e.KeyName,
		Message: e.Message,
	}
}

func (e InvalidInputError) ShouldLessThan(value interface{}) error {
	e.Message = fmt.Sprintf("should less than %v", value)
	return InvalidInputError{
		KeyName: e.KeyName,
		Message: e.Message,
	}
}

func (e InvalidInputError) ShouldBetweenRange(min, max interface{}) error {
	e.Message = fmt.Sprintf("should between range %v and %v (min < x < max)", min, max)
	return InvalidInputError{
		KeyName: e.KeyName,
		Message: e.Message,
	}
}

func (e InvalidInputError) ShouldInRange(min, max interface{}) error {
	e.Message = fmt.Sprintf("should in range %v and %v (min <= in <= max)", min, max)
	return InvalidInputError{
		KeyName: e.KeyName,
		Message: e.Message,
	}
}
