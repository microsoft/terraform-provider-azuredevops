//go:build all || utils || secretmemo
// +build all utils secretmemo

package secretmemo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNewHappyPath(t *testing.T) {
	result, memo, err := IsUpdating("mysecret", "")
	require.True(t, result)
	require.NotEmpty(t, memo)
	require.Nil(t, err)
}

func TestIsUpdatingHappyPath(t *testing.T) {
	firstResult, firstMemo, err1 := IsUpdating("mysecret", "")
	secondResult, secondMemo, err2 := IsUpdating("mychange", firstMemo)
	require.True(t, firstResult)
	require.True(t, secondResult)
	require.NotEqual(t, firstMemo, secondMemo)
	require.Nil(t, err1)
	require.Nil(t, err2)
}

func TestIsSameValueAsBeforeHappyPath(t *testing.T) {
	firstResult, firstMemo, err1 := IsUpdating("mysecret", "")
	secondResult, secondMemo, err2 := IsUpdating("mysecret", firstMemo)
	require.True(t, firstResult)
	require.False(t, secondResult)
	require.EqualValues(t, firstMemo, secondMemo)
	require.Nil(t, err1)
	require.Nil(t, err2)
}

func TestIsRottenMemo(t *testing.T) {
	result, memo, err := IsUpdating("mysecret", "!@#$")
	require.True(t, result)
	require.NotEmpty(t, memo)
	require.Nil(t, err)
}

func TestIsMissingSecret(t *testing.T) {
	result, memo, err := IsUpdating("", "anything")
	require.False(t, result)
	require.Equal(t, "anything", memo)
	require.Nil(t, err)
}

func TestIsValidMemo(t *testing.T) {
	require.False(t, isValidMemo("foo"))
	require.True(t, isValidMemo("$2a$"))
	require.True(t, isValidMemo("$2b$"))
	require.True(t, isValidMemo("$2y$"))
}

func isValidMemo(memo string) bool {
	validHashPrefixes := [4]string{"$2a$", "$2b$", "$2y$", "$arg"}
	for _, s := range validHashPrefixes {
		if strings.HasPrefix(memo, s) {
			return true
		}
	}
	return false
}
