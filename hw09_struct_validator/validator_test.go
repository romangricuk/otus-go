package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Profile struct {
		User User `validate:"nested"`
		Bio  string
	}

	ValidatorInvalidMin struct {
		Value int `validate:"min:aaa"`
	}
	ValidatorInvalidMax struct {
		Value int `validate:"max:aaa"`
	}
	ValidatorInvalidLen struct {
		Value string `validate:"len:aaa"`
	}
	UnknownValidator struct {
		Value string `validate:"aaa:aaa"`
	}
	ValidatorInvalidTagFormat struct {
		Value int `validate:"min"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: &ValidatorInvalidMin{
				Value: 123,
			},
			expectedErr: InvalidValidatorError{
				Field: "Value",
				Err: &strconv.NumError{
					Func: "Atoi",
					Num:  "aaa",
					Err:  strconv.ErrSyntax,
				},
			},
		},
		{
			in: &ValidatorInvalidMax{
				Value: 123,
			},
			expectedErr: InvalidValidatorError{
				Field: "Value",
				Err: &strconv.NumError{
					Func: "Atoi",
					Num:  "aaa",
					Err:  strconv.ErrSyntax,
				},
			},
		},
		{
			in: &ValidatorInvalidLen{
				Value: "123",
			},
			expectedErr: InvalidValidatorError{
				Field: "Value",
				Err: &strconv.NumError{
					Func: "Atoi",
					Num:  "aaa",
					Err:  strconv.ErrSyntax,
				},
			},
		},
		{
			in: &UnknownValidator{
				Value: "123",
			},
			expectedErr: InvalidValidatorError{
				Field: "Value",
				Err:   fmt.Errorf("unknown validator aaa"),
			},
		},
		{
			in: &ValidatorInvalidTagFormat{
				Value: 123,
			},
			expectedErr: InvalidValidatorError{
				Field: "Value",
				Err:   InvalidTagFormatError,
			},
		},

		{
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John Doe",
				Age:    25,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in: &User{
				ID:     "12345678-1234-1234-1234-1234567890123",
				Name:   "John Doe",
				Age:    17,
				Email:  "johndoeexample.com",
				Role:   "user",
				Phones: []string{"1234567890"},
			},
			expectedErr: ValidationErrors{
				ValidationError{"ID", fmt.Errorf("length must be 36")},
				ValidationError{"Age", fmt.Errorf("must be at least 18")},
				ValidationError{"Email", fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
				ValidationError{"Role", fmt.Errorf("must be one of admin, stuff")},
				ValidationError{"Phones", fmt.Errorf("length must be 11")},
			},
		},
		{
			in: &App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: &Response{
				Code: 403,
				Body: "Forbidden",
			},
			expectedErr: ValidationErrors{
				ValidationError{"Code", fmt.Errorf("must be one of 200, 404, 500")},
			},
		},
		{
			in: &Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: &Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: &Profile{
				User: User{
					ID:     "12345678-1234-1234-1234-123456789012",
					Name:   "John Doe",
					Age:    25,
					Email:  "johndoe@example.com",
					Role:   "admin",
					Phones: []string{"12345678901"},
				},
				Bio: "This is a bio.",
			},
			expectedErr: nil,
		},
		{
			in: &Profile{
				User: User{
					ID:     "12345678-1234-1234-1234-1234567890123",
					Name:   "John Doe",
					Age:    55,
					Email:  "john.doeexample.com",
					Role:   "user",
					Phones: []string{"1234567890"},
				},
				Bio: "This is a bio.",
			},
			expectedErr: ValidationErrors{
				ValidationError{"User.ID", fmt.Errorf("length must be 36")},
				ValidationError{"User.Age", fmt.Errorf("must be at most 50")},
				ValidationError{"User.Email", fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
				ValidationError{"User.Role", fmt.Errorf("must be one of admin, stuff")},
				ValidationError{"User.Phones", fmt.Errorf("length must be 11")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			//t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
