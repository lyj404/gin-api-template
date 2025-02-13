package util

import "golang.org/x/crypto/bcrypt"

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// ComparePassword 比较密码
func ComparePassword(hashedPassword, passwrod string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwrod))
}
