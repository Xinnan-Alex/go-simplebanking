package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrtsuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Random Integer generated between min and max
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

func RandomAmount() int64 {
	return RandomInt(-1000, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD, SGD, MYR}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}
