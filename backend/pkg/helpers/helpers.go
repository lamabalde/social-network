package helpers

import (
	"fmt"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHashToPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ShortMessage(message []byte) []byte {
	if len(message) > 100 {
		shortMessage := make([]byte, 100)
		copy(shortMessage, message[:97])
		shortMessage[97] = '.'
		shortMessage[98] = '.'
		shortMessage[99] = '.'
		return append(shortMessage, []byte("...")...)
	}
	return message
}

func FindFile(targetDir string, pattern string) string {

	matches, err := filepath.Glob(targetDir + pattern)
	if err != nil {
		fmt.Println(err)
	}
	if len(matches) != 0 {
		return matches[0]
	}
	return ""
}
