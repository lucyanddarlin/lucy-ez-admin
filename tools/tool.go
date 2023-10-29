package tools

import "golang.org/x/crypto/bcrypt"

func HashPwd(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
