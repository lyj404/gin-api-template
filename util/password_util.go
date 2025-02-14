package util

import (
	"gin-api-template/config"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	// 将密码加盐
	saltedPassword := config.CfgPassword.SaltPrefix + password + config.CfgPassword.SaltSuffix
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), config.CfgPassword.Cost)
	return string(hash), err
}

// ComparePassword 比较密码
func ComparePassword(hashedPassword, password string) error {
	saltedPassword := config.CfgPassword.SaltPrefix + password + config.CfgPassword.SaltSuffix
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
}
