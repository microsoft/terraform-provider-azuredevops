package secretmemo

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const isUpdating = true
const isNotUpdating = false
const isErr = false

func isBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isValidMemo(memo string) bool {
	validBcryptHashPrefixes := [3]string{"$2a$", "$2b$", "$2y$"}
	for _, s := range validBcryptHashPrefixes {
		if strings.HasPrefix(memo, s) {
			return true
		}
	}
	return false
}

func calcMementoForSecret(secret, memento string) (string, error) {
	secretAsBytes := []byte(secret)
	hash, err := bcrypt.GenerateFromPassword(secretAsBytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func doesMemoMatchSecret(secret, memento string) bool {
	if isBlank(memento) {
		return false
	}
	secretAsBytes := []byte(secret)
	mementoAsBytes := []byte(memento)
	err := bcrypt.CompareHashAndPassword(mementoAsBytes, secretAsBytes)
	if err != nil {
		// ignoring the err
		return false
	}
	return true
}

// IsUpdating is used to determine if the secret getting updated?
func IsUpdating(secret, oldMemo string) (bool, string, error) {
	if isBlank(secret) {
		return isNotUpdating, oldMemo, nil
	}

	if doesMemoMatchSecret(secret, oldMemo) {
		return isNotUpdating, oldMemo, nil
	}

	newMemo, err := calcMementoForSecret(secret, oldMemo)
	if err != nil {
		return isErr, "", err
	}

	return isUpdating, newMemo, nil

}
