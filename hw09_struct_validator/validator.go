package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidTagFormat = fmt.Errorf("invalid validation tag format")

// ValidationError represents a validation error for a specific field.
type ValidationError struct {
	Field string
	Err   error
}

// Err implements the error interface for ValidationError.
func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Err)
}

type InvalidValidatorError struct {
	Field string
	Err   error
}

// Err implements the error interface for InvalidValidatorError.
func (e InvalidValidatorError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Err)
}

// ValidationErrors is a slice of ValidationError.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	errs := make([]string, 0, len(ve))
	for _, err := range ve {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, ", ")
}

// Validate validates the fields of a struct based on the "validate" tag.
//
//nolint:gocognit
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	var validationErrors ValidationErrors
	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		validateTag := fieldType.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		if fieldType.Type.Kind() == reflect.Struct && validateTag == "nested" {
			if err := Validate(field.Addr().Interface()); err != nil {
				if ve, ok := err.(ValidationErrors); ok { //nolint:errorlint
					for i, e := range ve {
						ve[i].Field = fieldType.Name + "." + e.Field
					}
					validationErrors = append(validationErrors, ve...)
				} else {
					return InvalidValidatorError{
						Field: fieldType.Name,
						Err:   err,
					}
				}
			}
			continue
		}

		validations := strings.Split(validateTag, "|")
		for _, validation := range validations {
			validationParts := strings.SplitN(validation, ":", 2)
			if len(validationParts) != 2 {
				return InvalidValidatorError{
					Field: fieldType.Name,
					Err:   ErrInvalidTagFormat,
				}
			}
			validator := validationParts[0]
			param := validationParts[1]

			if err := validateField(field, validator, param); err != nil {
				if ve, ok := err.(*ValidationError); ok { //nolint:errorlint
					ve.Field = fieldType.Name
					validationErrors = append(validationErrors, *ve)
				} else {
					return InvalidValidatorError{
						Field: fieldType.Name,
						Err:   err,
					}
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(field reflect.Value, validator string, param string) error {
	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		return validateString(field.String(), validator, param)
	case reflect.Int:
		return validateInt(int(field.Int()), validator, param)
	case reflect.Slice:
		return validateSlice(field, validator, param)
	default:
		return nil
	}
}

func validateString(value string, validator string, param string) error {
	switch validator {
	case "len":
		length, err := strconv.Atoi(param)
		if err != nil {
			return err
		}
		if len(value) != length {
			return &ValidationError{
				Err: fmt.Errorf("length must be %d", length),
			}
		}
	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return err
		}
		if !re.MatchString(value) {
			return &ValidationError{
				Err: fmt.Errorf("must match regexp %s", param),
			}
		}
	case "in":
		options := strings.Split(param, ",")
		for _, option := range options {
			if value == option {
				return nil
			}
		}
		return &ValidationError{
			Err: fmt.Errorf("must be one of %s", strings.Join(options, ", ")),
		}
	default:
		return fmt.Errorf("unknown validator %s", validator)
	}
	return nil
}

func validateInt(value int, validator string, param string) error {
	switch validator {
	case "min":
		min, err := strconv.Atoi(param)
		if err != nil {
			return err
		}
		if value < min {
			return &ValidationError{
				Err: fmt.Errorf("must be at least %d", min),
			}
		}
	case "max":
		max, err := strconv.Atoi(param)
		if err != nil {
			return err
		}
		if value > max {
			return &ValidationError{
				Err: fmt.Errorf("must be at most %d", max),
			}
		}
	case "in":
		options := strings.Split(param, ",")
		for _, option := range options {
			opt, err := strconv.Atoi(option)
			if err != nil {
				return err
			}
			if value == opt {
				return nil
			}
		}
		return &ValidationError{
			Err: fmt.Errorf("must be one of %s", strings.Join(options, ", ")),
		}
	default:
		return fmt.Errorf("unknown validator %s", validator)
	}
	return nil
}

func validateSlice(field reflect.Value, validator string, param string) error {
	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i)
		if err := validateField(elem, validator, param); err != nil {
			return err
		}
	}
	return nil
}
