package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPassword(t *testing.T) {
	rawPassword := "Topsecret1"
	password, err := NewPassword(rawPassword)
	assert.NoError(t, err)
	assert.EqualValues(t, rawPassword, password)

	password, err = NewPassword("")
	assert.Error(t, err)
	assert.Empty(t, password)
}

func Test_Password_MarshalText(t *testing.T) {
	password := Password("Topsecret1")
	result, err := password.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, result, []byte(sensitiveOutput))
}

func Test_Password_MarshalJSON(t *testing.T) {
	password := Password("Topsecret1")
	result, err := json.Marshal(password)
	assert.NoError(t, err)
	assert.Equal(t, result, []byte(fmt.Sprintf(`"%s"`, sensitiveOutput)))
}

func Test_Password_UnmarshalText(t *testing.T) {
	var result Password

	input := []byte("Invalid")
	err := result.UnmarshalText(input)
	assert.ErrorIs(t, err, ErrInvalidPasswordLength)
	assert.ErrorIs(t, err, ErrMissingDigitCharacter)

	input = []byte(`Top"'\Secret1`)
	err = result.UnmarshalText(input)
	assert.NoError(t, err)
	assert.Equal(t, Password(`Top"'\Secret1`), result)
}

// Password does not implement encoding.Unmarshaler, as json package uses UnmarshalText if UnmarshalJSON is missing.
func Test_Password_UnmarshalJSON(t *testing.T) {
	var result Password

	input := []byte(`"Invalid"`)
	err := json.Unmarshal(input, &result)
	assert.ErrorIs(t, err, ErrInvalidPasswordLength)
	assert.ErrorIs(t, err, ErrMissingDigitCharacter)

	input = []byte(`"Top\"'\\Secret1"`)
	err = json.Unmarshal(input, &result)
	assert.NoError(t, err)
	assert.Equal(t, Password(`Top"'\Secret1`), result)
}

func Test_Password_String(t *testing.T) {
	password := Password("Topsecret1")
	result := password.String()
	assert.Equal(t, result, sensitiveOutput)
}

func Test_Password_Valid(t *testing.T) {
	tests := []struct {
		name            string
		password        Password
		expectedErr     error
		alsoExpectedErr error
	}{
		{
			name:        "success",
			password:    "Topsecret1",
			expectedErr: nil,
		},
		{
			name:        "failure:minimum-password-length",
			password:    "abc",
			expectedErr: ErrInvalidPasswordLength,
		},
		{
			name: "failure:maximum-password-length",
			// 256 chars
			password:    "ThisIsASuperLongAndVerySecurePasswordThatContainsManyDifferentWordsAndCharactersThatAreNotEasyToGuessOrCrack1234!ThisIsASuperLongAndVerySecurePasswordThatContainsManyDifferentWordsAndCharactersThatAreNotEasyToGuessOrCrack1234!ThisIsASuperLongAndVerySecureP",
			expectedErr: ErrInvalidPasswordLength,
		},
		{
			name:        "failure:missing-upper-case-character",
			password:    "topsecret1",
			expectedErr: ErrMissingUpperCaseCharacter,
		},
		{
			name:        "failure:missing-lower-case-character",
			password:    "TOPSECRET1",
			expectedErr: ErrMissingLowerCaseCharacter,
		},
		{
			name:        "failure:missing-digit-character",
			password:    "Topsecret,",
			expectedErr: ErrMissingDigitCharacter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.password.Valid()
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func Test_Password_Valid_MultiError(t *testing.T) {
	password := Password("TopSec")

	err := password.Valid()

	assert.ErrorIs(t, err, ErrInvalidPasswordLength)
	assert.ErrorIs(t, err, ErrMissingDigitCharacter)
	assert.NotErrorIs(t, err, ErrMissingUpperCaseCharacter)
	assert.NotErrorIs(t, err, ErrMissingLowerCaseCharacter)
}
