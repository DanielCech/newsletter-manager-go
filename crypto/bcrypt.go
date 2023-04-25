package crypto

import (
	"crypto/hmac"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

const (
	bcryptMaxPasswordLength = 72
)

// BcryptHasher is capable of creating and comparing hashes.
type BcryptHasher struct {
	cost   int
	pepper []byte
}

// NewBcryptHasher creates instance of BcryptHasher.
func NewBcryptHasher(cost int, pepper []byte) BcryptHasher {
	return BcryptHasher{
		cost:   cost,
		pepper: pepper,
	}
}

// NewDefaultBcryptHasher creates instance of BcryptHasher with default cost.
func NewDefaultBcryptHasher(pepper []byte) BcryptHasher {
	return NewBcryptHasher(bcrypt.DefaultCost, pepper)
}

// HashPassword creates a hash of the password.
func (b BcryptHasher) HashPassword(password []byte) ([]byte, error) {
	passwordHash, err := b.encodedSHA512(password)
	if err != nil {
		return nil, err
	}
	return bcrypt.GenerateFromPassword(passwordHash[:bcryptMaxPasswordLength], b.cost)
}

// CompareHashAndPassword compares password with the given hash.
func (b BcryptHasher) CompareHashAndPassword(hash, password []byte) bool {
	passwordHash, err := b.encodedSHA512(password)
	if err != nil {
		return false
	}
	if err = bcrypt.CompareHashAndPassword(hash, passwordHash[:bcryptMaxPasswordLength]); err != nil {
		return false
	}
	return true
}

func (b BcryptHasher) encodedSHA512(value []byte) ([]byte, error) {
	shaHasher := hmac.New(sha3.New512, b.pepper)
	if _, err := shaHasher.Write(value); err != nil {
		return nil, err
	}
	hash := shaHasher.Sum(nil)
	return []byte(base64.StdEncoding.EncodeToString(hash)), nil
}
