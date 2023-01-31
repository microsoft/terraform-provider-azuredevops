package secretmemo

import (
	"strings"

	"github.com/alexedwards/argon2id"
)

const isUpdating = true
const isNotUpdating = false
const isErr = false

func isBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func calcMementoForSecret(secret, memento string) (string, error) {
	hash, err := argon2id.CreateHash(secret, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func doesMemoMatchSecret(secret, memento string) bool {
	if isBlank(memento) {
		return false
	}
	match, err := argon2id.ComparePasswordAndHash(secret, memento)
	if err != nil {
		return false
	}
	return match
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