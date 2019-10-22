package secretmemo

import (
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
	firstResult, firstMemo, err := IsUpdating("mysecret", "")
	secondResult, secondMemo, err := IsUpdating("mychange", firstMemo)
	require.True(t, firstResult)
	require.True(t, secondResult)
	require.NotEqual(t, firstMemo, secondMemo)
	require.Nil(t, err)
}

func TestIsSameValueAsBeforeHappyPath(t *testing.T) {
	firstResult, firstMemo, err := IsUpdating("mysecret", "")
	secondResult, secondMemo, err := IsUpdating("mysecret", firstMemo)
	require.True(t, firstResult)
	require.False(t, secondResult)
	require.EqualValues(t, firstMemo, secondMemo)
	require.Nil(t, err)
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
