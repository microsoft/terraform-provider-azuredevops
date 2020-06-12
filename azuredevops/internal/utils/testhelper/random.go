package testhelper

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandIntSlice returns a random range of numbers as a slice within a range
//
// Usage:
// arr := random.RandIntSlice(2, 100, 3)
//
// Output: [40, 5, 81]
// Output: [12, 52, 31]
func RandIntSlice(min int, max int, n int) []int {
	arr := make([]int, n)
	for r := 0; r <= n-1; r++ {
		arr[r] = RandInt(min, max)
	}
	return arr
}

// RandInt return a random number within a range
func RandInt(min int, max int) int {
	return rand.Intn(max) + min
}
