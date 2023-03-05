package randutil

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var _ = seedRand()

func seedRand() error {
	rand.Seed(time.Now().UnixNano())
	return nil
}

func IntWithRange(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func StringWithQuantity(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func Owner() string {
	return StringWithQuantity(6)
}

func Money() int64 {
	return IntWithRange(0, 1000)
}

func Currency() string {
	currencies := []string{"USD", "VND", "EUR"}

	n := len(currencies)
	return currencies[rand.Intn(n)]
}
