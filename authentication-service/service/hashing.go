package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type passwordHasher struct {
	config argonParams
}

func NewPasswordHasher() *passwordHasher {
	timeCost := 1
	memoryCost := 64 * 1024
	parallelism := 4
	saltLength := 16
	hashLength := 32
	argonConfig := argonParams{
		memory:      uint32(memoryCost),
		iterations:  uint32(timeCost),
		parallelism: uint8(parallelism),
		saltLength:  uint32(saltLength),
		keyLength:   uint32(hashLength),
	}
	return &passwordHasher{argonConfig}
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func (p *passwordHasher) CreateEncodedHash(password string) (encodedPassword string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.config.saltLength)
	if err != nil {
		return "", err
	}
	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, p.config.iterations, p.config.memory, p.config.parallelism, p.config.keyLength)
	// Base64 encode the salt and hashed password.
	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)
	// Concatenate salt and hashed password for storage
	encodedPassword = fmt.Sprintf("%s$%s", b64Salt, b64Hash)

	return encodedPassword, nil
}

// VerifyPassword checks if the provided password matches the hashed password.
func (p *passwordHasher) VerifyPassword(hashedPassword, password string) (bool, error) {
	// Split the stored value into salt and hashed password
	parts := splitHashedPassword(hashedPassword)
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hashed password format")
	}

	// Decode the salt and hashed password from base64
	salt, err := decodeBase64(parts[0])
	if err != nil {
		return false, err
	}

	storedHashedPassword, err := decodeBase64(parts[1])
	if err != nil {
		return false, err
	}
	// Hash the provided password with the stored salt
	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	// Compare the computed hash with the stored hash
	return compareHashes(computedHash, storedHashedPassword), nil
}

// Helper function to decode base64 string to bytes
func decodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// Helper function to split the stored value into salt and hashed password parts
func splitHashedPassword(encodedPassword string) []string {
	return strings.Split(encodedPassword, "$")
}

// Helper function to compare two byte slices for equality
func compareHashes(hash1, hash2 []byte) bool {
	return subtle.ConstantTimeCompare(hash1, hash2) == 1
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
