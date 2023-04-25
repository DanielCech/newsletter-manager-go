package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewEmail(t *testing.T) {
	rawEmail := "jozko.dlouhy@gmail.com"
	email, err := NewEmail(rawEmail)
	assert.NoError(t, err)
	assert.EqualValues(t, rawEmail, email)

	rawEmail = "invalid@email"
	email, err = NewEmail(rawEmail)
	assert.Error(t, err)
	assert.Empty(t, email)
}

func Test_Email_UnmarshalText(t *testing.T) {
	rawEmail := "jozko.DLOUHY@gmail.com"
	expected := "jozko.dlouhy@gmail.com"
	email := Email("")
	err := email.UnmarshalText([]byte(rawEmail))
	assert.NoError(t, err)
	assert.EqualValues(t, expected, email)
}

func Test_Email_UnmarshalJSON(t *testing.T) {
	rawEmail := "jozko.DLOUHY@gmail.com"
	expected := "jozko.dlouhy@gmail.com"
	email := Email("")
	err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, rawEmail)), &email)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, email)
}

func Test_Email_String(t *testing.T) {
	rawEmail := "jozko.DLOUHY@gmail.com"
	expected := "jozko.dlouhy@gmail.com"
	email := Email(rawEmail)
	assert.Equal(t, expected, email.String())
}

func Test_Email_Valid(t *testing.T) {
	email := Email("jozko.dlouhy@gmail.com")
	assert.NoError(t, email.Valid())

	email = "invalid@email"
	assert.Error(t, email.Valid())

	email = "jozko.DLOUHY@gmail.com"
	assert.EqualError(t, email.Valid(), "email is not lowercase")
}
