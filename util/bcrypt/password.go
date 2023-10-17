package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// 비밀번호 암호화 함수
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errHashPassword(err)
	}

	return string(hashedPassword), nil
}

// 비밀번호 검증 함수
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
