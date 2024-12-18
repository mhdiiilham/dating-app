package credential

import "golang.org/x/crypto/bcrypt"

// Hasher struct is an instance for hashing plain password.
type Hasher struct {
}

// HashPassword function accept plain password and hash it using bcrypt.
func (h Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// ComparePassword function check plain and hashed password.
func (h Hasher) ComparePassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
