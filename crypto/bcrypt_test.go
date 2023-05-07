package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_BcryptHasher_Compatibility ensures compatibility between BcryptHasher.HashPassword and BcryptHasher.CompareHashAndPassword.
func Test_BcryptHasher_Compatibility(t *testing.T) {
	tests := []struct {
		name        string
		password    []byte
		pepper      []byte
		expectedErr error
	}{
		{
			name:        "success",
			password:    []byte("Topsecret1"),
			pepper:      []byte("57OEolMBzlyL1SMO40opczY5WEuk9BXzwI2p6K3qfWfnXUd3g699IlmvVW8Tpnpk"),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewDefaultBcryptHasher(tt.pepper)
			passwordHash, err := hasher.HashPassword(tt.password)
			assert.Equal(t, tt.expectedErr, err)
			assert.True(t, hasher.CompareHashAndPassword(passwordHash, tt.password))
		})
	}
}
