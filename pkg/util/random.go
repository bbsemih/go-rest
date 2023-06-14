package util

import (
	"fmt"
	"math/rand"
	"strings"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// Seed the random number generator
	rand.Seed(rand.Int63())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "TL"}
	l := len(currencies)

	return currencies[rand.Intn(l)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
