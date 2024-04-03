package tcpio

import (
	"math/rand"
	"time"
)

func uuidGen() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}

	return string(b)
}
